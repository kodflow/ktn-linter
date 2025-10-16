package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// runStructAnalyzerValidTest exécute un test qui ne devrait pas produire d'erreur.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
func runStructAnalyzerValidTest(t *testing.T, name, code string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				t.Errorf("Unexpected diagnostic: %s at %s", diag.Message, fset.Position(diag.Pos))
			},
		}

		_, err = analyzer.StructAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}
	})
}

// runStructAnalyzerErrorTest exécute un test qui devrait produire une erreur spécifique.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - expectedError: code d'erreur attendu
func runStructAnalyzerErrorTest(t *testing.T, name, code, expectedError string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		foundExpectedError := false
		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				if expectedError != "" && !foundExpectedError {
					foundExpectedError = true
					t.Logf("Found expected error: %s", diag.Message)
				}
			},
		}

		_, err = analyzer.StructAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}

		if expectedError != "" && !foundExpectedError {
			t.Errorf("Expected error containing %q, but no errors were reported", expectedError)
		}
	})
}

// TestStructAnalyzerNaming teste le nommage des structs.
//
// Params:
//   - t: instance de test
func TestStructAnalyzerNaming(t *testing.T) {
	runStructAnalyzerValidTest(t, "Valid struct naming", `package test

// UserConfig contient la configuration.
type UserConfig struct {}
`)

	runStructAnalyzerErrorTest(t, "Invalid struct naming with underscore", `package test

// user_config contient la configuration.
type user_config struct {}
`, "KTN-STRUCT-001")
}

// TestStructAnalyzerDocumentation teste la documentation des structs.
//
// Params:
//   - t: instance de test
func TestStructAnalyzerDocumentation(t *testing.T) {
	runStructAnalyzerValidTest(t, "Struct with godoc", `package test

// User représente un utilisateur.
type User struct {
	// Name est le nom de l'utilisateur
	Name string
}
`)

	runStructAnalyzerErrorTest(t, "Struct without godoc", `package test

type User struct {
	// Name est le nom
	Name string
}
`, "KTN-STRUCT-002")
}

// TestStructAnalyzerFields teste la documentation des champs exportés.
//
// Params:
//   - t: instance de test
func TestStructAnalyzerFields(t *testing.T) {
	runStructAnalyzerValidTest(t, "Exported fields with comments", `package test

// User représente un utilisateur.
type User struct {
	// Name est le nom de l'utilisateur
	Name string

	// Email est l'adresse email
	Email string

	// Champ privé, pas de commentaire requis
	age int
}
`)

	runStructAnalyzerErrorTest(t, "Exported field without comment", `package test

// User représente un utilisateur.
type User struct {
	Name string
}
`, "KTN-STRUCT-003")

	runStructAnalyzerValidTest(t, "Private field without comment (OK)", `package test

// config contient la configuration.
type config struct {
	port int
	host string
}
`)
}

// TestStructAnalyzerFieldCount teste le nombre de champs.
//
// Params:
//   - t: instance de test
func TestStructAnalyzerFieldCount(t *testing.T) {
	runStructAnalyzerValidTest(t, "Struct with acceptable field count", `package test

// Config contient la configuration.
type Config struct {
	// Port est le port d'écoute
	Port int

	// Host est l'hôte
	Host string

	// Timeout est le délai d'expiration
	Timeout int
}
`)

	// Créer une struct avec trop de champs
	tooManyFieldsCode := `package test

// LargeStruct a trop de champs.
type LargeStruct struct {
	// Field1 description
	Field1 string
	// Field2 description
	Field2 string
	// Field3 description
	Field3 string
	// Field4 description
	Field4 string
	// Field5 description
	Field5 string
	// Field6 description
	Field6 string
	// Field7 description
	Field7 string
	// Field8 description
	Field8 string
	// Field9 description
	Field9 string
	// Field10 description
	Field10 string
	// Field11 description
	Field11 string
	// Field12 description
	Field12 string
	// Field13 description
	Field13 string
	// Field14 description
	Field14 string
	// Field15 description
	Field15 string
	// Field16 description
	Field16 string
}
`

	runStructAnalyzerErrorTest(t, "Struct with too many fields", tooManyFieldsCode, "KTN-STRUCT-004")
}

// TestStructAnalyzerEdgeCases teste les cas limites.
//
// Params:
//   - t: instance de test
func TestStructAnalyzerEdgeCases(t *testing.T) {
	runStructAnalyzerValidTest(t, "Empty struct", `package test

// Empty est une struct vide.
type Empty struct {}
`)

	runStructAnalyzerValidTest(t, "Struct with embedded field", `package test

// Base est la structure de base.
type Base struct {
	// ID est l'identifiant
	ID int
}

// Extended étend Base.
type Extended struct {
	Base
	// Name est le nom
	Name string
}
`)

	runStructAnalyzerValidTest(t, "Non-struct types ignored", `package test

// MyInt est un alias.
type MyInt int

// MyFunc est un type fonction.
type MyFunc func() error
`)
}
