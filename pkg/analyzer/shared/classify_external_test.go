// External tests for classify.go.
package shared_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
)

// TestIsExportedIdent tests the IsExportedIdent function.
func TestIsExportedIdent(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Exported identifiers
		{"uppercase", "Foo", true},
		{"uppercase_underscore", "Foo_Bar", true},
		{"single_upper", "F", true},
		// Unexported identifiers
		{"lowercase", "foo", false},
		{"lowercase_underscore", "foo_bar", false},
		{"single_lower", "f", false},
		{"underscore_prefix", "_foo", false},
		// Edge cases
		{"empty", "", false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := shared.IsExportedIdent(tt.input)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("IsExportedIdent(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestClassifyFunc tests the ClassifyFunc function.
func TestClassifyFunc(t *testing.T) {
	// Define expected results
	tests := []struct {
		name         string
		code         string
		wantKind     shared.FuncKind
		wantVis      shared.Visibility
		wantReceiver string
	}{
		{"PublicFunc", "func PublicFunc() {}", shared.FuncTopLevel, shared.VisPublic, ""},
		{"privateFunc", "func privateFunc() {}", shared.FuncTopLevel, shared.VisPrivate, ""},
		{"PublicMethod", "func (s *Service) PublicMethod() {}", shared.FuncMethod, shared.VisPublic, "Service"},
		{"privateMethod", "func (s *Service) privateMethod() {}", shared.FuncMethod, shared.VisPrivate, "Service"},
		{"PublicOnPrivate", "func (s *service) PublicOnPrivate() {}", shared.FuncMethod, shared.VisPublic, "service"},
		{"privateOnPrivate", "func (s *service) privateOnPrivate() {}", shared.FuncMethod, shared.VisPrivate, "service"},
	}
	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create file set
			fset := token.NewFileSet()
			// Parse source
			src := "package test\n" + tt.code
			file, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				// Fail test
				t.Fatalf("Failed to parse: %v", err)
			}
			// Get function decl
			fn := file.Decls[0].(*ast.FuncDecl)
			// Classify function
			meta := shared.ClassifyFunc(fn)
			// Check kind
			if meta.Kind != tt.wantKind {
				// Report error
				t.Errorf("Kind = %v, want %v", meta.Kind, tt.wantKind)
			}
			// Check visibility
			if meta.Visibility != tt.wantVis {
				// Report error
				t.Errorf("Visibility = %v, want %v", meta.Visibility, tt.wantVis)
			}
			// Check receiver
			if meta.ReceiverName != tt.wantReceiver {
				// Report error
				t.Errorf("ReceiverName = %q, want %q", meta.ReceiverName, tt.wantReceiver)
			}
		})
	}
}

// TestBuildSuggestedTestName tests the BuildSuggestedTestName function.
func TestBuildSuggestedTestName(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		meta     *shared.FuncMeta
		expected string
	}{
		// Public top-level
		{"public_func", &shared.FuncMeta{Name: "Foo", Kind: shared.FuncTopLevel, Visibility: shared.VisPublic}, "TestFoo"},
		// Private top-level
		{"private_func", &shared.FuncMeta{Name: "foo", Kind: shared.FuncTopLevel, Visibility: shared.VisPrivate}, "Test_foo"},
		// Public method
		{"public_method", &shared.FuncMeta{Name: "Bar", ReceiverName: "Service", Kind: shared.FuncMethod, Visibility: shared.VisPublic}, "TestService_Bar"},
		// Private method
		{"private_method", &shared.FuncMeta{Name: "bar", ReceiverName: "Service", Kind: shared.FuncMethod, Visibility: shared.VisPrivate}, "TestService_bar"},
		// Public method on private type
		{"public_on_private", &shared.FuncMeta{Name: "Baz", ReceiverName: "service", Kind: shared.FuncMethod, Visibility: shared.VisPublic}, "Testservice_Baz"},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Build name
			got := shared.BuildSuggestedTestName(tt.meta)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("BuildSuggestedTestName(%+v) = %q, want %q", tt.meta, got, tt.expected)
			}
		})
	}
}

// TestParseTestName tests the ParseTestName function.
func TestParseTestName(t *testing.T) {
	// Define test cases
	tests := []struct {
		name         string
		testName     string
		wantOK       bool
		wantFunc     string
		wantReceiver string
		wantPrivate  bool
		wantMethod   bool
	}{
		// Valid patterns
		{"public_func", "TestFoo", true, "Foo", "", false, false},
		{"private_func", "Test_foo", true, "foo", "", true, false},
		{"method", "TestService_Bar", true, "Bar", "Service", false, true},
		{"private_method", "TestService_bar", true, "bar", "Service", true, true},
		{"method_on_private", "Testservice_Bar", true, "Bar", "service", false, true},
		// Private type with private method
		{"priv_type_method", "Test_Type_method", true, "method", "Type", true, true},
		// Invalid patterns
		{"no_test_prefix", "Foo", false, "", "", false, false},
		{"empty_body", "Test", false, "", "", false, false},
		{"empty_private", "Test_", false, "", "", false, false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Parse name
			target, ok := shared.ParseTestName(tt.testName)
			// Check parse success
			if ok != tt.wantOK {
				// Report error
				t.Errorf("ParseTestName(%q) ok = %v, want %v", tt.testName, ok, tt.wantOK)
				// Skip rest
				return
			}
			// Skip if not expected to parse
			if !tt.wantOK {
				// Done
				return
			}
			// Check func name
			if target.FuncName != tt.wantFunc {
				// Report error
				t.Errorf("ParseTestName(%q).FuncName = %q, want %q", tt.testName, target.FuncName, tt.wantFunc)
			}
			// Check receiver
			if target.ReceiverName != tt.wantReceiver {
				// Report error
				t.Errorf("ParseTestName(%q).ReceiverName = %q, want %q", tt.testName, target.ReceiverName, tt.wantReceiver)
			}
			// Check private
			if target.IsPrivate != tt.wantPrivate {
				// Report error
				t.Errorf("ParseTestName(%q).IsPrivate = %v, want %v", tt.testName, target.IsPrivate, tt.wantPrivate)
			}
			// Check method
			if target.IsMethod != tt.wantMethod {
				// Report error
				t.Errorf("ParseTestName(%q).IsMethod = %v, want %v", tt.testName, target.IsMethod, tt.wantMethod)
			}
		})
	}
}

// TestBuildTestLookupKey tests the BuildTestLookupKey function.
func TestBuildTestLookupKey(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		meta     *shared.FuncMeta
		expected string
	}{
		// Top-level function
		{"top_level", &shared.FuncMeta{Name: "Foo", Kind: shared.FuncTopLevel}, "Foo"},
		// Method
		{"method", &shared.FuncMeta{Name: "Bar", ReceiverName: "Service", Kind: shared.FuncMethod}, "Service_Bar"},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Build key
			got := shared.BuildTestLookupKey(tt.meta)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("BuildTestLookupKey(%+v) = %q, want %q", tt.meta, got, tt.expected)
			}
		})
	}
}

// TestBuildTestTargetKey tests the BuildTestTargetKey function.
func TestBuildTestTargetKey(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		target   shared.TestTarget
		expected string
	}{
		// Function
		{"function", shared.TestTarget{FuncName: "Foo", IsMethod: false}, "Foo"},
		// Method
		{"method", shared.TestTarget{FuncName: "Bar", ReceiverName: "Service", IsMethod: true}, "Service_Bar"},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Build key
			got := shared.BuildTestTargetKey(tt.target)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("BuildTestTargetKey(%+v) = %q, want %q", tt.target, got, tt.expected)
			}
		})
	}
}

// TestIsMockFile tests the IsMockFile function.
func TestIsMockFile(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		// Mock patterns
		{"suffix_mock", "user_mock.go", true},
		{"prefix_mock", "mock_user.go", true},
		{"ends_mock", "usermock.go", true},
		{"starts_mock", "mockuser.go", true},
		{"path_mock", "/path/to/mock_service.go", true},
		// Non-mock
		{"regular", "user.go", false},
		{"test", "user_test.go", false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := shared.IsMockFile(tt.filename)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("IsMockFile(%q) = %v, want %v", tt.filename, got, tt.expected)
			}
		})
	}
}

// TestIsMockName tests the IsMockName function.
func TestIsMockName(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Mock patterns
		{"prefix_mock", "MockService", true},
		{"prefix_fake", "FakeRepository", true},
		{"prefix_stub", "StubHandler", true},
		{"suffix_mock", "ServiceMock", true},
		{"suffix_fake", "RepositoryFake", true},
		{"contains_mock", "UserMockHelper", true},
		// Non-mock
		{"regular", "UserService", false},
		{"test_func", "TestService", false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := shared.IsMockName(tt.input)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("IsMockName(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsExemptTestFile tests the IsExemptTestFile function.
func TestIsExemptTestFile(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		// Exempt patterns
		{"helper", "helper_test.go", true},
		{"helpers", "helpers_test.go", true},
		{"integration", "integration_test.go", true},
		{"suite", "suite_test.go", true},
		{"main", "main_test.go", true},
		{"setup", "setup_test.go", true},
		{"fixtures", "fixtures_test.go", true},
		{"testutil", "testutil.go", true},
		{"testhelper", "testhelper.go", true},
		// Non-exempt
		{"regular", "user_test.go", false},
		{"external", "user_external_test.go", false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := shared.IsExemptTestFile(tt.filename)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("IsExemptTestFile(%q) = %v, want %v", tt.filename, got, tt.expected)
			}
		})
	}
}

// TestIsExemptTestName tests the IsExemptTestName function.
func TestIsExemptTestName(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		testName string
		expected bool
	}{
		// Exempt patterns
		{"testmain", "TestMain", true},
		{"mock_prefix", "MockService", true},
		{"setup", "TestSetup", true},
		{"helper", "TestHelper", true},
		// Non-exempt
		{"regular", "TestFoo", false},
		{"method", "TestService_Bar", false},
	}
	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			// Check result
			got := shared.IsExemptTestName(tt.testName)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("IsExemptTestName(%q) = %v, want %v", tt.testName, got, tt.expected)
			}
		})
	}
}

// TestExtractReceiverTypeName tests the ExtractReceiverTypeName function.
func TestExtractReceiverTypeName(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"value receiver", "func (s Service) Method() {}", "Service"},
		{"pointer receiver", "func (s *Service) Method() {}", "Service"},
		{"generic receiver", "func (s Service[T]) Method() {}", "Service"},
		{"generic pointer receiver", "func (s *Service[T, U]) Method() {}", "Service"},
	}
	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create file set
			fset := token.NewFileSet()
			// Parse source
			src := "package test\n" + tt.code
			file, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				// Fail test
				t.Fatalf("Failed to parse: %v", err)
			}
			// Get function decl
			fn := file.Decls[0].(*ast.FuncDecl)
			// Extract receiver name
			got := shared.ExtractReceiverTypeName(fn.Recv.List[0].Type)
			// Validate
			if got != tt.expected {
				// Report error
				t.Errorf("ReceiverTypeName = %q, want %q", got, tt.expected)
			}
		})
	}
}
