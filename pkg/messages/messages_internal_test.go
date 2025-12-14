// Package messages internal tests for private functions.
package messages

import (
	"testing"
)

// TestMessageFormatEmpty tests Format with empty verbose.
func TestMessageFormatEmpty(t *testing.T) {
	msg := Message{
		Code:    "TEST",
		Short:   "short %s",
		Verbose: "",
	}

	// Test verbose mode with empty verbose returns short
	result := msg.Format(true, "arg")
	// Vérification résultat
	if result != "short arg" {
		t.Errorf("Format(true) = %q, want %q", result, "short arg")
	}
}

// TestMessageFormatNoArgs tests Format without arguments.
func TestMessageFormatNoArgs(t *testing.T) {
	msg := Message{
		Code:    "TEST",
		Short:   "simple message",
		Verbose: "verbose message",
	}

	// Test short without args
	result := msg.Format(false)
	// Vérification résultat
	if result != "simple message" {
		t.Errorf("Format(false) = %q, want %q", result, "simple message")
	}

	// Test verbose without args
	result = msg.Format(true)
	// Vérification résultat
	if result != "verbose message" {
		t.Errorf("Format(true) = %q, want %q", result, "verbose message")
	}
}

// TestFormatShortWithVerboseHint tests FormatShort adds hint.
func TestFormatShortWithVerboseHint(t *testing.T) {
	msg := Message{
		Code:    "TEST",
		Short:   "short message",
		Verbose: "verbose message",
	}

	result := msg.FormatShort()
	// Vérification suffixe
	expected := "short message (--verbose pour détails)"
	if result != expected {
		t.Errorf("FormatShort() = %q, want %q", result, expected)
	}
}

// TestFormatShortWithoutVerbose tests FormatShort without verbose.
func TestFormatShortWithoutVerbose(t *testing.T) {
	msg := Message{
		Code:    "TEST",
		Short:   "short message",
		Verbose: "",
	}

	result := msg.FormatShort()
	// Vérification pas de suffixe
	if result != "short message" {
		t.Errorf("FormatShort() = %q, want %q", result, "short message")
	}
}

// TestRegistryInitialized tests that registry is initialized with messages.
func TestRegistryInitialized(t *testing.T) {
	// Vérification que le registre contient des messages
	if len(registry) == 0 {
		t.Error("registry is empty after init()")
	}

	// Vérification quelques règles essentielles
	essentialRules := []string{
		"KTN-FUNC-001",
		"KTN-VAR-001",
		"KTN-CONST-001",
		"KTN-COMMENT-001",
		"KTN-STRUCT-001",
		"KTN-TEST-004",
	}

	// Itération sur les règles essentielles
	for _, code := range essentialRules {
		// Vérification existence
		if _, ok := registry[code]; !ok {
			t.Errorf("registry missing essential rule %q", code)
		}
	}
}
