package ports

import (
	"context"

	"jobber/internal/core/entity"
)

type Scraper interface {
	Name() string
	Scrape(ctx context.Context) ([]entity.Job, error)
}
