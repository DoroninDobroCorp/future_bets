package parse

import (
	"livebets/parse_lobbet/internal/parse/suite"
	"livebets/shared"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Баскетбольный матч, в котором идет I четверть
	basketballFile1Name       = "basketball_match1.json"
	basketballMatch1Id  int64 = 22653761

	// Баскетбольный матч, в котором идет II четверть
	basketballFile2Name       = "basketball_match2.json"
	basketballMatch2Id  int64 = 22657449

	// Баскетбольный матч, в котором идет III четверть
	basketballFile3Name       = "basketball_match3.json"
	basketballMatch3Id  int64 = 22567555

	// Баскетбольный матч, в котором идет IV четверть
	basketballFile4Name       = "basketball_match4.json"
	basketballMatch4Id  int64 = 22572215
)

func TestBasketballMatch1Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "usm blida", resGame.HomeName, "HomeName")
	assert.Equal(t, "cs tlemcen", resGame.AwayName, "AwayName")
	assert.Equal(t, "algeria super division", resGame.LeagueName, "LeagueName")

	assert.Equal(t, 8., resGame.HomeScore, "HomeScore")
	assert.Equal(t, 15., resGame.AwayScore, "AwayScore")

	assert.Equal(t, basketballMatch1Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(basketballMatch1Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.BASKETBALL, resGame.SportName, "SportName")
}

func TestBasketballMatch2Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "hd marines", resGame.HomeName, "HomeName")
	assert.Equal(t, "rouiba cb", resGame.AwayName, "AwayName")
	assert.Equal(t, "algeria national a women", resGame.LeagueName, "LeagueName")

	assert.Equal(t, 14., resGame.HomeScore, "HomeScore")
	assert.Equal(t, 11., resGame.AwayScore, "AwayScore")

	assert.Equal(t, basketballMatch2Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(basketballMatch2Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.BASKETBALL, resGame.SportName, "SportName")
}

func TestBasketballMatch3Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "sydney", resGame.HomeName, "HomeName")
	assert.Equal(t, "illawarra hawks", resGame.AwayName, "AwayName")
	assert.Equal(t, "australia nbl", resGame.LeagueName, "LeagueName")

	assert.Equal(t, 74., resGame.HomeScore, "HomeScore")
	assert.Equal(t, 78., resGame.AwayScore, "AwayScore")

	assert.Equal(t, basketballMatch3Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(basketballMatch3Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.BASKETBALL, resGame.SportName, "SportName")
}

func TestBasketballMatch4Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "meralco bolts", resGame.HomeName, "HomeName")
	assert.Equal(t, "converge fiberxers", resGame.AwayName, "AwayName")
	assert.Equal(t, "philippines commissioners cup", resGame.LeagueName, "LeagueName")

	assert.Equal(t, 74., resGame.HomeScore, "HomeScore")
	assert.Equal(t, 90., resGame.AwayScore, "AwayScore")

	assert.Equal(t, basketballMatch4Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(basketballMatch4Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.BASKETBALL, resGame.SportName, "SportName")
}

func TestBasketballMatch1Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.7},
		WinNone: shared.Odd{Value: 0.},
		Win2:    shared.Odd{Value: 1.4},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestBasketballMatch2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.55},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 1.45},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestBasketballMatch3Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.65},
		WinNone: shared.Odd{Value: 0.},
		Win2:    shared.Odd{Value: 1.43},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestBasketballMatch4Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestBasketballMatch1Time1Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 8.5},
		WinNone: shared.Odd{Value: 13},
		Win2:    shared.Odd{Value: 1.1},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestBasketballMatch2Time2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.75},
		WinNone: shared.Odd{Value: 11},
		Win2:    shared.Odd{Value: 1.55},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestBasketballMatch3Time3Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 3.8},
		WinNone: shared.Odd{Value: 10},
		Win2:    shared.Odd{Value: 1.35},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time3Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time3Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time3Win1x2.Win2")
}

func TestBasketballMatch4Time4Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[4]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1.5},
		WinNone: shared.Odd{Value: 12},
		Win2:    shared.Odd{Value: 2.8},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time4Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time4Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time4Win1x2.Win2")
}

func TestBasketballMatch1Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"131.5": {WinLess: shared.Odd{Value: 2.55}, WinMore: shared.Odd{Value: 1.43}},
		"133.5": {WinLess: shared.Odd{Value: 2.25}, WinMore: shared.Odd{Value: 1.55}},
		"135.5": {WinLess: shared.Odd{Value: 1.98}, WinMore: shared.Odd{Value: 1.7}},
		"136.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.83}},
		"137.5": {WinLess: shared.Odd{Value: 1.7}, WinMore: shared.Odd{Value: 1.97}},
		"139.5": {WinLess: shared.Odd{Value: 1.55}, WinMore: shared.Odd{Value: 2.2}},
		"141.5": {WinLess: shared.Odd{Value: 1.45}, WinMore: shared.Odd{Value: 2.5}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestBasketballMatch2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"96.5":  {WinLess: shared.Odd{Value: 2.8}, WinMore: shared.Odd{Value: 1.35}},
		"98.5":  {WinLess: shared.Odd{Value: 2.35}, WinMore: shared.Odd{Value: 1.5}},
		"100.5": {WinLess: shared.Odd{Value: 2.}, WinMore: shared.Odd{Value: 1.68}},
		"101.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.85}},
		"102.5": {WinLess: shared.Odd{Value: 1.68}, WinMore: shared.Odd{Value: 2.}},
		"104.5": {WinLess: shared.Odd{Value: 1.5}, WinMore: shared.Odd{Value: 2.35}},
		"106.5": {WinLess: shared.Odd{Value: 1.35}, WinMore: shared.Odd{Value: 2.75}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestBasketballMatch3Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"217.5": {WinLess: shared.Odd{Value: 2.85}, WinMore: shared.Odd{Value: 1.35}},
		"218.5": {WinLess: shared.Odd{Value: 2.55}, WinMore: shared.Odd{Value: 1.43}},
		"219.5": {WinLess: shared.Odd{Value: 2.3}, WinMore: shared.Odd{Value: 1.5}},
		"220.5": {WinLess: shared.Odd{Value: 2.05}, WinMore: shared.Odd{Value: 1.65}},
		"221.5": {WinLess: shared.Odd{Value: 1.87}, WinMore: shared.Odd{Value: 1.8}},
		"222.5": {WinLess: shared.Odd{Value: 1.75}, WinMore: shared.Odd{Value: 1.95}},
		"223.5": {WinLess: shared.Odd{Value: 1.63}, WinMore: shared.Odd{Value: 2.05}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestBasketballMatch4Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"194.5": {WinLess: shared.Odd{Value: 2.1}, WinMore: shared.Odd{Value: 1.62}},
		"195.5": {WinLess: shared.Odd{Value: 1.83}, WinMore: shared.Odd{Value: 1.87}},
		"196.5": {WinLess: shared.Odd{Value: 1.6}, WinMore: shared.Odd{Value: 2.15}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestBasketballMatch1FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"65.5": {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.75}},
		"66.5": {WinLess: shared.Odd{Value: 1.72}, WinMore: shared.Odd{Value: 1.95}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestBasketballMatch2FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"48.5": {WinLess: shared.Odd{Value: 1.92}, WinMore: shared.Odd{Value: 1.75}},
		"49.5": {WinLess: shared.Odd{Value: 1.67}, WinMore: shared.Odd{Value: 2.}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestBasketballMatch3FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"108.5": {WinLess: shared.Odd{Value: 1.92}, WinMore: shared.Odd{Value: 1.75}},
		"109.5": {WinLess: shared.Odd{Value: 1.67}, WinMore: shared.Odd{Value: 2.}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestBasketballMatch4FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"89.5": {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.75}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestBasketballMatch1SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"70.5": {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.75}},
		"71.5": {WinLess: shared.Odd{Value: 1.72}, WinMore: shared.Odd{Value: 1.95}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestBasketballMatch2SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"51.5": {WinLess: shared.Odd{Value: 2.05}, WinMore: shared.Odd{Value: 1.63}},
		"52.5": {WinLess: shared.Odd{Value: 1.83}, WinMore: shared.Odd{Value: 1.83}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestBasketballMatch3SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"112.5": {WinLess: shared.Odd{Value: 1.92}, WinMore: shared.Odd{Value: 1.75}},
		"113.5": {WinLess: shared.Odd{Value: 1.68}, WinMore: shared.Odd{Value: 2.}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestBasketballMatch4SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"105.5": {WinLess: shared.Odd{Value: 1.93}, WinMore: shared.Odd{Value: 1.75}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestBasketballMatch1Time1Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	assert.NotEqual(t, 0, len(period.Totals), "There are no totals")

	expected := map[string]*shared.WinLessMore{
		"35.5": &shared.WinLessMore{WinLess: shared.Odd{Value: 1.82}, WinMore: shared.Odd{Value: 1.85}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestBasketballMatch2Time2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	assert.NotEqual(t, 0, len(period.Totals), "There are no totals")

	expected := map[string]*shared.WinLessMore{
		"21.5": &shared.WinLessMore{WinLess: shared.Odd{Value: 2.45}, WinMore: shared.Odd{Value: 1.45}},
		"23.5": &shared.WinLessMore{WinLess: shared.Odd{Value: 1.8}, WinMore: shared.Odd{Value: 1.85}},
		"25.5": &shared.WinLessMore{WinLess: shared.Odd{Value: 1.45}, WinMore: shared.Odd{Value: 2.4}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestBasketballMatch3Time3Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]

	assert.NotEqual(t, 0, len(period.Totals), "There are no totals")

	expected := map[string]*shared.WinLessMore{
		"55.5": &shared.WinLessMore{WinLess: shared.Odd{Value: 1.92}, WinMore: shared.Odd{Value: 1.75}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time3Totals")
}

func TestBasketballMatch4Time4Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[4]

	assert.NotEqual(t, 0, len(period.Totals), "There are no totals")

	expected := map[string]*shared.WinLessMore{
		"44.5": {WinLess: shared.Odd{Value: 1.82}, WinMore: shared.Odd{Value: 1.88}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time4Totals")
}

func TestBasketballMatch1Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 2.1}, Win2: shared.Odd{Value: 0}},
		"3.5":  {Win1: shared.Odd{Value: 2.}, Win2: shared.Odd{Value: 0}},
		"4.5":  {Win1: shared.Odd{Value: 1.87}, Win2: shared.Odd{Value: 0}},
		"5.5":  {Win1: shared.Odd{Value: 1.77}, Win2: shared.Odd{Value: 0}},
		"6.5":  {Win1: shared.Odd{Value: 1.65}, Win2: shared.Odd{Value: 0}},
		"-2.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.6}},
		"-3.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.7}},
		"-4.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.8}},
		"-5.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.93}},
		"-6.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 2.05}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch2Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"1.5":  {Win1: shared.Odd{Value: 2.25}, Win2: shared.Odd{Value: 0}},
		"2.5":  {Win1: shared.Odd{Value: 2.05}, Win2: shared.Odd{Value: 0}},
		"3.5":  {Win1: shared.Odd{Value: 1.9}, Win2: shared.Odd{Value: 0}},
		"4.5":  {Win1: shared.Odd{Value: 1.75}, Win2: shared.Odd{Value: 0}},
		"5.5":  {Win1: shared.Odd{Value: 1.6}, Win2: shared.Odd{Value: 0}},
		"-1.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.53}},
		"-2.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.65}},
		"-3.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.8}},
		"-4.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.95}},
		"-5.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 2.1}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch3Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 2.05}, Win2: shared.Odd{Value: 0}},
		"3.5":  {Win1: shared.Odd{Value: 1.93}, Win2: shared.Odd{Value: 0}},
		"4.5":  {Win1: shared.Odd{Value: 1.8}, Win2: shared.Odd{Value: 0}},
		"5.5":  {Win1: shared.Odd{Value: 1.68}, Win2: shared.Odd{Value: 0}},
		"6.5":  {Win1: shared.Odd{Value: 1.55}, Win2: shared.Odd{Value: 0}},
		"-2.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.65}},
		"-3.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.8}},
		"-4.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.92}},
		"-5.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 2.05}},
		"-6.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 2.2}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch4Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")
	expecteds := map[string]*shared.WinHandicap{
		"15.5":  {Win1: shared.Odd{Value: 1.95}, Win2: shared.Odd{Value: 0}},
		"16.5":  {Win1: shared.Odd{Value: 1.75}, Win2: shared.Odd{Value: 0}},
		"17.5":  {Win1: shared.Odd{Value: 1.62}, Win2: shared.Odd{Value: 0}},
		"-15.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.77}},
		"-16.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.95}},
		"-17.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.15}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch1HandicapTime1(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile1Name, basketballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"6.5":  {Win1: shared.Odd{Value: 1.77}, Win2: shared.Odd{Value: 0}},
		"-6.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.9}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch2HandicapTime2(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile2Name, basketballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"1.5":  {Win1: shared.Odd{Value: 2.}, Win2: shared.Odd{Value: 0}},
		"-1.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.68}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch3HandicapTime3(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile3Name, basketballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 2.05}, Win2: shared.Odd{Value: 0}},
		"-2.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.65}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestBasketballMatch4HandicapTime4(t *testing.T) {

	match := suite.GetLobbetMatch(t, basketballFile4Name, basketballMatch4Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[4]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")
	expecteds := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.95}},
		"-2.5": {Win1: shared.Odd{Value: 1.75}, Win2: shared.Odd{Value: 0}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}
