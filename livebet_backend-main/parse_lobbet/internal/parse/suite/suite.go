package suite

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"livebets/parse_lobbet/internal/entity"
	"livebets/shared"
	"os"
	"path"
	"testing"
)

const jsonPath = "test_json"

func getFile(fileName string) ([]byte, error) {
	fullFileName := path.Join(jsonPath, fileName)

	return os.ReadFile(fullFileName)
}

func GetLobbetMatch(t *testing.T, fileName string, matchID int64) entity.Match {

	body, err := getFile(fileName)
	if err != nil {
		t.Fatalf("[ERROR] Не удалось прочитать файл %s. Err: %s", fileName, err)
	}

	var apiResponse entity.ResponseMatchData
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("[ERROR] Не удалось Unmarshal файл %s. Err: %s", fileName, err)
	}

	matches := apiResponse.Live.Matches

	if len(matches) == 0 {
		t.Fatalf("[ERROR] Файл не содержит матчей. File: %s", fileName)
	}

	match := matches[0]
	if match.ID != matchID {
		t.Fatalf("[ERROR] Файл не содержит матча ID: %d. File: %s", matchID, fileName)
	}

	return *match
}

func CheckTotals(t *testing.T, expectedTotals, totals map[string]*shared.WinLessMore, caption string) {
	assert.Equal(t, len(expectedTotals), len(totals), "len("+caption+")")

	for key, total := range totals {
		expected, ok := expectedTotals[key]
		if !ok {
			expected = &shared.WinLessMore{}
		}
		assert.Equal(t, expected.WinLess, total.WinLess, caption+".WinLess <"+key)
		assert.Equal(t, expected.WinMore, total.WinMore, caption+".WinMore >"+key)
	}
}

func CheckHandicap(t *testing.T, expectedHandicap, matchHandicap map[string]*shared.WinHandicap, caption string) {
	assert.Equal(t, len(expectedHandicap), len(matchHandicap), "len("+caption+")")

	for key, handicap := range matchHandicap {
		expected, ok := expectedHandicap[key]
		if !ok {
			expected = &shared.WinHandicap{}
		}
		assert.Equal(t, expected.Win1, handicap.Win1, caption+".Win1 "+key)
		assert.Equal(t, expected.Win2, handicap.Win2, caption+".Win2 "+key)
	}
}

func CheckGamesWin1x2(t *testing.T, expectedWin1x2, gamesWin1x2 map[string]*shared.Win1x2Struct, caption string) {
	assert.Equal(t, len(expectedWin1x2), len(gamesWin1x2), "len("+caption+")")

	for key, win1x2 := range gamesWin1x2 {
		expected, ok := expectedWin1x2[key]
		if !ok {
			expected = &shared.Win1x2Struct{}
		}
		assert.Equal(t, expected.Win1, win1x2.Win1, caption+key+".Win1")
		assert.Equal(t, expected.Win2, win1x2.Win2, caption+key+".Win2")
	}
}
