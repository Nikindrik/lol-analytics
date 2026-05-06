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
	total := 0.0

	for _, it := range p.Items {
		items = append(items, models.ItemAnalytics{
			Name: it.DisplayName,
			Price: it.Price,
			Slot: it.Slot,
		})
		total += it.Price
	}

	events := []models.EventAnalytics{}
	for _, e := range data.Events.Events {
		events = append(events, models.EventAnalytics{
			Type: e.EventName,
			Time: e.EventTime,
			Killer: e.KillerName,
			Victim: e.VictimName,
		})
	}

	return &models.PlayerAnalytics{
		SummonerName: p.SummonerName,
		Champion: p.ChampionName,
		CurrentGold: data.ActivePlayer.CurrentGold,
		TotalItemGold: total,
		Kills: p.Scores.Kills,
		Deaths: p.Scores.Deaths,
		Assists: p.Scores.Assists,
		CS: p.Scores.CreepScore,
		WardScore: p.Scores.WardScore,
		Level: p.Level,
		Items: items,
		GameTime: data.GameData.GameTime,
	}, events
}
