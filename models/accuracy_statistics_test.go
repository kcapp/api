package models

import (
	"math"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

// TestGetAccuracyStatsLessThan170 checks that statistics is 0 when score is less than 170
func TestGetAccuracyStatsLessThan170(t *testing.T) {
	stats := new(AccuracyStatistics)
	stats.GetAccuracyStats(170, &Dart{Value: null.IntFrom(20), Multiplier: 1})
	assert.Equal(t, stats.attempts20, 0, "should have 0 attempts at 20")
	assert.Equal(t, stats.hits20, 0, "should have 0 hit at 20")
	assert.Equal(t, stats.Accuracy20.Float64, float64(0), "accuracy 20 should be 0")
}

// TestGetAccuracyStats checks that statistics is correctly calculated
func TestGetAccuracyStats(t *testing.T) {
	stats := new(AccuracyStatistics)
	for i := 1; i <= 20; i++ {
		stats.GetAccuracyStats(301, &Dart{Value: null.IntFrom(int64(i)), Multiplier: 1})
	}
	assert.Equal(t, stats.attempts19, 7, "should have 7 attempts at 19")
	assert.Equal(t, stats.hits19, 1, "should have 1 hit at 19")
	assert.Equal(t, stats.Accuracy19.Float64, float64(310), "accuracy 19 should be 310")

	assert.Equal(t, stats.attempts20, 7, "should have 7 attempts at 20")
	assert.Equal(t, stats.hits20, 1, "should have 1 hit at 20")
	assert.Equal(t, stats.Accuracy20.Float64, float64(310), "accuracy 20 should be 310")

	assert.Equal(t, stats.misses, 6, "misses should be 6")
}

// TestSetAccuracyStatsNoAttempts checks that no statistics is set
func TestSetAccuracyStatsNoAttempts(t *testing.T) {
	stats := new(AccuracyStatistics)

	stats.SetAccuracy()
	assert.Equal(t, stats.Accuracy19, null.NewFloat(0, false), "accuracy 19 should be 0")
	assert.Equal(t, stats.Accuracy20, null.NewFloat(0, false), "accuracy 20 should be 0")
	assert.Equal(t, math.IsNaN(stats.AccuracyOverall.Float64), true, "accuracy overall should be nan")
}

// TestSetAccuracyStatsNoAttempts checks that statistics is correctly set
func TestSetAccuracyStats(t *testing.T) {
	stats := new(AccuracyStatistics)
	stats.attempts19 = 7
	stats.Accuracy19 = null.FloatFrom(310)

	stats.attempts20 = 7
	stats.Accuracy20 = null.FloatFrom(310)

	stats.misses = 6

	stats.SetAccuracy()
	assert.Equal(t, stats.Accuracy19, null.FloatFrom(44.285714285714285), "accuracy 19 should be 44%")
	assert.Equal(t, stats.Accuracy20, null.FloatFrom(44.285714285714285), "accuracy 20 should be 44%")
	assert.Equal(t, stats.AccuracyOverall, null.FloatFrom(31), "accuracy overall should be 31%")
}
