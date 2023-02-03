package util

import (
	"net/http"
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

// SetHeaders will set the default headers used by all requests
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func SliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}
