package service

import (
	"context"
	"github.com/rs/zerolog"
	"livebets/parse_starcasino/cmd/config"
	"livebets/parse_starcasino/internal/api"
	"livebets/parse_starcasino/internal/entity"
	"livebets/parse_starcasino/internal/parse"
	"livebets/parse_starcasino/internal/sender"
	"livebets/shared"
	"sync"
	"time"
)

type Service struct {
	api    *api.API
	sender *sender.Sender
	logger *zerolog.Logger
}

func New(
	api *api.API,
	sender *sender.Sender,
	logger *zerolog.Logger,
) *Service {
	return &Service{
		api:    api,
		sender: sender,
		logger: logger,
	}
}
func (s *Service) RunLive(ctx context.Context, cfg config.APIConfig, mainWg *sync.WaitGroup) {
	defer mainWg.Done()

	wg := &sync.WaitGroup{}

	if cfg.Football {
		wg.Add(1)
		go s.runLiveSport(ctx, cfg, shared.SOCCER, wg)
	}

	if cfg.Tennis {
		wg.Add(1)
		go s.runLiveSport(ctx, cfg, shared.TENNIS, wg)
	}

	wg.Wait()
}

func (s *Service) runLiveSport(ctx context.Context, cfg config.APIConfig, sportName shared.SportName, wg *sync.WaitGroup) {
	defer wg.Done()

	eventsTicker := time.NewTicker(time.Duration(cfg.Live.EventsInterval) * time.Second)
	defer eventsTicker.Stop()

	oddsTicker := time.NewTicker(time.Duration(cfg.Live.OddsInterval) * time.Second)
	defer oddsTicker.Stop()

	var events []*entity.Event

	for {
		select {
		case <-eventsTicker.C:

			var err error
			events, err = s.api.GetLiveEvents(sportName)
			if err != nil {
				s.logger.Error().Err(err).Msgf("[Service.RunLive] error get all %s live events.", sportName)
				continue
			}

		case <-oddsTicker.C:

			for _, event := range events {

				go func(oneEvent *entity.Event) {
					match, err := s.api.GetLiveOneEvent(oneEvent.Id)
					if err != nil {
						s.logger.Error().Err(err).Msgf("[Service.RunLive] error get one %s event. Name: %s", sportName, oneEvent.Name)
						return
					}

					if match == nil {
						return
					}

					match.Scores = oneEvent.Scores

					responseGame := parse.StarCasinoToResponseGame(*match)
					responseGame.IsLive = true

					err = s.sender.SendMessage(*responseGame)

					if err != nil {
						s.logger.Error().Err(err).Msgf("[Service.RunLive] error send message.")
					}
				}(event)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) RunPreMatch(ctx context.Context, cfg config.APIConfig, mainWg *sync.WaitGroup) {
	defer mainWg.Done()

	wg := &sync.WaitGroup{}

	if cfg.Football {
		wg.Add(1)
		go s.runPreMatchSport(ctx, cfg, shared.SOCCER, wg)
	}

	if cfg.Tennis {
		wg.Add(1)
		go s.runPreMatchSport(ctx, cfg, shared.TENNIS, wg)
	}

	wg.Wait()
}

func (s *Service) runPreMatchSport(ctx context.Context, cfg config.APIConfig, sportName shared.SportName, wg *sync.WaitGroup) {
	defer wg.Done()

	eventsTicker := time.NewTicker(time.Duration(cfg.Prematch.EventsInterval) * time.Second)
	defer eventsTicker.Stop()

	for {
		select {
		case <-eventsTicker.C:

			const eventsPerPage = 100

			pageCount := 1
			for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
				preMatchData, err := s.api.GetPreMatchPage(sportName, pageNumber, eventsPerPage)
				if err != nil {
					s.logger.Error().Err(err).Msgf("[Service.RunLive] error get all %s live events.", sportName)
					continue
				}

				pageCount = preMatchData.PageCount

				gameDataSlice, stopPagination := parse.StarCasinoPreMatchToResponseGames(preMatchData)
				if stopPagination {
					// Обработаем оставшиеся данные на этой странице
					// и перестанем запрашивать следующие страницы
					pageCount = 1
				}

				for _, gameData := range gameDataSlice {

					/*
						if strings.HasPrefix(gameData.HomeName, "sydney") { // responseGame.HomeName == "tristan boyer" {
							fmt.Println("Score:", gameData.HomeScore, ":", gameData.AwayScore, "Match:", gameData.HomeName+" : "+gameData.AwayName)
							for i, period := range gameData.Periods {
								// 1x2
								fmt.Println("   Time", i, "Win1x2:", period.Win1x2.Win1, period.Win1x2.WinNone, period.Win1x2.Win2)
								// Game 1x2
								for line, game := range period.Games {
									fmt.Println("   Time", i, "Game", line, ":", "Win1:", game.Win1, ", Win2:", game.Win2)
								}
								// Totals
								for line, total := range period.Totals {
									fmt.Println("   Time", i, "Total", line, ":", " >", total.WinMore, ", <", total.WinLess)
								}
								// Totals - teams
								// FirstTeamTotals
								for line, total := range period.FirstTeamTotals {
									fmt.Println("   Time", i, "Total Team 1", line, ":", " >", total.WinMore, ", <", total.WinLess)
								}
								// SecondTeamTotals
								for line, total := range period.SecondTeamTotals {
									fmt.Println("   Time", i, "Total Team 2", line, ":", " >", total.WinMore, ", <", total.WinLess)
								}
								// Handicap
								for line, handicap := range period.Handicap {
									fmt.Println("   Time", i, "Handicap", line, ":", handicap.Win1, ",", handicap.Win2)
								}
							}
						}
					*/
					err = s.sender.SendMessage(*gameData)
					if err != nil {
						s.logger.Error().Err(err).Msgf("[Service.RunLive] error send message.")
					}
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
