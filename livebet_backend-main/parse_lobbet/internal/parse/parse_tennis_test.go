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
	// Теннисный матч, в котором идет SET1
	tennisFile1Name       = "tennis_match1.json"
	tennisMatch1Id  int64 = 22357517

	// Теннисный матч, в котором идет SET2
	tennisFile2Name       = "tennis_match2.json"
	tennisMatch2Id  int64 = 22357578

	// Теннисный матч, в котором идет SET3
	tennisFile3Name       = "tennis_match3.json"
	tennisMatch3Id  int64 = 22652515
)

func TestTennisMatch1Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "mohamed safwat", resGame.HomeName, "HomeName")
	assert.Equal(t, "marek gengel", resGame.AwayName, "AwayName")
	assert.Equal(t, "itf men sharm el sheikh", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch1Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch1Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
}

func TestTennisMatch2Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "jodie burrage", resGame.HomeName, "HomeName")
	assert.Equal(t, "amandine hesse", resGame.AwayName, "AwayName")
	assert.Equal(t, "itf women dubai", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch2Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch2Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
}

func TestTennisMatch3Name(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "madison keys", resGame.HomeName, "HomeName")
	assert.Equal(t, "clara tauson", resGame.AwayName, "AwayName")
	assert.Equal(t, "wta auckland", resGame.LeagueName, "LeagueName")

	assert.Equal(t, 1., resGame.HomeScore, "HomeScore")
	assert.Equal(t, 2., resGame.AwayScore, "AwayScore")

	assert.Equal(t, tennisMatch3Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch3Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
}

func TestTennisMatch1Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 6.4},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 1.08},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1.25},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 3.2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch1Time1Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 7.8},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 1.04},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestTennisMatch2Time1Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestTennisMatch1Time2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 3.95},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 1.2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestTennisMatch2Time2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]
	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.25},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 1.5},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestTennisMatch3Time2Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]
	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1.42},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 2.55},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestTennisMatch3Time3Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]
	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 3.42},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 4.55},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time3Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time3Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time3Win1x2.Win2")
}

func TestTennisMatch3Time4Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[4]
	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 4.42},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 5.55},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time4Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time4Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time4Win1x2.Win2")
}

func TestTennisMatch3Time5Win1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[5]
	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 5.42},
		WinNone: shared.Odd{Value: .0},
		Win2:    shared.Odd{Value: 6.55},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time5Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time5Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time5Win1x2.Win2")
}

func TestTennisMatch1GameWin1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.Win1x2Struct{
		"5": {Win1: shared.Odd{Value: 4.95}, Win2: shared.Odd{Value: 1.13}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch2GameWin1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.Win1x2Struct{
		"6":  {Win1: shared.Odd{Value: 2.4}, Win2: shared.Odd{Value: 1.45}},
		"10": {Win1: shared.Odd{Value: 2.41}, Win2: shared.Odd{Value: 1.451}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch3GameWin1x2(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]

	expected := map[string]*shared.Win1x2Struct{}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch1Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"17.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.85}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"27.5": {WinLess: shared.Odd{Value: 1.75}, WinMore: shared.Odd{Value: 1.85}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch1FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"5.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.85}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch2FirstTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.FirstTeamTotals), "There are no FirstTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"13.5": {WinLess: shared.Odd{Value: 1.8}, WinMore: shared.Odd{Value: 1.8}},
	}
	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch1SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"12.5": {WinLess: shared.Odd{Value: 1.43}, WinMore: shared.Odd{Value: 2.6}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch2SecondTeamTotals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.SecondTeamTotals), "There are no SecondTeamTotals")
	expected := map[string]*shared.WinLessMore{
		"12.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.75}},
	}
	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch1Time1Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	assert.NotEqual(t, 0, len(period.Totals), "There are no SET1 totals")

	expected := map[string]*shared.WinLessMore{
		"6.5":  {WinLess: shared.Odd{Value: 5.4}, WinMore: shared.Odd{Value: 1.11}},
		"7.5":  {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.8}},
		"8.5":  {WinLess: shared.Odd{Value: 1.48}, WinMore: shared.Odd{Value: 2.45}},
		"9.5":  {WinLess: shared.Odd{Value: 1.11}, WinMore: shared.Odd{Value: 5.4}},
		"10.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 7.85}},
		"12.5": {WinLess: shared.Odd{Value: 1.02}, WinMore: shared.Odd{Value: 9.85}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestTennisMatch2Time1Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	assert.Equal(t, 0, len(period.Totals), "There are no SET1 totals")

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestTennisMatch1Time2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	assert.Equal(t, 0, len(period.Totals), "There are no SET2 totals")

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestTennisMatch2Time2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET2 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5":  {WinLess: shared.Odd{Value: 2.7}, WinMore: shared.Odd{Value: 1.35}},
		"10.5": {WinLess: shared.Odd{Value: 1.37}, WinMore: shared.Odd{Value: 2.65}},
		"12.5": {WinLess: shared.Odd{Value: 1.17}, WinMore: shared.Odd{Value: 3.9}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestTennisMatch3Time2Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET2 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5":  {WinLess: shared.Odd{Value: 3.6}, WinMore: shared.Odd{Value: 1.23}},
		"10.5": {WinLess: shared.Odd{Value: 1.65}, WinMore: shared.Odd{Value: 2.05}},
		"12.5": {WinLess: shared.Odd{Value: 1.25}, WinMore: shared.Odd{Value: 3.35}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestTennisMatch3Time3Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[3]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET3 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5":  {WinLess: shared.Odd{Value: 3.6}, WinMore: shared.Odd{Value: 1.23}},
		"10.5": {WinLess: shared.Odd{Value: 3.65}, WinMore: shared.Odd{Value: 2.05}},
		"12.5": {WinLess: shared.Odd{Value: 1.25}, WinMore: shared.Odd{Value: 4.35}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time3Totals")
}

func TestTennisMatch3Time4Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[4]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET4 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5":  {WinLess: shared.Odd{Value: 4.6}, WinMore: shared.Odd{Value: 1.23}},
		"10.5": {WinLess: shared.Odd{Value: 1.65}, WinMore: shared.Odd{Value: 5.05}},
		"12.5": {WinLess: shared.Odd{Value: 6.25}, WinMore: shared.Odd{Value: 3.35}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time4Totals")
}

func TestTennisMatch3Time5Totals(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[5]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET5 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5":  {WinLess: shared.Odd{Value: 5.6}, WinMore: shared.Odd{Value: 1.23}},
		"10.5": {WinLess: shared.Odd{Value: 6.65}, WinMore: shared.Odd{Value: 2.05}},
		"12.5": {WinLess: shared.Odd{Value: 7.25}, WinMore: shared.Odd{Value: 3.35}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time5Totals")
}

func TestTennisMatch1Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")
	expecteds := map[string]*shared.WinHandicap{
		"6.5":  {Win1: shared.Odd{Value: 2.05}, Win2: shared.Odd{Value: 0}},
		"-6.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.65}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestTennisMatch2Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 2.05}},
		"-2.5": {Win1: shared.Odd{Value: 1.6}, Win2: shared.Odd{Value: 0}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}

func TestTennisMatch3Handicap(t *testing.T) {

	match := suite.GetLobbetMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expecteds := map[string]*shared.WinHandicap{
		"1.5":  {Win1: shared.Odd{Value: 1.9}, Win2: shared.Odd{Value: 0}},
		"-1.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.8}},
	}

	suite.CheckHandicap(t, expecteds, period.Handicap, "Handicap")
}
