package entity

import "time"

type MatchData struct {
	LeagueName string    `json:"leagueName"`
	HomeName   string    `json:"homeName"`
	AwayName   string    `json:"awayName"`
	MatchID    string    `json:"matchId"`
	Bookmaker  string    `json:"bookmaker"`
	SportName  string    `json:"sportName"`
	CreatedAt  time.Time `json:"createdAt"`
}
