package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_isVar037AppendCall tests detection of append calls.
func Test_isVar037AppendCall(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "not identifier",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "pkg"},
					Sel: &ast.Ident{Name: "append"},
				},
			},
			expected: false,
		},
		{
			name: "wrong function name",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "prepend"},
			},
			expected: false,
		},
		{
			name: "valid append call",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "append"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isVar037AppendCall(tt.call)
			// Verify result
			if result != tt.expected {
				t.Errorf("isVar037AppendCall() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isVar037BlankOrNil tests detection of blank or nil.
func Test_isVar037BlankOrNil(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "nil expression",
			expr:     nil,
			expected: true,
		},
		{
			name:     "blank identifier",
			expr:     &ast.Ident{Name: "_"},
			expected: true,
		},
		{
			name:     "normal identifier",
			expr:     &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name:     "non-identifier",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "1"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isVar037BlankOrNil(tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isVar037BlankOrNil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isSameIdent tests detection of same identifiers.
func Test_isSameIdent(t *testing.T) {
	tests := []struct {
		name     string
		expr1    ast.Expr
		expr2    ast.Expr
		expected bool
	}{
		{
			name:     "first not identifier",
			expr1:    &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			expr2:    &ast.Ident{Name: "x"},
			expected: false,
		},
		{
			name:     "second not identifier",
			expr1:    &ast.Ident{Name: "x"},
			expr2:    &ast.CallExpr{Fun: &ast.Ident{Name: "f"}},
			expected: false,
		},
		{
			name:     "different names",
			expr1:    &ast.Ident{Name: "x"},
			expr2:    &ast.Ident{Name: "y"},
			expected: false,
		},
		{
			name:     "same names",
			expr1:    &ast.Ident{Name: "x"},
			expr2:    &ast.Ident{Name: "x"},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isSameIdent(tt.expr1, tt.expr2)
			// Verify result
			if result != tt.expected {
				t.Errorf("isSameIdent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_validateRangeBody037 tests validation of range body.
func Test_validateRangeBody037(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		wantNil   bool
	}{
		{
			name: "nil body",
			rangeStmt: &ast.RangeStmt{
				Body: nil,
			},
			wantNil: true,
		},
		{
			name: "empty body",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
			wantNil: true,
		},
		{
			name: "multiple statements",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{},
					&ast.ExprStmt{},
				}},
			},
			wantNil: true,
		},
		{
			name: "single statement",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := validateRangeBody037(tt.rangeStmt)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("validateRangeBody037() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_extractAppendCall037 tests extraction of append call.
func Test_extractAppendCall037(t *testing.T) {
	tests := []struct {
		name    string
		stmt    ast.Stmt
		wantNil bool
	}{
		{
			name:    "not assignment",
			stmt:    &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			wantNil: true,
		},
		{
			name: "multiple RHS",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
			},
			wantNil: true,
		},
		{
			name: "multiple LHS",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}, &ast.Ident{Name: "y"}},
				Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "append"}}},
			},
			wantNil: true,
		},
		{
			name: "RHS not call",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
			},
			wantNil: true,
		},
		{
			name: "not append call",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.CallExpr{
					Fun:  &ast.Ident{Name: "other"},
					Args: []ast.Expr{&ast.Ident{Name: "x"}, &ast.Ident{Name: "v"}},
				}},
			},
			wantNil: true,
		},
		{
			name: "append with wrong arg count",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.CallExpr{
					Fun:  &ast.Ident{Name: "append"},
					Args: []ast.Expr{&ast.Ident{Name: "x"}},
				}},
			},
			wantNil: true,
		},
		{
			name: "valid append call",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.CallExpr{
					Fun:  &ast.Ident{Name: "append"},
					Args: []ast.Expr{&ast.Ident{Name: "x"}, &ast.Ident{Name: "v"}},
				}},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractAppendCall037(tt.stmt)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("extractAppendCall037() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_isKeyCollection tests detection of key collection.
func Test_isKeyCollection(t *testing.T) {
	tests := []struct {
		name          string
		rangeStmt     *ast.RangeStmt
		appendedValue ast.Expr
		expected      bool
	}{
		{
			name: "no key",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "_"},
			},
			appendedValue: &ast.Ident{Name: "k"},
			expected:      false,
		},
		{
			name: "value is used",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "v"},
			},
			appendedValue: &ast.Ident{Name: "k"},
			expected:      false,
		},
		{
			name: "appended value not key",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "_"},
			},
			appendedValue: &ast.Ident{Name: "other"},
			expected:      false,
		},
		{
			name: "valid key collection",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: nil,
			},
			appendedValue: &ast.Ident{Name: "k"},
			expected:      true,
		},
		{
			name: "valid key collection with blank value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "_"},
			},
			appendedValue: &ast.Ident{Name: "k"},
			expected:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isKeyCollection(tt.rangeStmt, tt.appendedValue)
			// Verify result
			if result != tt.expected {
				t.Errorf("isKeyCollection() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isValueCollection tests detection of value collection.
func Test_isValueCollection(t *testing.T) {
	tests := []struct {
		name          string
		rangeStmt     *ast.RangeStmt
		appendedValue ast.Expr
		expected      bool
	}{
		{
			name: "no value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: nil,
			},
			appendedValue: &ast.Ident{Name: "v"},
			expected:      false,
		},
		{
			name: "key is used",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "v"},
			},
			appendedValue: &ast.Ident{Name: "v"},
			expected:      false,
		},
		{
			name: "appended value not range value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
			},
			appendedValue: &ast.Ident{Name: "other"},
			expected:      false,
		},
		{
			name: "valid value collection",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			appendedValue: &ast.Ident{Name: "v"},
			expected:      true,
		},
		{
			name: "valid value collection with blank key",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
			},
			appendedValue: &ast.Ident{Name: "v"},
			expected:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isValueCollection(tt.rangeStmt, tt.appendedValue)
			// Verify result
			if result != tt.expected {
				t.Errorf("isValueCollection() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_detectCollectionType tests detection of collection type.
func Test_detectCollectionType(t *testing.T) {
	tests := []struct {
		name       string
		rangeStmt  *ast.RangeStmt
		appendCall *ast.CallExpr
		expected   string
	}{
		{
			name: "keys collection",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: nil,
			},
			appendCall: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "append"},
				Args: []ast.Expr{&ast.Ident{Name: "result"}, &ast.Ident{Name: "k"}},
			},
			expected: collectionTypeKeys,
		},
		{
			name: "values collection",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "_"},
				Value: &ast.Ident{Name: "v"},
			},
			appendCall: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "append"},
				Args: []ast.Expr{&ast.Ident{Name: "result"}, &ast.Ident{Name: "v"}},
			},
			expected: collectionTypeValues,
		},
		{
			name: "neither keys nor values",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "v"},
			},
			appendCall: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "append"},
				Args: []ast.Expr{&ast.Ident{Name: "result"}, &ast.Ident{Name: "other"}},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := detectCollectionType(tt.rangeStmt, tt.appendCall)
			// Verify result
			if result != tt.expected {
				t.Errorf("detectCollectionType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_detectMapCollectionPattern tests detection of map collection pattern.
func Test_detectMapCollectionPattern(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		expected  string
	}{
		{
			name: "nil body",
			rangeStmt: &ast.RangeStmt{
				Body: nil,
			},
			expected: "",
		},
		{
			name: "empty body",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
			expected: "",
		},
		{
			name: "not assignment in body",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			expected: "",
		},
		{
			name: "valid keys collection",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: nil,
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{&ast.Ident{Name: "result"}},
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun:  &ast.Ident{Name: "append"},
							Args: []ast.Expr{&ast.Ident{Name: "result"}, &ast.Ident{Name: "k"}},
						}},
					},
				}},
			},
			expected: collectionTypeKeys,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := detectMapCollectionPattern(tt.rangeStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("detectMapCollectionPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isRangingOverMap tests detection of ranging over map.
func Test_isRangingOverMap(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		expected  bool
	}{
		{
			name: "no type info",
			rangeStmt: &ast.RangeStmt{
				X: &ast.Ident{Name: "m"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isRangingOverMap(pass, tt.rangeStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("isRangingOverMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isRangingOverMap_withMapType tests isRangingOverMap with actual map type.
func Test_isRangingOverMap_withMapType(t *testing.T) {
	// Parse code with map range
	src := `package test

func foo() {
	m := map[string]int{"a": 1}
	for k := range m {
		_ = k
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Pkg:       pkg,
		TypesInfo: info,
	}

	// Find the range statement
	var rangeStmt *ast.RangeStmt
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for range statement
		if rs, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = rs
			return false
		}
		return true
	})

	// Check that rangeStmt was found
	if rangeStmt == nil {
		t.Fatal("Failed to find range statement")
	}

	// Test that it detects map type
	result := isRangingOverMap(pass, rangeStmt)
	// Verify result
	if !result {
		t.Error("isRangingOverMap() = false, want true for map range")
	}
}

// Test_isRangingOverMap_withSliceType tests isRangingOverMap with slice type.
func Test_isRangingOverMap_withSliceType(t *testing.T) {
	// Parse code with slice range
	src := `package test

func foo() {
	s := []int{1, 2, 3}
	for i := range s {
		_ = i
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Pkg:       pkg,
		TypesInfo: info,
	}

	// Find the range statement
	var rangeStmt *ast.RangeStmt
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for range statement
		if rs, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = rs
			return false
		}
		return true
	})

	// Check that rangeStmt was found
	if rangeStmt == nil {
		t.Fatal("Failed to find range statement")
	}

	// Test that it does not detect slice as map
	result := isRangingOverMap(pass, rangeStmt)
	// Verify result
	if result {
		t.Error("isRangingOverMap() = true, want false for slice range")
	}
}

// Test_runVar037_ruleDisabled tests runVar037 when rule is disabled.
func Test_runVar037_ruleDisabled(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar037]

	// Disable the rule
	cfg.Rules[ruleCodeVar037] = &config.RuleConfig{Enabled: config.Bool(false)}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar037)
		} else {
			cfg.Rules[ruleCodeVar037] = originalRule
		}
	}()

	// Parse code with map collection pattern
	src := `package test

func getKeys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  Analyzer037,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar037(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar037() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when disabled
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics when rule disabled, got %d", len(diagnostics))
	}
}

// Test_runVar037_fileExcluded tests runVar037 when file is excluded.
func Test_runVar037_fileExcluded(t *testing.T) {
	// Save the current config
	cfg := config.Get()
	// Initialize rules map if needed
	if cfg.Rules == nil {
		cfg.Rules = make(map[string]*config.RuleConfig)
	}
	// Save original state
	originalRule := cfg.Rules[ruleCodeVar037]

	// Set up rule with file exclusion
	cfg.Rules[ruleCodeVar037] = &config.RuleConfig{
		Exclude: []string{"excluded.go"},
	}
	// Ensure restoration at the end
	defer func() {
		// Restore original state
		if originalRule == nil {
			delete(cfg.Rules, ruleCodeVar037)
		} else {
			cfg.Rules[ruleCodeVar037] = originalRule
		}
	}()

	// Parse code with map collection pattern
	src := `package test

func getKeys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  Analyzer037,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar037(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar037() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when file excluded
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics when file excluded, got %d", len(diagnostics))
	}
}

// Test_runVar037_notRangingOverMap tests runVar037 with range over slice.
func Test_runVar037_notRangingOverMap(t *testing.T) {
	// Reset config
	config.Reset()

	// Parse code with slice range (not map)
	src := `package test

func foo() {
	s := []int{1, 2, 3}
	var result []int
	for _, v := range s {
		result = append(result, v)
	}
	_ = result
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  Analyzer037,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar037(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar037() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports for slice range
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics for slice range, got %d", len(diagnostics))
	}
}

// Test_runVar037_noPattern tests runVar037 with map range but no collection pattern.
func Test_runVar037_noPattern(t *testing.T) {
	// Reset config
	config.Reset()

	// Parse code with map range but not a simple collection pattern
	src := `package test

func foo() {
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		// Not a simple collection - uses both k and v
		println(k, v)
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  Analyzer037,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar037(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar037() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when no collection pattern
	if len(diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics when no pattern, got %d", len(diagnostics))
	}
}

// Test_runVar037_withReport tests runVar037 detects map collection pattern.
func Test_runVar037_withReport(t *testing.T) {
	// Reset config
	config.Reset()

	// Parse code with map keys collection pattern
	src := `package test

func getKeys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check for parsing errors
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Type check the file
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{}
	pkg, err := conf.Check("test", fset, []*ast.File{file}, info)
	// Check for type errors
	if err != nil {
		t.Fatalf("Failed to type check: %v", err)
	}

	// Create pass
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  Analyzer037,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	// Run the analyzer
	result, err := runVar037(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar037() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check that diagnostic was reported
	if len(diagnostics) != 1 {
		t.Errorf("expected 1 diagnostic for map keys collection, got %d", len(diagnostics))
	}
}
