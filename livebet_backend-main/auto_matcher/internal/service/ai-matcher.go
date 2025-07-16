package service

import (
	"context"
	"encoding/json"
	"fmt"
	"livebets/auto_matcher/cmd/config"
	"livebets/auto_matcher/internal/entity"
	"livebets/auto_matcher/internal/repository"
	"livebets/auto_matcher/pkg/cache"
	"livebets/auto_matcher/pkg/rdbms"
	"livebets/shared"
	"sync"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/rs/zerolog"
)

type AIMatcherService struct {
	txStorage            rdbms.TxStorage[repository.MatchStorage]
	leaguesLiveCache     cache.MemoryCacheInterface[int, entity.League]
	leaguesPrematchCache cache.MemoryCacheInterface[int, entity.League]
	teamsLiveCache       cache.MemoryCacheInterface[int, entity.UnMatchedTeam]
	teamsPrematchCache   cache.MemoryCacheInterface[int, entity.UnMatchedTeam]
	onlineMatcherService *OnlineMatcherService
	handMatchService     *HandMatcherService
	cfg                  config.AIMatcherConfig
	logger               *zerolog.Logger
	client               *anthropic.Client
}

func NewAIMatcherService(
	txStorage rdbms.TxStorage[repository.MatchStorage],
	leaguesLiveCache cache.MemoryCacheInterface[int, entity.League],
	leaguesPrematchCache cache.MemoryCacheInterface[int, entity.League],
	teamsLiveCache cache.MemoryCacheInterface[int, entity.UnMatchedTeam],
	teamsPrematchCache cache.MemoryCacheInterface[int, entity.UnMatchedTeam],
	onlineMatcherService *OnlineMatcherService,
	handMatchService *HandMatcherService,
	cfg config.AIMatcherConfig,
	logger *zerolog.Logger,
) *AIMatcherService {
	client := anthropic.NewClient(cfg.ApiKey, anthropic.WithBetaVersion(anthropic.BetaPromptCaching20240731))
	return &AIMatcherService{
		txStorage:            txStorage,
		leaguesLiveCache:     leaguesLiveCache,
		leaguesPrematchCache: leaguesPrematchCache,
		teamsLiveCache:       teamsLiveCache,
		teamsPrematchCache:   teamsPrematchCache,
		onlineMatcherService: onlineMatcherService,
		handMatchService:     handMatchService,
		cfg:                  cfg,
		logger:               logger,
		client:               client,
	}
}

func (s *AIMatcherService) Run(ctx context.Context, bookmakerPairs map[int64][2]string, wg *sync.WaitGroup) {
	defer wg.Done()

	prematchTicker := time.NewTicker(time.Duration(s.cfg.PrematchInterval) * time.Second)
	defer prematchTicker.Stop()

	liveTicker := time.NewTicker(time.Duration(s.cfg.LiveInterval) * time.Second)
	defer liveTicker.Stop()

	for {
		select {
		case <-liveTicker.C:
			s.logger.Info().Msg("[AIMatcherService.Run] start AI matching (live)")

			sports, err := s.txStorage.Storage().GetSports(ctx)
			if err != nil {
				s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get sports error")
				continue
			}

			for _, sport := range sports {
				for _, bookmakerPair := range bookmakerPairs {
					// Matching leagues
					unmatchedLeagues, err := s.onlineMatcherService.GetOnlineUnmatchLeagues(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get online unmatch leagues error")
						continue
					}
					s.logger.Info().Msgf("[LIVE] Got %d leagues for %s and %s", len(unmatchedLeagues), bookmakerPair[0], bookmakerPair[1])

					// Collecting leagues except cached leagues
					var leagues1, leagues2 []entity.League
					for _, league := range unmatchedLeagues {
						if _, ok := s.leaguesLiveCache.Read(int(league.ID)); ok {
							continue
						}

						if league.BookmakerName == bookmakerPair[0] {
							leagues1 = append(leagues1, league)
						} else {
							leagues2 = append(leagues2, league)
						}
					}
					s.logger.Info().Msgf("[LIVE] Collected %d leagues for %s and %d leagues for %s", len(leagues1), bookmakerPair[0], len(leagues2), bookmakerPair[1])

					if len(leagues1) >= 10 || len(leagues2) >= 10 && leagues1 != nil && leagues2 != nil {
						matchedLeagueIDs := make(map[int]struct{})
						matchedPairs, _ := s.sendLeaguesToClaude(leagues1, leagues2, bookmakerPair[1])
						s.logger.Info().Msgf("[LIVE] Pairs %d from Claude", len(matchedPairs))
						for _, pair := range matchedPairs {
							matchedLeagueIDs[int(pair.BK1LeagueID)] = struct{}{}
							matchedLeagueIDs[int(pair.BK2LeagueID)] = struct{}{}

							_, err := s.handMatchService.CreateLeaguesPair(ctx, pair.BK1LeagueID, pair.BK2LeagueID)
							if err != nil {
								s.logger.Error().Err(err).Msgf("[AIMatcherService.Run] create leagues pair error %d and %d", pair.BK1LeagueID, pair.BK2LeagueID)
								continue
							}
							s.logger.Info().Msg("[AIMatcherService.Run] create leagues pair")
						}
						// Caching leagues
						for _, league := range leagues1 {
							if _, exists := matchedLeagueIDs[int(league.ID)]; !exists {
								s.leaguesLiveCache.Write(int(league.ID), league)
							}
						}
						for _, league := range leagues2 {
							if _, exists := matchedLeagueIDs[int(league.ID)]; !exists {
								s.leaguesLiveCache.Write(int(league.ID), league)
							}
						}
					}

					// Matching teams
					matchedLeagues, err := s.onlineMatcherService.GetOnlineUnmatchTeams(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get online unmatch teams error")
						continue
					}
					s.logger.Info().Msgf("[LIVE] Got %d leagues for %s and %s", len(matchedLeagues), bookmakerPair[0], bookmakerPair[1])

					// Collecting teams except cached teams
					var teams1, teams2 []entity.UnMatchedTeam
					for _, matchedLeague := range matchedLeagues {
						var leagueTeams1, leagueTeams2 []entity.UnMatchedTeam

						for _, team1 := range matchedLeague.TeamsFirst {
							if _, ok := s.teamsLiveCache.Read(int(team1.TeamID)); ok {
								continue
							}
							leagueTeams1 = append(leagueTeams1, entity.UnMatchedTeam{
								TeamID:   team1.TeamID,
								TeamName: fmt.Sprintf("%s | %s", matchedLeague.LeagueNameFirst, team1.TeamName),
							})
						}
						for _, team2 := range matchedLeague.TeamsSecond {
							if _, ok := s.teamsLiveCache.Read(int(team2.TeamID)); ok {
								continue
							}
							leagueTeams2 = append(leagueTeams2, entity.UnMatchedTeam{
								TeamID:   team2.TeamID,
								TeamName: fmt.Sprintf("%s | %s", matchedLeague.LeagueNameSecond, team2.TeamName),
							})
						}

						if len(leagueTeams1) >= 10 || len(leagueTeams2) >= 10 && leagueTeams1 != nil && leagueTeams2 != nil {
							teams1 = append(teams1, leagueTeams1...)
							teams2 = append(teams2, leagueTeams2...)
						}
					}
					s.logger.Info().Msgf("[LIVE] Collected %d teams for %s and %d teams for %s", len(teams1), bookmakerPair[0], len(teams2), bookmakerPair[1])

					// Sending teams to Claude and creating pairs
					matchedTeamIDs := make(map[int]struct{})
					matchedTeams, _ := s.sendTeamsToClaude(teams1, teams2, bookmakerPair[1])
					s.logger.Info().Msgf("[LIVE] Pairs %d from Claude", len(matchedTeams))
					for _, pair := range matchedTeams {
						matchedTeamIDs[int(pair.BK1TeamID)] = struct{}{}
						matchedTeamIDs[int(pair.BK2TeamID)] = struct{}{}

						_, err := s.handMatchService.CreateTeamsPair(ctx, pair.BK1TeamID, pair.BK2TeamID)
						if err != nil {
							s.logger.Error().Err(err).Msg("[AIMatcherService.Run] create teams pair error")
							continue
						}
						s.logger.Info().Msg("[AIMatcherService.Run] create teams pair")
					}
					// Caching teams
					for _, team := range teams1 {
						if _, exists := matchedTeamIDs[int(team.TeamID)]; !exists {
							s.teamsLiveCache.Write(int(team.TeamID), team)
						}
					}
					for _, team := range teams2 {
						if _, exists := matchedTeamIDs[int(team.TeamID)]; !exists {
							s.teamsLiveCache.Write(int(team.TeamID), team)
						}
					}
				}
			}

		case <-prematchTicker.C:
			s.logger.Info().Msg("[AIMatcherService.Run] start AI matching (prematch)")

			sports, err := s.txStorage.Storage().GetSports(ctx)
			if err != nil {
				s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get sports error")
				return
			}

			for _, sport := range sports {
				s.logger.Info().Msgf("[PREMATCH] Sport %s", sport)
				for _, bookmakerPair := range bookmakerPairs {
					// Matching leagues
					unmatchedLeagues, err := s.onlineMatcherService.GetOnlineUnmatchLeaguesPrematch(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get online unmatch leagues error")
						continue
					}

					// Collecting leagues except cached leagues
					var leagues1, leagues2 []entity.League
					for _, league := range unmatchedLeagues {
						if _, ok := s.leaguesPrematchCache.Read(int(league.ID)); ok {
							continue
						}

						if league.BookmakerName == bookmakerPair[0] {
							leagues1 = append(leagues1, league)
						} else if league.BookmakerName == bookmakerPair[1] {
							leagues2 = append(leagues2, league)
						}
					}

					if len(leagues1) >= 10 || len(leagues2) >= 10 && leagues1 != nil && leagues2 != nil {
						batchedLeagues1 := splitIntoBatches(leagues1)
						batchedLeagues2 := splitIntoBatches(leagues2)

						for i, leagues1 := range batchedLeagues1 {
							if i >= len(batchedLeagues2) {
								break
							}
							leagues2 := batchedLeagues2[i]

							matchedLeaguesIDs := make(map[int]struct{})
							matchedPairs, err := s.sendLeaguesToClaude(leagues1, leagues2, bookmakerPair[1])
							if err != nil {
								s.logger.Error().Err(err).Msgf("[AIMatcherService.Run] send leagues to Claude error")
								continue
							}
							for _, pair := range matchedPairs {
								matchedLeaguesIDs[int(pair.BK1LeagueID)] = struct{}{}
								matchedLeaguesIDs[int(pair.BK2LeagueID)] = struct{}{}

								_, err := s.handMatchService.CreateLeaguesPair(ctx, pair.BK1LeagueID, pair.BK2LeagueID)
								if err != nil {
									s.logger.Error().Err(err).Msgf("[AIMatcherService.Run] create leagues pair error (%d, %d)", pair.BK1LeagueID, pair.BK2LeagueID)
									continue
								}
								s.logger.Info().Msgf("[AIMatcherService.Run] leagues pair %d -> %d", pair.BK1LeagueID, pair.BK2LeagueID)
							}
						}
					}

					// Matching teams
					matchedLeagues, err := s.onlineMatcherService.GetOnlineUnmatchTeamsPrematch(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						s.logger.Error().Err(err).Msg("[AIMatcherService.Run] get online unmatch teams error")
						continue
					}

					var teams1, teams2 []entity.UnMatchedTeam
					for _, matchedLeague := range matchedLeagues {
						var leagueTeams1, leagueTeams2 []entity.UnMatchedTeam

						for _, team1 := range matchedLeague.TeamsFirst {
							if _, ok := s.teamsPrematchCache.Read(int(team1.TeamID)); ok {
								continue
							}
							leagueTeams1 = append(leagueTeams1, entity.UnMatchedTeam{
								TeamID:   team1.TeamID,
								TeamName: fmt.Sprintf("%s | %s", matchedLeague.LeagueNameFirst, team1.TeamName),
							})
						}

						for _, team2 := range matchedLeague.TeamsSecond {
							if _, ok := s.teamsPrematchCache.Read(int(team2.TeamID)); ok {
								continue
							}
							leagueTeams2 = append(leagueTeams2, entity.UnMatchedTeam{
								TeamID:   team2.TeamID,
								TeamName: fmt.Sprintf("%s | %s", matchedLeague.LeagueNameSecond, team2.TeamName),
							})
						}

						if len(leagueTeams1) >= 10 || len(leagueTeams2) >= 10 && leagueTeams1 != nil && leagueTeams2 != nil {
							teams1 = append(teams1, leagueTeams1...)
							teams2 = append(teams2, leagueTeams2...)
						}
					}

					// Send teams to Claude and create pairs
					batchedTeams1 := splitIntoBatches(teams1)
					batchedTeams2 := splitIntoBatches(teams2)

					for i, teams1 := range batchedTeams1 {
						if i >= len(batchedTeams2) {
							break
						}
						teams2 := batchedTeams2[i]

						matchedTeamIDs := make(map[int]struct{})
						matchedTeams, err := s.sendTeamsToClaude(teams1, teams2, bookmakerPair[1])
						if err != nil {
							s.logger.Error().Err(err).Msg("[AIMatcherService.Run] send teams to Claude error")
							continue
						}
						for _, pair := range matchedTeams {
							matchedTeamIDs[int(pair.BK1TeamID)] = struct{}{}
							matchedTeamIDs[int(pair.BK2TeamID)] = struct{}{}

							_, err := s.handMatchService.CreateTeamsPair(ctx, pair.BK1TeamID, pair.BK2TeamID)
							if err != nil {
								s.logger.Error().Err(err).Msg("[AIMatcherService.Run] create teams pair error")
							}
							s.logger.Info().Msgf("[AIMatcherService.Run] create teams pair %d -> %d", pair.BK1TeamID, pair.BK2TeamID)
						}
						// Caching teams
						for _, team := range teams1 {
							if _, exists := matchedTeamIDs[int(team.TeamID)]; !exists {
								s.teamsPrematchCache.Write(int(team.TeamID), team)
							}
						}
						for _, team := range teams2 {
							if _, exists := matchedTeamIDs[int(team.TeamID)]; !exists {
								s.teamsPrematchCache.Write(int(team.TeamID), team)
							}
						}
					}
				}
			}

		case <-ctx.Done():
			liveTicker.Stop()
			prematchTicker.Stop()
			return
		}
	}
}

func (s *AIMatcherService) sendToClaude(query, systemMsg, bookmakerMsg string) (string, error) {
	resp, err := s.client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.Model(s.cfg.Claude.Model),
		MultiSystem: []anthropic.MessageSystemPart{
			{
				Type: "text",
				Text: systemMsg,
				CacheControl: &anthropic.MessageCacheControl{
					Type: anthropic.CacheControlTypeEphemeral,
				},
			},
			{
				Type: "text",
				Text: bookmakerMsg,
				CacheControl: &anthropic.MessageCacheControl{
					Type: anthropic.CacheControlTypeEphemeral,
				},
			},
		},
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(query),
		},
		MaxTokens: s.cfg.Claude.MaxTokens,
	})
	if err != nil {
		return "", err
	}

	return *resp.Content[0].Text, nil
}

func (s *AIMatcherService) sendLeaguesToClaude(leagues1, leagues2 []entity.League, bookmaker string) ([]entity.ResponsePairLeague, error) {
	bookmakerMsg := s.GetBookmakerMessage(bookmaker)

	req := entity.RequestLeagues{
		BK1Leagues: leagues1,
		BK2Leagues: leagues2,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Msgf("Sending leagues to Claude: %s", string(reqBytes))

	answer, err := s.sendToClaude(string(reqBytes), s.cfg.Claude.Messages.SystemLeaguesMsg, bookmakerMsg)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Msgf("Received leagues from Claude: %s", answer)

	var pairedLeagues []entity.ResponsePairLeague

	if err = json.Unmarshal([]byte(answer), &pairedLeagues); err != nil {
		return nil, err
	}

	return pairedLeagues, nil
}

func (s *AIMatcherService) sendTeamsToClaude(teams1, teams2 []entity.UnMatchedTeam, bookmaker string) ([]entity.ResponsePairTeam, error) {
	bookmakerMsg := s.GetBookmakerMessage(bookmaker)

	req := entity.RequestTeams{
		BK1Teams: teams1,
		BK2Teams: teams2,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Msgf("Sending teams to Claude: %s", string(reqBytes))

	answer, err := s.sendToClaude(string(reqBytes), s.cfg.Claude.Messages.SystemTeamsMsg, bookmakerMsg)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Msgf("Received teams from Claude: %s", answer)

	var pairedTeams []entity.ResponsePairTeam

	if err = json.Unmarshal([]byte(answer), &pairedTeams); err != nil {
		return nil, err
	}

	return pairedTeams, nil
}

func (s *AIMatcherService) GetBookmakerMessage(bookmaker string) string {
	if bookmaker == string(shared.BETCENTER) {
		return s.cfg.Claude.Messages.BetcenterMsg
	}

	if bookmaker == string(shared.FONBET) {
		return s.cfg.Claude.Messages.FonbetMsg
	}

	if bookmaker == string(shared.LADBROKES) {
		return s.cfg.Claude.Messages.LadbrokesMsg
	}

	if bookmaker == string(shared.LOBBET) {
		return s.cfg.Claude.Messages.LobbetMsg
	}

	if bookmaker == string(shared.MAXBET) {
		return s.cfg.Claude.Messages.MaxbetMsg
	}

	if bookmaker == string(shared.SANSABET) {
		return s.cfg.Claude.Messages.SansabetMsg
	}

	if bookmaker == string(shared.SBBET) {
		return s.cfg.Claude.Messages.SbbetMsg
	}

	if bookmaker == string(shared.STARCASINO) {
		return s.cfg.Claude.Messages.StarCasinoMsg
	}

	if bookmaker == string(shared.UNIBET) {
		return s.cfg.Claude.Messages.UnibetMsg
	}

	return "message"
}

func splitIntoBatches[T any](data []T) [][]T {
	batchSize := 100
	batches := make([][]T, 0)

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batches = append(batches, data[i:end])
	}

	return batches
}
