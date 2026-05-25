// lol-tracker/internal/transport/sender.go
package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"lol-tracker/internal/models"
)

// Структура для отправки на ваш сервер отображения
type MatchData struct {
	Timestamp int64      `json:"Timestamp"`
	Player    PlayerStat `json:"Player"`
	Opponent  PlayerStat `json:"Opponent"`
	Events    []Event    `json:"Events"`
}

type PlayerStat struct {
	SummonerName string  `json:"SummonerName"`
	Champion     string  `json:"Champion"`
	Kills        int     `json:"Kills"`
	Deaths       int     `json:"Deaths"`
	Assists      int     `json:"Assists"`
	TotalGold    float64 `json:"TotalGold"`
	CS           int     `json:"CS"`
	GameTime     float64 `json:"GameTime"`
}

type Event struct {
	Type   string  `json:"Type"`
	Killer string  `json:"Killer"`
	Victim string  `json:"Victim"`
	Time   float64 `json:"Time"`
}

type DataSender struct {
	url string
	c   *http.Client
}

func NewDataSender(url string) *DataSender {
	return &DataSender{
		url: url,
		c: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (d *DataSender) SendToBackend(
	player *models.PlayerAnalytics,
	opponent *models.PlayerAnalytics,
	events []models.EventAnalytics,
) error {

	// Конвертируем события
	convertedEvents := make([]Event, 0, len(events))
	for _, e := range events {
		convertedEvents = append(convertedEvents, Event{
			Type:   e.Type,
			Time:   e.Time,
			Killer: e.Killer,
			Victim: e.Victim,
		})
	}

	// Формируем данные для сервера отображения
	matchData := MatchData{
		Timestamp: time.Now().Unix(),
		Player: PlayerStat{
			SummonerName: player.SummonerName,
			Champion:     player.Champion,
			Kills:        player.Kills,
			Deaths:       player.Deaths,
			Assists:      player.Assists,
			TotalGold:    player.TotalGold,
			CS:           player.CS,
			GameTime:     player.GameTime,
		},
		Opponent: PlayerStat{
			SummonerName: opponent.SummonerName,
			Champion:     opponent.Champion,
			Kills:        opponent.Kills,
			Deaths:       opponent.Deaths,
			Assists:      opponent.Assists,
			TotalGold:    opponent.TotalGold,
			CS:           opponent.CS,
			GameTime:     opponent.GameTime,
		},
		Events: convertedEvents,
	}

	b, err := json.Marshal(matchData)
	if err != nil {
		return err
	}

	// Отправляем на сервер отображения
	resp, err := d.c.Post(
		d.url,
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}