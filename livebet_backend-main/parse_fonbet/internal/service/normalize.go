package service

import "strings"

func normalizeFootballLeague(league string) string {
	return normalizeAllName(league)
}

func normalizeTennisLeague(league string) string {
	return normalizeAllName(league)
}

func normalizeFootballTeam(team string) string {
	return normalizeAllName(team)
}

func normalizeTennisTeam(team string) string {
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
