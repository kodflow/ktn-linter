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

const (
	// FuncTopLevel is a package-level function.
	FuncTopLevel FuncKind = iota
	// FuncMethod is a method with a receiver.
	FuncMethod
)

// Visibility represents the visibility of a function.
type Visibility int

const (
	// VisPrivate is a private (unexported) function.
	VisPrivate Visibility = iota
	// VisPublic is a public (exported) function.
	VisPublic
)

// FuncMeta contains metadata about a function for test classification.
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

// TestTarget represents the target of a test function.
type TestTarget struct {
	// FuncName is the name of the function being tested.
	FuncName string
	// ReceiverName is the receiver type name (for method tests).
	ReceiverName string
	// IsPrivate indicates if this targets a private function.
	IsPrivate bool
	// IsMethod indicates if this targets a method.
	IsMethod bool
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
//   - FuncMeta: metadata about the function
func ClassifyFunc(funcDecl *ast.FuncDecl) FuncMeta {
	// Initialize metadata
	meta := FuncMeta{
		Name: funcDecl.Name.Name,
	}

	// Check if it's a method
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		meta.Kind = FuncMethod
		meta.ReceiverName = ExtractReceiverTypeName(funcDecl.Recv.List[0].Type)
		// Visibility is determined by METHOD name only
		if IsExportedIdent(meta.Name) {
			meta.Visibility = VisPublic
		} else {
			meta.Visibility = VisPrivate
		}
	} else {
		// Top-level function
		meta.Kind = FuncTopLevel
		// Visibility is determined by function name
		if IsExportedIdent(meta.Name) {
			meta.Visibility = VisPublic
		} else {
			meta.Visibility = VisPrivate
		}
	}

	return meta
}

// ExtractReceiverTypeName extracts the type name from a receiver expression.
//
// Params:
//   - expr: receiver type expression
//
// Returns:
//   - string: type name without pointer/slice decorators
func ExtractReceiverTypeName(expr ast.Expr) string {
	// Handle pointer types
	switch t := expr.(type) {
	case *ast.StarExpr:
		// Pointer receiver (*Type)
		return ExtractReceiverTypeName(t.X)
	case *ast.Ident:
		// Simple identifier
		return t.Name
	case *ast.IndexExpr:
		// Generic type (Type[T])
		return ExtractReceiverTypeName(t.X)
	case *ast.IndexListExpr:
		// Generic type with multiple params (Type[T, U])
		return ExtractReceiverTypeName(t.X)
	default:
		// Unknown type
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
func BuildSuggestedTestName(meta FuncMeta) string {
	// Handle methods
	if meta.Kind == FuncMethod {
		// Always TestType_Method for methods
		return "Test" + meta.ReceiverName + "_" + meta.Name
	}
	// Handle top-level functions
	if meta.Visibility == VisPrivate {
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
func BuildTestLookupKey(meta FuncMeta) string {
	// Handle methods
	if meta.Kind == FuncMethod {
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
	if !strings.HasPrefix(testName, "Test") {
		return TestTarget{}, false
	}

	// Remove Test prefix
	body := strings.TrimPrefix(testName, "Test")
	// Empty body
	if body == "" {
		return TestTarget{}, false
	}

	// Check for private function pattern (Test_foo)
	if strings.HasPrefix(body, "_") {
		privateName := strings.TrimPrefix(body, "_")
		// Empty after underscore
		if privateName == "" {
			return TestTarget{}, false
		}
		// Check if it's a method (Type_method pattern after underscore)
		if idx := strings.Index(privateName, "_"); idx > 0 {
			// Could be Test_Type_method
			typeName := privateName[:idx]
			methodName := privateName[idx+1:]
			// Check if method
			if methodName != "" {
				return TestTarget{
					FuncName:     methodName,
					ReceiverName: typeName,
					IsPrivate:    !IsExportedIdent(methodName),
					IsMethod:     true,
				}, true
			}
		}
		// Simple private function Test_foo
		return TestTarget{
			FuncName:  privateName,
			IsPrivate: true,
			IsMethod:  false,
		}, true
	}

	// Check for method pattern (TestType_Method)
	if idx := strings.Index(body, "_"); idx > 0 {
		typeName := body[:idx]
		methodName := body[idx+1:]
		// Valid method pattern
		if methodName != "" {
			return TestTarget{
				FuncName:     methodName,
				ReceiverName: typeName,
				IsPrivate:    !IsExportedIdent(methodName),
				IsMethod:     true,
			}, true
		}
	}

	// Public function TestFoo
	return TestTarget{
		FuncName:  body,
		IsPrivate: false,
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
	// Check setup/teardown patterns
	lower := strings.ToLower(testName)
	exemptPatterns := []string{
		"setup",
		"teardown",
		"init",
		"helper",
		"util",
		"fixture",
	}
	// Check patterns
	for _, pattern := range exemptPatterns {
		// Check if contains pattern
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	// Not exempt
	return false
}
