// Tests without assertions (passthrough tests).
package test012

import (
	"testing"
)

// TestBadProcessData est un test passthrough sans assertion.
//
// Params:
//   - t: contexte de test
func TestBadProcessData(t *testing.T) { // want "KTN-TEST-012: le test 'TestBadProcessData' est un test passthrough"
	// Test passthrough - pas d'assertion
	_ = ProcessData("hello")
}

// TestBadGetCount est un test passthrough sans assertion.
//
// Params:
//   - t: contexte de test
func TestBadGetCount(t *testing.T) { // want "KTN-TEST-012: le test 'TestBadGetCount' est un test passthrough"
	// Test passthrough - appelle la fonction mais ne vérifie rien
	GetCount()
}

// TestEmptyTest est un test vide.
//
// Params:
//   - t: contexte de test
func TestEmptyTest(t *testing.T) { // want "KTN-TEST-012: le test 'TestEmptyTest' est un test passthrough"
	// Corps vide - passthrough
}

// TestOnlyLog est un test avec seulement t.Log.
//
// Params:
//   - t: contexte de test
func TestOnlyLog(t *testing.T) { // want "KTN-TEST-012: le test 'TestOnlyLog' est un test passthrough"
	// t.Log n'est PAS une assertion
	t.Log("This is not a test")
}

// TestOnlyParallel est un test avec seulement t.Parallel.
//
// Params:
//   - t: contexte de test
func TestOnlyParallel(t *testing.T) { // want "KTN-TEST-012: le test 'TestOnlyParallel' est un test passthrough"
	// t.Parallel n'est pas une assertion
	t.Parallel()
}

// TestOnlyMarkAsAid est un test avec seulement t.Helper.
//
// Params:
//   - t: contexte de test
func TestOnlyMarkAsAid(t *testing.T) { // want "KTN-TEST-012: le test 'TestOnlyMarkAsAid' est un test passthrough"
	// t.Helper n'est pas une assertion
	t.Helper()
}

// TestOnlySkip est un test avec seulement t.Skip.
//
// Params:
//   - t: contexte de test
func TestOnlySkip(t *testing.T) { // want "KTN-TEST-012: le test 'TestOnlySkip' est un test passthrough"
	// t.Skip n'est pas une assertion
	t.Skip("skipping")
}

// TestOnlyCleanup est un test avec seulement t.Cleanup.
//
// Params:
//   - t: contexte de test
func TestOnlyCleanup(t *testing.T) { // want "KTN-TEST-012: le test 'TestOnlyCleanup' est un test passthrough"
	// t.Cleanup n'est pas une assertion
	t.Cleanup(func() {})
}
