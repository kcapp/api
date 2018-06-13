package models

import (
	"errors"
	"sort"
	"strings"

	"github.com/guregu/null"
	"github.com/kcapp/api/util"
)

// Visit struct used for storing legs
type Visit struct {
	ID          int    `json:"id"`
	LegID       int    `json:"leg_id"`
	PlayerID    int    `json:"player_id"`
	FirstDart   *Dart  `json:"first_dart"`
	SecondDart  *Dart  `json:"second_dart"`
	ThirdDart   *Dart  `json:"third_dart"`
	IsBust      bool   `json:"is_bust"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Count       int    `json:"count,omitempty"`
	DartsThrown int    `json:"darts_thrown,omitempty"`
}

//MakeListOfDarts does what the title says for a visit
func (visit Visit) MakeListOfDarts() []Dart {
	darts := []Dart{*visit.FirstDart, *visit.SecondDart, *visit.ThirdDart}
	return darts
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

// IsViliusVisit will check if this visit was a "Vilius Visit" (Two 20s and a Miss)
func (visit Visit) IsViliusVisit() bool {
	if visit.FirstDart.Multiplier != 1 || visit.SecondDart.Multiplier != 1 || visit.ThirdDart.Multiplier != 1 {
		return false
	}
	if (visit.FirstDart.Value.Int64 == 20 && visit.SecondDart.Value.Int64 == 0 && visit.ThirdDart.Value.Int64 == 20) ||
		(visit.FirstDart.Value.Int64 == 0 && visit.SecondDart.Value.Int64 == 20 && visit.ThirdDart.Value.Int64 == 20) ||
		(visit.FirstDart.Value.Int64 == 20 && visit.SecondDart.Value.Int64 == 20 && visit.ThirdDart.Value.Int64 == 0) {
		return true
	}
	return false
}

// IsFishAndChips will check if this visit was a Fish and Chips (20,5,1)
func (visit Visit) IsFishAndChips() bool {
	fishAndChipsVisit := new(Visit)
	fishAndChipsVisit.FirstDart = MakeDart(null.NewInt(20, true), SINGLE)
	fishAndChipsVisit.SecondDart = MakeDart(null.NewInt(5, true), SINGLE)
	fishAndChipsVisit.ThirdDart = MakeDart(null.NewInt(1, true), SINGLE)

	return visit.checkIfEquivalent(*fishAndChipsVisit)
}

// checkIfEquivalent sees if the input visit is the same as this visit
func (visit Visit) checkIfEquivalent(comparingVisit Visit) bool {
	comparingMatrix := makeComparingMatrix(visit, comparingVisit)
	rowSums := getSumsforRows(comparingMatrix)
	columnSums := getSumsforColumns(comparingMatrix)
	sort.Ints(rowSums)
	sort.Ints(columnSums)
	return util.Equal(rowSums, columnSums)
}

func makeComparingMatrix(visit Visit, comparingVisit Visit) [][]bool {
	comparingMatrix := make([][]bool, 3)
	for i, visitDart := range visit.MakeListOfDarts() {
		comparingMatrix[i] = make([]bool, 3)
		for j, comparingDart := range comparingVisit.MakeListOfDarts() {
			comparingMatrix[i][j] = visitDart == comparingDart
		}
	}
	return comparingMatrix
}

func getSumsforRows(matrix [][]bool) []int {
	retValue := make([]int, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matrix[i][j] {
				retValue[i]++
			}
		}
	}
	return retValue
}
func getSumsforColumns(matrix [][]bool) []int {
	retValue := make([]int, 3)
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			if matrix[i][j] {
				retValue[j]++
			}
		}
	}
	return retValue
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
