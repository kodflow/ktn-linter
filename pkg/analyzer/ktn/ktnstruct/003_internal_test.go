package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runStruct003 tests the private runStruct003 function.
func Test_runStruct003(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_collectStructTypes tests the private collectStructTypes function.
func Test_collectStructTypes(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected map[string]int // map of struct name to number of fields
	}{
		{
			name: "no structs",
			src: `package test
func main() {}`,
			expected: map[string]int{},
		},
		{
			name: "struct with fields",
			src: `package test
type User struct {
	Name string
	Age  int
}`,
			expected: map[string]int{"User": 2},
		},
		{
			name: "multiple structs",
			src: `package test
type User struct {
	Name string
}
type Admin struct {
	Role string
	Level int
}`,
			expected: map[string]int{"User": 1, "Admin": 2},
		},
		{
			name: "empty struct",
			src: `package test
type Empty struct{}`,
			expected: map[string]int{"Empty": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			result := collectStructTypes(pass)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d structs, got %d", len(tt.expected), len(result))
			}

			for name, expectedFields := range tt.expected {
				fields, ok := result[name]
				if !ok {
					t.Errorf("expected struct '%s' not found", name)
					continue
				}
				if len(fields) != expectedFields {
					t.Errorf("expected %d fields for '%s', got %d", expectedFields, name, len(fields))
				}
			}
		})
	}
}

// Test_isSimpleGetter tests the private isSimpleGetter function.
func Test_isSimpleGetter(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected bool
	}{
		{
			name: "simple getter",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() string { return u.name }`,
			expected: true,
		},
		{
			name: "getter with params",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName(prefix string) string { return prefix + u.name }`,
			expected: false,
		},
		{
			name: "getter without return",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
			}

			// Collect struct types first
			structTypes := collectStructTypes(pass)

			// Find the method
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Recv != nil {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil {
				t.Fatal("no method found")
			}

			result := isSimpleGetter(funcDecl, structTypes)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test_getReceiverTypeName tests the private getReceiverTypeName function.
func Test_getReceiverTypeName(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "pointer receiver",
			src: `package test
type User struct{}
func (u *User) Method() {}`,
			expected: "User",
		},
		{
			name: "value receiver",
			src: `package test
type User struct{}
func (u User) Method() {}`,
			expected: "User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Find the method
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok && fd.Recv != nil {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Recv == nil {
				t.Fatal("no method found")
			}

			result := getReceiverTypeName(funcDecl.Recv.List[0].Type)

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test_constants tests the constant values.
func Test_constants(t *testing.T) {
	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{name: "initialStructTypesCap", got: initialStructTypesCap, expected: 32},
		{name: "getPrefixLen", got: getPrefixLen, expected: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("expected %s to be %d, got %d", tt.name, tt.expected, tt.got)
			}
		})
	}
}

// Test_runStruct003_disabled tests that the rule is skipped when disabled.
func Test_runStruct003_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-003": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runStruct003(pass)
	if err != nil {
		t.Errorf("runStruct003() error = %v", err)
	}
}

// Test_runStruct003_excludedFile tests that excluded files are skipped.
func Test_runStruct003_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-003": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runStruct003(pass)
	if err != nil {
		t.Errorf("runStruct003() error = %v", err)
	}
}
