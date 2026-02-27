package config

import (
	"os"
	"strings"
)

type Config struct {
	DiscordWebhookURL string
	LinkedInCookies   string
	SeenJobsPath      string
	Keywords          []string
	Locations         []string
}

func Load() *Config {
	keywords := getEnv("KEYWORDS", "software engineer")
	locations := getEnv("LOCATIONS", "Vietnam,Ho Chi Minh,Remote")

	return &Config{
		DiscordWebhookURL: getEnv("DISCORD_WEBHOOK_URL", ""),
		LinkedInCookies:   getEnv("LINKEDIN_COOKIES", ""),
		SeenJobsPath:      getEnv("SEEN_JOBS_PATH", "seen-jobs.json"),
		Keywords:          splitAndTrim(keywords),
		Locations:         splitAndTrim(locations),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
