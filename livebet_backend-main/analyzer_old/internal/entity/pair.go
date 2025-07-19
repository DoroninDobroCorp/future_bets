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

type ResponseMatch struct {
	Bookmaker  string      `json:"bookmaker"`
	LeagueName string      `json:"leagueName"`
	HomeScore  float64     `json:"homeScore"`
	AwayScore  float64     `json:"awayScore"`
	HomeName   string      `json:"homeName"`
	AwayName   string      `json:"awayName"`
	MatchID    string      `json:"matchId"`
	CreatedAt  time.Time   `json:"createdAt"`
	Raw        interface{} `jsom:"raw"`
}

type ResponsePair struct {
	First     ResponseMatch `json:"first"`
	Second    ResponseMatch `json:"second"`
	Outcome   []Outcome     `json:"outcome"`
	SportName string        `json:"sportName"`
	IsLive    bool          `json:"isLive"`
	CreatedAt time.Time     `json:"createdAt"`
}
