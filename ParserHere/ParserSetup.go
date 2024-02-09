package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

/*
These are the structs where we are going to do a bunch of stuff to set up shit
idk what I'm doing tbh
actually just going to make this file the parser set up then pass it to other files.
This is also where are going to setup all the structs i guess
*/

type DemoParser struct {
	parser demoinfocs.Parser
	state  parsingState
	Match  *MatchInfo
}

type parsingState struct {
	round        int
	RoundOngoing bool
	TeamA        common.Team //maybe change this idk yet
	TeamB        common.Team
	WarmupKills  []events.Kill
}

type MatchInfo struct {
	Map      string
	WhoVsWho string
	Rounds   []RoundInformation
	Players  map[int64]playerStat
}

type RoundInformation struct {
	TeamNameA        string
	TeamNameB        string
	EconA            int
	EconB            int
	TypeofBuyA       string
	TypeofBuyB       string
	ScoreA           int
	ScoreB           int
	FirstKillCount   int      //need to think about trade kills later
	SurvivorsA       []string //need to get list of player names
	SurvivorsB       []string //need to get list of player names
	BombPlanted      bool
	PlayerPlanted    string
	RoundEndedReason string
	SideWon          string //need to change later
	KillARound       map[int]RoundKill
	Duration         time.Duration
}

type RoundKill struct {
	TimeOfKill       time.Duration
	Killer           string
	KillerId         uint64
	VictId           uint64
	Victim           string
	Assistor         string
	KillerTeamString string
	VictimTeamString string
	VictFlashDur     float32
	VictDmgTaken     int
	AttDmgTaken      int
	IsHeadshot       bool
	IsFlashed        bool
	Dist             float64
	KillerWeapon     int
	KillerTeam       int
	VictTeam         int
}

type playerStat struct {
	ImpactPerRnd     float64
	UserName         string
	SteamID          uint64
	Kills            int
	Deaths           int
	Assists          int
	HS               int
	HeadPercent      float64
	ADR              float64
	KAST             float64
	KDRatio          float64
	Firstkill        int
	FirstDeath       int
	FKDiff           int
	Round2k          int
	Round3k          int
	Round4k          int
	Round5k          int
	Totaldmg         int
	TradeKills       int
	TradeDeath       int
	CTkills          int
	Tkills           int
	EffectiveFlashes int
	AvgflshDuration  float64
	WeaponKill       map[int]int
	AvgDist          float64
	TotalDist        float64
	FlashesThrown    int
	ClanName         string
	TotalUtilDmg     int
	AvgKillsRnd      float64
	AvgDeathsRnd     float64
	AvgAssistsRnd    float64
	RoundSurvived    int
	RoundTraded      int
}

type playerPositions struct {
	User string
	X    float64
	Y    float64
	Z    float64
}

/*
This is Where are the parsing for the rounds should go
First we should set up the parser for the demos in this file as well
Then we should set of the functions on parsing the information on the rounds in here as well,
lastly we will find a way to pass the parser through two files. idk how tho
*/

func (p *DemoParser) startParsing(demoPath string) error {

	f, err := os.Open(demoPath)
	if err != nil {
		log.Panic("failed to open demo file: ", err)
		fmt.Println(err)
		// Ensure p.wg.Done() is called even if there's an error

	}
	defer f.Close()

	p.parser = dem.NewParser(f)
	defer p.parser.Close()
	p.Match = &MatchInfo{}

	p.parser.RegisterEventHandler(p.stateControler)
	p.parser.RegisterEventHandler(p.TeamSwitch)
	p.parser.RegisterEventHandler(p.RoundEcon)
	p.parser.RegisterEventHandler(p.ScoreUpdater)
	p.parser.RegisterEventHandler(p.MatchStartHandler)
	//p.parser.RegisterEventHandler(p.TeamSwitch)
	p.parser.RegisterEventHandler(p.PlayerAlive)
	p.parser.RegisterEventHandler(p.BombPlanted)
	p.parser.RegisterEventHandler(p.playerGetter)
	p.parser.RegisterEventHandler(p.KillHandler)
	p.parser.RegisterEventHandler(p.GetPresRoundKill)
	//p.parser.RegisterEventHandler(p.GernadesThrown)
	//p.parser.RegisterEventHandler(p.Inferno)
	p.parser.RegisterEventHandler(p.PlayerFlashed)

	err = p.parser.ParseToEnd()
	if err != nil {
		// Handle error, log, etc.
		//log.Fatal("Error Here", err)
		return nil
	}

	// Wait for the goroutine to finish

	return nil
}

func main() {
	demoDir := "C:\\Users\\iphon\\Desktop\\DEMOProject\\More_Demos"

	demoPaths, err := getDemoPaths(demoDir)
	if err != nil {
		log.Fatal("Error getting demo paths:", err)
	}

	for i, demoPath := range demoPaths {
		parser := &DemoParser{
			Match: &MatchInfo{Rounds: make([]RoundInformation, 0)},
		}
		fmt.Println("Parsing demo:", demoPath)

		err := parser.startParsing(demoPath)
		if err != nil {
			log.Printf("Error parsing demo %s: %v, Line 191\n", demoPath, err)
			continue
		}

		outputFileName := fmt.Sprintf("%d.json", i+1)

		//jsonDataMatchRounds, err := json.MarshalIndent(parser, "", "  ")
		jsonDataPlayers, err := json.MarshalIndent(parser, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Write JSON data to a file
		err = os.WriteFile(outputFileName, jsonDataPlayers, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		requestBodyRounds := bytes.NewBuffer(jsonDataPlayers)
		//PlayerRequestBody := bytes.NewBuffer(jsonDataPlayers)

		resp, err := http.Post("http://127.0.0.1:5000/MatchData", "application/json", requestBodyRounds)

		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		//defer respPlayers.Body.Close()
	}

}

func getDemoPaths(dir string) ([]string, error) {
	var demoPaths []string

	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue // Skip directories
		}

		demoPaths = append(demoPaths, filepath.Join(dir, item.Name()))
	}

	return demoPaths, nil
}
