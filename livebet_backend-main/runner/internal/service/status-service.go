package service

import (
	"context"
	"livebets/runner/cmd/config"
	"livebets/runner/internal/api"
	"livebets/runner/internal/entity"
	"livebets/runner/internal/storage"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type StatusService struct {
	cfg              config.StatusConfig
	parserAPI        *api.ParserAPI
	bookmakerStorage *storage.BookmakerStorage
	statusStorage    *storage.StatusStorage
	logger           *zerolog.Logger
}

func NewStatusService(
	cfg config.StatusConfig,
	parserAPI *api.ParserAPI,
	bookmakerStorage *storage.BookmakerStorage,
	statusStorage *storage.StatusStorage,
	logger *zerolog.Logger,
) *StatusService {
	return &StatusService{
		cfg:              cfg,
		parserAPI:        parserAPI,
		bookmakerStorage: bookmakerStorage,
		statusStorage:    statusStorage,
		logger:           logger,
	}
}

func (s *StatusService) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	statusInterval := time.Duration(time.Duration(s.cfg.Interval) * time.Second)
	statusTicker := time.NewTicker(statusInterval)

	for {
		select {
		case <-statusTicker.C:

			bookmakers := s.bookmakerStorage.ReadAll()
			for _, val := range bookmakers {

				status, err := s.parserAPI.GetOnlineMatchData(val.API)
				if err != nil || status != 200 {
					s.statusStorage.SetStatus(entity.StatusBookmaker{
						Name:      val.Name,
						Status:    entity.StatusOFF,
						CreatedAt: time.Now(),
					})
					continue
				}

				s.statusStorage.SetStatus(entity.StatusBookmaker{
					Name:      val.Name,
					Status:    entity.StatusON,
					CreatedAt: time.Now(),
				})

			}

		case <-ctx.Done():
			statusTicker.Stop()
			return
		}
	}
}

func (s *StatusService) GetStatuses(ctx context.Context) []entity.StatusBookmaker {
	var result []entity.StatusBookmaker
	statuses := s.statusStorage.ReadAll()
	for _, val := range statuses {
		result = append(result, val)
	}

	return result
}
