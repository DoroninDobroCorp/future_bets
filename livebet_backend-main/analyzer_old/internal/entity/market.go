package entity

import "time"

type MarketType struct {
	MatchScore     int       `json:"matchScore"`
	ChangedAt      time.Time `json:"changedAt"`
	IsChange       bool      `json:"isChange"`
	BookmakerScore float64   `json:"bookmakerScore"`
	CreatedAt      time.Time `json:"createdAt"`
	MarketType     int       `json:"marketType"`
	IsFallen       bool      `json:"isFallen"`
}
