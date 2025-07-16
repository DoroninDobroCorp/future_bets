package service

import (
	"fmt"
	"livebets/parse_maxbet/internal/entity"
	"livebets/shared"
	"strconv"
	"strings"
)

const (
	PeriodMatch = 0
	PeriodTime1 = 1
	PeriodTime2 = 2

	PeriodSet1 = 1
	PeriodSet2 = 2
	PeriodSet3 = 3
	PeriodSet4 = 4
	PeriodSet5 = 5
)

type OddMap struct {
	period int64
	team   string
}

var (
	win1X2Mappings = map[string]OddMap{
		// live

		// football
		"1": {PeriodMatch, "1"},
		"2": {PeriodMatch, "X"},
		"3": {PeriodMatch, "2"},

		"4": {PeriodTime1, "1"},
		"5": {PeriodTime1, "X"},
		"6": {PeriodTime1, "2"},

		"235": {PeriodTime2, "1"},
		"236": {PeriodTime2, "X"},
		"237": {PeriodTime2, "2"},

		// tennis
		"50510": {PeriodSet1, "1"},
		"50511": {PeriodSet1, "2"},

		"50512": {PeriodSet2, "1"},
		"50513": {PeriodSet2, "2"},

		"50514": {PeriodSet3, "1"},
		"50515": {PeriodSet3, "2"},

		"50516": {PeriodSet4, "1"},
		"50517": {PeriodSet4, "2"},

		"50518": {PeriodSet5, "1"},
		"50519": {PeriodSet5, "2"},
	}

	gamesMappings = map[string]OddMap{
		"50944": {PeriodSet1, "1"},
		"50945": {PeriodSet1, "2"},

		"50946": {PeriodSet2, "1"},
		"50947": {PeriodSet2, "2"},

		"50948": {PeriodSet3, "1"},
		"50949": {PeriodSet3, "2"},

		"50950": {PeriodSet4, "1"},
		"50951": {PeriodSet4, "2"},

		"50952": {PeriodSet5, "1"},
		"50953": {PeriodSet5, "2"},
	}

	totalsMappings = map[string]OddMap{
		//live

		// football
		"227": {PeriodMatch, "under"},
		"228": {PeriodMatch, "over"},

		"229": {PeriodTime1, "over"},
		"230": {PeriodTime1, "under"},

		// tennis
		"254": {PeriodMatch, "under"},
		"256": {PeriodMatch, "over"},

		"257": {PeriodSet1, "under"},
		"259": {PeriodSet1, "over"},

		"50520": {PeriodSet2, "under"},
		"50521": {PeriodSet2, "over"},
	}

	teamTotalsMappings = map[string]OddMap{
		// live

		// football
		"355": {PeriodMatch, "1 under"},
		"356": {PeriodMatch, "1 over"},
		"357": {PeriodMatch, "2 under"},
		"358": {PeriodMatch, "2 over"},

		"371": {PeriodTime1, "1 under"},
		"372": {PeriodTime1, "1 over"},
		"373": {PeriodTime1, "2 under"},
		"374": {PeriodTime1, "2 over"},
	}

	handicapsMappings = map[string]OddMap{
		// live

		// football
		"201": {PeriodMatch, "1 h"}, // handicap
		"202": {PeriodMatch, "X h"},
		"203": {PeriodMatch, "2 h"},

		"224": {PeriodTime1, "1 h"}, // handicap 1 time
		"225": {PeriodTime1, "X h"},
		"226": {PeriodTime1, "2 h"},

		"7": {PeriodMatch, "1 dc"}, // double chance
		"8": {PeriodMatch, "X dc"},
		"9": {PeriodMatch, "2 dc"},

		"397": {PeriodTime1, "1 dc"}, // double chance 1 time
		"398": {PeriodTime1, "X dc"},
		"399": {PeriodTime1, "2 dc"},

		"50915": {PeriodMatch, "1 ro"}, // rest of the match
		"50916": {PeriodMatch, "X ro"},
		"50917": {PeriodMatch, "2 ro"},

		// tennis
		"251": {PeriodMatch, "1 h"},
		"253": {PeriodMatch, "2 h"},
	}
)

func ensureMapEntry[T any](m map[string]*T, key string) {
	if _, ok := m[key]; !ok {
		m[key] = new(T)
	}
}

func processWin1x2(win1x2 *shared.Win1x2Struct, bet entity.Bet) {
	for key, coef := range bet.Coefs {
		oddMap := win1X2Mappings[key]

		switch oddMap.team {
		case "1":
			win1x2.Win1 = shared.Odd{Value: coef.Value}

		case "X":
			win1x2.WinNone = shared.Odd{Value: coef.Value}

		case "2":
			win1x2.Win2 = shared.Odd{Value: coef.Value}
		}
	}
}

func processGamesLive(games map[string]*shared.Win1x2Struct, bet entity.Bet) {
	for key, coef := range bet.Coefs {
		oddMap := gamesMappings[key]
		team, game := oddMap.team, bet.Line

		ensureMapEntry(games, game)

		switch team {
		case "1":
			games[game].Win1 = shared.Odd{Value: coef.Value}
		case "2":
			games[game].Win2 = shared.Odd{Value: coef.Value}
		}
	}
}

func processTotalLive(totals map[string]*shared.WinLessMore, bet entity.Bet, isTeamTotal bool) {
	line := strings.ReplaceAll(bet.Line, "total=", "")
	lineNum, _ := strconv.ParseFloat(line, 64)
	line = strconv.FormatFloat(lineNum, 'f', -1, 64)

	if _, ok := totals[line]; !ok {
		totals[line] = &shared.WinLessMore{}
	}

	for key, coef := range bet.Coefs {
		var oddMap OddMap
		if isTeamTotal {
			oddMap = teamTotalsMappings[key]
		} else {
			oddMap = totalsMappings[key]
		}

		totalType := oddMap.team

		if strings.HasSuffix(totalType, "over") {
			totals[line].WinMore = shared.Odd{Value: coef.Value}

		} else if strings.HasSuffix(totalType, "under") {
			totals[line].WinLess = shared.Odd{Value: coef.Value}
		}
	}
}

func processHandicapLive(handicaps map[string]*shared.WinHandicap, bet entity.Bet) {
	for key, coef := range bet.Coefs {
		oddMap := handicapsMappings[key]
		splitted := strings.Split(oddMap.team, " ")

		team, hcpType := splitted[0], splitted[1]
		line := bet.Line

		switch hcpType {
		case "h":
			if strings.Contains(line, ":") {
				scores := strings.Split(strings.ReplaceAll(line, "hcp=", ""), ":")
				fmt.Printf("line: %s\n", line)
				score1, _ := strconv.ParseFloat(scores[0], 64)
				score2, _ := strconv.ParseFloat(scores[1], 64)

				if team == "1" {
					lineStr := fmt.Sprintf("%.1f", (score2-score1)-0.5)
					ensureMapEntry(handicaps, lineStr)
					handicaps[lineStr].Win1 = shared.Odd{Value: coef.Value}

				} else if team == "2" {
					lineStr := fmt.Sprintf("%.1f", (score1-score2)-0.5)
					ensureMapEntry(handicaps, lineStr)
					handicaps[lineStr].Win2 = shared.Odd{Value: coef.Value}
				}
			} else {
				lineNum, _ := strconv.ParseFloat(strings.ReplaceAll(line, "hcp=", ""), 64)

				if team == "1" {
					lineStr := fmt.Sprintf("%.1f", lineNum-0.5)
					ensureMapEntry(handicaps, lineStr)
					handicaps[lineStr].Win1 = shared.Odd{Value: coef.Value}

				} else if team == "2" {
					lineStr := fmt.Sprintf("%.1f", (-1*lineNum)-0.5)
					ensureMapEntry(handicaps, lineStr)
					handicaps[lineStr].Win2 = shared.Odd{Value: coef.Value}
				}
			}

		case "dc":
			lineStr := "0.5"
			ensureMapEntry(handicaps, lineStr)

			if team == "1" {
				handicaps[lineStr].Win1 = shared.Odd{Value: coef.Value}
			} else if team == "2" {
				handicaps[lineStr].Win2 = shared.Odd{Value: coef.Value}
			}

		case "ro":
			score := strings.Split(strings.ReplaceAll(line, "score=", ""), ":")
			score1, _ := strconv.ParseFloat(score[0], 64)
			score2, _ := strconv.ParseFloat(score[1], 64)

			if team == "1" {
				lineStr := fmt.Sprintf("%.1f", (score2-score1)-0.5)
				ensureMapEntry(handicaps, lineStr)
				handicaps[lineStr].Win1 = shared.Odd{Value: coef.Value}

			} else if team == "2" {
				lineStr := fmt.Sprintf("%.1f", (score1-score2)-0.5)
				ensureMapEntry(handicaps, lineStr)
				handicaps[lineStr].Win2 = shared.Odd{Value: coef.Value}
			}
		}
	}
}
