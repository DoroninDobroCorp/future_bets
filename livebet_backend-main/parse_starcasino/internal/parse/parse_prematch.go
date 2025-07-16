package parse

import (
	"livebets/parse_starcasino/internal/entity"
	"livebets/shared"
	"strconv"
	"time"
)

func StarCasinoPreMatchToResponseGames(preMatchData *entity.PreMatchData) ([]*shared.GameData, bool) {

	stopPagination := false // Остановить дальнейшую загрузку страниц

	categories := itemSliceToMap(preMatchData.Categories)
	leagues := itemSliceToMap(preMatchData.Leagues)
	teams := itemSliceToMap(preMatchData.Teams)

	marketMap := marketSliceToMap(preMatchData.Markets)
	oddMap := oddSliceToMap(preMatchData.Odds)

	gameDataLice := make([]*shared.GameData, 0, len(preMatchData.Events))

	for _, event := range preMatchData.Events {

		if !dateIn48hoursInterval(event.StartDate) {
			// на этой странице есть пре-матчи, которые начинаются более чем через 48 часов
			stopPagination = true
			break
		}

		gameData := &shared.GameData{
			Pid:        event.Id,
			LeagueName: categories[event.CatId] + " " + leagues[event.LeagueId],
			MatchId:    strconv.FormatInt(event.Id, 10),
		}

		if len(event.TeamIds) == 2 {
			gameData.HomeName = teams[event.TeamIds[0]]
			gameData.AwayName = teams[event.TeamIds[1]]
		}

		switch event.SportId {
		case 66: // Football
			gameData.SportName = shared.SOCCER // Football
			gameData.Periods = parsePreMatchFootball(event, marketMap, oddMap)

			gameData.LeagueName = normalizeFootballLeague(gameData.LeagueName)
			gameData.HomeName = normalizeFootballTeam(gameData.HomeName)
			gameData.AwayName = normalizeFootballTeam(gameData.AwayName)

		case 68: // Tennis
			gameData.SportName = shared.TENNIS
			gameData.Periods = parsePreMatchTennis(event, marketMap, oddMap)

			gameData.LeagueName = normalizeTennisLeague(gameData.LeagueName)
			gameData.HomeName = normalizeTennisTeam(gameData.HomeName)
			gameData.AwayName = normalizeTennisTeam(gameData.AwayName)
		}

		// Add config data
		gameData.CreatedAt = time.Now()
		gameData.Source = shared.STARCASINO
		gameData.IsLive = false

		gameDataLice = append(gameDataLice, gameData)
	}

	return gameDataLice, stopPagination
}

func itemSliceToMap(itemSlice []*entity.Item) map[int64]string {
	itemMap := make(map[int64]string, len(itemSlice))
	for _, item := range itemSlice {
		itemMap[item.Id] = item.Name
	}
	return itemMap
}

func dateIn48hoursInterval(dateStr string) bool {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return false
	}
	return time.Until(date) < time.Hour*48
}

func marketSliceToMap(marketSlice []*entity.PreMatchMarket) map[int64]*entity.PreMatchMarket {
	marketMap := make(map[int64]*entity.PreMatchMarket, len(marketSlice))
	for _, market := range marketSlice {
		marketMap[market.Id] = market
	}
	return marketMap
}

func oddSliceToMap(oddSlice []*entity.Odd) map[int64]*entity.Odd {
	oddMap := make(map[int64]*entity.Odd, len(oddSlice))
	for _, odd := range oddSlice {
		oddMap[odd.Id] = odd
	}
	return oddMap
}

func parsePreMatchFootball(event *entity.Event, marketMap map[int64]*entity.PreMatchMarket, oddMap map[int64]*entity.Odd) []shared.PeriodData {

	periods := make([]shared.PeriodData, 0, 3)
	// periods[0] - общий результат матча
	// periods[1] - первый период матча
	// periods[2] - второй период матча
	periods = matchHas(2, periods)

	for _, marketId := range event.MarketsIds {

		market, ok := marketMap[marketId]
		if !ok {
			continue
		}

		for _, oddId := range market.OddIds {
			odd, ok := oddMap[oddId]
			if !ok {
				continue
			}
			switch market.Name {
			case "1x2":
				fillWin1x2(&periods[0].Win1x2, odd)

			case "1st half - 1x2":
				fillWin1x2(&periods[1].Win1x2, odd)

			case "2nd half - 1x2":
				fillWin1x2(&periods[2].Win1x2, odd)

			case "Handicap":
				fillHandicap(periods[0].Handicap, odd, market.Sv)

			case "1st half - handicap":
				fillHandicap(periods[1].Handicap, odd, market.Sv)

			case "2nd half - handicap":
				fillHandicap(periods[2].Handicap, odd, market.Sv)

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
				fillFootballTotals(&periods, odd, market.Name, market.Id, market.Sv)

				fillFootballRestOfTheMatch(&periods, odd, market.Name)

			}
		}
	}

	return periods
}

func parsePreMatchTennis(event *entity.Event, marketMap map[int64]*entity.PreMatchMarket, oddMap map[int64]*entity.Odd) []shared.PeriodData {

	periods := make([]shared.PeriodData, 0, 4)
	// periods[0] - общий результат матча
	// periods[1] - первый set матча
	// periods[2] - второй set матча
	periods = matchHas(2, periods)
	for _, marketId := range event.MarketsIds {

		market, ok := marketMap[marketId]
		if !ok {
			continue
		}

		for _, oddId := range market.OddIds {
			odd, ok := oddMap[oddId]
			if !ok {
				continue
			}

			fillTennisWin1x2(&periods, odd, market.Name, market.Sv)

			fillTennisTotals(&periods, odd, market.Name, market.Id, market.Sv)

			fillTennisHandicap(&periods, odd, market.Name, market.Sv)

		}
	}
	return periods
}
