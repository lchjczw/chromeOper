package chromeOper

import (
	"strconv"
	"time"
)

const (
	TIME_TYPE = time.Millisecond
)

func GetTime(t int) time.Duration {
	return time.Duration(t) * TIME_TYPE
}

func TransTime(t string) time.Duration {
	i, err := strconv.Atoi(t)
	if err != nil {
		i = 0
	}
	return GetTime(i)
}
