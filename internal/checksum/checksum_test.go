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
			wantHash: "da0a23fa895d44e039c77aaeec1c28bd",
		},
		{
			name: "sequence with 3 enormous elements",
			sequence: []uint64{
				^uint64(0),
				^uint64(0) - 1,
				^uint64(0) - 2,
			},
			wantHash: "fa65ce2585fc8be89c830a255c553eee",
		},
		{
			name:     "sequence with no elements",
			sequence: []uint64{},
			wantHash: "d41d8cd98f00b204e9800998ecf8427e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash := Make(tt.sequence)

			assert.Equal(t, tt.wantHash, gotHash)
		})
	}
}
