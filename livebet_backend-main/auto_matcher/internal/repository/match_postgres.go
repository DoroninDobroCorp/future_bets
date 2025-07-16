package repository

import (
	"context"
	"fmt"
	"livebets/auto_matcher/internal/entity"
	"livebets/auto_matcher/pkg/rdbms"

	"github.com/jackc/pgx/v5"
)

type MatchStorage interface {
	GetUnMatchedTeamsByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (teams []entity.UnMatchedTeamsByLeaguesPG, err error)
	GetMatchedTeamsByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (teams []entity.MatchedTeamsByLeaguesPG, err error)

	CheckTeams(ctx context.Context, firstTeamID, secondTeamID int64) (bool, error)
	CheckTeamsPair(ctx context.Context, firstTeamID, secondTeamID int64) (bool, error)
	InsertTeamsPair(ctx context.Context, firstTeamID, secondTeamID int64) error

	CheckLeagues(ctx context.Context, firstLeagueID, secondLeagueID int64) (bool, error)
	CheckLeaguesPair(ctx context.Context, firstLeagueID, secondLeagueID int64) (bool, error)
	InsertLeaguesPair(ctx context.Context, firstLeagueID, secondLeagueID int64) error

	GetUnMachedLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.League, err error)
	GetMatchedLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.LeagueMatchPG, err error)
	GetAllLeaguesByBookmaker(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.League, err error)

	GetBookmakers(ctx context.Context) (bookmakers []string, err error)
	GetSports(ctx context.Context) (sports []string, err error)

	GetUnMachedLeaguesByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string, inputLeagues, inputTeams []string) (leagues []entity.League, err error)
	GetUnMatchedTeamsByLeaguesByTeams(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string, inputLeagues, inputTeams []string) (teams []entity.UnMatchedTeamsByLeaguesPG, err error)
}

type MatchPGStorage struct {
	handler rdbms.Executor
}

func NewHandMatchPGStorage(handler rdbms.Executor) MatchStorage {
	return &MatchPGStorage{
		handler: handler,
	}
}

func (m *MatchPGStorage) GetBookmakers(ctx context.Context) (bookmakers []string, err error) {
	query := fmt.Sprintf(`SELECT DISTINCT bookmaker_name FROM %s;`, LeaguesTable)
	rows, err := m.handler.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookmaker string

		if err = rows.Scan(&bookmaker); err != nil {
			return nil, err
		}

		bookmakers = append(bookmakers, bookmaker)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetSports(ctx context.Context) (sports []string, err error) {
	query := fmt.Sprintf(`SELECT DISTINCT sport_name FROM %s;`, LeaguesTable)
	rows, err := m.handler.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sport string

		if err = rows.Scan(&sport); err != nil {
			return nil, err
		}

		sports = append(sports, sport)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetUnMachedLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.League, err error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT l.id, l.bookmaker_name, l.sport_name, l.league_name
		FROM %s AS l
        LEFT JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
        LEFT JOIN %s AS l2 ON l2.id = lm.league1_id OR l2.id = lm.league2_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3) 
    	AND ((l2.bookmaker_name <> $2 AND l2.bookmaker_name <> $3) OR lm.id IS NULL)
	`, LeaguesTable, LeaguesMergeTable, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var league entity.League

		err = rows.Scan(
			&league.ID,
			&league.BookmakerName,
			&league.SportName,
			&league.LeagueName,
		)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetMatchedLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.LeagueMatchPG, err error) {
	query := fmt.Sprintf(`
		SELECT l.id, l.bookmaker_name, l.sport_name, l.league_name, lm.id
		FROM %s AS l
		INNER JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3)
	`, LeaguesTable, LeaguesMergeTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var league entity.LeagueMatchPG

		err = rows.Scan(
			&league.ID,
			&league.BookmakerName,
			&league.SportName,
			&league.LeagueName,
			&league.LeagueMatchID,
		)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetAllLeaguesByBookmaker(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (leagues []entity.League, err error) {
	query := fmt.Sprintf(`
		SELECT l.id, l.bookmaker_name, l.sport_name, l.league_name FROM %s AS l
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3)
	`, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var league entity.League

		err = rows.Scan(
			&league.ID,
			&league.BookmakerName,
			&league.SportName,
			&league.LeagueName,
		)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetUnMatchedTeamsByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (teams []entity.UnMatchedTeamsByLeaguesPG, err error) {
	query := fmt.Sprintf(`
		SELECT l.id, l.bookmaker_name, l.sport_name, l.league_name, lm.id, t.id, t.team_name
		FROM %s AS l
		INNER JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
		INNER JOIN %s AS t ON l.id = t.league_id
		LEFT JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
        LEFT JOIN %s AS t2 ON t2.id = tm.team1_id OR t2.id = tm.team2_id
        LEFT JOIN %s AS l2 ON l2.id = t2.league_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3)
        AND ((l2.bookmaker_name <> $2 AND l2.bookmaker_name <> $3) OR tm.uuid IS NULL)
	`, LeaguesTable, LeaguesMergeTable, TeamsTable, TeamsMergeTable, TeamsTable, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var team entity.UnMatchedTeamsByLeaguesPG

		err = rows.Scan(
			&team.LeagueID,
			&team.BookmakerName,
			&team.SportName,
			&team.LeagueName,
			&team.LeagueMatchID,
			&team.TeamID,
			&team.TeamName,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetMatchedTeamsByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) (teams []entity.MatchedTeamsByLeaguesPG, err error) {
	query := fmt.Sprintf(`
		SELECT l.id, l.bookmaker_name, l.sport_name, l.league_name, lm.id, t.id, t.team_name, tm.uuid
		FROM %s AS l
		INNER JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
		INNER JOIN %s AS t ON l.id = t.league_id
		INNER JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3);
	`, LeaguesTable, LeaguesMergeTable, TeamsTable, TeamsMergeTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var team entity.MatchedTeamsByLeaguesPG

		err = rows.Scan(
			&team.LeagueID,
			&team.BookmakerName,
			&team.SportName,
			&team.LeagueName,
			&team.LeagueMatchID,
			&team.TeamID,
			&team.TeamName,
			&team.TeamMatch,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) CheckTeams(ctx context.Context, firstTeamID, secondTeamID int64) (bool, error) {
	query := fmt.Sprintf(`
		SELECT id FROM %s WHERE id IN ($1, $2)
	`, TeamsTable)

	rows, err := m.handler.Query(ctx, query, firstTeamID, secondTeamID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			return false, err
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if len(ids) != 2 {
		return false, nil
	}

	return true, nil
}

func (m *MatchPGStorage) CheckTeamsPair(ctx context.Context, firstTeamID, secondTeamID int64) (bool, error) {
	query := fmt.Sprintf(`
		SELECT uuid FROM %s
		WHERE (team1_id = $1 AND team2_id = $2) OR (team1_id = $2 AND team2_id = $1)
	`, TeamsMergeTable)

	rows, err := m.handler.Query(ctx, query, firstTeamID, secondTeamID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if counter == 0 {
		return false, nil
	}

	return true, nil
}

func (m *MatchPGStorage) InsertTeamsPair(ctx context.Context, firstTeamID, secondTeamID int64) error {
	query := fmt.Sprintf("INSERT INTO %s (team1_id, team2_id) VALUES ($1, $2)", TeamsMergeTable)
	_, err := m.handler.Exec(ctx, query, firstTeamID, secondTeamID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MatchPGStorage) CheckLeagues(ctx context.Context, firstLeagueID, secondLeagueID int64) (bool, error) {
	query := fmt.Sprintf(`
		SELECT id FROM %s WHERE id IN ($1, $2)
	`, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, firstLeagueID, secondLeagueID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			return false, err
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if len(ids) != 2 {
		return false, nil
	}

	return true, nil
}

func (m *MatchPGStorage) CheckLeaguesPair(ctx context.Context, firstLeagueID, secondLeagueID int64) (bool, error) {
	query := fmt.Sprintf(`
		SELECT id FROM %s
		WHERE (league1_id = $1 AND league2_id = $2) OR (league1_id = $2 AND league2_id = $1)
	`, LeaguesMergeTable)

	rows, err := m.handler.Query(ctx, query, firstLeagueID, secondLeagueID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if counter == 0 {
		return false, nil
	}

	return true, nil
}

func (m *MatchPGStorage) InsertLeaguesPair(ctx context.Context, firstLeagueID, secondLeagueID int64) error {
	query := fmt.Sprintf("INSERT INTO %s (league1_id, league2_id) VALUES ($1, $2)", LeaguesMergeTable)
	_, err := m.handler.Exec(ctx, query, firstLeagueID, secondLeagueID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MatchPGStorage) GetUnMachedLeaguesByLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string, inputLeagues, inputTeams []string) (leagues []entity.League, err error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT l.id, l.bookmaker_name, l.sport_name, l.league_name
		FROM %s AS l
		LEFT JOIN %s AS t ON l.id = t.league_id
		LEFT JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
		LEFT JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
		LEFT JOIN %s AS l2 ON l2.id = lm.league1_id OR l2.id = lm.league2_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3)
		AND (l2.bookmaker_name = $2 OR l2.bookmaker_name = $3 OR l2.bookmaker_name IS NULL)
		AND l.league_name = ANY ($4)
		AND t.team_name = ANY ($5)
		AND (l.league_name, t.team_name) NOT IN (
			SELECT DISTINCT l.league_name, t.team_name
			FROM %s AS l
			INNER JOIN %s AS t ON l.id = t.league_id
			INNER JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
			INNER JOIN %s AS t2 ON t2.id = tm.team1_id OR t2.id = tm.team2_id
			INNER JOIN %s AS l2 ON l2.id = t2.league_id
			WHERE l.sport_name = $1 AND ((l.bookmaker_name = $2 AND l2.bookmaker_name = $3) 
			OR (l.bookmaker_name = $3 AND l2.bookmaker_name = $2))
			AND l.league_name = ANY ($4)
			AND t.team_name = ANY ($5)
		)
	`, LeaguesTable, TeamsTable, TeamsMergeTable, LeaguesMergeTable, LeaguesTable,LeaguesTable, TeamsTable, TeamsMergeTable, TeamsTable, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName, inputLeagues, inputTeams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var league entity.League

		err = rows.Scan(
			&league.ID,
			&league.BookmakerName,
			&league.SportName,
			&league.LeagueName,
		)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func (m *MatchPGStorage) GetUnMatchedTeamsByLeaguesByTeams(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string, inputLeagues, inputTeams []string) (teams []entity.UnMatchedTeamsByLeaguesPG, err error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT l.id, l.bookmaker_name, l.sport_name, l.league_name, lm.id, t.id, t.team_name
		FROM %s AS l
		LEFT JOIN %s AS t ON l.id = t.league_id
		LEFT JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
		LEFT JOIN %s AS lm ON l.id = lm.league1_id OR l.id = lm.league2_id
		LEFT JOIN %s AS l2 ON l2.id = lm.league1_id OR l2.id = lm.league2_id
		WHERE l.sport_name = $1 AND (l.bookmaker_name = $2 OR l.bookmaker_name = $3)
		AND (l2.bookmaker_name = $2 OR l2.bookmaker_name = $3 OR l2.bookmaker_name IS NULL)
		AND (l2.bookmaker_name <> l.bookmaker_name OR l2.bookmaker_name IS NULL)
		AND l.league_name = ANY ($4)
		AND t.team_name = ANY ($5)
		AND (l.league_name, t.team_name) NOT IN (
			SELECT DISTINCT l.league_name, t.team_name
			FROM %s AS l
			INNER JOIN %s AS t ON l.id = t.league_id
			INNER JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
			INNER JOIN %s AS t2 ON t2.id = tm.team1_id OR t2.id = tm.team2_id
			INNER JOIN %s AS l2 ON l2.id = t2.league_id
			WHERE l.sport_name = $1 AND ((l.bookmaker_name = $2 AND l2.bookmaker_name = $3) 
			OR (l.bookmaker_name = $3 AND l2.bookmaker_name = $2))
			AND l.league_name = ANY ($4)
			AND t.team_name = ANY ($5)
		) AND lm.id IS NOT NULL
	`, LeaguesTable, TeamsTable, TeamsMergeTable, LeaguesMergeTable, LeaguesTable, LeaguesTable, TeamsTable, TeamsMergeTable, TeamsTable, LeaguesTable)

	rows, err := m.handler.Query(ctx, query, sportName, firstBookmakerName, secondBookmakerName, inputLeagues, inputTeams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var team entity.UnMatchedTeamsByLeaguesPG

		err = rows.Scan(
			&team.LeagueID,
			&team.BookmakerName,
			&team.SportName,
			&team.LeagueName,
			&team.LeagueMatchID,
			&team.TeamID,
			&team.TeamName,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}
