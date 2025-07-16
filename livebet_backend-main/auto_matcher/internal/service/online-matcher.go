package service

import (
	"context"
	"fmt"
	"livebets/auto_matcher/internal/api"
	"livebets/auto_matcher/internal/entity"
	"livebets/auto_matcher/internal/repository"
	"livebets/auto_matcher/pkg/rdbms"

	"github.com/rs/zerolog"
)

type OnlineMatcherService struct {
	txStorage           rdbms.TxStorage[repository.MatchStorage]
	analyzerAPI         *api.AnalizerAPI
	analyzerPrematchAPI *api.AnalizerPrematchAPI
	logger              *zerolog.Logger
}

func NewOnlineMatcherService(
	txStorage rdbms.TxStorage[repository.MatchStorage],
	analyzerAPI *api.AnalizerAPI,
	analyzerPrematchAPI *api.AnalizerPrematchAPI,
	logger *zerolog.Logger,
) *OnlineMatcherService {
	return &OnlineMatcherService{
		txStorage:           txStorage,
		analyzerAPI:         analyzerAPI,
		analyzerPrematchAPI: analyzerPrematchAPI,
		logger:              logger,
	}
}

func (o *OnlineMatcherService) GetOnlineUnmatchLeagues(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) ([]entity.League, error) {
	// Get match data from analyzer
	matchData, err := o.analyzerAPI.GetOnlineMatchData()
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get match data error")
		return nil, err
	}
	if len(matchData) == 0 {
		return nil, nil
	}

	// Get match by bookmaker
	var matchLeaguesByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchLeaguesByBookmakers = append(matchLeaguesByBookmakers, val.LeagueName)
		}
	}

	var matchTeamsByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.HomeName)
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.AwayName)
		}
	}

	leagues, err := o.txStorage.Storage().GetUnMachedLeaguesByLeagues(ctx, sportName, firstBookmakerName, secondBookmakerName, matchLeaguesByBookmakers, matchTeamsByBookmakers)
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get leagues error")
		return nil, err
	}

	return leagues, nil
}

func (o *OnlineMatcherService) GetOnlineUnmatchTeams(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) ([]entity.UnMatchedTeamsPairResponse, error) {
	// Get match data from analyzer
	matchData, err := o.analyzerAPI.GetOnlineMatchData()
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get match data error")
		return nil, err
	}
	if len(matchData) == 0 {
		return nil, nil
	}

	// Get match by bookmakentity.Mars
	// Get match by bookmaker
	var matchLeaguesByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchLeaguesByBookmakers = append(matchLeaguesByBookmakers, val.LeagueName)
		}
	}

	var matchTeamsByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.HomeName)
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.AwayName)
		}
	}

	unMatchedTeams, err := o.txStorage.Storage().GetUnMatchedTeamsByLeaguesByTeams(ctx, sportName, firstBookmakerName, secondBookmakerName, matchLeaguesByBookmakers, matchTeamsByBookmakers)
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get leagues error")
		return nil, err
	}

	pairs := make(map[string]entity.UnMatchedTeamsPair)
	for _, valTeam1 := range unMatchedTeams {
		for _, valTeam2 := range unMatchedTeams {
			if valTeam1.BookmakerName != valTeam2.BookmakerName && valTeam1.BookmakerName == firstBookmakerName &&
				valTeam1.LeagueMatchID == valTeam2.LeagueMatchID {

				// Create key for kill duplicates
				key := fmt.Sprintf("%s%s%s%s", valTeam1.BookmakerName, valTeam2.BookmakerName, valTeam1.LeagueName, valTeam2.LeagueName)

				teamsFirst := make(map[int64]entity.UnMatchedTeam)
				teamsSecond := make(map[int64]entity.UnMatchedTeam)

				_, ok := pairs[key]
				if ok {
					teamsFirst = pairs[key].TeamsFirst
					teamsSecond = pairs[key].TeamsSecond
				}

				teamsFirst[valTeam1.TeamID] = entity.UnMatchedTeam{TeamID: valTeam1.TeamID, TeamName: valTeam1.TeamName}
				teamsSecond[valTeam2.TeamID] = entity.UnMatchedTeam{TeamID: valTeam2.TeamID, TeamName: valTeam2.TeamName}

				pairs[key] = entity.UnMatchedTeamsPair{
					LeagueIDFirst:       valTeam1.LeagueID,
					LeagueIDSecond:      valTeam2.LeagueID,
					BookmakerNameFirst:  valTeam1.BookmakerName,
					BookmakerNameSecond: valTeam2.BookmakerName,
					LeagueNameFirst:     valTeam1.LeagueName,
					LeagueNameSecond:    valTeam2.LeagueName,
					TeamsFirst:          teamsFirst,
					TeamsSecond:         teamsSecond,
					SportName:           sportName,
				}
			}
		}
	}

	var result []entity.UnMatchedTeamsPairResponse
	for _, val := range pairs {
		var TeamsFirst []entity.UnMatchedTeam
		for _, firstT := range val.TeamsFirst {
			TeamsFirst = append(TeamsFirst, firstT)
		}

		var TeamsSecond []entity.UnMatchedTeam
		for _, secondT := range val.TeamsSecond {
			TeamsSecond = append(TeamsSecond, secondT)
		}

		result = append(result, entity.UnMatchedTeamsPairResponse{
			LeagueIDFirst:       val.LeagueIDFirst,
			LeagueIDSecond:      val.LeagueIDSecond,
			BookmakerNameFirst:  val.BookmakerNameFirst,
			BookmakerNameSecond: val.BookmakerNameSecond,
			LeagueNameFirst:     val.LeagueNameFirst,
			LeagueNameSecond:    val.LeagueNameSecond,
			TeamsFirst:          TeamsFirst,
			TeamsSecond:         TeamsSecond,
			SportName:           val.SportName,
		})
	}

	return result, nil
}

func (o *OnlineMatcherService) GetOnlineUnmatchLeaguesPrematch(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) ([]entity.League, error) {
	// Get match data from analyzer
	matchData, err := o.analyzerPrematchAPI.GetOnlineMatchData()
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get match data error")
		return nil, err
	}
	if len(matchData) == 0 {
		return nil, nil
	}

	// Get match by bookmaker
	var matchLeaguesByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchLeaguesByBookmakers = append(matchLeaguesByBookmakers, val.LeagueName)
		}
	}

	var matchTeamsByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.HomeName)
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.AwayName)
		}
	}

	leagues, err := o.txStorage.Storage().GetUnMachedLeaguesByLeagues(ctx, sportName, firstBookmakerName, secondBookmakerName, matchLeaguesByBookmakers, matchTeamsByBookmakers)
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get leagues error")
		return nil, err
	}

	return leagues, nil
}

func (o *OnlineMatcherService) GetOnlineUnmatchTeamsPrematch(ctx context.Context, sportName, firstBookmakerName, secondBookmakerName string) ([]entity.UnMatchedTeamsPairResponse, error) {
	// Get match data from analyzer
	matchData, err := o.analyzerPrematchAPI.GetOnlineMatchData()
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get match data error")
		return nil, err
	}
	if len(matchData) == 0 {
		return nil, nil
	}

	// Get match by bookmakentity.Mars
	// Get match by bookmaker
	var matchLeaguesByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchLeaguesByBookmakers = append(matchLeaguesByBookmakers, val.LeagueName)
		}
	}

	var matchTeamsByBookmakers []string
	for _, val := range matchData {
		if (val.Bookmaker == firstBookmakerName || val.Bookmaker == secondBookmakerName) && val.SportName == sportName {
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.HomeName)
			matchTeamsByBookmakers = append(matchTeamsByBookmakers, val.AwayName)
		}
	}

	unMatchedTeams, err := o.txStorage.Storage().GetUnMatchedTeamsByLeaguesByTeams(ctx, sportName, firstBookmakerName, secondBookmakerName, matchLeaguesByBookmakers, matchTeamsByBookmakers)
	if err != nil {
		o.logger.Error().Err(err).Msg("[OnlineMatcherService.GetOnlineUnmatchLeagues] get leagues error")
		return nil, err
	}

	pairs := make(map[string]entity.UnMatchedTeamsPair)
	for _, valTeam1 := range unMatchedTeams {
		for _, valTeam2 := range unMatchedTeams {
			if valTeam1.BookmakerName != valTeam2.BookmakerName && valTeam1.BookmakerName == firstBookmakerName &&
				valTeam1.LeagueMatchID == valTeam2.LeagueMatchID {

				// Create key for kill duplicates
				key := fmt.Sprintf("%s%s%s%s", valTeam1.BookmakerName, valTeam2.BookmakerName, valTeam1.LeagueName, valTeam2.LeagueName)

				teamsFirst := make(map[int64]entity.UnMatchedTeam)
				teamsSecond := make(map[int64]entity.UnMatchedTeam)

				_, ok := pairs[key]
				if ok {
					teamsFirst = pairs[key].TeamsFirst
					teamsSecond = pairs[key].TeamsSecond
				}

				teamsFirst[valTeam1.TeamID] = entity.UnMatchedTeam{TeamID: valTeam1.TeamID, TeamName: valTeam1.TeamName}
				teamsSecond[valTeam2.TeamID] = entity.UnMatchedTeam{TeamID: valTeam2.TeamID, TeamName: valTeam2.TeamName}

				pairs[key] = entity.UnMatchedTeamsPair{
					LeagueIDFirst:       valTeam1.LeagueID,
					LeagueIDSecond:      valTeam2.LeagueID,
					BookmakerNameFirst:  valTeam1.BookmakerName,
					BookmakerNameSecond: valTeam2.BookmakerName,
					LeagueNameFirst:     valTeam1.LeagueName,
					LeagueNameSecond:    valTeam2.LeagueName,
					TeamsFirst:          teamsFirst,
					TeamsSecond:         teamsSecond,
					SportName:           sportName,
				}
			}
		}
	}

	var result []entity.UnMatchedTeamsPairResponse
	for _, val := range pairs {
		var TeamsFirst []entity.UnMatchedTeam
		for _, firstT := range val.TeamsFirst {
			TeamsFirst = append(TeamsFirst, firstT)
		}

		var TeamsSecond []entity.UnMatchedTeam
		for _, secondT := range val.TeamsSecond {
			TeamsSecond = append(TeamsSecond, secondT)
		}

		result = append(result, entity.UnMatchedTeamsPairResponse{
			LeagueIDFirst:       val.LeagueIDFirst,
			LeagueIDSecond:      val.LeagueIDSecond,
			BookmakerNameFirst:  val.BookmakerNameFirst,
			BookmakerNameSecond: val.BookmakerNameSecond,
			LeagueNameFirst:     val.LeagueNameFirst,
			LeagueNameSecond:    val.LeagueNameSecond,
			TeamsFirst:          TeamsFirst,
			TeamsSecond:         TeamsSecond,
			SportName:           val.SportName,
		})
	}

	return result, nil
}
