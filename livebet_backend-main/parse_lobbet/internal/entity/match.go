package entity

type ResponseMatchData struct {
	Live Live `json:"IMatchLiveContainer"`
}

type Live struct {
	Matches []*Match `json:"matches"`
}

type Match struct {
	ID          int64  `json:"id"`
	HomeTeam    string `json:"home"`
	AwayTeam    string `json:"away"`
	LeagueID    int64  `json:"leagueId"`
	LeagueName  string `json:"leagueName"`
	SportLetter string `json:"sport"`
	MatchResult Result `json:"matchResult"`
	Bets        []Bet  `json:"bets"`
	LiveStatus  int64  `json:"liveStatus"`
	TimeStamp   int64  `json:"kickOffTime"`
}

type Bet struct {
	LiveBetCaption string `json:"liveBetCaption"`
	Picks          []Pick `json:"picks"`
}

type Result struct {
	CurrentScore Score  `json:"currentScore"`
	ResultType   string `json:"resultType"`
}

type Score struct {
	Home float64 `json:"h"`
	Away float64 `json:"a"`
}

type Pick struct {
	Caption          string  `json:"caption"`
	OddValue         float64 `json:"oddValue"`
	LiveBetPickLabel string  `json:"liveBetPickLabel"`
	SpecialValue     string  `json:"specialValue"`
}

type ResponsePrematchData struct {
	Bets []PrematchBet `json:"odBetPickGroups"`
}

type PrematchBet struct {
	Description        string        `json:"description"`
	Tips               []PrematchTip `json:"tipTypes"`
	HandicapParamValue string        `json:"handicapParamValue"`
}

type PrematchTip struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
