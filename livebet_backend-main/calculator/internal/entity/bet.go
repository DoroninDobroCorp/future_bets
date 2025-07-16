package entity

import (
	"time"
)

type CalculatedBet struct {
	OriginalAmount float64 `json:"originalAmount"`
	AdjustedAmount float64 `json:"adjustedAmount"`
	Percentage     float64 `json:"percentage"`
}

type CalculatedBetWithUsers struct {
	CalcBet    CalculatedBet `json:"calcBet"`
	UsersCount int           `json:"usersCount"`
}

type TotalPercent struct {
	TotalPercent float64   `json:"totalPercent"`
	CreatedAt    time.Time `json:"createdAt"`
}

type TotalPercentByKey struct {
	KeyMatch     string  `json:"keyMatch"`
	TotalPercent float64 `json:"totalPercent"`
}

type BetFile struct {
	MatchName string  `json:"match_name"`
	BetType   string  `json:"bet_type"`
	Amount    float64 `json:"amount"`
	Odds      float64 `json:"odds"`
}

type MissedBet struct {
	KeyMatch string         `json:"key_match"`
	Pair     PairOneOutcome `json:"pair"`
}
