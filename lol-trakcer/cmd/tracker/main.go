package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"lol-tracker/internal/analytics"
	"lol-tracker/internal/api"
	"lol-tracker/internal/config"
	"lol-tracker/internal/display"
	"lol-tracker/internal/transport"
)

func main() {
	cfg := config.Load()

	client := api.NewLeagueClient(cfg.LeagueBaseURL)
	tracker := analytics.NewGameTracker()
	formatter := display.NewFormatter()
	sender := transport.NewDataSender(cfg.ServerURL)

	fmt.Println("League Tracker started")
	fmt.Printf("Server: %s\n", cfg.ServerURL)

	dur, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !client.IsAvailable() {
				continue
			}

			game, err := client.GetGameData()
			if err != nil {
				log.Println(err)
				continue
			}

			player, events := tracker.GetPlayerAnalytics(game)
			formatter.PrintStatus(player)

			_ = sender.SendToBackend(player, events)
		}
	}
}
