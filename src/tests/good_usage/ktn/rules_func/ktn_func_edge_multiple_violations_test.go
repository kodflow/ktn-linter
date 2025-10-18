package rules_func

import (
	"testing"

)

// TestGoodFunctionName teste TODO.
//
// Params:
//   - t: contexte de test
func TestGoodFunctionName(t *testing.T) {
	cfg := MultiParamConfig{
		A: 1,
		B: 2,
		C: 3,
		D: 4,
		E: 5,
		F: 6,
	}

	result := GoodFunctionName(cfg)
	expected := 0 // Aucun nombre < 10 n'est multiple de 210
	if result != expected {
		t.Errorf("GoodFunctionName() = %v, want %v", result, expected)
	}
}

// TestShouldProcess teste TODO.
//
// Params:
//   - t: contexte de test
func TestShouldProcess(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want bool
	}{
		{"not multiple of 2", 1, false},
		{"not multiple of 3", 2, false},
		{"not multiple of 5", 3, false},
		{"not multiple of 7", 6, false},
		{"multiple of all", 210, true},
	}

	for _, tt := range tests {
		tt := tt // copie pour éviter capture de boucle
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldProcess(tt.i)
			if got != tt.want {
				t.Errorf("ShouldProcess(%v) = %v, want %v", tt.i, got, tt.want)
			}
		})
	}
}

// TestSumConfig teste TODO.
//
// Params:
//   - t: contexte de test
func TestSumConfig(t *testing.T) {
	cfg := MultiParamConfig{A: 1, B: 2, C: 3, D: 4, E: 5, F: 6}
	result := SumConfig(cfg)
	expected := 21
	if result != expected {
		t.Errorf("SumConfig() = %v, want %v", result, expected)
	}
}
