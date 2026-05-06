package analytics

import "lol-tracker/internal/models"

type GameTracker struct{}

func NewGameTracker() *GameTracker { return &GameTracker{} }

func (g *GameTracker) GetPlayerAnalytics(data *models.GameDataResponse) (*models.PlayerAnalytics, []models.EventAnalytics) {

	var p *models.Player

	for i := range data.AllPlayers {
		if data.AllPlayers[i].SummonerName == data.ActivePlayer.SummonerName {
			p = &data.AllPlayers[i]
			break
		}
	}

	if p == nil {
		return &models.PlayerAnalytics{}, nil
	}

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

	teamKills := 0
	teamGold := 0.0

	for _, pl := range data.AllPlayers {
		if pl.Team == p.Team {
			teamKills += pl.Scores.Kills

			for _, it := range pl.Items {
				teamGold += it.Price
			}
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

	return &models.PlayerAnalytics{
		SummonerName: p.SummonerName,
		Champion:     p.ChampionName,

		RiotID: p.RiotID,
		Role:   p.Position,

		CurrentGold:   data.ActivePlayer.CurrentGold,
		TotalItemGold: totalItemGold,
		TotalGold:     totalGold,

		Kills:   p.Scores.Kills,
		Deaths:  p.Scores.Deaths,
		Assists: p.Scores.Assists,

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
	}, events
}
