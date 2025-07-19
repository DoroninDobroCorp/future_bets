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
	"sync"
	"time"

	"github.com/rs/zerolog"
)

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
}

func NewPairsMatchingService(
	txStorage rdbms.TxStorage[repository.PairsMatchingStorage],
	redisClient *redis_client.Redis,
	receiveChan <-chan entity.ReceivedMsg,
	sendChan chan<- []entity.ResponsePair,
	priceStorage *priceStorage.PriceStorage,
	logger *zerolog.Logger,
) *PairsMatchingService {
	matchDataCache := cache.NewMemoryCache[string, entity.GameData]()
	matchKeysCache := cache.NewMemoryCache[string, cache.MemoryCacheInterface[string, bool]]()
	matchPairsCache := bikeymap.NewBiKeyMap[string, bool]()
	pairs := cache.NewMemoryCache[string, entity.ResponsePair]()
	return &PairsMatchingService{
		txStorage:       txStorage,
		redisClient:     redisClient,
		receiveChan:     receiveChan,
		matchDataCache:  matchDataCache,
		matchKeysCache:  matchKeysCache,
		matchPairsCache: matchPairsCache,
		pairs:           pairs,
		sendChan:        sendChan,
		priceStorage:    priceStorage,
		logger:          logger,
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

func (p *PairsMatchingService) cleanCaches(ctx context.Context, cfg config.PairsMatching, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()

	cleanCacheInterval := time.Duration(time.Duration(cfg.ClearCacheInterval) * time.Second)
	cleanCacheTicker := time.NewTicker(cleanCacheInterval)

	for {
		select {
		case <-cleanCacheTicker.C:

			// Delete element by time
			data := p.matchDataCache.ReadAll()
			for matchKey, matchValue := range data {
				if time.Since(matchValue.CreatedAt) > (time.Duration(cfg.MatchDataTimeout) * time.Second) {
					// Clear match data cache
					p.matchDataCache.Delete(matchKey)

					// Clear key cache
					p.matchKeysCache.Delete(matchKey)

					//Clear pair cache
					p.matchPairsCache.Delete(matchKey)
				}
			}

			// Delete element by keys
			keysMap := p.matchKeysCache.ReadAll()
			for key := range keysMap {
				_, ok := p.matchDataCache.Read(key)
				if !ok {
					p.matchKeysCache.Delete(key)

					//Clear pair cache
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

func (p *PairsMatchingService) workerMatchData(ctx context.Context, wgMatchWork *sync.WaitGroup) {
	defer wgMatchWork.Done()

	for {
		select {
		case msg := <-p.receiveChan:

			var gameData entity.GameData
			err := json.Unmarshal(msg, &gameData)
			if err != nil {
				p.logger.Error().Err(err).Msg("[PairsMatchingService.worker] game data unmarhall error")
				fmt.Println(string(msg))
			}

			if gameData.Periods == nil {
				continue
			}

			key := createKeyMatchData(gameData.Source, string(gameData.SportName), gameData.Pid)

			_, ok := p.matchDataCache.Read(key)
			if !ok {
				if gameData.Source == "" || gameData.SportName == "" || gameData.LeagueName == "" {
					continue
				}

				leagueID, err := p.txStorage.Storage().InsertLeague(ctx, gameData.Source, string(gameData.SportName), gameData.LeagueName)
				if err != nil {
					p.logger.Error().Err(err).Msg("[PairsMatchingService.worker] insert league error")
				} else {
					if leagueID != nil {
						if err := p.txStorage.Storage().InsertTeam(ctx, *leagueID, gameData.HomeName); err != nil {
							p.logger.Error().Err(err).Msg("[PairsMatchingService.worker] insert home team error")
						}
						if err := p.txStorage.Storage().InsertTeam(ctx, *leagueID, gameData.AwayName); err != nil {
							p.logger.Error().Err(err).Msg("[PairsMatchingService.worker] insert away team error")
						}
					}
				}
			}

			p.matchDataCache.Write(key, gameData)

			// Process match
			keyPair, ok := p.matchPairsCache.ReadKey(key)
			if ok {
				match1, ok1 := p.matchDataCache.Read(key)
				match2, ok2 := p.matchDataCache.Read(keyPair)

				if ok1 && ok2 {

					// Only for PINNACLE pairs
					var value1 entity.GameData = match1 // value1 always IS PINNACLE
					var value2 entity.GameData = match2
					if match2.Source == string(shared.PINNACLE) {
						value1 = match2
						value2 = match1
					}

					value2, err = reverseCoefs(fmt.Sprintf("%s %s", value1.HomeName, value1.AwayName), value2)
					if err != nil {
						continue
					}

					commonOutcomes := findCommonOutcomes(value2.Periods, value1.Periods, int(value1.HomeScore), int(value1.AwayScore)) // 2 arg PINNACLE important!!!
					if commonOutcomes == nil || len(commonOutcomes) == 0 {
						continue
					}

					filtered := p.calculateAndFilterCommonOutcomes(commonOutcomes, value2.Source, value1.SportName)
					if len(filtered) == 0 {
						continue
					}

					result := entity.ResponsePair{
						First: entity.ResponseMatch{
							Bookmaker:  value1.Source,
							LeagueName: value1.LeagueName,
							HomeScore:  value1.HomeScore,
							AwayScore:  value1.AwayScore,
							HomeName:   value1.HomeName,
							AwayName:   value1.AwayName,
							MatchID:    value1.MatchId,
							CreatedAt:  value1.CreatedAt,
							Raw:        value1.Raw,
						},
						Second: entity.ResponseMatch{
							Bookmaker:  value2.Source,
							LeagueName: value2.LeagueName,
							HomeScore:  value2.HomeScore,
							AwayScore:  value2.AwayScore,
							HomeName:   value2.HomeName,
							AwayName:   value2.AwayName,
							MatchID:    value2.MatchId,
							CreatedAt:  value2.CreatedAt,
							Raw:        value2.Raw,
						},
						Outcome:   filtered,
						IsLive:    value1.IsLive,
						SportName: string(value1.SportName),
						CreatedAt: time.Now(),
					}

					p.pairs.Write(value1.Source+string(value1.MatchId)+string(value1.SportName)+value2.Source+string(value2.MatchId), result)

					// Add data to price storage
					for _, out := range filtered {
						fullKey := utils.GenerateFullMatchKey(value1.Source, value2.Source, value1.MatchId, value2.MatchId, string(value1.SportName), out.Outcome)
						p.priceStorage.Write(fullKey, result.CreatedAt, entity.FullPriceRecord{
							First: entity.PriceRecord{
								Bookmaker: value1.Source,
								Score:     out.Score1.Value,
								CreatedAt: value1.CreatedAt,
							},
							Second: entity.PriceRecord{
								Bookmaker: value2.Source,
								Score:     out.Score2.Value,
								CreatedAt: value2.CreatedAt,
							},
							Outcome: out.Outcome,
							ROI:     out.ROI,
							Margin:  out.Margin,
						})
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
