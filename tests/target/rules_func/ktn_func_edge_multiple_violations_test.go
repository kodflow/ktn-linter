package rules_func_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_func"
)

func TestGoodFunctionName(t *testing.T) {
	cfg := rules_func.MultiParamConfig{
		A: 1,
		B: 2,
		C: 3,
		D: 4,
		E: 5,
		F: 6,
	}

	result := rules_func.GoodFunctionName(cfg)
	expected := 0 // Aucun nombre < 10 n'est multiple de 210
	if result != expected {
		t.Errorf("rules_func.GoodFunctionName() = %v, want %v", result, expected)
	}
}

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
		t.Run(tt.name, func(t *testing.T) {
			got := rules_func.ShouldProcess(tt.i)
			if got != tt.want {
				t.Errorf("rules_func.ShouldProcess(%v) = %v, want %v", tt.i, got, tt.want)
			}
		})
	}
}

func TestSumConfig(t *testing.T) {
	cfg := rules_func.MultiParamConfig{A: 1, B: 2, C: 3, D: 4, E: 5, F: 6}
	result := rules_func.SumConfig(cfg)
	expected := 21
	if result != expected {
		t.Errorf("rules_func.SumConfig() = %v, want %v", result, expected)
	}
}
