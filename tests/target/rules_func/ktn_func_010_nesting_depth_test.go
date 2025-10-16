package rules_func_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_func"
)

// TestDeeplyNestedGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestDeeplyNestedGood(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "zero value",
			value: 0,
			want:  0,
		},
		{
			name:  "small value",
			value: 5,
			want:  18,
		},
		{
			name:  "larger value",
			value: 10,
			want:  90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules_func.DeeplyNestedGood(tt.value)
			if got != tt.want {
				t.Errorf("rules_func.DeeplyNestedGood() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtremelyNestedGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestExtremelyNestedGood(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{
			name: "zeros",
			x:    0,
			y:    0,
			want: 0,
		},
		{
			name: "small values",
			x:    2,
			y:    3,
			want: 12,
		},
		{
			name: "x is zero",
			x:    0,
			y:    5,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules_func.ExtremelyNestedGood(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("rules_func.ExtremelyNestedGood() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestComplexNestedGood teste TODO.
//
// Params:
//   - t: contexte de test
func TestComplexNestedGood(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "empty slice",
			values: []int{},
			want:   []int{},
		},
		{
			name:   "mixed values",
			values: []int{12, 15, 18, -5, 0, 105},
			want:   []int{12, 18, 30},
		},
		{
			name:   "all negative",
			values: []int{-1, -2, -3},
			want:   []int{},
		},
		{
			name:   "modulo 0 under 100",
			values: []int{9, 12, 15},
			want:   []int{12, 15, 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules_func.ComplexNestedGood(tt.values)
			if len(got) != len(tt.want) {
				t.Errorf("rules_func.ComplexNestedGood() length = %v, want %v", len(got), len(tt.want))
				// Retourne pour arrêter l'exécution du test
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("rules_func.ComplexNestedGood()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
