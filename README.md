# KTN-Linter

Linter Go personnalisé pour appliquer les bonnes pratiques Kodflow.

## Vue d'ensemble

KTN-Linter vérifie automatiquement que votre code Go respecte les standards Kodflow.

**Formats de sortie :**
- **Format humain** (défaut) : Sortie colorée et structurée
- **Mode IA** (`-ai`) : Format Markdown pour Claude, ChatGPT
- **Mode simple** (`-simple`) : Une ligne par erreur pour IDE/VSCode
- **Sans couleurs** (`-no-color`) : Pour CI/CD et logs

**Règles implémentées :**
- ✅ **Constantes (package-level)** : Regroupement, documentation et typage explicite
- ✅ **Variables (package-level)** : Regroupement, documentation, typage et nommage
- ✅ **Fonctions (natives)** : Nommage, documentation stricte, complexité, longueur, profondeur
- ✅ **Structures** : Nommage MixedCaps, documentation, limitation champs
- ✅ **Interfaces** : Design interface-first, constructeurs obligatoires, fichiers dédiés
- ✅ **Tests** : Package naming, couverture fichiers, documentation complète

**Tests de validation :**
- 🎯 **tests/target/** : 0 violation - Code PARFAIT conforme à toutes les règles
- 🔴 **tests/source/** : 424 violations - Catalogue complet d'anti-patterns

---

## Installation

### Prérequis

- **Go 1.23+** (requis)
- **golangci-lint** (optionnel)

### Installation rapide

```bash
# 1. Vérifier Go
go version

# 2. Installer les dépendances
make deps

# 3. Compiler
make build

# 4. Tester
./builds/ktn-linter --help
```

---

## Utilisation

### Mode standalone

```bash
# Analyser un fichier
./builds/ktn-linter ./path/to/file.go

# Analyser un package
./builds/ktn-linter ./pkg/...

# Analyser tout le projet
./builds/ktn-linter ./...
```

### Options

```bash
# Mode IA (pour Claude/ChatGPT)
./builds/ktn-linter -ai ./...

# Mode simple (pour IDE/VSCode)
./builds/ktn-linter -simple ./...

# Sans couleurs (pour CI/CD)
./builds/ktn-linter -no-color ./...

# Verbose
./builds/ktn-linter -v ./...
```

### Avec VSCode (intégration automatique)

Le projet utilise un wrapper qui exécute uniquement KTN-Linter.

```bash
# Analyser avec le wrapper
./bin/golangci-lint-wrapper run ./...

# Dans VSCode, le wrapper est automatiquement utilisé
# Les erreurs apparaissent lors de la sauvegarde (Ctrl+S)
```

---

## Commandes Make

```bash
make help            # Aide
make deps            # Installer dépendances
make build           # Compiler linter
make lint            # Tester sur tests/
make test            # Tests unitaires
make clean           # Nettoyer binaires
make install-tools   # Installer golangci-lint
```

---

## Structure du projet

```
.
├── bin/
│   └── golangci-lint-wrapper    # Wrapper pour KTN-Linter
├── src/
│   ├── cmd/ktn-linter/          # Linter standalone
│   ├── pkg/analyzer/            # Analyseurs (const.go, var.go, ...)
│   │   └── formatter/           # Formatage sortie
│   ├── internal/                # Packages internes
│   │   ├── astutil/             # Utilitaires AST
│   │   ├── naming/              # Validation nommage
│   │   └── messageutil/         # Extraction messages
│   └── plugin/                  # Plugin module (pour future intégration)
├── tests/
│   ├── source/                  # Code avec violations - Anti-patterns
│   │   ├── README.md            # Guide des anti-patterns
│   │   ├── rules_const/         # Constantes mal déclarées
│   │   ├── rules_var/           # Variables anarchiques
│   │   ├── rules_func/          # Fonctions catastrophiques
│   │   ├── rules_struct/        # Structures incorrectes
│   │   ├── rules_interface/     # Design anti-patterns
│   │   └── rules_test/          # Tests inadéquats
│   └── target/                  # Code avec 0 violation - Perfection
│       ├── rules_const/         # Constantes parfaites
│       ├── rules_var/           # Variables optimales
│       ├── rules_func/          # Fonctions exemplaires
│       ├── rules_struct/        # Structures conformes
│       ├── rules_interface/     # Interface-first design
│       └── rules_test/          # Tests complets
├── .vscode/
│   ├── settings.json            # Config VSCode + wrapper
│   └── extensions.json          # Extension Go recommandée
├── .golangci.yml                # Config minimale (wrapper uniquement)
├── go.mod
├── Makefile
└── README.md
```

**Architecture des tests :**
- **Dualité parfaite** :
  - `tests/target/` : Code PARFAIT avec 0 violation (référence de qualité)
  - `tests/source/` : Code FOIREUX avec 405 violations (ce qu'il NE FAUT PAS faire)
- **Couverture complète** : Tous les scénarios, edge cases et anti-patterns
- **Validation bidirectionnelle** :
  - target/ prouve que le bon code passe ✅
  - source/ prouve que le mauvais code est détecté ❌

---

## Règles détaillées

### 📦 Constantes Package-Level (KTN-CONST-XXX)

Les constantes doivent être **regroupées** dans des blocs `const ()`, **documentées** et **typées explicitement**.

| Code | Description | Exemple |
|------|-------------|---------|
| `KTN-CONST-001` | Constante non groupée dans `const ()` | ❌ `const MaxRetries = 3`<br>✅ `const ( MaxRetries int = 3 )` |
| `KTN-CONST-002` | Groupe sans commentaire | ❌ `const ( ... )`<br>✅ `// Config constants`<br>`const ( ... )` |
| `KTN-CONST-003` | Constante sans commentaire individuel | ❌ `MaxRetries int = 3`<br>✅ `// MaxRetries ...`<br>`MaxRetries int = 3` |
| `KTN-CONST-004` | Constante sans type explicite | ❌ `MaxRetries = 3`<br>✅ `MaxRetries int = 3` |

**Exemple complet :**
```go
// Configuration constants.
// Define application limits and defaults.
const (
    // MaxConnections nombre maximum de connexions simultanées.
    MaxConnections int = 100

    // DefaultTimeout timeout par défaut en secondes.
    DefaultTimeout int = 30
)
```

**Exception iota :** Type explicite uniquement sur la première constante :
```go
// HTTP methods
const (
    // MethodGet représente GET
    MethodGet int = iota
    // MethodPost représente POST
    MethodPost
)
```

Documentation complète : [tests/source/rules_const/.README.md](./tests/source/rules_const/.README.md)

---

### 📝 Variables Package-Level (KTN-VAR-XXX)

Les variables doivent être **regroupées**, **documentées**, **typées explicitement** et suivre **MixedCaps**.

| Code | Description | Exemple |
|------|-------------|---------|
| `KTN-VAR-001` | Variable non groupée dans `var ()` | ❌ `var Port = 8080`<br>✅ `var ( Port int = 8080 )` |
| `KTN-VAR-002` | Groupe sans commentaire | ❌ `var ( ... )`<br>✅ `// HTTP config`<br>`var ( ... )` |
| `KTN-VAR-003` | Variable sans commentaire individuel | ❌ `Port int = 8080`<br>✅ `// Port ...`<br>`Port int = 8080` |
| `KTN-VAR-004` | Variable sans type explicite | ❌ `Port = 8080`<br>✅ `Port int = 8080` |
| `KTN-VAR-005` | Variable devrait être const | ❌ `var Pi = 3.14`<br>✅ `const Pi float64 = 3.14` |
| `KTN-VAR-006` | Multiple variables sur une ligne | ❌ `Host, Port = "localhost", 8080`<br>✅ Lignes séparées |
| `KTN-VAR-007` | Channel sans buffer size explicite | ❌ `Queue chan string`<br>✅ `// Queue (buffer=100)`<br>`Queue chan string = make(chan string, 100)` |
| `KTN-VAR-008` | Nom avec underscore | ❌ `max_retries`<br>✅ `maxRetries` |
| `KTN-VAR-009` | Nom en ALL_CAPS | ❌ `MAX_RETRIES`<br>✅ `MaxRetries` |

**Exemple complet :**
```go
// HTTP configuration
// Ces variables configurent le serveur HTTP
var (
    // Port est le port d'écoute du serveur
    Port int = 8080

    // Timeout est le délai d'expiration des requêtes
    Timeout int = 30

    // RequestQueue canal des requêtes entrantes (buffer=1000)
    RequestQueue chan Request = make(chan Request, 1000)
)
```

Documentation complète : [tests/source/rules_var/.README.md](./tests/source/rules_var/.README.md)

---

### ⚡ Fonctions Natives (KTN-FUNC-XXX)

Les fonctions doivent respecter des standards stricts de **nommage**, **documentation** et **complexité**.

| Code | Description | Seuil |
|------|-------------|-------|
| `KTN-FUNC-001` | Nom pas en MixedCaps/mixedCaps | ❌ snake_case interdit |
| `KTN-FUNC-002` | Fonction sans godoc | Toutes (exportées ET privées) |
| `KTN-FUNC-003` | Paramètres non documentés | Si > 0 params |
| `KTN-FUNC-004` | Retours non documentés | Si > 0 retours |
| `KTN-FUNC-005` | Trop de paramètres | > 5 paramètres |
| `KTN-FUNC-006` | Fonction trop longue | > 35 lignes |
| `KTN-FUNC-007` | Complexité cyclomatique trop élevée | ≥ 10 |
| `KTN-FUNC-008` | Return sans commentaire explicatif | Tous les returns |
| `KTN-FUNC-010` | Profondeur d'imbrication trop élevée | > 3 niveaux |

**Format godoc obligatoire (avec Params/Returns) :**
```go
// ProcessUser traite les données utilisateur et les valide.
//
// Params:
//   - user: les données utilisateur à traiter
//   - options: options de traitement
//
// Returns:
//   - *Result: résultat du traitement
//   - error: erreur si la validation échoue
func ProcessUser(user *User, options ProcessOptions) (*Result, error) {
    // Validation des données
    if err := validateUser(user); err != nil {
        // Erreur de validation (KTN-FUNC-008 : commentaire obligatoire)
        return nil, err
    }

    // Traitement
    result := &Result{
        Status: "processed",
        User:   user,
    }

    // Succès (KTN-FUNC-008 : commentaire obligatoire)
    return result, nil
}
```

**KTN-FUNC-008 : Commentaires sur returns (NOUVEAU) :**
```go
// ❌ INTERDIT - Returns sans commentaires
func Process() error {
    if err != nil {
        return err    // ❌ Pas de commentaire
    }
    return nil        // ❌ Pas de commentaire
}

// ✅ REQUIS - Tous les returns commentés
func Process() error {
    if err != nil {
        // Erreur de traitement
        return err
    }
    // Succès
    return nil
}
```

**Règles strictes :**
- **≤ 5 paramètres** : Utiliser struct si plus
- **≤ 35 lignes** : Extraire des sous-fonctions si plus
- **Complexité < 10** : Simplifier la logique (moins de if/for/switch)
- **Profondeur ≤ 3** : Utiliser early returns et helpers

Documentation complète : [tests/source/rules_func/.README.md](./tests/source/rules_func/.README.md)

---

### 🏗️ Structures (KTN-STRUCT-XXX)

Les structures doivent respecter le **nommage MixedCaps**, être **documentées** et avoir un **nombre raisonnable de champs**.

| Code | Description | Seuil |
|------|-------------|-------|
| `KTN-STRUCT-001` | Nom pas en MixedCaps/mixedCaps | ❌ snake_case interdit |
| `KTN-STRUCT-002` | Struct sans godoc | Toutes les structs |
| `KTN-STRUCT-003` | Champs exportés non documentés | Tous les champs exportés |
| `KTN-STRUCT-004` | Trop de champs | > 15 champs |

**Exemple complet conforme :**
```go
// UserConfig contient la configuration utilisateur pour l'application.
type UserConfig struct {
    // Host est l'hôte du serveur
    Host string
    // Port est le port d'écoute
    Port int
    // Timeout est le délai d'expiration en secondes
    Timeout int

    // Champs privés n'ont pas besoin de commentaires
    internalState bool
}

// Struct privée (mixedCaps)
type privateConfig struct {
    apiKey string
    secret string
}
```

**KTN-STRUCT-004 : Décomposition des structs larges :**
```go
// ❌ INTERDIT - Trop de champs (>15)
type ComplexUser struct {
    ID, Username, Email, FirstName, LastName string
    PhoneNumber, Address, City, ZipCode, Country string
    BirthDate string
    CreatedAt, UpdatedAt, LastLoginAt int64
    IsActive, IsVerified bool
    Roles, Permissions []string
    // ... 18 champs total > 15
}

// ✅ REQUIS - Décomposer en sous-structs
type UserBasicInfo struct {
    // ID est l'identifiant unique
    ID string
    // Username est le nom d'utilisateur
    Username string
    // Email est l'adresse email
    Email string
}

type UserContactInfo struct {
    // PhoneNumber est le numéro de téléphone
    PhoneNumber string
    // Address est l'adresse complète
    Address string
}

// User représente un utilisateur complet.
type User struct {
    // BasicInfo contient les informations de base
    BasicInfo UserBasicInfo
    // ContactInfo contient les informations de contact
    ContactInfo UserContactInfo
}
```

**Points clés :**
- ✅ **MixedCaps** pour structs exportées, **mixedCaps** pour privées
- ✅ **Godoc obligatoire** sur toutes les structs
- ✅ **Commentaires sur champs exportés** uniquement
- ✅ **Maximum 15 champs** : décomposer si plus

Documentation complète : [tests/source/rules_struct/.KTN-STRUCT.md](./tests/source/rules_struct/.KTN-STRUCT.md)

---

### 🔌 Interfaces (KTN-INTERFACE-XXX)

Design **interface-first** : types publics = interfaces, implémentations privées.

| Code | Description | Solution |
|------|-------------|----------|
| `KTN-INTERFACE-001` | Package sans fichier interfaces.go | Créer `interfaces.go` |
| `KTN-INTERFACE-002` | Type public struct au lieu d'interface | Exposer interface, struct privée |
| `KTN-INTERFACE-003` | Godoc incomplet sur interface | Ajouter Params/Returns |
| `KTN-INTERFACE-004` | Godoc incomplet sur méthode | Documenter params/returns |
| `KTN-INTERFACE-005` | Interface vide ou interface{} | Définir méthodes concrètes |
| `KTN-INTERFACE-006` | Interface sans constructeur New* | Créer `NewXxx()` |
| `KTN-INTERFACE-007` | Package sans types publics | Exposer au moins une interface |
| `KTN-INTERFACE-008` | **interfaces.go contient des structs** | **Uniquement interfaces dans interfaces.go** |
| `KTN-MOCK-001` | **mock.go manquant alors que interfaces existent** | **Créer mock.go avec build tag** |

**Pattern obligatoire (KTN-INTERFACE-008) :**
```go
// interfaces.go - UNIQUEMENT des interfaces (KTN-INTERFACE-008)
package myservice

// Service définit l'interface du service.
type Service interface {
    Process(data string) error
    GetStatus() string
}

// impl.go (même package)
package myservice

// service implémentation privée.
type service struct {
    db Database
}

// NewService crée une nouvelle instance de Service.
//
// Params:
//   - db: base de données à utiliser
//
// Returns:
//   - Service: instance du service
func NewService(db Database) Service {
    return &service{db: db}
}

// Process implémente Service.Process.
func (s *service) Process(data string) error {
    return s.db.Save(data)
}

// mock.go - Mocks réutilisables (KTN-MOCK-001)
//go:build test
// +build test

package myservice

// MockService mock réutilisable de Service.
type MockService struct {
    ProcessFunc    func(data string) error
    GetStatusFunc  func() string
}

func (m *MockService) Process(data string) error {
    if m.ProcessFunc != nil {
        return m.ProcessFunc(data)
    }
    return nil
}

func (m *MockService) GetStatus() string {
    if m.GetStatusFunc != nil {
        return m.GetStatusFunc()
    }
    return ""
}

var _ Service = (*MockService)(nil)
```

**Nouvelles règles strictes (depuis v1.1) :**
- ✅ **KTN-INTERFACE-008** : `interfaces.go` ne contient **QUE** des interfaces (pas de structs)
- ✅ **KTN-MOCK-001** : `mock.go` **obligatoire** si `interfaces.go` a des interfaces
- ✅ **mock.go avec `//go:build test`** pour exclure du build production
- ✅ Mocks réutilisables dans **tous les packages** de test

**Bénéfices :**
- ✅ **Testabilité** : Interfaces mockables + mocks centralisés
- ✅ **Découplage** : Dépendances sur contrats, pas implémentations
- ✅ **Flexibilité** : Implémentations interchangeables
- ✅ **Maintenabilité** : Un seul fichier mock.go par package

---

### 🧪 Tests (KTN-TEST-XXX)

Tests avec **package_test**, fichiers dédiés, **mocks réutilisables** et documentation complète.

| Code | Description | Solution |
|------|-------------|----------|
| `KTN-TEST-001` | Package incorrect | Utiliser `package xxx_test` |
| `KTN-TEST-002` | Fichier sans test | Créer `xxx_test.go` |
| `KTN-TEST-003` | Test sans fichier source | Créer fichier source ou déplacer test |
| `KTN-TEST-004` | Fonction test hors fichier `*_test.go` | Déplacer vers `*_test.go` |

**Architecture de tests :**
```
package myservice/
  ├── interfaces.go      # Interfaces publiques UNIQUEMENT
  ├── impl.go           # Implémentations privées
  ├── mock.go           # Mocks réutilisables (build tag test)
  ├── service_test.go   # Tests (package myservice_test)
  └── helper_test.go    # Tests helpers (package myservice_test)
```

**Fichier mock.go (mocks réutilisables) :**
```go
//go:build test
// +build test

package myservice

// MockService est un mock réutilisable de Service.
type MockService struct {
    ProcessFunc    func(data string) error
    GetStatusFunc  func() string
}

// Process implémente Service.Process.
func (m *MockService) Process(data string) error {
    if m.ProcessFunc != nil {
        return m.ProcessFunc(data)
    }
    return nil
}

// GetStatus implémente Service.GetStatus.
func (m *MockService) GetStatus() string {
    if m.GetStatusFunc != nil {
        return m.GetStatusFunc()
    }
    return ""
}

// Vérification à la compilation
var _ Service = (*MockService)(nil)
```

**Utilisation dans les tests :**
```go
// service_test.go
package myservice_test

import (
    "testing"
    "myproject/myservice"
)

// TestCreateUser vérifie la création d'utilisateur.
//
// Params:
//   - t: contexte de test
func TestProcessWithMock(t *testing.T) {
    // Arrange - utiliser le mock depuis mock.go
    mock := &myservice.MockService{
        ProcessFunc: func(data string) error {
            return nil
        },
    }

    // Act
    err := mock.Process("test data")

    // Assert
    if err != nil {
        t.Errorf("Process() error = %v", err)
    }
}
```

**Règles :**
- ✅ Package `xxx_test` pour isolation
- ✅ Un fichier `*_test.go` par fichier source
- ✅ **`mock.go` avec build tag `//go:build test`** pour mocks réutilisables
- ✅ **Exception : `mock.go` n'a pas besoin de `mock_test.go`**
- ✅ Godoc avec section Params sur tous les tests
- ✅ Mocks exportés (struct Mock*) utilisables depuis tous les packages

---

## Ajouter une règle

1. **Créer la structure de test :**
   ```bash
   mkdir -p tests/source/rules_<nom>
   mkdir -p tests/target/rules_<nom>
   ```

2. **Créer les fichiers :**
   - `tests/source/rules_<nom>/.README.md` : Documentation
   - `tests/source/rules_<nom>/package-level.go` : Code incorrect (ou autre nom descriptif)
   - `tests/target/rules_<nom>/package-level.go` : Code correct (ou autre nom descriptif)

3. **Implémenter l'analyseur :**
   - `src/pkg/analyzer/<nom>.go`

4. **Enregistrer l'analyseur :**
   - Dans `src/cmd/ktn-linter/main.go`
   - Dans `src/plugin/plugin.go`

5. **Tester :**
   ```bash
   make lint
   ```

Le Makefile analyse automatiquement tous les dossiers dans `tests/source/` et `tests/target/`.

---

## Intégration CI/CD

### GitHub Actions

```yaml
- name: Setup Go
  uses: actions/setup-go@v4
  with:
    go-version: '1.23'

- name: Run KTN-Linter
  run: |
    make build
    ./builds/ktn-linter ./...
```

### GitLab CI

```yaml
lint:
  script:
    - make build
    - ./builds/ktn-linter ./...
```

### Pre-commit hook

`.git/hooks/pre-commit` :

```bash
#!/bin/sh
./builds/ktn-linter ./... || exit 1
```

### VSCode Integration

Le projet est pré-configuré pour VSCode avec intégration automatique :

**Configuration :**
- `.vscode/settings.json` : Configure le wrapper KTN-Linter
- `.golangci.yml` : Configuration minimale
- Le wrapper exécute uniquement KTN-Linter (pas de linters golangci-lint)

**Utilisation :**
- Les erreurs apparaissent automatiquement lors de la sauvegarde (`Ctrl+S`)
- Le panel **PROBLEMS** affiche toutes les erreurs avec liens cliquables
- Les erreurs sont préfixées par `[KTN-CONST-XXX]`

**Prérequis :**
1. Installer l'extension Go : `Ctrl+Shift+X` → rechercher "Go"
2. Compiler le linter : `make build`

---

## Troubleshooting

**Go non installé :**
```bash
# Installer : https://go.dev/doc/install
go version
```

**golangci-lint non installé :**
```bash
make install-tools
```

**Wrapper ne trouve pas ktn-linter :**
```bash
make clean && make build
```

**Erreurs ne s'affichent pas dans VSCode :**
```bash
# Vérifier que golangci-lint v2+ est installé
golangci-lint --version

# Recompiler le linter
make build
```

---

## Licence

À définir
