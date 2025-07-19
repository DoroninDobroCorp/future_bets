package service

import (
	"fmt"
	"livebets/analazer/internal/entity"
	"livebets/shared"
	"sort"
	"strconv"
	"strings"
)

// Расчет MARGIN
func calculateMARGIN(outcomeName string, outcomes map[string]entity.OddsWithMarket) float64 {
	def := 1.08

	parallelOdds := []float64{}
	paralelNames := []string{outcomeName}

	splitedName := strings.Split(outcomeName, " ")

	switch len(splitedName) {
	case 1:
		names := []string{
			"1", "X", "2",
		}
		for _, name := range names {
			if name != paralelNames[0] {
				paralelNames = append(paralelNames, name)
			}
		}
	case 2:
		switch splitedName[0] {
		case "T>":
			paralelNames = append(paralelNames, fmt.Sprintf("T< %s", splitedName[1]))
		case "T<":
			paralelNames = append(paralelNames, fmt.Sprintf("T> %s", splitedName[1]))
		case "IT1>":
			paralelNames = append(paralelNames, fmt.Sprintf("IT1< %s", splitedName[1]))
		case "IT1<":
			paralelNames = append(paralelNames, fmt.Sprintf("IT1> %s", splitedName[1]))
		case "IT2>":
			paralelNames = append(paralelNames, fmt.Sprintf("IT2< %s", splitedName[1]))
		case "IT2<":
			paralelNames = append(paralelNames, fmt.Sprintf("IT2> %s", splitedName[1]))
		case "H1":
			line, _ := strconv.ParseFloat(splitedName[1], 64)
			paralelNames = append(paralelNames, fmt.Sprintf("H2 %.1f", line*-1))
		case "H2":
			line, _ := strconv.ParseFloat(splitedName[1], 64)
			paralelNames = append(paralelNames, fmt.Sprintf("H1 %.1f", line*-1))
		}
	case 3:
		switch splitedName[1] {
		case "T>":
			paralelNames = append(paralelNames, fmt.Sprintf("%s T< %s", splitedName[0], splitedName[2]))
		case "T<":
			paralelNames = append(paralelNames, fmt.Sprintf("%s T> %s", splitedName[0], splitedName[2]))
		case "IT1>":
			paralelNames = append(paralelNames, fmt.Sprintf("%s IT1< %s", splitedName[0], splitedName[2]))
		case "IT1<":
			paralelNames = append(paralelNames, fmt.Sprintf("%s IT1> %s", splitedName[0], splitedName[2]))
		case "IT2>":
			paralelNames = append(paralelNames, fmt.Sprintf("%s IT2< %s", splitedName[0], splitedName[2]))
		case "IT2<":
			paralelNames = append(paralelNames, fmt.Sprintf("%s IT2> %s", splitedName[0], splitedName[2]))
		case "H1":
			line, _ := strconv.ParseFloat(splitedName[2], 64)
			paralelNames = append(paralelNames, fmt.Sprintf("%s H2 %.1f", splitedName[0], line*-1))
		case "H2":
			line, _ := strconv.ParseFloat(splitedName[2], 64)
			paralelNames = append(paralelNames, fmt.Sprintf("%s H1 %.1f", splitedName[0], line*-1))
		case "1G":
			paralelNames = append(paralelNames, fmt.Sprintf("%s 2G %s", splitedName[0], splitedName[2]))
		case "2G":
			paralelNames = append(paralelNames, fmt.Sprintf("%s 1G %s", splitedName[0], splitedName[2]))
		}
	}

	for _, paparalelName := range paralelNames {
		odds, ok := outcomes[paparalelName]
		if ok {
			parallelOdds = append(parallelOdds, odds.Odds[1].Value) // PINNACLE = 1
		} else {
			return def
		}
	}

	sum := 0.0
	for _, odd := range parallelOdds {
		if odd != 0 {
			sum += 1.0 / odd
		} else {
			return def
		}
	}

	return sum
}

// Расчет и фильтрация общих исходов
func (p *PairsMatchingService) calculateAndFilterCommonOutcomes(commonOutcomes map[string]entity.OddsWithMarket, secondBookmakerName, sportName string) []entity.Outcome {
	var filtered []entity.Outcome

	for outcome, values := range commonOutcomes {

		margin := calculateMARGIN(outcome, commonOutcomes)
		roi := shared.CalculateROI(values.Odds[0].Value, values.Odds[1].Value, margin, values.MarketType, shared.Parser(secondBookmakerName), shared.SportName(sportName))
		if roi > 0 {
			filtered = append(filtered, entity.Outcome{
				Outcome:    outcome,
				ROI:        roi,
				Margin:     margin,
				Score1:     values.Odds[1], //PINNACLE
				Score2:     values.Odds[0],
				MarketType: values.MarketType,
			})
		}
	}
	// Сортируем результаты по убыванию ROI
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ROI > filtered[j].ROI
	})

	return filtered
}

var (
	timePrefixes = []string{
		"Time1 ",
		"Time2 ",
	}
	betTypeMapping = map[string]string{
		"home": "1",
		"draw": "X",
		"away": "2",
	}
	totalMapping = map[string]string{
		"Total More ":             "T>",
		"Total Less ":             "T<",
		"First Team Total More ":  "IT1>",
		"First Team Total Less ":  "IT1<",
		"Second Team Total More ": "IT2>",
		"Second Team Total Less ": "IT2<",
	}
)

func normalizeLine(total string) string {
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return total
	}
	return fmt.Sprintf("%.1f", value)
}

func processHandicap(betType, prefix string) string {
	parts := strings.SplitN(strings.TrimPrefix(betType, prefix+"Handicap "), " ", 2)
	if len(parts) != 2 {
		return betType
	}

	handicap := normalizeLine(parts[0])
	switch parts[1] {
	case "Win1":
		return prefix + "H1" + handicap
	case "Win2":
		return prefix + "H2" + handicap
	default:
		return betType
	}
}

func processTotal(betType, timePrefix string) string {
	for pattern, replacement := range totalMapping {
		if strings.HasPrefix(betType, timePrefix+pattern) {
			total := strings.TrimPrefix(betType, timePrefix+pattern)
			return timePrefix + replacement + normalizeLine(total)
		}
	}
	return betType
}

func NormalizeBetType(betType string) string {
	if normalized, ok := betTypeMapping[betType]; ok {
		return normalized
	}

	for _, timePrefix := range timePrefixes {
		if strings.HasPrefix(betType, timePrefix) {
			if normalized, ok := betTypeMapping[strings.ReplaceAll(betType, timePrefix, "")]; ok {
				return fmt.Sprintf("%s%s", timePrefix, normalized)
			}

			if strings.Contains(betType, "Handicap") {
				return processHandicap(betType, timePrefix)
			}

			return processTotal(betType, timePrefix)
		}
	}

	if strings.Contains(betType, "Handicap") {
		return processHandicap(betType, "")
	}

	return processTotal(betType, "")
}
