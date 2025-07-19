package service

import (
	"fmt"
	"livebets/analazer/internal/entity"
	"strconv"
)

const (
	MIN_VALUE = 1.2
	MAX_VALUE = 3.6
)

// Нормализация ключа тотала
func normalizeTotal(total string) string {
	// Пробуем преобразовать строку в число
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return total // Если не удалось преобразовать, возвращаем как есть
	}
	// Возвращаем строку с одним знаком после запятой
	return fmt.Sprintf("%.1f", value)
}

// Поиск общих исходов
func findCommonOutcomes(competitorData, pinnacleData []entity.PeriodData, homeScore, awayScore int) map[string]entity.OddsWithMarket {
	if len(competitorData) == 0 || len(pinnacleData) == 0 {
		return nil
	}

	common := make(map[string]entity.OddsWithMarket)

	if pinnacleData[0].Win1x2.Win1.Value >= MIN_VALUE && pinnacleData[0].Win1x2.Win1.Value <= MAX_VALUE {
		common["1"] = entity.OddsWithMarket{MarketType: checkMarketWinHome(homeScore, awayScore), Odds: [2]entity.Odd{competitorData[0].Win1x2.Win1, pinnacleData[0].Win1x2.Win1}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: home")
	}
	if pinnacleData[0].Win1x2.WinNone.Value >= MIN_VALUE && pinnacleData[0].Win1x2.WinNone.Value <= MAX_VALUE {
		common["X"] = entity.OddsWithMarket{MarketType: checkMarketWinNone(homeScore, awayScore), Odds: [2]entity.Odd{competitorData[0].Win1x2.WinNone, pinnacleData[0].Win1x2.WinNone}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: draw")
	}
	if pinnacleData[0].Win1x2.Win2.Value >= MIN_VALUE && pinnacleData[0].Win1x2.Win2.Value <= MAX_VALUE {
		common["2"] = entity.OddsWithMarket{MarketType: checkMarketWinAway(homeScore, awayScore), Odds: [2]entity.Odd{competitorData[0].Win1x2.Win2, pinnacleData[0].Win1x2.Win2}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: away")
	}

	// Создаем нормализованные мапы для тоталов
	normalizedCompetitorTotals := make(map[string]entity.WinLessMore)
	normalizedPinnTotals := make(map[string]entity.WinLessMore)

	// Нормализуем ключи для основных тоталов
	for key, value := range competitorData[0].Totals {
		normalizedKey := normalizeTotal(key)
		normalizedCompetitorTotals[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].Totals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnTotals[normalizedKey] = *value
	}

	// Проверяем тоталы с нормализованными ключами
	for key, competitorTotal := range normalizedCompetitorTotals {
		if pinnacleTotal, exists := normalizedPinnTotals[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["T> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["T< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для индивидуальных тоталов первой команды
	normalizedCompetitorFirstTeam := make(map[string]entity.WinLessMore)
	normalizedPinnFirstTeam := make(map[string]entity.WinLessMore)

	for key, value := range competitorData[0].FirstTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedCompetitorFirstTeam[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].FirstTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnFirstTeam[normalizedKey] = *value
	}

	// Проверяем индивидуальные тоталы первой команды с нормализованными ключами
	for key, competitorTotal := range normalizedCompetitorFirstTeam {
		if pinnacleTotal, exists := normalizedPinnFirstTeam[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["IT1> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход First Team Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["IT1< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход First Team Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для индивидуальных тоталов второй команды
	normalizedCompetitorSecondTeam := make(map[string]entity.WinLessMore)
	normalizedPinnSecondTeam := make(map[string]entity.WinLessMore)

	for key, value := range competitorData[0].SecondTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedCompetitorSecondTeam[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].SecondTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnSecondTeam[normalizedKey] = *value
	}

	// Проверяем индивидуальные тоталы второй команды с нормализованными ключами
	for key, competitorTotal := range normalizedCompetitorSecondTeam {
		if pinnacleTotal, exists := normalizedPinnSecondTeam[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["IT2> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход Second Team Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["IT2< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход Second Team Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для гандикапов
	normalizedCompetitorHandicap := make(map[string]entity.WinHandicap)
	normalizedPinnHandicap := make(map[string]entity.WinHandicap)

	for key, value := range competitorData[0].Handicap {
		normalizedKey := normalizeTotal(key)
		normalizedCompetitorHandicap[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].Handicap {
		normalizedKey := normalizeTotal(key)
		normalizedPinnHandicap[normalizedKey] = *value
	}

	// Проверяем гандикапы с нормализованными ключами
	for key, competitorHandicap := range normalizedCompetitorHandicap {
		if pinnacleHandicap, exists := normalizedPinnHandicap[key]; exists {
			// Проверяем Win1
			if pinnacleHandicap.Win1.Value >= MIN_VALUE && pinnacleHandicap.Win1.Value <= MAX_VALUE {
				common["H1 "+key] = entity.OddsWithMarket{MarketType: checkMarketHomeHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{competitorHandicap.Win1, pinnacleHandicap.Win1}}
				// log.Printf("[DEBUG] Найден общий исход Handicap Win1: %s", key)
			}
			// Проверяем Win2
			if pinnacleHandicap.Win2.Value >= MIN_VALUE && pinnacleHandicap.Win2.Value <= MAX_VALUE {
				common["H2 "+key] = entity.OddsWithMarket{MarketType: checkMarketAwayHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{competitorHandicap.Win2, pinnacleHandicap.Win2}}
				// log.Printf("[DEBUG] Найден общий исход Handicap Win2: %s", key)
			}
		}
	}

	var maxIndex int
	if len(pinnacleData) >= len(competitorData) {
		maxIndex = len(competitorData)
	} else {
		maxIndex = len(pinnacleData)
	}

	for i := 1; i < maxIndex; i++ {

		// Нормализуем геймы
		normalizedCompetitorGame := make(map[string]entity.Win1x2Struct)
		normalizedPinnGame := make(map[string]entity.Win1x2Struct)
		for key, value := range competitorData[i].Games {
			normalizedCompetitorGame[key] = *value
		}
		for key, value := range pinnacleData[i].Games {
			normalizedPinnGame[key] = *value
		}
		for key, competitorGame := range normalizedCompetitorGame {
			if pinnacleGame, exists := normalizedPinnGame[key]; exists {
				if pinnacleGame.Win1.Value >= MIN_VALUE && pinnacleGame.Win1.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" 1G "+key] = entity.OddsWithMarket{MarketType: 0, Odds: [2]entity.Odd{competitorGame.Win1, pinnacleGame.Win1}}
					// log.Printf("[DEBUG] Найден общий исход Win1x2: home")
				}
				if pinnacleGame.Win2.Value >= MIN_VALUE && pinnacleGame.Win2.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" 2G "+key] = entity.OddsWithMarket{MarketType: 0, Odds: [2]entity.Odd{competitorGame.Win2, pinnacleGame.Win2}}
					// log.Printf("[DEBUG] Найден общий исход Win1x2: away")
				}
			}
		}

		// Нормализуем тоталы первого тайма
		normalizedCompetitorTime1Totals := make(map[string]entity.WinLessMore)
		normalizedPinnTime1Totals := make(map[string]entity.WinLessMore)
		for key, value := range competitorData[i].Totals {
			normalizedKey := normalizeTotal(key)
			normalizedCompetitorTime1Totals[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].Totals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnTime1Totals[normalizedKey] = *value
		}

		for key, competitorTotal := range normalizedCompetitorTime1Totals {
			if pinnacleTotal, exists := normalizedPinnTime1Totals[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" T> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" T< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Totals Less: %s", key)
				}
			}
		}

		// Нормализуем индивидуальные тоталы первого тайма первой команды
		normalizedCompetitorTime1FirstTeam := make(map[string]entity.WinLessMore)
		normalizedPinnTime1FirstTeam := make(map[string]entity.WinLessMore)

		for key, value := range competitorData[i].FirstTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedCompetitorTime1FirstTeam[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].FirstTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnTime1FirstTeam[normalizedKey] = *value
		}

		for key, competitorTotal := range normalizedCompetitorTime1FirstTeam {
			if pinnacleTotal, exists := normalizedPinnTime1FirstTeam[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT1> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" First Team Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT1< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" First Team Totals Less: %s", key)
				}
			}
		}

		// Нормализуем индивидуальные тоталы первого тайма второй команды
		normalizedCompetitorTime1SecondTeam := make(map[string]entity.WinLessMore)
		normalizedPinnacleTime1SecondTeam := make(map[string]entity.WinLessMore)

		for key, value := range competitorData[i].SecondTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedCompetitorTime1SecondTeam[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].SecondTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnacleTime1SecondTeam[normalizedKey] = *value
		}

		for key, competitorTotal := range normalizedCompetitorTime1SecondTeam {
			if pinnacleTotal, exists := normalizedPinnacleTime1SecondTeam[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT2> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{competitorTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Second Team Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT2< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{competitorTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Totals Less: %s", key)
				}
			}
		}

		// Нормализуем гандикапы первого тайма
		normalizedCompetitorTime1Handicap := make(map[string]entity.WinHandicap)
		normalizedPinnacleTime1Handicap := make(map[string]entity.WinHandicap)

		for key, value := range competitorData[i].Handicap {
			normalizedKey := normalizeTotal(key)
			normalizedCompetitorTime1Handicap[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].Handicap {
			normalizedKey := normalizeTotal(key)
			normalizedPinnacleTime1Handicap[normalizedKey] = *value
		}

		for key, competitorHandicap := range normalizedCompetitorTime1Handicap {
			if pinnacleHandicap, exists := normalizedPinnacleTime1Handicap[key]; exists {
				// Проверяем Win1
				if pinnacleHandicap.Win1.Value >= MIN_VALUE && pinnacleHandicap.Win1.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" H1 "+key] = entity.OddsWithMarket{MarketType: checkMarketHomeHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{competitorHandicap.Win1, pinnacleHandicap.Win1}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Handicap Win1: %s", key)
				}
				// Проверяем Win2
				if pinnacleHandicap.Win2.Value >= MIN_VALUE && pinnacleHandicap.Win2.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" H2 "+key] = entity.OddsWithMarket{MarketType: checkMarketAwayHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{competitorHandicap.Win2, pinnacleHandicap.Win2}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Handicap Win2: %s", key)
				}
			}
		}
	}

	return common
}
