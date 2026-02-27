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

type ITViecScraper struct {
	browser  *browser.Browser
	keywords []string
}

func NewITViecScraper(b *browser.Browser, keywords []string) *ITViecScraper {
	return &ITViecScraper{
		browser:  b,
		keywords: keywords,
	}
}

func (s *ITViecScraper) Name() string {
	return constant.SourceITViec
}

func (s *ITViecScraper) Scrape(_ context.Context) ([]entity.Job, error) {
	var jobs []entity.Job

	for _, keyword := range s.keywords {
		keywordSlug := strings.ToLower(strings.ReplaceAll(keyword, " ", "-"))
		searchURL := fmt.Sprintf("https://itviec.com/it-jobs/%s", keywordSlug)
		log.Printf("  🔍 ITViec: keyword=%s", keyword)

		scraped, err := s.scrapePage(searchURL)
		if err != nil {
			log.Printf("  ⚠️  ITViec page error: %v", err)
			continue
		}
		jobs = append(jobs, scraped...)
	}

	return jobs, nil
}

func (s *ITViecScraper) scrapePage(url string) ([]entity.Job, error) {
	page := s.browser.StealthPage(url)
	defer page.MustClose()

	browser.RandomDelay(3*time.Second, 5*time.Second)

	titles := page.MustEval(`() => Array.from(document.querySelectorAll('div.job-card')).map(c => c.querySelector('h3')?.textContent.trim() || "")`).Arr()
	urls := page.MustEval(`() => Array.from(document.querySelectorAll('div.job-card')).map(c => c.querySelector('h3')?.getAttribute('data-url') || "")`).Arr()
	companies := page.MustEval(`() => Array.from(document.querySelectorAll('div.job-card')).map(c => c.querySelector('a.text-rich-grey')?.textContent.trim() || "Unknown")`).Arr()

	count := minOf(len(titles), len(urls))
	jobs := make([]entity.Job, count)

	for i := 0; i < count; i++ {
		company := "Unknown"
		if i < len(companies) {
			company = strings.TrimSpace(companies[i].Str())
		}

		jobs[i] = entity.Job{
			Title:     strings.TrimSpace(titles[i].Str()),
			Company:   company,
			URL:       urls[i].Str(),
			Source:    constant.SourceITViec,
			ScrapedAt: time.Now(),
		}
	}

	return jobs, nil
}
