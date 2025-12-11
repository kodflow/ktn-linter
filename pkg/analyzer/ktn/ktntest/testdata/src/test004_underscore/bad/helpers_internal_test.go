package bad

import "testing"

// This file exists but formatMessage is not tested
// Should trigger KTN-TEST-003 suggesting Test_formatMessage
func Test_someOtherFunction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Dummy test to make the package have test files
			t.Log("dummy test")

		})
	}
}
