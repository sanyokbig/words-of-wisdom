package solver

import (
	"log"

	"github.com/sanyokbig/words-of-wisdom/internal/checksum"
	mathpkg "github.com/sanyokbig/words-of-wisdom/internal/math"
)

// Solver tried to find x0 by building a tree of all possible sequences that could be used to achieve xk from x0 by
// applying F() K times.
// Assuming that calling an inverse F() costs more that building and accessing an inverted table of F() results,
// client is encouraged to use this map to solve the challenge, thus making it memory-bound.
type Solver struct {
	// xk is a value that was received from server and a base value that is used to find a solution
	xk uint64

	// targetDepth is an amount of F() applications that were used on solution
	targetDepth int

	// checksum is received from server, used to validate that found sequence is correct
	checksum string

	// inversionTable is a map that allows us to easily find a result on an InvertedF()
	// Given both:
	// 1. F(x) = y
	// 2. InvertedF(y) = x
	// We will store invTable[y] = x
	//
	// Since F(x) and F(x') can result in the same y, we should store the value as a slice:
	// inversionTable[y] = [x, x']
	inversionTable map[uint64][]uint64
}

type Method interface {
	F(uint64) uint64
}

func New(xk uint64, n, k int, expectedChecksum string, method Method) *Solver {
	return &Solver{
		targetDepth:    k,
		xk:             xk,
		checksum:       expectedChecksum,
		inversionTable: buildInversionTable(n, method),
	}
}

func (s *Solver) Solve() (uint64, bool) {
	log.Printf("looking for a value at depth %v with a sequnce checksum of %v", s.targetDepth, s.checksum)

	return s.findValueAtDepth(s.xk, []uint64{})
}

// Build a table for InvertedF() by applying F() to all integers in [0, 2^n)
func buildInversionTable(n int, method Method) map[uint64][]uint64 {
	max := mathpkg.Pow2(n)
	inversionTable := make(map[uint64][]uint64, max)

	for i := 0; i < int(max); i++ {
		v := uint64(i)

		fResult := method.F(v)

		row, ok := inversionTable[fResult]
		if !ok {
			row = []uint64{}
		}

		row = append(row, v)

		inversionTable[fResult] = row
	}

	log.Printf("built an inversionTable with a max %v and length %v", max, len(inversionTable))

	return inversionTable
}

// findValueAtDepth will run recursively and scan all leafs until required depth is achieved or there are no more leafs
// If required depth is achieved and sequence checksum matches one received from the server, then we've found a solution
func (s *Solver) findValueAtDepth(currValue uint64, sequence []uint64) (value uint64, ok bool) {
	currDepth := s.targetDepth - len(sequence)
	sequence = append(sequence, currValue)

	// Check for checksum when we've achieved required depth
	// In case sequence does not match, don't try to go deeper.
	if currDepth == 0 {
		gotChecksum := checksum.Make(sequence)

		if gotChecksum == s.checksum {
			return currValue, true
		}

		return 0, false
	}

	// CurrValue needs to by XORed with currDepth to match XOR done on the server
	leafs, ok := s.inversionTable[currValue^uint64(currDepth)]
	if !ok {
		return 0, false
	}

	// Call this func recursively for next leaves until solution is found
	for _, leaf := range leafs {
		value, ok := s.findValueAtDepth(leaf, sequence)
		if ok {
			return value, true
		}
	}

	return 0, false
}
