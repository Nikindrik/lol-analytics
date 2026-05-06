package display

import (
	"fmt"
	"lol-tracker/internal/models"
)

type Formatter struct{}

func NewFormatter() *Formatter { return &Formatter{} }

func (f *Formatter) PrintStatus(p *models.PlayerAnalytics) {
	if p == nil {
		return
	}

	fmt.Printf("\r%s | KDA %d/%d/%d | CS %d | Gold %.0f | Level %d | Ward Score %d",
		p.GameTime, p.Champion, p.Kills, p.Deaths, p.Assists, p.CS, p.CurrentGold, p.Level, p.WardScore)

}
