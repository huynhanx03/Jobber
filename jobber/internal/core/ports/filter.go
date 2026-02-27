package ports

import "jobber/internal/core/entity"

type Filter interface {
	Apply(jobs []entity.Job) []entity.Job
}
