package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Make returns a checksum of sequence
// Checksum is calculating by:
// 1. Transform every number of sequence to its hex representation with leading zeroes
// 2. Concatenating resulted hexes
func Make(sequence []uint64) string {
	// Uint64 consists of 64*8=512 bits. Since hex represents 16 bits, uint64 as hex will have at most 16 hexes.
	// Since we will be padding hexes with zeros to have a unique source string, it is safe to assume that we will need
	// a number of sequences times 16
	joinedHexes := make([]byte, 0, len(sequence)*16)

	for _, v := range sequence {
		joinedHexes = append(joinedHexes, []byte(fmt.Sprintf("%016x", v))...)
	}

	hexHash := sha256.Sum256(joinedHexes)
	hash := hex.EncodeToString(hexHash[:])

	return hash
}
