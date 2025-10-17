package rules_func_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/good_usage/rules_func"
)

// TestFiveParamsIsOk teste la fonction avec 5 paramètres.
//
// Params:
//   - t: contexte de test
func TestFiveParamsIsOk(t *testing.T) {
	result := rules_func.FiveParamsIsOk(1, 2, 3, 4, 5)
	expected := 15
	if result != expected {
		t.Errorf("fiveParamsIsOk() = %v, want %v", result, expected)
	}
}

// TestSixParamsWithConfig teste la fonction avec struct de config.
//
// Params:
//   - t: contexte de test
func TestSixParamsWithConfig(t *testing.T) {
	cfg := rules_func.SixParamsConfig{
		A: 1,
		B: 2,
		C: 3,
		D: 4,
		E: 5,
		F: 6,
	}
	result := rules_func.SixParamsWithConfig(cfg)
	expected := 21
	if result != expected {
		t.Errorf("sixParamsWithConfig() = %v, want %v", result, expected)
	}
}
