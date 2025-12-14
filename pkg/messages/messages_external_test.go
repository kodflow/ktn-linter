// Package messages_test provides black-box tests for the messages package.
package messages_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/messages"
)

// TestGet tests the Get function for retrieving messages.
func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantFound bool
	}{
		{
			name:      "existing rule FUNC-001",
			code:      "KTN-FUNC-001",
			wantFound: true,
		},
		{
			name:      "existing rule VAR-001",
			code:      "KTN-VAR-001",
			wantFound: true,
		},
		{
			name:      "existing rule COMMENT-001",
			code:      "KTN-COMMENT-001",
			wantFound: true,
		},
		{
			name:      "non-existing rule",
			code:      "KTN-INVALID-999",
			wantFound: false,
		},
		{
			name:      "empty code",
			code:      "",
			wantFound: false,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, found := messages.Get(tt.code)
			// Vérification du résultat
			if found != tt.wantFound {
				t.Errorf("Get(%q) found = %v, want %v", tt.code, found, tt.wantFound)
			}
			// Vérification du code si trouvé
			if found && msg.Code != tt.code {
				t.Errorf("Get(%q) msg.Code = %q, want %q", tt.code, msg.Code, tt.code)
			}
		})
	}
}

// TestMessageFormat tests the Format method of Message.
func TestMessageFormat(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		verbose  bool
		args     []any
		wantLen  int
	}{
		{
			name:    "short message without args",
			code:    "KTN-FUNC-001",
			verbose: false,
			args:    nil,
			wantLen: 1,
		},
		{
			name:    "verbose message without args",
			code:    "KTN-FUNC-001",
			verbose: true,
			args:    nil,
			wantLen: 10,
		},
		{
			name:    "short message with args",
			code:    "KTN-FUNC-001",
			verbose: false,
			args:    []any{3},
			wantLen: 1,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, found := messages.Get(tt.code)
			// Vérification existence
			if !found {
				t.Fatalf("Get(%q) not found", tt.code)
			}
			result := msg.Format(tt.verbose, tt.args...)
			// Vérification longueur minimale
			if len(result) < tt.wantLen {
				t.Errorf("Format() len = %d, want >= %d", len(result), tt.wantLen)
			}
		})
	}
}

// TestMessageFormatShort tests the FormatShort method.
func TestMessageFormatShort(t *testing.T) {
	msg, found := messages.Get("KTN-FUNC-001")
	// Vérification existence
	if !found {
		t.Fatal("Get(KTN-FUNC-001) not found")
	}

	result := msg.FormatShort(2)
	// Vérification contenu
	if len(result) == 0 {
		t.Error("FormatShort() returned empty string")
	}
}

// TestMessageFormatVerbose tests the FormatVerbose method.
func TestMessageFormatVerbose(t *testing.T) {
	msg, found := messages.Get("KTN-FUNC-001")
	// Vérification existence
	if !found {
		t.Fatal("Get(KTN-FUNC-001) not found")
	}

	result := msg.FormatVerbose(2)
	// Vérification contenu verbose
	if len(result) < 50 {
		t.Errorf("FormatVerbose() len = %d, want >= 50", len(result))
	}
}

// TestRegister tests the Register function.
func TestRegister(t *testing.T) {
	// Enregistrer un message de test
	testMsg := messages.Message{
		Code:    "KTN-TEST-999",
		Short:   "test message",
		Verbose: "verbose test message",
	}
	messages.Register(testMsg)

	// Vérifier qu'il est récupérable
	got, found := messages.Get("KTN-TEST-999")
	// Vérification existence
	if !found {
		t.Fatal("Register() message not found after registration")
	}
	// Vérification code
	if got.Code != testMsg.Code {
		t.Errorf("Register() Code = %q, want %q", got.Code, testMsg.Code)
	}
	// Vérification short
	if got.Short != testMsg.Short {
		t.Errorf("Register() Short = %q, want %q", got.Short, testMsg.Short)
	}
}
