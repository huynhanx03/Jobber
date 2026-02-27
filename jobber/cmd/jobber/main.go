package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"jobber/config"
	"jobber/internal/di"
)

func main() {
	cfg := config.Load()

	if cfg.DiscordWebhookURL == "" {
		log.Fatal("❌ DISCORD_WEBHOOK_URL is required")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	svc, browser := di.Wire(cfg)
	defer browser.Close()

	if err := svc.Run(ctx); err != nil {
		log.Fatalf("❌ Job hunt failed: %v", err)
	}
}
