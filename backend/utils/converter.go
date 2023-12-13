package utils

import "math"

// RoundMoney Copy from https://stackoverflow.com/a/29786394
func RoundMoney(amount float64) float64 {
	return toFixed(amount, 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
