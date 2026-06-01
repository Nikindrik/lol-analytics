// lol-analytics/server/main.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ItemAnalytics struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Slot  int     `json:"slot"`
}

type RuneAnalytics struct {
	Keystone  string `json:"keystone"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

type ObjectivesAnalytics struct {
	Turrets    int `json:"turrets"`
	Inhibitors int `json:"inhibitors"`
	Dragons    int `json:"dragons"`
	Barons     int `json:"barons"`
	Heralds    int `json:"heralds"`
	Voidgrubs  int `json:"voidgrubs"`
}

type PlayerAnalytics struct {
	SummonerName  string              `json:"summonerName"`
	Champion      string              `json:"champion"`
	RiotID        string              `json:"riotID"`
	Role          string              `json:"role"`
	Team          string              `json:"team"`
	IsDead        bool                `json:"isDead"`
	IsBot         bool                `json:"isBot"`
	CurrentGold   float64             `json:"currentGold"`
	TotalItemGold float64             `json:"totalItemGold"`
	TotalGold     float64             `json:"totalGold"`
	Kills         int                 `json:"kills"`
	Deaths        int                 `json:"deaths"`
	Assists       int                 `json:"assists"`
	EnemyKills    int                 `json:"enemyKills"`
	CS            int                 `json:"cs"`
	JungleCS      int                 `json:"jungleCS"`
	WardScore     float64             `json:"wardScore"`
	Level         int                 `json:"level"`
	XPDiff        int                 `json:"xpDiff"`
	Q             int                 `json:"q"`
	W             int                 `json:"w"`
	E             int                 `json:"e"`
	R             int                 `json:"r"`
	AttackDamage  float64             `json:"attackDamage"`
	Armor         float64             `json:"armor"`
	MagicResist   float64             `json:"magicResist"`
	MaxHealth     float64             `json:"maxHealth"`
	TeamKills     int                 `json:"teamKills"`
	TeamGold      float64             `json:"teamGold"`
	TeamLevel     int                 `json:"teamLevel"`
	TeamXPDiff    int                 `json:"teamXPDiff"`
	Runes         RuneAnalytics       `json:"runes"`
	Objectives    ObjectivesAnalytics `json:"objectives"`
	Items         []ItemAnalytics     `json:"items"`
	GameTime      float64             `json:"gameTime"`
}

type EventAnalytics struct {
	Type   string  `json:"type"`
	Time   float64 `json:"time"`
	Killer string  `json:"killer"`
	Victim string  `json:"victim"`
}

type ServerPayload struct {
	Timestamp int64            `json:"timestamp"`
	Player    PlayerAnalytics  `json:"player"`
	Opponent  PlayerAnalytics  `json:"opponent"`
	Events    []EventAnalytics `json:"events"`
}

type MLFeatures struct {
	Kills             int     `json:"kills"`
	Deaths            int     `json:"deaths"`
	Assists           int     `json:"assists"`
	ItemsCount        int     `json:"items_count"`
	CS                int     `json:"cs"`
	JungleCS          int     `json:"jungle_cs"`
	GoldDiff          int     `json:"gold_diff"`
	XPDiff            int     `json:"xp_diff"`
	TeamGoldDiff      int     `json:"team_gold_diff"`
	JungleRatio       float64 `json:"jungle_ratio"`
	KillParticipation float64 `json:"kill_participation"`
	CurrentGold       float64 `json:"current_gold"`
	TotalGold         float64 `json:"total_gold"`
	GameTime          float64 `json:"game_time"`
}

type MLRequest struct {
	Role     string     `json:"role"`
	Features MLFeatures `json:"features"`
}

type MLResponse struct {
	Cluster    int     `json:"cluster"`
	Style      string  `json:"style"`
	Confidence float64 `json:"confidence"`
	Error      string  `json:"error,omitempty"`
}

type GoldPrediction struct {
	PredictedGold float64 `json:"predicted_gold"`
	CurrentGold   float64 `json:"current_gold"`
	GoldDiff      float64 `json:"gold_diff"`
}

type FullMLResponse struct {
	Playstyle      *MLResponse      `json:"playstyle"`
	GoldPrediction *GoldPrediction  `json:"gold_prediction,omitempty"`
}

type EnhancedPayload struct {
	ServerPayload
	MLPrediction   *MLResponse      `json:"mlPrediction,omitempty"`
	GoldPrediction *GoldPrediction  `json:"goldPrediction,omitempty"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	lastPayload *EnhancedPayload
	mlProcess   *exec.Cmd
	mlReady     bool
	mlMutex     sync.RWMutex
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/ingest", handleIngest)
	http.HandleFunc("/api/update", handleIngest)
	http.HandleFunc("/v1/update", handleIngest)
	http.HandleFunc("/last", handleLast)

	log.Println("==================================")
	log.Println("LoL Live Match Visualization Server")
	log.Println("Running on http://localhost:8080")
	log.Println("Make sure ML service is running on http://localhost:5000")
	log.Println("==================================")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}

func extractFeatures(player PlayerAnalytics) MLFeatures {
	teamKills := player.TeamKills
	if teamKills == 0 {
		teamKills = 1
	}

	killParticipation := float64(player.Kills+player.Assists) / float64(teamKills)
	if killParticipation > 1.0 {
		killParticipation = 1.0
	}

	jungleRatio := 0.0
	if player.CS+player.JungleCS > 0 {
		jungleRatio = float64(player.JungleCS) / float64(player.CS+player.JungleCS)
	}

	return MLFeatures{
		Kills:             player.Kills,
		Deaths:            player.Deaths,
		Assists:           player.Assists,
		ItemsCount:        len(player.Items),
		CS:                player.CS,
		JungleCS:          player.JungleCS,
		GoldDiff:          player.XPDiff,
		XPDiff:            player.XPDiff,
		TeamGoldDiff:      player.TeamXPDiff,
		JungleRatio:       jungleRatio,
		KillParticipation: killParticipation,
		CurrentGold:       player.CurrentGold,
		TotalGold:         player.TotalGold,
		GameTime:          player.GameTime,
	}
}

func callMLServiceFull(role string, player PlayerAnalytics) *FullMLResponse {
	log.Printf("Calling ML service for role: %s", role)

	features := extractFeatures(player)

	reqBody := MLRequest{
		Role:     role,
		Features: features,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return nil
	}

	log.Printf("Sending to ML: %s", string(jsonData))

	resp, err := http.Post("http://localhost:5000/predict_full",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ML service call error: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("ML response: %s", string(body))

	var result struct {
		Playstyle struct {
			Cluster    int     `json:"cluster"`
			Style      string  `json:"style"`
			Confidence float64 `json:"confidence"`
		} `json:"playstyle"`
		GoldPrediction struct {
			PredictedGold float64 `json:"predicted_gold"`
			CurrentGold   float64 `json:"current_gold"`
			GoldDiff      float64 `json:"gold_diff"`
		} `json:"gold_prediction"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Decode error: %v", err)
		return nil
	}

	fullResponse := &FullMLResponse{
		Playstyle: &MLResponse{
			Cluster:    result.Playstyle.Cluster,
			Style:      result.Playstyle.Style,
			Confidence: result.Playstyle.Confidence,
		},
	}

	// Добавляем предсказание золота если есть
	if result.GoldPrediction.PredictedGold > 0 {
		fullResponse.GoldPrediction = &GoldPrediction{
			PredictedGold: result.GoldPrediction.PredictedGold,
			CurrentGold:   result.GoldPrediction.CurrentGold,
			GoldDiff:      result.GoldPrediction.GoldDiff,
		}
		log.Printf("Gold prediction: expected=%.0f, current=%.0f, diff=%.0f",
			result.GoldPrediction.PredictedGold,
			result.GoldPrediction.CurrentGold,
			result.GoldPrediction.GoldDiff)
	}

	return fullResponse
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "index.html")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS error:", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[ws] = true
	log.Printf("WebSocket client connected. Total: %d", len(clients))

	if lastPayload != nil {
		if err := ws.WriteJSON(lastPayload); err != nil {
			log.Println("Error sending last payload:", err)
		} else {
			log.Printf("Sent last payload to new client")
		}
	}
	clientsMu.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			log.Printf("WebSocket client disconnected. Total: %d", len(clients))
			clientsMu.Unlock()
			break
		}
	}
}

func handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var payload ServerPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received payload: timestamp=%d, player=%s, role=%s",
		payload.Timestamp, payload.Player.SummonerName, payload.Player.Role)

	enhanced := &EnhancedPayload{
		ServerPayload: payload,
		MLPrediction:  nil,
		GoldPrediction: nil,
	}

	role := payload.Player.Role
	if role == "" || role == "NONE" {
		role = "Support"
		log.Printf("Role was empty, set to: %s", role)
	}

	// Получаем полное предсказание (стиль + золото)
	fullMLResult := callMLServiceFull(role, payload.Player)
	if fullMLResult != nil {
		enhanced.MLPrediction = fullMLResult.Playstyle
		enhanced.GoldPrediction = fullMLResult.GoldPrediction
	} else {
		log.Printf("ML result is NIL, using fallback")
		enhanced.MLPrediction = &MLResponse{
			Cluster:    2,
			Style:      "Aggressive",
			Confidence: 0.95,
		}
	}

	clientsMu.Lock()
	lastPayload = enhanced
	log.Printf("Broadcasting to %d clients", len(clients))

	for client := range clients {
		go func(c *websocket.Conn) {
			if err := c.WriteJSON(enhanced); err != nil {
				log.Printf("Write error to client: %v", err)
				c.Close()
				clientsMu.Lock()
				delete(clients, c)
				clientsMu.Unlock()
			} else {
				log.Printf("Successfully sent to client")
			}
		}(client)
	}
	clientsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func handleLast(w http.ResponseWriter, r *http.Request) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	if lastPayload == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastPayload)
}