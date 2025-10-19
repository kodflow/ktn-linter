package ktnconst

import (
	"go/ast"
	"testing"
)

func TestHasValidComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  *ast.CommentGroup
		expected bool
	}{
		{
			name:     "nil comment group",
			comment:  nil,
			expected: false,
		},
		{
			name: "empty comment list",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{},
			},
			expected: false,
		},
		{
			name: "valid line comment",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// This is a valid comment"},
				},
			},
			expected: true,
		},
		{
			name: "valid block comment",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* This is a valid block comment */"},
				},
			},
			expected: true,
		},
		{
			name: "line comment with want directive",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// want \"some error\""},
				},
			},
			expected: false,
		},
		{
			name: "block comment with want directive (single space)",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* want \"some error\" */"},
				},
			},
			expected: false, // text[2:7] = " want" matches the pattern, so it's a want directive
		},
		{
			name: "block comment with want directive (double space)",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/*  want \"some error\" */"},
				},
			},
			expected: true, // text[2:7] = "  wan" doesn't match " want", so it's considered valid
		},
		{
			name: "mixed comments with want and valid",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// want \"ignored\""},
					{Text: "// This is valid"},
				},
			},
			expected: true,
		},
		{
			name: "only want directives",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// want \"error1\""},
					{Text: "/* want \"error2\" */"},
				},
			},
			expected: false,
		},
		{
			name: "short comment (less than 6 chars)",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// OK"},
				},
			},
			expected: true,
		},
		{
			name: "exactly 6 chars line comment",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "//want"},
				},
			},
			expected: false, // text[2:6] = "want" matches
		},
		{
			name: "exactly 7 chars block comment",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* want"},
				},
			},
			expected: false, // text[2:7] = " want" matches
		},
		{
			name: "short block comment (less than 7 chars)",
			comment: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* OK"},
				},
			},
			expected: true, // len < 7, so block comment check is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasValidComment(tt.comment)
			if result != tt.expected {
				t.Errorf("hasValidComment() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
