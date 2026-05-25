package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Структуры под ваш JSON
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

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/api/update", handleParserUpdate) // Сюда парсер должен слать POST запросы

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка сервера: ", err)
	}
}

// Отдача HTML страницы
func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
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
	clientsMu.Unlock()

	// Держим соединение открытым
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			break
		}
	}
}

// Прием данных от вашего парсера
func handleParserUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	var data MatchData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Рассылаем полученные данные всем подключенным браузерам
	clientsMu.Lock()
	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
	clientsMu.Unlock()

	w.WriteHeader(http.StatusOK)
}