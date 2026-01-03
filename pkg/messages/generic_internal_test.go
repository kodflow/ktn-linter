package messages

import "testing"

// Test_registerGenericMessages tests that all generic messages are properly registered.
func Test_registerGenericMessages(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "KTN-GENERIC-001_registered", code: "KTN-GENERIC-001"},
		{name: "KTN-GENERIC-002_registered", code: "KTN-GENERIC-002"},
		{name: "KTN-GENERIC-003_registered", code: "KTN-GENERIC-003"},
		{name: "KTN-GENERIC-005_registered", code: "KTN-GENERIC-005"},
		{name: "KTN-GENERIC-006_registered", code: "KTN-GENERIC-006"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Verifier que le message est enregistre
			if !ok {
				t.Errorf("message %s not found in registry", tt.code)
				// Retour anticipe
				return
			}
			// Verifier que le message n'est pas vide
			if msg.Short == "" {
				t.Errorf("message %s has empty Short", tt.code)
			}
			// Verifier que le message verbose n'est pas vide
			if msg.Verbose == "" {
				t.Errorf("message %s has empty Verbose", tt.code)
			}
		})
	}
}

// Test_registerGeneric001 tests that KTN-GENERIC-001 message is registered correctly.
func Test_registerGeneric001(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "message_registered", code: "KTN-GENERIC-001"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Check registration
			if !ok {
				t.Fatalf("%s not registered", tt.code)
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Error("Short message is empty")
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Error("Verbose message is empty")
			}
		})
	}
}

// Test_registerGeneric002 tests that KTN-GENERIC-002 message is registered correctly.
func Test_registerGeneric002(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "message_registered", code: "KTN-GENERIC-002"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Check registration
			if !ok {
				t.Fatalf("%s not registered", tt.code)
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Error("Short message is empty")
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Error("Verbose message is empty")
			}
		})
	}
}

// Test_registerGeneric003 tests that KTN-GENERIC-003 message is registered correctly.
func Test_registerGeneric003(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "message_registered", code: "KTN-GENERIC-003"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Check registration
			if !ok {
				t.Fatalf("%s not registered", tt.code)
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Error("Short message is empty")
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Error("Verbose message is empty")
			}
		})
	}
}

// Test_registerGeneric005 tests that KTN-GENERIC-005 message is registered correctly.
func Test_registerGeneric005(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "message_registered", code: "KTN-GENERIC-005"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Check registration
			if !ok {
				t.Fatalf("%s not registered", tt.code)
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Error("Short message is empty")
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Error("Verbose message is empty")
			}
		})
	}
}

// Test_registerGeneric006 tests that KTN-GENERIC-006 message is registered correctly.
func Test_registerGeneric006(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "message_registered", code: "KTN-GENERIC-006"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := Get(tt.code)
			// Check registration
			if !ok {
				t.Fatalf("%s not registered", tt.code)
			}
			// Verify code matches
			if msg.Code != tt.code {
				t.Errorf("Code = %q, want %q", msg.Code, tt.code)
			}
			// Verify short message not empty
			if msg.Short == "" {
				t.Error("Short message is empty")
			}
			// Verify verbose message not empty
			if msg.Verbose == "" {
				t.Error("Verbose message is empty")
			}
		})
	}
}
