package math

import "math"

func Pow2(n int) uint64 {
	return uint64(math.Pow(2, float64(n)))
}
