package solver

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildInversionTable(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want map[uint64][]uint64
	}{
		{
			name: "small table",
			args: args{
				n: 4,
			},
			want: map[uint64][]uint64{
				0: {1, 8, 15},
				1: {0, 7, 14},
				2: {6, 13},
				3: {5, 12},
				4: {4, 11},
				5: {3, 10},
				6: {2, 9},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := stubMethod{
				max: uint64(math.Pow(2, float64(tt.args.n))) - 1,
			}

			got := buildInversionTable(tt.args.n, method)
			assert.Equalf(t, tt.want, got, "want: %+v\ngot: %+v", tt.want, got)
		})
	}
}

type stubMethod struct {
	max uint64
}

// Calculated so that there are several cases of the same result on different x
func (s stubMethod) F(x uint64) uint64 {
	return (s.max - x) % (s.max / 2)
}

func TestSolver_Solve(t *testing.T) {
	type fields struct {
		xk             uint64
		targetDepth    int
		checksum       string
		inversionTable map[uint64][]uint64
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue uint64
		wantOk    bool
	}{
		{
			name: "value found",
			fields: fields{
				xk:          6,
				targetDepth: 3,
				checksum:    "d5c9c3aa15c6442d033e59c755260eff",
				inversionTable: map[uint64][]uint64{
					0:  {1, 8, 15},
					1:  {0, 7, 14},
					2:  {6, 13},
					3:  {5, 12},
					4:  {4, 11},
					5:  {3, 10},
					6:  {2, 9},
					9:  {11, 12, 13},
					12: {13},
					13: {2, 4, 5},
				},
			},
			wantValue: 9,
			wantOk:    true,
		},
		{
			name: "no valid checksum",
			fields: fields{
				xk:          6,
				targetDepth: 3,
				checksum:    "i-am-always-invalid",
				inversionTable: map[uint64][]uint64{
					0:  {1, 8, 15},
					1:  {0, 7, 14},
					2:  {6, 13},
					3:  {5, 12},
					4:  {4, 11},
					5:  {3, 10},
					6:  {2, 9},
					9:  {11, 12, 13},
					12: {13},
					13: {2, 4, 5},
				},
			},
			wantValue: 0,
			wantOk:    false,
		},
		{
			name: "no valid path in 3 steps",
			fields: fields{
				xk:          1,
				targetDepth: 3,
				checksum:    "i-won't-be-used",
				inversionTable: map[uint64][]uint64{
					1: {2},
					2: {3},
					3: {4},
					4: {5},
				},
			},
			wantValue: 0,
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solver{
				xk:             tt.fields.xk,
				targetDepth:    tt.fields.targetDepth,
				checksum:       tt.fields.checksum,
				inversionTable: tt.fields.inversionTable,
			}
			gotValue, gotOk := s.Solve()

			assert.Equal(t, tt.wantValue, gotValue)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}