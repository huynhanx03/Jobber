package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"jobber/internal/constant"
	"jobber/internal/core/entity"
	"jobber/internal/infrastructure/browser"
)

type IndeedScraper struct {
	browser   *browser.Browser
	keywords  []string
	locations []string
}

func NewIndeedScraper(b *browser.Browser, keywords, locations []string) *IndeedScraper {
	return &IndeedScraper{
		browser:   b,
		keywords:  keywords,
		locations: locations,
	}
}

func (s *IndeedScraper) Name() string {
	return constant.SourceIndeed
}

func (s *IndeedScraper) Scrape(_ context.Context) ([]entity.Job, error) {
	var jobs []entity.Job
	seen := make(map[string]bool)

	for _, keyword := range s.keywords {
		for _, loc := range s.locations {
			searchURL := fmt.Sprintf("https://vn.indeed.com/jobs?q=%s&l=%s&sort=date", keyword, loc)
			log.Printf("  🔍 Indeed: keyword=%s location=%s", keyword, loc)

			scraped, err := s.scrapePage(searchURL)
			if err != nil {
				log.Printf("  ⚠️  Indeed page error: %v", err)
				continue
			}

			for _, job := range scraped {
				key := strings.ToLower(job.Title + "|" + job.Company)
				if !seen[key] {
					seen[key] = true
					jobs = append(jobs, job)
				}
			}
		}
	}
	return jobs, nil
}

func (s *IndeedScraper) scrapePage(url string) ([]entity.Job, error) {
	page := s.browser.StealthPage(url)
	defer page.MustClose()

	browser.RandomDelay(1*time.Second, 2*time.Second)

	titles := page.MustEval(`() => Array.from(document.querySelectorAll('h2.jobTitle span[title]')).map(e => e.getAttribute('title'))`).Arr()
	companies := page.MustEval(`() => Array.from(document.querySelectorAll('[data-testid="company-name"]')).map(e => e.textContent)`).Arr()
	locations := page.MustEval(`() => Array.from(document.querySelectorAll('[data-testid="text-location"]')).map(e => e.textContent)`).Arr()
	urls := page.MustEval(`() => Array.from(document.querySelectorAll('h2.jobTitle a')).map(e => e.href)`).Arr()
	salaries := page.MustEval(`() => Array.from(document.querySelectorAll('.salary-snippet-container')).map(e => e.textContent)`).Arr()

	count := minOf(len(titles), len(companies), len(locations), len(urls))
	jobs := make([]entity.Job, count)

	for i := 0; i < count; i++ {
		salary := ""
		if i < len(salaries) {
			salary = strings.TrimSpace(salaries[i].Str())
		}

		jobs[i] = entity.Job{
			Title:     strings.TrimSpace(titles[i].Str()),
			Company:   strings.TrimSpace(companies[i].Str()),
			URL:       urls[i].Str(),
			Location:  strings.TrimSpace(locations[i].Str()),
			Salary:    salary,
			Source:    constant.SourceIndeed,
			ScrapedAt: time.Now(),
		}
	}

	return jobs, nil
}

func minOf(vals ...int) int {
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
