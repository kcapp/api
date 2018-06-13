package util

import (
	"strconv"
	"strings"
)

// StringToIntArray will convert the given comma separated string into a int array
func StringToIntArray(s string) []int {
	strs := strings.Split(s, ",")
	ints := make([]int, len(strs))
	for i, v := range strs {
		ints[i], _ = strconv.Atoi(v)
	}
	return ints
}

// Equal compares two int slices and returns true if they are equal
func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
