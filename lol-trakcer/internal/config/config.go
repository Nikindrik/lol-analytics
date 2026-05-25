package config

import "os"

type Config struct {
	ServerURL     string
	LeagueBaseURL string
	PollInterval  string
}

func Load() Config {
	return Config{
		ServerURL:     get("LOL_SERVER_URL", "http://localhost:8080"),
		LeagueBaseURL: get("LOL_LEAGUE_URL", "https://127.0.0.1:2999"),
		PollInterval:  get("LOL_POLL_INTERVAL", "3s"),
	}
}

func get(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
