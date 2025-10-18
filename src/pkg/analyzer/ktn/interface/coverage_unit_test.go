package ktn_interface_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
)

// TestNeedsInterfacesFile teste les différents cas de needsInterfacesFile.
func TestNeedsInterfacesFile(t *testing.T) {
	tests := []struct {
		name     string
		info     *ktn_interface.PackageInfo
		expected bool
	}{
		{
			name: "package avec struct publique non autorisée",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "MyService", FileName: "service.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: true,
		},
		{
			name: "package avec struct publique autorisée (Config)",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "MyConfig", FileName: "config.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec struct publique autorisée (ID)",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "UserID", FileName: "types.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec struct publique autorisée (Type)",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "ErrorType", FileName: "types.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec interface publique",
			info: &ktn_interface.PackageInfo{
				PublicStructs:    []ktn_interface.PublicStruct{},
				PublicInterfaces: []ktn_interface.PublicInterface{{Name: "MyInterface"}},
			},
			expected: true,
		},
		{
			name: "package sans types publics",
			info: &ktn_interface.PackageInfo{
				PublicStructs:    []ktn_interface.PublicStruct{},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec multiple structs dont une non autorisée",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "MyConfig", FileName: "config.go"},
					{Name: "MyService", FileName: "service.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: true,
		},
		{
			name: "package avec structs toutes autorisées",
			info: &ktn_interface.PackageInfo{
				PublicStructs: []ktn_interface.PublicStruct{
					{Name: "MyConfig", FileName: "config.go"},
					{Name: "UserID", FileName: "types.go"},
					{Name: "ErrorType", FileName: "types.go"},
				},
				PublicInterfaces: []ktn_interface.PublicInterface{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.NeedsInterfacesFile(tt.info)
			if got != tt.expected {
				t.Errorf("ktn_interface.NeedsInterfacesFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestIsAllowedPublicType teste tous les suffixes autorisés.
func TestIsAllowedPublicType(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		// Autorisés
		{"ID suffix", "UserID", true},
		{"Id suffix", "UserId", true},
		{"Type suffix", "ErrorType", true},
		{"Kind suffix", "EventKind", true},
		{"Status suffix", "OrderStatus", true},
		{"State suffix", "GameState", true},
		{"Count suffix", "ItemCount", true},
		{"Size suffix", "BufferSize", true},
		{"Index suffix", "StartIndex", true},
		{"Name suffix", "UserName", true},
		{"Title suffix", "BookTitle", true},
		{"Config suffix", "AppConfig", true},
		{"Options suffix", "BuildOptions", true},
		{"Settings suffix", "UserSettings", true},
		{"Data suffix", "MetaData", true},

		// Non autorisés
		{"Service", "UserService", false},
		{"Handler", "RequestHandler", false},
		{"Manager", "SessionManager", false},
		{"Random", "RandomStruct", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsAllowedPublicType(tt.typeName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsAllowedPublicType(%q) = %v, want %v", tt.typeName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage teste tous les packages exemptés.
func TestIsExemptedPackage(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		expected bool
	}{
		{"main package", "main", true},
		{"main_test package", "main_test", true},
		{"test suffix", "mypackage_test", true},
		{"normal package", "mypackage", false},
		{"internal package", "internal", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsExemptedPackage(tt.pkgName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsExemptedPackage(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage002 teste tous les packages exemptés pour Rule002.
func TestIsExemptedPackage002(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		expected bool
	}{
		{"main package", "main", true},
		{"main_test package", "main_test", true},
		{"test suffix", "mypackage_test", true},
		{"normal package", "mypackage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsExemptedPackage002(tt.pkgName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsExemptedPackage002(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage003 teste tous les packages exemptés pour Rule003.
func TestIsExemptedPackage003(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		expected bool
	}{
		{"main package", "main", true},
		{"main_test package", "main_test", true},
		{"test suffix", "mypackage_test", true},
		{"normal package", "mypackage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsExemptedPackage003(tt.pkgName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsExemptedPackage003(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage004 teste tous les packages exemptés pour Rule004.
func TestIsExemptedPackage004(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		expected bool
	}{
		{"main package", "main", true},
		{"main_test package", "main_test", true},
		{"test suffix", "mypackage_test", true},
		{"normal package", "mypackage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsExemptedPackage004(tt.pkgName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsExemptedPackage004(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage005 teste tous les packages exemptés pour Rule005.
func TestIsExemptedPackage005(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		expected bool
	}{
		{"main package", "main", true},
		{"main_test package", "main_test", true},
		{"test suffix", "mypackage_test", true},
		{"normal package", "mypackage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ktn_interface.IsExemptedPackage005(tt.pkgName)
			if got != tt.expected {
				t.Errorf("ktn_interface.IsExemptedPackage005(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestCollectPackageInfo teste la collecte d'informations de package.
func TestCollectPackageInfo(t *testing.T) {
	tests := []struct {
		name                      string
		src                       string
		expectedHasInterfacesFile bool
		expectedPublicStructs     int
		expectedPublicInterfaces  int
	}{
		{
			name: "package with interfaces.go and public struct",
			src: `package test

type PublicStruct struct {
	Field int
}

type PublicInterface interface {
	Method() error
}
`,
			expectedHasInterfacesFile: false,
			expectedPublicStructs:     1,
			expectedPublicInterfaces:  1,
		},
		{
			name: "package with only private types",
			src: `package test

type privateStruct struct {
	field int
}

type privateInterface interface {
	method() error
}
`,
			expectedHasInterfacesFile: false,
			expectedPublicStructs:     0,
			expectedPublicInterfaces:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, parser.ParseComments)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Pkg:   types.NewPackage("test", "test"),
			}

			info := ktn_interface.CollectPackageInfo(pass)

			if len(info.PublicStructs) != tt.expectedPublicStructs {
				t.Errorf("expected %d public structs, got %d", tt.expectedPublicStructs, len(info.PublicStructs))
			}

			if len(info.PublicInterfaces) != tt.expectedPublicInterfaces {
				t.Errorf("expected %d public interfaces, got %d", tt.expectedPublicInterfaces, len(info.PublicInterfaces))
			}
		})
	}
}

// TestRunRule001EdgeCases teste les cas edge de RunRule001.
func TestRunRule001EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := ktn_interface.RunRule001(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("package with only functions", func(t *testing.T) {
		src := `package test

// PublicFunction is a public exported function for testing purposes.
func PublicFunction() {}
func privateFunction() {}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for function-only package") },
		}

		_, err := ktn_interface.RunRule001(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule002EdgeCases teste les cas edge de RunRule002.
func TestRunRule002EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := ktn_interface.RunRule002(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("allowed public type", func(t *testing.T) {
		src := `package test

type UserID struct {
	value string
}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for allowed type") },
		}

		_, err := ktn_interface.RunRule002(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule003EdgeCases teste les cas edge de RunRule003.
func TestRunRule003EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := ktn_interface.RunRule003(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("interfaces.go file", func(t *testing.T) {
		src := `package test

type PublicInterface interface {
	Method() error
}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for interfaces.go") },
		}

		_, err := ktn_interface.RunRule003(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule004EdgeCases teste les cas edge de RunRule004.
func TestRunRule004EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := ktn_interface.RunRule004(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("interface with constructor", func(t *testing.T) {
		src := `package test

type Service interface {
	Method() error
}

// NewService creates and returns a new Service instance.
func NewService() Service {
	return nil
}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report when constructor exists") },
		}

		_, err := ktn_interface.RunRule004(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule005EdgeCases teste les cas edge de RunRule005.
func TestRunRule005EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := ktn_interface.RunRule005(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no interfaces.go file", func(t *testing.T) {
		src := `package test

type Service interface {
	Method() error
}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "service.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report when no interfaces.go exists") },
		}

		_, err := ktn_interface.RunRule005(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("interfaces.go with public interface", func(t *testing.T) {
		src := `package test

type Service interface {
	Method() error
}
`
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("test", "test"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report when interfaces.go has public interface") },
		}

		_, err := ktn_interface.RunRule005(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule001WindowsPaths teste les chemins Windows.
func TestRunRule001WindowsPaths(t *testing.T) {
	src := `package test

type PublicService struct {
	field int
}
`
	fset := token.NewFileSet()
	// Simuler un chemin Windows avec backslashes
	file, _ := parser.ParseFile(fset, "C:\\tests\\target\\service.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for tests/target path on Windows") },
	}

	_, err := ktn_interface.RunRule001(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule002WindowsPaths teste les chemins Windows pour Rule002.
func TestRunRule002WindowsPaths(t *testing.T) {
	src := `package test

type PublicService struct {
	field int
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "C:\\tests\\target\\service.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for tests/target path on Windows") },
	}

	_, err := ktn_interface.RunRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule002TestFile teste les fichiers _test.go.
func TestRunRule002TestFile(t *testing.T) {
	src := `package test

type PublicHelper struct {
	field int
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "helper_test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for _test.go files") },
	}

	_, err := ktn_interface.RunRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestCollectPackageInfoNonTypeDecl teste les déclarations non-type.
func TestCollectPackageInfoNonTypeDecl(t *testing.T) {
	src := `package test

const PublicConst = 42
var PublicVar = "test"

// PublicFunc is a public exported function for testing purposes.
func PublicFunc() {}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
	}

	info := ktn_interface.CollectPackageInfo(pass)

	if len(info.PublicStructs) != 0 {
		t.Errorf("expected 0 public structs, got %d", len(info.PublicStructs))
	}

	if len(info.PublicInterfaces) != 0 {
		t.Errorf("expected 0 public interfaces, got %d", len(info.PublicInterfaces))
	}
}

// TestCollectPackageInfoTypeAliases teste les alias de type.
func TestCollectPackageInfoTypeAliases(t *testing.T) {
	src := `package test

type PublicAlias = string
type PublicInt = int
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
	}

	info := ktn_interface.CollectPackageInfo(pass)

	// Les alias ne sont ni des structs ni des interfaces
	if len(info.PublicStructs) != 0 {
		t.Errorf("expected 0 public structs for aliases, got %d", len(info.PublicStructs))
	}

	if len(info.PublicInterfaces) != 0 {
		t.Errorf("expected 0 public interfaces for aliases, got %d", len(info.PublicInterfaces))
	}
}

// TestRunRule003PrivateInterface teste les interfaces privées.
func TestRunRule003PrivateInterface(t *testing.T) {
	src := `package test

type privateInterface interface {
	method() error
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "service.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for private interface") },
	}

	_, err := ktn_interface.RunRule003(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule004EmptyInterfaceNoConstructor teste les interfaces vides.
func TestRunRule004EmptyInterfaceNoConstructor(t *testing.T) {
	src := `package test

type EmptyInterface interface{}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for empty interface") },
	}

	_, err := ktn_interface.RunRule004(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule004NilMethods teste les interfaces avec Methods nil.
func TestRunRule004NilMethods(t *testing.T) {
	src := `package test

type MarkerInterface interface{}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for marker interface") },
	}

	_, err := ktn_interface.RunRule004(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule005PrivateInterfaceInInterfacesFile teste les interfaces privées dans interfaces.go.
func TestRunRule005PrivateInterfaceInInterfacesFile(t *testing.T) {
	src := `package test

type privateInterface interface {
	method() error
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)

	reported := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) {
			reported = true
		},
	}

	_, err := ktn_interface.RunRule005(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("should report when interfaces.go exists but has no public interfaces")
	}
}

// TestRunRule005PrivateStructInInterfacesFile teste les structs privées dans interfaces.go.
func TestRunRule005PrivateStructInInterfacesFile(t *testing.T) {
	src := `package test

type privateStruct struct {
	field int
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)

	reported := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) {
			reported = true
		},
	}

	_, err := ktn_interface.RunRule005(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("should report when interfaces.go exists but has no public interfaces")
	}
}

// TestRunRule002NonStructType teste les types non-struct.
func TestRunRule002NonStructType(t *testing.T) {
	src := `package test

type PublicAlias = string
type PublicFunc func() error
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "types.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for non-struct types") },
	}

	_, err := ktn_interface.RunRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule003NonInterfaceType teste les types non-interface.
func TestRunRule003NonInterfaceType(t *testing.T) {
	src := `package test

type PublicStruct struct {
	field int
}

type PublicAlias = string
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "types.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Pkg:    types.NewPackage("test", "test"),
		Report: func(diag analysis.Diagnostic) { t.Error("should not report for non-interface types") },
	}

	_, err := ktn_interface.RunRule003(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestCollectPackageInfoInterfacesFile teste la détection du fichier interfaces.go.
func TestCollectPackageInfoInterfacesFile(t *testing.T) {
	src := `package test

type PublicInterface interface {
	Method() error
}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "interfaces.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
	}

	info := ktn_interface.CollectPackageInfo(pass)

	if !info.HasInterfacesFile {
		t.Error("should detect interfaces.go file")
	}

	if len(info.PublicInterfaces) != 1 {
		t.Errorf("expected 1 public interface, got %d", len(info.PublicInterfaces))
	}
}
