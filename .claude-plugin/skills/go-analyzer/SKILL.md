# Go Code Analyzer Skill

## Description

Analyse automatique du code Go avec ktn-linter et golangci-lint pour détecter les violations de conventions, bugs potentiels, et problèmes de design.

## Capabilities

- Analyse statique du code Go
- Détection de violations de conventions
- Suggestions de design patterns appropriés
- Auto-correction des problèmes simples
- Vérification de la coverage des tests

## Usage

Cette skill s'active automatiquement lorsque vous travaillez avec du code Go.

## Workflow

1. **Détection** : Identifie les fichiers `.go` modifiés
2. **Analyse** : Exécute ktn-linter et golangci-lint
3. **Rapport** : Présente les violations avec sévérité (ERROR/WARNING/INFO)
4. **Suggestions** : Propose corrections et design patterns appropriés
5. **Auto-fix** : Corrige automatiquement quand possible

## Configuration

### Prerequisites

```bash
# Go 1.23+ requis
go version

# Installation golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Si le projet a ktn-linter
make build
```

### Linter Configuration

Si `.golangci.yml` n'existe pas, créer :

```yaml
run:
  timeout: 5m
  go: "1.25"

linters:
  enable:
    - errcheck      # Vérifie erreurs ignorées
    - gosimple      # Simplifications possibles
    - govet         # Examine le code
    - ineffassign   # Détecte assignations inutiles
    - staticcheck   # Analyse statique avancée
    - unused        # Code inutilisé
    - gofmt         # Formatage
    - goimports     # Imports triés
    - misspell      # Fautes d'orthographe
    - gocritic      # Suggestions critiques

linters-settings:
  errcheck:
    check-blank: true
  govet:
    check-shadowing: true
  gofmt:
    simplify: true
```

## Execution

### Auto-Lint après modification

```bash
# 1. Compiler le code
go build ./...

# 2. Linter (ktn + golangci)
make lint

# 3. Tests avec coverage
make test

# 4. Si échecs → CORRIGER avant de continuer
```

### Commandes Disponibles

```bash
# Analyse complète
make lint

# Tests avec coverage
make test

# Build
make build

# Tout en une fois
make all
```

## Severity Levels

### ERROR (Red ✖)
**Blocants** - bugs potentiels, violations graves
- KTN-VAR-003 : Variable réassignée (paramètre)
- KTN-VAR-004 : Variable globale mutable
- KTN-FUNC-006 : Erreur pas en dernière position
- KTN-FUNC-008 : Context pas en premier paramètre
- KTN-FUNC-012 : else après return (dead code)

**Action** : CORRIGER IMMÉDIATEMENT

### WARNING (Orange ⚠)
**Maintenabilité** - problèmes de qualité/conventions
- KTN-CONST-001 : Nommage ALL_CAPS manquant
- KTN-FUNC-001 : Fonction trop longue (>35 lignes)
- KTN-FUNC-002 : Trop de paramètres (>5)
- KTN-STRUCT-004 : Documentation manquante

**Action** : Corriger avant commit

### INFO (Blue ℹ)
**Style** - recommandations, optimisations
- KTN-CONST-002 : Groupement iota/valeur
- KTN-VAR-006 : Utiliser := au lieu de var
- KTN-FUNC-003 : Extraire constantes magiques

**Action** : Améliorer progressivement

## Design Pattern Detection

La skill détecte automatiquement les cas où un design pattern devrait être appliqué :

### Détection : Functional Options
```go
// ❌ Détecté : trop de paramètres
func NewServer(host string, port int, timeout time.Duration, maxConn int) *Server

// ✅ Suggestion : Functional Options Pattern
func NewServer(opts ...Option) *Server
```

### Détection : Builder
```go
// ❌ Détecté : construction complexe
s := &Server{}
s.SetHost("localhost")
s.SetPort(8080)
s.SetTimeout(30)

// ✅ Suggestion : Builder Pattern
s := NewServerBuilder().
    Host("localhost").
    Port(8080).
    Timeout(30).
    Build()
```

### Détection : Strategy
```go
// ❌ Détecté : switch sur type d'algorithme
func Process(data []byte, algo string) {
    switch algo {
    case "gzip": // ...
    case "zlib": // ...
    }
}

// ✅ Suggestion : Strategy Pattern
type Compressor interface {
    Compress([]byte) ([]byte, error)
}
```

## Auto-Fix Capabilities

La skill peut auto-corriger :

1. **Imports** : Trier et nettoyer
2. **Formatting** : gofmt, goimports
3. **Naming** : snake_case → camelCase (variables)
4. **Documentation** : Templates de commentaires
5. **Error handling** : Ajout de checks manquants

## Integration avec Autres Tools

### golangci-lint
Exécuté automatiquement via `make lint`

### ktn-linter
Règles custom strictes du projet

### go test
Coverage reporting automatique

## Output Format

```
📁 File: /path/to/file.go (3 issues)
────────────────────────────────────────────────

[1] /path/to/file.go:42:1
  ✖ Code: KTN-FUNC-006
  ▶ error parameter should be last in function signature

[2] /path/to/file.go:55:2
  ⚠ Code: KTN-FUNC-001
  ▶ function too long (42 lines), consider extracting to smaller functions

[3] /path/to/file.go:78:5
  ℹ Code: KTN-VAR-006
  ▶ use short variable declaration ':=' instead of 'var'

════════════════════════════════════════════════
📊 Total: 3 issue(s) to fix
════════════════════════════════════════════════
```

## Best Practices Enforcement

### Before Commit
```bash
# 1. Lint
make lint
# Expected: ✅ No issues found! Code is compliant.

# 2. Test
make test
# Expected: ok, PASS, coverage ≥90%

# 3. Build
make build
# Expected: ✅ Binaire créé
```

### CI/CD Integration
```yaml
# .github/workflows/ci.yml
- name: Lint
  run: make lint

- name: Test
  run: make test

- name: Build
  run: make build
```

## Troubleshooting

### "golangci-lint not found"
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### "ktn-linter not found"
```bash
# Si le projet contient ktn-linter
make build

# Sinon, cloner le repo
git clone https://github.com/kodflow/ktn-linter
cd ktn-linter && make build
```

### "Tests failing"
```bash
# Verbose output
go test -v ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Learning Mode

En mode apprentissage, la skill explique :
- **Pourquoi** une règle existe
- **Quel** design pattern appliquer
- **Comment** refactorer le code
- **Exemples** de code correct

## Continuous Improvement

La skill vérifie régulièrement :
- Nouvelles versions de Go (go.dev/doc/devel/release)
- Updates golangci-lint
- Évolution des best practices
- Nouveaux design patterns

## Performance

- Analyse < 1s pour fichiers individuels
- Analyse complète projet < 10s
- Cache des résultats pour fichiers non modifiés
