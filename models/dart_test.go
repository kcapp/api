package models

import (
	"errors"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

// TestIsBust will check that the given dart is bust
func TestIsBust(t *testing.T) {
	dart := Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsBust(20, OUTSHOTDOUBLE), true, "should be bust")

	dart = Dart{Value: null.IntFrom(10), Multiplier: 2}
	assert.Equal(t, dart.IsBust(20, OUTSHOTDOUBLE), false, "should not be bust")

	dart = Dart{Value: null.NewInt(-1, false), Multiplier: 1}
	assert.Equal(t, dart.IsBust(301, OUTSHOTDOUBLE), false, "should be bust")
	assert.Equal(t, dart.Value.Valid, true, "should be valid")
	assert.Equal(t, dart.Value.Int64, int64(0), "should be 0")

	dart = Dart{Value: null.IntFrom(3), Multiplier: 1}
	assert.Equal(t, dart.IsBust(4, OUTSHOTDOUBLE), true, "should be bust")

	dart = Dart{Value: null.IntFrom(3), Multiplier: 1}
	assert.Equal(t, dart.IsBust(4, OUTSHOTANY), false, "should not be bust")
}

func TestIsBust_ScoreIsZero(t *testing.T) {
	dart := Dart{Value: null.IntFrom(3), Multiplier: 1}
	assert.Equal(t, dart.IsBust(3, OUTSHOTANY), false, "should not be bust")

	dart = Dart{Value: null.IntFrom(3), Multiplier: 2}
	assert.Equal(t, dart.IsBust(6, OUTSHOTDOUBLE), false, "should not be bust")

	dart = Dart{Value: null.IntFrom(3), Multiplier: 3}
	assert.Equal(t, dart.IsBust(9, OUTSHOTMASTER), false, "should not be bust")

	dart = Dart{Value: null.IntFrom(3), Multiplier: 2}
	assert.Equal(t, dart.IsBust(6, OUTSHOTMASTER), false, "should not be bust")
}

// TestIsBustAbove will check that the given dart is bust above
func TestIsBustAbove(t *testing.T) {
	dart := Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsBustAbove(181, 200), true, "should be bust")

	dart = Dart{Value: null.NewInt(-1, false), Multiplier: 1}
	assert.Equal(t, dart.IsBustAbove(180, 200), false, "should be bust")
	assert.Equal(t, dart.Value.Valid, true, "should be valid")
	assert.Equal(t, dart.Value.Int64, int64(0), "should be 0")
}

// TestValidateInput will check that the input is valid
func TestValidateInput(t *testing.T) {
	// Invalid value
	dart := Dart{Value: null.IntFrom(-1), Multiplier: 1}
	err := dart.ValidateInput()
	assert.Equal(t, err, errors.New("value cannot be less than 0"), "should be equal")

	// Invalid value
	dart = Dart{Value: null.IntFrom(26), Multiplier: 1}
	err = dart.ValidateInput()
	assert.Equal(t, err, errors.New("value has to be less than 21 (or 25 (bull))"), "should be equal")

	// Invalid multiplier
	dart = Dart{Value: null.IntFrom(20), Multiplier: 4}
	err = dart.ValidateInput()
	assert.Equal(t, err, errors.New("multiplier has to be one of 1 (single), 2 (douhle), 3 (triple)"), "should be equal")

	// Invalid multiplier
	dart = Dart{Value: null.IntFrom(20), Multiplier: -1}
	err = dart.ValidateInput()
	assert.Equal(t, err, errors.New("multiplier has to be one of 1 (single), 2 (douhle), 3 (triple)"), "should be equal")

	// Mulitplier changed
	dart = Dart{Value: null.IntFrom(0), Multiplier: 3}
	err = dart.ValidateInput()
	assert.Equal(t, err == nil, true, "should be nil")
	assert.Equal(t, dart.Multiplier, int64(1), "should be equal")
}

// TestGetScore will check that score is correct
func TestGetScore(t *testing.T) {
	dart := Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.GetScore(), 60, "score should be 60")
}

// TestGetBermudaTriangleScore will check that BermudaTriangle score is correctly calculated
func TestGetBermudaTriangleScore(t *testing.T) {
	// Any 12
	dart := Dart{Value: null.IntFrom(12), Multiplier: 3}
	score := dart.GetBermudaTriangleScore(TargetsBermudaTriangle[0])
	assert.Equal(t, score, 36, "score should be 36")

	// Any Double
	dart = Dart{Value: null.IntFrom(16), Multiplier: 2}
	score = dart.GetBermudaTriangleScore(TargetsBermudaTriangle[3])
	assert.Equal(t, score, 32, "score should be 32")

	// Single Bull
	dart = Dart{Value: null.IntFrom(25), Multiplier: 2}
	score = dart.GetBermudaTriangleScore(TargetsBermudaTriangle[11])
	assert.Equal(t, score, 25, "score should be 25")

	// Miss
	dart = Dart{Value: null.IntFrom(1), Multiplier: 1}
	score = dart.GetBermudaTriangleScore(TargetsBermudaTriangle[5])
	assert.Equal(t, score, 0, "score should be 0")
}

// TestGet420Score will check that 420 score is correctly calculated
func TestGet420Score(t *testing.T) {
	// Double 1
	dart := Dart{Value: null.IntFrom(1), Multiplier: 2}
	score := dart.Get420Score(Targets420[0])
	assert.Equal(t, score, 2, "score should be 2")

	// Miss Single
	dart = Dart{Value: null.IntFrom(1), Multiplier: 1}
	score = dart.Get420Score(Targets420[0])
	assert.Equal(t, score, 0, "score should be 0")

	// Miss Triple
	dart = Dart{Value: null.IntFrom(1), Multiplier: 3}
	score = dart.Get420Score(Targets420[0])
	assert.Equal(t, score, 0, "score should be 0")
}

// TestGetJDCPracticeScore will check that JDC Practice score is correctly calculated
func TestGetJDCPracticeScore(t *testing.T) {
	// Single 10
	dart := Dart{Value: null.IntFrom(10), Multiplier: 1}
	score := dart.GetJDCPracticeScore(TargetsJDCPractice[0])
	assert.Equal(t, score, 10, "score should be 10")

	// Double 10
	dart = Dart{Value: null.IntFrom(10), Multiplier: 2}
	score = dart.GetJDCPracticeScore(TargetsJDCPractice[0])
	assert.Equal(t, score, 20, "score should be 20")

	// Triple 10
	dart = Dart{Value: null.IntFrom(10), Multiplier: 3}
	score = dart.GetJDCPracticeScore(TargetsJDCPractice[0])
	assert.Equal(t, score, 30, "score should be 30")

	// Miss
	dart = Dart{Value: null.IntFrom(1), Multiplier: 1}
	score = dart.GetJDCPracticeScore(TargetsJDCPractice[0])
	assert.Equal(t, score, 0, "score should be 0")
}

// TestIsCheckoutAttempt_Doubles will check if the given dart is a checkout attempt
func TestIsCheckoutAttempt_Doubles(t *testing.T) {
	// Invalid dart
	dart := Dart{Value: null.NewInt(0, false), Multiplier: 2}
	assert.Equal(t, dart.IsCheckoutAttempt(301, 1, OUTSHOTDOUBLE), false, "should be false")

	// Not checkout
	dart = Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsCheckoutAttempt(301, 1, OUTSHOTDOUBLE), false, "should be false")

	// Successful checkout
	dart = Dart{Value: null.IntFrom(20), Multiplier: 2}
	assert.Equal(t, dart.IsCheckoutAttempt(40, 1, OUTSHOTDOUBLE), true, "should be true")

	// Checkout attempt
	dart = Dart{Value: null.IntFrom(8), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(32, 1, OUTSHOTDOUBLE), true, "should be true")

	// Checkout attempt bull
	dart = Dart{Value: null.IntFrom(18), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(50, 3, OUTSHOTDOUBLE), true, "should be true")
}

// TestIsCheckoutAttempt_Any will check if the given dart is a checkout attempt
func TestIsCheckoutAttempt_Any(t *testing.T) {
	dart := Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsCheckoutAttempt(20, 1, OUTSHOTANY), true, "should be true")

	dart = Dart{Value: null.IntFrom(20), Multiplier: 2}
	assert.Equal(t, dart.IsCheckoutAttempt(40, 1, OUTSHOTANY), true, "should be true")

	dart = Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(60, 1, OUTSHOTANY), true, "should be true")
}

// TestIsCheckoutAttempt_Master will check if the given dart is a checkout attempt
func TestIsCheckoutAttempt_Master(t *testing.T) {
	dart := Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(23, 1, OUTSHOTMASTER), false, "should not be true")

	dart = Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(40, 1, OUTSHOTMASTER), true, "should be true")

	dart = Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsCheckoutAttempt(60, 1, OUTSHOTMASTER), true, "should be true")
}

// TestGetString will check that dart string is created correctly
func TestGetString(t *testing.T) {
	// 3-20
	dart := Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.GetString(), "3-20", "string should be 3-20")

	// 1-NULL
	dart = Dart{Value: null.NewInt(-1, false), Multiplier: 1}
	assert.Equal(t, dart.GetString(), "1-NULL", "string should be 1-NULL")
}

// TestNewDart will check that a new dart is created correctly
func TestNewDart(t *testing.T) {
	assert.Equal(t, NewDart(null.IntFrom(20), 2), &Dart{Value: null.IntFrom(20), Multiplier: 2}, "darts should be equal")
}

// TestIsSingle will check that the dart is single
func TestIsSingle(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsSingle(), true, "dart should be single")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 2}
	assert.Equal(t, dart.IsSingle(), false, "dart should not be single")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsSingle(), false, "dart should not be single")
}

// TestIsDouble will check that the dart is double
func TestIsDouble(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsDouble(), false, "dart should not be double")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 2}
	assert.Equal(t, dart.IsDouble(), true, "dart should be double")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsDouble(), false, "dart should not be double")
}

// TestIsTriple will check that the dart is triple
func TestIsTriple(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsTriple(), false, "dart should not be triple")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 2}
	assert.Equal(t, dart.IsTriple(), false, "dart should not be triple")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsTriple(), true, "dart should be triple")
}

// TestIsBull will check that the dart is bull
func TestIsBull(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(20), Multiplier: 1}
	assert.Equal(t, dart.IsBull(), false, "dart should not be bull")

	dart = &Dart{Value: null.IntFrom(25), Multiplier: 1}
	assert.Equal(t, dart.IsBull(), true, "dart should be bull")

	dart = &Dart{Value: null.IntFrom(25), Multiplier: 2}
	assert.Equal(t, dart.IsBull(), true, "dart should be bull")
}

// TestIsMiss will check that the dart is miss
func TestIsMiss(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(0), Multiplier: 1}
	assert.Equal(t, dart.IsMiss(), true, "dart should be miss")

	dart = &Dart{Value: null.IntFrom(19), Multiplier: 2}
	assert.Equal(t, dart.IsMiss(), false, "dart should not be miss")
}

// TestIsCricketMiss will check that the dart is a cricket miss
func TestIsCricketMiss(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(1), Multiplier: 1}
	assert.Equal(t, dart.IsCricketMiss(), true, "dart should be miss")

	for _, num := range []int{15, 16, 17, 18, 19, 20, 25} {
		dart = &Dart{Value: null.IntFrom(int64(num)), Multiplier: 1}
		assert.Equal(t, dart.IsCricketMiss(), false, "dart should not be miss")
	}
}

// TestValueRaw will check that the dart returns correct value
func TestValueRaw(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(17), Multiplier: 1}
	assert.Equal(t, dart.ValueRaw(), 17, "dart should be 17")

	dart = &Dart{Value: null.NewInt(0, false), Multiplier: 1}
	assert.Equal(t, dart.ValueRaw(), 0, "dart should be 0")
}

// TestGetMarksHit will check the number of darts hit
func TestGetMarksHit(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(17), Multiplier: 3}

	// Number open
	marks := dart.GetMarksHit(make(map[int]int64), true)
	assert.Equal(t, marks, int64(3), "marks should be 3")

	// Number closed
	hits := make(map[int]int64)
	hits[17] = 3
	marks = dart.GetMarksHit(hits, false)
	assert.Equal(t, marks, int64(0), "marks should be 0")

	// Already hit three times
	hits = make(map[int]int64)
	hits[17] = 4
	marks = dart.GetMarksHit(hits, true)
	assert.Equal(t, marks, int64(3), "marks should be 3")

	// Already hit three times closed
	hits = make(map[int]int64)
	hits[17] = 2
	marks = dart.GetMarksHit(hits, false)
	assert.Equal(t, marks, int64(1), "marks should be 1")
}

// TestCalculateCricketScore will check that the calculated score is correct
func TestCalculateCricketScore(t *testing.T) {
	// Invalid dart
	dart := &Dart{Value: null.NewInt(-1, false), Multiplier: 1}
	score := dart.CalculateCricketScore(1, make(map[int]*Player2Leg))
	assert.Equal(t, score, 0, "score should be 0")

	// Not cricket dart
	dart = &Dart{Value: null.IntFrom(1), Multiplier: 1}
	score = dart.CalculateCricketScore(1, make(map[int]*Player2Leg))
	assert.Equal(t, score, 0, "score should be 0")

	// No hits
	scores := make(map[int]*Player2Leg)
	p2l := new(Player2Leg)
	p2l.Hits = make(map[int]*Hits)
	scores[1] = p2l

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 1}
	score = dart.CalculateCricketScore(1, scores)
	assert.Equal(t, score, 0, "score should be 0")

	// Already closed
	scores = make(map[int]*Player2Leg)
	p2l = new(Player2Leg)
	p2l.Hits = make(map[int]*Hits)

	hits := new(Hits)
	hits.Total = 3

	p2l.Hits[20] = hits
	scores[1] = p2l
	scores[2] = new(Player2Leg)

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	score = dart.CalculateCricketScore(1, scores)
	assert.Equal(t, score, 60, "score should be 60")

	// Closed by all players
	scores = make(map[int]*Player2Leg)
	p2l = new(Player2Leg)
	p2l.Hits = make(map[int]*Hits)

	hits = new(Hits)
	hits.Total = 3
	p2l.Hits[20] = hits

	scores[1] = p2l
	scores[2] = p2l

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	score = dart.CalculateCricketScore(1, scores)
	assert.Equal(t, score, 0, "score should be 0")
}

// TestIsHit will check that the dart hit
func TestIsHit(t *testing.T) {
	dart := &Dart{Value: null.IntFrom(17), Multiplier: 1}
	assert.Equal(t, dart.IsHit([]int{5, 17, 3, 8}), true, "dart should be hit")

	dart = &Dart{Value: null.IntFrom(20), Multiplier: 3}
	assert.Equal(t, dart.IsHit([]int{5, 17, 3, 8}), false, "dart should be miss")
}
