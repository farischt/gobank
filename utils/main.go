package utils

import "strconv"

func Uint8ToFloat(u []uint8) float64 {
	b, _ := strconv.ParseFloat(string(u), 64)
	return b
}
