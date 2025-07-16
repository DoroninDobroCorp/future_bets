package parse

import (
	"livebets/parse_starcasino/internal/entity"
	"livebets/shared"
	"strconv"
	"time"
)

func StarCasinoToResponseGame(match entity.Match) *shared.GameData {
	resGame := &shared.GameData{
		Pid:        match.Id,
		LeagueName: match.Category.Name + " " + match.League.Name,
		MatchId:    strconv.FormatInt(match.Id, 10),
	}

	if len(match.Teams) == 2 {
		resGame.HomeName = match.Teams[0].Name
		resGame.AwayName = match.Teams[1].Name
	}

	if len(match.Scores) == 2 {
		resGame.HomeScore = match.Scores[0]
		resGame.AwayScore = match.Scores[1]
	}

	switch match.Sport.Id {
	case 66: // Football
		resGame.SportName = shared.SOCCER
		resGame.Periods = parseFootball(match)

		resGame.LeagueName = normalizeFootballLeague(resGame.LeagueName)
		resGame.HomeName = normalizeFootballTeam(resGame.HomeName)
		resGame.AwayName = normalizeFootballTeam(resGame.AwayName)

	case 68: // Tennis
		resGame.SportName = shared.TENNIS
		resGame.Periods = parseTennis(match)

		resGame.LeagueName = normalizeTennisLeague(resGame.LeagueName)
		resGame.HomeName = normalizeTennisTeam(resGame.HomeName)
		resGame.AwayName = normalizeTennisTeam(resGame.AwayName)
	}

	// Add config data
	resGame.CreatedAt = time.Now()
	resGame.Source = shared.STARCASINO

	return resGame
}

func parseFootball(match entity.Match) []shared.PeriodData {

	bets := getBets(match.Markets)

	periods := make([]shared.PeriodData, 0, 3)
	// periods[0] - общий результат матча
	// periods[1] - первый период матча
	// periods[2] - второй период матча
	periods = matchHas(2, periods)

	for _, odd := range match.Odds {

		if odd.IsBB {
			continue
		}

		bet, ok := bets[odd.Id]
		if !ok {
			continue
		}

		switch bet.Name {
		case "1x2":
			fillWin1x2(&periods[0].Win1x2, odd)

		case "1st half - 1x2":
			fillWin1x2(&periods[1].Win1x2, odd)

		case "2nd half - 1x2":
			fillWin1x2(&periods[2].Win1x2, odd)

		case "Handicap":
			fillHandicap(periods[0].Handicap, odd, bet.Sv)

		case "1st half - handicap":
			fillHandicap(periods[1].Handicap, odd, bet.Sv)

		case "2nd half - handicap":
			fillHandicap(periods[2].Handicap, odd, bet.Sv)

		case "Double chance":
			fillFootballDoubleChance(periods[0].Handicap, odd)

		case "1st half - double chance":
			fillFootballDoubleChance(periods[1].Handicap, odd)

		case "2nd half - double chance":
			fillFootballDoubleChance(periods[2].Handicap, odd)

		case "Draw no bet":
			fillFootballDrawNoBet(periods[0].Handicap, odd)

		case "1st half - draw no bet":
			fillFootballDrawNoBet(periods[1].Handicap, odd)

		case "2nd half - draw no bet":
			fillFootballDrawNoBet(periods[2].Handicap, odd)

		default:
			fillFootballTotals(&periods, odd, bet.Name, bet.TypeId, bet.Sv)

			fillFootballRestOfTheMatch(&periods, odd, bet.Name)

			// fmt.Println(odd.TypeId)
		}
	}

	return periods
}

func parseTennis(match entity.Match) []shared.PeriodData {

	bets := getBets(match.Markets)

	periods := make([]shared.PeriodData, 0, 4)
	// periods[0] - общий результат матча
	// periods[1] - первый set матча
	// periods[2] - второй set матча
	periods = matchHas(2, periods)

	for _, odd := range match.Odds {

		if odd.IsBB {
			continue
		}

		bet, ok := bets[odd.Id]
		if !ok {
			continue
		}

		fillTennisWin1x2(&periods, odd, bet.Name, bet.Sv)

		fillTennisTotals(&periods, odd, bet.Name, bet.TypeId, bet.Sv)

		fillTennisHandicap(&periods, odd, bet.Name, bet.Sv)

	}

	return periods
}

func getBets(markets []*entity.Market) map[int64]*entity.Market {
	bets := make(map[int64]*entity.Market, len(markets)*4)
	for _, market := range markets {
		for _, oddIds := range market.OddIds {
			for _, oddId := range oddIds {
				bets[oddId] = market
			}
		}
	}
	return bets
}
