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
		tt := tt // Capture range variable
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

// TestMessage_Format tests the Format method of Message.
func TestMessage_Format(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		verbose bool
		args    []any
		wantLen int
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
		tt := tt // Capture range variable
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

// TestMessage_FormatShort tests the FormatShort method.
func TestMessage_FormatShort(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		args       []any
		wantMinLen int
	}{
		{
			name:       "FUNC-001 short message",
			code:       "KTN-FUNC-001",
			args:       []any{2},
			wantMinLen: 1,
		},
		{
			name:       "VAR-001 short message",
			code:       "KTN-VAR-001",
			args:       nil,
			wantMinLen: 1,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := messages.Get(tt.code)
			// Verify message exists
			if !found {
				t.Fatalf("Get(%q) not found", tt.code)
			}
			result := msg.FormatShort(tt.args...)
			// Verify result is not empty
			if len(result) < tt.wantMinLen {
				t.Errorf("FormatShort() len = %d, want >= %d", len(result), tt.wantMinLen)
			}
		})
	}
}

// TestMessage_FormatVerbose tests the FormatVerbose method.
func TestMessage_FormatVerbose(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		args       []any
		wantMinLen int
	}{
		{
			name:       "FUNC-001 verbose message",
			code:       "KTN-FUNC-001",
			args:       []any{2},
			wantMinLen: 50,
		},
		{
			name:       "VAR-001 verbose message",
			code:       "KTN-VAR-001",
			args:       nil,
			wantMinLen: 10,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, found := messages.Get(tt.code)
			// Verify message exists
			if !found {
				t.Fatalf("Get(%q) not found", tt.code)
			}
			result := msg.FormatVerbose(tt.args...)
			// Verify verbose result has minimum length
			if len(result) < tt.wantMinLen {
				t.Errorf("FormatVerbose() len = %d, want >= %d", len(result), tt.wantMinLen)
			}
		})
	}
}

// TestRegister tests the Register function.
func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		msg     messages.Message
		wantErr bool
	}{
		{
			name: "register new message",
			msg: messages.Message{
				Code:    "KTN-TEST-999",
				Short:   "test message",
				Verbose: "verbose test message",
			},
			wantErr: false,
		},
		{
			name: "register another message",
			msg: messages.Message{
				Code:    "KTN-TEST-998",
				Short:   "another test",
				Verbose: "another verbose",
			},
			wantErr: false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Register the message
			messages.Register(tt.msg)
			// Verify it can be retrieved
			got, found := messages.Get(tt.msg.Code)
			// Verify message exists after registration
			if !found {
				t.Fatal("Register() message not found after registration")
			}
			// Verify code matches
			if got.Code != tt.msg.Code {
				t.Errorf("Register() Code = %q, want %q", got.Code, tt.msg.Code)
			}
			// Verify short matches
			if got.Short != tt.msg.Short {
				t.Errorf("Register() Short = %q, want %q", got.Short, tt.msg.Short)
			}
		})
	}
}

// TestNewMessage tests the NewMessage function.
func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		short   string
		verbose string
	}{
		{
			name:    "simple message",
			code:    "KTN-TEST-001",
			short:   "short message",
			verbose: "verbose message",
		},
		{
			name:    "message with empty verbose",
			code:    "KTN-TEST-002",
			short:   "short only",
			verbose: "",
		},
		{
			name:    "message with format verbs",
			code:    "KTN-TEST-003",
			short:   "value is %d",
			verbose: "The value %d exceeds %d",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg := messages.NewMessage(tt.code, tt.short, tt.verbose)
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("NewMessage() Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short matches
			if msg.Short != tt.short {
				t.Errorf("NewMessage() Short = %q, want %q", msg.Short, tt.short)
			}
			// Verify verbose matches
			if msg.Verbose != tt.verbose {
				t.Errorf("NewMessage() Verbose = %q, want %q", msg.Verbose, tt.verbose)
			}
		})
	}
}

// TestGetAll tests the GetAll function.
func TestGetAll(t *testing.T) {
	tests := []struct {
		name        string
		wantMinLen  int
		wantNonEmpty bool
	}{
		{name: "returns_messages", wantMinLen: 1, wantNonEmpty: true},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			allMsgs := messages.GetAll()
			// Verify we have messages
			if len(allMsgs) < tt.wantMinLen {
				t.Fatalf("GetAll() returned %d messages, want >= %d", len(allMsgs), tt.wantMinLen)
			}
			// Verify all messages have codes if required
			if tt.wantNonEmpty {
				for _, msg := range allMsgs {
					// Check code is not empty
					if msg.Code == "" {
						t.Error("GetAll() contains message with empty Code")
					}
				}
			}
		})
	}
}

// TestClear tests the Clear function.
func TestClear(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		short    string
		verbose  string
		wantEmpty bool
	}{
		{
			name:     "clears_registry",
			code:     "KTN-CLEAR-TEST-001",
			short:    "test short",
			verbose:  "test verbose",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Register a test message first
			testMsg := messages.Message{Code: tt.code, Short: tt.short, Verbose: tt.verbose}
			messages.Register(testMsg)
			// Verify it exists
			_, found := messages.Get(testMsg.Code)
			// Check registration
			if !found {
				t.Fatal("Message not registered before Clear test")
			}
			// Clear and verify registry is empty
			messages.Clear()
			allMsgs := messages.GetAll()
			// Verify clear worked
			if tt.wantEmpty && len(allMsgs) != 0 {
				t.Errorf("Clear() did not empty registry, got %d messages", len(allMsgs))
			}
			// Re-register essential messages for other tests (reinit)
			messages.Register(messages.Message{Code: "KTN-FUNC-001", Short: "test", Verbose: "test"})
			messages.Register(messages.Message{Code: "KTN-VAR-001", Short: "test", Verbose: "test"})
		})
	}
}

// TestUnregister tests the Unregister function.
func TestUnregister(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		short   string
		verbose string
	}{
		{
			name:    "unregisters_message",
			code:    "KTN-UNREG-TEST-001",
			short:   "test short",
			verbose: "test verbose",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Register a test message first
			testMsg := messages.Message{Code: tt.code, Short: tt.short, Verbose: tt.verbose}
			messages.Register(testMsg)
			// Verify it exists
			_, found := messages.Get(testMsg.Code)
			// Check registration succeeded
			if !found {
				t.Fatal("Message not registered before Unregister test")
			}
			// Unregister it
			messages.Unregister(testMsg.Code)
			// Verify it no longer exists
			_, found = messages.Get(testMsg.Code)
			// Check unregistration succeeded
			if found {
				t.Error("Unregister() did not remove message")
			}
		})
	}
}
