package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mathpkg "github.com/sanyokbig/words-of-wisdom/internal/math"
)

func TestValuesInValidRange(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "n is 4",
			n:    4,
		},
		{
			name: "n is 20",
			n:    20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simpleMethod := New(tt.n)

			max := mathpkg.Pow2(tt.n)

			for i := uint64(0); i < max; i++ {
				got := simpleMethod.F(i)

				assert.GreaterOrEqual(t, got, uint64(0))
				assert.Less(t, got, max)
			}
		})
	}
}

func TestHasCollidingResults(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "n is 4",
			n:    4,
		},
		{
			name: "n is 20",
			n:    20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simpleMethod := New(tt.n)

			max := mathpkg.Pow2(tt.n)

			results := map[uint64]int{}

			for i := uint64(0); i < max; i++ {
				got := simpleMethod.F(i)

				results[got]++
			}

			multiple := 0
			for _, r := range results {
				if r > 1 {
					multiple++
				}
			}

			t.Logf("results: %v, multple: %v", len(results), multiple)

			assert.Greater(t, multiple, 0)
		})
	}
}
