package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"jobber/internal/constant"
)

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func (f *FileStorage) LoadSeenJobs() (map[string]bool, error) {
	seen := make(map[string]bool)

	data, err := os.ReadFile(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			return seen, nil
		}
		return nil, err
	}

	var urls []string
	if err := json.Unmarshal(data, &urls); err != nil {
		return nil, err
	}

	for _, u := range urls {
		seen[u] = true
	}
	return seen, nil
}

func (f *FileStorage) SaveSeenJobs(urls []string) error {
	existing, err := f.LoadSeenJobs()
	if err != nil {
		existing = make(map[string]bool)
	}

	for _, u := range urls {
		existing[u] = true
	}

	var merged []string
	for u := range existing {
		merged = append(merged, u)
	}

	if len(merged) > constant.MaxSeenJobs {
		merged = merged[len(merged)-constant.MaxSeenJobs:]
	}

	dir := filepath.Dir(f.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.Marshal(merged)
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o644)
}
