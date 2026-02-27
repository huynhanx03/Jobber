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

type TopDevScraper struct {
	browser  *browser.Browser
	keywords []string
}

func NewTopDevScraper(b *browser.Browser, keywords []string) *TopDevScraper {
	return &TopDevScraper{
		browser:  b,
		keywords: keywords,
	}
}

func (s *TopDevScraper) Name() string {
	return constant.SourceTopDev
}

func (s *TopDevScraper) Scrape(_ context.Context) ([]entity.Job, error) {
	var jobs []entity.Job

	for _, keyword := range s.keywords {
		searchURL := fmt.Sprintf("https://topdev.vn/jobs/search?keyword=%s&page=1", keyword)
		log.Printf("  🔍 TopDev: keyword=%s", keyword)

		scraped, err := s.scrapePage(searchURL)
		if err != nil {
			log.Printf("  ⚠️  TopDev page error: %v", err)
			continue
		}
		jobs = append(jobs, scraped...)
	}

	return jobs, nil
}

func (s *TopDevScraper) scrapePage(url string) ([]entity.Job, error) {
	page := s.browser.StealthPage(url)
	defer page.MustClose()

	browser.RandomDelay(2*time.Second, 3*time.Second)

	titles := page.MustEval(`() => Array.from(document.querySelectorAll('a[href*="/detail-jobs/"]')).map(e => e.textContent.trim())`).Arr()
	companies := page.MustEval(`() => Array.from(document.querySelectorAll('a[href*="/companies/"]')).map(e => e.textContent.trim())`).Arr()
	urls := page.MustEval(`() => Array.from(document.querySelectorAll('a[href*="/detail-jobs/"]')).map(e => e.href)`).Arr()

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
			Source:    constant.SourceTopDev,
			ScrapedAt: time.Now(),
		}
	}

	return jobs, nil
}
