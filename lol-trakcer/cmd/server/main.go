package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"lol-tracker/internal/models"
)

var (
	last models.ServerPayload
	mu sync.RWMutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/v1/update", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			mu.RLock()
			defer mu.RUnlock()
			json.NewEncoder(w).Encode(last)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}

		var p models.ServerPayload
		json.NewDecoder(r.Body).Decode(&p)

		mu.Lock()
		last = p
		mu.Unlock()

		log.Println("received payload")
	})

	log.Println("server :8080")
	http.ListenAndServe(":8080", mux)
}
