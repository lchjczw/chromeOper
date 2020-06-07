package util

import (
	"strconv"
	"strings"
)

func SplitNum(str string, sep string) []int {
	a := strings.Split(str, sep)
	var r []int
	for _, value := range a {
		i, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		r = append(r, i)
	}

	return r
}
