// Internal tests for classify.go - private function tests.
package shared

import (
	"testing"
)

// Test_getVisibility tests the getVisibility function.
//
// Params:
//   - t: testing context
func Test_getVisibility(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want Visibility
	}{
		{
			name: "public uppercase",
			id:   "Foo",
			want: VisPublic,
		},
		{
			name: "private lowercase",
			id:   "foo",
			want: VisPrivate,
		},
		{
			name: "public uppercase acronym",
			id:   "HTTPHandler",
			want: VisPublic,
		},
		{
			name: "private underscore prefix",
			id:   "_private",
			want: VisPrivate,
		},
		{
			name: "empty string",
			id:   "",
			want: VisPrivate,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVisibility(tt.id)
			// Check result
			if got != tt.want {
				t.Errorf("getVisibility(%q) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

// Test_parsePrivateTestName tests the parsePrivateTestName function.
//
// Params:
//   - t: testing context
func Test_parsePrivateTestName(t *testing.T) {
	tests := []struct {
		name        string
		privateName string
		wantTarget  TestTarget
		wantOK      bool
	}{
		{
			name:        "empty string",
			privateName: "",
			wantTarget:  TestTarget{},
			wantOK:      false,
		},
		{
			name:        "simple private function",
			privateName: "foo",
			wantTarget: TestTarget{
				FuncName:  "foo",
				IsPrivate: true,
				IsMethod:  false,
			},
			wantOK: true,
		},
		{
			name:        "private method Type_method",
			privateName: "Type_method",
			wantTarget: TestTarget{
				FuncName:     "method",
				ReceiverName: "Type",
				IsPrivate:    true,
				IsMethod:     true,
			},
			wantOK: true,
		},
		{
			name:        "private method with public method name",
			privateName: "Type_Method",
			wantTarget: TestTarget{
				FuncName:     "Method",
				ReceiverName: "Type",
				IsPrivate:    false,
				IsMethod:     true,
			},
			wantOK: true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePrivateTestName(tt.privateName)
			// Check success status
			if ok != tt.wantOK {
				t.Errorf("parsePrivateTestName(%q) ok = %v, want %v", tt.privateName, ok, tt.wantOK)
				return
			}
			// Check target if successful
			if ok && got != tt.wantTarget {
				t.Errorf("parsePrivateTestName(%q) = %+v, want %+v", tt.privateName, got, tt.wantTarget)
			}
		})
	}
}
