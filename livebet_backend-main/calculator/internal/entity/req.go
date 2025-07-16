package entity

type AcceptBet struct {
	Pair   PairOneOutcome         `json:"pair"`
	Bet    CalculatedBetWithUsers `json:"bet"`
	Sum    float64                `json:"sum"`
	Coef   float64                `json:"coef"`
	Time   string                 `json:"time"`
	UserId int                    `json:"userId"`
}
