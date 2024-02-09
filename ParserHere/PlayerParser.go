package main

import (
	"fmt"
	"math"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var (
	playerMap = make(map[string]int)
)

func (p *DemoParser) GetActivePlayer(c []*common.Player) {
	for _, player := range c {
		steamId := player.SteamID64

		if p.Match.Players == nil {
			p.Match.Players = make(map[int64]playerStat)
		}

		if _, exists := p.Match.Players[int64(steamId)]; exists {
			return
		} else {
			fmt.Println(player)
			p.Match.Players[int64(steamId)] = p.ThePlayer(player)
		}
	}
}

func (p *DemoParser) ThePlayer(player *common.Player) playerStat {

	return playerStat{
		UserName:        player.Name,
		SteamID:         player.SteamID64,
		Kills:           0,
		Deaths:          0,
		Assists:         0,
		HS:              0,
		HeadPercent:     0,
		ADR:             0,
		KAST:            0,
		KDRatio:         0,
		Firstkill:       0,
		FirstDeath:      0,
		FKDiff:          0,
		Round2k:         0,
		Round3k:         0,
		Round4k:         0,
		Round5k:         0,
		Totaldmg:        0,
		TradeKills:      0,
		TradeDeath:      0,
		CTkills:         0,
		Tkills:          0,
		AvgflshDuration: 0,
		WeaponKill:      p.allweapons(),
		ClanName:        player.TeamState.ClanName(),
		TotalUtilDmg:    0,
	}
}

func (p *DemoParser) allweapons() map[int]int {

	return make(map[int]int)
}

func (p *DemoParser) playerGetter(e events.RoundEnd) {

	roundTime := &timeRound{}

	TeamA := p.parser.GameState().TeamCounterTerrorists().Members()
	TeamB := p.parser.GameState().TeamTerrorists().Members()

	p.statSetter(TeamA)
	p.statSetter(TeamB)

	roundTime.RoundEndTime = p.parser.CurrentTime()

	roundTime.RoundLength = roundTime.RoundEndTime - roundTime.roundStartTime

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		p.Match.Rounds[p.state.round-1].Duration = roundTime.RoundLength
	}

}

func (p *DemoParser) statSetter(c []*common.Player) {
	//need to get dmg, hs percentage, ADR, kd Ratio, and multi kills, kills, assists, deaths
	//have to get headshot percentage later because it's in the kill handler
	gs := p.parser.GameState()
	for i := range c {
		steamId := c[i].SteamID64
		playerStat, exists := p.Match.Players[int64(steamId)]

		if !exists {
			continue
		}

		playerStat.Kills = c[i].Kills()
		playerStat.Deaths = c[i].Deaths()
		playerStat.Assists = c[i].Assists()
		playerStat.Totaldmg = c[i].TotalDamage()
		playerStat.ADR = math.Round(p.calcADR(playerStat.Totaldmg)*100) / 100
		playerStat.KDRatio = math.Round(p.calcKDRatio(playerStat.Kills, playerStat.Deaths)*100) / 100
		playerStat.HeadPercent = math.Round(p.calcHSPercent(playerStat.Kills, playerStat.HS)*100) / 100
		playerStat.TotalUtilDmg = c[i].UtilityDamage()
		playerStat.AvgKillsRnd = math.Round(float64(playerStat.Kills)/float64(gs.TotalRoundsPlayed())*100) / 100
		playerStat.AvgDeathsRnd = math.Round(float64(playerStat.Deaths)/float64(gs.TotalRoundsPlayed())*100) / 100
		playerStat.AvgAssistsRnd = math.Round(float64(playerStat.Assists)/float64(gs.TotalRoundsPlayed())*100) / 100
		playerStat.ImpactPerRnd = math.Round((2.13*playerStat.AvgKillsRnd+0.42*playerStat.AvgAssistsRnd-0.41)*100) / 100
		//playerStat.KAST = math.Round(((float64(playerStat.Kills)+float64(playerStat.Assits)+float64(playerStat.RoundSurvived)+float64(playerStat.RoundTraded))/float64(gs.TotalRoundsPlayed()))*100) / 100

		playerName := c[i].Name

		multiKillCheck := c[i].Kills() - playerMap[playerName]

		switch {
		case multiKillCheck == 2:
			playerStat.Round2k++
		case multiKillCheck == 3:
			playerStat.Round3k++
		case multiKillCheck == 4:
			playerStat.Round4k++
		case multiKillCheck == 5:
			playerStat.Round5k++
		}

		p.Match.Players[int64(steamId)] = playerStat
		//fmt.Println(p.Match.Players[int64(steamId)].UserName, p.Match.Players[int64(steamId)].Kills, p.Match.Players[int64(steamId)].Deaths)
	}
}

func (p *DemoParser) GetPresRoundKill(e events.RoundFreezetimeEnd) {

	TeamA := p.parser.GameState().TeamCounterTerrorists().Members()
	TeamB := p.parser.GameState().TeamTerrorists().Members()

	p.printThis(TeamA)
	p.printThis(TeamB)

}

func (p *DemoParser) printThis(c []*common.Player) {

	for _, v := range c {
		playerMap[v.Name] = v.Kills()
	}

}

func (p *DemoParser) calcADR(dmg int) float64 {
	roundsPlayed := p.parser.GameState().TotalRoundsPlayed()

	adr := float64(dmg) / float64(roundsPlayed)

	return adr
}

func (p *DemoParser) calcKDRatio(kills, deaths int) float64 {

	ratio := float64(kills) / float64(deaths)

	return ratio
}

func (p *DemoParser) calcHSPercent(kills, headshots int) float64 {
	return float64(headshots) / float64(kills)
}

func (p *DemoParser) KillHandler(e events.Kill) {

	if e.Killer == nil || e.Victim == nil {
		return
	}

	if p.parser.GameState().IsWarmupPeriod() {
		p.state.WarmupKills = append(p.state.WarmupKills, e)
		return
	}

	if e.IsHeadshot {
		p.AddHeadshot(e.Killer)
	}
	var assistorName string
	if e.Assister != nil {
		assistorName = e.Assister.Name
	}

	if e.Killer.ActiveWeapon() == nil {
		return
	}

	if p.state.round > 0 && p.state.round <= len(p.Match.Rounds) {
		if p.Match.Rounds[p.state.round-1].KillARound == nil {
			p.Match.Rounds[p.state.round-1].KillARound = make(map[int]RoundKill)
		}
		count := len(p.Match.Rounds[p.state.round-1].KillARound) + 1
		if _, exists := p.Match.Rounds[p.state.round-1].KillARound[count]; exists {
			return
		} else {
			p.Match.Rounds[p.state.round-1].KillARound[count] = RoundKill{
				TimeOfKill:       p.parser.CurrentTime(),
				Killer:           e.Killer.Name,
				KillerId:         e.Killer.SteamID64,
				VictId:           e.Victim.SteamID64,
				Victim:           e.Victim.Name,
				Assistor:         assistorName,
				KillerTeamString: e.Killer.TeamState.ClanName(),
				VictimTeamString: e.Victim.TeamState.ClanName(),
				IsHeadshot:       e.IsHeadshot,
				IsFlashed:        e.Victim.IsBlinded(),
				VictFlashDur:     e.Victim.GetFlashDuration(),
				//yikes
				//Dist:         math.Round(DistForm(e.Killer.Position(), e.Victim.Position())*100) / 100,
				KillerTeam: int(e.Killer.TeamState.Team()),
				VictTeam:   int(e.Victim.TeamState.Team()),
			}
			count++
		}

		if p.Match.Rounds[p.state.round-1].FirstKillCount == 0 {
			p.addFirst(e.Killer, e.Victim)
			p.Match.Rounds[p.state.round-1].FirstKillCount++
		}

		p.updateWeaponKills(e.Killer, e.Weapon.Type)
		//p.IsFlashed(e.Victim)
	}
}

func (p *DemoParser) IsFlashed(c *common.Player) {
	playerId := c.SteamID64

	playerStat, exists := p.Match.Players[int64(playerId)]

	if !exists {
		return
	}

	flashDurationInSeconds := time.Duration(c.FlashDuration) * time.Second

	if flashDurationInSeconds >= 2*time.Second {
		playerStat.EffectiveFlashes++
		p.Match.Players[int64(playerId)] = playerStat
	}

}

func (p *DemoParser) AddHeadshot(c *common.Player) {

	playerId := c.SteamID64
	playerStat, exists := p.Match.Players[int64(playerId)]

	if !exists {
		return
	}

	playerStat.HS++

	p.Match.Players[int64(playerId)] = playerStat

}

func (p *DemoParser) updateWeaponKills(c *common.Player, weaponType common.EquipmentType) {
	playerId := c.SteamID64

	playerStat, exists := p.Match.Players[int64(playerId)]

	if !exists {
		return
	}

	playerStat.WeaponKill[int(weaponType)]++

	p.Match.Players[int64(playerId)] = playerStat
}

func (p *DemoParser) addFirst(c *common.Player, c2 *common.Player) {

	//c is killer and c2 is Victim
	playeridKiller := c.SteamID64
	playerIdVict := c2.SteamID64

	playerStatKiller, exists := p.Match.Players[int64(playeridKiller)]

	if !exists {
		return
	}

	playerStatKiller.Firstkill++
	p.Match.Players[int64(playeridKiller)] = playerStatKiller

	playerStatVict, exists2 := p.Match.Players[int64(playerIdVict)]

	if !exists2 {
		return
	}

	playerStatVict.FirstDeath++
	p.Match.Players[int64(playerIdVict)] = playerStatVict
}

/*
func (p *DemoParser) DistForm(killerPos, victPos r3.Vector, c *common.Player) {

	//need to get the dist for the kill then add it to the player or killer or something
	for i := range p.Match.Players {
		if p.Match.Players[i].UserName == c.Name {
			dist := math.Sqrt(math.Pow((killerPos.X-victPos.X), 2) + math.Pow((killerPos.Y-victPos.Y), 2) + math.Pow((killerPos.Z-victPos.Z), 2))
			p.Match.Players[i].TotalDist += dist
			p.Match.Players[i].AvgDist = dist / float64(c.Kills())
		}
	}
}

func DistForm(killerPos, victPos r3.Vector) float64 {

	return math.Sqrt(math.Pow((killerPos.X-victPos.X), 2) + math.Pow((killerPos.Y-victPos.Y), 2) + math.Pow((killerPos.Z-victPos.Z), 2))
}
*/
