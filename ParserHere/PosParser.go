package main

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func (p *DemoParser) GetPlayerPos(e events.FrameDone) {
	for _, gameStatePlayer := range p.parser.GameState().Participants().Playing() {
		PlayerWhere := playerPositions{
			User: gameStatePlayer.Name,
			X:    gameStatePlayer.Position().X,
			Y:    gameStatePlayer.Position().Y,
			Z:    gameStatePlayer.Position().Z,
		}
		if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
			p.Match.Rounds[p.state.round-1].Positions = append(p.Match.Rounds[p.state.round-1].Positions, PlayerWhere)
		}
	}
}
