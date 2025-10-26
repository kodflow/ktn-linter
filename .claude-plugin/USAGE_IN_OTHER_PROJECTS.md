# Utiliser le Plugin dans Vos Projets Go

Ce guide explique comment intégrer le plugin `go-expert-ktn` dans n'importe quel projet Go.

## 🎯 Scénarios d'Utilisation

### Scénario 1 : Nouveau Projet Go

Vous démarrez un nouveau projet et voulez les best practices dès le début.

```bash
# 1. Créer le projet
mkdir my-api && cd my-api
go mod init github.com/user/my-api

# 2. Installer le plugin
git clone https://github.com/kodflow/ktn-linter /tmp/ktn-linter
cp -r /tmp/ktn-linter/.claude-plugin ~/.claude/plugins/go-expert-ktn

# 3. Ouvrir dans Claude Code
code .

# 4. Le plugin détecte le projet Go et propose la configuration
```

**Résultat attendu** :

```
🤖 Go Expert KTN activé

Nouveau projet Go détecté !

Je vais configurer :
  ✅ .golangci.yml (linting strict)
  ✅ Makefile (targets lint/test/build)
  ✅ CLAUDE.md (instructions projet)

Installation golangci-lint...
  ✅ golangci-lint v1.62.2 installé

Configuration terminée ! Prêt à coder. 🚀
```

### Scénario 2 : Projet Existant (Legacy)

Vous héritez d'un vieux projet Go avec de la dette technique.

```bash
cd legacy-project

# Plugin installé globalement
code .
```

**Plugin s'active** :

```
🔍 Analyse du projet...

Détections:
  ⚠ Go 1.18 (outdated, Go 1.25 disponible)
  ❌ golangci-lint non configuré
  ❌ Pas de Makefile
  ⚠ 847 violations détectées

🤖 Plan d'action proposé:

  Étape 1: Upgrade Go 1.18 → 1.25
    - Modifier go.mod
    - go mod tidy

  Étape 2: Configuration linters
    - Créer .golangci.yml
    - Installer golangci-lint
    - Créer Makefile

  Étape 3: Analyse initiale
    - make lint (détection violations)
    - Tri par sévérité (ERROR/WARNING/INFO)

  Étape 4: Plan de correction
    - Jour 1: 127 ERRORS (blocants)
    - Jour 2-3: 453 WARNINGS (maintenabilité)
    - Jour 4-5: 267 INFO (style)

Accepter le plan ? (y/n)
```

**Après acceptation, le plugin guide étape par étape** :

```
✅ Étape 1/4 complétée
   go.mod: go 1.18 → go 1.25
   Dépendances à jour

✅ Étape 2/4 complétée
   .golangci.yml créé
   golangci-lint installé

📊 Étape 3/4 - Première analyse

make lint
───────────────────────────────────────────────

📁 File: auth/handler.go (23 issues)

[1] auth/handler.go:42:1
  ✖ Code: KTN-FUNC-006
  ▶ error parameter should be last

[2] auth/handler.go:58:5
  ✖ Code: KTN-VAR-004
  ▶ global mutable variable 'sessionStore'

... (21 more)

───────────────────────────────────────────────
📊 Total projet: 127 ERRORS, 453 WARNINGS, 267 INFO

🎯 Commencer par les ERRORS ? (y/n)
```

### Scénario 3 : Contribution Open Source

Vous contribuez à un projet Go open source et voulez respecter leurs conventions.

```bash
# 1. Fork & clone
git clone https://github.com/org/project
cd project

# 2. Le plugin détecte le projet existant
code .
```

**Plugin s'adapte** :

```
🤖 Projet existant détecté

Configuration trouvée:
  ✅ .golangci.yml (détecté)
  ✅ Makefile (détecté)
  ✅ Go 1.24 (ok)

Analyse des conventions du projet...
  📋 50 fichiers analysés
  📊 Patterns détectés:
     - Functional Options (12 utilisations)
     - Builder (5 utilisations)
     - Worker Pool (3 utilisations)

🎓 Mode Learning activé

Je vais respecter les conventions existantes:
  - Naming: camelCase receivers (détecté)
  - Errors: wrapped avec %w (100% du projet)
  - Contexts: toujours en 1er paramètre
  - Tests: table-driven (style dominant)

Hooks adaptés aux conventions du projet.
Prêt à contribuer ! 🚀
```

## 📁 Structure Recommandée

Le plugin fonctionne mieux avec cette structure :

```
your-project/
├── cmd/
│   └── app/
│       └── main.go
├── pkg/
│   └── module/
│       ├── module.go
│       └── module_test.go
├── internal/
│   └── helper/
├── go.mod
├── go.sum
├── Makefile                # Généré par plugin si manquant
├── .golangci.yml          # Généré par plugin si manquant
├── CLAUDE.md              # Optionnel, conventions projet
└── README.md
```

## ⚙️ Configuration Projet-Spécifique

### Créer CLAUDE.md (Optionnel)

Pour des règles spécifiques à votre projet :

```markdown
# Mon Projet API

## Conventions Spécifiques

- Tous les handlers doivent retourner `(Response, error)`
- Utiliser `zap` pour le logging (pas `log`)
- Contexts timeout: 30s pour API, 5s pour DB
- Nommage: `xxxHandler` pour HTTP, `xxxService` pour business logic

## Design Patterns Préférés

- **Configuration** : Functional Options
- **Services** : Dependency Injection (constructeur)
- **Async** : Worker Pool (10 workers max)

## Tests

- Coverage minimale: 85%
- Toujours tester cas d'erreur
- Mocker interfaces externes
```

Le plugin lira ce fichier et adaptera ses suggestions.

### Configuration .golangci.yml Custom

Si vous avez des besoins spécifiques :

```yaml
run:
  timeout: 10m  # Projet large
  go: "1.25"

linters:
  enable:
    - errcheck
    - gosimple
    # ... vos linters
  disable:
    - unused  # Désactiver si besoin

linters-settings:
  errcheck:
    exclude-functions:
      - (io.Closer).Close  # Ignorer Close dans defer
```

Le plugin respectera cette configuration.

## 🔄 Workflow Typique

### Développement Feature

```bash
# 1. Créer branche
git checkout -b feat/user-auth

# 2. Coder (plugin auto-lint à chaque save)
vim pkg/auth/service.go

# Sauvegarde → Auto-lint
📁 File: pkg/auth/service.go (2 issues)

[1] pkg/auth/service.go:25:1
  ⚠ Code: KTN-FUNC-001
  ▶ function ValidateToken is 42 lines (max 35)

  💡 Suggestion: Extract validation logic to separate function

[2] pkg/auth/service.go:48:5
  ℹ Code: KTN-VAR-006
  ▶ use := instead of var

  💡 Auto-fix available

# 3. Corriger violations
# 4. Tests automatiques après modification *_test.go
vim pkg/auth/service_test.go

# Sauvegarde → Auto-test
make test
=== RUN   TestValidateToken
--- PASS: TestValidateToken (0.00s)
PASS
Coverage: 92.3%

# 5. Pre-commit hook avant commit
git add .
git commit -m "feat: add token validation"

# Hook pre-commit s'active
🔍 Pre-commit checks...
  ✅ make lint (0 issues)
  ✅ make test (100% PASS)
  ✅ make build (OK)

[feat/user-auth 3c2a1b4] feat: add token validation
 3 files changed, 85 insertions(+)
```

### Review Pull Request

```bash
# Reviewer utilise le plugin
git checkout pr/123

# Plugin analyse automatiquement
🔍 PR Analysis: feat/user-auth

Files changed: 3
Lines: +85 -12

📊 Quality Check:
  ✅ Lint: 0 issues
  ✅ Tests: 100% PASS
  ✅ Coverage: 92.3% (+2.1%)
  ✅ Design: Functional Options utilisé ✓

🎯 Suggestions:
  ℹ auth/service.go:45
    Consider adding context timeout to DB calls

  ℹ auth/service_test.go:78
    Add test case for expired token

Overall: ✅ APPROVE (minor suggestions)
```

## 🚀 Cas d'Usage Avancés

### Multi-Module Monorepo

```
monorepo/
├── services/
│   ├── api/
│   │   └── go.mod
│   └── worker/
│       └── go.mod
└── libs/
    └── shared/
        └── go.mod
```

Le plugin détecte chaque module et applique les règles indépendamment.

### Microservices

Chaque microservice a le plugin :

```bash
# Service 1
cd user-service
# Plugin actif avec règles user-service

# Service 2
cd order-service
# Plugin actif avec règles order-service
```

Les conventions peuvent différer entre services (adaptabilité).

### CI/CD Integration

**.github/workflows/go.yml** :

```yaml
name: Go CI

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install golangci-lint
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

      - name: Lint
        run: make lint

      - name: Test
        run: make test

      - name: Build
        run: make build

      # Les mêmes checks que le plugin local !
```

## 📊 Métriques & Rapports

Le plugin peut générer des rapports :

```bash
# Coverage report
make test
# → COVERAGE.MD généré

# Complexity report
golangci-lint run --enable=gocyclo

# Security scan
golangci-lint run --enable=gosec
```

## 🆘 Troubleshooting

### Plugin ne détecte pas le projet Go

```bash
# Vérifier go.mod existe
ls go.mod

# Si manquant
go mod init github.com/user/project
```

### Trop de violations initiales

```bash
# Désactiver temporairement certaines règles
# .golangci.yml
linters:
  disable:
    - gosimple  # Réactiver progressivement
    - unused
```

### Conflit avec CI/CD existant

```bash
# Aligner configuration
# Copier .golangci.yml du CI dans le projet
# Le plugin utilisera la même config
```

## 🎓 Best Practices

1. **Activer le plugin dès le début** du projet
2. **Lire CLAUDE.md** pour conventions spécifiques
3. **Corriger les ERRORS immédiatement** (ne pas accumuler)
4. **Refactor progressif** pour legacy (ERRORS → WARNINGS → INFO)
5. **Documenter patterns** utilisés dans le projet
6. **Partager config** (.golangci.yml) avec l'équipe

## 🔗 Ressources

- **Documentation** : [README.md](./.README.md)
- **Exemples détaillés** : [EXAMPLES.md](./EXAMPLES.md)
- **Installation** : [INSTALL.md](./INSTALL.md)
- **Issues** : https://github.com/kodflow/ktn-linter/issues

---

**Le plugin s'adapte à VOTRE projet, pas l'inverse !** 🎯
