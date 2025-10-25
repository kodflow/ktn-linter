package test005_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test005"
)

// TestCalculatorSingleCase teste un seul cas simple (BIEN - une seule assertion)
func TestCalculatorSingleCase(t *testing.T) {
	result, err := test005.Calculator("+", 2, 3)
	// Vérification pas d'erreur
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Vérification résultat
	if result != 5 {
		t.Errorf("got %d, want 5", result)
	}
}

// TestCalculatorAllOperations teste toutes les opérations avec table-driven (BIEN)
func TestCalculatorAllOperations(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a       int
		b       int
		want    int
		wantErr bool
	}{
		{
			name:    "addition simple",
			op:      "+",
			a:       2,
			b:       3,
			want:    5,
			wantErr: false,
		},
		{
			name:    "addition nombres négatifs",
			op:      "+",
			a:       -5,
			b:       -3,
			want:    -8,
			wantErr: false,
		},
		{
			name:    "soustraction",
			op:      "-",
			a:       10,
			b:       3,
			want:    7,
			wantErr: false,
		},
		{
			name:    "multiplication",
			op:      "*",
			a:       4,
			b:       5,
			want:    20,
			wantErr: false,
		},
		{
			name:    "division",
			op:      "/",
			a:       20,
			b:       4,
			want:    5,
			wantErr: false,
		},
		{
			name:    "division par zéro",
			op:      "/",
			a:       10,
			b:       0,
			want:    0,
			wantErr: true,
		},
		{
			name:    "opération invalide",
			op:      "%",
			a:       10,
			b:       3,
			want:    0,
			wantErr: true,
		},
	}

	// Parcours des cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := test005.Calculator(tt.op, tt.a, tt.b)
			// Vérification erreur attendue
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculator() error = %v, wantErr %v", err, tt.wantErr)
				// Retour anticipé
				return
			}
			// Vérification résultat
			if got != tt.want {
				t.Errorf("Calculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateEmailCases teste la validation d'email avec table-driven (BIEN)
func TestValidateEmailCases(t *testing.T) {
	testcases := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"john.doe@company.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"user@", false},
		{"user@@example.com", false},
		{"", false},
		{"a@b.c", true},
	}

	// Itération sur les cas
	for _, tc := range testcases {
		result := test005.ValidateEmail(tc.email)
		// Vérification résultat
		if result != tc.valid {
			t.Errorf("ValidateEmail(%q) = %v, want %v", tc.email, result, tc.valid)
		}
	}
}

// TestParseIntVariousCases teste ParseInt avec table-driven (BIEN)
func TestParseIntVariousCases(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    int
		shouldError bool
	}{
		{"nombre positif", "123", 123, false},
		{"nombre négatif", "-456", -456, false},
		{"zéro", "0", 0, false},
		{"string vide", "", 0, true},
		{"caractère non numérique", "12a3", 0, true},
		{"plusieurs tirets", "--123", 0, true},
	}

	// Parcours des cas
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			result, err := test005.ParseInt(c.input)
			// Vérification erreur
			if (err != nil) != c.shouldError {
				t.Errorf("ParseInt(%q) error = %v, shouldError %v", c.input, err, c.shouldError)
			}
			// Vérification résultat si pas d'erreur attendue
			if !c.shouldError && result != c.expected {
				t.Errorf("ParseInt(%q) = %d, expected %d", c.input, result, c.expected)
			}
		})
	}
}

// TestFactorialSmallNumbers teste factorielle avec sous-tests (BIEN)
func TestFactorialSmallNumbers(t *testing.T) {
	tests := map[string]struct {
		n    int
		want int
	}{
		"factorial 0": {0, 1},
		"factorial 1": {1, 1},
		"factorial 2": {2, 2},
		"factorial 3": {3, 6},
		"factorial 4": {4, 24},
		"factorial 5": {5, 120},
	}

	// Parcours des tests
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := test005.Factorial(tt.n)
			// Vérification pas d'erreur
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Vérification résultat
			if got != tt.want {
				t.Errorf("Factorial(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
