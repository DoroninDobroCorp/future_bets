package entity

type EventsData struct {
	LiveSports  []*LiveSport `json:"liveSports"`
	PageCount   int          `json:"pageCount"`
	Page        int          `json:"page"`
	Events      []*Event     `json:"events"`
	Competitors []*Item      `json:"competitors"`
}

type PreMatchData struct {
	PageCount  int               `json:"pageCount"`
	Page       int               `json:"page"`
	Markets    []*PreMatchMarket `json:"markets"`
	Odds       []*Odd            `json:"odds"`
	Events     []*Event          `json:"events"`
	Categories []*Item           `json:"categories"` // Country
	Leagues    []*Item           `json:"champs"`
	Teams      []*Item           `json:"competitors"`
}

type PreMatchMarket struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	OddIds []int64 `json:"oddIds"`
	TypeId int64   `json:"typeId"`
	Sv     string  `json:"sv"`
	IsMB   bool    `json:"isMB"`
}

type Market struct {
	Id     int64     `json:"id"`
	Name   string    `json:"name"`
	OddIds [][]int64 `json:"desktopOddIds"`
	TypeId int64     `json:"typeId"`
	Sv     string    `json:"sv"`
	IsMB   bool      `json:"isMB"`
	IsBB   bool      `json:"isBB"`
}

type Event struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     int64     `json:"status"`
	SportId    int64     `json:"sportId"`
	CatId      int64     `json:"campId"`        // Championship category (Country)
	LeagueId   int64     `json:"champId"`       // Championship
	TeamIds    []int64   `json:"competitorIds"` // Teams
	MarketsIds []int64   `json:"marketIds"`
	Scores     []float64 `json:"score"`
	IsBooked   bool      `json:"isBooked"`
	StartDate  string    `json:"startDate"`
}

type Item struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type LiveSport struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IconName string `json:"iconName"`
	Count    int64  `json:"count"`
}

type Match struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	LiveTime string    `json:"liveTime"`
	Sport    LiveSport `json:"sport"`
	League   Item      `json:"champ"`
	Category Item      `json:"category"`
	Teams    []*Item   `json:"competitors"`
	Markets  []*Market `json:"markets"`
	Odds     []*Odd    `json:"odds"`
	Scores   []float64
}

type Odd struct {
	TypeId    int64   `json:"typeId"`
	Price     float64 `json:"price"`
	IsMB      bool    `json:"isMB"`
	IsBB      bool    `json:"isBB"`
	Sv        string  `json:"sv"`
	OddStatus int64   `json:"oddStatus"`
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
}
