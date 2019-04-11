package common

import (
	"log"
	"math"
	"strconv"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func truncate(num float64) int {
	return int(num + math.Copysign(0.0, num))
}

func ToFixed(num float64, precision int) float64 {
	// output := math.Pow(10, float64(precision))
	// return float64(round(num*output)) / output
	output := math.Pow(10, float64(precision))
	return float64(truncate(num*output)) / output

}

func TrailingFloatZero(f float64, maxDigit int) string {
	fs := strconv.FormatFloat(f, 'f', maxDigit, 64)
	ff, err := strconv.ParseFloat(fs, 64)
	if err != nil {
		log.Printf(" TrailingFloatZero(f float64,maxDigit int) error: %v ", err)
	}
	fst := strconv.FormatFloat(ff, 'f', -1, 64)
	return fst
}
