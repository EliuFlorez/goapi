package utils

import (
	"math/rand"
	"strconv"
)

func CodeRandom() int {
	min := 100000
	max := 999999
	return rand.Intn(max-min) + min
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
