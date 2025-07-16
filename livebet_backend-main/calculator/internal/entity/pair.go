package entity

import "time"

type Outcome struct {
	Outcome    string  `json:"outcome"`
	ROI        float64 `json:"roi"`
	Margin     float64 `json:"margin"`
	Score1     Odd     `json:"score1"`
	Score2     Odd     `json:"score2"`
	MarketType int     `json:"marketType"`
}

type Odd struct {
	Value float64     `json:"value"`
	Raw   interface{} `json:"raw"`
}

type Match struct {
	Bookmaker  string    `json:"bookmaker"`
	LeagueName string    `json:"leagueName"`
	HomeScore  float64   `json:"homeScore"`
	AwayScore  float64   `json:"awayScore"`
	HomeName   string    `json:"homeName"`
	AwayName   string    `json:"awayName"`
	MatchID    string    `json:"matchId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Pair struct {
	First     Match     `json:"first"`
	Second    Match     `json:"second"`
	Outcome   []Outcome `json:"outcome"`
	IsLive    bool      `json:"isLive"`
	SportName string    `json:"sportName"`
	CreatedAt time.Time `json:"createdAt"`
}

type PairOneOutcome struct {
	First     Match     `json:"first"`
	Second    Match     `json:"second"`
	Outcome   Outcome   `json:"outcome"`
	IsLive    bool      `json:"isLive"`
	SportName string    `json:"sportName"`
	CreatedAt time.Time `json:"createdAt"`
}
