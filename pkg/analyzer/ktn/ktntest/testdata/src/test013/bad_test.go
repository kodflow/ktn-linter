// Tests without assertions (passthrough tests).
package test013_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test013"
)

// TestBadProcessData est un test passthrough sans assertion.
//
// Params:
//   - t: contexte de test
func TestBadProcessData(t *testing.T) { // want "KTN-TEST-013: le test 'TestBadProcessData' est un test passthrough"
	// Test passthrough - pas d'assertion
	_ = test013.ProcessData("hello")
}

// TestBadGetCount est un test passthrough sans assertion.
//
// Params:
//   - t: contexte de test
func TestBadGetCount(t *testing.T) { // want "KTN-TEST-013: le test 'TestBadGetCount' est un test passthrough"
	// Test passthrough - appelle la fonction mais ne vérifie rien
	test013.GetCount()
}

// TestEmptyTest est un test vide.
//
// Params:
//   - t: contexte de test
func TestEmptyTest(t *testing.T) { // want "KTN-TEST-013: le test 'TestEmptyTest' est un test passthrough"
	// Corps vide - passthrough
}
