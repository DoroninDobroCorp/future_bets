package parse

import (
	"livebets/parse_lobbet/internal/entity"
	"livebets/shared"
	"strconv"
	"strings"
)

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

func fillWin1x2(win1x2 *shared.Win1x2Struct, picks []entity.Pick, caption string) {
	for _, pick := range picks {
		switch pick.Caption {
		case caption + " 1":
			win1x2.Win1.Value = pick.OddValue
		case caption + " X":
			win1x2.WinNone.Value = pick.OddValue
		case caption + " 2":
			win1x2.Win2.Value = pick.OddValue
		}
	}
}

func fillTennisGameWin1x2(periods *[]shared.PeriodData, picks []entity.Pick) {
	win1x2 := shared.Win1x2Struct{}
	for _, pick := range picks {
		if len(pick.Caption) > 14 { // Caption = "10. gem 2. seta 1"
			gameNumber := pick.SpecialValue                           // "10"
			setString := strings.TrimPrefix(pick.Caption, gameNumber) // "10. gem 2. seta 1" -> ". gem 2. seta 1"
			setNumber := setString[6] - 48                            // ". gem 2. seta 1" -> 2
			if setNumber > 0 && setNumber <= 5 {
				*periods = matchHas(int(setNumber), *periods)

				pickCaption := setString[14:]
				switch pickCaption {
				case "1":
					win1x2.Win1.Value = pick.OddValue
				case "2":
					win1x2.Win2.Value = pick.OddValue
				}

				(*periods)[setNumber].Games[gameNumber] = &win1x2
			}
		}
	}
}

func fillBasketballHalfWin1x2(win1x2 *shared.Win1x2Struct, picks []entity.Pick) {
	for _, pick := range picks {
		switch pick.Caption {
		case "1":
			win1x2.Win1.Value = pick.OddValue
		case "X":
			win1x2.WinNone.Value = pick.OddValue
		case "2":
			win1x2.Win2.Value = pick.OddValue
		}
	}
}

func fillTotals(totals map[string]*shared.WinLessMore, picks []entity.Pick, caption string) {
	captionLess := caption + "<"
	captionMore := caption + ">"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		line := pick.SpecialValue
		totalsLine := getTotalLine(totals, line)

		switch pick.LiveBetPickLabel {
		case captionLess:
			totalsLine.WinLess.Value = pick.OddValue
		case captionMore:
			totalsLine.WinMore.Value = pick.OddValue
		}
	}
}

func fillBasketballTotals(totals map[string]*shared.WinLessMore, picks []entity.Pick, caption string) {
	captionLess := caption + "<"
	captionMore := caption + ">"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		line := pick.SpecialValue
		totalsLine := getTotalLine(totals, line)

		switch pick.Caption {
		case captionLess:
			totalsLine.WinLess.Value = pick.OddValue
		case captionMore:
			totalsLine.WinMore.Value = pick.OddValue
		}
	}
}

func fillHandicap(handicap map[string]*shared.WinHandicap, picks []entity.Pick, caption string) {
	caption1 := caption + " 1"
	caption2 := caption + " 2"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		lineValue, _ := strconv.ParseFloat(pick.SpecialValue, 64)
		switch pick.Caption {
		case caption1:
			line := floatToLine(lineValue - 0.5)
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win1.Value = pick.OddValue
		case caption2:
			line := floatToLine(-lineValue - 0.5)
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win2.Value = pick.OddValue
		}
	}
}

func fillHandicapScore(handicap map[string]*shared.WinHandicap, picks []entity.Pick, caption string, score1, score2 float64) {
	caption1 := caption + " 1"
	caption2 := caption + " 2"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		switch pick.Caption {
		case caption1:
			line := floatToLine((score2 - score1) - 0.5)
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win1.Value = pick.OddValue
		case caption2:
			line := floatToLine((score1 - score2) - 0.5)
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win2.Value = pick.OddValue
		}
	}
}

func fillHandicapDoubleChance(handicap map[string]*shared.WinHandicap, picks []entity.Pick, caption string) {
	caption1 := caption + " 1X"
	caption2 := caption + " X2"
	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		line := "0.5"
		handicapLine := getHandicapLine(handicap, line)

		switch pick.Caption {
		case caption1:
			handicapLine.Win1.Value = pick.OddValue
		case caption2:
			handicapLine.Win2.Value = pick.OddValue
		}
	}
}

func fillPreMatchFootballTotals(totals map[string]*shared.WinLessMore, picks []entity.Pick, caption string) {
	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		pickCaption := strings.TrimPrefix(pick.Caption, caption)
		if strings.HasPrefix(pickCaption, "0-") {
			lineNum, _ := strconv.ParseFloat(strings.TrimPrefix(pickCaption, "0-"), 64)
			line := floatToLine(lineNum + 0.5)
			totalsLine := getTotalLine(totals, line)
			totalsLine.WinLess.Value = pick.OddValue
		} else if strings.HasSuffix(pickCaption, "+") {
			lineNum, _ := strconv.ParseFloat(strings.TrimSuffix(pickCaption, "+"), 64)
			line := floatToLine(lineNum - 0.5)
			totalsLine := getTotalLine(totals, line)
			totalsLine.WinMore.Value = pick.OddValue
		} else if pickCaption == "0" {
			line := "0.5"
			totalsLine := getTotalLine(totals, line)
			totalsLine.WinLess.Value = pick.OddValue
		}
	}
}

func fillPreMatchTennisTotals(totals map[string]*shared.WinLessMore, picks []entity.Pick, caption string) {
	captionLess := caption + "<"
	captionMore := caption + ">"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		line := pick.SpecialValue
		if line == "" {
			continue
		}
		totalsLine := getTotalLine(totals, line)
		switch pick.Caption {
		case captionLess:
			totalsLine.WinLess.Value = pick.OddValue
		case captionMore:
			totalsLine.WinMore.Value = pick.OddValue
		}
	}
}

func fillTennisBasketballHandicap(handicap map[string]*shared.WinHandicap, picks []entity.Pick, caption string) {
	caption1 := caption + " 1"
	caption2 := caption + " 2"

	for _, pick := range picks {
		if pick.OddValue == 0 {
			continue
		}

		line := pick.SpecialValue
		switch pick.Caption {
		case caption1:
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win1.Value = pick.OddValue
		case caption2:
			// Меняем знак
			if strings.HasPrefix(line, "-") {
				line = strings.TrimPrefix(line, "-")
			} else {
				line = "-" + line
			}
			handicapLine := getHandicapLine(handicap, line)
			handicapLine.Win2.Value = pick.OddValue
		}
	}
}

func getTotalLine(totals map[string]*shared.WinLessMore, line string) *shared.WinLessMore {
	totalsLine, ok := totals[line]
	if !ok {
		totalsLine = &shared.WinLessMore{}
		totals[line] = totalsLine
	}
	return totalsLine
}

func getHandicapLine(handicap map[string]*shared.WinHandicap, line string) *shared.WinHandicap {
	handicapLine, ok := handicap[line]
	if !ok {
		handicapLine = &shared.WinHandicap{}
		handicap[line] = handicapLine
	}
	return handicapLine
}

func floatToLine(value float64) string {
	return strconv.FormatFloat(value, 'f', 1, 64)
}
