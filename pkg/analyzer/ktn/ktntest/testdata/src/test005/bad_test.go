package test005_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test005"
)

// TestStringLengthMultipleCases teste plusieurs cas sans table-driven (PAS BIEN)
func TestStringLengthMultipleCases(t *testing.T) {
	// Test cas 1
	if test005.StringLength("hello") != 5 {
		t.Error("longueur de 'hello' devrait être 5")
	}

	// Test cas 2
	if test005.StringLength("") != 0 {
		t.Error("longueur de '' devrait être 0")
	}

	// Test cas 3 - déclenche la règle (>= 2 assertions)
	if test005.StringLength("a") != 1 {
		t.Error("longueur de 'a' devrait être 1")
	}

	// Test cas 4
	if test005.StringLength("test string") != 11 {
		t.Error("longueur de 'test string' devrait être 11")
	}
}

// TestIsEmptyRepeatedAssertions teste avec assertions répétitives (PAS BIEN)
func TestIsEmptyRepeatedAssertions(t *testing.T) {
	// Première assertion
	result1 := test005.IsEmpty("")
	if !result1 {
		t.Errorf("IsEmpty('') devrait retourner true")
	}

	// Deuxième assertion
	result2 := test005.IsEmpty("hello")
	if result2 {
		t.Errorf("IsEmpty('hello') devrait retourner false")
	}

	// Troisième assertion - déclenche la règle
	result3 := test005.IsEmpty(" ")
	if result3 {
		t.Errorf("IsEmpty(' ') devrait retourner false")
	}

	// Quatrième assertion
	result4 := test005.IsEmpty("   ")
	if result4 {
		t.Errorf("IsEmpty('   ') devrait retourner false")
	}
}

// TestToUpperManyScenarios teste de nombreux scénarios sans structure (PAS BIEN)
func TestToUpperManyScenarios(t *testing.T) {
	// Scénario 1: string minuscule
	if test005.ToUpper("hello") != "HELLO" {
		t.Fatal("conversion de 'hello' échouée")
	}

	// Scénario 2: string déjà majuscule
	if test005.ToUpper("WORLD") != "WORLD" {
		t.Fatal("conversion de 'WORLD' échouée")
	}

	// Scénario 3: string mixte - déclenche la règle
	if test005.ToUpper("HeLLo") != "HELLO" {
		t.Fatal("conversion de 'HeLLo' échouée")
	}

	// Scénario 4: avec chiffres
	if test005.ToUpper("test123") != "TEST123" {
		t.Fatal("conversion de 'test123' échouée")
	}

	// Scénario 5: avec caractères spéciaux
	if test005.ToUpper("hello@world.com") != "HELLO@WORLD.COM" {
		t.Fatal("conversion de 'hello@world.com' échouée")
	}
}

// TestContainsManyChecks teste avec beaucoup de vérifications (PAS BIEN)
func TestContainsManyChecks(t *testing.T) {
	// Check 1
	if !test005.Contains("hello world", "world") {
		t.Error("devrait contenir 'world'")
	}

	// Check 2
	if test005.Contains("hello", "bye") {
		t.Error("ne devrait pas contenir 'bye'")
	}

	// Check 3 - déclenche la règle
	if !test005.Contains("test", "test") {
		t.Error("devrait contenir 'test'")
	}

	// Check 4
	if !test005.Contains("abc", "a") {
		t.Error("devrait contenir 'a'")
	}

	// Check 5
	if test005.Contains("", "x") {
		t.Error("string vide ne contient rien")
	}

	// Check 6
	if !test005.Contains("abcdef", "cde") {
		t.Error("devrait contenir 'cde'")
	}
}

// TestCountWordsMultipleInputs teste plusieurs inputs sans pattern (PAS BIEN)
func TestCountWordsMultipleInputs(t *testing.T) {
	// Input 1
	count1 := test005.CountWords("")
	if count1 != 0 {
		t.Errorf("string vide devrait avoir 0 mots, got %d", count1)
	}

	// Input 2
	count2 := test005.CountWords("hello")
	if count2 != 1 {
		t.Errorf("'hello' devrait avoir 1 mot, got %d", count2)
	}

	// Input 3 - déclenche la règle
	count3 := test005.CountWords("hello world")
	if count3 != 2 {
		t.Errorf("'hello world' devrait avoir 2 mots, got %d", count3)
	}

	// Input 4
	count4 := test005.CountWords("one two three")
	if count4 != 3 {
		t.Errorf("'one two three' devrait avoir 3 mots, got %d", count4)
	}

	// Input 5
	count5 := test005.CountWords("  spaces  around  ")
	if count5 != 2 {
		t.Errorf("'  spaces  around  ' devrait avoir 2 mots, got %d", count5)
	}
}
