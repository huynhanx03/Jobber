package entity

import "time"

type Job struct {
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	URL         string    `json:"url"`
	Description string    `json:"description,omitempty"`
	Salary      string    `json:"salary,omitempty"`
	Location    string    `json:"location,omitempty"`
	Source      string    `json:"source"`
	TechStack   string    `json:"tech_stack,omitempty"`
	PostedDate  string    `json:"posted_date,omitempty"`
	ScrapedAt   time.Time `json:"scraped_at"`
}
