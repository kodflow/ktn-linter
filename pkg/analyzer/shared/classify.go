// Shared utilities for function classification and test name parsing.
package shared

import (
	"go/ast"
	"path/filepath"
	"strings"
	"unicode"
)

// FuncKind represents the kind of function (top-level or method).
type FuncKind int

// Visibility represents the visibility of a function.
type Visibility int

const (
	// FUNC_TOP_LEVEL is a package-level function.
	FUNC_TOP_LEVEL FuncKind = iota
	// FUNC_METHOD is a method with a receiver.
	FUNC_METHOD

	// VIS_PRIVATE is a private (unexported) function.
	VIS_PRIVATE Visibility = iota
	// VIS_PUBLIC is a public (exported) function.
	VIS_PUBLIC
)

// FuncMeta contains metadata about a function for test classification.
// It captures information needed to determine expected test names and
// visibility rules for analyzer rules like KTN-TEST-004, 009, 010.
type FuncMeta struct {
	// Name is the function name.
	Name string
	// ReceiverName is the receiver type name (empty for top-level functions).
	ReceiverName string
	// Kind indicates if it's a top-level function or method.
	Kind FuncKind
	// Visibility indicates if the function is public or private.
	Visibility Visibility
}

// IsExportedIdent checks if an identifier is exported (starts with uppercase).
//
// Params:
//   - name: identifier name to check
//
// Returns:
//   - bool: true if exported
func IsExportedIdent(name string) bool {
	// Check if name is empty
	if name == "" {
		return false
	}
	// Check first character
	return unicode.IsUpper(rune(name[0]))
}

// ClassifyFunc classifies a function declaration for test purposes.
// For methods, visibility is determined by the METHOD name only,
// not the receiver type.
//
// Params:
//   - funcDecl: function declaration to classify
//
// Returns:
//   - *FuncMeta: metadata about the function
func ClassifyFunc(funcDecl *ast.FuncDecl) *FuncMeta {
	// Initialize metadata
	meta := &FuncMeta{
		Name: funcDecl.Name.Name,
	}

	// Check if it's a method
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		// Method with receiver
		meta.Kind = FUNC_METHOD
		meta.ReceiverName = ExtractReceiverTypeName(funcDecl.Recv.List[0].Type)
		// Visibility is determined by METHOD name only
		meta.Visibility = getVisibility(meta.Name)
	} else {
		// Top-level function
		meta.Kind = FUNC_TOP_LEVEL
		// Visibility is determined by function name
		meta.Visibility = getVisibility(meta.Name)
	}
	// Return metadata
	return meta
}

// getVisibility returns the visibility based on identifier name.
//
// Params:
//   - name: identifier name
//
// Returns:
//   - Visibility: public or private
func getVisibility(name string) Visibility {
	// Check if exported
	if IsExportedIdent(name) {
		// Public identifier
		return VIS_PUBLIC
	}
	// Private identifier
	return VIS_PRIVATE
}

// ExtractReceiverTypeName extracts the type name from a receiver expression.
//
// Params:
//   - expr: receiver type expression
//
// Returns:
//   - string: type name without pointer/slice decorators
func ExtractReceiverTypeName(expr ast.Expr) string {
	// Handle different expression types using type switch
	switch t := expr.(type) {
	// Case: pointer receiver (*Type)
	case *ast.StarExpr:
		// Recurse to unwrap pointer
		return ExtractReceiverTypeName(t.X)
	// Case: simple identifier (Type)
	case *ast.Ident:
		// Base case - return name
		return t.Name
	// Case: generic type with single param (Type[T])
	case *ast.IndexExpr:
		// Recurse to get base type
		return ExtractReceiverTypeName(t.X)
	// Case: generic type with multiple params (Type[T, U])
	case *ast.IndexListExpr:
		// Recurse to get base type
		return ExtractReceiverTypeName(t.X)
	// Case: unknown expression type
	default:
		// Return empty for unhandled cases
		return ""
	}
}

// BuildSuggestedTestName builds the expected test name for a function.
// Rules:
//   - Public top-level: TestFoo
//   - Private top-level: Test_foo
//   - Method (any): TestType_Method
//
// Params:
//   - meta: function metadata
//
// Returns:
//   - string: suggested test name
func BuildSuggestedTestName(meta *FuncMeta) string {
	// Handle methods
	if meta.Kind == FUNC_METHOD {
		// Always TestType_Method for methods
		return "Test" + meta.ReceiverName + "_" + meta.Name
	}
	// Handle top-level functions
	if meta.Visibility == VIS_PRIVATE {
		// Private: Test_foo
		return "Test_" + meta.Name
	}
	// Public: TestFoo
	return "Test" + meta.Name
}

// BuildTestLookupKey builds the key used to match tests with functions.
// For methods: ReceiverName_MethodName
// For functions: FunctionName
//
// Params:
//   - meta: function metadata
//
// Returns:
//   - string: lookup key
func BuildTestLookupKey(meta *FuncMeta) string {
	// Handle methods
	if meta.Kind == FUNC_METHOD {
		// Method key format
		return meta.ReceiverName + "_" + meta.Name
	}
	// Top-level function
	return meta.Name
}

// ParseTestName parses a test function name to extract its target.
// Patterns:
//   - TestFoo -> targets Foo (public function)
//   - Test_foo -> targets foo (private function)
//   - TestType_Method -> targets Type.Method
//   - Testtype_Method -> targets type.Method (method on private type)
//
// Params:
//   - testName: name of the test function
//
// Returns:
//   - TestTarget: parsed target info
//   - bool: true if parsing succeeded
func ParseTestName(testName string) (TestTarget, bool) {
	// Must start with Test
	body, hasPrefix := strings.CutPrefix(testName, "Test")
	// No Test prefix
	if !hasPrefix || body == "" {
		return TestTarget{}, false
	}

	// Check for private function pattern (Test_foo)
	if privateName, isPrivate := strings.CutPrefix(body, "_"); isPrivate {
		// Parse private function test
		return parsePrivateTestName(privateName)
	}

	// Check for method pattern (TestType_Method)
	if typeName, methodName, hasMethod := strings.Cut(body, "_"); hasMethod && methodName != "" {
		// Valid method pattern
		return TestTarget{
			FuncName:     methodName,
			ReceiverName: typeName,
			IsPrivate:    !IsExportedIdent(methodName),
			IsMethod:     true,
		}, true
	}

	// Public function TestFoo
	return TestTarget{
		FuncName:  body,
		IsPrivate: false,
		IsMethod:  false,
	}, true
}

// parsePrivateTestName parses the body of a private test name (after Test_).
//
// Params:
//   - privateName: test name body after "Test_"
//
// Returns:
//   - TestTarget: parsed target info
//   - bool: true if parsing succeeded
func parsePrivateTestName(privateName string) (TestTarget, bool) {
	// Empty after underscore
	if privateName == "" {
		return TestTarget{}, false
	}
	// Check if it's a method (Type_method pattern after underscore)
	if typeName, methodName, hasMethod := strings.Cut(privateName, "_"); hasMethod && methodName != "" {
		// Test_Type_method pattern
		return TestTarget{
			FuncName:     methodName,
			ReceiverName: typeName,
			IsPrivate:    !IsExportedIdent(methodName),
			IsMethod:     true,
		}, true
	}
	// Simple private function Test_foo
	return TestTarget{
		FuncName:  privateName,
		IsPrivate: true,
		IsMethod:  false,
	}, true
}

// BuildTestTargetKey builds a lookup key from a test target.
//
// Params:
//   - target: test target
//
// Returns:
//   - string: lookup key
func BuildTestTargetKey(target TestTarget) string {
	// Handle methods
	if target.IsMethod && target.ReceiverName != "" {
		// Method key format
		return target.ReceiverName + "_" + target.FuncName
	}
	// Top-level function
	return target.FuncName
}

// IsMockFile checks if a file is a mock file based on its name.
//
// Params:
//   - filename: path to the file
//
// Returns:
//   - bool: true if it's a mock file
func IsMockFile(filename string) bool {
	// Get base name
	base := filepath.Base(filename)
	baseLower := strings.ToLower(base)
	// Check common mock patterns
	return strings.Contains(baseLower, "_mock") ||
		strings.Contains(baseLower, "mock_") ||
		strings.HasSuffix(baseLower, "mock.go") ||
		strings.HasPrefix(baseLower, "mock")
}

// IsMockName checks if a function/type name indicates a mock.
//
// Params:
//   - name: function or type name
//
// Returns:
//   - bool: true if it's a mock name
func IsMockName(name string) bool {
	// Lowercase for comparison
	lower := strings.ToLower(name)
	// Check common mock patterns
	return strings.HasPrefix(lower, "mock") ||
		strings.HasPrefix(lower, "fake") ||
		strings.HasPrefix(lower, "stub") ||
		strings.HasSuffix(lower, "mock") ||
		strings.HasSuffix(lower, "fake") ||
		strings.HasSuffix(lower, "stub") ||
		strings.Contains(lower, "mock") ||
		strings.Contains(lower, "fake")
}

// IsExemptTestFile checks if a test file should be exempt from certain rules.
//
// Params:
//   - filename: path to the file
//
// Returns:
//   - bool: true if exempt
func IsExemptTestFile(filename string) bool {
	// Get base name
	base := filepath.Base(filename)
	baseLower := strings.ToLower(base)
	// Exempt patterns
	exemptPatterns := []string{
		"helper_test.go",
		"helpers_test.go",
		"integration_test.go",
		"suite_test.go",
		"main_test.go",
		"setup_test.go",
		"fixtures_test.go",
		"testutil",
		"testhelper",
	}
	// Check patterns
	for _, pattern := range exemptPatterns {
		// Check if contains pattern
		if strings.Contains(baseLower, pattern) {
			return true
		}
	}
	// Not exempt
	return false
}

// IsExemptTestName checks if a test function name should be exempt.
//
// Params:
//   - testName: name of the test function
//
// Returns:
//   - bool: true if exempt
func IsExemptTestName(testName string) bool {
	// TestMain is always exempt
	if testName == "TestMain" {
		return true
	}
	// Check mock-related names
	if IsMockName(testName) {
		return true
	}
	// Check exact exempt test names (case-insensitive)
	// TestSetup, Test_setup, TestTeardown, Test_teardown, TestHelper, etc.
	lower := strings.ToLower(testName)
	// Match word boundaries - not substring to avoid false positives
	exemptSuffixes := []string{
		"setup",
		"teardown",
		"init",
		"helper",
		"helpers",
		"util",
		"utils",
		"fixture",
		"fixtures",
	}
	// Extract the body after "Test" or "Test_"
	body := strings.TrimPrefix(lower, "test")
	body = strings.TrimPrefix(body, "_")
	// Check if body starts with an exempt suffix followed by end or underscore
	for _, suffix := range exemptSuffixes {
		// Check if body equals suffix or starts with suffix_
		if body == suffix || strings.HasPrefix(body, suffix+"_") {
			return true
		}
	}
	// Not exempt
	return false
}
