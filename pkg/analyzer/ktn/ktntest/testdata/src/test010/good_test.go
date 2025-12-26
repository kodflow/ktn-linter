// Tests with proper assertions.
package test012

import (
	"testing"
)

// TestProcessData teste ProcessData avec des assertions.
func TestProcessData(t *testing.T) {
	// Test avec assertion
	result := ProcessData("hello")
	// Vérification avec assertion
	if result != "processed:hello" {
		t.Errorf("ProcessData() = %v, want %v", result, "processed:hello")
	}
}

// TestGetCount teste GetCount avec des assertions.
func TestGetCount(t *testing.T) {
	// Test avec assertion
	got := GetCount()
	// Vérification avec assertion
	if got != 42 {
		t.Fatalf("GetCount() = %d, want 42", got)
	}
}

// TestWithComparison teste avec une comparaison.
func TestWithComparison(t *testing.T) {
	// La comparaison == est un signal de validation
	result := GetCount()
	// Comparaison
	_ = result == 42
}

// TestWithSubtest teste avec t.Run.
func TestWithSubtest(t *testing.T) {
	// t.Run est un signal de validation
	t.Run("subtest", func(t *testing.T) {
		// Sous-test
		t.Log("in subtest")
	})
}

// TestWithError teste avec t.Error.
func TestWithError(t *testing.T) {
	// t.Error est une assertion
	if false {
		t.Error("should not happen")
	}
}

// TestWithFatal teste avec t.Fatal.
func TestWithFatal(t *testing.T) {
	// t.Fatal est une assertion
	if false {
		t.Fatal("should not happen")
	}
}

// TestWithFail teste avec t.Fail.
func TestWithFail(t *testing.T) {
	// t.Fail est une assertion
	if false {
		t.Fail()
	}
}

// TestWithFailNow teste avec t.FailNow.
func TestWithFailNow(t *testing.T) {
	// t.FailNow est une assertion
	if false {
		t.FailNow()
	}
}

// TestWithErrCheck teste la vérification d'erreur.
func TestWithErrCheck(t *testing.T) {
	// Vérification d'erreur
	err := GetError(false)
	// La comparaison != nil est un signal de validation
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestWithHelper teste avec un helper de test.
func TestWithHelper(t *testing.T) {
	// Appeler un helper qui prend t comme premier argument
	helperFunc(t)
}

// helperFunc est un helper de test.
func helperFunc(t *testing.T) {
	// Helper
	t.Helper()
}

// TestTableDriven est un test table-driven.
func TestTableDriven(t *testing.T) {
	// Test table-driven
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "processed:"},
		{"hello", "hello", "processed:hello"},
	}

	// Itérer sur les cas
	for _, tc := range cases {
		// Sous-test
		t.Run(tc.name, func(t *testing.T) {
			// Vérification
			got := ProcessData(tc.input)
			// Assertion
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
