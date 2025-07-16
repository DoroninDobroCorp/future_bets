package parse

import (
	"regexp"
	"strings"
)

func NormalizeFootballLeague(league string) string {
	return normalizeAllName(league)
}

func NormalizeTennisLeague(league string) string {
	league = strings.ReplaceAll(league, "Qual.", "")
	// Слова (\b) Doubles, clay, STB, hard, grass, carpet
	// отдельное слово если после букв M или W идёт число
	// отдельное слово из цифр
	// всё что в скобках и сами скобки
	re := regexp.MustCompile(`\b(Doubles|Clay|STB|hard|grass|carpet|[M|W]\d+|\d+)\b|\(.*\)`)
	league = re.ReplaceAllString(league, "")
	return normalizeAllName(league)
}

func NormalizeBasketballLeague(league string) string {
	return normalizeAllName(league)
}

func NormalizeFootballTeam(team string) string {
	re := regexp.MustCompile(`\b(FC|SC|FK|CF|CD|NK|LK|U\d+)\b`) // отдельные слова: FC,SC,FK,CF,CD,NK,LK или U затем цифры
	team = re.ReplaceAllString(team, "")
	return normalizeAllName(team)
}

func NormalizeTennisTeam(team string) string {
	if strings.Contains(team, "/") {
		// парный матч
		teams := strings.Split(team, "/")
		if len(teams) == 2 {
			for i := range teams {
				// нормализуем каждого игрока
				teams[i] = NormalizeTennisTeam(teams[i])
			}
			return strings.Join(teams, "/")
		}
	}

	// Имя игрока может содержать инициалы с точками
	parts := strings.Fields(team)
	if len(parts) > 1 { // "Agustin Gomez F." => "F Agustin Gomez"
		newParts := make([]string, 0, len(parts))
		for _, part := range parts {
			if strings.Contains(part, ".") {
				// Инициал: убираем точку и добавляем в начало
				cleaned := strings.ReplaceAll(part, ".", "")
				newParts = append([]string{cleaned}, newParts...)
			} else {
				// Фамилия: добавляем в конец
				newParts = append(newParts, part)
			}
		}
		team = strings.Join(newParts, " ")
	}

	return normalizeAllName(team)
}

func NormalizeBasketballTeam(team string) string {
	return normalizeAllName(team)
}

func normalizeAllName(name string) string {
	// Удаляем запятые и дефисы
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, "-", "")

	// Удаляем двойные пробелы и пробелы в начале и конце
	name = strings.Join(strings.Fields(name), " ")

	// Переводим строку в нижний регистр
	name = strings.ToLower(name)

	return name
}
