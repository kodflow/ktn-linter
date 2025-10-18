package ktn_pool

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test trackPoolGetAssignment avec différents cas
func TestTrackPoolGetAssignment(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantVars int // nombre de variables trackées
	}{
		{
			name: "single assignment - pool.Get()",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get()
}`,
			wantVars: 1,
		},
		{
			name: "multiple LHS - should skip",
			code: `package test
var pool sync.Pool
func f() {
	a, b := pool.Get(), pool.Get()
}`,
			wantVars: 0,
		},
		{
			name: "multiple RHS - should skip",
			code: `package test
var pool sync.Pool
func f() {
	buf := func() interface{} { return pool.Get() }()
}`,
			wantVars: 0,
		},
		{
			name: "type assertion",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get().([]byte)
}`,
			wantVars: 1,
		},
		{
			name: "not pool.Get()",
			code: `package test
func f() {
	x := otherFunc()
}
func otherFunc() int { return 0 }`,
			wantVars: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			poolVars := make(map[string]ast.Expr)
			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: nil, // Test sans info de type
			}

			ast.Inspect(file, func(n ast.Node) bool {
				if stmt, ok := n.(*ast.AssignStmt); ok {
					trackPoolGetAssignment(stmt, poolVars, pass)
				}
				return true
			})

			if len(poolVars) != tt.wantVars {
				t.Errorf("got %d poolVars, want %d", len(poolVars), tt.wantVars)
			}
		})
	}
}

// Test trackDeferPut avec différents cas
func TestTrackDeferPut(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantDefer int // nombre de defer Put trackés
	}{
		{
			name: "valid defer pool.Put()",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get()
	defer pool.Put(buf)
}`,
			wantDefer: 1,
		},
		{
			name: "defer without call - should skip",
			code: `package test
func f() {
	defer func() {}()
}`,
			wantDefer: 0,
		},
		{
			name: "defer other method",
			code: `package test
func f() {
	defer cleanup()
}
func cleanup() {}`,
			wantDefer: 0,
		},
		{
			name: "defer pool.Put() without args",
			code: `package test
var pool sync.Pool
func f() {
	defer pool.Put()
}`,
			wantDefer: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			deferredPuts := make(map[string]bool)

			ast.Inspect(file, func(n ast.Node) bool {
				if stmt, ok := n.(*ast.DeferStmt); ok {
					trackDeferPut(stmt, deferredPuts)
				}
				return true
			})

			if len(deferredPuts) != tt.wantDefer {
				t.Errorf("got %d deferredPuts, want %d", len(deferredPuts), tt.wantDefer)
			}
		})
	}
}

// Test trackDeferPut avec DeferStmt ayant Call nil (cas edge case)
func TestTrackDeferPut_NilCall(t *testing.T) {
	// Créer manuellement un DeferStmt avec Call nil
	deferredPuts := make(map[string]bool)
	stmt := &ast.DeferStmt{
		Call: nil, // Call explicitement nil
	}

	// Ne devrait pas paniquer et ne rien ajouter
	trackDeferPut(stmt, deferredPuts)

	if len(deferredPuts) != 0 {
		t.Errorf("expected 0 deferredPuts with nil Call, got %d", len(deferredPuts))
	}
}

// Test isPoolGetCall avec différents cas
func TestIsPoolGetCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "pool.Get() call",
			code: `package test
var pool sync.Pool
func f() {
	_ = pool.Get()
}`,
			want: true,
		},
		{
			name: "not a call expression",
			code: `package test
func f() {
	x := 42
}`,
			want: false,
		},
		{
			name: "not a selector expression",
			code: `package test
func f() {
	_ = someFunc()
}
func someFunc() int { return 0 }`,
			want: false,
		},
		{
			name: "wrong method name",
			code: `package test
var pool sync.Pool
func f() {
	_ = pool.Put(nil)
}`,
			want: false,
		},
		{
			name: "non-pool variable",
			code: `package test
type Custom struct{}
func (c Custom) Get() interface{} { return nil }
var obj Custom
func f() {
	_ = obj.Get()
}`,
			want: false,
		},
		{
			name: "chained method call - not ident",
			code: `package test
type Container struct{}
func (c Container) GetPool() interface{} { return nil }
func (c Container) Get() interface{} { return nil }
var container Container
func f() {
	_ = container.GetPool().Get()
}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			pass := &analysis.Pass{
				Fset:      fset,
				TypesInfo: nil,
			}

			var result bool
			ast.Inspect(file, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if isPoolGetCall(call, pass) {
						result = true
					}
				}
				return true
			})

			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

// Test isPoolGetCall avec selExpr.X qui n'est pas un *ast.Ident
func TestIsPoolGetCall_NonIdentSelector(t *testing.T) {
	// Créer manuellement un CallExpr avec un SelectorExpr dont X n'est pas un Ident
	// Par exemple: something().Get() où X est un *ast.CallExpr
	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.CallExpr{ // Pas un *ast.Ident
				Fun: &ast.Ident{Name: "getPool"},
			},
			Sel: &ast.Ident{Name: "Get"},
		},
	}

	pass := &analysis.Pass{
		TypesInfo: nil,
	}

	result := isPoolGetCall(callExpr, pass)
	if result != false {
		t.Errorf("expected false for non-ident selector, got %v", result)
	}
}

// Test isPoolPutCall avec différents cas
func TestIsPoolPutCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "pool.Put() call",
			code: `package test
var pool sync.Pool
func f() {
	pool.Put(nil)
}`,
			want: true,
		},
		{
			name: "not a selector expression",
			code: `package test
func f() {
	someFunc()
}
func someFunc() {}`,
			want: false,
		},
		{
			name: "wrong method name",
			code: `package test
var pool sync.Pool
func f() {
	_ = pool.Get()
}`,
			want: false,
		},
		{
			name: "non-pool variable",
			code: `package test
type Custom struct{}
func (c Custom) Put(v interface{}) {}
var obj Custom
func f() {
	obj.Put(nil)
}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var result bool
			ast.Inspect(file, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if isPoolPutCall(call) {
						result = true
					}
				}
				return true
			})

			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

// Test isPoolPutCall avec selExpr.X qui n'est pas un *ast.Ident
func TestIsPoolPutCall_NonIdentSelector(t *testing.T) {
	// Créer manuellement un CallExpr avec un SelectorExpr dont X n'est pas un Ident
	// Par exemple: something().Put() où X est un *ast.CallExpr
	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.CallExpr{ // Pas un *ast.Ident
				Fun: &ast.Ident{Name: "getPool"},
			},
			Sel: &ast.Ident{Name: "Put"},
		},
	}

	result := isPoolPutCall(callExpr)
	if result != false {
		t.Errorf("expected false for non-ident selector, got %v", result)
	}
}

// Test extractVarName avec différents cas
func TestExtractVarName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		testLHS  bool
		testRHS  bool
		wantLHS  string
		wantRHS  string
	}{
		{
			name: "simple identifier on LHS",
			code: `package test
func f() {
	buf := 42
}`,
			testLHS: true,
			wantLHS: "buf",
		},
		{
			name: "call expression on RHS - should return empty",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get()
}`,
			testRHS: true,
			wantRHS: "",
		},
		{
			name: "type assertion on RHS - should return empty",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get().([]byte)
}`,
			testRHS: true,
			wantRHS: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			ast.Inspect(file, func(n ast.Node) bool {
				if stmt, ok := n.(*ast.AssignStmt); ok {
					if tt.testLHS && len(stmt.Lhs) > 0 {
						result := extractVarName(stmt.Lhs[0])
						if result != tt.wantLHS {
							t.Errorf("LHS: got %q, want %q", result, tt.wantLHS)
						}
						return false
					}
					if tt.testRHS && len(stmt.Rhs) > 0 {
						result := extractVarName(stmt.Rhs[0])
						if result != tt.wantRHS {
							t.Errorf("RHS: got %q, want %q", result, tt.wantRHS)
						}
						return false
					}
				}
				return true
			})
		})
	}
}

// Test unwrapTypeAssertion
func TestUnwrapTypeAssertion(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantUnwrap bool
	}{
		{
			name: "with type assertion",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get().([]byte)
}`,
			wantUnwrap: true,
		},
		{
			name: "without type assertion",
			code: `package test
var pool sync.Pool
func f() {
	buf := pool.Get()
}`,
			wantUnwrap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			var hasTypeAssert bool
			ast.Inspect(file, func(n ast.Node) bool {
				if stmt, ok := n.(*ast.AssignStmt); ok {
					if len(stmt.Rhs) > 0 {
						unwrapped := unwrapTypeAssertion(stmt.Rhs[0])
						_, isTypeAssert := stmt.Rhs[0].(*ast.TypeAssertExpr)
						hasTypeAssert = isTypeAssert && unwrapped != stmt.Rhs[0]
						return false
					}
				}
				return true
			})

			if hasTypeAssert != tt.wantUnwrap {
				t.Errorf("got unwrap=%v, want unwrap=%v", hasTypeAssert, tt.wantUnwrap)
			}
		})
	}
}
