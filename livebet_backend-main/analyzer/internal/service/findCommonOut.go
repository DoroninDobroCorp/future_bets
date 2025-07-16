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
func findCommonOutcomes(sansabetData, pinnacleData []entity.PeriodData, homeScore, awayScore int) map[string]entity.OddsWithMarket {
	if len(sansabetData) == 0 || len(pinnacleData) == 0 {
		return nil
	}

	common := make(map[string]entity.OddsWithMarket)

	if pinnacleData[0].Win1x2.Win1.Value >= MIN_VALUE && pinnacleData[0].Win1x2.Win1.Value <= MAX_VALUE {
		common["1"] = entity.OddsWithMarket{MarketType: checkMarketWinHome(homeScore, awayScore), Odds: [2]entity.Odd{sansabetData[0].Win1x2.Win1, pinnacleData[0].Win1x2.Win1}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: home")
	}
	if pinnacleData[0].Win1x2.WinNone.Value >= MIN_VALUE && pinnacleData[0].Win1x2.WinNone.Value <= MAX_VALUE {
		common["X"] = entity.OddsWithMarket{MarketType: checkMarketWinNone(homeScore, awayScore), Odds: [2]entity.Odd{sansabetData[0].Win1x2.WinNone, pinnacleData[0].Win1x2.WinNone}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: draw")
	}
	if pinnacleData[0].Win1x2.Win2.Value >= MIN_VALUE && pinnacleData[0].Win1x2.Win2.Value <= MAX_VALUE {
		common["2"] = entity.OddsWithMarket{MarketType: checkMarketWinAway(homeScore, awayScore), Odds: [2]entity.Odd{sansabetData[0].Win1x2.Win2, pinnacleData[0].Win1x2.Win2}}
		// log.Printf("[DEBUG] Найден общий исход Win1x2: away")
	}

	// Создаем нормализованные мапы для тоталов
	normalizedSansaTotals := make(map[string]entity.WinLessMore)
	normalizedPinnTotals := make(map[string]entity.WinLessMore)

	// Нормализуем ключи для основных тоталов
	for key, value := range sansabetData[0].Totals {
		normalizedKey := normalizeTotal(key)
		normalizedSansaTotals[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].Totals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnTotals[normalizedKey] = *value
	}

	// Проверяем тоталы с нормализованными ключами
	for key, sansabetTotal := range normalizedSansaTotals {
		if pinnacleTotal, exists := normalizedPinnTotals[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["T> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["T< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для индивидуальных тоталов первой команды
	normalizedSansaFirstTeam := make(map[string]entity.WinLessMore)
	normalizedPinnFirstTeam := make(map[string]entity.WinLessMore)

	for key, value := range sansabetData[0].FirstTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedSansaFirstTeam[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].FirstTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnFirstTeam[normalizedKey] = *value
	}

	// Проверяем индивидуальные тоталы первой команды с нормализованными ключами
	for key, sansabetTotal := range normalizedSansaFirstTeam {
		if pinnacleTotal, exists := normalizedPinnFirstTeam[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["IT1> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход First Team Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["IT1< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход First Team Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для индивидуальных тоталов второй команды
	normalizedSansaSecondTeam := make(map[string]entity.WinLessMore)
	normalizedPinnSecondTeam := make(map[string]entity.WinLessMore)

	for key, value := range sansabetData[0].SecondTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedSansaSecondTeam[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].SecondTeamTotals {
		normalizedKey := normalizeTotal(key)
		normalizedPinnSecondTeam[normalizedKey] = *value
	}

	// Проверяем индивидуальные тоталы второй команды с нормализованными ключами
	for key, sansabetTotal := range normalizedSansaSecondTeam {
		if pinnacleTotal, exists := normalizedPinnSecondTeam[key]; exists {
			// Проверяем WinMore
			if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
				common["IT2> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
				// log.Printf("[DEBUG] Найден общий исход Second Team Totals More: %s", key)
			}
			// Проверяем WinLess
			if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
				common["IT2< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
				// log.Printf("[DEBUG] Найден общий исход Second Team Totals Less: %s", key)
			}
		}
	}

	// Нормализуем ключи для гандикапов
	normalizedSansaHandicap := make(map[string]entity.WinHandicap)
	normalizedPinnHandicap := make(map[string]entity.WinHandicap)

	for key, value := range sansabetData[0].Handicap {
		normalizedKey := normalizeTotal(key)
		normalizedSansaHandicap[normalizedKey] = *value
	}
	for key, value := range pinnacleData[0].Handicap {
		normalizedKey := normalizeTotal(key)
		normalizedPinnHandicap[normalizedKey] = *value
	}

	// Проверяем гандикапы с нормализованными ключами
	for key, sansabetHandicap := range normalizedSansaHandicap {
		if pinnacleHandicap, exists := normalizedPinnHandicap[key]; exists {
			// Проверяем Win1
			if pinnacleHandicap.Win1.Value >= MIN_VALUE && pinnacleHandicap.Win1.Value <= MAX_VALUE {
				common["H1 "+key] = entity.OddsWithMarket{MarketType: checkMarketHomeHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{sansabetHandicap.Win1, pinnacleHandicap.Win1}}
				// log.Printf("[DEBUG] Найден общий исход Handicap Win1: %s", key)
			}
			// Проверяем Win2
			if pinnacleHandicap.Win2.Value >= MIN_VALUE && pinnacleHandicap.Win2.Value <= MAX_VALUE {
				common["H2 "+key] = entity.OddsWithMarket{MarketType: checkMarketAwayHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{sansabetHandicap.Win2, pinnacleHandicap.Win2}}
				// log.Printf("[DEBUG] Найден общий исход Handicap Win2: %s", key)
			}
		}
	}

	var maxIndex int 
	if len(pinnacleData) >= len(sansabetData) {
		maxIndex = len(sansabetData)
	} else {
		maxIndex = len(pinnacleData)
	}

	for i := 1; i < maxIndex; i++ {

		// Нормализуем геймы 
		normalizedSansaGame := make(map[string]entity.Win1x2Struct)
		normalizedPinnGame := make(map[string]entity.Win1x2Struct)
		for key, value := range sansabetData[i].Games {
			normalizedSansaGame[key] = *value
		}
		for key, value := range pinnacleData[i].Games {
			normalizedPinnGame[key] = *value
		}
		for key, sansabetGame := range normalizedSansaGame {
			if pinnacleGame, exists := normalizedPinnGame[key]; exists {
				if pinnacleGame.Win1.Value >= MIN_VALUE && pinnacleGame.Win1.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" 1G "+key] = entity.OddsWithMarket{MarketType: 0, Odds: [2]entity.Odd{sansabetGame.Win1, pinnacleGame.Win1}}
					// log.Printf("[DEBUG] Найден общий исход Win1x2: home")
				}
				if pinnacleGame.Win2.Value >= MIN_VALUE && pinnacleGame.Win2.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" 2G "+key] = entity.OddsWithMarket{MarketType: 0, Odds: [2]entity.Odd{sansabetGame.Win2, pinnacleGame.Win2}}
					// log.Printf("[DEBUG] Найден общий исход Win1x2: away")
				}
			}
		}

		// Нормализуем тоталы первого тайма
		normalizedSansaTime1Totals := make(map[string]entity.WinLessMore)
		normalizedPinnTime1Totals := make(map[string]entity.WinLessMore)
		for key, value := range sansabetData[i].Totals {
			normalizedKey := normalizeTotal(key)
			normalizedSansaTime1Totals[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].Totals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnTime1Totals[normalizedKey] = *value
		}

		for key, sansabetTotal := range normalizedSansaTime1Totals {
			if pinnacleTotal, exists := normalizedPinnTime1Totals[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" T> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" T< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Totals Less: %s", key)
				}
			}
		}

		// Нормализуем индивидуальные тоталы первого тайма первой команды
		normalizedSansaTime1FirstTeam := make(map[string]entity.WinLessMore)
		normalizedPinnTime1FirstTeam := make(map[string]entity.WinLessMore)

		for key, value := range sansabetData[i].FirstTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedSansaTime1FirstTeam[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].FirstTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnTime1FirstTeam[normalizedKey] = *value
		}

		for key, sansabetTotal := range normalizedSansaTime1FirstTeam {
			if pinnacleTotal, exists := normalizedPinnTime1FirstTeam[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT1> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" First Team Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT1< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" First Team Totals Less: %s", key)
				}
			}
		}

		// Нормализуем индивидуальные тоталы первого тайма второй команды
		normalizedSansaTime1SecondTeam := make(map[string]entity.WinLessMore)
		normalizedPinnacleTime1SecondTeam := make(map[string]entity.WinLessMore)

		for key, value := range sansabetData[i].SecondTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedSansaTime1SecondTeam[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].SecondTeamTotals {
			normalizedKey := normalizeTotal(key)
			normalizedPinnacleTime1SecondTeam[normalizedKey] = *value
		}

		for key, sansabetTotal := range normalizedSansaTime1SecondTeam {
			if pinnacleTotal, exists := normalizedPinnacleTime1SecondTeam[key]; exists {
				// Проверяем WinMore
				if pinnacleTotal.WinMore.Value >= MIN_VALUE && pinnacleTotal.WinMore.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT2> "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, ">"), Odds: [2]entity.Odd{sansabetTotal.WinMore, pinnacleTotal.WinMore}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Second Team Totals More: %s", key)
				}
				// Проверяем WinLess
				if pinnacleTotal.WinLess.Value >= MIN_VALUE && pinnacleTotal.WinLess.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" IT2< "+key] = entity.OddsWithMarket{MarketType: chechMarketTotal(homeScore, awayScore, "<"), Odds: [2]entity.Odd{sansabetTotal.WinLess, pinnacleTotal.WinLess}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Second Team Totals Less: %s", key)
				}
			}
		}

		// Нормализуем гандикапы первого тайма
		normalizedSansaTime1Handicap := make(map[string]entity.WinHandicap)
		normalizedPinnacleTime1Handicap := make(map[string]entity.WinHandicap)

		for key, value := range sansabetData[i].Handicap {
			normalizedKey := normalizeTotal(key)
			normalizedSansaTime1Handicap[normalizedKey] = *value
		}
		for key, value := range pinnacleData[i].Handicap {
			normalizedKey := normalizeTotal(key)
			normalizedPinnacleTime1Handicap[normalizedKey] = *value
		}

		for key, sansabetHandicap := range normalizedSansaTime1Handicap {
			if pinnacleHandicap, exists := normalizedPinnacleTime1Handicap[key]; exists {
				// Проверяем Win1
				if pinnacleHandicap.Win1.Value >= MIN_VALUE && pinnacleHandicap.Win1.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" H1 "+key] = entity.OddsWithMarket{MarketType: checkMarketHomeHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{sansabetHandicap.Win1, pinnacleHandicap.Win1}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Handicap Win1: %s", key)
				}
				// Проверяем Win2
				if pinnacleHandicap.Win2.Value >= MIN_VALUE && pinnacleHandicap.Win2.Value <= MAX_VALUE {
					common["P"+strconv.Itoa(i)+" H2 "+key] = entity.OddsWithMarket{MarketType: checkMarketAwayHandicap(homeScore, awayScore, key), Odds: [2]entity.Odd{sansabetHandicap.Win2, pinnacleHandicap.Win2}}
					// log.Printf("[DEBUG] Найден общий исход Time"+string(i)+" Handicap Win2: %s", key)
				}
			}
		}
	}

	return common
}
