package challenger

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChallenger_Prepare(t *testing.T) {
	type args struct {
		k int
		n int
	}
	tests := []struct {
		name     string
		args     args
		randSeed int64
		want     *Challenge
	}{
		{
			name: "simple challenge",
			args: args{
				k: 5,
				n: 4,
			},
			randSeed: 11,
			want: &Challenge{
				X0:       9,
				Xk:       8,
				K:        5,
				N:        4,
				Checksum: "9631482d2c592fc903aafd3a5229fc79",
			},
		},
		{
			name: "hard challenge",
			args: args{
				k: 32,
				n: 21,
			},
			randSeed: 1,
			want: &Challenge{
				X0:       1899858,
				Xk:       1899890,
				K:        32,
				N:        21,
				Checksum: "fab9387142a8bce59bb374af29cecb54",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(tt.randSeed)

			method := stubMethod{}

			assert.Equalf(t, tt.want, New().Prepare(method, tt.args.n, tt.args.k), "New(%v, %v)", tt.args.n, tt.args.k)
		})
	}
}

type stubMethod struct{}

func (s stubMethod) F(u uint64) uint64 {
	return u
}
