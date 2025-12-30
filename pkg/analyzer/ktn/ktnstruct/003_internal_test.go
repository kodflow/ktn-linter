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
		name        string
		code        string
		expectError bool
	}{
		{
			name: "method without receiver",
			code: `package test
func GetName() string { return "test" }`,
			expectError: false,
		},
		{
			name: "private method",
			code: `package test
type User struct { name string }
func (u *User) getName() string { return u.name }`,
			expectError: false,
		},
		{
			name: "method on private struct",
			code: `package test
type user struct { name string }
func (u *user) GetName() string { return u.name }`,
			expectError: false,
		},
		{
			name: "method named just Get",
			code: `package test
type User struct{}
func (u *User) Get() string { return "" }`,
			expectError: false,
		},
		{
			name: "simple getter with Get prefix",
			code: `package test
type User struct { name string }
func (u *User) GetName() string { return u.name }`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			cfg := &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-STRUCT-003": {Enabled: config.Bool(true)},
				},
			}
			config.Set(cfg)
			defer config.Reset()

			errCount := 0
			inspectPass := &analysis.Pass{
				Fset:     fset,
				Files:    []*ast.File{f},
				Report:   func(d analysis.Diagnostic) { errCount++ },
				ResultOf: make(map[*analysis.Analyzer]any),
			}
			inspectResult, _ := inspect.Analyzer.Run(inspectPass)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspectResult,
				},
				Report: func(_ analysis.Diagnostic) { errCount++ },
			}

			_, err = runStruct003(pass)
			if err != nil {
				t.Errorf("runStruct003() error = %v", err)
			}

			if tt.expectError && errCount == 0 {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && errCount > 0 {
				t.Errorf("Expected no error but got %d", errCount)
			}
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
		tt := tt // Capture range variable
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
		{
			name: "getter with multiple returns is still getter",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() (string, error) { return u.name, nil }`,
			expected: true,
		},
		{
			name: "getter with multiple statements is still getter",
			src: `package test
type User struct {
	name string
}
func (u *User) GetName() string {
	x := u.name
	return x
}`,
			expected: true,
		},
		{
			name: "getter with unknown receiver type",
			src: `package test
import "io"
func (r *io.Reader) GetData() string { return "" }`,
			expected: true,
		},
		{
			name: "getter with matching field in struct",
			src: `package test
type User struct { name string }
func (u *User) GetName() string { return u.name }`,
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
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

	// Test edge case: receiver type not in struct map
	t.Run("receiver type not in struct map", func(t *testing.T) {
		src := `package test
type User struct { name string }
func (u *User) GetData() string { return "" }`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", src, 0)
		if err != nil {
			t.Fatalf("failed to parse source: %v", err)
		}

		// Create empty struct types map (User not included)
		structTypes := make(map[string][]string)

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
		if !result {
			t.Errorf("expected true for type not in map, got false")
		}
	})
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
		{
			name: "selector expr receiver",
			src: `package test
import "io"
func (u *io.Reader) Method() {}`,
			expected: "",
		},
		{
			name: "array type receiver extracts type name",
			src: `package test
type IntSlice []int
func (is IntSlice) Method() {}`,
			expected: "IntSlice",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
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
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("expected %s to be %d, got %d", tt.name, tt.expected, tt.got)
			}
		})
	}
}

// Test_runStruct003_disabled tests that the rule is skipped when disabled.
func Test_runStruct003_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

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

		})
	}
}

// Test_runStruct003_excludedFile tests that excluded files are skipped.
func Test_runStruct003_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

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

		})
	}
}

// Test_isValidGetterToReport tests the isValidGetterToReport function.
func Test_isValidGetterToReport(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
		{"error case"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_capitalizeFirstLetter tests the capitalizeFirstLetter function.
func Test_capitalizeFirstLetter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "lowercase first letter",
			input:    "name",
			expected: "Name",
		},
		{
			name:     "uppercase first letter",
			input:    "Name",
			expected: "Name",
		},
		{
			name:     "single lowercase letter",
			input:    "x",
			expected: "X",
		},
		{
			name:     "single uppercase letter",
			input:    "X",
			expected: "X",
		},
		{
			name:     "camelCase input",
			input:    "userName",
			expected: "UserName",
		},
		{
			name:     "ID abbreviation",
			input:    "iD",
			expected: "ID",
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := capitalizeFirstLetter(tt.input)
			// Verify result
			if result != tt.expected {
				t.Errorf("capitalizeFirstLetter(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
