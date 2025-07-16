package calculator

import (
	"livebets/results/internal/api/model"
	"log"
	"strings"
)

var periodBasketballMapping = map[string]string{
	"P1": "P3",
	"P2": "P4",
	"P3": "P5",
	"P4": "P6",
	"P5": "P1",
	"P6": "P2",
}

func CalculateBasketballOutcome(fixture *model.Fixture, outcomeStr string, betSum, coef float64) OutcomeResult {
	for k, v := range periodBasketballMapping {
		if strings.Contains(outcomeStr, k) {
			log.Printf("Period %s -> %s", k, v)
			outcomeStr = strings.ReplaceAll(outcomeStr, k, v)
			break
		}
	}

	return CalculateOutcome(fixture, outcomeStr, betSum, coef)
}
