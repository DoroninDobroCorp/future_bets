package priceStorage

import (
	"context"
	"livebets/analazer/cmd/config"
	"livebets/analazer/internal/entity"
	"sync"
	"time"
)

type PriceStorage struct {
	sync.RWMutex
	prices map[string]map[time.Time]entity.FullPriceRecord
}

func NewPriceStorage() *PriceStorage {
	return &PriceStorage{
		prices: make(map[string]map[time.Time]entity.FullPriceRecord),
	}
}

func (p *PriceStorage) Write(key string, timeR time.Time, record entity.FullPriceRecord) {
	p.Lock()
	defer p.Unlock()

	saveRecord := p.prices[key]
	if saveRecord == nil {
		p.prices[key] = make(map[time.Time]entity.FullPriceRecord)
	}

	p.prices[key][timeR] = record
}

func (p *PriceStorage) ReadByKey(key string) []entity.ResponsePriceRecord {
	p.RLock()
	defer p.RUnlock()

	priceV, ok := p.prices[key]
	if !ok {
		return nil
	}

	var records []entity.ResponsePriceRecord
	for time, record := range priceV {
		records = append(records, entity.ResponsePriceRecord{
			Key:       key,
			First:     record.First,
			Second:    record.Second,
			Outcome:   record.Outcome,
			ROI:       record.ROI,
			Margin:    record.Margin,
			CreatedAt: time,
		})
	}

	return records
}

func (p *PriceStorage) CleanByTimeout(ctx context.Context, cfg config.PriceStorage, wg *sync.WaitGroup) {
	defer wg.Done()

	cleanInterval := time.Duration(time.Duration(cfg.ClearInterval) * time.Second)
	cleanTicker := time.NewTicker(cleanInterval)

	for {
		select {
		case <-cleanTicker.C:
			p.Lock()

			for priceK, priceV := range p.prices {
				if len(priceV) == 0 {
					delete(p.prices, priceK)
				}
				for timeK := range priceV {
					if time.Since(timeK) > (time.Duration(cfg.DataTimeout) * time.Second) {
						delete(p.prices[priceK], timeK)
					}
				}
			}

			p.Unlock()
		case <-ctx.Done():
			cleanTicker.Stop()
			return
		}
	}
}
