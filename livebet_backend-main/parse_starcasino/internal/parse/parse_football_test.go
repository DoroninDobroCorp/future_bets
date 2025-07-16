package parse

import (
	"github.com/stretchr/testify/assert"
	"livebets/parse_starcasino/internal/parse/suite"
	"livebets/shared"
	"strconv"
	"testing"
)

const (
	// Футбольный матч1
	footballFile1Name       = "football_match1.json"
	footballMatch1Id  int64 = 11322681
	// Футбольный матч2
	footballFile2Name       = "football_match2.json"
	footballMatch2Id  int64 = 11174143
	// Футбольный матч3
	footballFile3Name       = "football_match3.json"
	footballMatch3Id  int64 = 11335344
	// Футбольный матч4 v.3.0
	footballFile4Name       = "football_match4.json"
	footballMatch4Id  int64 = 31030
	// Футбольный матч5 v.3.0
	footballFile5Name       = "football_match5.json"
	footballMatch5Id  int64 = 39003
	// Футбольный матч6 v.3.0
	footballFile6Name       = "football_match6.json"
	footballMatch6Id  int64 = 16633
	// Футбольный матч7 v.3.0 + "1ST HALF - HANDICAP"
	footballFile7Name       = "football_match7.json"
	footballMatch7Id  int64 = 6549
	// Футбольный матч8 prematch v.3.0
	footballFile8Name       = "football_match8.json"
	footballMatch8Id  int64 = 8170
)

func TestFootballMatch1Name(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "melbourne knights", resGame.HomeName, "HomeName")
	assert.Equal(t, "heidelberg united", resGame.AwayName, "AwayName")
	assert.Equal(t, "australia victoria npl", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch1Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch1Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestFootballMatch2Name(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "aparecidense/go", resGame.HomeName, "HomeName")
	assert.Equal(t, "crac/go", resGame.AwayName, "AwayName")
	assert.Equal(t, "brazil goiano", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch2Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch2Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestFootballMatch3Name(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	assert.Equal(t, "blooming", resGame.HomeName, "HomeName")
	assert.Equal(t, "el nacional quito", resGame.AwayName, "AwayName")
	assert.Equal(t, "americas copa libertadores", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch3Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch3Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
	assert.False(t, resGame.CreatedAt.IsZero(), "CreatedAt.IsZero()")
	assert.Equal(t, shared.STARCASINO, resGame.Source, "Source")
}

func TestFootballMatch1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 4.5},
		WinNone: shared.Odd{Value: 1.8182},
		Win2:    shared.Odd{Value: 3.25},
	}

	period := resGame.Periods[0]

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 5.6667},
		WinNone: shared.Odd{Value: 3.5},
		Win2:    shared.Odd{Value: 1.6154},
	}

	period := resGame.Periods[0]

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch3Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 5.6667},
		WinNone: shared.Odd{Value: 3.5},
		Win2:    shared.Odd{Value: 1.6154},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch1Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch2Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 26},
		WinNone: shared.Odd{Value: 4.5},
		Win2:    shared.Odd{Value: 1.2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch3Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1},
		WinNone: shared.Odd{Value: 41},
		Win2:    shared.Odd{Value: 151},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch1Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch2Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.5455},
		WinNone: shared.Odd{Value: 2.1429},
		Win2:    shared.Odd{Value: 4.2},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestFootballMatch3Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 3.1429},
		WinNone: shared.Odd{Value: 2.3637},
		Win2:    shared.Odd{Value: 2.9},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"4.5": {WinLess: shared.Odd{Value: 2.1}, WinMore: shared.Odd{Value: 1.6667}},
		"5.5": {WinLess: shared.Odd{Value: 1.1819}, WinMore: shared.Odd{Value: 4.2}},
		"6.5": {WinLess: shared.Odd{Value: 1.02}, WinMore: shared.Odd{Value: 9}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 4.}, WinMore: shared.Odd{Value: 1.2}},
		"2.5": {WinLess: shared.Odd{Value: 1.8}, WinMore: shared.Odd{Value: 1.9}},
		"3.5": {WinLess: shared.Odd{Value: 1.2223}, WinMore: shared.Odd{Value: 3.75}},
		"4.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 8.}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch3Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"4.5": {WinLess: shared.Odd{Value: 4}, WinMore: shared.Odd{Value: 1.2223}},
		"5.5": {WinLess: shared.Odd{Value: 1.75}, WinMore: shared.Odd{Value: 2.05}},
		"6.5": {WinLess: shared.Odd{Value: 1.2}, WinMore: shared.Odd{Value: 4.2}},
		"7.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 8.5}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch1Time1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch2Time1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.4546}, WinMore: shared.Odd{Value: 2.4}},
		"2.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 7.}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch3Time1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"4.5": {WinLess: shared.Odd{Value: 1.1667}, WinMore: shared.Odd{Value: 4.75}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch1Time2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestFootballMatch2Time2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.3637}, WinMore: shared.Odd{Value: 2.6667}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestFootballMatch3Time2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 3.4}, WinMore: shared.Odd{Value: 1.3}},
		"1.5": {WinLess: shared.Odd{Value: 1.5556}, WinMore: shared.Odd{Value: 2.4}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"2.5": {WinLess: shared.Odd{Value: 1.3334}, WinMore: shared.Odd{Value: 3.}},
		"3.5": {WinLess: shared.Odd{Value: 1.02}, WinMore: shared.Odd{Value: 9.}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestFootballMatch2FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 2.4}, WinMore: shared.Odd{Value: 1.5}},
		"1.5": {WinLess: shared.Odd{Value: 1.2728}, WinMore: shared.Odd{Value: 3.4}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestFootballMatch3FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"3.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.75}},
		"4.5": {WinLess: shared.Odd{Value: 1.1539}, WinMore: shared.Odd{Value: 4.5}},
		"5.5": {WinLess: shared.Odd{Value: 1.01}, WinMore: shared.Odd{Value: 10.}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestFootballMatch1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"2.5": {WinLess: shared.Odd{Value: 1.5}, WinMore: shared.Odd{Value: 2.4}},
		"3.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 8.}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestFootballMatch2SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.75}, WinMore: shared.Odd{Value: 1.95}},
		"2.5": {WinLess: shared.Odd{Value: 1.1112}, WinMore: shared.Odd{Value: 5.3334}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestFootballMatch3SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 2.05}, WinMore: shared.Odd{Value: 1.7}},
		"2.5": {WinLess: shared.Odd{Value: 1.1667}, WinMore: shared.Odd{Value: 4.3334}},
		"3.5": {WinLess: shared.Odd{Value: 1.02}, WinMore: shared.Odd{Value: 9.}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestFootballMatch1Time1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "Time1FirstTeamTotals")
}

func TestFootballMatch1Time1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "Time1SecondTeamTotals")
}

func TestFootballMatch2Time1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.25}, WinMore: shared.Odd{Value: 3.25}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "Time1FirstTeamTotals")
}

func TestFootballMatch2Time1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.1429}, WinMore: shared.Odd{Value: 4.5}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "Time1SecondTeamTotals")
}

func TestFootballMatch3Time2FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.7}, WinMore: shared.Odd{Value: 2.05}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "Time2FirstTeamTotals")
}

func TestFootballMatch3Time2SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.8}, WinMore: shared.Odd{Value: 1.9}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "Time2SecondTeamTotals")
}

func TestFootballMatch1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	// current score = 1:0
	expected := map[string]*shared.WinHandicap{
		"-0.5": {Win1: shared.Odd{Value: 4.5}, Win2: shared.Odd{Value: 3.25}},
		"0.0":  {Win1: shared.Odd{Value: 2.2}, Win2: shared.Odd{Value: 1.6}},
		"0.5":  {Win1: shared.Odd{Value: 1.3}, Win2: shared.Odd{Value: 1.1667}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	// score = 2:0
	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 3.6}},
		"-0.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 1.6154}},
		"0.0":  {Win1: shared.Odd{Value: 4.}, Win2: shared.Odd{Value: 1.2}},
		"0.5":  {Win1: shared.Odd{Value: 2.1667}, Win2: shared.Odd{Value: 1.1112}},
		"1.5":  {Win1: shared.Odd{Value: 1.25}, Win2: shared.Odd{Value: 0.}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch3Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[0]

	// current score = 2:1
	expected := map[string]*shared.WinHandicap{
		"-3.5": {Win1: shared.Odd{Value: 7.5}, Win2: shared.Odd{Value: 0.}},
		"-2.5": {Win1: shared.Odd{Value: 2.9}, Win2: shared.Odd{Value: 0.}},
		"-1.5": {Win1: shared.Odd{Value: 1.4546}, Win2: shared.Odd{Value: 0.}},
		"-0.5": {Win1: shared.Odd{Value: 1.091}, Win2: shared.Odd{Value: 0.}},
		"0.0":  {Win1: shared.Odd{Value: 1.01}, Win2: shared.Odd{Value: 12.5}},
		"0.5":  {Win1: shared.Odd{Value: 1.}, Win2: shared.Odd{Value: 7.}},
		"1.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.6667}},
		"2.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.4}},
		"3.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.0625}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch1Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch2Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 5.3334}},
		"-0.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.2}},
		"0.5":  {Win1: shared.Odd{Value: 3.75}, Win2: shared.Odd{Value: 1.}},
		"1.5":  {Win1: shared.Odd{Value: 1.1}, Win2: shared.Odd{Value: 0.}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch3Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{
		"-2.5": {Win1: shared.Odd{Value: 8.}, Win2: shared.Odd{Value: 0.}},
		"-1.5": {Win1: shared.Odd{Value: 1.0625}, Win2: shared.Odd{Value: 0.}},
		"0.0":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 10.}},
		"0.5":  {Win1: shared.Odd{Value: 1.}, Win2: shared.Odd{Value: 12.}},
		"1.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 7.}},
		"2.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 1.}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch1Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile1Name, footballMatch1Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch2Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile2Name, footballMatch2Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{
		"0.0": {Win1: shared.Odd{Value: 1.5}, Win2: shared.Odd{Value: 2.3}},
		"0.5": {Win1: shared.Odd{Value: 1.1539}, Win2: shared.Odd{Value: 1.4167}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch3Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetStarCasinoMatch(t, footballFile3Name, footballMatch3Id)

	resGame := StarCasinoToResponseGame(*match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{
		"-0.5": {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 2.9}},
		"0.0":  {Win1: shared.Odd{Value: 1.9}, Win2: shared.Odd{Value: 1.8}},
		"0.5":  {Win1: shared.Odd{Value: 1.3637}, Win2: shared.Odd{Value: 1.3077}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}
