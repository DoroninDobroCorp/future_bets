package calculator

import (
	"livebets/results/internal/api/model"
	"livebets/results/internal/entity"
	"log"
	"strconv"
	"strings"
)

type OutcomeResult struct {
	Status string
	Amount float64
}

func GetOutcomeResult(bet *entity.LogBetAccept, fixture *model.Fixture, outcomeStr string, betSum, coef float64, sport string) OutcomeResult {
	log.Printf("[%s] Ставка: '%s' | Исход: '%s'\n", sport, bet.KeyOutcome, outcomeStr)

	pair := bet.Data["pair"].(map[string]interface{})
	first := pair["first"].(map[string]interface{})
	log.Printf("Матч (%s): %s-%s\n", first["matchId"], first["homeName"], first["awayName"])

	var result OutcomeResult

	if sport == "Soccer" {
		result = CalculateFootballOutcome(fixture, outcomeStr, betSum, coef)

	} else if sport == "Tennis" {
		result = CalculateTennisOutcome(fixture, outcomeStr, betSum, coef)

	} else if sport == "Basketball" {
		result = CalculateBasketballOutcome(fixture, outcomeStr, betSum, coef)
	}

	log.Printf("Результат ставки: %s (Amount: %f)\n", result.Status, result.Amount)
	return result
}

func CalculateOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	outcomeStr = strings.TrimSpace(outcomeStr)

	// Обработка периода
	if strings.HasPrefix(outcomeStr, "P") {
		return CalculateHalfOutcome(fixture, outcomeStr, betSum, coef)
	}

	// Обработка победы
	if outcomeStr == "1" || outcomeStr == "2" || outcomeStr == "X" {
		return Calculate1X2Outcome(fixture, outcomeStr, betSum, coef)
	}

	// Обработка гандикапов
	if strings.HasPrefix(outcomeStr, "H") {
		return CalculateHandicapOutcome(fixture, outcomeStr, betSum, coef)
	}

	// Обработка тоталов
	if strings.HasPrefix(outcomeStr, "T") {
		return CalculateTotalOutcome(fixture, outcomeStr, betSum, coef)
	}

	// Обработка индивидуальных тоталов
	if strings.HasPrefix(outcomeStr, "IT") {
		return CalculateIndividualTotalOutcome(fixture, outcomeStr, betSum, coef)
	}

	return OutcomeResult{Status: "unknown", Amount: 0}
}

func CalculateHalfOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
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

	return CalculateOutcome(halfFixture, outcomeStr, betSum, coef)
}

func Calculate1X2Outcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	home, away := fixture.FinalScore()

	switch outcomeStr {
	case "1":
		if home > away {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if home == away {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}

	case "2":
		if away > home {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if home == away {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}

	case "X":
		if home == away {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}

	default:
		return OutcomeResult{Status: "unknown", Amount: 0}
	}
}

func CalculateHandicapOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	side := outcomeStr[:2]
	handicapValue, err := strconv.ParseFloat(strings.TrimSpace(outcomeStr[2:]), 64)
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}

	homeScore, awayScore := fixture.FinalScore()

	if side == "H1" {
		adjustedScore := float64(homeScore) + handicapValue

		if adjustedScore > float64(awayScore) {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if adjustedScore == float64(awayScore) {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}

	} else if side == "H2" {
		adjustedScore := float64(awayScore) + handicapValue

		if adjustedScore > float64(homeScore) {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if adjustedScore == float64(homeScore) {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}
	}

	return OutcomeResult{Status: "unknown", Amount: 0}
}

func CalculateTotalOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	outcomeRest := strings.TrimSpace(outcomeStr[1:])
	operator := string(outcomeRest[0])
	totalThreshold, err := strconv.ParseFloat(strings.TrimSpace(outcomeRest[1:]), 64)
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}

	homeScore, awayScore := fixture.FinalScore()
	totalScore := float64(homeScore + awayScore)

	if operator == ">" {
		if totalScore > totalThreshold {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if totalScore == totalThreshold {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}

	} else if operator == "<" {
		if totalScore < totalThreshold {
			return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
		} else if totalScore == totalThreshold {
			return OutcomeResult{Status: "void", Amount: 0}
		} else {
			return OutcomeResult{Status: "lose", Amount: -betSum}
		}
	}

	return OutcomeResult{Status: "unknown", Amount: 0}
}

func CalculateIndividualTotalOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	side := outcomeStr[:3]
	rest := strings.TrimSpace(outcomeStr[3:])
	operator := string(rest[0])
	totalThreshold, err := strconv.ParseFloat(strings.TrimSpace(rest[1:]), 64)
	if err != nil {
		return OutcomeResult{Status: "error", Amount: 0}
	}

	homeScore, awayScore := fixture.FinalScore()

	if side == "IT1" {
		if operator == "<" {
			if float64(homeScore) < totalThreshold {
				return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
			} else if float64(homeScore) == totalThreshold {
				return OutcomeResult{Status: "void", Amount: 0}
			} else {
				return OutcomeResult{Status: "lose", Amount: -betSum}
			}

		} else if operator == ">" {
			if float64(homeScore) > totalThreshold {
				return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
			} else if float64(homeScore) == totalThreshold {
				return OutcomeResult{Status: "void", Amount: 0}
			} else {
				return OutcomeResult{Status: "lose", Amount: -betSum}
			}
		}

	} else if side == "IT2" {
		if operator == "<" {
			if float64(awayScore) < totalThreshold {
				return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
			} else if float64(awayScore) == totalThreshold {
				return OutcomeResult{Status: "void", Amount: 0}
			} else {
				return OutcomeResult{Status: "lose", Amount: -betSum}
			}

		} else if operator == ">" {
			if float64(awayScore) > totalThreshold {
				return OutcomeResult{Status: "win", Amount: betSum * (coef - 1)}
			} else if float64(awayScore) == totalThreshold {
				return OutcomeResult{Status: "void", Amount: 0}
			} else {
				return OutcomeResult{Status: "lose", Amount: -betSum}
			}
		}
	}

	return OutcomeResult{Status: "unknown", Amount: 0}
}
