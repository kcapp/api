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
