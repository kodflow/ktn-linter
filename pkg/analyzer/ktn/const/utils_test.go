package ktnconst

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestHasValidComment teste la fonction hasValidComment.
func TestHasValidComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected bool
	}{
		{
			name:     "nil comment group",
			comment:  "",
			expected: false,
		},
		{
			name:     "valid comment",
			comment:  "// This is a valid comment",
			expected: true,
		},
		{
			name:     "want directive line comment",
			comment:  "// want something",
			expected: false,
		},
		{
			name:     "want directive block comment with space",
			comment:  "/* want something */",
			expected: false,
		},
		{
			name:     "want directive without space",
			comment:  "//want no space",
			expected: false,
		},
		{
			name:     "want directive block no space",
			comment:  "/*want block*/",
			expected: false,
		},
		{
			name:     "mixed valid and want",
			comment:  "// Valid\n// want directive",
			expected: true,
		},
		{
			name:     "only valid comments",
			comment:  "// First\n// Second",
			expected: true,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cg *ast.CommentGroup
			// Vérification si le commentaire n'est pas vide
			if tt.comment != "" {
				// Parsing du commentaire
				src := "package test\n" + tt.comment + "\nconst X = 1"
				fset := token.NewFileSet()
				file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
				// Vérification d'erreur de parsing
				if err != nil {
					t.Fatalf("Failed to parse: %v", err)
				}
				// Vérification que le fichier a des commentaires
				if len(file.Comments) > 0 {
					cg = file.Comments[0]
				}
			}

			result := hasValidComment(cg)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("hasValidComment() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsValidConstantName teste la fonction isValidConstantName.
func TestIsValidConstantName(t *testing.T) {
	tests := []struct {
		name     string
		constName string
		expected  bool
	}{
		{
			name:      "valid SCREAMING_SNAKE_CASE",
			constName: "MAX_SIZE",
			expected:  true,
		},
		{
			name:      "valid single letter",
			constName: "A",
			expected:  true,
		},
		{
			name:      "valid acronym",
			constName: "HTTP",
			expected:  true,
		},
		{
			name:      "valid with numbers",
			constName: "HTTP2",
			expected:  true,
		},
		{
			name:      "valid with underscores and numbers",
			constName: "TLS1_2_VERSION",
			expected:  true,
		},
		{
			name:      "invalid camelCase",
			constName: "maxSize",
			expected:  false,
		},
		{
			name:      "invalid PascalCase",
			constName: "MaxSize",
			expected:  false,
		},
		{
			name:      "invalid snake_case",
			constName: "max_size",
			expected:  false,
		},
		{
			name:      "invalid mixed case",
			constName: "Max_Size",
			expected:  false,
		},
		{
			name:      "invalid starts with number",
			constName: "1MAX",
			expected:  false,
		},
		{
			name:      "invalid with special chars",
			constName: "MAX-SIZE",
			expected:  false,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidConstantName(tt.constName)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isValidConstantName(%q) = %v, want %v", tt.constName, result, tt.expected)
			}
		})
	}
}
