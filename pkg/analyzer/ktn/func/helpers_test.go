package ktnfunc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Tests for countPureCodeLines
func TestCountPureCodeLines(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "simple function",
			code: `package test
func Simple() {
	x := 1
	y := 2
	z := x + y
	_ = z
}`,
			expected: 4,
		},
		{
			name: "function with comments",
			code: `package test
func WithComments() {
	// This is a comment
	x := 1
	// Another comment
	y := 2
	z := x + y
	_ = z
}`,
			expected: 4,
		},
		{
			name: "function with block comments",
			code: `package test
func WithBlockComments() {
	/* Block comment
	   spanning multiple lines
	*/
	x := 1
	y := 2
	z := x + y
	_ = z
}`,
			expected: 4,
		},
		{
			name: "empty function",
			code: `package test
func Empty() {
}`,
			expected: 0,
		},
		{
			name: "nil body",
			code: "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code == "" {
				result := countPureCodeLines(nil, nil)
				if result != tt.expected {
					t.Errorf("countPureCodeLines(nil) = %d, want %d", result, tt.expected)
				}
				return
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Body == nil {
				t.Fatal("No function found in code")
			}

			pass := &analysis.Pass{
				Fset: fset,
				ReadFile: func(filename string) ([]byte, error) {
					return []byte(tt.code), nil
				},
			}

			result := countPureCodeLines(pass, funcDecl.Body)
			if result != tt.expected {
				t.Errorf("countPureCodeLines() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// Test for countPureCodeLines with ReadFile error
func TestCountPureCodeLinesReadFileError(t *testing.T) {
	fset := token.NewFileSet()
	code := `package test
func Test() {
	x := 1
}`
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	// Pass with ReadFile returning error
	pass := &analysis.Pass{
		Fset: fset,
		ReadFile: func(filename string) ([]byte, error) {
			return nil, &testError{}
		},
	}

	result := countPureCodeLines(pass, funcDecl.Body)
	if result != 0 {
		t.Errorf("countPureCodeLines with ReadFile error = %d, want 0", result)
	}
}

// Test for countPureCodeLines with nil ReadFile
func TestCountPureCodeLinesNilReadFile(t *testing.T) {
	fset := token.NewFileSet()
	code := `package test
func Test() {
	x := 1
}`
	file, _ := parser.ParseFile(fset, "test.go", code, 0)

	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	// Pass with nil ReadFile
	pass := &analysis.Pass{
		Fset:     fset,
		ReadFile: nil,
	}

	result := countPureCodeLines(pass, funcDecl.Body)
	if result != 0 {
		t.Errorf("countPureCodeLines with nil ReadFile = %d, want 0", result)
	}
}

// Test for countPureCodeLines with invalid line index
func TestCountPureCodeLinesInvalidIndex(t *testing.T) {
	fset := token.NewFileSet()
	// Create a realistic source code
	sourceCode := `package test

func LongFunction() {
	x := 1
	y := 2
	z := 3
	a := 4
	b := 5
	c := 6
}
`
	// Add a file and parse it
	file := fset.AddFile("test.go", fset.Base(), len(sourceCode))
	file.SetLinesForContent([]byte(sourceCode))

	// Parse to get a real function
	astFile, _ := parser.ParseFile(fset, "test.go", sourceCode, 0)
	var funcDecl *ast.FuncDecl
	ast.Inspect(astFile, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	pass := &analysis.Pass{
		Fset: fset,
		ReadFile: func(filename string) ([]byte, error) {
			// Return FEWER lines than the source expects
			// This simulates a file that was modified/truncated
			return []byte("package test\nfunc LongFunction() {\n"), nil
		},
	}

	// This should handle the invalid index gracefully (lines beyond file length)
	result := countPureCodeLines(pass, funcDecl.Body)
	// The function should skip invalid indices and return what it can count
	if result < 0 {
		t.Errorf("countPureCodeLines with truncated file = %d, should be >= 0", result)
	}
}

type testError struct{}

func (e *testError) Error() string { return "test error" }

// Test for isLineToSkip
func TestIsLineToSkip(t *testing.T) {
	tests := []struct {
		name           string
		trimmed        string
		inBlockComment bool
		want           bool
		wantBlockState bool
	}{
		{
			name:           "empty line",
			trimmed:        "",
			inBlockComment: false,
			want:           true,
			wantBlockState: false,
		},
		{
			name:           "line comment",
			trimmed:        "// comment",
			inBlockComment: false,
			want:           true,
			wantBlockState: false,
		},
		{
			name:           "opening brace",
			trimmed:        "{",
			inBlockComment: false,
			want:           true,
			wantBlockState: false,
		},
		{
			name:           "closing brace",
			trimmed:        "}",
			inBlockComment: false,
			want:           true,
			wantBlockState: false,
		},
		{
			name:           "start block comment",
			trimmed:        "/* comment",
			inBlockComment: false,
			want:           true,
			wantBlockState: true,
		},
		{
			name:           "end block comment",
			trimmed:        "comment */",
			inBlockComment: true,
			want:           true,
			wantBlockState: false,
		},
		{
			name:           "inside block comment",
			trimmed:        "comment line",
			inBlockComment: true,
			want:           true,
			wantBlockState: true,
		},
		{
			name:           "code line",
			trimmed:        "x := 1",
			inBlockComment: false,
			want:           false,
			wantBlockState: false,
		},
		{
			name:           "single line block comment",
			trimmed:        "/* comment */ x := 1",
			inBlockComment: false,
			want:           true,
			wantBlockState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inBlock := tt.inBlockComment
			got := isLineToSkip(tt.trimmed, &inBlock)
			if got != tt.want {
				t.Errorf("isLineToSkip(%q, %v) = %v, want %v", tt.trimmed, tt.inBlockComment, got, tt.want)
			}
			if inBlock != tt.wantBlockState {
				t.Errorf("isLineToSkip(%q, %v) block state = %v, want %v", tt.trimmed, tt.inBlockComment, inBlock, tt.wantBlockState)
			}
		})
	}
}

// Tests for isTestFunction
func TestIsTestFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     bool
	}{
		{"Test function", "TestSomething", true},
		{"Benchmark function", "BenchmarkSomething", true},
		{"Example function", "ExampleSomething", true},
		{"Fuzz function", "FuzzSomething", true},
		{"Regular function", "RegularFunction", false},
		{"test lowercase", "testSomething", false},
		{"Empty name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTestFunction(tt.funcName)
			if result != tt.want {
				t.Errorf("isTestFunction(%q) = %v, want %v", tt.funcName, result, tt.want)
			}
		})
	}
}

// Tests for hasNamedReturns
func TestHasNamedReturns(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "no returns",
			code: "func Test() {}",
			want: false,
		},
		{
			name: "unnamed returns",
			code: "func Test() int { return 0 }",
			want: false,
		},
		{
			name: "named returns",
			code: "func Test() (result int) { return }",
			want: true,
		},
		{
			name: "multiple named returns",
			code: "func Test() (a int, err error) { return }",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			code := "package test\n" + tt.code
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			var results *ast.FieldList
			if funcDecl != nil && funcDecl.Type != nil {
				results = funcDecl.Type.Results
			}

			got := hasNamedReturns(results)
			if got != tt.want {
				t.Errorf("hasNamedReturns() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test hasNamedReturns with nil
func TestHasNamedReturnsNil(t *testing.T) {
	if hasNamedReturns(nil) {
		t.Error("hasNamedReturns(nil) should return false")
	}
}

// Tests for extractFirstWord
func TestExtractFirstWord(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"GetUser", "Get"},
		{"SetValue", "Set"},
		{"IsValid", "Is"},
		{"HTTPServer", "H"}, // First uppercase letter before next uppercase
		{"A", "A"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFirstWord(tt.name)
			if got != tt.want {
				t.Errorf("extractFirstWord(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

// Tests for isGetter
func TestIsGetter(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"GetValue", true},
		{"IsValid", true},
		{"HasData", true},
		{"SetValue", false},
		{"UpdateData", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isGetter(tt.name)
			if got != tt.want {
				t.Errorf("isGetter(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// Tests for hasSideEffect
func TestHasSideEffect(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "selector expression",
			code: "m.field",
			want: true,
		},
		{
			name: "index with selector",
			code: "m.array[0]",
			want: true,
		},
		{
			name: "simple identifier",
			code: "x",
			want: false,
		},
		{
			name: "index with ident",
			code: "arr[0]",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			got := hasSideEffect(expr)
			if got != tt.want {
				t.Errorf("hasSideEffect(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// Tests for isContextType
func TestIsContextType(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "context.Context",
			code: "context.Context",
			want: true,
		},
		{
			name: "other.Context",
			code: "other.Context",
			want: false,
		},
		{
			name: "context.Other",
			code: "context.Other",
			want: false,
		},
		{
			name: "simple ident",
			code: "int",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			got := isContextType(expr)
			if got != tt.want {
				t.Errorf("isContextType(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// Tests for startsWithVerb
func TestStartsWithVerb(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"GetUser", true},
		{"SetValue", true},
		{"CreateAccount", true},
		{"UserData", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startsWithVerb(tt.name)
			if got != tt.want {
				t.Errorf("startsWithVerb(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// Tests for isErrorType
func TestIsErrorTypeWithNilTypesInfo(t *testing.T) {
	ident := &ast.Ident{Name: "error"}
	pass := &analysis.Pass{
		TypesInfo: nil,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("isErrorType with nil TypesInfo should panic")
		}
	}()

	isErrorType(pass, ident)
}

func TestIsErrorTypeWithEmptyTypesInfo(t *testing.T) {
	ident := &ast.Ident{Name: "error"}
	pass := &analysis.Pass{
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	result := isErrorType(pass, ident)
	if result {
		t.Error("isErrorType with empty TypesInfo should return false")
	}
}

func TestIsErrorTypeWithNonIdent(t *testing.T) {
	fset := token.NewFileSet()
	code := "package test\nfunc Test() *int { return nil }"
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	if funcDecl == nil {
		t.Fatal("No function found")
	}

	pass := &analysis.Pass{
		Fset: fset,
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
	}

	returnType := funcDecl.Type.Results.List[0].Type
	result := isErrorType(pass, returnType)
	if result {
		t.Error("isErrorType(*int) should return false")
	}
}

// Tests for calculateComplexity
func TestCalculateComplexity(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "simple function",
			code: `package test
func Simple() {
	x := 1
}`,
			expected: 1,
		},
		{
			name: "one if",
			code: `package test
func OneIf() {
	if true {
		x := 1
		_ = x
	}
}`,
			expected: 2,
		},
		{
			name: "one for loop",
			code: `package test
func OneFor() {
	for i := 0; i < 10; i++ {
		_ = i
	}
}`,
			expected: 2,
		},
		{
			name: "range loop",
			code: `package test
func RangeLoop() {
	arr := []int{1, 2, 3}
	for _, v := range arr {
		_ = v
	}
}`,
			expected: 2,
		},
		{
			name: "switch with cases",
			code: `package test
func SwitchCase() {
	switch 1 {
	case 1:
		x := 1
		_ = x
	case 2:
		y := 2
		_ = y
	default:
		z := 3
		_ = z
	}
}`,
			expected: 3,
		},
		{
			name: "logical operators",
			code: `package test
func LogicalOps() {
	if true && false {
		x := 1
		_ = x
	}
	if true || false {
		y := 2
		_ = y
	}
}`,
			expected: 5,
		},
		{
			name: "select with cases",
			code: `package test
func SelectCase() {
	ch := make(chan int)
	select {
	case x := <-ch:
		_ = x
	case ch <- 1:
		y := 2
		_ = y
	default:
		z := 3
		_ = z
	}
}`,
			expected: 3,
		},
		{
			name: "nested complexity",
			code: `package test
func NestedComplexity() {
	for i := 0; i < 10; i++ {
		if i > 5 && i < 8 {
			switch i {
			case 6:
				x := 1
				_ = x
			case 7:
				y := 2
				_ = y
			}
		}
	}
}`,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			var funcDecl *ast.FuncDecl
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fd
					return false
				}
				return true
			})

			if funcDecl == nil || funcDecl.Body == nil {
				t.Fatal("No function found")
			}

			got := calculateComplexity(funcDecl.Body)
			if got != tt.expected {
				t.Errorf("calculateComplexity() = %d, want %d", got, tt.expected)
			}
		})
	}
}

// Test for validateDocFormat edge cases
func TestValidateDocFormatEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		comments   []string
		funcName   string
		hasParams  bool
		hasReturns bool
		wantError  bool
	}{
		{
			name:       "blank line between sections",
			comments:   []string{"// Test func", "//", "// Params:", "// - x: value"},
			funcName:   "Test",
			hasParams:  true,
			hasReturns: false,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateDocFormat(tt.comments, tt.funcName, tt.hasParams, tt.hasReturns)
			hasError := result != ""
			if hasError != tt.wantError {
				t.Errorf("validateDocFormat() error = %v, wantError %v (message: %s)", hasError, tt.wantError, result)
			}
		})
	}
}

