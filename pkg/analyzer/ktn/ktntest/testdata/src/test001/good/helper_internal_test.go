package good

import "testing"

// This file follows the correct naming convention for internal tests

func testPrivateHelper() {
	// Test private function
}

func TestInternalHelper(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testPrivateHelper()
			t.Log("test internal helper")

		})
	}
}
