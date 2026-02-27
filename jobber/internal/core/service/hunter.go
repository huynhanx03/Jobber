package service

import (
	"context"
	"log"
	"sync"
	"time"

	"jobber/internal/core/entity"
	"jobber/internal/core/ports"
)

type HunterService struct {
	scrapers  []ports.Scraper
	filter    ports.Filter
	storage   ports.Storage
	notifiers []ports.Notifier
}

func NewHunterService(
	scrapers []ports.Scraper,
	filter ports.Filter,
	storage ports.Storage,
	notifiers []ports.Notifier,
) *HunterService {
	return &HunterService{
		scrapers:  scrapers,
		filter:    filter,
		storage:   storage,
		notifiers: notifiers,
	}
}

func (h *HunterService) Run(ctx context.Context) error {
	log.Println("🚀 Starting job hunt...")
	start := time.Now()

	var allJobs []entity.Job
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, s := range h.scrapers {
		wg.Add(1)
		go func(scraper ports.Scraper) {
			defer wg.Done()
			log.Printf("⏳ Scraping %s...", scraper.Name())
			jobs, err := scraper.Scrape(ctx)
			if err != nil {
				log.Printf("⚠️  %s scrape failed: %v", scraper.Name(), err)
				return
			}
			log.Printf("✅ %s: found %d jobs", scraper.Name(), len(jobs))

			mu.Lock()
			allJobs = append(allJobs, jobs...)
			mu.Unlock()
		}(s)
	}
	wg.Wait()
	log.Printf("📦 Total raw jobs: %d", len(allJobs))

	filtered := h.filter.Apply(allJobs)
	log.Printf("🧹 After keyword filter: %d jobs", len(filtered))

	seenJobs, err := h.storage.LoadSeenJobs()
	if err != nil {
		log.Printf("⚠️  Failed to load seen jobs: %v", err)
		seenJobs = make(map[string]bool)
	}

	var newJobs []entity.Job
	for _, job := range filtered {
		if !seenJobs[job.URL] {
			newJobs = append(newJobs, job)
		}
	}
	log.Printf("🔍 New jobs: %d", len(newJobs))

	if len(newJobs) == 0 {
		log.Println("ℹ️  No new jobs to send.")
		return nil
	}

	for _, n := range h.notifiers {
		if err := n.Send(ctx, newJobs); err != nil {
			log.Printf("⚠️  Notifier failed: %v", err)
		}
	}

	var sentURLs []string
	for _, job := range newJobs {
		sentURLs = append(sentURLs, job.URL)
	}
	if err := h.storage.SaveSeenJobs(sentURLs); err != nil {
		log.Printf("⚠️  Failed to save seen jobs: %v", err)
	}

	log.Printf("🏁 Done in %s. Sent %d jobs.", time.Since(start).Round(time.Millisecond), len(newJobs))
	return nil
}
