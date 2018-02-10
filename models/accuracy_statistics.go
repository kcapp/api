package models

import "github.com/guregu/null"

// AccuracyStatistics struct used for storing accuracy statistics
type AccuracyStatistics struct {
	AccuracyOverall null.Float
	Accuracy20      null.Float
	attempts20      int
	hits20          int
	Accuracy19      null.Float
	attempts19      int
	hits19          int
	misses          int
}

// GetAccuracyStats will add statistics for the given dart
func (stats *AccuracyStatistics) GetAccuracyStats(remainingScore int, dart *Dart) {
	if remainingScore-dart.GetScore() < 171 {
		// We only want to calculate accuracy stats when player has a remaining score over 170
		return
	}

	score := dart.Value.Int64
	switch score {
	case 20:
		stats.hits20++
		stats.attempts20++
		stats.Accuracy20.Float64 += 100
	case 5, 1:
		stats.attempts20++
		stats.Accuracy20.Float64 += 70
	case 12, 18:
		stats.attempts20++
		stats.Accuracy20.Float64 += 30
	case 9, 4:
		stats.attempts20++
		stats.Accuracy20.Float64 += 5
	case 19:
		stats.hits19++
		stats.attempts19++
		stats.Accuracy19.Float64 += 100
	case 7, 3:
		stats.attempts19++
		stats.Accuracy19.Float64 += 70
	case 16, 17:
		stats.attempts19++
		stats.Accuracy19.Float64 += 30
	case 8, 2:
		stats.attempts19++
		stats.Accuracy19.Float64 += 5
	default:
		stats.misses++
	}
}

// SetAccuracy will set the accuracy based on hits and attempts
func (stats *AccuracyStatistics) SetAccuracy() {
	stats.AccuracyOverall = null.FloatFrom((stats.Accuracy20.Float64 + stats.Accuracy19.Float64) /
		float64((stats.attempts20 + stats.attempts19 + stats.misses)))

	if stats.attempts20 == 0 {
		stats.Accuracy20 = null.NewFloat(0, false)
	} else {
		stats.Accuracy20 = null.FloatFrom(stats.Accuracy20.Float64 / float64(stats.attempts20))
	}

	if stats.attempts19 == 0 {
		stats.Accuracy19 = null.NewFloat(0, false)
	} else {
		stats.Accuracy19 = null.FloatFrom(stats.Accuracy19.Float64 / float64(stats.attempts19))
	}
}
