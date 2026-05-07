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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			fmt.Println("\nStopping tracker...")
			return

		case <-ticker.C:

			if !client.IsAvailable() {
				fmt.Print("\rWaiting for game...")
				continue
			}

			// Получение live game data
			game, err := client.GetGameData()
			if err != nil {
				log.Println("get game data:", err)
				continue
			}

			// Аналитика игрока + вражина + events
			player, opponent, events := tracker.GetPlayerAnalytics(game)

			// Дебаг костко
			formatter.PrintStatus(player)

			// Отправка на backend/mock-server
			err = sender.SendToBackend(
				player,
				opponent,
				events,
			)

			if err != nil {
				log.Println("send error:", err)
			}
		}
	}
}
