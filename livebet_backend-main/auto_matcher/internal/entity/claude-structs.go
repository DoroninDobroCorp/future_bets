package entity

type RequestLeagues struct {
	BK1Leagues []League `json:"BK1_leagues"`
	BK2Leagues []League `json:"BK2_leagues"`
}

type ResponsePairLeague struct {
	BK1LeagueID int64 `json:"BK1_league_id"`
	BK2LeagueID int64 `json:"BK2_league_id"`
}

type RequestTeams struct {
	BK1Teams []UnMatchedTeam `json:"BK1_teams"`
	BK2Teams []UnMatchedTeam `json:"BK2_teams"`
}

type ResponsePairTeam struct {
	BK1TeamID int64 `json:"BK1_team_id"`
	BK2TeamID int64 `json:"BK2_team_id"`
}
