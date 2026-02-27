package di

import (
	"jobber/config"
	"jobber/internal/adapters/driven/filter"
	"jobber/internal/adapters/driven/notifier"
	"jobber/internal/adapters/driven/scraper"
	"jobber/internal/adapters/driven/storage"
	"jobber/internal/core/ports"
	"jobber/internal/core/service"
	"jobber/internal/infrastructure/browser"
)

func Wire(cfg *config.Config) (*service.HunterService, *browser.Browser) {
	b := browser.New()

	scrapers := []ports.Scraper{
		scraper.NewIndeedScraper(b, cfg.Keywords, cfg.Locations),
		scraper.NewITViecScraper(b, cfg.Keywords),
		scraper.NewTopDevScraper(b, cfg.Keywords),
	}

	jobFilter := filter.NewJobFilter()
	fileStorage := storage.NewFileStorage(cfg.SeenJobsPath)
	discord := notifier.NewDiscordNotifier(cfg.DiscordWebhookURL)

	notifiers := []ports.Notifier{discord}

	svc := service.NewHunterService(scrapers, jobFilter, fileStorage, notifiers)
	return svc, b
}
