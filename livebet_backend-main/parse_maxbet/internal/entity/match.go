package entity

type RequestedData struct {
	LiveEvents  []Event       `json:"liveHeaders"`
	LiveBets    []Bet         `json:"liveBets"`
	LiveResults []EventResult `json:"liveResults"`
}

type Event struct {
	Id         int64  `json:"id"`
	HomeName   string `json:"h"`
	AwayName   string `json:"a"`
	LeagueName string `json:"lg"`
	SportId    string `json:"s"`
}

type Bet struct {
	EventId int64           `json:"mId"`
	Line    string          `json:"sv"`
	Coefs   map[string]Coef `json:"om"`
}

type Coef struct {
	Value float64 `json:"ov"`
}

type EventResult struct {
	EventId   int64 `json:"mi"`
	HomeScore Score `json:"hs"`
	AwayScore Score `json:"as"`
}

type Score struct {
	FullTime   int `json:"FULLTIME"`
	FirstHalf  int `json:"FIRST_HALF"`
	SecondHalf int `json:"SECOND_HALF"`
}
