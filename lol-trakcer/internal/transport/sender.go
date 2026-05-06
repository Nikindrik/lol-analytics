package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"lol-tracker/internal/models"
)

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

	body := models.ServerPayload{
		Timestamp: time.Now().Unix(),
		Player:    *player,
		Opponent:  *opponent,
		Events:    events,
	}

	b, _ := json.Marshal(body)

	_, err := d.c.Post(
		d.url,
		"application/json",
		bytes.NewBuffer(b),
	)

	return err
}
