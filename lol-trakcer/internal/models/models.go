package models

type GameDataResponse struct {
	ActivePlayer ActivePlayer `json:"activePlayer"`
	AllPlayers   []Player     `json:"allPlayers"`
	Events       EventsData   `json:"events"`
	GameData     GameInfo     `json:"gameData"`
}

type ActivePlayer struct {
	SummonerName string  `json:"summonerName"`
	CurrentGold  float64 `json:"currentGold"`

	Abilities Abilities     `json:"abilities"`
	Stats     ChampionStats `json:"championStats"`
}

type Abilities struct {
	Q Ability `json:"Q"`
	W Ability `json:"W"`
	E Ability `json:"E"`
	R Ability `json:"R"`
}

type Ability struct {
	AbilityLevel int `json:"abilityLevel"`
}

type ChampionStats struct {
	AttackDamage float64 `json:"attackDamage"`
	Armor        float64 `json:"armor"`
	MagicResist  float64 `json:"magicResist"`
	MaxHealth    float64 `json:"maxHealth"`
}

type Player struct {
	SummonerName string       `json:"summonerName"`
	ChampionName string       `json:"championName"`
	Level        int          `json:"level"`
	Items        []Item       `json:"items"`
	Scores       PlayerScores `json:"scores"`

	Position string `json:"position"`
	Team     string `json:"team"`
	RiotID   string `json:"riotId"`

	IsDead bool `json:"isDead"`
	IsBot  bool `json:"isBot"`
}

type Item struct {
	DisplayName string  `json:"displayName"`
	Price       float64 `json:"price"`
	Slot        int     `json:"slot"`
}

type PlayerScores struct {
	Kills      int     `json:"kills"`
	Deaths     int     `json:"deaths"`
	Assists    int     `json:"assists"`
	CreepScore int     `json:"creepScore"`
	WardScore  float64 `json:"wardScore"`
}

type EventsData struct {
	Events []GameEvent `json:"Events"`
}

type GameEvent struct {
	EventName  string  `json:"EventName"`
	EventTime  float64 `json:"EventTime"`
	KillerName string  `json:"KillerName"`
	VictimName string  `json:"VictimName"`
}

type GameInfo struct {
	GameTime float64 `json:"gameTime"`
}

type PlayerAnalytics struct {
	SummonerName string
	Champion     string

	RiotID string
	Role   string
	Team   string

	IsDead bool
	IsBot  bool

	CurrentGold   float64
	TotalItemGold float64
	TotalGold     float64

	Kills   int
	Deaths  int
	Assists int

	EnemyKills int

	CS        int
	JungleCS  int
	WardScore float64

	Level int

	Q int
	W int
	E int
	R int

	AttackDamage float64
	Armor        float64
	MagicResist  float64
	MaxHealth    float64

	TeamKills int
	TeamGold  float64

	Items []ItemAnalytics

	GameTime float64
}

type ItemAnalytics struct {
	Name  string
	Price float64
	Slot  int
}

type EventAnalytics struct {
	Type   string
	Time   float64
	Killer string
	Victim string
}

type ServerPayload struct {
	Timestamp int64

	Player   PlayerAnalytics
	Opponent PlayerAnalytics

	Events []EventAnalytics
}
