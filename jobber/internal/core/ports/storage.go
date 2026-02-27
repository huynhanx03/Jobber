package ports

type Storage interface {
	LoadSeenJobs() (map[string]bool, error)
	SaveSeenJobs(urls []string) error
}
