# Go Expert KTN Plugin

## 🎯 Le Plugin IA Go Ultime

Transforme Claude Code en **expert Go** avec :
- ✅ Auto-linting après chaque modification
- ✅ 13+ design patterns intégrés
- ✅ Connaissance Go 1.25+ à jour
- ✅ Zéro dette technique garantie
- ✅ Configuration automatique

## 🚀 Installation

### Via Marketplace (quand disponible)

```bash
# Dans Claude Code
/plugins install go-expert-ktn
```

### Installation Manuelle

1. Cloner le repo :
```bash
git clone https://github.com/kodflow/ktn-linter
```

2. Copier le plugin :
```bash
mkdir -p ~/.claude/plugins
cp -r ktn-linter/.claude-plugin ~/.claude/plugins/go-expert-ktn
```

3. Redémarrer Claude Code

## 📋 Features

### 1. Auto-Linting Réflexe

**Chaque fois que vous modifiez du code Go** :
- ✅ `make lint` s'exécute automatiquement
- ✅ Violations affichées avec couleurs (ERROR/WARNING/INFO)
- ✅ Suggestions de corrections
- ✅ Blocage si erreurs critiques

```bash
# Exemple de sortie
📁 File: /workspace/main.go (2 issues)
────────────────────────────────────────────────

[1] /workspace/main.go:15:1
  ✖ Code: KTN-FUNC-006
  ▶ error parameter should be last in function signature

[2] /workspace/main.go:23:2
  ⚠ Code: KTN-FUNC-001
  ▶ function too long (42 lines), extract to smaller functions
```

### 2. Design Patterns Intelligence

Le plugin **détecte automatiquement** quand utiliser :

#### 🏗️ Functional Options Pattern
```go
// ❌ Avant : trop de paramètres
func NewServer(host string, port int, timeout time.Duration, maxConn int)

// ✅ Après : suggestions du plugin
func NewServer(opts ...Option) *Server {
    // ...
}
```

#### 🔧 Strategy Pattern
```go
// ❌ Avant : switch sur type
func Compress(data []byte, algo string)

// ✅ Après : suggestions du plugin
type Compressor interface {
    Compress([]byte) ([]byte, error)
}
```

**13 patterns disponibles** : Factory, Builder, Adapter, Decorator, Observer, Worker Pool, Pipeline, Fan-Out/Fan-In, et plus...

### 3. Go Version Tracking

Le plugin **vérifie automatiquement** :
- ✅ Version Go installée (≥1.23 requis)
- ✅ Dernières features Go 1.24/1.25
- ✅ Release notes (go.dev/doc/devel/release)
- ✅ Breaking changes

```go
// Utilise automatiquement les nouvelles features
for i := range 10 {  // Go 1.22+
    fmt.Println(i)
}
```

### 4. Zero Debt Enforcement

**Workflow strict** :
1. Écrire code
2. Auto-lint → corrige TOUTES les violations
3. Auto-test → coverage ≥90%
4. Si échec → **BLOQUE** jusqu'à correction

**Impossible d'accumuler de la dette !**

### 5. Configuration Intelligente

Le plugin **configure automatiquement** :
- golangci-lint (si manquant)
- .golangci.yml (config stricte)
- Makefile (targets lint/test/build)
- Pre-commit hooks

## 🎓 Usage

### Cas d'Usage 1 : Nouveau Projet

```bash
# Claude détecte : pas de linter configuré
# Claude propose :
"Je vois que ce projet n'a pas de linter configuré.
Je vais créer .golangci.yml et Makefile. OK ?"

# ✅ Auto-configuration
# ✅ Installation golangci-lint
# ✅ Premier lint du code existant
```

### Cas d'Usage 2 : Refactoring

```bash
# Vous : "Refactor cette fonction trop longue"

# Claude :
# 1. Analyse la fonction (98 lignes)
# 2. Détecte : KTN-FUNC-001 (>35 lignes)
# 3. Identifie : Builder pattern approprié
# 4. Propose refactor avec 4 fonctions < 25 lignes
# 5. Execute make lint → ✅ 0 issues
# 6. Execute make test → ✅ PASS (coverage 95%)
```

### Cas d'Usage 3 : Code Review

```bash
# Vous : "Review ce code"

# Claude analyse avec :
# - ktn-linter (28 règles strictes)
# - golangci-lint (20+ linters)
# - Design pattern detection
# - Performance suggestions
# - Security checks

# Retour structuré :
# ✖ 2 ERRORS (blocants)
# ⚠ 5 WARNINGS (maintenabilité)
# ℹ 3 INFO (suggestions)
```

## ⚙️ Configuration

### Personnaliser le Plugin

Éditer `.claude-plugin/agents/go-expert.md` pour :
- Ajouter vos propres règles
- Modifier les seuils (lignes max, params max, etc.)
- Activer/désactiver certains patterns
- Changer le niveau de strictness

### Intégration CI/CD

```yaml
# .github/workflows/ci.yml
name: Go CI

on: [push, pull_request]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

      - name: Lint
        run: make lint

      - name: Test
        run: make test

      - name: Build
        run: make build
```

## 🔧 Hooks Disponibles

| Hook | Trigger | Action |
|------|---------|--------|
| `auto-lint-go` | Après édition .go | `make lint` |
| `auto-test-go` | Après édition *_test.go | `make test` |
| `pre-commit-check` | Avant commit | Lint + Test + Build |
| `go-version-check` | Au démarrage | Vérifie Go ≥1.23 |
| `setup-linters` | Au démarrage | Install golangci-lint |

## 📊 Métriques Qualité

Le plugin garantit :
- ✅ **0 warnings** lint
- ✅ **100% tests** PASS
- ✅ **≥90% coverage**
- ✅ **Documentation** complète
- ✅ **Design patterns** appropriés

## 🎯 Réponse au Post Reddit

Ce plugin résout **EXACTEMENT** les problèmes mentionnés :

### ❌ Problème : "Claude oublie les conventions Go"
✅ **Solution** : Agent Go expert avec conventions Go 1.25+ intégrées

### ❌ Problème : "Pas de contexte projet persistant"
✅ **Solution** : Hooks auto-lint après chaque modification

### ❌ Problème : "Répéter les règles à chaque prompt"
✅ **Solution** : Règles KTN + golangci-lint toujours actives

### ❌ Problème : "Linter externe puis feedback à Claude"
✅ **Solution** : Linter intégré, exécution automatique, feedback instantané

### ❌ Problème : "Claude hallucine ou downgrade la qualité"
✅ **Solution** : Règles machine-readable (KTN), blocage sur violations

## 🚦 Severity System

### ✖ ERROR (Rouge)
**Bugs potentiels, violations graves**
- Action : CORRIGER IMMÉDIATEMENT
- Exemples : variables globales mutables, error mal positionné

### ⚠ WARNING (Orange)
**Maintenabilité, conventions**
- Action : Corriger avant commit
- Exemples : fonction trop longue, doc manquante

### ℹ INFO (Bleu)
**Style, optimisations**
- Action : Améliorer progressivement
- Exemples : utiliser :=, grouper constantes

## 📚 Ressources

- **Repo** : https://github.com/kodflow/ktn-linter
- **Issues** : https://github.com/kodflow/ktn-linter/issues
- **Docs Go** : https://go.dev/doc/
- **Design Patterns** : Intégrés dans l'agent

## 🤝 Contribution

```bash
# Ajouter un design pattern
# Éditer : .claude-plugin/agents/go-expert.md

# Ajouter une règle KTN
# Voir : pkg/analyzer/ktn/

# Tester le plugin
make test && make lint
```

## 📝 License

MIT - Voir LICENSE

## 🎉 Résultat

**Avant le plugin** :
```
❌ Claude oublie les conventions
❌ Répète les mêmes erreurs
❌ Pas de contexte projet
❌ Qualité variable
```

**Avec le plugin** :
```
✅ Conventions Go 1.25+ automatiques
✅ Auto-correction réflexe
✅ Contexte projet persistant
✅ Qualité production garantie
✅ 0 dette technique
```

---

**Made with ❤️ by kodflow**
