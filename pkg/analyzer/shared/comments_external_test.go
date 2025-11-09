package shared_test

import (
	"go/ast"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
)

func TestHasValidComment(t *testing.T) {
	tests := []struct {
		name     string
		comments *ast.CommentGroup
		want     bool
	}{
		{
			name:     "nil comment group",
			comments: nil,
			want:     false,
		},
		{
			name: "empty comment group",
			comments: &ast.CommentGroup{
				List: []*ast.Comment{},
			},
			want: false,
		},
		{
			name: "valid comment",
			comments: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// This is a valid comment"},
				},
			},
			want: true,
		},
		{
			name: "want directive only",
			comments: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// want \"error\""},
				},
			},
			want: false,
		},
		{
			name: "want directive with space",
			comments: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "/* want error */"},
				},
			},
			want: false,
		},
		{
			name: "mixed comments with want",
			comments: &ast.CommentGroup{
				List: []*ast.Comment{
					{Text: "// want \"error\""},
					{Text: "// This is valid"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shared.HasValidComment(tt.comments)
			if got != tt.want {
				t.Errorf("HasValidComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
