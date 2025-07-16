package entity

import (
	"time"
)

type LogBetAccept struct {
	KeyMatch    string                 `json:"key_match"`
	KeyOutcome  string                 `json:"key_outcome"`
	Data        map[string]interface{} `json:"data"`
	CorrectData map[string]interface{} `json:"correct_data"`
	Percent     float64                `json:"percent"`
	CreatedAt   time.Time              `json:"created_at"`
	EVProfit    *float64               `json:"ev_profit"`   // Ожидаемая прибыль (может быть NULL в БД)
	RealProfit  *float64               `json:"real_profit"` // Реальная прибыль (может быть NULL в БД)
}
