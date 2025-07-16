package parse

import (
	"livebets/parse_lobbet/internal/parse/suite"
	"livebets/shared"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// Футбольный матч1
	footballFile1Name       = "football_match1.json"
	footballMatch1Id  int64 = 22487824
	// Футбольный матч2
	footballFile2Name       = "football_match2.json"
	footballMatch2Id  int64 = 22487825
	// Футбольный матч3
	footballFile3Name       = "football_match3.json"
	footballMatch3Id  int64 = 22509407
)

func TestFootballMatch1Name(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "huesca", resGame.HomeName, "HomeName")
	assert.Equal(t, "tenerife", resGame.AwayName, "AwayName")
	assert.Equal(t, "spain laliga2", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch1Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch1Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
}

func TestFootballMatch2Name(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "dep. la coruna", resGame.HomeName, "HomeName")
	assert.Equal(t, "castellon", resGame.AwayName, "AwayName")
	assert.Equal(t, "spain laliga2", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch2Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch2Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
}

func TestFootballMatch3Name(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	assert.Equal(t, "cardiff", resGame.HomeName, "HomeName")
	assert.Equal(t, "watford", resGame.AwayName, "AwayName")
	assert.Equal(t, "england premier league cup", resGame.LeagueName, "LeagueName")

	assert.Equal(t, footballMatch3Id, resGame.Pid, "Pid")
	assert.Equal(t, strconv.FormatInt(footballMatch3Id, 10), resGame.MatchId, "MatchId")

	assert.Equal(t, shared.SOCCER, resGame.SportName, "SportName")
}

func TestFootballMatch1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1.3},
		WinNone: shared.Odd{Value: 4.75},
		Win2:    shared.Odd{Value: 13.5},
	}

	period := resGame.Periods[0]

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 3.95},
		WinNone: shared.Odd{Value: 3.6},
		Win2:    shared.Odd{Value: 1.9},
	}

	period := resGame.Periods[0]

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch3Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 0},
		WinNone: shared.Odd{Value: 0},
		Win2:    shared.Odd{Value: 0},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Win1x2.Win2")
}

func TestFootballMatch1Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 1.06},
		WinNone: shared.Odd{Value: 8},
		Win2:    shared.Odd{Value: 68},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch2Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 30},
		WinNone: shared.Odd{Value: 5.3},
		Win2:    shared.Odd{Value: 1.17},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch3Time1Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 0},
		WinNone: shared.Odd{Value: 0},
		Win2:    shared.Odd{Value: 0},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time1Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time1Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time1Win1x2.Win2")
}

func TestFootballMatch1Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2.85},
		WinNone: shared.Odd{Value: 2.13},
		Win2:    shared.Odd{Value: 3.8},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestFootballMatch2Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 2},
		WinNone: shared.Odd{Value: 2.9},
		Win2:    shared.Odd{Value: 4.15},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestFootballMatch3Time2Win1x2(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := shared.Win1x2Struct{
		Win1:    shared.Odd{Value: 0},
		WinNone: shared.Odd{Value: 0},
		Win2:    shared.Odd{Value: 0},
	}

	assert.Equal(t, expected.Win1, period.Win1x2.Win1, "Time2Win1x2.Win1")
	assert.Equal(t, expected.WinNone, period.Win1x2.WinNone, "Time2Win1x2.WinNone")
	assert.Equal(t, expected.Win2, period.Win1x2.Win2, "Time2Win1x2.Win2")
}

func TestFootballMatch1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 3.3}, WinMore: shared.Odd{Value: 1.3}},
		"2.5": {WinLess: shared.Odd{Value: 1.5}, WinMore: shared.Odd{Value: 2.45}},
		"3.5": {WinLess: shared.Odd{Value: 1.11}, WinMore: shared.Odd{Value: 6}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 7.1}, WinMore: shared.Odd{Value: 1.07}},
		"2.5": {WinLess: shared.Odd{Value: 2.5}, WinMore: shared.Odd{Value: 1.5}},
		"3.5": {WinLess: shared.Odd{Value: 1.5}, WinMore: shared.Odd{Value: 2.55}},
		"4.5": {WinLess: shared.Odd{Value: 1.15}, WinMore: shared.Odd{Value: 5}},
		"5.5": {WinLess: shared.Odd{Value: 1.02}, WinMore: shared.Odd{Value: 11}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch3Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"4.5": {WinLess: shared.Odd{Value: 1.12}, WinMore: shared.Odd{Value: 4.85}},
	}
	suite.CheckTotals(t, expected, period.Totals, "Totals")
}

func TestFootballMatch1Time1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.2}, WinMore: shared.Odd{Value: 4.25}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch2Time1Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.43}, WinMore: shared.Odd{Value: 2.65}},
		"2.5": {WinLess: shared.Odd{Value: 1.01}, WinMore: shared.Odd{Value: 11}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time1Totals")
}

func TestFootballMatch1Time2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 2.65}, WinMore: shared.Odd{Value: 1.45}},
		"1.5": {WinLess: shared.Odd{Value: 1.33}, WinMore: shared.Odd{Value: 3.2}},
		"2.5": {WinLess: shared.Odd{Value: 1.04}, WinMore: shared.Odd{Value: 9.4}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestFootballMatch2Time2Totals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 5.05}, WinMore: shared.Odd{Value: 1.15}},
		"1.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.82}},
		"2.5": {WinLess: shared.Odd{Value: 1.25}, WinMore: shared.Odd{Value: 3.65}},
		"3.5": {WinLess: shared.Odd{Value: 1.05}, WinMore: shared.Odd{Value: 8.8}},
	}

	suite.CheckTotals(t, expected, period.Totals, "Time2Totals")
}

func TestFootballMatch1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.95}, WinMore: shared.Odd{Value: 1.8}},
		"2.5": {WinLess: shared.Odd{Value: 1.15}, WinMore: shared.Odd{Value: 4.95}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestFootballMatch2FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 3.45}, WinMore: shared.Odd{Value: 1.3}},
		"1.5": {WinLess: shared.Odd{Value: 1.55}, WinMore: shared.Odd{Value: 2.37}},
		"2.5": {WinLess: shared.Odd{Value: 1.13}, WinMore: shared.Odd{Value: 5.3}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "FirstTeamTotals")
}

func TestFootballMatch1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.65}, WinMore: shared.Odd{Value: 2.15}},
		"1.5": {WinLess: shared.Odd{Value: 1.07}, WinMore: shared.Odd{Value: 6.9}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestFootballMatch2SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 2.05}, WinMore: shared.Odd{Value: 1.73}},
		"2.5": {WinLess: shared.Odd{Value: 1.17}, WinMore: shared.Odd{Value: 4.55}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "SecondTeamTotals")
}

func TestFootballMatch1Time1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.11}, WinMore: shared.Odd{Value: 5.7}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "Time1FirstTeamTotals")
}

func TestFootballMatch2Time1FirstTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.25}, WinMore: shared.Odd{Value: 3.75}},
	}

	suite.CheckTotals(t, expected, period.FirstTeamTotals, "Time1FirstTeamTotals")
}

func TestFootballMatch1Time1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"0.5": {WinLess: shared.Odd{Value: 1.07}, WinMore: shared.Odd{Value: 7.2}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "Time1SecondTeamTotals")
}

func TestFootballMatch2Time1SecondTeamTotals(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinLessMore{
		"1.5": {WinLess: shared.Odd{Value: 1.15}, WinMore: shared.Odd{Value: 4.95}},
	}

	suite.CheckTotals(t, expected, period.SecondTeamTotals, "Time1SecondTeamTotals")
}

func TestFootballMatch1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 2.75}, Win2: shared.Odd{Value: 0}},
		"0.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 3.45}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinHandicap{
		"0.5":  {Win1: shared.Odd{Value: 1.87}, Win2: shared.Odd{Value: 0}},
		"-1.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 4.25}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch3Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[0]

	expected := map[string]*shared.WinHandicap{
		"1.5":  {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 11}},
		"-2.5": {Win1: shared.Odd{Value: 11}, Win2: shared.Odd{Value: 0}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Handicap")
}

func TestFootballMatch1Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 8}, Win2: shared.Odd{Value: 0}},
		"0.5":  {Win1: shared.Odd{Value: 0.}, Win2: shared.Odd{Value: 9.2}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Time1Handicap")
}

func TestFootballMatch2Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{
		"-1.5": {Win1: shared.Odd{Value: 0}, Win2: shared.Odd{Value: 6.9}},
		"0.5":  {Win1: shared.Odd{Value: 5.1}, Win2: shared.Odd{Value: 0}},
	}

	suite.CheckHandicap(t, expected, period.Handicap, "Time1Handicap")
}

func TestFootballMatch3Time1Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[1]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Time1Handicap")
}

func TestFootballMatch1Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile1Name, footballMatch1Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Time2Handicap")
}

func TestFootballMatch2Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile2Name, footballMatch2Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Time2Handicap")
}

func TestFootballMatch3Time2Handicap(t *testing.T) {
	t.Parallel()

	match := suite.GetLobbetMatch(t, footballFile3Name, footballMatch3Id)

	resGame := LiveToResponseGame(match)

	period := resGame.Periods[2]

	expected := map[string]*shared.WinHandicap{}

	suite.CheckHandicap(t, expected, period.Handicap, "Time2Handicap")
}
