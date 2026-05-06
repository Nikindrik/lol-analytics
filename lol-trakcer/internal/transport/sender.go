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
		c:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (d *DataSender) SendToBackend(p *models.PlayerAnalytics, e []models.EventAnalytics) error {

	body := models.ServerPayload{
		Timestamp: time.Now().Unix(),
		Player:    *p,
		Events:    e,
	}

	b, _ := json.Marshal(body)

	_, err := d.c.Post(d.url, "application/json", bytes.NewBuffer(b))
	return err
}
