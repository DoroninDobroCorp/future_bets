package calculator

import (
	"livebets/results/internal/api/model"
)

func CalculateFootballOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	return CalculateOutcome(fixture, outcomeStr, betSum, coef)
}
