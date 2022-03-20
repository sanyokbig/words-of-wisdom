package checksum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Make(t *testing.T) {
	tests := []struct {
		name     string
		sequence []uint64
		wantHash string
	}{
		{
			name: "sequence with 3 tiny elements",
			sequence: []uint64{
				1, 2, 3,
			},
			wantHash: "fb59f40ac082d9075b753ee3dede9b654b696b478e33c9701fcc7da5687290e0",
		},
		{
			name: "sequence with 3 enormous elements",
			sequence: []uint64{
				^uint64(0),
				^uint64(0) - 1,
				^uint64(0) - 2,
			},
			wantHash: "d873983c0a6822fd8c0fce4004c3d729698dc6abb04f1b48652d165a7769ae89",
		},
		{
			name:     "sequence with no elements",
			sequence: []uint64{},
			wantHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash := Make(tt.sequence)

			assert.Equal(t, tt.wantHash, gotHash)
		})
	}
}
