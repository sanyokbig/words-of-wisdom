package cryptorand

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
)

func Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Panicf("failed to ready from crypto/rand reader: %v", err)
	}

	return v
}
