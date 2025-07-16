package entity

type League struct {
	ID            int64  `json:"id"`
	BookmakerName string `json:"bookmakerName"`
	SportName     string `json:"sportName"`
	LeagueName    string `json:"leagueName"`
}

type LeagueMatchPG struct {
	ID            int64  `json:"id"`
	BookmakerName string `json:"bookmakerName"`
	SportName     string `json:"sportName"`
	LeagueName    string `json:"leagueName"`
	LeagueMatchID int64  `json:"leaguesMatchID"`
}

type LeaguesMatchPair struct {
	LeagueIDFirst       int64  `json:"leagueIDFirst"`
	LeagueIDSecond      int64  `json:"leagueIDSecond"`
	BookmakerNameFirst  string `json:"bookmakerNameFirst"`
	BookmakerNameSecond string `json:"bookmakerNameSecond"`
	LeagueNameFirst     string `json:"leagueNameFirst"`
	LeagueNameSecond    string `json:"leagueNameSecond"`
	SportName           string `json:"sportName"`
}
