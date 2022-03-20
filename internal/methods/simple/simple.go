package simple

import (
	"log"
	"math"
)

type Simple struct {
	max uint64
}

func New(n int) *Simple {
	log.Printf("prepared simple method with n: %v", n)
	return &Simple{
		max: uint64(math.Pow(2, float64(n))) - 1,
	}
}

// F is a straightforward function that when received uint in a range [0, max) and will return another uint in the same
// range. Since this is going to be used on client side to generate a tree, we need it to allow F(x) = F(x') so that
// some nodes have multiple leafs requiring client to walk through all of them and access memory more often.
func (s Simple) F(x uint64) uint64 {
	result := math.Sin(float64(x))
	if result < 0 {
		result = -result
	}

	result = result * float64(s.max)

	// Do not allow F(x) = x
	if uint64(result) == x {
		result += float64(s.max) / 2
		result = float64(uint64(result) % s.max)
	}

	return uint64(result)
}