package service

import (
	"context"
	"github.com/rs/zerolog"
	"livebets/parse_lobbet/cmd/config"
	"livebets/parse_lobbet/internal/api"
	"livebets/parse_lobbet/internal/parse"
	"livebets/parse_lobbet/internal/sender"
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

const (
	liveStatus     = 1
	preMatchStatus = 0
)

func (s *Service) RunLive(ctx context.Context, cfg config.APIConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(cfg.Live.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			matches, err := s.api.GetAllMatches(liveStatus)
			if err != nil {
				s.logger.Error().Err(err).Msg("[Service.RunLive] error get all matches.")
				break
			}

			elapsed := time.Since(start)
			s.logger.Info().Msgf("Получено %d матчей за %s", len(matches), elapsed)

			for _, match := range matches {
				responseGame := parse.LiveToResponseGame(*match)

				err = s.sender.SendMessage(ctx, *responseGame)
				if err != nil {
					s.logger.Error().Err(err).Msgf("[Service.RunLive] error send message.")
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) RunPrematch(ctx context.Context, cfg config.APIConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(cfg.Prematch.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			matches, err := s.api.GetAllMatches(preMatchStatus)
			if err != nil {
				s.logger.Error().Err(err).Msg("[Service.RunPrematch] error get all matches.")
				break
			}

			elapsed := time.Since(start)
			s.logger.Info().Msgf("Получено %d матчей за %s", len(matches), elapsed)

			start = time.Now()
			matchCounter := 0
			for _, match := range matches {
				err = s.api.GetMatchOdds(match)
				if err != nil {
					s.logger.Error().Err(err).Msg("[Service.RunPrematch] error get match data.")
					continue
				}

				if len(match.Bets) == 0 {
					continue
				}

				responseGame := parse.PrematchToResponseGame(*match)

				err = s.sender.SendMessage(ctx, *responseGame)
				if err != nil {
					s.logger.Error().Err(err).Msgf("[Service.RunPrematch] error send message.")
				}

				matchCounter++
			}

			elapsed = time.Since(start)
			s.logger.Info().Msgf("В анализатор отправлено %d матчей за %s", matchCounter, elapsed)

		case <-ctx.Done():
			return
		}
	}
}
