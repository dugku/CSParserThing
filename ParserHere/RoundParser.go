package main

import (
	"fmt"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type timeRound struct {
	roundStartTime time.Duration
	RoundEndTime   time.Duration
	RoundLength    time.Duration
}

func (p *DemoParser) stateControler(e events.RoundStart) {
	p.state.RoundOngoing = true

	p.state.round++

	round := RoundInformation{}

	p.Match.Rounds = append(p.Match.Rounds, round)
}

func (p *DemoParser) MatchStartHandler(e events.MatchStart) {
	ActivePlayers := p.parser.GameState().Participants().Playing()
	p.GetActivePlayer(ActivePlayers)
	roundTime := &timeRound{}

	p.state.TeamA = common.TeamCounterTerrorists
	p.state.TeamB = common.TeamTerrorists

	CountTeam := p.parser.GameState().TeamCounterTerrorists().ClanName()
	TerrTeam := p.parser.GameState().TeamTerrorists().ClanName()

	p.Match.WhoVsWho = CountTeam + " vs " + TerrTeam

	p.Match.Map = p.parser.Header().MapName
	fmt.Println(p.parser.Header().MapName)

	roundTime.roundStartTime = p.parser.CurrentTime()
	//these two lines broke the code so just commenting them out just in case i need it again
	//or not
	/*
		p.Match.TeamAPlayers = p.parser.GameState().TeamCounterTerrorists().Members()
		p.Match.TeamBPlayers = p.parser.GameState().TeamTerrorists().Members()
	*/

}

func (p *DemoParser) TeamSwitch(e events.TeamSideSwitch) {

	p.state.TeamA = common.TeamTerrorists
	p.state.TeamB = common.TeamCounterTerrorists
}

func (p *DemoParser) ScoreUpdater(e events.ScoreUpdated) {

	//TeamA is Always CT and Team B is T
	//Then we will just switch then when they switch sides. or just keep it since I'm not tracking sides just yet
	ATeam := p.parser.GameState().TeamCounterTerrorists().ClanName()
	BTeam := p.parser.GameState().TeamTerrorists().ClanName()

	AScore := p.parser.GameState().TeamCounterTerrorists().Score()
	Bscore := p.parser.GameState().TeamTerrorists().Score()

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {

		p.Match.Rounds[p.state.round-1].ScoreA = AScore
		p.Match.Rounds[p.state.round-1].ScoreB = Bscore
		p.Match.Rounds[p.state.round-1].TeamNameA = ATeam
		p.Match.Rounds[p.state.round-1].TeamNameB = BTeam

	}
}

func (p *DemoParser) RoundEcon(e events.RoundFreezetimeEnd) {

	//Need to get the equipment value of each team then assess if it is a full buy or not
	TeamAEcon := p.parser.GameState().Team(common.TeamCounterTerrorists).CurrentEquipmentValue()
	TeamBEcon := p.parser.GameState().Team(common.TeamTerrorists).CurrentEquipmentValue()

	FullBuy := 20000
	SemiBuy := 10000
	SemiEco := 5000

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		//I forget pointer exist
		roundInfo := &p.Match.Rounds[p.state.round-1]

		roundInfo.EconA = TeamAEcon
		roundInfo.EconB = TeamBEcon

		roundInfo.TypeofBuyA = AssessBuytype(TeamAEcon, FullBuy, SemiBuy, SemiEco)

		roundInfo.TypeofBuyB = AssessBuytype(TeamBEcon, FullBuy, SemiBuy, SemiEco)
	}
}

func AssessBuytype(econValue, FullBuy, SemiBuy, SemiEco int) string {
	switch {
	case econValue >= FullBuy:
		return "Full Buy"
	case econValue >= SemiBuy && econValue < FullBuy:
		return "Semi Buy"
	case econValue >= SemiEco && econValue < SemiBuy:
		return "Force Buy"
	default:
		return "Eco"
	}
}

func (p *DemoParser) PlayerAlive(e events.RoundEnd) {

	ReasonsMap := map[int]string{
		1: "TargetBombed",
		7: "BombDefused",
		8: "CTWin",
		9: "TWin",
	}

	WinnerMap := map[int]string{
		2: "Terrorists",
		3: "Counter Terrorists",
	}

	Reason := e.Reason
	SideWon := e.Winner

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {

		roundInfo := &p.Match.Rounds[p.state.round-1]
		roundInfo.roundEndedReason = ReasonsMap[int(Reason)]
		roundInfo.SideWon = WinnerMap[int(SideWon)]
	}

	PlayersTeamA := p.parser.GameState().TeamCounterTerrorists().Members()
	PlayersTeamB := p.parser.GameState().TeamTerrorists().Members()

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		roundInfo := &p.Match.Rounds[p.state.round-1]
		for _, v := range PlayersTeamA {
			if v.IsAlive() {
				roundInfo.SurvivorsA = append(roundInfo.SurvivorsA, v.String())
			}
		}

		for _, v := range PlayersTeamB {
			if v.IsAlive() {
				roundInfo.SurvivorsB = append(roundInfo.SurvivorsB, v.String())
			}
		}
	}

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		roundInfo := &p.Match.Rounds[p.state.round-1]

		//man I hope this logic is right
		for key, _ := range roundInfo.KillARound {

			if key+1 < len(roundInfo.KillARound) {
				nextValue := roundInfo.KillARound[key+1]

				if roundInfo.KillARound[key].Killer == nextValue.Victim && ((nextValue.TimeOfKill - roundInfo.KillARound[key].TimeOfKill) < (6 * time.Second)) {
					TradeKillerId := nextValue.KillerId
					playerStat, exists := p.Match.Players[int64(TradeKillerId)]

					if !exists {
						return
					}
					playerStat.TradeKills++
					p.Match.Players[int64(TradeKillerId)] = playerStat
				}
			}
		}
	}
}

func (p *DemoParser) BombPlanted(e events.BombPlanted) {

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		roundInfo := &p.Match.Rounds[p.state.round-1]

		roundInfo.BombPlanted = true

		roundInfo.playerPlanted = e.Player.Name
	}
}
