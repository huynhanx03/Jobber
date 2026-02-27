package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"jobber/internal/core/entity"
)

type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type discordEmbed struct {
	Title     string         `json:"title"`
	URL       string         `json:"url,omitempty"`
	Color     int            `json:"color"`
	Fields    []discordField `json:"fields,omitempty"`
	Footer    *discordFooter `json:"footer,omitempty"`
	Timestamp string         `json:"timestamp,omitempty"`
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordFooter struct {
	Text string `json:"text"`
}

type discordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []discordEmbed `json:"embeds"`
}

func (d *DiscordNotifier) Send(ctx context.Context, jobs []entity.Job) error {
	for _, job := range jobs {
		embed := d.buildEmbed(job)
		payload := discordPayload{Embeds: []discordEmbed{embed}}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal discord payload: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.webhookURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("create discord request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := d.client.Do(req)
		if err != nil {
			return fmt.Errorf("send discord webhook: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(2 * time.Second)
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (d *DiscordNotifier) buildEmbed(job entity.Job) discordEmbed {
	return discordEmbed{
		Title: fmt.Sprintf("💼 %s", job.Title),
		URL:   job.URL,
		Color: sourceColor(job.Source),
		Fields: []discordField{
			{Name: "🏢 Company", Value: job.Company, Inline: false},
			{Name: "📍 Location", Value: job.Location, Inline: false},
			{Name: "💰 Salary", Value: nonEmpty(job.Salary, "Negotiable"), Inline: false},
			{Name: " Posted", Value: nonEmpty(job.PostedDate, "Recent"), Inline: false},
		},
		Footer:    &discordFooter{Text: job.Source},
		Timestamp: job.ScrapedAt.Format(time.RFC3339),
	}
}

func sourceColor(source string) int {
	colors := map[string]int{
		"Indeed": 0x003A9B,
		"ITViec": 0xE74C3C,
		"TopDev": 0x8E44AD,
	}
	if c, ok := colors[source]; ok {
		return c
	}
	return 0x95A5A6
}

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
