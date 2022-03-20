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
		name string
		args args
		rand RandUint64
		want *Challenge
	}{
		{
			name: "simple challenge",
			args: args{
				k: 5,
				n: 4,
			},
			rand: stubRand{u: 9}.Uint64,
			want: &Challenge{
				X0:       9,
				Xk:       8,
				K:        5,
				N:        4,
				Checksum: "c62c24a07a66dd420465bf9913c37bdadc2145ebb131c52588fab625956a5cc3",
			},
		},
		{
			name: "hard challenge",
			args: args{
				k: 32,
				n: 21,
			},
			rand: stubRand{u: 1899858}.Uint64,
			want: &Challenge{
				X0:       1899858,
				Xk:       1899890,
				K:        32,
				N:        21,
				Checksum: "369e061fc3cca37028aacf96a6b38973dfa8f9381e9bcfe8e78eb8e0214a5850",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := stubMethod{}

			got := New(tt.rand).Prepare(method, tt.args.n, tt.args.k)

			assert.Equalf(t, tt.want, got, "New(%v, %v)", tt.args.n, tt.args.k)
		})
	}
}

type stubMethod struct{}

func (s stubMethod) F(u uint64) uint64 {
	return u
}

type stubRand struct {
	u uint64
}

func (s stubRand) Uint64() uint64 {
	return s.u
}

func BenchmarkChallenger_Prepare(b *testing.B) {
	type args struct {
		k int
		n int
	}
	benchmarks := []struct {
		name string
		args args
		rand RandUint64
	}{
		{
			name: "n:4, k:5",
			args: args{
				k: 5,
				n: 4,
			},
			rand: rand.Uint64,
		},
		{
			name: "n:21, k:32",
			args: args{
				k: 32,
				n: 21,
			},
			rand: rand.Uint64,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			method := stubMethod{}

			for i := 0; i < b.N; i++ {
				New(bm.rand).Prepare(method, bm.args.n, bm.args.k)
			}
		})
	}
}
