package service

import (
	"fmt"
	"livebets/parse_sansabet/internal/entity"
	"strconv"
)

const (
	// Period indices
	PeriodMatch = 0
	PeriodTime1 = 1
	PeriodTime2 = 2

	PeriodSet1 = 1
	PeriodSet2 = 2
	PeriodSet3 = 3
	PeriodSet4 = 4
	PeriodSet5 = 5

	// Outcome types
	OutcomeWin1    = "Win1"
	OutcomeWinNone = "WinNone"
	OutcomeWin2    = "Win2"
	OutcomeMore    = "WinMore"
	OutcomeLess    = "WinLess"
)

type OddMapping struct {
	periodIndex int
	oddType     string
	team        string
}

var (
	win1x2Mappings = map[int64]OddMapping{
		// Football
		1:  {PeriodMatch, OutcomeWin1, ""},
		2:  {PeriodMatch, OutcomeWinNone, ""},
		10: {PeriodMatch, OutcomeWin2, ""},
		93: {PeriodTime1, OutcomeWin1, ""},
		94: {PeriodTime1, OutcomeWinNone, ""},
		95: {PeriodTime1, OutcomeWin2, ""},
		96: {PeriodTime2, OutcomeWin1, ""},
		97: {PeriodTime2, OutcomeWinNone, ""},
		98: {PeriodTime2, OutcomeWin2, ""},

		// Tennis
		691: {PeriodSet1, OutcomeWin1, ""},
		692: {PeriodSet1, OutcomeWin2, ""},
		693: {PeriodSet2, OutcomeWin1, ""},
		694: {PeriodSet2, OutcomeWin2, ""},
		695: {PeriodSet3, OutcomeWin1, ""},
		696: {PeriodSet3, OutcomeWin2, ""},
		697: {PeriodSet4, OutcomeWin1, ""},
		698: {PeriodSet4, OutcomeWin2, ""},
		699: {PeriodSet5, OutcomeWin1, ""},
		700: {PeriodSet5, OutcomeWin2, ""},
	}

	gamesMappings = map[int64]OddMapping{
		667: {PeriodSet1, OutcomeWin1, ""},
		668: {PeriodSet1, OutcomeWin2, ""},
		669: {PeriodSet2, OutcomeWin1, ""},
		670: {PeriodSet2, OutcomeWin2, ""},
		671: {PeriodSet3, OutcomeWin1, ""},
		672: {PeriodSet3, OutcomeWin2, ""},
		673: {PeriodSet4, OutcomeWin1, ""},
		674: {PeriodSet4, OutcomeWin2, ""},
		675: {PeriodSet5, OutcomeWin1, ""},
		676: {PeriodSet5, OutcomeWin2, ""},
	}

	totalsMappings = map[int64]OddMapping{
		// Football
		105: {PeriodMatch, OutcomeMore, ""},
		103: {PeriodMatch, OutcomeLess, ""},
		167: {PeriodTime1, OutcomeMore, ""},
		165: {PeriodTime1, OutcomeLess, ""},
		755: {PeriodTime2, OutcomeMore, ""},
		754: {PeriodTime2, OutcomeLess, ""},

		// Tennis
		666: {PeriodMatch, OutcomeMore, ""},
		665: {PeriodMatch, OutcomeLess, ""},
		658: {PeriodSet1, OutcomeMore, ""},
		657: {PeriodSet1, OutcomeLess, ""},
		660: {PeriodSet2, OutcomeMore, ""},
		659: {PeriodSet2, OutcomeLess, ""},
		662: {PeriodSet3, OutcomeMore, ""},
		661: {PeriodSet3, OutcomeLess, ""},
		664: {PeriodSet4, OutcomeMore, ""},
		663: {PeriodSet4, OutcomeLess, ""},
		// 666: {PeriodSet5, OutcomeMore, ""},
		// 665: {PeriodSet5, OutcomeLess, ""},
	}

	teamTotalsMappings = map[int64]OddMapping{
		// Football
		168: {PeriodMatch, OutcomeMore, "first"},
		169: {PeriodMatch, OutcomeLess, "first"},
		170: {PeriodMatch, OutcomeMore, "second"},
		171: {PeriodMatch, OutcomeLess, "second"},
		747: {PeriodTime1, OutcomeMore, "first"},
		746: {PeriodTime1, OutcomeLess, "first"},
		749: {PeriodTime1, OutcomeMore, "second"},
		748: {PeriodTime1, OutcomeLess, "second"},
	}
)

func ensureMapEntry[T any](m map[string]*T, key string) {
	if _, ok := m[key]; !ok {
		m[key] = new(T)
	}
}

func setWin1x2Value(win1x2 *entity.Win1x2Struct, oddType string, value float64) {
	switch oddType {
	case OutcomeWin1:
		win1x2.Win1 = entity.OddValue{Value: value}
	case OutcomeWinNone:
		win1x2.WinNone = entity.OddValue{Value: value}
	case OutcomeWin2:
		win1x2.Win2 = entity.OddValue{Value: value}
	}
}

func setTotalValue(total *entity.WinLessMore, oddType string, value float64) {
	switch oddType {
	case OutcomeMore:
		total.WinMore = entity.OddValue{Value: value}
	case OutcomeLess:
		total.WinLess = entity.OddValue{Value: value}
	}
}

func processHandicap(oddN int64, line string, oddValue float64, periods *[]entity.ResponsePeriod, homeScore, awayScore float64) {
	hcpLine, _ := strconv.ParseFloat(line, 64)

	handicapMappings := map[int64]struct {
		periodIndex int
		line        float64
		isWin1      bool
	}{
		// Football
		121: {PeriodMatch, hcpLine - 0.5, true},
		123: {PeriodMatch, (-1 * hcpLine) - 0.5, false},
		83:  {PeriodMatch, 0.5, true},
		85:  {PeriodMatch, 0.5, false},
		734: {PeriodMatch, (awayScore - homeScore) - 0.5, true},
		736: {PeriodMatch, (homeScore - awayScore) - 0.5, false},
		737: {PeriodTime1, (awayScore - homeScore) - 0.5, true},
		739: {PeriodTime1, (homeScore - awayScore) - 0.5, false},

		// Tennis
		1193: {PeriodMatch, hcpLine - 0.5, true},
		1194: {PeriodMatch, (-1 * hcpLine) - 0.5, false},
	}

	if mapping, ok := handicapMappings[oddN]; ok {
		lineStr := fmt.Sprintf("%.1f", mapping.line)
		ensureMapEntry((*periods)[mapping.periodIndex].Handicap, lineStr)

		if mapping.isWin1 {
			(*periods)[mapping.periodIndex].Handicap[lineStr].Win1 = entity.OddValue{
				Value: oddValue,
				Raw:   entity.OddRaw{Line: line, BetNum: oddN},
			}
		} else {
			(*periods)[mapping.periodIndex].Handicap[lineStr].Win2 = entity.OddValue{
				Value: oddValue,
				Raw:   entity.OddRaw{Line: line, BetNum: oddN},
			}
		}
	}
}
