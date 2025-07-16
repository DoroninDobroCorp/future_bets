package entity

import "time"

const ParserName = "Sansabet"

type SportName string

const SportSoccer SportName = "Soccer"
const SportTennis SportName = "Tennis"

type ResponseGame struct {
	// Get from matches
	Pid        int64  `json:"Pid"`
	LeagueName string `json:"LeagueName"`
	HomeName   string `json:"homeName"`
	AwayName   string `json:"awayName"`
	MatchId    string `json:"MatchId"`
	IsLive     bool   `json:"isLive"`

	// Get from odds
	HomeScore float64          `json:"HomeScore"`
	AwayScore float64          `json:"AwayScore"`
	Periods   []ResponsePeriod `json:"Periods"`

	// Get from config
	Source    string    `json:"Source"`
	SportName SportName `json:"SportName"`
	CreatedAt time.Time `json:"CreatedAt"`

	Raw EventRaw `json:"Raw"`
}

type ResponsePeriod struct {
	Win1x2           Win1x2Struct             `json:"Win1x2"`
	Games            map[string]*Win1x2Struct `json:"Games"`
	Totals           map[string]*WinLessMore  `json:"Totals"`
	Handicap         map[string]*WinHandicap  `json:"Handicap"`
	FirstTeamTotals  map[string]*WinLessMore  `json:"FirstTeamTotals"`
	SecondTeamTotals map[string]*WinLessMore  `json:"SecondTeamTotals"`
}

type WinHandicap struct {
	Win1 OddValue `json:"Win1"`
	Win2 OddValue `json:"Win2"`
}

type WinLessMore struct {
	WinMore OddValue `json:"WinMore"`
	WinLess OddValue `json:"WinLess"`
}

type Win1x2Struct struct {
	Win1    OddValue `json:"Win1"`
	WinNone OddValue `json:"WinNone"`
	Win2    OddValue `json:"Win2"`
}

type OddValue struct {
	Value float64     `json:"value"`
	Raw   interface{} `json:"raw"`
}

type EventRaw struct {
	MatchName string `json:"match_name"`
}

type OddRaw struct {
	Line   string `json:"line"`
	BetNum int64  `json:"bet_num"`
}
