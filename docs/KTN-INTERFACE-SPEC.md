# Spécification KTN-INTERFACE

## Vision

Tous les packages doivent définir leurs types publics comme interfaces pour garantir la testabilité via mocks.

## Principe Fondamental

**"Tout doit être mockable"** - Chaque package doit exposer des interfaces, pas des types concrets.

## Règles

### KTN-INTERFACE-001: Fichier interfaces.go obligatoire

**Sévérité**: Erreur
**Description**: Chaque package doit contenir un fichier `interfaces.go`

**Détection**:
```
package mypackage/
  ├── interfaces.go  ← OBLIGATOIRE
  ├── impl.go
  └── other.go
```

**Message d'erreur**:
```
[KTN-INTERFACE-001] Package 'mypackage' sans fichier interfaces.go
Créez interfaces.go pour définir les interfaces publiques du package.
Exemple:
  // interfaces.go
  package mypackage

  // MyService définit le contrat du service.
  type MyService interface {
      DoSomething(input string) (string, error)
  }
```

### KTN-INTERFACE-002: Types concrets publics interdits

**Sévérité**: Erreur
**Description**: Les types publics (exportés) doivent être des interfaces, pas des structs

**Détection**:
```go
// ❌ MAUVAIS - struct publique
type MyService struct {
    db *Database
}

// ✅ BON - interface publique
type MyService interface {
    DoSomething(input string) (string, error)
}

// ✅ BON - implémentation privée
type myServiceImpl struct {
    db Database
}
```

**Message d'erreur**:
```
[KTN-INTERFACE-002] Type public 'MyService' défini comme struct au lieu d'interface
Les types publics doivent être des interfaces dans interfaces.go.
Déplacez la struct vers une implémentation privée et créez l'interface.
Exemple:
  // interfaces.go
  type MyService interface {
      DoSomething(input string) (string, error)
  }

  // impl.go
  type myServiceImpl struct { ... }
```

### KTN-INTERFACE-003: Interface non utilisée

**Sévérité**: Warning
**Description**: Interface définie dans interfaces.go mais sans implémentation

**Détection**:
```go
// interfaces.go
type UnusedService interface {  // ⚠️ Aucune implémentation trouvée
    Method()
}
```

**Message d'erreur**:
```
[KTN-INTERFACE-003] Interface 'UnusedService' définie mais non implémentée
Supprimez l'interface ou créez son implémentation.
```

### KTN-INTERFACE-004: Implémentation publique interdite

**Sévérité**: Erreur
**Description**: Les implémentations doivent être privées (nom en minuscule)

**Détection**:
```go
// ❌ MAUVAIS - implémentation publique
type MyServiceImpl struct {
    db Database
}

// ✅ BON - implémentation privée
type myServiceImpl struct {
    db Database
}
```

**Message d'erreur**:
```
[KTN-INTERFACE-004] Implémentation 'MyServiceImpl' doit être privée
Les implémentations doivent commencer par une minuscule.
Renommez: MyServiceImpl → myServiceImpl
```

### KTN-INTERFACE-005: Interface mal placée

**Sévérité**: Erreur
**Description**: Interface définie ailleurs que dans interfaces.go

**Détection**:
```go
// service.go ❌ MAUVAIS
type MyService interface {
    Method()
}
```

**Message d'erreur**:
```
[KTN-INTERFACE-005] Interface 'MyService' définie dans service.go
Les interfaces publiques doivent être dans interfaces.go.
Déplacez cette interface vers interfaces.go.
```

### KTN-INTERFACE-006: Constructeur manquant

**Sévérité**: Warning
**Description**: Interface sans constructeur New{InterfaceName}()

**Détection**:
```go
// interfaces.go
type MyService interface {
    Method()
}

// ⚠️ Aucune fonction NewMyService() trouvée
```

**Message d'erreur**:
```
[KTN-INTERFACE-006] Interface 'MyService' sans constructeur
Ajoutez un constructeur qui retourne l'interface.
Exemple:
  // impl.go ou constructor.go
  func NewMyService(deps ...) MyService {
      return &myServiceImpl{...}
  }
```

## Exceptions

### Types Natifs Go Autorisés

Ces types peuvent être publics sans être des interfaces:
- Types primitifs: `type UserID string`
- Aliases de types: `type Count int`
- Types fonctions: `type Handler func(context.Context) error`
- Enums (iota)

### Packages Exemptés

Certains packages peuvent être exemptés:
- `main` package
- Packages de tests `*_test`
- Packages internes spécifiques (à définir)

## Architecture Recommandée

```
package mypackage/
  ├── interfaces.go      ← Toutes les interfaces publiques
  ├── impl.go           ← Implémentations privées (struct + méthodes)
  ├── constructor.go    ← Constructeurs New*() retournant interfaces
  ├── types.go          ← Types auxiliaires (si nécessaire)
  └── mypackage_test.go ← Tests avec mocks
```

## Exemple Complet

### interfaces.go
```go
package analyzer

// Analyzer définit le contrat d'un analyseur de code.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: résultat de l'analyse
//   - error: erreur éventuelle
type Analyzer interface {
    Run(pass *analysis.Pass) (interface{}, error)
    Name() string
    Doc() string
}

// FileReader lit les fichiers du système.
type FileReader interface {
    Read(path string) ([]byte, error)
}
```

### impl.go
```go
package analyzer

import "golang.org/x/tools/go/analysis"

// constAnalyzer implémente Analyzer pour les constantes.
type constAnalyzer struct {
    analyzer *analysis.Analyzer
}

func (c *constAnalyzer) Run(pass *analysis.Pass) (interface{}, error) {
    return c.analyzer.Run(pass)
}

func (c *constAnalyzer) Name() string {
    return c.analyzer.Name
}

func (c *constAnalyzer) Doc() string {
    return c.analyzer.Doc
}
```

### constructor.go
```go
package analyzer

import "golang.org/x/tools/go/analysis"

// NewConstAnalyzer crée un nouvel analyseur de constantes.
//
// Returns:
//   - Analyzer: l'analyseur configuré
func NewConstAnalyzer() Analyzer {
    return &constAnalyzer{
        analyzer: &analysis.Analyzer{
            Name: "ktnconst",
            Doc:  "Vérifie les constantes",
            Run:  runConstAnalyzer,
        },
    }
}
```

### Tests avec Mocks
```go
package analyzer_test

import "testing"

type mockAnalyzer struct {
    runFunc func(*analysis.Pass) (interface{}, error)
}

func (m *mockAnalyzer) Run(pass *analysis.Pass) (interface{}, error) {
    return m.runFunc(pass)
}

func TestWithMock(t *testing.T) {
    mock := &mockAnalyzer{
        runFunc: func(pass *analysis.Pass) (interface{}, error) {
            return nil, nil
        },
    }

    // Test avec le mock
    result, err := mock.Run(nil)
    // assertions...
}
```

## Bénéfices

1. **Testabilité**: Tous les types peuvent être mockés
2. **Découplage**: Les packages dépendent d'interfaces, pas d'implémentations
3. **Maintenabilité**: Contrats clairs et documentés
4. **Flexibilité**: Implémentations interchangeables
5. **SOLID**: Respecte l'Inversion de Dépendance (DIP)

## Migration

Pour migrer un package existant:

1. Créer `interfaces.go`
2. Extraire les types publics en interfaces
3. Renommer implémentations en privé
4. Créer constructeurs retournant interfaces
5. Mettre à jour les dépendances
6. Ajouter tests avec mocks

## Configuration

Dans `.ktn-linter.yaml`:
```yaml
interface:
  enabled: true
  strict: true
  exempted_packages:
    - main
    - internal/config
  require_constructors: true
  min_interface_methods: 1
```

## Références

- [Go Proverbs](https://go-proverbs.github.io/): "Accept interfaces, return structs"
- [Effective Go](https://golang.org/doc/effective_go#interfaces): Interfaces
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID): Dependency Inversion
