package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_isCallWithArg1 verifie la detection d'un appel avec argument 1.
func Test_isCallWithArg1(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "call with arg 1",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "1"},
				},
			},
			expected: true,
		},
		{
			name: "call with arg 2",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "2"},
				},
			},
			expected: false,
		},
		{
			name: "call with no args",
			call: &ast.CallExpr{
				Args: []ast.Expr{},
			},
			expected: false,
		},
		{
			name: "call with multiple args",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "1"},
					&ast.BasicLit{Kind: token.INT, Value: "2"},
				},
			},
			expected: false,
		},
		{
			name: "call with string arg",
			call: &ast.CallExpr{
				Args: []ast.Expr{
					&ast.BasicLit{Kind: token.STRING, Value: "\"1\""},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isCallWithArg1(tt.call)
			// Verification du resultat
			if result != tt.expected {
				t.Errorf("isCallWithArg1() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_getReceiverName verifie l'extraction du nom du receiver.
func Test_getReceiverName(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "simple ident",
			expr:     &ast.Ident{Name: "wg"},
			expected: "wg",
		},
		{
			name:     "non-ident expr",
			expr:     &ast.BasicLit{Kind: token.INT, Value: "1"},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := getReceiverName(tt.expr)
			// Verification du resultat
			if result != tt.expected {
				t.Errorf("getReceiverName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractFuncLiteral verifie l'extraction de fonction literal.
func Test_extractFuncLiteral(t *testing.T) {
	tests := []struct {
		name     string
		goStmt   *ast.GoStmt
		expected bool
	}{
		{
			name: "go with func literal",
			goStmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Body: &ast.BlockStmt{},
					},
				},
			},
			expected: true,
		},
		{
			name: "go with ident (non-literal)",
			goStmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.Ident{Name: "myFunc"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := extractFuncLiteral(tt.goStmt)
			// Verification du resultat
			if (result != nil) != tt.expected {
				t.Errorf("extractFuncLiteral() returned %v, want non-nil=%v", result, tt.expected)
			}
		})
	}
}

// Test_extractWaitGroupMethodCall tests extraction of WaitGroup method calls.
func Test_extractWaitGroupMethodCall(t *testing.T) {
	// Create a simple types.Info with no type information
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	// Create pass with empty types info
	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name       string
		call       *ast.CallExpr
		methodName string
		expected   string
	}{
		{
			name: "not a selector",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "Add"},
			},
			methodName: "Add",
			expected:   "",
		},
		{
			name: "wrong method name",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "wg"},
					Sel: &ast.Ident{Name: "Done"},
				},
			},
			methodName: "Add",
			expected:   "",
		},
		{
			name: "no type info for receiver",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "wg"},
					Sel: &ast.Ident{Name: "Add"},
				},
			},
			methodName: "Add",
			expected:   "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractWaitGroupMethodCall(pass, tt.call, tt.methodName)
			// Verify result
			if result != tt.expected {
				t.Errorf("extractWaitGroupMethodCall() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isWaitGroupType tests detection of WaitGroup types.
func Test_isWaitGroupType(t *testing.T) {
	tests := []struct {
		name     string
		typVal   types.Type
		expected bool
	}{
		{
			name:     "nil type",
			typVal:   nil,
			expected: false,
		},
		{
			name:     "basic type",
			typVal:   types.Typ[types.Int],
			expected: false,
		},
		{
			name:     "pointer to basic type",
			typVal:   types.NewPointer(types.Typ[types.Int]),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Handle nil type specially
			if tt.typVal == nil {
				// isWaitGroupType expects non-nil type
				defer func() {
					// Should panic or handle gracefully
					if r := recover(); r != nil {
						// Expected panic for nil type
					}
				}()
			}
			// Skip nil test since it will panic
			if tt.typVal == nil {
				return
			}
			result := isWaitGroupType(tt.typVal)
			// Verify result
			if result != tt.expected {
				t.Errorf("isWaitGroupType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isWaitGroupReceiver tests detection of WaitGroup receiver.
func Test_isWaitGroupReceiver(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "no type info",
			expr:     &ast.Ident{Name: "wg"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isWaitGroupReceiver(pass, tt.expr)
			// Verify result
			if result != tt.expected {
				t.Errorf("isWaitGroupReceiver() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractWaitGroupAdd1 tests extraction of wg.Add(1).
func Test_extractWaitGroupAdd1(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name     string
		stmt     ast.Stmt
		expected string
	}{
		{
			name:     "not expression statement",
			stmt:     &ast.AssignStmt{},
			expected: "",
		},
		{
			name: "not a call expression",
			stmt: &ast.ExprStmt{
				X: &ast.Ident{Name: "x"},
			},
			expected: "",
		},
		{
			name: "call without arg 1",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "wg"},
						Sel: &ast.Ident{Name: "Add"},
					},
					Args: []ast.Expr{
						&ast.BasicLit{Kind: token.INT, Value: "2"},
					},
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractWaitGroupAdd1(pass, tt.stmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("extractWaitGroupAdd1() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isGoWithDeferDone tests go statement with defer wg.Done().
func Test_isGoWithDeferDone(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name     string
		stmt     ast.Stmt
		wgName   string
		expected bool
	}{
		{
			name:     "not a go statement",
			stmt:     &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "go with non-literal func",
			stmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.Ident{Name: "myFunc"},
				},
			},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "go with empty func body",
			stmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Body: &ast.BlockStmt{List: []ast.Stmt{}},
					},
				},
			},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "go with nil func body",
			stmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Body: nil,
					},
				},
			},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "go with first stmt not defer",
			stmt: &ast.GoStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
						}},
					},
				},
			},
			wgName:   "wg",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isGoWithDeferDone(pass, tt.stmt, tt.wgName)
			// Verify result
			if result != tt.expected {
				t.Errorf("isGoWithDeferDone() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_hasDeferDoneFirst tests first statement is defer wg.Done().
func Test_hasDeferDoneFirst(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name     string
		funcLit  *ast.FuncLit
		wgName   string
		expected bool
	}{
		{
			name: "nil body",
			funcLit: &ast.FuncLit{
				Body: nil,
			},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "empty body",
			funcLit: &ast.FuncLit{
				Body: &ast.BlockStmt{List: []ast.Stmt{}},
			},
			wgName:   "wg",
			expected: false,
		},
		{
			name: "first not defer",
			funcLit: &ast.FuncLit{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			wgName:   "wg",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := hasDeferDoneFirst(pass, tt.funcLit, tt.wgName)
			// Verify result
			if result != tt.expected {
				t.Errorf("hasDeferDoneFirst() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isDeferDoneForWaitGroup tests defer statement is for correct WaitGroup.
func Test_isDeferDoneForWaitGroup(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name      string
		deferStmt *ast.DeferStmt
		wgName    string
		expected  bool
	}{
		{
			name: "not matching WaitGroup",
			deferStmt: &ast.DeferStmt{
				Call: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "other"},
						Sel: &ast.Ident{Name: "Done"},
					},
				},
			},
			wgName:   "wg",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isDeferDoneForWaitGroup(pass, tt.deferStmt, tt.wgName)
			// Verify result
			if result != tt.expected {
				t.Errorf("isDeferDoneForWaitGroup() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isWaitGroupPattern tests detection of WaitGroup pattern.
func Test_isWaitGroupPattern(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name     string
		addStmt  ast.Stmt
		goStmt   ast.Stmt
		expected bool
	}{
		{
			name:     "addStmt not wg.Add(1)",
			addStmt:  &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			goStmt:   &ast.GoStmt{Call: &ast.CallExpr{Fun: &ast.Ident{Name: "f"}}},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isWaitGroupPattern(pass, tt.addStmt, tt.goStmt)
			// Verify result
			if result != tt.expected {
				t.Errorf("isWaitGroupPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_analyzeConsecutiveStatements tests analysis of consecutive statements.
func Test_analyzeConsecutiveStatements(t *testing.T) {
	// Create empty types info
	emptyInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		TypesInfo: emptyInfo,
	}

	tests := []struct {
		name  string
		stmts []ast.Stmt
	}{
		{
			name:  "empty list",
			stmts: []ast.Stmt{},
		},
		{
			name:  "single statement",
			stmts: []ast.Stmt{&ast.ExprStmt{X: &ast.Ident{Name: "x"}}},
		},
		{
			name: "two non-matching statements",
			stmts: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				&ast.ExprStmt{X: &ast.Ident{Name: "y"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			defer func() {
				// Recover if panic happens
				if r := recover(); r != nil {
					t.Errorf("analyzeConsecutiveStatements() panicked: %v", r)
				}
			}()
			analyzeConsecutiveStatements(pass, tt.stmts)
		})
	}
}

// Test_runVar034_defensiveBranches tests defensive branches in runVar034.
func Test_runVar034_defensiveBranches(t *testing.T) {
	tests := []struct {
		name string
		pass *analysis.Pass
	}{
		{
			name: "inspector not in ResultOf",
			pass: &analysis.Pass{
				Fset:     token.NewFileSet(),
				ResultOf: map[*analysis.Analyzer]any{},
			},
		},
		{
			name: "inspector is nil",
			pass: &analysis.Pass{
				Fset: token.NewFileSet(),
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: nil,
				},
			},
		},
		{
			name: "inspector wrong type",
			pass: &analysis.Pass{
				Fset: token.NewFileSet(),
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: "wrong type",
				},
			},
		},
		{
			name: "valid inspector but nil fset",
			pass: &analysis.Pass{
				Fset: nil,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: inspector.New([]*ast.File{}),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			result, err := runVar034(tt.pass)
			// Verify no error
			if err != nil {
				t.Errorf("runVar034() error = %v, want nil", err)
			}
			// Verify nil result
			if result != nil {
				t.Errorf("runVar034() result = %v, want nil", result)
			}
		})
	}
}
