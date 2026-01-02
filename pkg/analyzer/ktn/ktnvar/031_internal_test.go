package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_extractRangeKeyValue031 tests extraction of key/value from range.
func Test_extractRangeKeyValue031(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		wantKey   bool
		wantVal   bool
	}{
		{
			name: "both key and value identifiers",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "v"},
			},
			wantKey: true,
			wantVal: true,
		},
		{
			name: "nil key",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			wantKey: false,
			wantVal: false,
		},
		{
			name: "nil value",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: nil,
			},
			wantKey: false,
			wantVal: false,
		},
		{
			name: "key not identifier",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.BasicLit{Kind: token.INT, Value: "0"},
				Value: &ast.Ident{Name: "v"},
			},
			wantKey: false,
			wantVal: false,
		},
		{
			name: "value not identifier",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.BasicLit{Kind: token.INT, Value: "0"},
			},
			wantKey: false,
			wantVal: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			key, val := extractRangeKeyValue031(tt.rangeStmt)
			gotKey := key != nil
			gotVal := val != nil
			// Verify key result
			if gotKey != tt.wantKey {
				t.Errorf("extractRangeKeyValue031() key = %v, want %v", gotKey, tt.wantKey)
			}
			// Verify value result
			if gotVal != tt.wantVal {
				t.Errorf("extractRangeKeyValue031() val = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

// Test_validateIndexExpr031 tests validation of index expressions.
func Test_validateIndexExpr031(t *testing.T) {
	tests := []struct {
		name      string
		indexExpr *ast.IndexExpr
		keyName   string
		mapMakes  map[string]token.Pos
		expected  bool
	}{
		{
			name: "valid index expression",
			indexExpr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "clone"},
				Index: &ast.Ident{Name: "k"},
			},
			keyName:  "k",
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: true,
		},
		{
			name: "X not identifier",
			indexExpr: &ast.IndexExpr{
				X:     &ast.BasicLit{Kind: token.STRING, Value: "test"},
				Index: &ast.Ident{Name: "k"},
			},
			keyName:  "k",
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "clone not in mapMakes",
			indexExpr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "notClone"},
				Index: &ast.Ident{Name: "k"},
			},
			keyName:  "k",
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "index not identifier",
			indexExpr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "clone"},
				Index: &ast.BasicLit{Kind: token.INT, Value: "0"},
			},
			keyName:  "k",
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "index name mismatch",
			indexExpr: &ast.IndexExpr{
				X:     &ast.Ident{Name: "clone"},
				Index: &ast.Ident{Name: "other"},
			},
			keyName:  "k",
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := validateIndexExpr031(tt.indexExpr, tt.keyName, tt.mapMakes)
			// Verify result
			if result != tt.expected {
				t.Errorf("validateIndexExpr031() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractMakeMapAssign tests extraction of make(map) assignments.
func Test_extractMakeMapAssign(t *testing.T) {
	tests := []struct {
		name         string
		stmt         ast.Stmt
		expectedName string
	}{
		{
			name:         "not an assignment",
			stmt:         &ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			expectedName: "",
		},
		{
			name: "multiple LHS",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
				Rhs: []ast.Expr{&ast.CallExpr{}},
			},
			expectedName: "",
		},
		{
			name: "LHS not identifier",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
				Rhs: []ast.Expr{&ast.CallExpr{}},
			},
			expectedName: "",
		},
		{
			name: "RHS not call",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
			},
			expectedName: "",
		},
		{
			name: "not a make call",
			stmt: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.CallExpr{
					Fun:  &ast.Ident{Name: "someFunc"},
					Args: []ast.Expr{},
				}},
			},
			expectedName: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			name, _ := extractMakeMapAssign(tt.stmt)
			// Verify result
			if name != tt.expectedName {
				t.Errorf("extractMakeMapAssign() = %v, want %v", name, tt.expectedName)
			}
		})
	}
}

// Test_isMakeMapCallExpr tests detection of make(map) calls.
func Test_isMakeMapCallExpr(t *testing.T) {
	tests := []struct {
		name     string
		call     *ast.CallExpr
		expected bool
	}{
		{
			name: "valid make(map)",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{&ast.MapType{}},
			},
			expected: true,
		},
		{
			name: "not make",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "other"},
				Args: []ast.Expr{&ast.MapType{}},
			},
			expected: false,
		},
		{
			name: "fun not identifier",
			call: &ast.CallExpr{
				Fun:  &ast.SelectorExpr{X: &ast.Ident{Name: "pkg"}, Sel: &ast.Ident{Name: "make"}},
				Args: []ast.Expr{&ast.MapType{}},
			},
			expected: false,
		},
		{
			name: "no args",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{},
			},
			expected: false,
		},
		{
			name: "first arg not map type",
			call: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "make"},
				Args: []ast.Expr{&ast.ArrayType{}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isMakeMapCallExpr(tt.call)
			// Verify result
			if result != tt.expected {
				t.Errorf("isMakeMapCallExpr() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractRangeBodyAssign031 tests extraction of assignment from range body.
func Test_extractRangeBodyAssign031(t *testing.T) {
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
					&ast.AssignStmt{},
					&ast.ExprStmt{},
				}},
			},
			wantNil: true,
		},
		{
			name: "not assignment statement",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
				}},
			},
			wantNil: true,
		},
		{
			name: "assignment with wrong token",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.AssignStmt{
						Tok: token.DEFINE,
						Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
						Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
					},
				}},
			},
			wantNil: true,
		},
		{
			name: "multiple LHS",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.AssignStmt{
						Tok: token.ASSIGN,
						Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
						Rhs: []ast.Expr{&ast.Ident{Name: "x"}},
					},
				}},
			},
			wantNil: true,
		},
		{
			name: "valid assignment",
			rangeStmt: &ast.RangeStmt{
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.AssignStmt{
						Tok: token.ASSIGN,
						Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
						Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
					},
				}},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := extractRangeBodyAssign031(tt.rangeStmt)
			gotNil := result == nil
			// Verify result
			if gotNil != tt.wantNil {
				t.Errorf("extractRangeBodyAssign031() = nil? %v, want nil? %v", gotNil, tt.wantNil)
			}
		})
	}
}

// Test_validateCloneAssignment031 tests validation of clone assignment.
func Test_validateCloneAssignment031(t *testing.T) {
	tests := []struct {
		name     string
		assign   *ast.AssignStmt
		keyIdent *ast.Ident
		valIdent *ast.Ident
		mapMakes map[string]token.Pos
		expected bool
	}{
		{
			name: "LHS not index expression",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
			},
			keyIdent: &ast.Ident{Name: "k"},
			valIdent: &ast.Ident{Name: "v"},
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "RHS not identifier",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.IndexExpr{
					X:     &ast.Ident{Name: "clone"},
					Index: &ast.Ident{Name: "k"},
				}},
				Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}},
			},
			keyIdent: &ast.Ident{Name: "k"},
			valIdent: &ast.Ident{Name: "v"},
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "RHS wrong name",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.IndexExpr{
					X:     &ast.Ident{Name: "clone"},
					Index: &ast.Ident{Name: "k"},
				}},
				Rhs: []ast.Expr{&ast.Ident{Name: "other"}},
			},
			keyIdent: &ast.Ident{Name: "k"},
			valIdent: &ast.Ident{Name: "v"},
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: false,
		},
		{
			name: "valid clone assignment",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.IndexExpr{
					X:     &ast.Ident{Name: "clone"},
					Index: &ast.Ident{Name: "k"},
				}},
				Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
			},
			keyIdent: &ast.Ident{Name: "k"},
			valIdent: &ast.Ident{Name: "v"},
			mapMakes: map[string]token.Pos{"clone": 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := validateCloneAssignment031(tt.assign, tt.keyIdent, tt.valIdent, tt.mapMakes)
			// Verify result
			if result != tt.expected {
				t.Errorf("validateCloneAssignment031() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_isSimpleMapClone tests detection of simple map clone patterns.
func Test_isSimpleMapClone(t *testing.T) {
	tests := []struct {
		name      string
		rangeStmt *ast.RangeStmt
		mapMakes  map[string]token.Pos
		stmtIndex int
		stmts     []ast.Stmt
		expected  bool
	}{
		{
			name: "nil key",
			rangeStmt: &ast.RangeStmt{
				Key:   nil,
				Value: &ast.Ident{Name: "v"},
			},
			mapMakes:  map[string]token.Pos{},
			stmtIndex: 0,
			stmts:     []ast.Stmt{},
			expected:  false,
		},
		{
			name: "nil body",
			rangeStmt: &ast.RangeStmt{
				Key:   &ast.Ident{Name: "k"},
				Value: &ast.Ident{Name: "v"},
				Body:  nil,
			},
			mapMakes:  map[string]token.Pos{},
			stmtIndex: 0,
			stmts:     []ast.Stmt{},
			expected:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := isSimpleMapClone(tt.rangeStmt, tt.mapMakes, tt.stmtIndex, tt.stmts)
			// Verify result
			if result != tt.expected {
				t.Errorf("isSimpleMapClone() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_analyzeBlockForMapClone tests analysis of block for map clone.
func Test_analyzeBlockForMapClone(t *testing.T) {
	// Create a minimal pass for testing
	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Fset: fset,
	}

	tests := []struct {
		name  string
		block *ast.BlockStmt
	}{
		{
			name:  "empty block",
			block: &ast.BlockStmt{List: []ast.Stmt{}},
		},
		{
			name: "non-range statement only",
			block: &ast.BlockStmt{List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "x"}},
			}},
		},
		{
			name: "assignment not make(map)",
			block: &ast.BlockStmt{List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
					Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
				},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Just verify no panic
			defer func() {
				// Recover if panic happens
				if r := recover(); r != nil {
					t.Errorf("analyzeBlockForMapClone() panicked: %v", r)
				}
			}()
			analyzeBlockForMapClone(pass, tt.block)
		})
	}
}
