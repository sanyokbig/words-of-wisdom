package challenger

import (
	"log"
	"math"
	"math/rand"

	"github.com/sanyokbig/word-of-wisdom/internal/checksum"
)

// Challenge contains all info required to solve the challenge
// Inspired by https://users.soe.ucsc.edu/~abadi/Papers/memory-longer-acm.pdf
type Challenge struct {
	// X0 is a number that the client needs to find to solve the challenge
	// X0 is in range [0, 2^BitSize)
	X0 uint64

	// Xk is a number that was generated after applying F() K times
	Xk uint64

	// K is a number of F() applications, can also be called a depth
	K int

	// N is a bit size of max value of X0
	N int

	// Checksum is a hash to sequence Xk,...,X0.
	// Used to find a correct solution
	Checksum string
}

// Method will be applied to X0 to generate Xk
type Method interface {
	F(uint64) uint64
}

// Challenger is a constructor of Challenge using Prepare() method
type Challenger struct{}

func New() *Challenger {
	return &Challenger{}
}

// Prepare will prepare a Challenge with selected Method and difficulty
func (c *Challenger) Prepare(method Method, n, k int) *Challenge {
	// Max possible value for x0
	max := uint64(math.Pow(2, float64(n))) - 1

	// Initial value, which the client will try to find
	// Remainder from division is used to limit by max
	x0 := rand.Uint64() % (max + 1)

	// xk will be the value that client will receive, to be modified later
	xk := x0

	// Sequence represents the sequence xk,...,x0 as a slice of bytes
	// Will be used later to generate sequence checksum for client to validate found solution
	sequence := make([]uint64, k+1)
	sequence[k] = x0

	log.Printf("x0: %v\n", x0)

	// Apply F() k times to x0 while also applying XOR on every iteration to make calculation dependent on a step number
	for i := uint64(1); i <= uint64(k); i++ {
		xk = method.F(xk) ^ i
		log.Printf("x%v: %v\n", i, xk)
		sequence[uint64(k)-i] = xk
	}

	log.Printf("sequence: %v\n", sequence)

	seqChecksum := checksum.Make(sequence)

	return &Challenge{
		X0:       x0,
		Xk:       xk,
		K:        k,
		N:        n,
		Checksum: seqChecksum,
	}
}