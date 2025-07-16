package service

import (
	"context"
	"fmt"
	"livebets/calculator/internal/entity"
	"livebets/calculator/pkg/utils"
	"strconv"
	"strings"
)

func (l *LogsService) LogTestBetAccept(ctx context.Context, pairAccept entity.AcceptBet) error {
	keyMatch := utils.GenerateFullMatchKey(pairAccept.Pair.First.Bookmaker, pairAccept.Pair.First.LeagueName, pairAccept.Pair.First.HomeName, pairAccept.Pair.First.AwayName, pairAccept.Pair.SportName, "")
	keyOutcome := utils.GenerateFullMatchKey(pairAccept.Pair.First.Bookmaker, pairAccept.Pair.Second.Bookmaker, pairAccept.Pair.First.MatchID, pairAccept.Pair.Second.MatchID, pairAccept.Pair.SportName, pairAccept.Pair.Outcome.Outcome)

	// Set percent
	percent := pairAccept.Sum / pairAccept.Bet.CalcBet.OriginalAmount * 100
	// per, ok := l.percentCache.Read(keyMatch)
	// if !ok {
	// 	l.percentCache.Write(keyMatch, entity.TotalPercent{TotalPercent: percent, CreatedAt: time.Now()})
	// } else {
	// 	per.TotalPercent += percent
	// 	per.CreatedAt = time.Now()
	// 	l.percentCache.Write(keyMatch, per)
	// }

	// Parse time
	strs := strings.Split(pairAccept.Time, ":")
	if len(strs) != 2 {
		err := fmt.Errorf("split time correct error")
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] split time correct error")
		return err
	}

	minutes, err := strconv.Atoi(strs[0])
	if err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] parse string to int error")
		return err
	}

	seconds, err := strconv.Atoi(strs[1])
	if err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] parse string to int error")
		return err
	}

	bookmakerForPrices := pairAccept.Pair.Second.Bookmaker
	if bookmakerForPrices == "Ladbrokes2" {
		bookmakerForPrices = "Ladbrokes"
	}

	var priceRecods *entity.ResponsePriceRecords
	// Go to analyzer correct
	if pairAccept.Pair.IsLive {
		priceRecods, err = l.analyzerAPI.GeTPricesByTimeout(entity.RequestPriceRecordsByTime{
			Bookmaker1: pairAccept.Pair.First.Bookmaker,
			Bookmaker2: bookmakerForPrices,
			MatchID1:   pairAccept.Pair.First.MatchID,
			MatchID2:   pairAccept.Pair.Second.MatchID,
			SportName:  pairAccept.Pair.SportName,
			Outcome:    pairAccept.Pair.Outcome.Outcome,

			Minutes:  minutes,
			Seconds:  seconds,
			LongTime: 120,
		})
		if err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices live error")
		}
	} else {
		priceRecods, err = l.analyzerPrematchAPI.GeTPricesByTimeout(entity.RequestPriceRecordsByTime{
			Bookmaker1: pairAccept.Pair.First.Bookmaker,
			Bookmaker2: bookmakerForPrices,
			MatchID1:   pairAccept.Pair.First.MatchID,
			MatchID2:   pairAccept.Pair.Second.MatchID,
			SportName:  pairAccept.Pair.SportName,
			Outcome:    pairAccept.Pair.Outcome.Outcome,

			Minutes:  minutes,
			Seconds:  seconds,
			LongTime: 1200,
		})
		if err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices prematch error")
		}
	}

	if priceRecods == nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices nil error")
		if err = l.txStorage.Storage().InsertLogTestBetAccept(ctx, keyMatch, keyOutcome, pairAccept, nil, percent); err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
			return err
		}
		return nil
	}
	if len(priceRecods.Records) <= priceRecods.ISave {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices length records error")
		if err = l.txStorage.Storage().InsertLogTestBetAccept(ctx, keyMatch, keyOutcome, pairAccept, nil, percent); err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
			return err
		}
		return nil
	}

	if err = l.txStorage.Storage().InsertLogTestBetAccept(ctx, keyMatch, keyOutcome, pairAccept, &priceRecods.Records[priceRecods.ISave], percent); err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
		return err
	}

	if pairAccept.Pair.Outcome.ROI > 3 && pairAccept.Pair.Outcome.ROI < 15 {
		if err = sendMissedBet(pairAccept.Pair, keyMatch); err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] send missed bet error")
		}
	}

	correctROI := priceRecods.Records[priceRecods.ISave].ROI
	go l.GetPricesForFlie(ctx, pairAccept, minutes, seconds, correctROI, true)

	return nil
}
