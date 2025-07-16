package service

import (
	"strconv"
)

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func checkMarketWinHome(homeScore, awayScore int) int {
	if homeScore > awayScore {
		return -1
	}
	return 0
}

func checkMarketWinNone(homeScore, awayScore int) int {
	if homeScore == awayScore {
		return -1
	}
	return 0
}

func checkMarketWinAway(homeScore, awayScore int) int {
	if awayScore > homeScore {
		return -1
	}
	return 0
}

func chechMarketTotal(homeScore, awayScore int, value string) int {
	if value == "<" {
		return -1
	}
	return 0
}

func checkMarketHomeHandicap(homeScore, awayScore int, value string) int {
	handicap, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	if float64(homeScore-awayScore) > handicap*-1 {
		return -1
	}

	return 0
}

func checkMarketAwayHandicap(homeScore, awayScore int, value string) int {
	handicap, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	if float64(awayScore-homeScore) > handicap*-1 {
		return -1
	}

	return 0
}
