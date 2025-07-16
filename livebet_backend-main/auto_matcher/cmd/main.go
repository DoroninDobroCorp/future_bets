package main

import (
	"context"
	"livebets/auto_matcher/cmd/config"
	"livebets/auto_matcher/internal/api"
	"livebets/auto_matcher/internal/entity"
	"livebets/auto_matcher/internal/handler"
	"livebets/auto_matcher/internal/repository"
	"livebets/auto_matcher/internal/service"
	"livebets/auto_matcher/pkg/cache"
	"livebets/auto_matcher/pkg/pgsql"
	"livebets/auto_matcher/pkg/rdbms"
	"livebets/auto_matcher/pkg/server"
	"livebets/shared"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Auto_Matcher")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	// Connect to postgres
	postgres, err := pgsql.New(appConfig.PostgresConfig.ConnectionString())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	logger.Info().Msg("Connected to Postgres")
	defer postgres.Close()

	wg := &sync.WaitGroup{}

	leagueCandidatesCache := cache.NewMemoryCache[string, entity.LeagueCandidatePair]()
	teamCandidatesCache := cache.NewMemoryCache[string, entity.TeamCandidatePair]()

	leaguesLiveCache := cache.NewMemoryCache[int, entity.League]()
	leaguesPrematchCache := cache.NewMemoryCache[int, entity.League]()
	teamsLiveCache := cache.NewMemoryCache[int, entity.UnMatchedTeam]()
	teamsPrematchCache := cache.NewMemoryCache[int, entity.UnMatchedTeam]()

	analizerAPI := api.NewAnalizerAPI(appConfig.AnalyzerAPI)
	analizerPrematchAPI := api.NewAnalizerPrematchAPI(appConfig.AnalyzerPrematchAPI)

	matchTxStorage := rdbms.NewPgTxStorage(postgres.Pool, repository.NewHandMatchPGStorage)
	handMatchService := service.NewHandMatcherService(matchTxStorage, leagueCandidatesCache, teamCandidatesCache, &logger)
	onlineMatchService := service.NewOnlineMatcherService(matchTxStorage, analizerAPI, analizerPrematchAPI, &logger)

	// TODO: add to config
	bookmakerPairs := map[int64][2]string{
		0: {string(shared.PINNACLE), string(shared.FONBET)},
		1: {string(shared.PINNACLE), string(shared.LADBROKES)},
		2: {string(shared.PINNACLE), string(shared.LOBBET)},
		3: {string(shared.PINNACLE), string(shared.MAXBET)},
		4: {string(shared.PINNACLE), string(shared.SANSABET)},
		5: {string(shared.PINNACLE), string(shared.SBBET)},
		6: {string(shared.PINNACLE), string(shared.STARCASINO)},
		7: {string(shared.PINNACLE), string(shared.UNIBET)},
		8: {string(shared.PINNACLE), string(shared.SERGE)},
	}

	autoMatcherService := service.NewAutoMatcherService(matchTxStorage, leagueCandidatesCache, teamCandidatesCache, handMatchService, &logger)
	wg.Add(1)
	go autoMatcherService.Run(ctx, appConfig.AutoMatcherConfig, bookmakerPairs, wg)

	aiMatcherService := service.NewAIMatcherService(matchTxStorage, leaguesLiveCache, leaguesPrematchCache, teamsLiveCache, teamsPrematchCache, onlineMatchService, handMatchService, appConfig.AIMatcherConfig, &logger)
	wg.Add(1)
	go aiMatcherService.Run(ctx, bookmakerPairs, wg)

	handlers := handler.NewHandler(handMatchService, autoMatcherService, onlineMatchService)

	srv := new(server.Server)
	go func() {
		logger.Info().Msgf("starting server on port = %s", "7001")
		if err := srv.Run("7001", handlers.InitRoutes()); err != nil {
			logger.Error().Err(err).Msg("error occured while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancelFunc()
	wg.Wait()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("error occured on server shutting down")
	}
	logger.Info().Msg(">> Stopping Auto_Matcher")
}
