package service

import (
	"fmt"
	"livebets/analazer/internal/entity"
	"strconv"

	fuzz "github.com/paul-mannino/go-fuzzywuzzy"
)

// Нормализация ключа тотала
func reverseHandicap(total string) string {
	// Пробуем преобразовать строку в число
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return total // Если не удалось преобразовать, возвращаем как есть
	}
	// Возвращаем строку с одним знаком после запятой
	return fmt.Sprintf("%.1f", value * -1)
}


func reverseCoefs(teams1 string, gamedata entity.GameData) (entity.GameData, error) {
	if fuzz.Ratio(teams1, fmt.Sprintf("%s %s", gamedata.HomeName, gamedata.AwayName)) > 75 {
		return gamedata, nil
	}

	if fuzz.Ratio(teams1, fmt.Sprintf("%s %s", gamedata.AwayName, gamedata.HomeName)) > 75 {
		newGameData := gamedata

		newGameData.HomeName = gamedata.AwayName
		newGameData.AwayName = gamedata.HomeName

		newGameData.AwayScore = gamedata.HomeScore
		newGameData.HomeScore = gamedata.AwayScore

		for i, val := range newGameData.Periods {
			newGameData.Periods[i].Win1x2.Win1 = val.Win1x2.Win2
			newGameData.Periods[i].Win1x2.Win2 = val.Win1x2.Win1

			for key, valueMap := range val.Games {
				newGameData.Periods[i].Games[key] = &entity.Win1x2Struct{Win1: valueMap.Win2, WinNone: valueMap.WinNone, Win2: valueMap.Win1}
			}

			for key, valueMap := range val.Handicap {
				newGameData.Periods[i].Handicap[reverseHandicap(key)] = &entity.WinHandicap{Win1: valueMap.Win2, Win2: valueMap.Win1}
			}

			newGameData.Periods[i].FirstTeamTotals = val.SecondTeamTotals
			newGameData.Periods[i].SecondTeamTotals = val.FirstTeamTotals
		}

		return newGameData, nil
	}

	return gamedata, fmt.Errorf("error reverse coefs")
}
