# Go Expert KTN Plugin

Plugin Claude Code pour Go 1.25+. Auto-lint, patterns, 0 dette technique.

## Installation

```bash
# Clone repo
git clone https://github.com/kodflow/ktn-linter

# Plugin dans .claude-plugin/ détecté automatiquement par Claude Code
```

## Features

- Auto-lint après chaque .go modifié
- Auto-test après chaque *_test.go modifié
- Pre-commit hook (lint + test + build)
- Règle stricte: 0 issues = 0 issues (même INFO)

## Hooks

1. `auto-lint-go` → `make lint` après édition .go
2. `auto-test-go` → `make test` après édition test
3. `pre-commit-check` → lint + test + build
4. `go-version-check` → Go ≥1.25 requis
5. `setup-linters` → install golangci-lint si manquant

## Patterns

- Trop de params (>3) → Functional Options
- Construction complexe → Builder
- Algo interchangeable → Strategy
- I/O parallèle → Worker Pool
- Pipeline → Channels

## Règle Absolue

**STOP si > 0 issues**. Corriger IMMÉDIATEMENT. Pas d'exception.

## Workflow

```bash
# Chaque modification
make lint && make test
# Si échec → STOP et corriger
```

---

Repo: https://github.com/kodflow/ktn-linter
