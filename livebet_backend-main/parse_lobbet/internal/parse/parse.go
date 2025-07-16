package parse

import (
	"livebets/parse_lobbet/internal/entity"
	"livebets/shared"
	"strconv"
)

func LiveToResponseGame(match entity.Match) *shared.GameData {
	resGame := &shared.GameData{
		Pid:     match.ID,
		MatchId: strconv.FormatInt(match.ID, 10),
		IsLive:  true,
	}

	switch match.SportLetter {
	case "S": // Football
		resGame.SportName = shared.SOCCER
		resGame.LeagueName = normalizeFootballLeague(match.LeagueName)
		resGame.HomeName = normalizeFootballTeam(match.HomeTeam)
		resGame.AwayName = normalizeFootballTeam(match.AwayTeam)
	case "T": // Tennis
		resGame.SportName = shared.TENNIS
		resGame.LeagueName = normalizeTennisLeague(match.LeagueName)
		resGame.HomeName = normalizeTennisTeam(match.HomeTeam)
		resGame.AwayName = normalizeTennisTeam(match.AwayTeam)
	case "B": // Basketball
		resGame.SportName = shared.BASKETBALL
		resGame.LeagueName = normalizeBasketballLeague(match.LeagueName)
		resGame.HomeName = normalizeBasketballsTeam(match.HomeTeam)
		resGame.AwayName = normalizeBasketballsTeam(match.AwayTeam)
	}

	resGame.HomeScore = match.MatchResult.CurrentScore.Home
	resGame.AwayScore = match.MatchResult.CurrentScore.Away

	periods := make([]shared.PeriodData, 0, 3)
	// periods[0] - общий результат матча
	// periods[1] - первый период матча
	// periods[2] - второй период матча
	periods = matchHas(2, periods)

	for _, bet := range match.Bets {

		switch bet.LiveBetCaption {
		case "KONAČAN ISHOD": // 1X2
			fillWin1x2(&periods[0].Win1x2, bet.Picks, "ki")

		case "I POLUVRIJEME": // 1X2
			if match.SportLetter == "B" { // Basketball first half
				periods = matchHas(5, periods)
				fillBasketballHalfWin1x2(&periods[5].Win1x2, bet.Picks)
			} else { // Football 1 time
				fillWin1x2(&periods[1].Win1x2, bet.Picks, "Ip")
			}

		case "II POLUVRIJEME": // 1X2 2 time
			if match.SportLetter == "B" { // Basketball second half
				periods = matchHas(6, periods)
				fillBasketballHalfWin1x2(&periods[6].Win1x2, bet.Picks)
			} else { // Football 2 time
				fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIp")
			}

		case "UKUPNO GOLOVA UŽIVO": // Totals
			fillTotals(periods[0].Totals, bet.Picks, "gol. uk")

		case "TIM1 UKUPNO GOLOVA": // Individual totals 1 team
			fillTotals(periods[0].FirstTeamTotals, bet.Picks, "tim1gol uk")

		case "TIM2 UKUPNO GOLOVA": // Indivudual totals 2 team
			fillTotals(periods[0].SecondTeamTotals, bet.Picks, "tim2gol uk")

		case "TIM1 GOLOVA I POLUVRIJEME": // Individual totals 1 team 1 time
			fillTotals(periods[1].FirstTeamTotals, bet.Picks, "tim1gol Ip")

		case "TIM2 GOLOVA I POLUVRIJEME": // Individual totals 2 team 1 time
			fillTotals(periods[1].SecondTeamTotals, bet.Picks, "tim2gol Ip")

		case "UKUPNO GOLOVA I POLUVRIJEME UŽIVO": // Totals 1 time
			fillTotals(periods[1].Totals, bet.Picks, "gol. Ip")

		case "UKUPNO GOLOVA II POLUVRIJEME UŽIVO": // Totals 2 time
			fillTotals(periods[2].Totals, bet.Picks, "gol. IIp")

		case "HENDIKEP": // Handicap
			fillHandicap(periods[0].Handicap, bet.Picks, "h")

		case "DUPLA ŠANSA": // DOUBLE CHANCE
			fillHandicapDoubleChance(periods[0].Handicap, bet.Picks, "ds")

		case "OSTATAK MEČA": // REST OF THE MATCH
			fillHandicapScore(periods[0].Handicap, bet.Picks, "ft ost.", resGame.HomeScore, resGame.AwayScore)

		case "I POLUVRIJEME DUPLA ŠANSA": // DOUBLE CHANGE 1 time
			fillHandicapDoubleChance(periods[1].Handicap, bet.Picks, "Ip ds")

		case "OSTATAK I POLUVREMENA": // REST OF THE MATCH 1 time
			fillHandicapScore(periods[1].Handicap, bet.Picks, "Ip ost.", resGame.HomeScore, resGame.AwayScore)

		// ------------------- TENNIS ----------------------
		case "I SET": // 1X2 SET1 - теннис
			fillWin1x2(&periods[1].Win1x2, bet.Picks, "Is")

		case "II SET": // 1X2 SET2 - теннис
			fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIs")

		case "III SET": // 1X2 SET3 - теннис
			periods = matchHas(3, periods)
			fillWin1x2(&periods[3].Win1x2, bet.Picks, "IIIs")

		case "IV SET": // 1X2 SET4 - теннис
			periods = matchHas(4, periods)
			fillWin1x2(&periods[4].Win1x2, bet.Picks, "IVs")

		case "V SET": // 1X2 SET5 - теннис
			periods = matchHas(5, periods)
			fillWin1x2(&periods[5].Win1x2, bet.Picks, "Vs")

		case "OSVAJA GEM": // 1X2 Game - теннис
			fillTennisGameWin1x2(&periods, bet.Picks)

		case "UKUPNO GEMOVA": // Totals - теннис
			fillTotals(periods[0].Totals, bet.Picks, "ug")

		case "TIM1 UKUPNO GEMOVA": // Individual totals 1 team - теннис
			fillTotals(periods[0].FirstTeamTotals, bet.Picks, "tim1 ug")

		case "TIM2 UKUPNO GEMOVA": // Individual totals 2 team - теннис
			fillTotals(periods[0].SecondTeamTotals, bet.Picks, "tim2 ug")

		case "I SET UKUPNO GEMOVA": // Totals SET1 - теннис
			fillTotals(periods[1].Totals, bet.Picks, "Isg")

		case "II SET UKUPNO GEMOVA": // Totals SET2 - теннис
			fillTotals(periods[2].Totals, bet.Picks, "IIsg")

		case "III SET UKUPNO GEMOVA": // Totals SET3 - теннис
			periods = matchHas(3, periods)
			fillTotals(periods[3].Totals, bet.Picks, "IIIsg")

		case "IV SET UKUPNO GEMOVA": // Totals SET4 - теннис
			periods = matchHas(4, periods)
			fillTotals(periods[4].Totals, bet.Picks, "IVsg")

		case "V SET UKUPNO GEMOVA": // Totals SET5 - теннис
			periods = matchHas(5, periods)
			fillTotals(periods[5].Totals, bet.Picks, "Vsg")

		case "HENDIKEP U GEMOVIMA": // Handicap - теннис
			fillTennisBasketballHandicap(periods[0].Handicap, bet.Picks, "hg")

		// ------------------- BASKETBALL ----------------------
		case "KONAČAN ISHOD - SA PRODUŽECIMA": // 1X2 - баскетбол
			fillWin1x2(&periods[0].Win1x2, bet.Picks, "kisp")

		//case "KONAČAN ISHOD - BEZ PRODUŽETAKA": // 1X2 - баскетбол
		//	fillWin1x2(&periods[0].Win1x2, bet.Picks, "ki")
		//
		case "I ČETVRTINA": // 1X2 1 четверть - баскетбол
			fillWin1x2(&periods[1].Win1x2, bet.Picks, "Ic")

		case "II ČETVRTINA": // 1X2 2 четверть - баскетбол
			fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIc")

		case "III ČETVRTINA": // 1X2 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillWin1x2(&periods[3].Win1x2, bet.Picks, "IIIc")

		case "IV ČETVRTINA": // 1X2 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillWin1x2(&periods[4].Win1x2, bet.Picks, "IVc")

		case "UKUPNO POENA - SA PRODUŽECIMA": // Totals - баскетбол
			fillTotals(periods[0].Totals, bet.Picks, "upsp")

		case "TIM1 UKUPNO POENA SA PRODUŽECIMA": // Individual totals 1 team - баскетбол
			fillTotals(periods[0].FirstTeamTotals, bet.Picks, "tim1p uk")

		case "TIM2 UKUPNO POENA SA PRODUŽECIMA": // Individual totals 2 team - баскетбол
			fillTotals(periods[0].SecondTeamTotals, bet.Picks, "tim2p uk")

		case "I ČETVRTINA UKUPNO POENA": // Totals 1 четверть - баскетбол
			fillTotals(periods[1].Totals, bet.Picks, "up Ic")

		case "II ČETVRTINA UKUPNO POENA": // Totals 2 четверть - баскетбол
			fillTotals(periods[2].Totals, bet.Picks, "up IIc")

		case "III ČETVRTINA UKUPNO POENA": // Totals 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillTotals(periods[3].Totals, bet.Picks, "up IIIc")

		case "IV ČETVRTINA UKUPNO POENA": // Totals 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillTotals(periods[4].Totals, bet.Picks, "up IVc")

		case "I POLUVRIJEME UKUPNO POENA": // Totals 1 половина - баскетбол
			periods = matchHas(5, periods)
			fillTotals(periods[5].Totals, bet.Picks, "up Ip")

		case "HENDIKEP - SA PRODUŽECIMA": // Handicap - баскетбол
			fillTennisBasketballHandicap(periods[0].Handicap, bet.Picks, "hksp")

		case "I POLUVRIJEME HENDIKEP": // Handicap - баскетбол первая половина матча
			periods = matchHas(5, periods)
			fillTennisBasketballHandicap(periods[5].Handicap, bet.Picks, "h Ip")

		case "II POLUVRIJEME HENDIKEP": // Handicap - баскетбол вторая половина матча
			periods = matchHas(6, periods)
			fillTennisBasketballHandicap(periods[6].Handicap, bet.Picks, "h IIp")

		case "I ČETVRTINA HENDIKEP": // Handicap 1 четверть - баскетбол
			fillTennisBasketballHandicap(periods[1].Handicap, bet.Picks, "h Ic")

		case "II ČETVRTINA HENDIKEP": // Handicap 2 четверть - баскетбол
			fillTennisBasketballHandicap(periods[2].Handicap, bet.Picks, "h IIc")

		case "III ČETVRTINA HENDIKEP": // Handicap 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillTennisBasketballHandicap(periods[3].Handicap, bet.Picks, "h IIIc")

		case "IV ČETVRTINA HENDIKEP": // Handicap 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillTennisBasketballHandicap(periods[4].Handicap, bet.Picks, "h IVc")

		}
	}

	resGame.Periods = periods

	return resGame
}

func PrematchToResponseGame(match entity.Match) *shared.GameData {
	resGame := &shared.GameData{
		Pid:     match.ID,
		MatchId: strconv.FormatInt(match.ID, 10),
	}

	switch match.SportLetter {
	case "S": // Football
		resGame.SportName = shared.SOCCER
		resGame.LeagueName = normalizeFootballLeague(match.LeagueName)
		resGame.HomeName = normalizeFootballTeam(match.HomeTeam)
		resGame.AwayName = normalizeFootballTeam(match.AwayTeam)
	case "T": // Tennis
		resGame.SportName = shared.TENNIS
		resGame.LeagueName = normalizeTennisLeague(match.LeagueName)
		resGame.HomeName = normalizeTennisTeam(match.HomeTeam)
		resGame.AwayName = normalizeTennisTeam(match.AwayTeam)
	case "B": // Basketball
		resGame.SportName = shared.BASKETBALL
		resGame.LeagueName = normalizeBasketballLeague(match.LeagueName)
		resGame.HomeName = normalizeBasketballsTeam(match.HomeTeam)
		resGame.AwayName = normalizeBasketballsTeam(match.AwayTeam)
	}

	resGame.HomeScore = match.MatchResult.CurrentScore.Home
	resGame.AwayScore = match.MatchResult.CurrentScore.Away

	periods := make([]shared.PeriodData, 0, 3)
	// periods[0] - общий результат матча
	// periods[1] - первый период матча
	// periods[2] - второй период матча
	periods = matchHas(2, periods)

	for _, bet := range match.Bets {

		switch bet.LiveBetCaption {
		case "KONAČAN ISHOD": // 1X2
			fillWin1x2(&periods[0].Win1x2, bet.Picks, "ki")

		case "I POLUVRIJEME": // 1X2
			if match.SportLetter == "B" { // Basketball first half
				periods = matchHas(5, periods)
				fillWin1x2(&periods[5].Win1x2, bet.Picks, "Ip")
			} else { // Football 1 time
				fillWin1x2(&periods[1].Win1x2, bet.Picks, "Ip")
			}

		case "II POLUVRIJEME": // 1X2
			if match.SportLetter == "B" { // Basketball second half
				periods = matchHas(6, periods)
				fillWin1x2(&periods[6].Win1x2, bet.Picks, "IIp")
			} else { // Football 2 time
				fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIp")
			}

		case "UKUPNO GOLOVA": // Totals
			fillPreMatchFootballTotals(periods[0].Totals, bet.Picks, "ug ")

		case "TIM1 UKUPNO GOLOVA": // Individual totals 1 team
			fillPreMatchFootballTotals(periods[0].FirstTeamTotals, bet.Picks, "ug tim1 ")

		case "TIM2 UKUPNO GOLOVA": // Indivudual totals 2 team
			fillPreMatchFootballTotals(periods[0].SecondTeamTotals, bet.Picks, "ug tim2 ")

		// case "TIM1 GOLOVA I POLUVRIJEME": // Individual totals 1 team 1 time
		// 	fillTotals(periods[1].FirstTeamTotals, bet.Picks, "tim1gol Ip")

		// case "TIM2 GOLOVA I POLUVRIJEME": // Individual totals 2 team 1 time
		// 	fillTotals(periods[1].SecondTeamTotals, bet.Picks, "tim2gol Ip")

		case "UKUPNO GOLOVA I POLUVRIJEME": // Totals 1 time
			fillPreMatchFootballTotals(periods[1].Totals, bet.Picks, "Ip ")

		case "UKUPNO GOLOVA II POLUVRIJEME": // Totals 2 time
			fillPreMatchFootballTotals(periods[2].Totals, bet.Picks, "IIp ")

		case "HENDIKEP A": // Handicap
			fillHandicap(periods[0].Handicap, bet.Picks, "h")

		case "HENDIKEP I POLUVRIJEME": // Handicap 1 time
			fillHandicap(periods[1].Handicap, bet.Picks, "h Ip")

		case "DUPLA ŠANSA": // DOUBLE CHANCE
			fillHandicapDoubleChance(periods[0].Handicap, bet.Picks, "ds")

		case "I POLUVRIJEME DUPLA ŠANSA": // DOUBLE CHANGE
			if match.SportLetter == "B" { // Basketball first half
				periods = matchHas(5, periods)
				fillHandicapDoubleChance(periods[5].Handicap, bet.Picks, "Ip ds")
			} else { // Football  1 time
				fillHandicapDoubleChance(periods[1].Handicap, bet.Picks, "Ip ds")
			}

		case "II POLUVRIJEME DUPLA ŠANSA": // DOUBLE CHANGE
			if match.SportLetter == "B" { // Basketball second half
				periods = matchHas(6, periods)
				fillHandicapDoubleChance(periods[6].Handicap, bet.Picks, "IIp ds")
			} else { // Football 2 time
				fillHandicapDoubleChance(periods[2].Handicap, bet.Picks, "IIp ds")
			}

		case "OSTATAK MEČA": // REST OF THE MATCH
			fillHandicapScore(periods[0].Handicap, bet.Picks, "ft ost.", resGame.HomeScore, resGame.AwayScore)

		case "OSTATAK I POLUVREMENA": // REST OF THE MATCH 1 time
			fillHandicapScore(periods[1].Handicap, bet.Picks, "Ip ost.", resGame.HomeScore, resGame.AwayScore)

		// ------------------- TENNIS ----------------------
		case "I SET": // 1X2 SET1 - теннис
			fillWin1x2(&periods[1].Win1x2, bet.Picks, "Is")

		case "II SET": // 1X2 SET2 - теннис
			fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIs")

		case "III SET": // 1X2 SET3 - теннис
			periods = matchHas(3, periods)
			fillWin1x2(&periods[3].Win1x2, bet.Picks, "IIIs")

		case "IV SET": // 1X2 SET4 - теннис
			periods = matchHas(4, periods)
			fillWin1x2(&periods[4].Win1x2, bet.Picks, "IVs")

		case "V SET": // 1X2 SET5 - теннис
			periods = matchHas(5, periods)
			fillWin1x2(&periods[5].Win1x2, bet.Picks, "Vs")

		case "OSVAJA GEM": // 1X2 Game - теннис
			fillTennisGameWin1x2(&periods, bet.Picks)

		case "TIM1 UKUPNO GEMOVA": // Individual totals 1 team - теннис
			fillPreMatchTennisTotals(periods[0].FirstTeamTotals, bet.Picks, "tim1 ug")

		case "TIM2 UKUPNO GEMOVA": // Individual totals 2 team - теннис
			fillPreMatchTennisTotals(periods[0].SecondTeamTotals, bet.Picks, "tim2 ug")

		case "UKUPNO GEMOVA": // Totals - теннис
			fillPreMatchTennisTotals(periods[0].Totals, bet.Picks, "ug")

		case "UKUPNO GEMOVA ALTERNATIVNI (A)": // Totals - теннис
			fillPreMatchTennisTotals(periods[0].Totals, bet.Picks, "ug")

		case "I SET UKUPNO GEMOVA": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "I SET GEMOVA ALTERNATIVNI (A)": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "I SET GEMOVA ALTERNATIVNI (B)": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "I SET GEMOVA ALTERNATIVNI (C)": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "I SET GEMOVA ALTERNATIVNI (D)": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "I SET GEMOVA ALTERNATIVNI (E)": // Totals SET1 - теннис
			fillPreMatchTennisTotals(periods[1].Totals, bet.Picks, "Isg")

		case "II SET GEMOVA": // Totals SET2 - теннис
			fillPreMatchTennisTotals(periods[2].Totals, bet.Picks, "IIsg")

		case "III SET GEMOVA": // Totals SET3 - теннис
			periods = matchHas(3, periods)
			fillPreMatchTennisTotals(periods[3].Totals, bet.Picks, "IIIsg")

		case "IV SET GEMOVA": // Totals SET4 - теннис
			periods = matchHas(4, periods)
			fillPreMatchTennisTotals(periods[4].Totals, bet.Picks, "IVsg")

		case "V SET GEMOVA": // Totals SET5 - теннис
			periods = matchHas(5, periods)
			fillPreMatchTennisTotals(periods[5].Totals, bet.Picks, "Vsg")

		case "HENDIKEP U GEMOVIMA": // Handicap - теннис
			fillTennisBasketballHandicap(periods[0].Handicap, bet.Picks, "hg")

		// ------------------- BASKETBALL ----------------------
		case "HENDIKEP": // Handicap - баскетбол
			fillTennisBasketballHandicap(periods[0].Handicap, bet.Picks, "hksp")

		case "I POLUVRIJEME HENDIKEP": // Handicap - баскетбол первая половина матча
			periods = matchHas(5, periods)
			fillTennisBasketballHandicap(periods[5].Handicap, bet.Picks, "h Ip")

		case "II POLUVRIJEME HENDIKEP": // Handicap - баскетбол вторая половина матча
			periods = matchHas(6, periods)
			fillTennisBasketballHandicap(periods[6].Handicap, bet.Picks, "h IIp")

		case "UKUPNO POENA": // Totals - баскетбол
			fillBasketballTotals(periods[0].Totals, bet.Picks, "upsp")

		case "TIM1 UKUPNO POENA": // Individual totals 1 team - баскетбол
			fillBasketballTotals(periods[0].FirstTeamTotals, bet.Picks, "tim1p uk")

		case "TIM2 UKUPNO POENA": // Individual totals 2 team - баскетбол
			fillBasketballTotals(periods[0].SecondTeamTotals, bet.Picks, "tim2p uk")

		case "I POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(5, periods)
			fillBasketballTotals(periods[5].Totals, bet.Picks, "up Ip")

		case "II POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(6, periods)
			fillBasketballTotals(periods[6].Totals, bet.Picks, "up IIp")

		case "TIM1 I POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(5, periods)
			fillBasketballTotals(periods[5].FirstTeamTotals, bet.Picks, "tim1p Ip")

		case "TIM2 I POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(5, periods)
			fillBasketballTotals(periods[5].SecondTeamTotals, bet.Picks, "tim2p Ip")

		case "TIM1 II POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(6, periods)
			fillBasketballTotals(periods[6].FirstTeamTotals, bet.Picks, "tim1p IIp")

		case "TIM2 II POLUVRIJEME POENA": // Totals - баскетбол
			periods = matchHas(6, periods)
			fillBasketballTotals(periods[6].SecondTeamTotals, bet.Picks, "tim2p IIp")

		case "TIM1 I ČETVRTINA POENA": // Totals - баскетбол
			fillBasketballTotals(periods[1].FirstTeamTotals, bet.Picks, "tim1p Ic")

		case "TIM2 I ČETVRTINA POENA": // Totals - баскетбол
			fillBasketballTotals(periods[1].SecondTeamTotals, bet.Picks, "tim2p Ic")

		case "TIM1 II ČETVRTINA POENA": // Totals - баскетбол
			fillBasketballTotals(periods[2].FirstTeamTotals, bet.Picks, "tim1p IIc")

		case "TIM2 II ČETVRTINA POENA": // Totals - баскетбол
			fillBasketballTotals(periods[2].SecondTeamTotals, bet.Picks, "tim2p IIc")

		case "TIM1 III ČETVRTINA POENA": // Totals - баскетбол
			periods = matchHas(3, periods)
			fillBasketballTotals(periods[3].FirstTeamTotals, bet.Picks, "tim1p IIIc")

		case "TIM2 III ČETVRTINA POENA": // Totals - баскетбол
			periods = matchHas(3, periods)
			fillBasketballTotals(periods[3].SecondTeamTotals, bet.Picks, "tim2p IIIc")

		case "TIM1 IV ČETVRTINA POENA": // Totals - баскетбол
			periods = matchHas(4, periods)
			fillBasketballTotals(periods[4].FirstTeamTotals, bet.Picks, "tim1p IVc")

		case "TIM2 IV ČETVRTINA POENA": // Totals - баскетбол
			periods = matchHas(4, periods)
			fillBasketballTotals(periods[4].SecondTeamTotals, bet.Picks, "tim2p IVc")

		case "I ČETVRTINA": // 1X2 1 четверть - баскетбол
			fillWin1x2(&periods[1].Win1x2, bet.Picks, "Ic")

		case "II ČETVRTINA": // 1X2 2 четверть - баскетбол
			fillWin1x2(&periods[2].Win1x2, bet.Picks, "IIc")

		case "III ČETVRTINA": // 1X2 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillWin1x2(&periods[3].Win1x2, bet.Picks, "IIIc")

		case "IV ČETVRTINA": // 1X2 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillWin1x2(&periods[4].Win1x2, bet.Picks, "IVc")

		case "I ČETVRTINA POENA": // Totals 1 четверть - баскетбол
			fillBasketballTotals(periods[1].Totals, bet.Picks, "up Ic")

		case "II ČETVRTINA POENA": // Totals 2 четверть - баскетбол
			fillBasketballTotals(periods[2].Totals, bet.Picks, "up IIc")

		case "III ČETVRTINA POENA": // Totals 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillBasketballTotals(periods[3].Totals, bet.Picks, "up IIIc")

		case "IV ČETVRTINA POENA": // Totals 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillBasketballTotals(periods[4].Totals, bet.Picks, "up IVc")

		case "I ČETVRTINA DUPLA ŠANSA": // DOUBLE CHANCE
			fillHandicapDoubleChance(periods[1].Handicap, bet.Picks, "Ic ds")

		case "II ČETVRTINA DUPLA ŠANSA": // DOUBLE CHANCE
			fillHandicapDoubleChance(periods[2].Handicap, bet.Picks, "IIc ds")

		case "III ČETVRTINA DUPLA ŠANSA": // DOUBLE CHANCE
			periods = matchHas(3, periods)
			fillHandicapDoubleChance(periods[3].Handicap, bet.Picks, "IIIc ds")

		case "IV ČETVRTINA DUPLA ŠANSA": // DOUBLE CHANCE
			periods = matchHas(4, periods)
			fillHandicapDoubleChance(periods[4].Handicap, bet.Picks, "IVc ds")

		case "I ČETVRTINA HENDIKEP": // Handicap 1 четверть - баскетбол
			fillTennisBasketballHandicap(periods[1].Handicap, bet.Picks, "h Ic")

		case "II ČETVRTINA HENDIKEP": // Handicap 2 четверть - баскетбол
			fillTennisBasketballHandicap(periods[2].Handicap, bet.Picks, "h IIc")

		case "III ČETVRTINA HENDIKEP": // Handicap 3 четверть - баскетбол
			periods = matchHas(3, periods)
			fillTennisBasketballHandicap(periods[3].Handicap, bet.Picks, "h IIIc")

		case "IV ČETVRTINA HENDIKEP": // Handicap 4 четверть - баскетбол
			periods = matchHas(4, periods)
			fillTennisBasketballHandicap(periods[4].Handicap, bet.Picks, "h IVc")
		}
	}

	resGame.Periods = periods

	return resGame
}
