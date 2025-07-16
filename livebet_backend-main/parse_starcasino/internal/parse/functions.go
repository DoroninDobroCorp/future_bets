package parse

import (
	"livebets/parse_starcasino/internal/entity"
	"livebets/shared"
	"strconv"
	"strings"
)

// Periods

func matchHas(number int, periods []shared.PeriodData) []shared.PeriodData {

	if number < len(periods) {
		return periods
	}

	for i := len(periods); i <= number; i++ {
		newPeriod := shared.PeriodData{
			Games:            make(map[string]*shared.Win1x2Struct),
			Totals:           make(map[string]*shared.WinLessMore),
			Handicap:         make(map[string]*shared.WinHandicap),
			FirstTeamTotals:  make(map[string]*shared.WinLessMore),
			SecondTeamTotals: make(map[string]*shared.WinLessMore),
		}

		periods = append(periods, newPeriod)
	}

	return periods
}

// Win1x2

func fillWin1x2(win1x2 *shared.Win1x2Struct, odd *entity.Odd) {
	switch odd.TypeId {
	case 1:
		win1x2.Win1.Value = odd.Price
	case 2:
		win1x2.WinNone.Value = odd.Price
	case 3:
		win1x2.Win2.Value = odd.Price
	}
}

func fillTennisWin1x2(periods *[]shared.PeriodData, odd *entity.Odd, betName string, betSv string) {
	if betName == "Winner" {
		fillWin1x2(&(*periods)[0].Win1x2, odd)

	} else if strings.HasSuffix(betName, " - winner") {
		// betName = "First set - winner" or "First set game 7 - winner"

		setNumber := betSv[0] - 48 // "1" -> 1
		if setNumber < 1 || setNumber > 5 {
			return
		}
		*periods = matchHas(int(setNumber), *periods)

		if strings.HasSuffix(betName, " set - winner") {
			// Set winner: betName = "First set - winner"
			fillWin1x2(&(*periods)[setNumber].Win1x2, odd)

		} else if index := strings.Index(betName, " game "); index >= 0 {
			// Game winner: betName = "First set game 7 - winner"
			gameNumber := betSv[2:] // "1|7" => "7"

			win1x2Pointer, ok := (*periods)[setNumber].Games[gameNumber]
			if !ok {
				win1x2Pointer = &shared.Win1x2Struct{}
				(*periods)[setNumber].Games[gameNumber] = win1x2Pointer
			}

			fillWin1x2(win1x2Pointer, odd)
		}
	}
}

// Totals

func fillFootballTotals(periods *[]shared.PeriodData, odd *entity.Odd, betName string, betTypeId int64, betSv string) {
	if betName == "Total" {
		fillTotalsMap((*periods)[0].Totals, odd, betSv)

	} else if betName == "1st half - total" {
		fillTotalsMap((*periods)[1].Totals, odd, betSv)

	} else if betName == "2nd half - total" {
		fillTotalsMap((*periods)[2].Totals, odd, betSv)

	} else if strings.HasSuffix(betName, "total goals") {
		// Teams
		switch betTypeId {
		case 19:
			fillTotalsMap((*periods)[0].FirstTeamTotals, odd, betSv)
		case 20:
			fillTotalsMap((*periods)[0].SecondTeamTotals, odd, betSv)
		}

	} else if strings.HasSuffix(betName, "total") {
		// Half and teams
		half := betName[:8]
		if half == "1st half" {
			switch betTypeId {
			case 69:
				fillTotalsMap((*periods)[1].FirstTeamTotals, odd, betSv)
			case 70:
				fillTotalsMap((*periods)[1].SecondTeamTotals, odd, betSv)
			}
		} else if half == "2nd half" {
			switch betTypeId {
			case 91:
				fillTotalsMap((*periods)[2].FirstTeamTotals, odd, betSv)
			case 92:
				fillTotalsMap((*periods)[2].SecondTeamTotals, odd, betSv)
			}
		}
	}
}

func fillTotalsMap(totals map[string]*shared.WinLessMore, odd *entity.Odd, betSv string) {
	line := odd.Sv
	if len(line) == 0 {
		line = betSv
	}

	totalsLine, ok := totals[line]
	if !ok {
		totalsLine = &shared.WinLessMore{}
		totals[line] = totalsLine
	}

	switch odd.TypeId {
	case 12:
		totalsLine.WinMore.Value = odd.Price
	case 13:
		totalsLine.WinLess.Value = odd.Price
	}
}

func fillTennisTotals(periods *[]shared.PeriodData, odd *entity.Odd, betName string, betTypeId int64, betSv string) {
	if betName == "Total games" {
		fillTotalsMap((*periods)[0].Totals, odd, betSv)

		// } else if betName == "Total sets" {

	} else if strings.HasSuffix(betName, " - total games") {
		// Set total: betName = "Second set - total games"
		setNumber := betSv[0] - 48 // "2|9.5" -> 2
		if setNumber < 1 || setNumber > 5 {
			return
		}
		*periods = matchHas(int(setNumber), *periods)

		fillTotalsMap((*periods)[setNumber].Totals, odd, betSv)

	} else if strings.HasSuffix(betName, " total games") {
		// Team total: betName = "Alcaraz, Carlos total games"
		switch betTypeId {
		case 190:
			fillTotalsMap((*periods)[0].FirstTeamTotals, odd, betSv)
		case 191:
			fillTotalsMap((*periods)[0].SecondTeamTotals, odd, betSv)
		}
	}
}

// Handicap

func fillHandicap(handicap map[string]*shared.WinHandicap, odd *entity.Odd, betSv string) {
	line := odd.Sv
	if len(line) == 0 {
		line = betSv
	}
	line = strings.TrimPrefix(line, "+") // "+1.5" => "1.5"

	switch odd.TypeId {
	case 1714:
		handicapLine := getHandicapLine(handicap, line)
		handicapLine.Win1.Value = max(odd.Price, handicapLine.Win1.Value)

	case 1715:
		// line - меняем знак
		if line[0] == '-' {
			line = strings.TrimPrefix(line, "-")
		} else {
			line = "-" + line
		}

		handicapLine := getHandicapLine(handicap, line)
		handicapLine.Win2.Value = max(odd.Price, handicapLine.Win2.Value)
	}
}

func fillFootballDoubleChance(handicap map[string]*shared.WinHandicap, odd *entity.Odd) {
	line := "0.5"
	handicapLine := getHandicapLine(handicap, line)

	switch odd.TypeId {
	case 9:
		handicapLine.Win1.Value = max(odd.Price, handicapLine.Win1.Value)
	case 11:
		handicapLine.Win2.Value = max(odd.Price, handicapLine.Win2.Value)
	}
}

func fillFootballDrawNoBet(handicap map[string]*shared.WinHandicap, odd *entity.Odd) {
	line := "0.0"
	handicapLine := getHandicapLine(handicap, line)

	switch odd.TypeId {
	case 1:
		handicapLine.Win1.Value = max(odd.Price, handicapLine.Win1.Value)
	case 3:
		handicapLine.Win2.Value = max(odd.Price, handicapLine.Win2.Value)
	}
}

func getHandicapLine(handicap map[string]*shared.WinHandicap, line string) *shared.WinHandicap {
	handicapLine, ok := handicap[line]
	if !ok {
		handicapLine = &shared.WinHandicap{}
		handicap[line] = handicapLine
	}
	return handicapLine
}

func fillFootballRestOfTheMatch(periods *[]shared.PeriodData, odd *entity.Odd, betName string) {
	if strings.HasPrefix(betName, "Which team wins the rest of the match") {
		fillFootballRestOfTheMatchHandicap((*periods)[0].Handicap, odd)

	} else if strings.HasPrefix(betName, "1st half - which team wins the rest") {
		fillFootballRestOfTheMatchHandicap((*periods)[1].Handicap, odd)

	} else if strings.HasPrefix(betName, "2nd half - which team wins the rest") {
		fillFootballRestOfTheMatchHandicap((*periods)[2].Handicap, odd)
	}
}

func fillFootballRestOfTheMatchHandicap(handicap map[string]*shared.WinHandicap, odd *entity.Odd) {
	scores := strings.Split(odd.Sv, ":")
	if len(scores) != 2 {
		return
	}

	score1, _ := strconv.ParseFloat(strings.TrimSpace(scores[0]), 64)
	score2, _ := strconv.ParseFloat(strings.TrimSpace(scores[1]), 64)

	switch odd.TypeId {
	case 1:
		scoreLine := (score2 - score1) - 0.5
		line := strconv.FormatFloat(scoreLine, 'f', 1, 64)
		handicapLine := getHandicapLine(handicap, line)
		handicapLine.Win1.Value = max(odd.Price, handicapLine.Win1.Value)
	case 3:
		scoreLine := (score1 - score2) - 0.5
		line := strconv.FormatFloat(scoreLine, 'f', 1, 64)
		handicapLine := getHandicapLine(handicap, line)
		handicapLine.Win2.Value = max(odd.Price, handicapLine.Win2.Value)
	}
}

func fillTennisHandicap(periods *[]shared.PeriodData, odd *entity.Odd, betName string, betSv string) {
	if betName == "Game handicap" {
		fillHandicap((*periods)[0].Handicap, odd, betSv)

		//} else if betName == "Set handicap" {

	} else if strings.HasSuffix(betName, " - game handicap") {
		// Set handicap: betName = "First set - game handicap"
		setNumber := betSv[0] - 48 // "1|-1.5" -> 1
		if setNumber < 1 || setNumber > 5 {
			return
		}
		*periods = matchHas(int(setNumber), *periods)

		fillHandicap((*periods)[setNumber].Handicap, odd, betSv)
	}
}
