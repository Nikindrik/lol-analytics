package analytics

import "lol-tracker/internal/models"

type GameTracker struct{}

func NewGameTracker() *GameTracker {
	return &GameTracker{}
}

func calcItemGold(items []models.Item) float64 {

	total := 0.0

	for _, it := range items {
		total += it.Price
	}

	return total
}

func stringContains(s string, sub string) bool {

	for i := 0; i <= len(s)-len(sub); i++ {

		if s[i:i+len(sub)] == sub {
			return true
		}
	}

	return false
}

func contains(s string, sub string) bool {

	return len(s) >= len(sub) &&
		(s == sub ||
			len(s) > len(sub) &&
				stringContains(s, sub))
}

func calculateObjectives(
	data *models.GameDataResponse,
	team string,
) models.ObjectivesAnalytics {

	obj := models.ObjectivesAnalytics{}

	for _, e := range data.Events.Events {

		switch e.EventName {

		case "TurretKilled":

			if team == "ORDER" &&
				contains(e.KillerName, "T100") {

				obj.Turrets++
			}

			if team == "CHAOS" &&
				contains(e.KillerName, "T200") {

				obj.Turrets++
			}

		case "InhibKilled":

			if team == "ORDER" &&
				contains(e.KillerName, "T100") {

				obj.Inhibitors++
			}

			if team == "CHAOS" &&
				contains(e.KillerName, "T200") {

				obj.Inhibitors++
			}

		case "DragonKill":

			if contains(e.KillerName, team) {
				obj.Dragons++
			}

		case "BaronKill":

			if contains(e.KillerName, team) {
				obj.Barons++
			}

		case "HeraldKill":

			if contains(e.KillerName, team) {
				obj.Heralds++
			}

		case "VoidGrubKill":

			if contains(e.KillerName, team) {
				obj.Voidgrubs++
			}
		}
	}

	return obj
}

func buildPlayerAnalytics(
	data *models.GameDataResponse,
	p *models.Player,
	opponent *models.Player,

	teamKills int,
	teamGold float64,
	teamLevel int,

	enemyKills int,
	enemyLevel int,

	isActive bool,
) models.PlayerAnalytics {

	items := []models.ItemAnalytics{}

	totalItemGold := calcItemGold(p.Items)

	for _, it := range p.Items {

		items = append(items, models.ItemAnalytics{
			Name:  it.DisplayName,
			Price: it.Price,
			Slot:  it.Slot,
		})
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

	xpDiff := 0

	if opponent != nil {
		xpDiff = p.Level - opponent.Level
	}

	objectives := calculateObjectives(data, p.Team)

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

		XPDiff: xpDiff,

		TeamKills: teamKills,
		TeamGold:  teamGold,
		TeamLevel: teamLevel,

		TeamXPDiff: teamLevel - enemyLevel,

		Runes: models.RuneAnalytics{
			Keystone:  p.Runes.Keystone.DisplayName,
			Primary:   p.Runes.PrimaryRuneTree.DisplayName,
			Secondary: p.Runes.SecondaryRuneTree.DisplayName,
		},

		Objectives: objectives,

		Items: items,

		GameTime: data.GameData.GameTime,
	}

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
	playerTeamLevel := 0

	enemyTeamKills := 0
	enemyTeamLevel := 0

	for _, p := range data.AllPlayers {

		total := calcItemGold(p.Items)

		if p.Team == player.Team {

			playerTeamKills += p.Scores.Kills
			playerTeamGold += total
			playerTeamLevel += p.Level

		} else {

			enemyTeamKills += p.Scores.Kills
			enemyTeamLevel += p.Level
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
		opponent,

		playerTeamKills,
		playerTeamGold,
		playerTeamLevel,

		enemyTeamKills,
		enemyTeamLevel,

		true,
	)

	var enemy models.PlayerAnalytics

	if opponent != nil {

		enemy = buildPlayerAnalytics(
			data,
			opponent,
			player,

			enemyTeamKills,
			calcItemGold(opponent.Items),
			enemyTeamLevel,

			playerTeamKills,
			playerTeamLevel,

			false,
		)
	}

	return &me, &enemy, events
}
