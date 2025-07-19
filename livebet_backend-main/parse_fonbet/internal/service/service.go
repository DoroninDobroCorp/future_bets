package service

import (
	"context"
	"fmt"
	"livebets/parse_fonbet/cmd/config"
	"livebets/parse_fonbet/internal/api"
	"livebets/parse_fonbet/internal/entity"
	"livebets/shared"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type GeneralService struct {
	pAPI         *api.FonbetAPI
	sendChan     chan<- entity.ResponseGame
	footballData map[int64]*entity.ResponseGame
	tennisData   map[int64]*entity.ResponseGame
	logger       *zerolog.Logger
}

func NewGeneralService(
	pAPI *api.FonbetAPI,
	sendChan chan<- entity.ResponseGame,
	logger *zerolog.Logger,
) *GeneralService {
	return &GeneralService{
		pAPI:         pAPI,
		sendChan:     sendChan,
		footballData: make(map[int64]*entity.ResponseGame),
		tennisData:   make(map[int64]*entity.ResponseGame),
		logger:       logger,
	}
}

func (s *GeneralService) RunFootball(ctx context.Context, cfg config.FonbetConfig, wg *sync.WaitGroup, isLive bool) {
	defer wg.Done()
	sportID := entity.FootballID
	place := "line"
	if isLive {
		place = "live"
	}

	matchInterval := time.Duration(time.Duration(cfg.IntervalMatch) * time.Second)
	matchTicker := time.NewTicker(matchInterval)

	oddsInterval := time.Duration(time.Duration(cfg.IntervalODDS) * time.Second)
	oddsTicker := time.NewTicker(oddsInterval)

	for {
		select {
		case <-matchTicker.C:
			eventsData, err := s.pAPI.GetAllData()
			if err != nil {
				s.logger.Error().Err(err).Msg("[Service.RunFootball] error get all live events.")
				continue
			}

			for _, event := range eventsData.Events {
				if event.Place != place {
					continue
				}

				var eventLeague entity.League
				for _, league := range eventsData.Leagues {
					if league.Id == event.SportId {
						eventLeague = league
						break
					}
				}

				if eventLeague.SportId != int64(sportID) {
					continue
				}

				var homeScore, awayScore float64
				for _, eventMisc := range eventsData.EventMiscs {
					if eventMisc.Id == event.Id {
						homeScore = float64(eventMisc.Score1)
						awayScore = float64(eventMisc.Score2)
						break
					}
				}

				if event.Name == "" {
					s.footballData[event.Id] = &entity.ResponseGame{
						Pid:        event.Id,
						LeagueName: normalizeFootballLeague(eventLeague.Name),
						HomeName:   normalizeFootballTeam(event.HomeName),
						AwayName:   normalizeFootballTeam(event.AwayName),
						MatchId:    fmt.Sprintf("%d", event.Id),
						IsLive:     isLive,

						HomeScore: homeScore,
						AwayScore: awayScore,

						Source:    string(shared.FONBET),
						SportName: string(shared.SOCCER),
						CreatedAt: time.Now(),
					}
				}
			}

		case <-oddsTicker.C:
			matchIds := make([]int64, 0, len(s.footballData))
			for _, match := range s.footballData {
				matchIds = append(matchIds, match.Pid)
			}

			for _, matchId := range matchIds {
				go func(matchId int64) {
					eventData, err := s.pAPI.GetMatchODDS(matchId)
					if err != nil {
						s.logger.Error().Err(err).Msg("failed to get odds")
						return
					}

					eventsParse := make(map[int64]entity.Event)

					for _, eventData := range eventData.Events {
						eventsParse[eventData.Id] = eventData
					}

					periods := make([]entity.ResponsePeriod, 3)

					for i := range periods {
						periods[i].Win1x2 = entity.Win1x2Struct{}
						periods[i].Games = make(map[string]*entity.Win1x2Struct)
						periods[i].Totals = make(map[string]*entity.WinLessMore)
						periods[i].Handicap = make(map[string]*entity.WinHandicap)
						periods[i].FirstTeamTotals = make(map[string]*entity.WinLessMore)
						periods[i].SecondTeamTotals = make(map[string]*entity.WinLessMore)
					}

					for _, eventOdds := range eventData.CustomFactors {
						event, exist := eventsParse[eventOdds.EventId]
						if !exist {
							continue
						}

						periodIndex := 0
						if event.Name == "" {
							periodIndex = 0
						} else if event.Name == "1st half" {
							periodIndex = 1
						} else if event.Name == "2nd half" {
							periodIndex = 2
						} else {
							continue
						}

						for _, outcome := range eventOdds.Factors {
							// Process 1x2
							if mapping, ok := win1x2Mappings[outcome.FactorId]; ok {
								setWin1x2Value(&periods[periodIndex].Win1x2, mapping.oddType, outcome)
								continue
							}

							// Process totals
							if mapping, ok := totalsMappings[outcome.FactorId]; ok {
								ensureMapEntry(periods[periodIndex].Totals, normalizeLine(outcome.Line))
								setTotalValue(periods[periodIndex].Totals[normalizeLine(outcome.Line)], mapping.oddType, outcome)
								continue
							}

							// Process team totals
							if mapping, ok := teamTotalsMappings[outcome.FactorId]; ok {
								var totalsMap map[string]*entity.WinLessMore
								if mapping.team == "first" {
									totalsMap = periods[periodIndex].FirstTeamTotals
								} else {
									totalsMap = periods[periodIndex].SecondTeamTotals
								}

								ensureMapEntry(totalsMap, normalizeLine(outcome.Line))
								setTotalValue(totalsMap[normalizeLine(outcome.Line)], mapping.oddType, outcome)
								continue
							}

							// Process handicaps
							processHandicap(&periods, periodIndex, outcome)
						}
					}

					s.footballData[matchId].Periods = periods

					s.sendChan <- *s.footballData[matchId]
				}(matchId)
			}

		case <-ctx.Done():
			matchTicker.Stop()
			oddsTicker.Stop()
			return
		}
	}
}

func (s *GeneralService) RunTennis(ctx context.Context, cfg config.FonbetConfig, wg *sync.WaitGroup, isLive bool) {
	defer wg.Done()
	sportID := entity.TennisID
	place := "line"
	if isLive {
		place = "live"
	}

	matchInterval := time.Duration(time.Duration(cfg.IntervalMatch) * time.Second)
	matchTicker := time.NewTicker(matchInterval)

	oddsInterval := time.Duration(time.Duration(cfg.IntervalODDS) * time.Second)
	oddsTicker := time.NewTicker(oddsInterval)

	for {
		select {
		case <-matchTicker.C:
			eventsData, err := s.pAPI.GetAllData()
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Service.Run] error get all matches. sportID - %d", sportID)
				continue
			}

			for _, event := range eventsData.Events {
				if event.Place != place {
					continue
				}

				var eventLeague entity.League
				for _, league := range eventsData.Leagues {
					if league.Id == event.SportId {
						eventLeague = league
						break
					}
				}

				if eventLeague.SportId != int64(sportID) {
					continue
				}

				var homeScore, awayScore float64
				for _, eventMisc := range eventsData.EventMiscs {
					if eventMisc.Id == event.Id {
						homeScore = float64(eventMisc.Score1)
						awayScore = float64(eventMisc.Score2)
						break
					}
				}

				if event.Name == "" {
					s.tennisData[event.Id] = &entity.ResponseGame{
						Pid:        event.Id,
						LeagueName: normalizeTennisLeague(eventLeague.Name),
						HomeName:   normalizeTennisTeam(event.HomeName),
						AwayName:   normalizeTennisTeam(event.AwayName),
						MatchId:    fmt.Sprintf("%d", event.Id),
						IsLive:     isLive,

						HomeScore: homeScore,
						AwayScore: awayScore,

						Source:    string(shared.FONBET),
						SportName: string(shared.TENNIS),
						CreatedAt: time.Now(),
					}
				}
			}

		case <-oddsTicker.C:
			matchIds := make([]int64, 0, len(s.tennisData))
			for _, match := range s.tennisData {
				matchIds = append(matchIds, match.Pid)
			}

			for _, matchId := range matchIds {
				go func(matchId int64) {
					eventData, err := s.pAPI.GetMatchODDS(matchId)
					if err != nil {
						s.logger.Error().Err(err).Msg("[Service.RunTennis] error get match odds.")
						return
					}

					eventsParse := make(map[int64]entity.Event)

					for _, eventData := range eventData.Events {
						eventsParse[eventData.Id] = eventData
					}

					periods := make([]entity.ResponsePeriod, 6)

					for i := range periods {
						periods[i].Win1x2 = entity.Win1x2Struct{}
						periods[i].Games = make(map[string]*entity.Win1x2Struct)
						periods[i].Totals = make(map[string]*entity.WinLessMore)
						periods[i].Handicap = make(map[string]*entity.WinHandicap)
						periods[i].FirstTeamTotals = make(map[string]*entity.WinLessMore)
						periods[i].SecondTeamTotals = make(map[string]*entity.WinLessMore)
					}

					for _, eventOdds := range eventData.CustomFactors {
						event, exist := eventsParse[eventOdds.EventId]
						if !exist {
							continue
						}

						periodIndex := 0
						if event.Name == "" {
							periodIndex = 0
						} else if event.Name == "1st set" {
							periodIndex = 1
						} else if event.Name == "2nd set" {
							periodIndex = 2
						} else if event.Name == "3rd set" {
							periodIndex = 3
						} else if event.Name == "4th set" {
							periodIndex = 4
						} else if event.Name == "5th set" {
							periodIndex = 5
						} else {
							continue
						}

						for _, outcome := range eventOdds.Factors {
							// Process 1x2
							if mapping, ok := win1x2Mappings[outcome.FactorId]; ok {
								setWin1x2Value(&periods[periodIndex].Win1x2, mapping.oddType, outcome)
								continue
							}

							// Process games
							if mapping, ok := gamesMapping[outcome.FactorId]; ok {
								ensureMapEntry(periods[periodIndex].Games, normalizeLine(outcome.Line))
								setWin1x2Value(periods[periodIndex].Games[normalizeLine(outcome.Line)], mapping.oddType, outcome)
								continue
							}

							// Process totals
							if mapping, ok := totalsMappings[outcome.FactorId]; ok {
								ensureMapEntry(periods[periodIndex].Totals, normalizeLine(outcome.Line))
								setTotalValue(periods[periodIndex].Totals[normalizeLine(outcome.Line)], mapping.oddType, outcome)
								continue
							}

							// Process team totals
							if mapping, ok := teamTotalsMappings[outcome.FactorId]; ok {
								var totalsMap map[string]*entity.WinLessMore
								if mapping.team == "first" {
									totalsMap = periods[periodIndex].FirstTeamTotals
								} else {
									totalsMap = periods[periodIndex].SecondTeamTotals
								}

								ensureMapEntry(totalsMap, normalizeLine(outcome.Line))
								setTotalValue(totalsMap[normalizeLine(outcome.Line)], mapping.oddType, outcome)
								continue
							}

							// Process handicaps
							processHandicap(&periods, periodIndex, outcome)
						}
					}

					s.tennisData[matchId].Periods = periods

					s.sendChan <- *s.tennisData[matchId]
				}(matchId)
			}

		case <-ctx.Done():
			matchTicker.Stop()
			oddsTicker.Stop()
			return
		}
	}
}
