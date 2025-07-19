package service

import (
	"livebets/analazer/internal/entity"
	priceStorage "livebets/analazer/internal/price-storage"
	"livebets/analazer/pkg/utils"
	"math"
	"sort"

	"github.com/rs/zerolog"
)

type PriceService struct {
	priceStorage *priceStorage.PriceStorage
	logger       *zerolog.Logger
}

func NewPriceService(
	priceStorage *priceStorage.PriceStorage,
	logger *zerolog.Logger,
) *PriceService {
	return &PriceService{
		priceStorage: priceStorage,
		logger:       logger,
	}
}

func (p *PriceService) GetPriceRecordsByTime(data entity.ReqGetPriceRecordsByTime) (int, []entity.ResponsePriceRecord) {
	key := utils.GenerateFullMatchKey(data.Bookmaker1, data.Bookmaker2, data.MatchID1, data.MatchID2, data.SportName, data.Outcome)
	var iSave int = -1

	// Get prices from storage
	records := p.priceStorage.ReadByKey(key)
	if records == nil {
		return iSave, nil
	}

	// Sorting by time
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})

	secs := data.Minutes * 60 + data.Seconds

	// Find near time
	var minFirst int = math.MaxInt32
	var minSecond int = math.MaxInt32
	var minPair int = math.MaxInt32
	for i, record := range records {
		subFirst := record.First.CreatedAt.Minute() * 60 + record.First.CreatedAt.Second() - secs
		subSecond := record.Second.CreatedAt.Minute() * 60 + record.Second.CreatedAt.Second() - secs
		if subFirst < 0 { // ABS
			subFirst = -subFirst
		}
		if subSecond < 0 { //ABS
			subSecond = -subSecond
		}

		if subFirst < minFirst {
			minFirst = subFirst
		}

		if subSecond < minSecond {
			minSecond = subSecond
		}

		if minPair >  minFirst + minSecond {
			minPair = minFirst + minSecond
			iSave = i
		}
	}

	if iSave == -1 {
		return iSave, nil
	}

	var results []entity.ResponsePriceRecord
	
	for _, record := range records {
		sub := record.CreatedAt.Unix() - records[iSave].CreatedAt.Unix()
		if sub < 0 {
			sub = -sub
		}
		if sub <= data.LongTime {
			results = append(results, record)
		}
	}

	var newISave int = -1
	for i, res := range results {
		if res.CreatedAt == records[iSave].CreatedAt {
			newISave = i
		}
	}

	return newISave, results
}