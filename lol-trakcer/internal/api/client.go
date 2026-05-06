package api

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"lol-tracker/internal/models"
)

type LeagueClient struct {
	base string
	c *http.Client
}

func NewLeagueClient(base string) *LeagueClient {
	return &LeagueClient{
		base: base,
		c: &http.Client{
			Timeout: 5*time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (l *LeagueClient) IsAvailable() bool {
	resp, err := l.c.Get(l.base + "/liveclientdata/gamestats")
	if err != nil { return false }
	return resp.StatusCode == 200
}

func (l *LeagueClient) GetGameData() (*models.GameDataResponse, error) {
	resp, err := l.c.Get(l.base + "/liveclientdata/allgamedata")
	if err != nil { return nil, err }
	defer resp.Body.Close()

	b,_ := io.ReadAll(resp.Body)

	var g models.GameDataResponse
	json.Unmarshal(b,&g)
	return &g,nil
}
