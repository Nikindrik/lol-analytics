package models

type GameDataResponse struct {
	ActivePlayer ActivePlayer `json:"activePlayer"`
	AllPlayers []Player `json:"allPlayers"`
	Events EventsData `json:"events"`
	GameData GameInfo `json:"gameData"`
}

type ActivePlayer struct {
	SummonerName string `json:"summonerName"`
	CurrentGold float64 `json:"currentGold"`
}

type Player struct {
	SummonerName string
	ChampionName string
	Level int
	Items []Item
	Scores PlayerScores
}

type Item struct {
	DisplayName string
	Price float64
	Slot int
}

type PlayerScores struct {
	Kills int
	Deaths int
	Assists int
	CreepScore int
	WardScore float64
}

type EventsData struct {
	Events []GameEvent
}

type GameEvent struct {
	EventName string
	EventTime float64
	KillerName string
	VictimName string
}

type GameInfo struct {
	GameTime float64
}

type PlayerAnalytics struct {
	SummonerName string
	Champion string
	CurrentGold float64
	TotalItemGold float64
	Kills int
	Deaths int
	Assists int
	CS int
	WardScore float64
	Level int
	Items []ItemAnalytics
	GameTime float64
}

type ItemAnalytics struct {
	Name string
	Price float64
	Slot int
}

type EventAnalytics struct {
	Type string
	Time float64
	Killer string
	Victim string
}

type ServerPayload struct {
	Timestamp int64
	Player PlayerAnalytics
	Events []EventAnalytics
}
