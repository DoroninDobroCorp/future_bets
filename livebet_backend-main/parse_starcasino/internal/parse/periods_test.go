package parse

import (
	"github.com/stretchr/testify/assert"
	"livebets/shared"
	"testing"
)

func TestPeriod0(t *testing.T) {
	var periods []shared.PeriodData

	periods = matchHas(0, periods)

	assert.Equal(t, 1, len(periods), "len(periods)")

	assert.NotNil(t, periods[0].Totals, "periods[0].Totals")
}

func TestPeriod0Win(t *testing.T) {
	periods := make([]shared.PeriodData, 0, 5)

	periods = matchHas(0, periods)

	assert.Equal(t, 1, len(periods), "len(periods)")
	assert.Equal(t, 5, cap(periods), "cap(periods)")

	const roundValue = 1.175
	periods[0].Win1x2.Win1.Value = roundValue

	assert.Equal(t, roundValue, periods[0].Win1x2.Win1.Value, "periods[0].Win1x2.Win1")
}

func TestPeriod1Win(t *testing.T) {
	var periods []shared.PeriodData

	periods = matchHas(0, periods)

	assert.Equal(t, 1, len(periods), "len(periods)")

	const roundValue = 1.175
	periods[0].Win1x2.Win1.Value = roundValue

	periods = matchHas(1, periods)

	assert.Equal(t, 2, len(periods), "len(periods)")

	periods[1].Win1x2.Win2.Value = roundValue

	assert.Equal(t, roundValue, periods[0].Win1x2.Win1.Value, "periods[0].Win1x2.Win1")
	assert.Equal(t, roundValue, periods[1].Win1x2.Win2.Value, "periods[1].Win1x2.Win2")

	assert.NotNil(t, periods[0].Totals, "periods[0].Totals")
	assert.NotNil(t, periods[1].Totals, "periods[1].Totals")
}

func TestPeriod4(t *testing.T) {
	var periods []shared.PeriodData

	periods = matchHas(2, periods)

	assert.Equal(t, 3, len(periods), "len(periods)")

	const roundValue = 1.175
	periods[1].Win1x2.Win1.Value = roundValue

	periods[2].Totals["16.5"] = &shared.WinLessMore{WinLess: shared.Odd{Value: roundValue}}

	periods = matchHas(4, periods)

	assert.Equal(t, 5, len(periods), "len(periods)")

	periods[4].Win1x2.Win2.Value = roundValue

	assert.Equal(t, roundValue, periods[1].Win1x2.Win1.Value, "periods[1].Win1x2.Win1")
	assert.Equal(t, roundValue, periods[4].Win1x2.Win2.Value, "periods[4].Win1x2.Win2")

	assert.Equal(t, roundValue, periods[2].Totals["16.5"].WinLess.Value, "periods[2].Totals[16.5].WinLess")
}

func TestPeriod5(t *testing.T) {
	var periods []shared.PeriodData

	periods = matchHas(2, periods)

	assert.Equal(t, 3, len(periods), "len(periods)")

	periods = matchHas(5, periods)

	const roundValue = 1.175
	periods[5].Win1x2.Win1.Value = roundValue

	periods[5].Totals["16.5"] = &shared.WinLessMore{WinLess: shared.Odd{Value: roundValue}}

	// 5 periods -> 3 periods
	periods = matchHas(3, periods)

	periods[3].Win1x2.Win2.Value = roundValue

	assert.Equal(t, 6, len(periods), "len(periods)")

	assert.Equal(t, roundValue, periods[5].Win1x2.Win1.Value, "periods[5].Win1x2.Win1")
	assert.Equal(t, roundValue, periods[3].Win1x2.Win2.Value, "periods[3].Win1x2.Win2")

	assert.Equal(t, roundValue, periods[5].Totals["16.5"].WinLess.Value, "periods[5].Totals[16.5].WinLess")
}
