package analytics

import "lol-tracker/internal/models"

type GameTracker struct{}

func NewGameTracker() *GameTracker { return &GameTracker{} }

func buildPlayerAnalytics(
	data *models.GameDataResponse,
	p *models.Player,
	teamKills int,
	teamGold float64,
	enemyKills int,
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

	totalGold := totalItemGold + data.ActivePlayer.CurrentGold

	jungleCS := int(float64(p.Scores.CreepScore) * 0.3)

	return models.PlayerAnalytics{
		SummonerName: p.SummonerName,
		Champion:     p.ChampionName,

		RiotID: p.RiotID,
		Role:   p.Position,
		Team:   p.Team,

		CurrentGold:   data.ActivePlayer.CurrentGold,
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

		Q: data.ActivePlayer.Abilities.Q.AbilityLevel,
		W: data.ActivePlayer.Abilities.W.AbilityLevel,
		E: data.ActivePlayer.Abilities.E.AbilityLevel,
		R: data.ActivePlayer.Abilities.R.AbilityLevel,

		AttackDamage: data.ActivePlayer.Stats.AttackDamage,
		Armor:        data.ActivePlayer.Stats.Armor,
		MagicResist:  data.ActivePlayer.Stats.MagicResist,
		MaxHealth:    data.ActivePlayer.Stats.MaxHealth,

		TeamKills: teamKills,
		TeamGold:  teamGold,

		Items:    items,
		GameTime: data.GameData.GameTime,
	}
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

	teamKills := 0
	teamGold := 0.0

	enemyKills := 0

	for _, p := range data.AllPlayers {

		if p.Team == player.Team {

			teamKills += p.Scores.Kills

			for _, it := range p.Items {
				teamGold += it.Price
			}

		} else {
			enemyKills += p.Scores.Kills
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
		teamKills,
		teamGold,
		enemyKills,
	)

	var enemy models.PlayerAnalytics

	if opponent != nil {

		enemy = buildPlayerAnalytics(
			data,
			opponent,
			enemyKills,
			teamGold,
			teamKills,
		)
	}

	return &me, &enemy, events
}