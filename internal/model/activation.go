package model

import "math"

func Softmax(x []float64) []float64 {
	max := x[0]
	for _, v := range x {
		if v > max {
			max = v
		}
	}

	sum := 0.0
	result := make([]float64, len(x))
	for i, v := range x {
		result[i] = math.Exp(v - max)
		sum += result[i]
	}

	for i := range result {
		result[i] /= sum
	}
	return result
}
