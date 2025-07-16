package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"livebets/calculator/internal/entity"
	"livebets/calculator/pkg/rdbms"

	"github.com/jackc/pgx/v5"
)

type LogsStorage interface {
	InsertLogBetAccept(ctx context.Context, keyMatch, keyOutcome string, pair entity.AcceptBet, priceRecord *entity.PriceRecord, percent float64, userId int, isLive bool, sport, bookmaker string) error
	GetInitializeCalcBet(ctx context.Context) (percents []entity.TotalPercentByKey, err error)
	InsertLogTestBetAccept(ctx context.Context, keyMatch, keyOutcome string, pair entity.AcceptBet, priceRecord *entity.PriceRecord, percent float64) error
}

type LogsPGStorage struct {
	handler rdbms.Executor
}

func NewLogsPGStorage(handler rdbms.Executor) LogsStorage {
	return &LogsPGStorage{
		handler: handler,
	}
}

func (l *LogsPGStorage) InsertLogBetAccept(ctx context.Context, keyMatch, keyOutcome string, pair entity.AcceptBet, priceRecord *entity.PriceRecord, percent float64, userId int, isLive bool, sport, bookmaker string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (key_match, key_outcome, data, correct_data, percent, user_id, is_live, sport, bookmaker) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, LogBetAccept)

	bytePair, err := json.Marshal(pair)
	if err != nil {
		return err
	}

	bytePriceRecord, err := json.Marshal(priceRecord)
	if err != nil {
		return err
	}

	if _, err := l.handler.Exec(ctx, query, keyMatch, keyOutcome, bytePair, bytePriceRecord, percent, userId, isLive, sport, bookmaker); err != nil {
		return err
	}

	return nil
}

func (l *LogsPGStorage) InsertLogTestBetAccept(ctx context.Context, keyMatch, keyOutcome string, pair entity.AcceptBet, priceRecord *entity.PriceRecord, percent float64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (key_match, key_outcome, data, correct_data, percent) VALUES ($1, $2, $3, $4, $5)
	`, LogTestBetAccept)

	bytePair, err := json.Marshal(pair)
	if err != nil {
		return err
	}

	bytePriceRecord, err := json.Marshal(priceRecord)
	if err != nil {
		return err
	}

	if _, err := l.handler.Exec(ctx, query, keyMatch, keyOutcome, bytePair, bytePriceRecord, percent); err != nil {
		return err
	}

	return nil
}

func (l *LogsPGStorage) GetInitializeCalcBet(ctx context.Context) (percents []entity.TotalPercentByKey, err error) {
	query := fmt.Sprintf(`
		SELECT key_match, sum(percent) as totalPercent FROM %s 
		WHERE created_at >= NOW() - INTERVAL '2 HOURS'
		GROUP BY key_match
	`, LogBetAccept)

	rows, err := l.handler.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var percent entity.TotalPercentByKey

		if err = rows.Scan(&percent.KeyMatch, &percent.TotalPercent); err != nil {
			return nil, err
		}

		percents = append(percents, percent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}
