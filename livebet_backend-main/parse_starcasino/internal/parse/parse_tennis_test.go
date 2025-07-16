package parse

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"livebets/parse_starcasino/internal/parse/suite"
	"livebets/shared"
	"strconv"
	"testing"
)

const (
	// Теннисный матч, в котором идет SET1
	tennisFile1Name       = "tennis_match1.json"
	tennisMatch1Id  int64 = 11756621

	// Теннисный матч, в котором идет SET2
	tennisFile2Name       = "tennis_match2.json"
	tennisMatch2Id  int64 = 11757084

	// Теннисный матч, в котором идет SET3
	tennisFile3Name       = "tennis_match3.json"
	tennisMatch3Id  int64 = 11757087

	// Теннисный матч
	tennisFile4Name       = "tennis_match4.json"
	tennisMatch4Id  int64 = 11763584

	// Теннисный матч
	tennisFile5Name       = "tennis_match5.json"
	tennisMatch5Id  int64 = 11763769

	// Теннисный матч 6, v.3.0 - all bets in file
	tennisFile6Name       = "tennis_match6.json"
	tennisMatch6Id  int64 = 8052

	// Теннисный матч 7, prematch v.3.0 - all bets in file
	tennisFile7Name       = "tennis_match7.json"
	tennisMatch7Id  int64 = 8072
)

func TestTennisMatch1Name(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "carlos alcaraz", resGame.HomeName, "HomeName")
	assert.Equal(t, "hubert hurkacz", resGame.AwayName, "AwayName")
	assert.Equal(t, "atp rotterdam netherlands men", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch1Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch1Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestTennisMatch2Name(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "sumit nagal", resGame.HomeName, "HomeName")
	assert.Equal(t, "juan pablo ficovich", resGame.AwayName, "AwayName")
	assert.Equal(t, "atp buenos aires argentina men", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch2Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch2Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestTennisMatch3Name(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "thiago monteiro", resGame.HomeName, "HomeName")
	assert.Equal(t, "felipe alves", resGame.AwayName, "AwayName")
	assert.Equal(t, "atp buenos aires argentina men", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch3Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch3Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestTennisMatch4Name(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "matteo gigante", resGame.HomeName, "HomeName")
	assert.Equal(t, "hugo grenier", resGame.AwayName, "AwayName")
	assert.Equal(t, "atp marseille france men", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch4Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch4Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestTennisMatch5Name(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "clement tabur", resGame.HomeName, "HomeName")
	assert.Equal(t, "miguel damas", resGame.AwayName, "AwayName")
	assert.Equal(t, "challenger atp tenerife spain men", resGame.LeagueName, "LeagueName")

	assert.Equal(t, tennisMatch5Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(tennisMatch5Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.TENNIS, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestTennisMatch1Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.7},
		Win2: shared.Odd{Value: 2.2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch2Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.7},
		Win2: shared.Odd{Value: 2.1},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch3Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.75},
		Win2: shared.Odd{Value: 2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch4Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 2.55},
		Win2: shared.Odd{Value: 1.45},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch5Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.2},
		Win2: shared.Odd{Value: 4.15},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestTennisMatch1Set1Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 3.75},
		Win2: shared.Odd{Value: 1.25},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set1Win1x2.Win2")
}

func TestTennisMatch1Set2Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.45},
		Win2: shared.Odd{Value: 2.65},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set2Win1x2.Win2")
}

func TestTennisMatch2Set2Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]
	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.02},
		Win2: shared.Odd{Value: 10.5},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set2Win1x2.Win2")
}

func TestTennisMatch3Set3Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[3]
	expected := shared.Win1x2Struct{}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set3Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set3Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set3Win1x2.Win2")
}

func TestTennisMatch4Set1Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]
	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 2.55},
		Win2: shared.Odd{Value: 1.45},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set1Win1x2.Win2")
}

func TestTennisMatch4Set2Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]
	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 2.2},
		Win2: shared.Odd{Value: 1.6},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set2Win1x2.Win2")
}

func TestTennisMatch5Set1Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]
	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.14},
		Win2: shared.Odd{Value: 4.95},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set1Win1x2.Win2")
}

func TestTennisMatch5Set2Win1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]
	expected := shared.Win1x2Struct{
		Win1: shared.Odd{Value: 1.35},
		Win2: shared.Odd{Value: 2.85},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Set2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Set2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Set2Win1x2.Win2")
}

func TestTennisMatch1GamesWin1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]
	expected := map[string]*shared.Win1x2Struct{
		"7": {Win1: shared.Odd{Value: 3.4}, Win2: shared.Odd{Value: 1.2858}},
		"8": {Win1: shared.Odd{Value: 1.1}, Win2: shared.Odd{Value: 6.5}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch2GamesWin1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]
	expected := map[string]*shared.Win1x2Struct{
		"8": {Win1: shared.Odd{Value: 1.3572}, Win2: shared.Odd{Value: 3}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set2Game")
}

func TestTennisMatch3GamesWin1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[3]
	expected := map[string]*shared.Win1x2Struct{
		"5": {Win1: shared.Odd{Value: 1.1}, Win2: shared.Odd{Value: 6.5}},
		"6": {Win1: shared.Odd{Value: 6}, Win2: shared.Odd{Value: 1.1112}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch4GamesWin1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]
	expected := map[string]*shared.Win1x2Struct{
		"4": {Win1: shared.Odd{Value: 1.3334}, Win2: shared.Odd{Value: 3.1}},
		"5": {Win1: shared.Odd{Value: 4.3}, Win2: shared.Odd{Value: 1.1905}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch5GamesWin1x2(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]
	expected := map[string]*shared.Win1x2Struct{
		"4": {Win1: shared.Odd{Value: 1.1667}, Win2: shared.Odd{Value: 4.5}},
		"5": {Win1: shared.Odd{Value: 3}, Win2: shared.Odd{Value: 1.3572}},
	}

	suite.CheckGamesWin1x2(t, expected, period.Games, "Set1Game")
}

func TestTennisMatch1Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"24.5": {WinLess: shared.Odd{Value: 2.1}, WinMore: shared.Odd{Value: 1.6667}},
		"25.5": {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.8}},
		"26.5": {WinLess: shared.Odd{Value: 1.6667}, WinMore: shared.Odd{Value: 2.1}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch2Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"30.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.75}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch3Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"31.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.75}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch4Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"22.5": {WinLess: shared.Odd{Value: 1.8334}, WinMore: shared.Odd{Value: 1.8334}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch5Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]
	require.NotEqual(t, 0, len(period.Totals), "There are no match totals")

	expected := map[string]*shared.WinLessMore{
		"20.5": {WinLess: shared.Odd{Value: 1.8334}, WinMore: shared.Odd{Value: 1.8334}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestTennisMatch1Set1Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET3 totals")

	expected := map[string]*shared.WinLessMore{
		"8.5": {WinLess: shared.Odd{Value: 7}, WinMore: shared.Odd{Value: 1.0625}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Set1Totals")
}

func TestTennisMatch2Set2Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET3 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5": {WinLess: shared.Odd{Value: 1.2}, WinMore: shared.Odd{Value: 4}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Set2Totals")
}

func TestTennisMatch3Set3Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[3]

	require.Equal(t, 0, len(period.Totals), "There are no SET3 totals")

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.Totals, "Set3Totals")
}

func TestTennisMatch4Set1Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET4 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5": {WinLess: shared.Odd{Value: 2.1}, WinMore: shared.Odd{Value: 1.6667}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Set1Totals")
}

func TestTennisMatch5Set1Totals(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Totals), "There are no SET5 totals")

	expected := map[string]*shared.WinLessMore{
		"9.5": {WinLess: shared.Odd{Value: 1.9}, WinMore: shared.Odd{Value: 1.8}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Set5Totals")
}

func TestTennisMatch1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"13.5": {WinLess: shared.Odd{Value: 1.8334}, WinMore: shared.Odd{Value: 1.8334}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"12.5": {WinLess: shared.Odd{Value: 2.05}, WinMore: shared.Odd{Value: 1.7}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch2FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"17.5": {WinLess: shared.Odd{Value: 2.3}, WinMore: shared.Odd{Value: 1.5}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch2SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"14.5": {WinLess: shared.Odd{Value: 1.85}, WinMore: shared.Odd{Value: 1.8}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch3FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"16.5": {WinLess: shared.Odd{Value: 1.3}, WinMore: shared.Odd{Value: 3}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch3SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"16.5": {WinLess: shared.Odd{Value: 2.6667}, WinMore: shared.Odd{Value: 1.3637}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch4FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"11.5": {WinLess: shared.Odd{Value: 2}, WinMore: shared.Odd{Value: 1.6667}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch4SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"12.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.7}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch5FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"12.5": {WinLess: shared.Odd{Value: 1.7}, WinMore: shared.Odd{Value: 1.95}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestTennisMatch5SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"8.5": {WinLess: shared.Odd{Value: 1.8}, WinMore: shared.Odd{Value: 1.85}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestTennisMatch1Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"2.5":  {Win1: shared.Odd{Value: 1.4546}, Win2: shared.Odd{Value: 0.}},
		"1.5":  {Win1: shared.Odd{Value: 1.5556}, Win2: shared.Odd{Value: 0.}},
		"0.5":  {Win1: shared.Odd{Value: 1.6667}, Win2: shared.Odd{Value: 0.}},
		"-2.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.5455}},
		"-1.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.2858}},
		"-0.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.1}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch2Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-3.5": {Win1: shared.Odd{Value: 1.8334}, Win2: shared.Odd{Value: 0.}},
		"3.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.8334}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch3Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")
	expected := map[string]*shared.WinHandicap{
		"0.5":  {Win1: shared.Odd{Value: 1.75}, Win2: shared.Odd{Value: 0.}},
		"-0.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.95}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch4Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"3.5":  {Win1: shared.Odd{Value: 1.7}, Win2: shared.Odd{Value: 0.}},
		"-3.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.05}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch5Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-4.5": {Win1: shared.Odd{Value: 1.8334}, Win2: shared.Odd{Value: 0.}},
		"4.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.8334}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch1Set1Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile1Name, tennisMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 7.}, Win2: shared.Odd{Value: 0.}},
		"1.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.0625}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch2Set2Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile2Name, tennisMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-2.5": {Win1: shared.Odd{Value: 1.2}, Win2: shared.Odd{Value: 0.}},
		"2.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 4}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch3Set3Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile3Name, tennisMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[3]

	require.Equal(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch4Set1Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile4Name, tennisMatch4Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.75}},
		"1.5":  {Win1: shared.Odd{Value: 1.95}, Win2: shared.Odd{Value: 0.}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestTennisMatch5Set1Handicap(t *testing.T) {

	match := suite.GetStarCasinoMatch(t, tennisFile5Name, tennisMatch5Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	require.NotEqual(t, 0, len(period.Handicap), "There is no Handicap")

	expected := map[string]*shared.WinHandicap{
		"-2.5": {Win1: shared.Odd{Value: 2.05}, Win2: shared.Odd{Value: 0.}},
		"2.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.7}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}
