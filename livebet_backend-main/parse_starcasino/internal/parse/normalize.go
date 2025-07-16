package parse

import (
	"regexp"
	"strings"
)

func normalizeFootballLeague(league string) string {
	return normalizeAllName(league)
}

func normalizeFootballTeam(team string) string {
	re := regexp.MustCompile(`\b(FC|SC|FK|CF|CD|NK|LK|U\d+)\b`) // отдельные слова: FC,SC,FK,CF,CD,NK,LK или U затем цифры
	team = re.ReplaceAllString(team, "")
	return normalizeAllName(team)
}

func normalizeTennisLeague(league string) string {
	// Слова (\b) Singles, Double, Doubles, Qualification, clay, stb, hard, grass, carpet
	// отдельное слово если после букв M или W идёт число
	// отдельное слово из чисел
	// всё что в скобках и сами скобки
	re := regexp.MustCompile(`\b(Singles|Double|Doubles|Qualification|Clay|Stb|Hard|Grass|Carpet|[M|W]\d+|\d+)\b|\(.*\)`)
	league = re.ReplaceAllString(league, "")

	// Удаляем повторяющиеся слова
	seen := make(map[string]bool)
	re = regexp.MustCompile(`\b(\w+)\b`)
	league = re.ReplaceAllStringFunc(league, func(match string) string {
		if seen[match] {
			return ""
		}
		seen[match] = true
		return match
	})

	return normalizeAllName(league)
}

func normalizeTennisTeam(team string) string {

	if strings.Contains(team, " / ") {
		// парный матч
		teams := strings.Split(team, " / ")
		if len(teams) == 2 {
			for i := range teams {
				// нормализуем каждого игрока
				teams[i] = normalizeTennisTeam(teams[i])
			}
			return strings.Join(teams, "/")
		}
	}

	// "Barkova, Tatiana" => "Tatiana Barkova"

	split := strings.Split(team, ", ")

	if len(split) == 2 {
		lastName := split[0]
		firstName := split[1] // имя игрока
		team = firstName + " " + lastName
	}

	return normalizeAllName(team)
}

func normalizeBasketballLeague(league string) string {
	return normalizeAllName(league)
}

func normalizeBasketballsTeam(team string) string {
	if strings.HasSuffix(team, ".") {
		team = strings.TrimSuffix(team, ".")
	}
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
