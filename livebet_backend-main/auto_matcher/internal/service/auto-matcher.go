package service

import (
	"context"
	"livebets/auto_matcher/cmd/config"
	"livebets/auto_matcher/internal/entity"
	"livebets/auto_matcher/internal/repository"
	"livebets/auto_matcher/pkg/cache"
	"livebets/auto_matcher/pkg/rdbms"
	"livebets/auto_matcher/pkg/utils"
	"strings"
	"sync"
	"time"

	fuzz "github.com/paul-mannino/go-fuzzywuzzy"
	"github.com/rs/zerolog"
)

const (
	MATCHPERCENT float64 = 72
	PERCENT float64 = 67
)

var (
	symbols = []string{"-", ",", " "}
)

type AutoMatcherService struct {
	txStorage             rdbms.TxStorage[repository.MatchStorage]
	leagueCandidatesCache cache.MemoryCacheInterface[string, entity.LeagueCandidatePair]
	teamCandidatesCache   cache.MemoryCacheInterface[string, entity.TeamCandidatePair]
	handMatchService      *HandMatcherService
	logger                *zerolog.Logger
}

func NewAutoMatcherService(
	txStorage rdbms.TxStorage[repository.MatchStorage],
	leagueCandidatesCache cache.MemoryCacheInterface[string, entity.LeagueCandidatePair],
	teamCandidatesCache cache.MemoryCacheInterface[string, entity.TeamCandidatePair],
	handMatchService *HandMatcherService,
	logger *zerolog.Logger,
) *AutoMatcherService {
	return &AutoMatcherService{
		txStorage:             txStorage,
		leagueCandidatesCache: leagueCandidatesCache,
		teamCandidatesCache:   teamCandidatesCache,
		handMatchService:      handMatchService,
		logger:                logger,
	}
}

func (a *AutoMatcherService) GetLeagueCandidates() map[string]entity.LeagueCandidatePair {
	return a.leagueCandidatesCache.ReadAll()
}

func (a *AutoMatcherService) GetTeamCandidates() map[string]entity.TeamCandidatePair {
	return a.teamCandidatesCache.ReadAll() 
}

func (a *AutoMatcherService) Run(ctx context.Context, cfg config.AutoMatcherConfig, bookmakerPairs map[int64][2]string, wg *sync.WaitGroup) {
	defer wg.Done()

	leaguesTicker := time.NewTicker(time.Duration(cfg.IntervalMatchingLeagues) * time.Second)
	defer leaguesTicker.Stop()

	teamsTicker := time.NewTicker(time.Duration(cfg.IntervalMatchingTeams) * time.Second)
	defer teamsTicker.Stop()

	for {
		select {
		case <-leaguesTicker.C:

			sports, err := a.txStorage.Storage().GetSports(ctx)
			if err != nil {
				a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] get sports error")
				continue
			}
			a.leagueCandidatesCache.Clean()

			for _, sport := range sports {
				for _, bookmakerPair := range bookmakerPairs {
					unmatchedLeagues, err := a.handMatchService.GetUnMachedLeagues(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] get unmatched leagues error")
						continue
					}
					leagues1, leagues2 := splitLeagues(unmatchedLeagues, bookmakerPair[0], bookmakerPair[1])

					for _, league := range leagues1 {
						foundLeague, similarity := findSimilarLeagues(league, leagues2)

						if similarity == 100 {
							ok, err := a.handMatchService.CreateLeaguesPair(ctx, league.ID, foundLeague.ID)
							if err != nil {
								a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] create leagues pair error")
								continue
							}
							if ok {
								a.logger.Info().Msgf("[AutoMatcherService.Run] successfully created leagues pair (%s - %s - %s): %s -> %s",
									league.BookmakerName, sport, foundLeague.BookmakerName,
									league.LeagueName, foundLeague.LeagueName)
							}
							continue
						} else if similarity >= PERCENT {
							key := utils.GenerateKeyForCandidate(league.ID, foundLeague.ID)

							a.leagueCandidatesCache.Write(key, entity.LeagueCandidatePair{
								First: entity.LeagueCandidate{
									BookmakerName: league.BookmakerName,
									LeagueName:    league.LeagueName,
									LeagueID:      league.ID,
								},
								Second: entity.LeagueCandidate{
									BookmakerName: foundLeague.BookmakerName,
									LeagueName:    foundLeague.LeagueName,
									LeagueID:      foundLeague.ID,
								},
								SportName:  sport,
								Similarity: similarity,
							})
						}
					}
				}
			}

		case <-teamsTicker.C:

			sports, err := a.txStorage.Storage().GetSports(ctx)
			if err != nil {
				a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] get sports error")
			}
			a.teamCandidatesCache.Clean()

			for _, sport := range sports {
				for _, bookmakerPair := range bookmakerPairs {
					unmatchedTeamsByLeagues, err := a.handMatchService.GetUnMatchedTeamsByLeagues(ctx, sport, bookmakerPair[0], bookmakerPair[1])
					if err != nil {
						a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] get unmatched teams error")
						continue
					}

					for _, unmatchedTeams := range unmatchedTeamsByLeagues {
						teams1, teams2 := unmatchedTeams.TeamsFirst, unmatchedTeams.TeamsSecond

						for _, team := range teams1 {
							foundTeam, similarity := findSimilarTeams(team, teams2)

							if similarity > MATCHPERCENT {
								ok, err := a.handMatchService.CreateTeamsPair(ctx, team.TeamID, foundTeam.TeamID)
								if err != nil {
									a.logger.Error().Err(err).Msg("[AutoMatcherService.Run] get unmatched teams error")
									continue
								}
								if ok {
									a.logger.Info().Msgf("[AutoMatcherService.Run] successfully created team pair (%s - %s - %s): %s (%s) -> %s (%s)",
										unmatchedTeams.BookmakerNameFirst, sport, unmatchedTeams.BookmakerNameSecond,
										team.TeamName, unmatchedTeams.LeagueNameFirst,
										foundTeam.TeamName, unmatchedTeams.LeagueNameSecond)
								}
								continue
							} else if similarity >= PERCENT {
								key := utils.GenerateKeyForCandidate(team.TeamID, foundTeam.TeamID)

								a.teamCandidatesCache.Write(key, entity.TeamCandidatePair{
									First: entity.TeamCandidate{
										BookmakerName: unmatchedTeams.BookmakerNameFirst,
										LeagueName:    unmatchedTeams.LeagueNameFirst,
										LeagueID:      unmatchedTeams.LeagueIDFirst,
										TeamName:      team.TeamName,
										TeamID:        team.TeamID,
									},
									Second: entity.TeamCandidate{
										BookmakerName: unmatchedTeams.BookmakerNameSecond,
										LeagueName:    unmatchedTeams.LeagueNameSecond,
										LeagueID:      unmatchedTeams.LeagueIDSecond,
										TeamName:      foundTeam.TeamName,
										TeamID:        foundTeam.TeamID,
									},
									SportName:  sport,
									Similarity: similarity,
								})
							}
						}
					}
				}
			}

		case <-ctx.Done():
			leaguesTicker.Stop()
			teamsTicker.Stop()
			return
		}
	}
}

func splitLeagues(leagues []entity.League, firstBookmakerName, secondBookmakerName string) ([]entity.League, []entity.League) {
	var firstBookmakerLeagues []entity.League
	var secondBookmakerLeagues []entity.League

	for _, league := range leagues {
		if league.BookmakerName == firstBookmakerName {
			firstBookmakerLeagues = append(firstBookmakerLeagues, league)
		} else if league.BookmakerName == secondBookmakerName {
			secondBookmakerLeagues = append(secondBookmakerLeagues, league)
		}
	}

	return firstBookmakerLeagues, secondBookmakerLeagues
}

func findSimilarLeagues(league entity.League, otherLeagues []entity.League) (entity.League, float64) {
	var foundLeague entity.League
	var bestSimilarity float64

	for _, otherLeague := range otherLeagues {
		similatiry := compareLeagues(league.LeagueName, otherLeague.LeagueName)
		if similatiry > bestSimilarity {
			bestSimilarity = similatiry
			foundLeague = otherLeague
		}
	}

	if bestSimilarity >= 0 {
		return foundLeague, bestSimilarity
	}

	return entity.League{}, 0
}

func findSimilarTeams(team entity.UnMatchedTeam, otherTeams []entity.UnMatchedTeam) (entity.UnMatchedTeam, float64) {
	var foundTeam entity.UnMatchedTeam
	var bestSimilarity float64

	for _, otherTeam := range otherTeams {
		similatiry := compareTeams(team.TeamName, otherTeam.TeamName)
		if similatiry > bestSimilarity {
			bestSimilarity = similatiry
			foundTeam = otherTeam
		}
	}

	if bestSimilarity >= 0 {
		return foundTeam, bestSimilarity
	}

	return entity.UnMatchedTeam{}, 0
}

func compareLeagues(league1, league2 string) float64 {
	league1 = strings.ToLower(strings.TrimSpace(league1))
	league2 = strings.ToLower(strings.TrimSpace(league2))
	for _, symbol := range symbols {
		league1 = strings.ReplaceAll(league1, symbol, "")
		league2 = strings.ReplaceAll(league2, symbol, "")
	}

	if league1 == league2 {
		return 100.0
	}

	return float64(fuzz.Ratio(league1, league2))
}

func compareTeams(team1, team2 string) float64 {
	team1 = strings.ToLower(strings.TrimSpace(team1))
	team2 = strings.ToLower(strings.TrimSpace(team2))
	splited_team1, splited_team2 := strings.Split(team1, " "), strings.Split(team2, " ")
	var filtered1, filtered2 []string

	for _, part := range splited_team1 {
		if len(part) > 3 {
			filtered1 = append(filtered1, part)
		}
	}
	for _, part := range splited_team2 {
		if len(part) > 3 {
			filtered2 = append(filtered2, part)
		}
	}

	team1, team2 = strings.Join(filtered1, " "), strings.Join(filtered2, " ")

	if team1 == team2 {
		return 100.0
	}

	return float64(fuzz.Ratio(team1, team2))
}
