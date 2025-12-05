// Internal tests for 005.go - struct documentation analyzer.
package ktncomment

import (
	"go/ast"
	"testing"
)

// Test_hasValidDocumentation tests the hasValidDocumentation function.
//
// Params:
//   - t: testing context
func Test_hasValidDocumentation(t *testing.T) {
	tests := []struct {
		name       string
		doc        *ast.CommentGroup
		structName string
		want       bool
	}{
		{
			name:       "nil documentation",
			doc:        nil,
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "empty comment list",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "single line doc",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct does something"},
				},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "valid two line doc",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct represents a structure."},
					{Text: "// It provides functionality for testing."},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
		{
			name: "doc not starting with struct name",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// This struct does something."},
					{Text: "// More description here."},
				},
			},
			structName: "MyStruct",
			want:       false,
		},
		{
			name: "doc with empty lines",
			doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// MyStruct is a test structure."},
					{Text: "//"},
					{Text: "// More info here."},
				},
			},
			structName: "MyStruct",
			want:       true,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasValidDocumentation(tt.doc, tt.structName)
			// Check result
			if got != tt.want {
				t.Errorf("hasValidDocumentation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runComment005 tests the runComment005 function configuration.
//
// Params:
//   - t: testing context
func Test_runComment005(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "runComment005 is tested via analysistest"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that Analyzer005 is properly configured
			if Analyzer005 == nil {
				t.Error("Analyzer005 should not be nil")
				return
			}
			// Check analyzer name
			if Analyzer005.Name != "ktncomment005" {
				t.Errorf("Analyzer005.Name = %q, want %q", Analyzer005.Name, "ktncomment005")
			}
		})
	}
}
