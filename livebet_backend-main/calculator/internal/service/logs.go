package service

import (
	"context"
	"fmt"
	"livebets/calculator/cmd/config"
	"livebets/calculator/internal/api"
	"livebets/calculator/internal/entity"
	"livebets/calculator/internal/repository"
	"livebets/calculator/pkg/cache"
	"livebets/calculator/pkg/rdbms"
	"livebets/calculator/pkg/utils"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const (
	edge = 5.0
	risk = 12.0
	bank = 10000.0
)

type UserCache struct {
	sync.RWMutex
	data map[string]map[string]entity.UserIDCache
}

type LogsService struct {
	txStorage           rdbms.TxStorage[repository.LogsStorage]
	analyzerAPI         *api.AnalizerAPI
	analyzerPrematchAPI *api.AnalizerPrematchAPI
	percentCache        cache.MemoryCacheInterface[string, entity.TotalPercent]
	usersCache          *UserCache
	logger              *zerolog.Logger
}

func NewLogsService(
	txStorage rdbms.TxStorage[repository.LogsStorage],
	analyzerAPI *api.AnalizerAPI,
	analyzerPrematchAPI *api.AnalizerPrematchAPI,
	logger *zerolog.Logger,
) *LogsService {
	percentCache := cache.NewMemoryCache[string, entity.TotalPercent]()
	usersCache := &UserCache{
		data:    make(map[string]map[string]entity.UserIDCache),
		RWMutex: sync.RWMutex{},
	}
	return &LogsService{
		txStorage:           txStorage,
		analyzerAPI:         analyzerAPI,
		analyzerPrematchAPI: analyzerPrematchAPI,
		percentCache:        percentCache,
		usersCache:          usersCache,
		logger:              logger,
	}
}

func (l *LogsService) InitializeTotalBetPercents(ctx context.Context) error {
	percents, err := l.txStorage.Storage().GetInitializeCalcBet(ctx)
	if err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.InitializeTotalBetPercents] get saved percent error")
		return err
	}

	for _, percent := range percents {
		l.percentCache.Write(percent.KeyMatch, entity.TotalPercent{TotalPercent: percent.TotalPercent, CreatedAt: time.Now()})
	}

	return nil
}

func (l *LogsService) CleanCaches(ctx context.Context, cfg config.LogsService, wg *sync.WaitGroup) {
	defer wg.Done()

	percentCacheInterval := time.Duration(time.Duration(cfg.PercentCacheInterval) * time.Second)
	percentCacheTicker := time.NewTicker(percentCacheInterval)

	usersCacheInterval := time.Duration(time.Duration(cfg.UsersCacheInterval) * time.Second)
	usersCacheTicker := time.NewTicker(usersCacheInterval)

	for {
		select {
		case <-percentCacheTicker.C:
			percentCache := l.percentCache.ReadAll()

			for key, value := range percentCache {
				if time.Since(value.CreatedAt) > (time.Duration(cfg.PercentCacheTimeout) * time.Second) {
					l.percentCache.Delete(key)
				}
			}

		case <-usersCacheTicker.C:
			l.usersCache.Lock()

			for key, users := range l.usersCache.data {
				for userK, user := range users {
					if time.Since(user.CreatedAt) > (time.Duration(cfg.UsersCacheTimeout) * time.Second) {
						delete(l.usersCache.data[key], userK)
					}
				}

				if len(l.usersCache.data[key]) == 0 {
					delete(l.usersCache.data, key)
				}
			}

			l.usersCache.Unlock()

		case <-ctx.Done():
			percentCacheTicker.Stop()
			usersCacheTicker.Stop()
			return
		}
	}
}

func (l *LogsService) LogBetAccept(ctx context.Context, pairAccept entity.AcceptBet) error {
	keyMatch := utils.GenerateFullMatchKey(pairAccept.Pair.First.Bookmaker, pairAccept.Pair.First.LeagueName, pairAccept.Pair.First.HomeName, pairAccept.Pair.First.AwayName, pairAccept.Pair.SportName, "")
	keyOutcome := utils.GenerateFullMatchKey(pairAccept.Pair.First.Bookmaker, pairAccept.Pair.Second.Bookmaker, pairAccept.Pair.First.MatchID, pairAccept.Pair.Second.MatchID, pairAccept.Pair.SportName, pairAccept.Pair.Outcome.Outcome)

	// Set percent
	percent := pairAccept.Sum / pairAccept.Bet.CalcBet.OriginalAmount * 100
	per, ok := l.percentCache.Read(keyMatch)
	if !ok {
		l.percentCache.Write(keyMatch, entity.TotalPercent{TotalPercent: percent, CreatedAt: time.Now()})
	} else {
		per.TotalPercent += percent
		per.CreatedAt = time.Now()
		l.percentCache.Write(keyMatch, per)
	}

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

	var priceRecods *entity.ResponsePriceRecords
	// Go to analyzer correct
	if pairAccept.Pair.IsLive {
		priceRecods, err = l.analyzerAPI.GeTPricesByTimeout(entity.RequestPriceRecordsByTime{
			Bookmaker1: pairAccept.Pair.First.Bookmaker,
			Bookmaker2: pairAccept.Pair.Second.Bookmaker,
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
			Bookmaker2: pairAccept.Pair.Second.Bookmaker,
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
		if err = l.txStorage.Storage().InsertLogBetAccept(ctx, keyMatch, keyOutcome, pairAccept, nil, percent, pairAccept.UserId, pairAccept.Pair.IsLive, pairAccept.Pair.SportName, pairAccept.Pair.Second.Bookmaker); err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
			return err
		}
		return nil
	}
	if len(priceRecods.Records) <= priceRecods.ISave {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices length records error")
		if err = l.txStorage.Storage().InsertLogBetAccept(ctx, keyMatch, keyOutcome, pairAccept, nil, percent, pairAccept.UserId, pairAccept.Pair.IsLive, pairAccept.Pair.SportName, pairAccept.Pair.Second.Bookmaker); err != nil {
			l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
			return err
		}
		return nil
	}

	if err = l.txStorage.Storage().InsertLogBetAccept(ctx, keyMatch, keyOutcome, pairAccept, &priceRecods.Records[priceRecods.ISave], percent, pairAccept.UserId, pairAccept.Pair.IsLive, pairAccept.Pair.SportName, pairAccept.Pair.Second.Bookmaker); err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] insert log bet accept error")
		return err
	}

	correctROI := priceRecods.Records[priceRecods.ISave].ROI
	go l.GetPricesForFlie(ctx, pairAccept, minutes, seconds, correctROI, false)

	return nil
}

func (l *LogsService) GetPricesForFlie(ctx context.Context, pairAccept entity.AcceptBet, minutes, seconds int, correctROI float64, isTest bool) error {
	if pairAccept.Pair.IsLive {
		time.Sleep(120 * time.Second)
	} else {
		time.Sleep(1200 * time.Second)
	}

	bookmakerForPrices := pairAccept.Pair.Second.Bookmaker
	if bookmakerForPrices == "Ladbrokes2" {
		bookmakerForPrices = "Ladbrokes"
	}

	var priceRecods *entity.ResponsePriceRecords
	var err error
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
		return nil
	}
	if len(priceRecods.Records) <= priceRecods.ISave {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices length records error")
		return nil
	}

	record := priceRecods.Records[priceRecods.ISave]
	//otherCoef := getPriceForSecond(priceRecods, record.CreatedAt, minutes, seconds, pairAccept.Pair.IsLive, pairAccept.Coef)
	//roi := shared.CalculateROI(otherCoef, record.First.Score, record.Margin, pairAccept.Pair.Outcome.MarketType, shared.Parser(pairAccept.Pair.Second.Bookmaker), shared.SportName(pairAccept.Pair.SportName))
	salary := pairAccept.Sum * (correctROI / 100) * 0.5 / 1.5

	var matchNameRaw string
	if isTest {
		matchNameRaw = fmt.Sprintf("%s_vs_%s_%s_%.2f_%s_%d", pairAccept.Pair.Second.HomeName, pairAccept.Pair.Second.AwayName, pairAccept.Pair.Second.Bookmaker, salary, record.Outcome, int64(pairAccept.Pair.Outcome.ROI))
	} else {
		matchNameRaw = fmt.Sprintf("%s_vs_%s_%s_%.2f_%s", pairAccept.Pair.Second.HomeName, pairAccept.Pair.Second.AwayName, pairAccept.Pair.Second.Bookmaker, salary, record.Outcome)
	}
	matchName := removeSpecialChars(replaceDotsInFileName(matchNameRaw))
	if err := l.createBetFile(matchName, *priceRecods, isTest); err != nil {
		l.logger.Error().Err(err).Msgf("[LogsService.LogBetAccept] get prices length records error")
		return nil
	}

	return err
}

func (l *LogsService) CalcSumBet(ctx context.Context, userID string, pair entity.PairOneOutcome) (entity.CalculatedBet, int) {
	keyMatch := utils.GenerateFullMatchKey(pair.First.Bookmaker, pair.First.LeagueName, pair.First.HomeName, pair.First.AwayName, pair.SportName, "")

	// IMPORTANT !!! ANALYZER SAVE PINNACLE TO FIRST
	// Calculation bet
	calcBet := l.calculateBet(keyMatch, pair.Outcome.Score1.Value)

	// Add new user to cache and get count
	l.usersCache.Lock()
	_, ok := l.usersCache.data[keyMatch]
	if !ok {
		l.usersCache.data[keyMatch] = make(map[string]entity.UserIDCache)
	}
	l.usersCache.data[keyMatch][userID] = entity.UserIDCache{UserID: userID, CreatedAt: time.Now()}

	usersCount := len(l.usersCache.data[keyMatch])

	l.usersCache.Unlock()

	calcBet.AdjustedAmount = calcBet.AdjustedAmount / float64(usersCount)

	return calcBet, usersCount
}

func (l *LogsService) calculateBet(keyMatch string, odds float64) entity.CalculatedBet {
	// Рассчитываем размер ставки
	originalAmount := l.getBetSize(odds)
	adjustedAmount := l.calculateAdjustedBetSize(keyMatch, originalAmount)

	// Рассчитываем процент оставшейся суммы
	percentage := 100.0
	if originalAmount > 0 {
		percentage = (adjustedAmount / originalAmount) * 100
		if percentage < 0 {
			percentage = 0
		} else if percentage > 100 {
			percentage = 100
		}
	} else {
		percentage = 0
	}

	return entity.CalculatedBet{
		OriginalAmount: originalAmount,
		AdjustedAmount: adjustedAmount,
		Percentage:     percentage,
	}
}

// getBetSize рассчитывает оптимальный размер ставки на основе критерия Келли
func (l *LogsService) getBetSize(odds float64) float64 {
	if edge < 0 {
		return 0
	}

	// Преобразуем edge из процентов в десятичную дробь
	edgeDecimal := edge / 100

	// Рассчитываем фактор внутри логарифма
	logFactor := 1 - (1 / (odds / (1 + edgeDecimal)))

	// Рассчитываем процент от банкролла для ставки
	betSizePercent := math.Log10(logFactor) / math.Log10(math.Pow(10, -risk))

	// Проверяем, что результат имеет смысл
	if betSizePercent < 0 || betSizePercent > 1 {
		return 0
	}

	// Рассчитываем фактический размер ставки
	betSize := betSizePercent * bank

	// Округляем до ближайшего числа, кратного 5
	roundedBetSize := math.Round(betSize/5) * 5

	return roundedBetSize
}

func (l *LogsService) calculateAdjustedBetSize(keyMatch string, baseBetSize float64) float64 {
	// Получаем процент уже поставленных денег на матч
	totalBetPercent, ok := l.percentCache.Read(keyMatch)
	if !ok {
		totalBetPercent.TotalPercent = 0
	}

	// Вычисляем оставшийся процент от базовой суммы ставки
	remainingPercentage := 100.0
	if baseBetSize > 0 {
		remainingPercentage -= totalBetPercent.TotalPercent
		if remainingPercentage < 0 {
			remainingPercentage = 0
		} else if remainingPercentage > 100.0 {
			remainingPercentage = 100.0
		}
	}

	// Корректируем размер ставки
	adjustedBetSize := baseBetSize * remainingPercentage / 100.0

	// Округляем до ближайшего числа, кратного 5
	adjustedBetSize = math.Round(adjustedBetSize/5) * 5

	return adjustedBetSize
}

func (l *LogsService) createBetFile(name string, records entity.ResponsePriceRecords, isTest bool) error {
	// Sort
	var sortRecords entity.ResponsePriceRecords
	for _, record := range records.Records {
		sortRecords.Records = append(sortRecords.Records, record)
	}

	// Sorting by time
	sort.Slice(sortRecords.Records, func(i, j int) bool {
		return sortRecords.Records[i].CreatedAt.Before(sortRecords.Records[j].CreatedAt)
	})

	for i, val := range sortRecords.Records {
		if val == records.Records[records.ISave] {
			sortRecords.ISave = i
		}
	}

	beforePrices := sortRecords.Records[:sortRecords.ISave]
	afterPrices := sortRecords.Records[sortRecords.ISave:]

	// Формируем строки CSV
	var csvBuilder strings.Builder
	csvBuilder.WriteString("Section,Time,Price\n") // Заголовок с новым столбцом

	// Добавляем "До ставки" цены
	if len(beforePrices) > 0 {
		csvBuilder.WriteString("Before Bet\n")
		csvBuilder.WriteString("Time,Price,Section\n") // Повторный заголовок для секции
		for _, price := range beforePrices {
			csvBuilder.WriteString(fmt.Sprintf("%s,%.2f,%s,%.2f,Before Bet\n", price.First.CreatedAt.Format(time.RFC3339), price.First.Score, price.Second.CreatedAt.Format(time.RFC3339), price.Second.Score))
		}
	} else {
		csvBuilder.WriteString("Before Bet\n")
		csvBuilder.WriteString("Time,Price,Section\n")
		csvBuilder.WriteString("No prices found\n")
	}

	// Добавляем "После ставки" цены
	if len(afterPrices) > 0 {
		csvBuilder.WriteString("After Bet\n")
		csvBuilder.WriteString("Time,Price,Section\n") // Повторный заголовок для секции
		for _, price := range afterPrices {
			csvBuilder.WriteString(fmt.Sprintf("%s,%.2f,%s,%.2f,After Bet\n", price.First.CreatedAt.Format(time.RFC3339), price.First.Score, price.Second.CreatedAt.Format(time.RFC3339), price.Second.Score))
		}
	} else {
		csvBuilder.WriteString("After Bet\n")
		csvBuilder.WriteString("Time,Price,Section\n")
		csvBuilder.WriteString("No prices found\n")
	}

	// Определяем имя файла
	var fileName string
	if isTest {
		fileName = fmt.Sprintf("logs/testbets/%s.csv", name)
	} else {
		fileName = fmt.Sprintf("logs/bets/%s.csv", name)
	}

	// Создаем директорию, если она отсутствует
	if err := os.MkdirAll("logs/bets", 0755); err != nil {
		log.Printf("[HISTORY] Ошибка создания директории logs: %v", err)
		return err
	}

	if err := os.WriteFile(fileName, []byte(csvBuilder.String()), 0644); err != nil {
		log.Printf("[HISTORY] Ошибка записи файла %s: %v", fileName, err)
		return err
	}

	return nil
}
