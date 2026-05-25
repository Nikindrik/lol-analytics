package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	// OpenAPI-compliant payload
	payload := models.ServerPayload{
		Timestamp: time.Now().UnixMilli(),
		Player:    *player,
		Opponent:  *opponent,
		Events:    events,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// Determine target URL: if user provided a path (e.g. /v1/update), use it as-is.
	target := d.url
	if u, perr := url.Parse(d.url); perr == nil {
		if u.Path == "" || u.Path == "/" {
			// no path -> use /ingest
			target = d.url + "/ingest"
		}
	}

	resp, err := d.c.Post(
		target,
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return fmt.Errorf("post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	return nil
}