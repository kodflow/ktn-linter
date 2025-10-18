package ktn_interface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TestNeedsInterfacesFile teste les différents cas de needsInterfacesFile
func TestNeedsInterfacesFile(t *testing.T) {
	tests := []struct {
		name     string
		info     *packageInfo
		expected bool
	}{
		{
			name: "package avec struct publique non autorisée",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "MyService", fileName: "service.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: true,
		},
		{
			name: "package avec struct publique autorisée (Config)",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "MyConfig", fileName: "config.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec struct publique autorisée (ID)",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "UserID", fileName: "types.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec struct publique autorisée (Type)",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "ErrorType", fileName: "types.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec interface publique",
			info: &packageInfo{
				publicStructs:    []publicStruct{},
				publicInterfaces: []publicInterface{{name: "MyInterface"}},
			},
			expected: true,
		},
		{
			name: "package sans types publics",
			info: &packageInfo{
				publicStructs:    []publicStruct{},
				publicInterfaces: []publicInterface{},
			},
			expected: false,
		},
		{
			name: "package avec multiple structs dont une non autorisée",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "MyConfig", fileName: "config.go"},
					{name: "MyService", fileName: "service.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: true,
		},
		{
			name: "package avec structs toutes autorisées",
			info: &packageInfo{
				publicStructs: []publicStruct{
					{name: "MyConfig", fileName: "config.go"},
					{name: "UserID", fileName: "types.go"},
					{name: "ErrorType", fileName: "types.go"},
				},
				publicInterfaces: []publicInterface{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsInterfacesFile(tt.info)
			if got != tt.expected {
				t.Errorf("needsInterfacesFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestIsAllowedPublicType teste tous les suffixes autorisés
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
			got := isAllowedPublicType(tt.typeName)
			if got != tt.expected {
				t.Errorf("isAllowedPublicType(%q) = %v, want %v", tt.typeName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage teste tous les packages exemptés
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
			got := isExemptedPackage(tt.pkgName)
			if got != tt.expected {
				t.Errorf("isExemptedPackage(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage002 teste tous les packages exemptés pour Rule002
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
			got := isExemptedPackage002(tt.pkgName)
			if got != tt.expected {
				t.Errorf("isExemptedPackage002(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage003 teste tous les packages exemptés pour Rule003
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
			got := isExemptedPackage003(tt.pkgName)
			if got != tt.expected {
				t.Errorf("isExemptedPackage003(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage004 teste tous les packages exemptés pour Rule004
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
			got := isExemptedPackage004(tt.pkgName)
			if got != tt.expected {
				t.Errorf("isExemptedPackage004(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestIsExemptedPackage005 teste tous les packages exemptés pour Rule005
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
			got := isExemptedPackage005(tt.pkgName)
			if got != tt.expected {
				t.Errorf("isExemptedPackage005(%q) = %v, want %v", tt.pkgName, got, tt.expected)
			}
		})
	}
}

// TestCollectPackageInfo teste la collecte d'informations de package
func TestCollectPackageInfo(t *testing.T) {
	tests := []struct {
		name                     string
		src                      string
		expectedHasInterfacesFile bool
		expectedPublicStructs    int
		expectedPublicInterfaces int
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

			info := collectPackageInfo(pass)

			if len(info.publicStructs) != tt.expectedPublicStructs {
				t.Errorf("expected %d public structs, got %d", tt.expectedPublicStructs, len(info.publicStructs))
			}

			if len(info.publicInterfaces) != tt.expectedPublicInterfaces {
				t.Errorf("expected %d public interfaces, got %d", tt.expectedPublicInterfaces, len(info.publicInterfaces))
			}
		})
	}
}

// TestRunRule001_EdgeCases teste les cas edge de runRule001
func TestRunRule001_EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := runRule001(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("package with only functions", func(t *testing.T) {
		src := `package test

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

		_, err := runRule001(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule002_EdgeCases teste les cas edge de runRule002
func TestRunRule002_EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := runRule002(pass)
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

		_, err := runRule002(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule003_EdgeCases teste les cas edge de runRule003
func TestRunRule003_EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := runRule003(pass)
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

		_, err := runRule003(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule004_EdgeCases teste les cas edge de runRule004
func TestRunRule004_EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := runRule004(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("interface with constructor", func(t *testing.T) {
		src := `package test

type Service interface {
	Method() error
}

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

		_, err := runRule004(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule005_EdgeCases teste les cas edge de runRule005
func TestRunRule005_EdgeCases(t *testing.T) {
	t.Run("exempted package", func(t *testing.T) {
		fset := token.NewFileSet()
		file, _ := parser.ParseFile(fset, "test.go", "package main", parser.ParseComments)

		pass := &analysis.Pass{
			Fset:   fset,
			Files:  []*ast.File{file},
			Pkg:    types.NewPackage("main", "main"),
			Report: func(diag analysis.Diagnostic) { t.Error("should not report for main package") },
		}

		_, err := runRule005(pass)
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

		_, err := runRule005(pass)
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

		_, err := runRule005(pass)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestRunRule001_WindowsPaths teste les chemins Windows
func TestRunRule001_WindowsPaths(t *testing.T) {
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

	_, err := runRule001(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule002_WindowsPaths teste les chemins Windows pour Rule002
func TestRunRule002_WindowsPaths(t *testing.T) {
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

	_, err := runRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule002_TestFile teste les fichiers _test.go
func TestRunRule002_TestFile(t *testing.T) {
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

	_, err := runRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestCollectPackageInfo_NonTypeDecl teste les déclarations non-type
func TestCollectPackageInfo_NonTypeDecl(t *testing.T) {
	src := `package test

const PublicConst = 42
var PublicVar = "test"

func PublicFunc() {}
`
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("test", "test"),
	}

	info := collectPackageInfo(pass)

	if len(info.publicStructs) != 0 {
		t.Errorf("expected 0 public structs, got %d", len(info.publicStructs))
	}

	if len(info.publicInterfaces) != 0 {
		t.Errorf("expected 0 public interfaces, got %d", len(info.publicInterfaces))
	}
}

// TestCollectPackageInfo_TypeAliases teste les alias de type
func TestCollectPackageInfo_TypeAliases(t *testing.T) {
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

	info := collectPackageInfo(pass)

	// Les alias ne sont ni des structs ni des interfaces
	if len(info.publicStructs) != 0 {
		t.Errorf("expected 0 public structs for aliases, got %d", len(info.publicStructs))
	}

	if len(info.publicInterfaces) != 0 {
		t.Errorf("expected 0 public interfaces for aliases, got %d", len(info.publicInterfaces))
	}
}

// TestRunRule003_PrivateInterface teste les interfaces privées
func TestRunRule003_PrivateInterface(t *testing.T) {
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

	_, err := runRule003(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule004_EmptyInterfaceNoConstructor teste les interfaces vides
func TestRunRule004_EmptyInterfaceNoConstructor(t *testing.T) {
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

	_, err := runRule004(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule004_NilMethods teste les interfaces avec Methods nil
func TestRunRule004_NilMethods(t *testing.T) {
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

	_, err := runRule004(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule005_PrivateInterfaceInInterfacesFile teste les interfaces privées dans interfaces.go
func TestRunRule005_PrivateInterfaceInInterfacesFile(t *testing.T) {
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

	_, err := runRule005(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("should report when interfaces.go exists but has no public interfaces")
	}
}

// TestRunRule005_PrivateStructInInterfacesFile teste les structs privées dans interfaces.go
func TestRunRule005_PrivateStructInInterfacesFile(t *testing.T) {
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

	_, err := runRule005(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("should report when interfaces.go exists but has no public interfaces")
	}
}

// TestRunRule002_NonStructType teste les types non-struct
func TestRunRule002_NonStructType(t *testing.T) {
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

	_, err := runRule002(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRunRule003_NonInterfaceType teste les types non-interface
func TestRunRule003_NonInterfaceType(t *testing.T) {
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

	_, err := runRule003(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestCollectPackageInfo_InterfacesFile teste la détection du fichier interfaces.go
func TestCollectPackageInfo_InterfacesFile(t *testing.T) {
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

	info := collectPackageInfo(pass)

	if !info.hasInterfacesFile {
		t.Error("should detect interfaces.go file")
	}

	if len(info.publicInterfaces) != 1 {
		t.Errorf("expected 1 public interface, got %d", len(info.publicInterfaces))
	}
}
