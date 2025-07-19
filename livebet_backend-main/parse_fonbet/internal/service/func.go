package service

import (
	"fmt"
	"livebets/parse_fonbet/internal/entity"
	"strconv"
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

	OutcomeWin1    = "Win1"
	OutcomeWinNone = "WinNone"
	OutcomeWin2    = "Win2"
	OutcomeMore    = "WinMore"
	OutcomeLess    = "WinLess"
)

type OddMapping struct {
	oddType string
	team    string
}

var (
	win1x2Mappings = map[int64]OddMapping{
		921: {OutcomeWin1, ""},
		922: {OutcomeWinNone, ""},
		923: {OutcomeWin2, ""},
	}

	gamesMapping = map[int64]OddMapping{
		1747: {OutcomeWin1, ""},
		1748: {OutcomeWin2, ""},
		1750: {OutcomeWin1, ""},
		1751: {OutcomeWin2, ""},
		1753: {OutcomeWin1, ""},
		1754: {OutcomeWin2, ""},
		9961: {OutcomeWin1, ""},
		9962: {OutcomeWin2, ""},
	}

	totalsMappings = map[int64]OddMapping{
		930:  {OutcomeMore, ""},
		940:  {OutcomeMore, ""},
		1848: {OutcomeMore, ""},
		1696: {OutcomeMore, ""},
		1727: {OutcomeMore, ""},
		1730: {OutcomeMore, ""},
		1733: {OutcomeMore, ""},
		1736: {OutcomeMore, ""},
		1739: {OutcomeMore, ""},
		1793: {OutcomeMore, ""},
		1796: {OutcomeMore, ""},
		1799: {OutcomeMore, ""},
		1802: {OutcomeMore, ""},
		1805: {OutcomeMore, ""},

		931:  {OutcomeLess, ""},
		941:  {OutcomeLess, ""},
		1849: {OutcomeLess, ""},
		1697: {OutcomeLess, ""},
		1728: {OutcomeLess, ""},
		1731: {OutcomeLess, ""},
		1734: {OutcomeLess, ""},
		1737: {OutcomeLess, ""},
		1791: {OutcomeLess, ""},
		1794: {OutcomeLess, ""},
		1797: {OutcomeLess, ""},
		1800: {OutcomeLess, ""},
		1803: {OutcomeLess, ""},
		1806: {OutcomeLess, ""},
	}

	teamTotalsMappings = map[int64]OddMapping{
		974:  {OutcomeMore, "first"},
		1809: {OutcomeMore, "first"},
		1812: {OutcomeMore, "first"},
		1815: {OutcomeMore, "first"},
		1818: {OutcomeMore, "first"},
		1821: {OutcomeMore, "first"},
		1824: {OutcomeMore, "first"},
		1827: {OutcomeMore, "first"},
		1830: {OutcomeMore, "first"},
		2203: {OutcomeMore, "first"},

		976:  {OutcomeLess, "first"},
		1810: {OutcomeLess, "first"},
		1813: {OutcomeLess, "first"},
		1816: {OutcomeLess, "first"},
		1819: {OutcomeLess, "first"},
		1822: {OutcomeLess, "first"},
		1825: {OutcomeLess, "first"},
		1828: {OutcomeLess, "first"},
		1831: {OutcomeLess, "first"},
		2204: {OutcomeLess, "first"},

		978:  {OutcomeMore, "second"},
		1854: {OutcomeMore, "second"},
		1873: {OutcomeMore, "second"},
		1880: {OutcomeMore, "second"},
		1883: {OutcomeMore, "second"},
		1886: {OutcomeMore, "second"},
		1893: {OutcomeMore, "second"},
		1896: {OutcomeMore, "second"},
		1899: {OutcomeMore, "second"},
		2209: {OutcomeMore, "second"},

		980:  {OutcomeLess, "second"},
		1871: {OutcomeLess, "second"},
		1874: {OutcomeLess, "second"},
		1881: {OutcomeLess, "second"},
		1884: {OutcomeLess, "second"},
		1887: {OutcomeLess, "second"},
		1894: {OutcomeLess, "second"},
		1897: {OutcomeLess, "second"},
		1900: {OutcomeLess, "second"},
		2210: {OutcomeLess, "second"},
	}
)

func ensureMapEntry[T any](m map[string]*T, key string) {
	if _, ok := m[key]; !ok {
		m[key] = new(T)
	}
}

func normalizeLine(line int64) string {
	return strconv.FormatFloat(float64(line)/100, 'f', 1, 64)
}

func setWin1x2Value(win1x2 *entity.Win1x2Struct, oddType string, outcome entity.Factor) {
	switch oddType {
	case OutcomeWin1:
		win1x2.Win1 = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
	case OutcomeWinNone:
		win1x2.WinNone = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
	case OutcomeWin2:
		win1x2.Win2 = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
	}
}

func setTotalValue(total *entity.WinLessMore, oddType string, outcome entity.Factor) {
	switch oddType {
	case OutcomeMore:
		total.WinMore = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
	case OutcomeLess:
		total.WinLess = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
	}
}

func processHandicap(periods *[]entity.ResponsePeriod, periodIndex int, outcome entity.Factor) {
	line, _ := strconv.ParseFloat(normalizeLine(outcome.Line), 64)

	handicapMappings := map[int64]struct {
		line   float64
		isWin1 bool
	}{
		910:  {line, true},
		927:  {line, true},
		937:  {line, true},
		989:  {line, true},
		1569: {line, true},
		1672: {line, true},
		1677: {line, true},
		1680: {line, true},
		1683: {line, true},
		1686: {line, true},
		1689: {line, true},
		1692: {line, true},
		1723: {line, true},
		1845: {line, true},

		912:  {line, false},
		928:  {line, false},
		938:  {line, false},
		991:  {line, false},
		1572: {line, false},
		1675: {line, false},
		1678: {line, false},
		1681: {line, false},
		1684: {line, false},
		1687: {line, false},
		1690: {line, false},
		1718: {line, false},
		1724: {line, false},
		1846: {line, false},

		924: {0.5, true},
		925: {0.5, false},
	}

	if mapping, ok := handicapMappings[outcome.FactorId]; ok {
		lineStr := fmt.Sprintf("%.1f", mapping.line)
		ensureMapEntry((*periods)[periodIndex].Handicap, lineStr)

		if mapping.isWin1 {
			(*periods)[periodIndex].Handicap[lineStr].Win1 = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
		} else {
			(*periods)[periodIndex].Handicap[lineStr].Win2 = entity.Odd{Value: outcome.Odds, Raw: entity.RawOdds{FactorId: outcome.FactorId}}
		}
	}
}
