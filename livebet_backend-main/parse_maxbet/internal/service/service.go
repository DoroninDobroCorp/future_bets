package service

import (
	"context"
	"fmt"
	"livebets/parse_maxbet/cmd/config"
	"livebets/parse_maxbet/internal/api"
	"livebets/parse_maxbet/internal/entity"
	"livebets/parse_maxbet/internal/parse"
	"livebets/parse_maxbet/internal/sender"
	"livebets/shared"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var SPORTS = []string{
	"S", // football
	"T", // tennis
	// "BB", // basketball
}

type GeneralService struct {
	api    *api.API
	sender *sender.Sender
	logger *zerolog.Logger
}

func New(
	api *api.API,
	sender *sender.Sender,
	logger *zerolog.Logger,
) *GeneralService {
	return &GeneralService{
		api:    api,
		sender: sender,
		logger: logger,
	}
}

func (s *GeneralService) RunLive(ctx context.Context, cfg config.APIConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(cfg.Live.IntervalEvents) * time.Second)
	defer ticker.Stop()

	sportFilter := makeSportFilter(cfg.SportConfig)

	for {
		select {
		case <-ticker.C:
			datas, err := s.api.GetLiveData()
			if err != nil {
				s.logger.Error().Err(err).Msg("[Service.RunLive] error get all leagues.")
				continue
			}

			var responseGames []shared.GameData

			for _, data := range datas {
				for _, event := range data.LiveEvents {
					if _, ok := sportFilter[event.SportId]; !ok {
						continue
					}

					var sportName shared.SportName
					var leagueName, homeName, awayName string
					var periods []shared.PeriodData
					if event.SportId == "S" {
						sportName = shared.SOCCER
						leagueName = parse.NormalizeFootballLeague(event.LeagueName)
						homeName = parse.NormalizeFootballTeam(event.HomeName)
						awayName = parse.NormalizeFootballTeam(event.AwayName)
						periods = make([]shared.PeriodData, 3)
					} else if event.SportId == "T" {
						sportName = shared.TENNIS
						leagueName = parse.NormalizeTennisLeague(event.LeagueName)
						homeName = parse.NormalizeTennisTeam(event.HomeName)
						awayName = parse.NormalizeTennisTeam(event.AwayName)
						periods = make([]shared.PeriodData, 6)
					} else if event.SportId == "BB" {
						sportName = shared.BASKETBALL
						leagueName = parse.NormalizeBasketballLeague(event.LeagueName)
						homeName = parse.NormalizeBasketballTeam(event.HomeName)
						awayName = parse.NormalizeBasketballTeam(event.AwayName)
						periods = make([]shared.PeriodData, 4)
					}
					homeScore, awayScore := findScores(data.LiveResults, event.Id)

					responseGame := shared.GameData{
						Pid:        event.Id,
						LeagueName: leagueName,
						HomeName:   homeName,
						AwayName:   awayName,
						MatchId:    fmt.Sprintf("%d", event.Id),
						IsLive:     true,

						HomeScore: homeScore,
						AwayScore: awayScore,

						Source:    shared.MAXBET,
						SportName: sportName,
					}

					for i := range periods {
						periods[i].Win1x2 = shared.Win1x2Struct{}
						periods[i].Games = make(map[string]*shared.Win1x2Struct)
						periods[i].Totals = make(map[string]*shared.WinLessMore)
						periods[i].Handicap = make(map[string]*shared.WinHandicap)
						periods[i].FirstTeamTotals = make(map[string]*shared.WinLessMore)
						periods[i].SecondTeamTotals = make(map[string]*shared.WinLessMore)
					}

					eventBets := findBets(data.LiveBets, event.Id)
					if len(eventBets) == 0 {
						break
					}

					for _, bet := range eventBets {
						for key, _ := range bet.Coefs {
							if oddMap, ok := win1X2Mappings[key]; ok {
								processWin1x2(&periods[oddMap.period].Win1x2, bet)
								break
							}

							if oddMap, ok := gamesMappings[key]; ok {
								processGamesLive(periods[oddMap.period].Games, bet)
								break
							}

							if oddMap, ok := totalsMappings[key]; ok {
								processTotalLive(periods[oddMap.period].Totals, bet, false)
								break
							}

							if oddMap, ok := teamTotalsMappings[key]; ok {
								if strings.HasPrefix(oddMap.team, "1") {
									processTotalLive(periods[oddMap.period].FirstTeamTotals, bet, true)
								} else {
									processTotalLive(periods[oddMap.period].SecondTeamTotals, bet, true)
								}
								break
							}

							if oddMap, ok := handicapsMappings[key]; ok {
								processHandicapLive(periods[oddMap.period].Handicap, bet)
								break
							}

							break
						}
					}

					responseGame.Periods = periods

					responseGames = append(responseGames, responseGame)

					err := s.sender.SendMessage(ctx, responseGame)
					if err != nil {
						s.logger.Error().Err(err).Msgf("[Service.RunLive] error send message.")
					}
				}
			}

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *GeneralService) RunPrematch(ctx context.Context, cfg config.APIConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(cfg.Prematch.IntervalEvents) * time.Second)
	defer ticker.Stop()

	//sportFilter := makeSportFilter(cfg.SportConfig)

	for {
		select {
		case <-ticker.C:

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func makeSportFilter(sportConfig config.SportConfig) map[string]shared.SportName {
	filter := make(map[string]shared.SportName)

	if sportConfig.Football {
		filter["S"] = shared.SOCCER
	}

	if sportConfig.Tennis {
		filter["T"] = shared.TENNIS
	}

	if sportConfig.Basketball {
		filter["BB"] = shared.BASKETBALL
	}

	return filter
}

func findScores(results []entity.EventResult, eventId int64) (float64, float64) {
	for _, result := range results {
		if result.EventId == eventId {
			return float64(result.HomeScore.FullTime), float64(result.AwayScore.FullTime)
		}
	}
	return 0, 0
}

func findBets(bets []entity.Bet, eventId int64) []entity.Bet {
	var result []entity.Bet

	for _, bet := range bets {
		if bet.EventId == eventId {
			result = append(result, bet)
		}
	}

	return result
}
