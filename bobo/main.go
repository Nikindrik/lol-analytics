package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// OpenAPI-compliant structures
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
	SummonerName string  `json:"summonerName"`
	Champion     string  `json:"champion"`
	RiotID       string  `json:"riotID"`
	Role         string  `json:"role"`
	Team         string  `json:"team"`
	IsDead       bool    `json:"isDead"`
	IsBot        bool    `json:"isBot"`
	CurrentGold  float64 `json:"currentGold"`
	TotalItemGold float64 `json:"totalItemGold"`
	TotalGold    float64 `json:"totalGold"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	EnemyKills   int     `json:"enemyKills"`
	CS           int     `json:"cs"`
	JungleCS     int     `json:"jungleCS"`
	WardScore    float64 `json:"wardScore"`
	Level        int     `json:"level"`
	XPDiff       int     `json:"xpDiff"`
	Q            int     `json:"q"`
	W            int     `json:"w"`
	E            int     `json:"e"`
	R            int     `json:"r"`
	AttackDamage float64 `json:"attackDamage"`
	Armor        float64 `json:"armor"`
	MagicResist  float64 `json:"magicResist"`
	MaxHealth    float64 `json:"maxHealth"`
	TeamKills    int     `json:"teamKills"`
	TeamGold     float64 `json:"teamGold"`
	TeamLevel    int     `json:"teamLevel"`
	TeamXPDiff   int     `json:"teamXPDiff"`
	Runes        RuneAnalytics `json:"runes"`
	Objectives   ObjectivesAnalytics `json:"objectives"`
	Items        []ItemAnalytics `json:"items"`
	GameTime     float64 `json:"gameTime"`
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

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	lastPayload *ServerPayload
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/ingest", handleIngest)
	http.HandleFunc("/api/update", handleIngest)
	http.HandleFunc("/v1/update", handleIngest)
	http.HandleFunc("/last", handleLast)

	log.Println("================================")
	log.Println("🚀 LoL Live Match Visualization Server запущен на http://localhost:8080")
	log.Println("📁 Working directory:", ".")
	log.Println("================================")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("❌ Ошибка сервера: ", err)
	}
}

// Отдача HTML страницы
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /", "path:", r.URL.Path)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "index.html")
	log.Println("HTML served")
}

// Подключение вебсокетов (для браузера)
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка WS:", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clients[ws] = true
	log.Println("WebSocket client connected. Total clients:", len(clients))
	
	// Отправляем последний payload если он есть
	if lastPayload != nil {
		if err := ws.WriteJSON(lastPayload); err != nil {
			log.Println("Ошибка отправки последнего payload:", err)
		}
	}
	clientsMu.Unlock()

	// Читаем сообщения (для keep-alive), игнорируем ошибки
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			log.Println("WebSocket client disconnected. Total clients:", len(clients))
			clientsMu.Unlock()
			break
		}
	}
}

// Прием данных от парсера (OpenAPI endpoint /ingest)
func handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	var payload ServerPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("received payload", payload.Timestamp, "player:", payload.Player.SummonerName)

	// Сохраняем последний payload
	clientsMu.Lock()
	lastPayload = &payload
	
	// Рассылаем полученные данные всем подключенным браузерам
	for client := range clients {
		go func(c *websocket.Conn) {
			if err := c.WriteJSON(payload); err != nil {
				c.Close()
				delete(clients, c)
			}
		}(client)
	}
	clientsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// Debug endpoint to view last payload
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
