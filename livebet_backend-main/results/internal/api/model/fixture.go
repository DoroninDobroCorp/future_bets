package model

// Period описывает период матча.
type Period struct {
	Number       int    `json:"number"`
	Status       int    `json:"status"`
	SettlementId int64  `json:"settlementId"`
	SettledAt    string `json:"settledAt"`
	Team1Score   int    `json:"team1Score"`
	Team2Score   int    `json:"team2Score"`
	// CancellationReason можно добавить при необходимости.
}

// Fixture описывает событие (матч) из Pinnacle API.
type Fixture struct {
	ID      int      `json:"id"`
	Periods []Period `json:"periods"`
}

// FinalScore возвращает финальный счёт матча.
// Сначала ищем период с Number == 0, если не найден – суммируем все периоды.
func (f *Fixture) FinalScore() (int, int) {
	for _, p := range f.Periods {
		if p.Number == 0 {
			return p.Team1Score, p.Team2Score
		}
	}
	home, away := 0, 0
	for _, p := range f.Periods {
		home += p.Team1Score
		away += p.Team2Score
	}
	return home, away
}

// ScoreByPeriod возвращает счёт для указанного периода и флаг, найден ли такой период.
func (f *Fixture) ScoreByPeriod(period int) (int, int, bool) {
	for _, p := range f.Periods {
		if p.Number == period {
			return p.Team1Score, p.Team2Score, true
		}
	}
	return 0, 0, false
}
