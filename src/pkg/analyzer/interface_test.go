package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

// TestInterfaceAnalyzerMissingInterfacesFile teste la détection de l'absence de interfaces.go.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerMissingInterfacesFile(t *testing.T) {
	code := `package myservice

type MyService struct {
	db Database
}

func (s *MyService) DoSomething() error {
	return nil
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "service.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	foundError := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   types.NewPackage("example.com/myservice", "myservice"),
		Report: func(diag analysis.Diagnostic) {
			if !foundError {
				t.Logf("Found expected error: %s", diag.Message)
				foundError = true
			}
		},
	}

	_, err = analyzer.InterfaceAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Analyzer returned error: %v", err)
	}

	if !foundError {
		t.Error("Expected KTN-INTERFACE-001 error for missing interfaces.go, but got none")
	}
}

// TestInterfaceAnalyzerPublicStruct teste la détection de structs publiques.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerPublicStruct(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldError bool
		description string
	}{
		{
			name: "Public struct should error",
			code: `package myservice

type MyService struct {
	field string
}
`,
			shouldError: true,
			description: "Public struct without interface",
		},
		{
			name: "Config struct allowed",
			code: `package myservice

type MyConfig struct {
	host string
	port int
}
`,
			shouldError: false,
			description: "Config suffix is allowed",
		},
		{
			name: "ID type allowed",
			code: `package myservice

type UserID string
`,
			shouldError: false,
			description: "ID suffix is allowed",
		},
		{
			name: "Private struct OK",
			code: `package myservice

type myServiceImpl struct {
	field string
}
`,
			shouldError: false,
			description: "Private struct is OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			// Ajouter un fichier interfaces.go vide pour éviter INTERFACE-001
			interfacesFile, _ := parser.ParseFile(fset, "interfaces.go", "package myservice", parser.ParseComments)
			file, err := parser.ParseFile(fset, "service.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			foundError := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{interfacesFile, file},
				Pkg:   types.NewPackage("example.com/myservice", "myservice"),
				Report: func(diag analysis.Diagnostic) {
					// Ignorer INTERFACE-001 et INTERFACE-007 (fichier vide pour les tests)
					if !containsInterface(diag.Message, "KTN-INTERFACE-001") &&
						!containsInterface(diag.Message, "KTN-INTERFACE-007") {
						foundError = true
						t.Logf("Found error: %s", diag.Message)
					}
				},
			}

			_, err = analyzer.InterfaceAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}

			if tt.shouldError && !foundError {
				t.Errorf("Expected error for %s, but got none", tt.description)
			} else if !tt.shouldError && foundError {
				t.Errorf("Unexpected error for %s", tt.description)
			}
		})
	}
}

// TestInterfaceAnalyzerInterfaceInWrongFile teste la détection d'interfaces mal placées.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerInterfaceInWrongFile(t *testing.T) {
	code := `package myservice

type MyService interface {
	DoSomething() error
}
`

	fset := token.NewFileSet()
	// Ajouter interfaces.go vide
	interfacesFile, _ := parser.ParseFile(fset, "interfaces.go", "package myservice", parser.ParseComments)
	// Interface dans le mauvais fichier
	file, err := parser.ParseFile(fset, "service.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	foundError := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{interfacesFile, file},
		Pkg:   types.NewPackage("example.com/myservice", "myservice"),
		Report: func(diag analysis.Diagnostic) {
			if containsInterface(diag.Message, "KTN-INTERFACE-005") {
				foundError = true
				t.Logf("Found expected error: %s", diag.Message)
			}
		},
	}

	_, err = analyzer.InterfaceAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Analyzer returned error: %v", err)
	}

	if !foundError {
		t.Error("Expected KTN-INTERFACE-005 error for interface in wrong file, but got none")
	}
}

// TestInterfaceAnalyzerMissingConstructor teste la détection de constructeurs manquants.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerMissingConstructor(t *testing.T) {
	tests := []struct {
		name         string
		interfaceCode string
		implCode      string
		shouldWarn    bool
	}{
		{
			name: "Interface with methods needs constructor",
			interfaceCode: `package myservice

type MyService interface {
	DoSomething() error
}
`,
			implCode: `package myservice

type myServiceImpl struct {}

func (s *myServiceImpl) DoSomething() error {
	return nil
}
`,
			shouldWarn: true,
		},
		{
			name: "Constructor exists - OK",
			interfaceCode: `package myservice

type MyService interface {
	DoSomething() error
}
`,
			implCode: `package myservice

type myServiceImpl struct {}

func NewMyService() MyService {
	return &myServiceImpl{}
}

func (s *myServiceImpl) DoSomething() error {
	return nil
}
`,
			shouldWarn: false,
		},
		{
			name: "Empty interface no warning",
			interfaceCode: `package myservice

type MyMarker interface {}
`,
			implCode:   `package myservice`,
			shouldWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			interfacesFile, _ := parser.ParseFile(fset, "interfaces.go", tt.interfaceCode, parser.ParseComments)
			implFile, _ := parser.ParseFile(fset, "impl.go", tt.implCode, parser.ParseComments)

			foundWarning := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{interfacesFile, implFile},
				Pkg:   types.NewPackage("example.com/myservice", "myservice"),
				Report: func(diag analysis.Diagnostic) {
					if containsInterface(diag.Message, "KTN-INTERFACE-006") {
						foundWarning = true
						t.Logf("Found warning: %s", diag.Message)
					}
				},
			}

			_, err := analyzer.InterfaceAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}

			if tt.shouldWarn && !foundWarning {
				t.Error("Expected KTN-INTERFACE-006 warning, but got none")
			} else if !tt.shouldWarn && foundWarning {
				t.Error("Unexpected warning")
			}
		})
	}
}

// TestInterfaceAnalyzerCompliantPackage teste un package conforme.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerCompliantPackage(t *testing.T) {
	interfacesCode := `package myservice

// MyService définit le contrat du service.
type MyService interface {
	DoSomething(input string) (string, error)
}

// Repository gère la persistance.
type Repository interface {
	Save(data string) error
	Load(id string) (string, error)
}
`

	implCode := `package myservice

type myServiceImpl struct {
	repo Repository
}

func NewMyService(repo Repository) MyService {
	return &myServiceImpl{repo: repo}
}

func (s *myServiceImpl) DoSomething(input string) (string, error) {
	return "result", nil
}

type repositoryImpl struct {
	db string
}

func NewRepository(db string) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Save(data string) error {
	return nil
}

func (r *repositoryImpl) Load(id string) (string, error) {
	return "", nil
}
`

	fset := token.NewFileSet()
	interfacesFile, _ := parser.ParseFile(fset, "interfaces.go", interfacesCode, parser.ParseComments)
	implFile, _ := parser.ParseFile(fset, "impl.go", implCode, parser.ParseComments)

	hasError := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{interfacesFile, implFile},
		Pkg:   types.NewPackage("example.com/myservice", "myservice"),
		Report: func(diag analysis.Diagnostic) {
			// Ne devrait pas y avoir d'erreurs
			hasError = true
			t.Errorf("Unexpected diagnostic: %s", diag.Message)
		},
	}

	_, err := analyzer.InterfaceAnalyzer.Run(pass)
	if err != nil {
		t.Errorf("Analyzer returned error: %v", err)
	}

	if hasError {
		t.Error("Compliant package should not have errors")
	}
}

// TestInterfaceAnalyzerExemptedPackages teste les packages exemptés.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerExemptedPackages(t *testing.T) {
	tests := []struct {
		pkgName     string
		shouldCheck bool
	}{
		{"main", false},
		{"myservice", true},
		{"myservice_test", false},
		{"main_test", false},
	}

	for _, tt := range tests {
		t.Run(tt.pkgName, func(t *testing.T) {
			code := `package ` + tt.pkgName + `

type MyService struct {
	field string
}
`

			fset := token.NewFileSet()
			file, _ := parser.ParseFile(fset, "service.go", code, parser.ParseComments)

			hasError := false
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{file},
				Pkg:   types.NewPackage("example.com/"+tt.pkgName, tt.pkgName),
				Report: func(diag analysis.Diagnostic) {
					hasError = true
				},
			}

			_, err := analyzer.InterfaceAnalyzer.Run(pass)
			if err != nil {
				t.Errorf("Analyzer returned error: %v", err)
			}

			if tt.shouldCheck && !hasError {
				t.Errorf("Expected errors for package %s, but got none", tt.pkgName)
			} else if !tt.shouldCheck && hasError {
				t.Errorf("Package %s should be exempted, but got errors", tt.pkgName)
			}
		})
	}
}

// TestInterfaceAnalyzerRealWorldScenarios teste des scénarios réels.
//
// Params:
//   - t: instance de test
func TestInterfaceAnalyzerRealWorldScenarios(t *testing.T) {
	t.Run("HTTP Handler Pattern", func(t *testing.T) {
		interfacesCode := `package handler

// Handler traite les requêtes HTTP.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// UserService gère les utilisateurs.
type UserService interface {
	GetUser(id string) (*UserData, error)
	CreateUser(user *UserData) error
}

// UserData représente les données d'un utilisateur (struct de données avec suffixe autorisé).
type UserData struct {
	ID   string
	Name string
}
`

		implCode := `package handler

type handlerImpl struct {
	userService UserService
}

func NewHandler(userService UserService) Handler {
	return &handlerImpl{userService: userService}
}

func (h *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// implementation
}

type userServiceImpl struct {}

func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) GetUser(id string) (*UserData, error) {
	return nil, nil
}

func (u *userServiceImpl) CreateUser(user *UserData) error {
	return nil
}
`

		fset := token.NewFileSet()
		interfacesFile, _ := parser.ParseFile(fset, "interfaces.go", interfacesCode, parser.ParseComments)
		implFile, _ := parser.ParseFile(fset, "impl.go", implCode, parser.ParseComments)

		hasError := false
		pass := &analysis.Pass{
			Fset:  fset,
			Files: []*ast.File{interfacesFile, implFile},
			Pkg:   types.NewPackage("example.com/handler", "handler"),
			Report: func(diag analysis.Diagnostic) {
				hasError = true
				t.Errorf("Unexpected error: %s", diag.Message)
			},
		}

		_, err := analyzer.InterfaceAnalyzer.Run(pass)
		if err != nil {
			t.Errorf("Analyzer returned error: %v", err)
		}

		if hasError {
			t.Error("HTTP handler pattern should be valid")
		}
	})
}

// containsInterface vérifie si une chaîne contient une sous-chaîne.
//
// Params:
//   - s: la chaîne à analyser
//   - substr: la sous-chaîne recherchée
//
// Returns:
//   - bool: true si substr est trouvé dans s
func containsInterface(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && containsInterfaceHelper(s, substr)
}

// containsInterfaceHelper est une fonction helper pour rechercher une sous-chaîne.
//
// Params:
//   - s: la chaîne à analyser
//   - substr: la sous-chaîne recherchée
//
// Returns:
//   - bool: true si substr est trouvé dans s
func containsInterfaceHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
