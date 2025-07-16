package marketStorage

import (
	"context"
	"livebets/analazer/cmd/config"
	"livebets/analazer/internal/entity"
	"sync"
	"time"
)

type MarketStorage struct {
	sync.RWMutex
	markets map[string]entity.MarketType
}

func NewMarketStorage() *MarketStorage {
	return &MarketStorage{
		markets: make(map[string]entity.MarketType),
	}
}

// Write data and get type of market
func (m *MarketStorage) Write(key string, matchScore int, bookmakerScore float64, CreatedAt time.Time) int {
	m.Lock()
	defer m.Unlock()

	marketType := 0

	// If market not exist, type of market is norlmally
	oldMarket, ok := m.markets[key]
	if !ok {
		m.markets[key] = entity.MarketType{
			MatchScore:     -1,
			BookmakerScore: bookmakerScore,
			IsChange:       false,
			ChangedAt:      CreatedAt,
			CreatedAt:      CreatedAt,
			MarketType:     marketType,
			IsFallen:       false,
		}
		return marketType
	}

	// If changed match score, type of market is normally
	if oldMarket.MatchScore != matchScore {
		m.markets[key] = entity.MarketType{
			MatchScore:     matchScore,
			ChangedAt:      CreatedAt,
			IsChange:       true,
			BookmakerScore: bookmakerScore,
			CreatedAt:      CreatedAt,
			MarketType:     marketType,
			IsFallen:       false,
		}
		return marketType
	}

	// Check falling market after match score changed
	if oldMarket.IsChange {

		if oldMarket.BookmakerScore < bookmakerScore {
			m.markets[key] = entity.MarketType{
				MatchScore:     matchScore,
				ChangedAt:      oldMarket.ChangedAt,
				IsChange:       false,
				BookmakerScore: bookmakerScore,
				CreatedAt:      CreatedAt,
				MarketType:     marketType,
				IsFallen:       false,
			}
			return marketType
		}

		// Set falling type for market
		if time.Since(oldMarket.ChangedAt) > (120*time.Second) && oldMarket.IsFallen {
			marketType = -1
			m.markets[key] = entity.MarketType{
				MatchScore:     matchScore,
				ChangedAt:      oldMarket.ChangedAt,
				IsChange:       false,
				BookmakerScore: bookmakerScore,
				CreatedAt:      CreatedAt,
				MarketType:     marketType,
				IsFallen:       false,
			}
			return marketType
		}
	}

	isFallen := oldMarket.IsFallen
	if oldMarket.BookmakerScore > bookmakerScore {
		isFallen = true
	}

	m.markets[key] = entity.MarketType{
		MatchScore:     matchScore,
		ChangedAt:      oldMarket.ChangedAt,
		IsChange:       oldMarket.IsChange,
		BookmakerScore: bookmakerScore,
		CreatedAt:      CreatedAt,
		MarketType:     oldMarket.MarketType,
		IsFallen:       isFallen,
	}

	return oldMarket.MarketType
}

func (m *MarketStorage) CleanByTimeout(ctx context.Context, cfg config.MarketStorage, wg *sync.WaitGroup) {
	defer wg.Done()

	cleanInterval := time.Duration(time.Duration(cfg.ClearInterval) * time.Second)
	cleanTicker := time.NewTicker(cleanInterval)

	for {
		select {
		case <-cleanTicker.C:
			m.Lock()

			for marketK, marketV := range m.markets {
				if time.Since(marketV.CreatedAt) > (time.Duration(cfg.DataTimeout) * time.Second) {
					delete(m.markets, marketK)
				}
			}

			m.Unlock()
		case <-ctx.Done():
			cleanTicker.Stop()
			return
		}
	}
}
