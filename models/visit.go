package models

import (
	"errors"
	"sort"
	"strings"

	"github.com/guregu/null"
)

// Visit struct used for storing legs
type Visit struct {
	ID          int         `json:"id"`
	LegID       int         `json:"leg_id"`
	PlayerID    int         `json:"player_id"`
	FirstDart   *Dart       `json:"first_dart"`
	SecondDart  *Dart       `json:"second_dart"`
	ThirdDart   *Dart       `json:"third_dart"`
	IsBust      bool        `json:"is_bust"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Count       int         `json:"count,omitempty"`
	DartsThrown int         `json:"darts_thrown,omitempty"`
	Score       int         `json:"score"`
	Scores      map[int]int `json:"scores"`
}

type comparingMatrix [][]bool

// GetDarts returns all darts for the given visit
func (visit Visit) GetDarts() []Dart {
	darts := []Dart{*visit.FirstDart, *visit.SecondDart, *visit.ThirdDart}
	return darts
}

// GetLastDart will return the last non-miss dart from the visit
func (visit Visit) GetLastDart() *Dart {
	if visit.ThirdDart.IsMiss() {
		if visit.SecondDart.IsMiss() {
			return visit.FirstDart
		}
		return visit.SecondDart
	}
	return visit.ThirdDart
}

// ValidateInput will verify the input does not containg any errors
func (visit Visit) ValidateInput() error {
	if visit.FirstDart == nil {
		return errors.New("First dart cannot be null")
	}
	err := visit.FirstDart.ValidateInput()
	if err != nil {
		return err
	}
	err = visit.SecondDart.ValidateInput()
	if err != nil {
		return err
	}
	err = visit.ThirdDart.ValidateInput()
	if err != nil {
		return err
	}
	return nil
}

// SetIsBust will set IsBust for the given visit
func (visit *Visit) SetIsBust(currentScore int) {
	isBust := false
	isBust = visit.FirstDart.IsBust(currentScore)
	currentScore = currentScore - visit.FirstDart.GetScore()
	if !isBust && currentScore > 0 {
		isBust = visit.SecondDart.IsBust(currentScore)
		currentScore = currentScore - visit.SecondDart.GetScore()
		if !isBust && currentScore > 0 {
			isBust = visit.ThirdDart.IsBust(currentScore)
		} else {
			// Invalidate third dart if second was bust
			visit.ThirdDart.Value = null.IntFromPtr(nil)
		}
	} else {
		// Invalidate second/third dart if first was bust
		visit.SecondDart.Value = null.IntFromPtr(nil)
		visit.ThirdDart.Value = null.IntFromPtr(nil)
	}

	if !isBust && currentScore > 0 {
		// If this visit was not a bust, make sure that darts are set
		// as 0 (miss) instead of 'nil' (not thrown)
		if !visit.FirstDart.Value.Valid {
			visit.FirstDart.Value = null.IntFrom(0)
		}
		if !visit.SecondDart.Value.Valid {
			visit.SecondDart.Value = null.IntFrom(0)
		}
		if !visit.ThirdDart.Value.Valid {
			visit.ThirdDart.Value = null.IntFrom(0)
		}
	}

	visit.IsBust = isBust
}

// IsCheckout will check if the given visit is a checkout (remaining score is 0 and last dart thrown is a double)
func (visit Visit) IsCheckout(currentScore int) bool {
	remaining := currentScore - visit.GetScore()
	if remaining == 0 {
		if visit.ThirdDart.Value.Valid {
			return visit.ThirdDart.IsDouble()
		} else if visit.SecondDart.Value.Valid {
			return visit.SecondDart.IsDouble()
		} else {
			return visit.FirstDart.IsDouble()
		}
	}
	return false
}

// IsViliusVisit will check if this visit was a "Vilius Visit" (Two 20s and a Miss)
func (visit Visit) IsViliusVisit() bool {
	viliusVisit := new(Visit)
	viliusVisit.FirstDart = NewDart(null.IntFrom(20), SINGLE)
	viliusVisit.SecondDart = NewDart(null.IntFrom(20), SINGLE)
	viliusVisit.ThirdDart = NewDart(null.IntFrom(0), SINGLE)

	return visit.isEqualTo(*viliusVisit)
}

// IsFishAndChips will check if this visit was a Fish and Chips (20,5,1)
func (visit Visit) IsFishAndChips() bool {
	fishAndChipsVisit := new(Visit)
	fishAndChipsVisit.FirstDart = NewDart(null.IntFrom(20), SINGLE)
	fishAndChipsVisit.SecondDart = NewDart(null.IntFrom(5), SINGLE)
	fishAndChipsVisit.ThirdDart = NewDart(null.IntFrom(1), SINGLE)

	return visit.isEqualTo(*fishAndChipsVisit)
}

// checkIfEquivalent sees if the input visit is the same as this visit
func (visit Visit) isEqualTo(comparingVisit Visit) bool {
	return visit.GetScore() == comparingVisit.GetScore() && visit.makeComparingMatrix(comparingVisit).isMatrixEqual()
}

// GetVisitString will return a (sorted) string based on the darts thrown. This will make sure common visits will be the same
func (visit Visit) GetVisitString() string {
	strs := []string{visit.FirstDart.GetString(), visit.SecondDart.GetString(), visit.ThirdDart.GetString()}
	sort.Strings(strs)
	return strings.Join(strs, " ")
}

// GetHitsMap will return a map where key is dart and value is count of single,double,triple hits
func GetHitsMap(visits []*Visit) (map[int64]*Hits, int) {
	hitsMap := make(map[int64]*Hits)
	// Populate the map with hits for each value (miss, 1-20, bull)
	for i := 0; i <= 20; i++ {
		hitsMap[int64(i)] = new(Hits)
	}
	hitsMap[25] = new(Hits)

	var dartsThrown int
	for _, visit := range visits {
		if visit.IsBust {
			continue
		}
		if visit.FirstDart.Value.Valid {
			hit := hitsMap[visit.FirstDart.Value.Int64]
			if visit.FirstDart.IsSingle() {
				hit.Singles++
			}
			if visit.FirstDart.IsDouble() {
				hit.Doubles++
			}
			if visit.FirstDart.IsTriple() {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.SecondDart.Value.Valid {
			hit := hitsMap[visit.SecondDart.Value.Int64]
			if visit.SecondDart.IsSingle() {
				hit.Singles++
			}
			if visit.SecondDart.IsDouble() {
				hit.Doubles++
			}
			if visit.SecondDart.IsTriple() {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.ThirdDart.Value.Valid {
			hit := hitsMap[visit.ThirdDart.Value.Int64]
			if visit.ThirdDart.IsSingle() {
				hit.Singles++
			}
			if visit.ThirdDart.IsDouble() {
				hit.Doubles++
			}
			if visit.ThirdDart.IsTriple() {
				hit.Triples++
			}
			dartsThrown++
		}
	}
	return hitsMap, dartsThrown
}

// GetScore will return the total points scored during the given visit
func (visit Visit) GetScore() int {
	return visit.FirstDart.GetScore() + visit.SecondDart.GetScore() + visit.ThirdDart.GetScore()
}

// GetDartsThrown will return the actual number of darts thrown during this visit
func (visit Visit) GetDartsThrown() int {
	thrown := 1
	if visit.SecondDart.Value.Valid {
		thrown++
	}
	if visit.ThirdDart.Value.Valid {
		thrown++
	}
	return thrown
}

// makeComparingMatrix will create a comparing matrix for the two visits
func (visit Visit) makeComparingMatrix(comparingVisit Visit) comparingMatrix {
	comparingMatrix := make([][]bool, 3)
	for i, visitDart := range visit.GetDarts() {
		comparingMatrix[i] = make([]bool, 3)
		for j, comparingDart := range comparingVisit.GetDarts() {
			comparingMatrix[i][j] = visitDart == comparingDart
		}
	}
	return comparingMatrix
}

// isMatrixEqual will check if the values in the matrix are equal
func (matrix comparingMatrix) isMatrixEqual() bool {
	rows := make([]int, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matrix[i][j] {
				rows[i]++
			}
		}
		if rows[i] == 0 {
			return false
		}
	}
	columns := make([]int, 3)
	for cIndex := 0; cIndex < 3; cIndex++ {
		for rIndex := 0; rIndex < 3; rIndex++ {
			if matrix[rIndex][cIndex] {
				columns[cIndex]++
			}
		}
		if columns[cIndex] == 0 {
			return false
		}
	}
	return true
}

// GetMarksHit will return the number of marks for the given visit
// It will only match against the slice of darts, and when other players have not closed it
func (visit Visit) GetMarksHit(darts []int, hitsMap map[int]map[int]int64) int {
	pid := visit.PlayerID
	hits := hitsMap[pid]
	marks := int64(0)

	open, self := isMarkOpen(pid, visit.FirstDart, darts, hitsMap)
	if open || self {
		marks += visit.FirstDart.GetMarksHit(hits, open)
	}
	open, self = isMarkOpen(pid, visit.SecondDart, darts, hitsMap)
	if open || self {
		marks += visit.SecondDart.GetMarksHit(hits, open)
	}
	open, self = isMarkOpen(pid, visit.ThirdDart, darts, hitsMap)
	if open || self {
		marks += visit.ThirdDart.GetMarksHit(hits, open)
	}
	return int(marks)
}

// find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// isMarkOpen will check if the given value is still open for other players and current player
func isMarkOpen(playerID int, dart *Dart, darts []int, hitsMap map[int]map[int]int64) (bool, bool) {
	val := dart.ValueRaw()

	_, found := find(darts, val)
	if found {
		// Check if number is cloed by us
		self := hitsMap[playerID][val] < 3
		others := false
		// Check if number is closed by other players
		for id, playerHits := range hitsMap {
			if playerID != id {
				if playerHits[val] < 3 {
					others = true
				}
			}
		}
		return others, self
	}
	return false, false
}

// CalculateCricketScore will calculate the score for each player for the given visit
func (visit *Visit) CalculateCricketScore(scores map[int]*Player2Leg) int {
	points := visit.FirstDart.CalculateCricketScore(visit.PlayerID, scores)
	points += visit.SecondDart.CalculateCricketScore(visit.PlayerID, scores)
	points += visit.ThirdDart.CalculateCricketScore(visit.PlayerID, scores)
	return points
}

// CalculateAroundTheClockScore will calculate the score for the given visit
func (visit *Visit) CalculateAroundTheClockScore(currentScore int) int {
	score := 0
	if visit.FirstDart.ValueRaw() == currentScore+1 && visit.FirstDart.IsSingle() || (currentScore+1 == 21 && visit.FirstDart.IsBull()) {
		score++
		currentScore++
	}
	if visit.SecondDart.ValueRaw() == currentScore+1 && visit.SecondDart.IsSingle() || (currentScore+1 == 21 && visit.SecondDart.IsBull()) {
		score++
		currentScore++
	}
	if visit.ThirdDart.ValueRaw() == currentScore+1 && visit.ThirdDart.IsSingle() || (currentScore+1 == 21 && visit.ThirdDart.IsBull()) {
		score++
		currentScore++
	}
	return score
}

// CalculateAroundTheWorldScore will calculate the score for the given visit
func (visit *Visit) CalculateAroundTheWorldScore(round int) int {
	score := 0
	if round == visit.FirstDart.ValueRaw() || (round == 21 && visit.FirstDart.IsBull()) {
		score += visit.FirstDart.GetScore()
	}
	if round == visit.SecondDart.ValueRaw() || (round == 21 && visit.SecondDart.IsBull()) {
		score += visit.SecondDart.GetScore()
	}
	if round == visit.ThirdDart.ValueRaw() || (round == 21 && visit.ThirdDart.IsBull()) {
		score += visit.ThirdDart.GetScore()
	}
	return score
}

// CalculateBermudaTriangleScore will calculate the score for the given visit
func (visit *Visit) CalculateBermudaTriangleScore(round int) int {
	score := 0

	target := TargetsBermudaTriangle[round]
	score += visit.FirstDart.GetBermudaTriangleScore(target)
	score += visit.SecondDart.GetBermudaTriangleScore(target)
	score += visit.ThirdDart.GetBermudaTriangleScore(target)
	return score
}

// Calculate420Score will calculate the score for the given visit
func (visit *Visit) Calculate420Score(round int) int {
	score := 0

	target := Targets420[round]
	score += visit.FirstDart.Get420Score(target)
	score += visit.SecondDart.Get420Score(target)
	score += visit.ThirdDart.Get420Score(target)
	return score
}

// IsShanghai will check if the given visit is a "Shanghai". A Shanghai visit is one where a single, double and triple multipler is hit with each dart
func (visit *Visit) IsShanghai() bool {
	first := visit.FirstDart
	second := visit.SecondDart
	third := visit.ThirdDart

	if first.ValueRaw() == second.ValueRaw() && first.ValueRaw() == third.ValueRaw() && first.ValueRaw() != 0 &&
		((first.Multiplier == 1 && second.Multiplier == 2 && third.Multiplier == 3) ||
			(first.Multiplier == 2 && second.Multiplier == 3 && third.Multiplier == 1) ||
			(first.Multiplier == 3 && second.Multiplier == 1 && third.Multiplier == 2) ||
			(first.Multiplier == 3 && second.Multiplier == 2 && third.Multiplier == 1) ||
			(first.Multiplier == 1 && second.Multiplier == 3 && third.Multiplier == 2) ||
			(first.Multiplier == 2 && second.Multiplier == 1 && third.Multiplier == 3)) {
		return true
	}
	return false
}
