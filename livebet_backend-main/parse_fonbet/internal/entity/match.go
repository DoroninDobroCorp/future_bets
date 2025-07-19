package entity

type Requested struct {
	Events        []Event     `json:"events"`
	EventMiscs    []EventMisc `json:"eventMiscs"`
	Leagues       []League    `json:"sports"`
	CustomFactors []EventOdds `json:"customFactors"`
}

type Event struct {
	Id       int64  `json:"id"`
	HomeName string `json:"team1"`
	AwayName string `json:"team2"`
	Place    string `json:"place"`
	SportId  int64  `json:"sportId"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentId"`
}

type EventMisc struct {
	Id     int64 `json:"id"`
	Score1 int64 `json:"score1"`
	Score2 int64 `json:"score2"`
}

type League struct {
	Id      int64  `json:"id"`
	SportId int64  `json:"parentId"`
	Name    string `json:"name"`
}

type EventOdds struct {
	EventId int64    `json:"e"`
	Factors []Factor `json:"factors"`
}

type Factor struct {
	FactorId   int64   `json:"f"`
	Odds       float64 `json:"v"`
	Line       int64   `json:"p"`
	LineString string  `json:"pt"`
}
