package calculator

import (
	"livebets/results/internal/api/model"
	"strconv"
	"strings"
)

func CalculateTennisOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	outcomeStr = strings.TrimSpace(outcomeStr)

	// Обработка геймов
	if strings.Contains(outcomeStr, "G") {
		return CalculateTennisGameOutcome(fixture, outcomeStr, betSum, coef)
	}

	// обработка периодов
	if strings.HasPrefix(outcomeStr, "P") {
		return CalculateTennisHalfOutcome(fixture, outcomeStr, betSum, coef)
	}

	// Обработка 1X2
	if outcomeStr == "1" || outcomeStr == "2" {
		return CalculateTennis1X2Outcome(fixture, outcomeStr, betSum, coef)
	}

	return OutcomeResult{Status: "unknown", Amount: 0}
}

func CalculateTennisHalfOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	halfPrefix := outcomeStr[:2]
	outcomeStr = strings.TrimSpace(outcomeStr[2:])

	period, err := strconv.Atoi(strings.ReplaceAll(halfPrefix, "P", ""))
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}

	homeScore, awayScore, found := fixture.ScoreByPeriod(period)
	if !found {
		return OutcomeResult{Status: "unknown", Amount: 0}
	}

	halfFixture := &model.Fixture{
		Periods: []model.Period{
			{
				Number:     0,
				Team1Score: homeScore,
				Team2Score: awayScore,
			},
		},
	}

	return CalculateTennisOutcome(halfFixture, outcomeStr, betSum, coef)
}

func CalculateTennisGameOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	splitted := strings.Split(outcomeStr, " ")

	periodNumber, err := strconv.Atoi(strings.ReplaceAll(splitted[0], "P", ""))
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}
	gameNumber, err := strconv.Atoi(splitted[2])
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}
	team := strings.ReplaceAll(splitted[1], "G", "")

	period := (periodNumber-1)*13 + (gameNumber - 1) + 6

	home, away, found := fixture.ScoreByPeriod(period)
	if !found {
		return OutcomeResult{Status: "unknown", Amount: 0}
	}

	switch team {
	case "1":
		if home > away {
			return OutcomeResult{Status: "win", Amount: betSum * coef}
		} else {
			return OutcomeResult{Status: "lose", Amount: betSum * coef}
		}

	case "2":
		if away > home {
			return OutcomeResult{Status: "win", Amount: betSum * coef}
		} else {
			return OutcomeResult{Status: "lose", Amount: betSum * coef}
		}

	default:
		return OutcomeResult{Status: "unknown", Amount: 0}
	}
}

func CalculateTennis1X2Outcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	home, away := fixture.FinalScore()

	switch outcomeStr {
	case "1":
		if home > away {
			return OutcomeResult{Status: "win", Amount: betSum * coef}
		} else {
			return OutcomeResult{Status: "lose", Amount: betSum * coef}
		}

	case "2":
		if away > home {
			return OutcomeResult{Status: "win", Amount: betSum * coef}
		} else {
			return OutcomeResult{Status: "lose", Amount: betSum * coef}
		}

	default:
		return OutcomeResult{Status: "unknown", Amount: 0}
	}
}
