package repository

import (
	"context"
	"fmt"
	"livebets/analazer/pkg/rdbms"

	"github.com/jackc/pgx/v5"
)

type PairsMatchingStorage interface {
	InsertLeague(ctx context.Context, bookmaker, sport, league string) (*int64, error)
	InsertTeam(ctx context.Context, leagueID int64, team string) error
	GetUUIDKeys(ctx context.Context, bookmaker, sport, league, homeTeam, awayTeam string) (uuids []string, err error)
}

type PairsMatchingPGStorage struct {
	handler rdbms.Executor
}

func NewPairsMatchingPGStorage(handler rdbms.Executor) PairsMatchingStorage {
	return &PairsMatchingPGStorage{
		handler: handler,
	}
}

func (p *PairsMatchingPGStorage) InsertLeague(ctx context.Context, bookmaker, sport, league string) (*int64, error) {
	query := fmt.Sprintf(`WITH new_league AS (
		INSERT INTO %s(bookmaker_name, sport_name, league_name) VALUES ($1, $2, $3)
		ON CONFLICT("bookmaker_name", "sport_name", "league_name") DO NOTHING
		RETURNING id
	) SELECT COALESCE(
		(SELECT id FROM new_league),
		(SELECT id FROM %s WHERE bookmaker_name = $1 AND sport_name = $2 AND league_name = $3)
	) result;`, LeaguesTable, LeaguesTable)

	var id *int64
	row := p.handler.QueryRow(ctx, query, bookmaker, sport, league)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return id, nil
}

func (p *PairsMatchingPGStorage) InsertTeam(ctx context.Context, leagueID int64, team string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s(league_id, team_name) VALUES ($1, $2)
		ON CONFLICT("league_id", "team_name") DO NOTHING;`, TeamsTable)

	if _, err := p.handler.Exec(ctx, query, leagueID, team); err != nil {
		return err
	}

	return nil
}

func (p *PairsMatchingPGStorage) GetUUIDKeys(ctx context.Context, bookmaker, sport, league, homeTeam, awayTeam string) (uuids []string, err error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT tm.uuid FROM %s AS l
		INNER JOIN %s AS t ON l.id = t.league_id
		INNER JOIN %s AS tm ON t.id = tm.team1_id OR t.id = tm.team2_id
		WHERE l.bookmaker_name = $1 AND l.sport_name = $2 AND l.league_name = $3 
		AND (t.team_name = $4 OR t.team_name = $5)
	`, LeaguesTable, TeamsTable, TeamsMergeTable)

	rows, err := p.handler.Query(ctx, query, bookmaker, sport, league, homeTeam, awayTeam)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var uuid string

		if err = rows.Scan(&uuid); err != nil {
			return nil, err
		}

		uuids = append(uuids, uuid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}
