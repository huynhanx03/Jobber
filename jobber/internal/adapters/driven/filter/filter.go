package filter

import (
	"strings"

	"jobber/internal/constant"
	"jobber/internal/core/entity"
)

type JobFilter struct{}

func NewJobFilter() *JobFilter {
	return &JobFilter{}
}

func (f *JobFilter) Apply(jobs []entity.Job) []entity.Job {
	var result []entity.Job
	for _, job := range jobs {
		text := strings.ToLower(job.Title + " " + job.Description)
		if constant.KeywordRegex.MatchString(text) {
			result = append(result, job)
		}
	}
	return result
}
