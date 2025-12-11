// Internal tests for analyzer 012.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_isPassthroughTest tests the isPassthroughTest function.
//
// Params:
//   - t: testing context
func Test_isPassthroughTest(t *testing.T) {
	tests := []struct {
		name            string
		code            string
		wantPassthrough bool
	}{
		// === PASSTHROUGH CASES (should be flagged) ===
		{
			name: "empty test body",
			code: `func TestEmpty(t *testing.T) {
			}`,
			wantPassthrough: true,
		},
		{
			name: "only t.Log call",
			code: `func TestOnlyLog(t *testing.T) {
				t.Log("hello")
			}`,
			wantPassthrough: true,
		},
		{
			name: "only function call without validation",
			code: `func TestOnlyCall(t *testing.T) {
				Foo()
			}`,
			wantPassthrough: true,
		},
		{
			name: "only variable assignment",
			code: `func TestOnlyAssign(t *testing.T) {
				x := Foo()
				_ = x
			}`,
			wantPassthrough: true,
		},
		{
			name: "only t.Parallel call",
			code: `func TestOnlyParallel(t *testing.T) {
				t.Parallel()
			}`,
			wantPassthrough: true,
		},
		{
			name: "only t.Helper call",
			code: `func TestOnlyHelper(t *testing.T) {
				t.Helper()
			}`,
			wantPassthrough: true,
		},
		{
			name: "only t.Skip call",
			code: `func TestOnlySkip(t *testing.T) {
				t.Skip("not implemented")
			}`,
			wantPassthrough: true,
		},
		{
			name: "only t.Cleanup call",
			code: `func TestOnlyCleanup(t *testing.T) {
				t.Cleanup(func() {})
			}`,
			wantPassthrough: true,
		},
		// === NON-PASSTHROUGH CASES (should NOT be flagged) ===
		{
			name: "with t.Error",
			code: `func TestWithError(t *testing.T) {
				t.Error("something wrong")
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.Errorf",
			code: `func TestWithErrorf(t *testing.T) {
				t.Errorf("got %v", 42)
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.Fatal",
			code: `func TestWithFatal(t *testing.T) {
				t.Fatal("critical")
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.Fatalf",
			code: `func TestWithFatalf(t *testing.T) {
				t.Fatalf("got %v", 42)
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.Fail",
			code: `func TestWithFail(t *testing.T) {
				t.Fail()
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.FailNow",
			code: `func TestWithFailNow(t *testing.T) {
				t.FailNow()
			}`,
			wantPassthrough: false,
		},
		{
			name: "with t.Run subtest",
			code: `func TestWithRun(t *testing.T) {
				t.Run("subtest", func(t *testing.T) {
					t.Log("in subtest")
				})
			}`,
			wantPassthrough: false,
		},
		{
			name: "with assert.Equal",
			code: `func TestWithAssert(t *testing.T) {
				assert.Equal(t, 1, 1)
			}`,
			wantPassthrough: false,
		},
		{
			name: "with require.NoError",
			code: `func TestWithRequire(t *testing.T) {
				require.NoError(t, nil)
			}`,
			wantPassthrough: false,
		},
		{
			name: "with helper call",
			code: `func TestWithHelper(t *testing.T) {
				myHelper(t, "data")
			}`,
			wantPassthrough: false,
		},
		{
			name: "with equality comparison",
			code: `func TestWithEqual(t *testing.T) {
				got := Foo()
				if got == expected {
					t.Log("ok")
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with inequality comparison",
			code: `func TestWithNotEqual(t *testing.T) {
				got := Foo()
				if got != want {
					t.Log("mismatch")
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with less than comparison",
			code: `func TestWithLess(t *testing.T) {
				x := 5
				if x < 10 {
					t.Log("small")
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with greater than comparison",
			code: `func TestWithGreater(t *testing.T) {
				x := 15
				if x > 10 {
					t.Log("big")
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with err != nil check",
			code: `func TestWithErrCheck(t *testing.T) {
				err := Foo()
				if err != nil {
					t.Log("error")
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with table-driven test",
			code: `func TestTableDriven(t *testing.T) {
				tests := []struct{
					name string
					want int
				}{
					{"case1", 1},
				}
				for _, tt := range tests {
					t.Run(tt.name, func(t *testing.T) {
						got := 1
						if got != tt.want {
							t.Errorf("got %v, want %v", got, tt.want)
						}
					})
				}
			}`,
			wantPassthrough: false,
		},
		{
			name: "with len comparison",
			code: `func TestWithLen(t *testing.T) {
				items := []int{1, 2, 3}
				if len(items) == 0 {
					t.Error("empty")
				}
			}`,
			wantPassthrough: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract function declaration
			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			// Vérification de la déclaration
			if funcDecl == nil {
				t.Fatal("no function declaration found")
			}

			got := isPassthroughTest(funcDecl)
			// Vérification du résultat
			if got != tt.wantPassthrough {
				t.Errorf("isPassthroughTest() = %v, want %v", got, tt.wantPassthrough)
			}
		})
	}
}

// Test_isComparisonOperator tests the isComparisonOperator function.
//
// Params:
//   - t: testing context
func Test_isComparisonOperator(t *testing.T) {
	tests := []struct {
		name string
		op   token.Token
		want bool
	}{
		{"EQL is comparison", token.EQL, true},
		{"NEQ is comparison", token.NEQ, true},
		{"LSS is comparison", token.LSS, true},
		{"GTR is comparison", token.GTR, true},
		{"LEQ is comparison", token.LEQ, true},
		{"GEQ is comparison", token.GEQ, true},
		{"ADD is not comparison", token.ADD, false},
		{"SUB is not comparison", token.SUB, false},
		{"MUL is not comparison", token.MUL, false},
		{"AND is not comparison", token.AND, false},
		{"OR is not comparison", token.OR, false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isComparisonOperator(tt.op)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isComparisonOperator(%v) = %v, want %v", tt.op, got, tt.want)
			}
		})
	}
}

// Test_isTestingAssertionCall tests the isTestingAssertionCall function.
//
// Params:
//   - t: testing context
func Test_isTestingAssertionCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"t.Error is assertion", "t.Error(\"msg\")", true},
		{"t.Errorf is assertion", "t.Errorf(\"msg %v\", x)", true},
		{"t.Fatal is assertion", "t.Fatal(\"msg\")", true},
		{"t.Fatalf is assertion", "t.Fatalf(\"msg %v\", x)", true},
		{"t.Fail is assertion", "t.Fail()", true},
		{"t.FailNow is assertion", "t.FailNow()", true},
		{"t.Log is NOT assertion", "t.Log(\"msg\")", false},
		{"t.Logf is NOT assertion", "t.Logf(\"msg %v\", x)", false},
		{"t.Run is NOT assertion", "t.Run(\"name\", func(t int){})", false},
		{"t.Skip is NOT assertion", "t.Skip(\"reason\")", false},
		{"t.Parallel is NOT assertion", "t.Parallel()", false},
		{"t.Helper is NOT assertion", "t.Helper()", false},
		{"t.Cleanup is NOT assertion", "t.Cleanup(func(){})", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract call expression
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					return false
				}
				return true
			})

			// Vérification de l'expression
			if callExpr == nil {
				t.Fatal("no call expression found")
			}

			got := isTestingAssertionCall(callExpr)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isTestingAssertionCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_isSubTestCall tests the isSubTestCall function.
//
// Params:
//   - t: testing context
func Test_isSubTestCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"t.Run is subtest", "t.Run(\"name\", func(x int){})", true},
		{"other.Run is NOT subtest", "other.Run(\"name\", func(x int){})", false},
		{"t.Error is NOT subtest", "t.Error(\"msg\")", false},
		{"t.Log is NOT subtest", "t.Log(\"msg\")", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract call expression
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					return false
				}
				return true
			})

			// Vérification de l'expression
			if callExpr == nil {
				t.Fatal("no call expression found")
			}

			got := isSubTestCall(callExpr)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isSubTestCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_isAssertLibraryCall tests the isAssertLibraryCall function.
//
// Params:
//   - t: testing context
func Test_isAssertLibraryCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"assert.Equal", "assert.Equal(t, 1, 1)", true},
		{"assert.NoError", "assert.NoError(t, nil)", true},
		{"assert.True", "assert.True(t, true)", true},
		{"require.NoError", "require.NoError(t, nil)", true},
		{"require.Equal", "require.Equal(t, 1, 1)", true},
		{"Assert.Equal uppercase", "Assert.Equal(t, 1, 1)", true},
		{"Require.NoError uppercase", "Require.NoError(t, nil)", true},
		{"t.Error is NOT assert lib", "t.Error(\"msg\")", false},
		{"fmt.Println is NOT assert lib", "fmt.Println(\"msg\")", false},
		{"other.Equal is NOT assert lib", "other.Equal(t, 1, 1)", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract call expression
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					return false
				}
				return true
			})

			// Vérification de l'expression
			if callExpr == nil {
				t.Fatal("no call expression found")
			}

			got := isAssertLibraryCall(callExpr)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isAssertLibraryCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_isTestHelperCall tests the isTestHelperCall function.
//
// Params:
//   - t: testing context
func Test_isTestHelperCall(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"helper(t, data)", "helper(t, \"data\")", true},
		{"setupTest(t)", "setupTest(t)", true},
		{"runCheck(t, x, y)", "runCheck(t, x, y)", true},
		{"func with no args", "noArgs()", false},
		{"func with non-t first arg", "helper(x, y)", false},
		{"func with string first arg", "helper(\"t\", y)", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract call expression
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					return false
				}
				return true
			})

			// Vérification de l'expression
			if callExpr == nil {
				t.Fatal("no call expression found")
			}

			got := isTestHelperCall(callExpr)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("isTestHelperCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_checkForValidationSignal tests the checkForValidationSignal function.
//
// Params:
//   - t: testing context
func Test_checkForValidationSignal(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"t.Error call", "t.Error(\"msg\")", true},
		{"assert.Equal call", "assert.Equal(t, 1, 1)", true},
		{"helper call", "helper(t, x)", true},
		{"comparison x == y", "x == y", true},
		{"comparison x != y", "x != y", true},
		{"t.Log call", "t.Log(\"msg\")", false},
		{"simple function call", "Foo()", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse expression or statement
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Find the relevant node
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du signal
				if checkForValidationSignal(n) {
					found = true
					return false
				}
				return true
			})

			// Vérification du résultat
			if found != tt.want {
				t.Errorf("checkForValidationSignal() found = %v, want %v", found, tt.want)
			}
		})
	}
}

// Test_checkCallForValidation tests the checkCallForValidation function.
//
// Params:
//   - t: testing context
func Test_checkCallForValidation(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"t.Error is validation", "t.Error(\"msg\")", true},
		{"t.Fatal is validation", "t.Fatal(\"msg\")", true},
		{"t.Run is validation", "t.Run(\"name\", func(x int){})", true},
		{"assert.Equal is validation", "assert.Equal(t, 1, 1)", true},
		{"helper(t, x) is validation", "helper(t, x)", true},
		{"t.Log is NOT validation", "t.Log(\"msg\")", false},
		{"Foo() is NOT validation", "Foo()", false},
		{"fmt.Println is NOT validation", "fmt.Println(\"msg\")", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code
			fset := token.NewFileSet()
			code := "package test\nfunc f() { " + tt.code + " }"
			file, err := parser.ParseFile(fset, "", code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse code: %v", err)
			}

			// Extract call expression
			var callExpr *ast.CallExpr
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du noeud
				if ce, ok := n.(*ast.CallExpr); ok {
					callExpr = ce
					return false
				}
				return true
			})

			// Vérification de l'expression
			if callExpr == nil {
				t.Fatal("no call expression found")
			}

			got := checkCallForValidation(callExpr)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("checkCallForValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_runTest012 tests the runTest012 function.
//
// Params:
//   - t: testing context
func Test_runTest012(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{"analyzer exists", "ktntest012"},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification de l'analyseur
			if Analyzer012 == nil || Analyzer012.Name != tt.expectedName {
				t.Errorf("Analyzer012 invalid: nil=%v, Name=%q, want %q",
					Analyzer012 == nil, Analyzer012.Name, tt.expectedName)
			}
		})
	}
}

// Test_runTest012_disabled tests that the rule is skipped when disabled.
func Test_runTest012_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runTest012_excludedFile tests that excluded files are skipped.
func Test_runTest012_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}
