package ports

import (
	"context"

	"jobber/internal/core/entity"
)

type Notifier interface {
	Send(ctx context.Context, jobs []entity.Job) error
}
