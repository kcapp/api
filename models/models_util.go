package models

import (
	"strconv"
	"strings"
)

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
			if visit.FirstDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.FirstDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.FirstDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.SecondDart.Value.Valid {
			hit := hitsMap[visit.SecondDart.Value.Int64]
			if visit.SecondDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.SecondDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.SecondDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.ThirdDart.Value.Valid {
			hit := hitsMap[visit.ThirdDart.Value.Int64]
			if visit.ThirdDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.ThirdDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.ThirdDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
	}
	return hitsMap, dartsThrown
}

// StringToIntArray will convert the given comma separated string into a int array
func stringToIntArray(s string) []int {
	strs := strings.Split(s, ",")
	ints := make([]int, len(strs))
	for i, v := range strs {
		ints[i], _ = strconv.Atoi(v)
	}
	return ints
}
