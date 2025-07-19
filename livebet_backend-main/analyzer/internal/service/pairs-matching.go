// Файл: livebet_backend-main/analyzer/internal/service/pairs-matching.go
// ФИНАЛЬНАЯ ЧИСТАЯ ВЕРСИЯ

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"livebets/analazer/cmd/config"
	"livebets/analazer/internal/entity"
	priceStorage "livebets/analazer/internal/price-storage"
	"livebets/analazer/internal/repository"
	bikeymap "livebets/analazer/pkg/bikey-map"
	"livebets/analazer/pkg/cache"
	"livebets/analazer/pkg/rdbms"
	redis_client "livebets/analazer/pkg/redis"
	"livebets/analazer/pkg/utils"
	"livebets/shared"
	"strconv"
	"strings"
	"sync"
	"time"

	fuzz "github.com/paul-mannino/go-fuzzywuzzy"
	"github.com/rs/zerolog"
)

type MatchGroup struct {
	ID      string
	Matches []entity.GameData
}

type BookmakerOdds struct {
	Bookmaker string
	Odds      float64
}

type PairsMatchingService struct {
	txStorage       rdbms.TxStorage[repository.PairsMatchingStorage]
	redisClient     *redis_client.Redis
	matchDataCache  cache.MemoryCacheInterface[string, entity.GameData]
	matchKeysCache  cache.MemoryCacheInterface[string, cache.MemoryCacheInterface[string, bool]]
	matchPairsCache bikeymap.BiKeyMapInterface[string, bool]
	pairs           cache.MemoryCacheInterface[string, entity.ResponsePair]
	receiveChan     <-chan entity.ReceivedMsg
	sendChan        chan<- []entity.ResponsePair
	priceStorage    *priceStorage.PriceStorage
	logger          *zerolog.Logger
	groupsCache     cache.MemoryCacheInterface[string, MatchGroup]
}

func NewPairsMatchingService(
	txStorage rdbms.TxStorage[repository.PairsMatchingStorage],
	redisClient *redis_client.Redis,
	receiveChan <-chan entity.ReceivedMsg,
	sendChan chan<- []entity.ResponsePair,
	priceStorage *priceStorage.PriceStorage,
	logger *zerolog.Logger,
) *PairsMatchingService {
	return &PairsMatchingService{
		txStorage:       txStorage,
		redisClient:     redisClient,
		receiveChan:     receiveChan,
		matchDataCache:  cache.NewMemoryCache[string, entity.GameData](),
		matchKeysCache:  cache.NewMemoryCache[string, cache.MemoryCacheInterface[string, bool]](),
		matchPairsCache: bikeymap.NewBiKeyMap[string, bool](),
		pairs:           cache.NewMemoryCache[string, entity.ResponsePair](),
		sendChan:        sendChan,
		priceStorage:    priceStorage,
		logger:          logger,
		groupsCache:     cache.NewMemoryCache[string, MatchGroup](),
	}
}

func (p *PairsMatchingService) Run(ctx context.Context, cfg config.PairsMatching, wg *sync.WaitGroup) {
	defer wg.Done()
	wgMatchWork := &sync.WaitGroup{}

	for i := 0; i < cfg.ReceiveWorkersCount; i++ {
		wgMatchWork.Add(1)
		go p.workerMatchData(ctx, wgMatchWork)
	}
	wgMatchWork.Add(1)
	go p.buildGroupsPeriodically(ctx, wgMatchWork)
	wgMatchWork.Add(1)
	go p.cleanCaches(ctx, cfg, wgMatchWork)
	wgMatchWork.Add(1)
	go p.updateKeysCache(ctx, cfg, wgMatchWork)
	wgMatchWork.Add(1)
	go p.updatePairsCache(ctx, cfg, wgMatchWork)
	wgMatchWork.Add(1)
	go p.send(ctx, cfg, wgMatchWork)

	wgMatchWork.Wait()
	p.logger.Info().Msg("[PairsMatchingService.Run] workers stopped")
}

func normalizeOutcomeKey(key string) string {
	value, err := strconv.ParseFloat(key, 64)
	if err != nil {
		return key
	}
	return fmt.Sprintf("%.1f", value)
}

func extractAllOutcomes(match entity.GameData) map[string]float64 {
	outcomes := make(map[string]float64)
	for i, period := range match.Periods {
		prefix := ""
		if i > 0 {
			prefix = fmt.Sprintf("P%d ", i)
		}
		if period.Win1x2.Win1.Value > 0 {
			outcomes[prefix+"1"] = period.Win1x2.Win1.Value
		}
		if period.Win1x2.WinNone.Value > 0 {
			outcomes[prefix+"X"] = period.Win1x2.WinNone.Value
		}
		if period.Win1x2.Win2.Value > 0 {
			outcomes[prefix+"2"] = period.Win1x2.Win2.Value
		}
		for key, total := range period.Totals {
			normalizedKey := normalizeOutcomeKey(key)
			if total.WinMore.Value > 0 {
				outcomes[fmt.Sprintf("%sT> %s", prefix, normalizedKey)] = total.WinMore.Value
			}
			if total.WinLess.Value > 0 {
				outcomes[fmt.Sprintf("%sT< %s", prefix, normalizedKey)] = total.WinLess.Value
			}
		}
		for key, handicap := range period.Handicap {
			normalizedKey := normalizeOutcomeKey(key)
			if handicap.Win1.Value > 0 {
				outcomes[fmt.Sprintf("%sH1 %s", prefix, normalizedKey)] = handicap.Win1.Value
			}
			if handicap.Win2.Value > 0 {
				outcomes[fmt.Sprintf("%sH2 %s", prefix, normalizedKey)] = handicap.Win2.Value
			}
		}
		for key, total := range period.FirstTeamTotals {
			normalizedKey := normalizeOutcomeKey(key)
			if total.WinMore.Value > 0 {
				outcomes[fmt.Sprintf("%sIT1> %s", prefix, normalizedKey)] = total.WinMore.Value
			}
			if total.WinLess.Value > 0 {
				outcomes[fmt.Sprintf("%sIT1< %s", prefix, normalizedKey)] = total.WinLess.Value
			}
		}
		for key, total := range period.SecondTeamTotals {
			normalizedKey := normalizeOutcomeKey(key)
			if total.WinMore.Value > 0 {
				outcomes[fmt.Sprintf("%sIT2> %s", prefix, normalizedKey)] = total.WinMore.Value
			}
			if total.WinLess.Value > 0 {
				outcomes[fmt.Sprintf("%sIT2< %s", prefix, normalizedKey)] = total.WinLess.Value
			}
		}
		for key, game := range period.Games {
			if game.Win1.Value > 0 {
				outcomes[fmt.Sprintf("%s1G %s", prefix, key)] = game.Win1.Value
			}
			if game.Win2.Value > 0 {
				outcomes[fmt.Sprintf("%s2G %s", prefix, key)] = game.Win2.Value
			}
		}
	}
	return outcomes
}

func (p *PairsMatchingService) buildGroupsPeriodically(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			groupsByPinnacleId := make(map[string]*MatchGroup)
			allMatches := p.matchDataCache.ReadAll()
			allPairedKeys, _ := p.matchPairsCache.ReadAll()
			for key1, key2 := range allPairedKeys {
				match1, ok1 := allMatches[key1]
				match2, ok2 := allMatches[key2]
				if !ok1 || !ok2 {
					continue
				}
				var pinnacleMatch, partnerMatch entity.GameData
				if match1.Source == string(shared.PINNACLE) {
					pinnacleMatch = match1
					partnerMatch = match2
				} else if match2.Source == string(shared.PINNACLE) {
					pinnacleMatch = match2
					partnerMatch = match1
				} else {
					continue
				}
				pinnacleId := pinnacleMatch.MatchId
				if _, exists := groupsByPinnacleId[pinnacleId]; !exists {
					groupsByPinnacleId[pinnacleId] = &MatchGroup{
						ID:      pinnacleId,
						Matches: []entity.GameData{pinnacleMatch},
					}
				}
				alreadyInGroup := false
				for _, m := range groupsByPinnacleId[pinnacleId].Matches {
					if m.Source == partnerMatch.Source {
						alreadyInGroup = true
						break
					}
				}
				if !alreadyInGroup {
					groupsByPinnacleId[pinnacleId].Matches = append(groupsByPinnacleId[pinnacleId].Matches, partnerMatch)
				}
			}
			p.groupsCache.Lock()
			p.groupsCache.CleanUnsafe()
			for pinnacleId, group := range groupsByPinnacleId {
				p.groupsCache.WriteUnsafe(pinnacleId, *group)
				if len(group.Matches) >= 3 {
					outcomesAcrossGroup := make(map[string][]BookmakerOdds)
					for _, match := range group.Matches {
						if time.Since(match.CreatedAt) > 10*time.Second {
							continue
						}
						allOutcomes := extractAllOutcomes(match)
						for outcomeStr, odds := range allOutcomes {
							outcomesAcrossGroup[outcomeStr] = append(outcomesAcrossGroup[outcomeStr], BookmakerOdds{Bookmaker: match.Source, Odds: odds})
						}
					}
					for outcomeStr, prices := range outcomesAcrossGroup {
						if len(prices) >= 3 {
							var pricesLog strings.Builder
							for _, pr := range prices {
								pricesLog.WriteString(fmt.Sprintf("%s: %.2f, ", pr.Bookmaker, pr.Odds))
							}
							finalPricesStr := strings.TrimSuffix(pricesLog.String(), ", ")
							p.logger.Info().
								Str("pinnacleMatchId", group.ID).
								Str("outcome", outcomeStr).
								Msgf("[GROUP PRICE] Prices: [%s]", finalPricesStr)
						}
					}
				}
			}
			p.groupsCache.Unlock()
		case <-ctx.Done():
			p.logger.Info().Msg("[PairsMatchingService.buildGroupsPeriodically] context cancelled")
			return
		}
	}
}

func (p *PairsMatchingService) cleanCaches(ctx context.Context, cfg config.PairsMatching, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()
	cleanCacheInterval := time.Duration(time.Duration(cfg.ClearCacheInterval) * time.Second)
	cleanCacheTicker := time.NewTicker(cleanCacheInterval)
	for {
		select {
		case <-cleanCacheTicker.C:
			data := p.matchDataCache.ReadAll()
			for matchKey, matchValue := range data {
				if time.Since(matchValue.CreatedAt) > (time.Duration(cfg.MatchDataTimeout) * time.Second) {
					p.matchDataCache.Delete(matchKey)
					p.matchKeysCache.Delete(matchKey)
					p.matchPairsCache.Delete(matchKey)
					if matchValue.Source == string(shared.PINNACLE) {
						p.groupsCache.Delete(matchValue.MatchId)
					}
				}
			}
			keysMap := p.matchKeysCache.ReadAll()
			for key := range keysMap {
				_, ok := p.matchDataCache.Read(key)
				if !ok {
					p.matchKeysCache.Delete(key)
					p.matchPairsCache.Delete(key)
				}
			}
			keysCachePair, _ := p.matchPairsCache.ReadAll()
			for key := range keysCachePair {
				_, ok := p.matchDataCache.Read(key)
				if !ok {
					p.matchPairsCache.Delete(key)
				}
			}
		case <-ctx.Done():
			cleanCacheTicker.Stop()
			return
		}
	}
}

func (p *PairsMatchingService) processAndCachePair(pinnacleMatch, partnerMatch entity.GameData) {
	partnerMatch, err := reverseCoefs(fmt.Sprintf("%s %s", pinnacleMatch.HomeName, pinnacleMatch.AwayName), partnerMatch)
	if err != nil {
		return
	}
	commonOutcomes := findCommonOutcomes(partnerMatch.Periods, pinnacleMatch.Periods, int(pinnacleMatch.HomeScore), int(pinnacleMatch.AwayScore))
	if len(commonOutcomes) == 0 {
		return
	}
	filteredOutcomes := p.calculateAndFilterCommonOutcomes(commonOutcomes, partnerMatch.Source, pinnacleMatch.SportName)
	if len(filteredOutcomes) == 0 {
		return
	}
	result := entity.ResponsePair{
		First: entity.ResponseMatch{
			Bookmaker:  pinnacleMatch.Source,
			LeagueName: pinnacleMatch.LeagueName,
			HomeScore:  pinnacleMatch.HomeScore,
			AwayScore:  pinnacleMatch.AwayScore,
			HomeName:   pinnacleMatch.HomeName,
			AwayName:   pinnacleMatch.AwayName,
			MatchID:    pinnacleMatch.MatchId,
			CreatedAt:  pinnacleMatch.CreatedAt,
			Raw:        pinnacleMatch.Raw,
		},
		Second: entity.ResponseMatch{
			Bookmaker:  partnerMatch.Source,
			LeagueName: partnerMatch.LeagueName,
			HomeScore:  partnerMatch.HomeScore,
			AwayScore:  partnerMatch.AwayScore,
			HomeName:   partnerMatch.HomeName,
			AwayName:   partnerMatch.AwayName,
			MatchID:    partnerMatch.MatchId,
			CreatedAt:  partnerMatch.CreatedAt,
			Raw:        partnerMatch.Raw,
		},
		Outcome:   filteredOutcomes,
		IsLive:    pinnacleMatch.IsLive,
		SportName: string(pinnacleMatch.SportName),
		CreatedAt: time.Now(),
	}
	pairKey := pinnacleMatch.Source + string(pinnacleMatch.MatchId) + string(pinnacleMatch.SportName) + partnerMatch.Source + string(partnerMatch.MatchId)
	p.pairs.Write(pairKey, result)
	for _, out := range filteredOutcomes {
		fullKey := utils.GenerateFullMatchKey(pinnacleMatch.Source, partnerMatch.Source, pinnacleMatch.MatchId, partnerMatch.MatchId, string(pinnacleMatch.SportName), out.Outcome)
		p.priceStorage.Write(fullKey, result.CreatedAt, entity.FullPriceRecord{
			First: entity.PriceRecord{
				Bookmaker: pinnacleMatch.Source,
				Score:     out.Score1.Value,
				CreatedAt: pinnacleMatch.CreatedAt,
			},
			Second: entity.PriceRecord{
				Bookmaker: partnerMatch.Source,
				Score:     out.Score2.Value,
				CreatedAt: partnerMatch.CreatedAt,
			},
			Outcome: out.Outcome,
			ROI:     out.ROI,
			Margin:  out.Margin,
		})
	}
}

func getPinnaclePartner(match1, match2 entity.GameData) (pinnacle, partner entity.GameData, isPinnaclePair bool) {
	if match1.Source == string(shared.PINNACLE) {
		return match1, match2, true
	}
	if match2.Source == string(shared.PINNACLE) {
		return match2, match1, true
	}
	return entity.GameData{}, entity.GameData{}, false
}

func (p *PairsMatchingService) workerMatchData(ctx context.Context, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()
	for {
		select {
		case msg := <-p.receiveChan:
			var gameData entity.GameData
			err := json.Unmarshal(msg, &gameData)
			if err != nil {
				p.logger.Error().Err(err).Msg("[PairsMatchingService.worker] game data unmarhall error")
				continue
			}
			if gameData.Periods == nil {
				continue
			}
			key := createKeyMatchData(gameData.Source, string(gameData.SportName), gameData.Pid)
			p.matchDataCache.Write(key, gameData)

			keyPair, ok := p.matchPairsCache.ReadKey(key)
			if ok {
				match1, ok1 := p.matchDataCache.Read(key)
				match2, ok2 := p.matchDataCache.Read(keyPair)
				if ok1 && ok2 {
					pinnacleMatch, partnerMatch, isPinnaclePair := getPinnaclePartner(match1, match2)
					if isPinnaclePair {
						p.processAndCachePair(pinnacleMatch, partnerMatch)
					}
				}
				continue
			}

			if gameData.Source != string(shared.PINNACLE) {
				allMatches := p.matchDataCache.ReadAll()
				for _, pinnacleCandidate := range allMatches {
					if pinnacleCandidate.Source == string(shared.PINNACLE) &&
						pinnacleCandidate.SportName == gameData.SportName &&
						fuzz.Ratio(
							fmt.Sprintf("%s %s", gameData.HomeName, gameData.AwayName),
							fmt.Sprintf("%s %s", pinnacleCandidate.HomeName, pinnacleCandidate.AwayName),
						) > 85 {

						p.logger.Info().Msgf("[TEMP MATCH] Found temporary pair for %s with Pinnacle by name similarity.", gameData.Source)
						p.processAndCachePair(pinnacleCandidate, gameData)
						break
					}
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (p *PairsMatchingService) send(ctx context.Context, cfg config.PairsMatching, wgWork *sync.WaitGroup) {
	defer wgWork.Done()
	interval := time.Duration(time.Duration(cfg.SendInterval) * time.Millisecond)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			pairs := p.pairs.ReadAll()
			var results []entity.ResponsePair
			for key, val := range pairs {
				if time.Since(val.CreatedAt) > (time.Duration(cfg.PairTimeout) * time.Second) {
					p.pairs.Delete(key)
				} else {
					results = append(results, val)
					msg, err := json.Marshal(val)
					if err != nil {
						p.logger.Error().Err(err).Msg("[PairsMatchingService.send] value marshall error")
					}
					redisKey := shared.GetRKeyPairs(val.IsLive, val.First.Bookmaker, val.Second.Bookmaker)
					err = p.redisClient.Publish(ctx, redisKey, msg)
					if err != nil {
						p.logger.Error().Err(err).Msg("[PairsMatchingService.send] write msg to redis error")
					}
				}
			}
			p.sendChan <- results
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (p *PairsMatchingService) updateKeysCache(ctx context.Context, cfg config.PairsMatching, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()
	updateKeysCacheInterval := time.Duration(time.Duration(cfg.UpdateKeysCacheInterval) * time.Second)
	updateKeysCacheTicker := time.NewTicker(updateKeysCacheInterval)
	for {
		select {
		case <-updateKeysCacheTicker.C:
			data := p.matchDataCache.ReadAll()
			for keyMatch, valueMatch := range data {
				uuids, err := p.txStorage.Storage().GetUUIDKeys(ctx, valueMatch.Source, string(valueMatch.SportName),
					valueMatch.LeagueName, valueMatch.HomeName, valueMatch.AwayName)
				if err != nil {
					p.logger.Error().Err(err).Msg("[PairsMatchingService.updateKeysCache] get uuid keys error")
				}
				if len(uuids) < 2 {
					continue
				}
				newKeys := cache.NewMemoryCache[string, bool]()
				for _, uuid := range uuids {
					newKeys.Write(uuid, true)
				}
				p.matchKeysCache.Write(keyMatch, newKeys)
			}
		case <-ctx.Done():
			updateKeysCacheTicker.Stop()
			return
		}
	}
}

func (p *PairsMatchingService) updatePairsCache(ctx context.Context, cfg config.PairsMatching, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()
	updatePairsCacheInterval := time.Duration(time.Duration(cfg.UpdatePairsCacheInterval) * time.Second)
	updatePairsCacheTicker := time.NewTicker(updatePairsCacheInterval)
	for {
		select {
		case <-updatePairsCacheTicker.C:
			matchKeys := p.matchKeysCache.ReadAll()
			for key1 := range matchKeys {
				for key2 := range matchKeys {
					if key1 != key2 {
						uuidsKey1 := matchKeys[key1]
						uuidsKey2 := matchKeys[key2]
						if uuidsKey1 == nil || uuidsKey2 == nil {
							continue
						}
						uuids1 := uuidsKey1.ReadAll()
						uuids2 := uuidsKey2.ReadAll()
						counter := 0
						for uuid := range uuids1 {
							uuidEqual := uuids2[uuid]
							if !uuidEqual {
								continue
							}
							counter++
						}
						if counter == 2 {
							p.matchPairsCache.WriteBothKeys(key1, key2, true)
						}
					}
				}
			}
			fmt.Printf("Match Data - %d\n", p.matchDataCache.Len())
			fmt.Printf("Match Keys - %d\n", p.matchKeysCache.Len())
			keyC, valueC := p.matchPairsCache.Len()
			fmt.Printf("Match Pairs - %d : %d\n", keyC, valueC)
		case <-ctx.Done():
			updatePairsCacheTicker.Stop()
			return
		}
	}
}

func (p *PairsMatchingService) GetMatchData(ctx context.Context) map[string]entity.GameData {
	return p.matchDataCache.ReadAll()
}

func (p *PairsMatchingService) GetCacheKeys(ctx context.Context) map[string]map[string]bool {
	cacheKeys := p.matchKeysCache.ReadAll()
	newMap := make(map[string]map[string]bool)
	for key, value := range cacheKeys {
		newMap[key] = value.ReadAll()
	}
	return newMap
}

func (p *PairsMatchingService) GetCachePairs(ctx context.Context) (map[string]string, map[string]bool) {
	return p.matchPairsCache.ReadAll()
}

func (p *PairsMatchingService) GetPairs(ctx context.Context) map[string]entity.ResponsePair {
	return p.pairs.ReadAll()
}

func createKeyMatchData(bookmaker, sportName string, pid int64) string {
	return fmt.Sprintf("%s%s%d", bookmaker, sportName, pid)
}

func (p *PairsMatchingService) GetOnlineMatchData(ctx context.Context) []entity.MatchData {
	matchDataMap := p.matchDataCache.ReadAll()
	var matchData []entity.MatchData
	for _, match := range matchDataMap {
		matchData = append(matchData, entity.MatchData{
			LeagueName: match.LeagueName,
			HomeName:   match.HomeName,
			AwayName:   match.AwayName,
			MatchID:    match.MatchId,
			Bookmaker:  match.Source,
			SportName:  string(match.SportName),
			CreatedAt:  match.CreatedAt,
		})
	}
	return matchData
}
