package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar021 tests the private runVar021 function.
func Test_runVar021(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_extractReceiverType tests the private extractReceiverType function.
func Test_extractReceiverType(t *testing.T) {
	tests := []struct {
		name         string
		expr         ast.Expr
		expectedName string
		expectedPtr  bool
	}{
		{
			name:         "pointer to ident",
			expr:         &ast.StarExpr{X: &ast.Ident{Name: "MyStruct"}},
			expectedName: "MyStruct",
			expectedPtr:  true,
		},
		{
			name:         "pointer to non-ident",
			expr:         &ast.StarExpr{X: &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}}},
			expectedName: "",
			expectedPtr:  false,
		},
		{
			name:         "ident type",
			expr:         &ast.Ident{Name: "MyStruct"},
			expectedName: "MyStruct",
			expectedPtr:  false,
		},
		{
			name:         "selector expr (not identifiable)",
			expr:         &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "Type"}},
			expectedName: "",
			expectedPtr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			name, isPtr := extractReceiverType(tt.expr)
			// Verification du nom
			if name != tt.expectedName {
				t.Errorf("extractReceiverType() name = %q, expected %q", name, tt.expectedName)
			}
			// Verification du pointeur
			if isPtr != tt.expectedPtr {
				t.Errorf("extractReceiverType() isPointer = %v, expected %v", isPtr, tt.expectedPtr)
			}
		})
	}
}

// Test_checkReceiverConsistency tests the private checkReceiverConsistency function.
func Test_checkReceiverConsistency(t *testing.T) {
	tests := []struct {
		name          string
		typeReceivers map[string][]receiverInfo
		expectReports int
	}{
		{
			name: "single method - no report",
			typeReceivers: map[string][]receiverInfo{
				"MyType": {
					{isPointer: true, pos: &ast.Ident{}},
				},
			},
			expectReports: 0,
		},
		{
			name: "two methods - consistent pointers",
			typeReceivers: map[string][]receiverInfo{
				"MyType": {
					{isPointer: true, pos: &ast.Ident{}},
					{isPointer: true, pos: &ast.Ident{}},
				},
			},
			expectReports: 0,
		},
		{
			name: "two methods - inconsistent",
			typeReceivers: map[string][]receiverInfo{
				"MyType": {
					{isPointer: true, pos: &ast.Ident{}},
					{isPointer: false, pos: &ast.Ident{}},
				},
			},
			expectReports: 1,
		},
		{
			name: "three methods - two inconsistent",
			typeReceivers: map[string][]receiverInfo{
				"MyType": {
					{isPointer: false, pos: &ast.Ident{}},
					{isPointer: true, pos: &ast.Ident{}},
					{isPointer: true, pos: &ast.Ident{}},
				},
			},
			expectReports: 2,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reports := 0
			mockPass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reports++
				},
			}

			checkReceiverConsistency(mockPass, tt.typeReceivers)

			// Verification du nombre de rapports
			if reports != tt.expectReports {
				t.Errorf("checkReceiverConsistency() reported %d issues, expected %d", reports, tt.expectReports)
			}
		})
	}
}

// Test_reportInconsistency tests the private reportInconsistency function.
func Test_reportInconsistency(t *testing.T) {
	tests := []struct {
		name              string
		dominantIsPointer bool
	}{
		{
			name:              "dominant is pointer",
			dominantIsPointer: true,
		},
		{
			name:              "dominant is value",
			dominantIsPointer: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reported := false
			mockPass := &analysis.Pass{
				Report: func(_d analysis.Diagnostic) {
					reported = true
				},
			}

			pos := &ast.Ident{}
			reportInconsistency(mockPass, pos, "MyType", tt.dominantIsPointer)

			// Verification du rapport
			if !reported {
				t.Error("reportInconsistency() did not report")
			}
		})
	}
}

// Test_collectReceivers tests the private collectReceivers function.
func Test_collectReceivers(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectedLen  int
		expectedType string
	}{
		{
			name: "method with pointer receiver",
			code: `package test
func (s *MyStruct) Method() {}`,
			expectedLen:  1,
			expectedType: "MyStruct",
		},
		{
			name: "method with value receiver",
			code: `package test
func (s MyStruct) Method() {}`,
			expectedLen:  1,
			expectedType: "MyStruct",
		},
		{
			name: "regular function - no receiver",
			code: `package test
func RegularFunc() {}`,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Verification de l'erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{
				Fset: fset,
			}
			cfg := config.Get()

			result := collectReceivers(pass, insp, cfg)

			// Verification du nombre de types
			if len(result) != tt.expectedLen {
				t.Errorf("collectReceivers() returned %d types, expected %d", len(result), tt.expectedLen)
			}

			// Verification du type attendu
			if tt.expectedLen > 0 && tt.expectedType != "" {
				if _, ok := result[tt.expectedType]; !ok {
					t.Errorf("collectReceivers() missing expected type %q", tt.expectedType)
				}
			}
		})
	}
}

// Test_runVar021_disabled tests runVar021 with disabled rule.
func Test_runVar021_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-021": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar021(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar021() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar021() reported %d issues, expected 0 when disabled", reportCount)
			}
		})
	}
}

// Test_runVar021_nilInspector tests runVar021 with nil inspector.
func Test_runVar021_nilInspector(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			fset := token.NewFileSet()
			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: nil, // nil inspector
				},
			}

			result, err := runVar021(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar021() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar021() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar021_nilFset tests runVar021 with nil Fset.
func Test_runVar021_nilFset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			code := `package test`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{
				Fset: nil, // nil Fset
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
			}

			result, err := runVar021(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar021() error = %v", err)
			}
			// Result should be nil
			if result != nil {
				t.Errorf("runVar021() = %v, expected nil", result)
			}
		})
	}
}

// Test_runVar021_fileExcluded tests runVar021 with excluded file.
func Test_runVar021_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-021": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse code with inconsistent receivers
			code := `package test
type MyStruct struct{}
func (s *MyStruct) MethodPtr() {}
func (s MyStruct) MethodVal() {}
`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar021(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar021() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar021() reported %d issues, expected 0 when file excluded", reportCount)
			}
		})
	}
}

// Test_collectReceivers_emptyRecvList tests with empty receiver list.
func Test_collectReceivers_emptyRecvList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			config.Reset()

			// Create a function with Recv but empty List
			funcDecl := &ast.FuncDecl{
				Name: &ast.Ident{Name: "Method"},
				Recv: &ast.FieldList{List: []*ast.Field{}}, // Empty list
				Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{}}},
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			}

			file := &ast.File{
				Name:  &ast.Ident{Name: "test"},
				Decls: []ast.Decl{funcDecl},
			}

			fset := token.NewFileSet()
			insp := inspector.New([]*ast.File{file})
			pass := &analysis.Pass{
				Fset: fset,
			}
			cfg := config.Get()

			result := collectReceivers(pass, insp, cfg)

			// Should not have any receivers
			if len(result) != 0 {
				t.Errorf("collectReceivers() returned %d types, expected 0", len(result))
			}
		})
	}
}
