package service

import (
	"context"
	"fmt"
	"livebets/parse_sansabet/cmd/config"
	"livebets/parse_sansabet/internal/api"
	"livebets/parse_sansabet/internal/entity"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type GeneralService struct {
	sAPI     *api.SansabetAPI
	sendChan chan<- entity.ResponseGame
	data     map[int64]*entity.ResponseGame
	logger   *zerolog.Logger
}

func NewGeneralService(
	sAPI *api.SansabetAPI,
	sendChan chan<- entity.ResponseGame,
	logger *zerolog.Logger,
) *GeneralService {
	data := make(map[int64]*entity.ResponseGame)
	return &GeneralService{
		sAPI:     sAPI,
		sendChan: sendChan,
		data:     data,
		logger:   logger,
	}
}

func (s *GeneralService) Run(ctx context.Context, cfg config.SansabetConfig, sportID entity.Sport, wg *sync.WaitGroup) {
	defer wg.Done()
	matchInterval := time.Duration(time.Duration(cfg.IntervalMatch) * time.Second)
	matchTicker := time.NewTicker(matchInterval)

	oddsInterval := time.Duration(time.Duration(cfg.IntervalODDS) * time.Second)
	oddsTicker := time.NewTicker(oddsInterval)

	for {
		select {
		case <-matchTicker.C:
			events, err := s.sAPI.GetAllMatches()
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Service.Run] error get all matches.")
			}

			for _, event := range *events {

				// Фильтр по виду спорта
				if event.H.SportId != string(sportID) {
					continue
				}

				// Фильтр только лайв
				if event.H.MS != "IP" {
					continue
				}

				splitedName := strings.Split(event.H.MatchName, " : ")
				if len(splitedName) != 2 {
					continue
				}
				homeName, awayName := splitedName[0], splitedName[1]

				s.data[event.H.ID] = &entity.ResponseGame{
					Pid:        event.H.ID,
					LeagueName: event.H.LeagueName,
					HomeName:   homeName,
					AwayName:   awayName,
					MatchId:    fmt.Sprintf("%d", event.H.ID),
					Raw: entity.EventRaw{
						MatchName: event.H.MatchName,
					},
				}
			}
		case <-oddsTicker.C:
			matchIds := make([]int64, 0, len(s.data))
			for _, match := range s.data {
				matchIds = append(matchIds, match.Pid)
			}

			eventsOdds, err := s.sAPI.GetAllMatchesODDS(matchIds)
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Service.Run] error get odds for all matches.")
			}

			for _, event := range *eventsOdds {
				// If not exist match
				if _, ok := s.data[event.H.ID]; !ok {
					continue
				}

				// Setting sport
				var sport entity.SportName
				if event.H.SportId == string(entity.FootballID) {
					sport = entity.SportSoccer
				} else if event.H.SportId == string(entity.TennisID) {
					sport = entity.SportTennis
				}

				// Add score
				scores := strings.Split(event.R["G"].(string), "-")
				homeScore, _ := strconv.ParseFloat(scores[0], 64)
				awayScore, _ := strconv.ParseFloat(scores[1], 64)

				s.data[event.H.ID].HomeScore = homeScore
				s.data[event.H.ID].AwayScore = awayScore

				// Parse outcomes
				var periods []entity.ResponsePeriod
				if sport == entity.SportSoccer {
					periods = make([]entity.ResponsePeriod, 3)
				} else if sport == entity.SportTennis {
					periods = make([]entity.ResponsePeriod, 6)
				}

				for i := range periods {
					periods[i] = entity.ResponsePeriod{
						Win1x2:           entity.Win1x2Struct{},
						Games:            make(map[string]*entity.Win1x2Struct),
						Totals:           make(map[string]*entity.WinLessMore),
						Handicap:         make(map[string]*entity.WinHandicap),
						FirstTeamTotals:  make(map[string]*entity.WinLessMore),
						SecondTeamTotals: make(map[string]*entity.WinLessMore),
					}
				}

				for _, outcome := range event.M {
					for _, odd := range outcome.S {
						// Process 1x2
						if mapping, ok := win1x2Mappings[odd.N]; ok {
							setWin1x2Value(&periods[mapping.periodIndex].Win1x2, mapping.oddType, odd.O)
							continue
						}

						// Process games
						if mapping, ok := gamesMappings[odd.N]; ok {
							ensureMapEntry(periods[mapping.periodIndex].Games, outcome.B)
							setWin1x2Value(periods[mapping.periodIndex].Games[outcome.B], mapping.oddType, odd.O)
							continue
						}

						// Process totals
						if mapping, ok := totalsMappings[odd.N]; ok {
							ensureMapEntry(periods[mapping.periodIndex].Totals, outcome.B)
							setTotalValue(periods[mapping.periodIndex].Totals[outcome.B], mapping.oddType, odd.O)
							continue
						}

						// Process team totals
						if mapping, ok := teamTotalsMappings[odd.N]; ok {
							var totalsMap map[string]*entity.WinLessMore
							if mapping.team == "first" {
								totalsMap = periods[mapping.periodIndex].FirstTeamTotals
							} else {
								totalsMap = periods[mapping.periodIndex].SecondTeamTotals
							}
							ensureMapEntry(totalsMap, outcome.B)
							setTotalValue(totalsMap[outcome.B], mapping.oddType, odd.O)
							continue
						}

						// Process handicaps
						processHandicap(odd.N, outcome.B, odd.O, &periods, homeScore, awayScore)
					}
				}

				s.data[event.H.ID].Periods = periods

				// Add config data
				s.data[event.H.ID].SportName = sport
				s.data[event.H.ID].CreatedAt = time.Now()
				s.data[event.H.ID].Source = entity.ParserName
				s.data[event.H.ID].IsLive = true

				s.sendChan <- *s.data[event.H.ID]
			}
		case <-ctx.Done():
			matchTicker.Stop()
			oddsTicker.Stop()
			return
		}
	}
}
