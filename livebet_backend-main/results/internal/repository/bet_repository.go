package repository

import (
	"encoding/json"
	"time"

	"livebets/results/internal/entity"
)

func (p *PostgresClient) GetYesterdayBets() ([]*entity.LogBetAccept, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	query := `
        SELECT key_match, key_outcome, data, correct_data, percent, created_at, ev_profit, real_profit
        FROM Calculator.log_bet_accept
        WHERE DATE(created_at) = $1`
	rows, err := p.DB.Query(query, yesterday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []*entity.LogBetAccept
	for rows.Next() {
		var bet entity.LogBetAccept
		var dataBytes, correctDataBytes []byte

		if err := rows.Scan(&bet.KeyMatch, &bet.KeyOutcome, &dataBytes, &correctDataBytes, &bet.Percent, &bet.CreatedAt, &bet.EVProfit, &bet.RealProfit); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dataBytes, &bet.Data); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(correctDataBytes, &bet.CorrectData); err != nil {
			return nil, err
		}
		bets = append(bets, &bet)
	}
	return bets, nil
}

func (p *PostgresClient) UpdateBetProfits(keyOutcome string, evProfit, realProfit float64) error {
	query := `
		UPDATE Calculator.log_bet_accept
		SET ev_profit = $2, real_profit = $3
		WHERE key_outcome = $1
	`
	_, err := p.DB.Exec(query, keyOutcome, evProfit, realProfit)
	return err
}

func (p *PostgresClient) UpdateBetTime(keyOutcome string, newTime time.Time) error {
	query := `
		UPDATE Calculator.log_bet_accept
		SET created_at = $3
		WHERE key_outcome = $1
	`
	_, err := p.DB.Exec(query, keyOutcome, newTime)
	return err
}

func (p *PostgresClient) GetYesterdayTestBets() ([]*entity.LogBetAccept, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	query := `
        SELECT key_match, key_outcome, data, correct_data, percent, created_at, ev_profit, real_profit
        FROM Calculator.log_test_bet_accept
        WHERE DATE(created_at) = $1
    `
	rows, err := p.DB.Query(query, yesterday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []*entity.LogBetAccept
	for rows.Next() {
		var bet entity.LogBetAccept
		var dataBytes, correctDataBytes []byte

		if err := rows.Scan(&bet.KeyMatch, &bet.KeyOutcome, &dataBytes, &correctDataBytes, &bet.Percent, &bet.CreatedAt, &bet.EVProfit, &bet.RealProfit); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dataBytes, &bet.Data); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(correctDataBytes, &bet.CorrectData); err != nil {
			return nil, err
		}
		bets = append(bets, &bet)
	}
	return bets, nil
}

func (p *PostgresClient) UpdateTestBetProfits(keyOutcome string, evProfit, realProfit float64) error {
	query := `
		UPDATE Calculator.log_test_bet_accept
		SET ev_profit = $2, real_profit = $3
		WHERE key_outcome = $1
	`
	_, err := p.DB.Exec(query, keyOutcome, evProfit, realProfit)
	return err
}

func (p *PostgresClient) GetTestBets() ([]entity.LogBetAccept, error) {
	query := "SELECT key_match, key_outcome, created_at, data, ev_profit, real_profit FROM Calculator.log_test_bet_accept"
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []entity.LogBetAccept
	for rows.Next() {
		var bet entity.LogBetAccept
		var dataBytes []byte

		if err := rows.Scan(&bet.KeyMatch, &bet.KeyOutcome, &bet.CreatedAt, &dataBytes, &bet.EVProfit, &bet.RealProfit); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dataBytes, &bet.Data); err != nil {
		}

		pair := bet.Data["pair"].(map[string]interface{})
		outcome := pair["outcome"].(map[string]interface{})

		if bet.RealProfit == nil || outcome["roi"].(float64) > 12 || !pair["isLive"].(bool) {
			continue
		}

		bets = append(bets, bet)
	}
	return bets, nil
}

func (p *PostgresClient) FixTestDB() error {
	bets, err := p.GetTestBets()
	if err != nil {
		return err
	}

	query := "DELETE FROM Calculator.log_test_bet_accept WHERE key_outcome = $1"

	for _, bet := range bets {
		pair := bet.Data["pair"].(map[string]interface{})

		if pair["sportName"].(string) == "Tennis" {
			_, err := p.DB.Exec(query, bet.KeyOutcome)
			if err != nil {
				return err
			}
		}
	}

	return err
}

type BetRepository interface {
	GetYesterdayBets() ([]*entity.LogBetAccept, error)
	UpdateBetProfits(keyOutcome string, evProfit, realProfit float64) error
	UpdateBetTime(keyOutcome string, newTime time.Time) error
	GetYesterdayTestBets() ([]*entity.LogBetAccept, error)
	UpdateTestBetProfits(keyOutcome string, evProfit, realProfit float64) error
	GetTestBets() ([]entity.LogBetAccept, error)
}
