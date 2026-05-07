package analytics

import "lol-tracker/internal/models"

type GameTracker struct{}

func NewGameTracker() *GameTracker {
	return &GameTracker{}
}

func buildPlayerAnalytics(
	data *models.GameDataResponse,
	p *models.Player,
	teamKills int,
	teamGold float64,
	enemyKills int,
	isActive bool,
) models.PlayerAnalytics {

	items := []models.ItemAnalytics{}
	totalItemGold := 0.0

	for _, it := range p.Items {
		items = append(items, models.ItemAnalytics{
			Name:  it.DisplayName,
			Price: it.Price,
			Slot:  it.Slot,
		})

		totalItemGold += it.Price
	}

	currentGold := 0.0

	if isActive {
		currentGold = data.ActivePlayer.CurrentGold
	}

	totalGold := totalItemGold + currentGold

	jungleCS := 0

	if p.Position == "JUNGLE" {
		jungleCS = int(float64(p.Scores.CreepScore) * 0.7)
	}

	playerData := models.PlayerAnalytics{
		SummonerName: p.SummonerName,
		Champion:     p.ChampionName,

		RiotID: p.RiotID,
		Role:   p.Position,
		Team:   p.Team,

		IsDead: p.IsDead,
		IsBot:  p.IsBot,

		CurrentGold:   currentGold,
		TotalItemGold: totalItemGold,
		TotalGold:     totalGold,

		Kills:   p.Scores.Kills,
		Deaths:  p.Scores.Deaths,
		Assists: p.Scores.Assists,

		EnemyKills: enemyKills,

		CS:        p.Scores.CreepScore,
		JungleCS:  jungleCS,
		WardScore: p.Scores.WardScore,

		Level: p.Level,

		TeamKills: teamKills,
		TeamGold:  teamGold,

		Items: items,

		GameTime: data.GameData.GameTime,
	}

	// Только для active player доступны championStats/abilities
	if isActive {

		playerData.Q = data.ActivePlayer.Abilities.Q.AbilityLevel
		playerData.W = data.ActivePlayer.Abilities.W.AbilityLevel
		playerData.E = data.ActivePlayer.Abilities.E.AbilityLevel
		playerData.R = data.ActivePlayer.Abilities.R.AbilityLevel

		playerData.AttackDamage = data.ActivePlayer.Stats.AttackDamage
		playerData.Armor = data.ActivePlayer.Stats.Armor
		playerData.MagicResist = data.ActivePlayer.Stats.MagicResist
		playerData.MaxHealth = data.ActivePlayer.Stats.MaxHealth
	}

	return playerData
}

func (g *GameTracker) GetPlayerAnalytics(
	data *models.GameDataResponse,
) (*models.PlayerAnalytics, *models.PlayerAnalytics, []models.EventAnalytics) {

	var player *models.Player

	for i := range data.AllPlayers {

		if data.AllPlayers[i].SummonerName == data.ActivePlayer.SummonerName {
			player = &data.AllPlayers[i]
			break
		}
	}

	if player == nil {
		return &models.PlayerAnalytics{}, &models.PlayerAnalytics{}, nil
	}

	var opponent *models.Player

	for i := range data.AllPlayers {

		p := &data.AllPlayers[i]

		if p.Team != player.Team &&
			p.Position == player.Position {

			opponent = p
			break
		}
	}

	playerTeamKills := 0
	playerTeamGold := 0.0

	enemyTeamKills := 0
	enemyTeamGold := 0.0

	for _, p := range data.AllPlayers {

		total := 0.0

		for _, it := range p.Items {
			total += it.Price
		}

		if p.Team == player.Team {

			playerTeamKills += p.Scores.Kills
			playerTeamGold += total

		} else {

			enemyTeamKills += p.Scores.Kills
			enemyTeamGold += total
		}
	}

	events := []models.EventAnalytics{}

	for _, e := range data.Events.Events {

		events = append(events, models.EventAnalytics{
			Type:   e.EventName,
			Time:   e.EventTime,
			Killer: e.KillerName,
			Victim: e.VictimName,
		})
	}

	me := buildPlayerAnalytics(
		data,
		player,
		playerTeamKills,
		playerTeamGold,
		enemyTeamKills,
		true,
	)

	var enemy models.PlayerAnalytics

	if opponent != nil {

		enemy = buildPlayerAnalytics(
			data,
			opponent,
			enemyTeamKills,
			enemyTeamGold,
			playerTeamKills,
			false,
		)
	}

	return &me, &enemy, events
}
