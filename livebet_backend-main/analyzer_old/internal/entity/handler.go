package entity

type ReqGetPriceRecordsByTime struct {
	Bookmaker1 string `json:"bookmaker1"`
	Bookmaker2 string `json:"bookmaker2"`
	MatchID1   string `json:"matchID1"`
	MatchID2   string `json:"matchID2"`
	SportName  string `json:"sportName"`
	Outcome    string `json:"outcome"`

	Minutes  int   `json:"minutes"`
	Seconds  int   `json:"seconds"`
	LongTime int64 `json:"longTime"`
}
