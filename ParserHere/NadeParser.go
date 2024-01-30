package main

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func (p *DemoParser) PlayerFlashed(e events.PlayerFlashed) {
	playerId := e.Attacker.SteamID64

	playerStat, exists := p.Match.Players[int64(playerId)]

	if !exists {
		return
	}

	if e.FlashDuration().Milliseconds() >= 2000 {
		playerStat.EffectiveFlashes++
		p.Match.Players[int64(playerId)] = playerStat
	}

}

/*
func (p *DemoParser) Inferno(e events.InfernoStart) {
	Thrower := e.Inferno.Thrower().SteamID64

	playerStat, exists := p.Match.Players[int64(Thrower)]

	if !exists {
		return
	}

	fmt.Println(e.Inferno)

	playerStat.NadeThrowen[503]++

	p.Match.Players[int64(Thrower)] = playerStat
}

func (p *DemoParser) GernadesThrown(e events.GrenadeEventIf) {
	GerBase := e.Base()

	if GerBase.Thrower.String() == "GOTV" {
		return
	}

	fmt.Println(GerBase)
	playerId := GerBase.Thrower.SteamID64

	playerStat, exists := p.Match.Players[int64(playerId)]

	if !exists {
		return
	}

	playerStat.NadeThrowen[int(GerBase.GrenadeType)]++

	if GerBase.GrenadeType == common.EqFlash {
		playerStat.FlashesThrown++
	}

	p.Match.Players[int64(playerId)] = playerStat

}
*/
