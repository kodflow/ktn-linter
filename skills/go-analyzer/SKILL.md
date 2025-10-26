# Go Code Analyzer Skill

Analyse auto du code Go avec ktn-linter + golangci-lint.

## Workflow

1. Détection fichiers .go modifiés
2. Exécution `make lint`
3. Rapport violations (ERROR/WARNING/INFO)
4. STOP si > 0 issues

## Severity

- **ERROR** (✖ rouge) : Blocant, bug potentiel
- **WARNING** (⚠ orange) : Maintenabilité
- **INFO** (ℹ bleu) : Style

## Règle

**0 issues = 0 issues**. Même INFO doit être corrigé.

## Setup Auto

```bash
# Installe golangci-lint si manquant
# Crée .golangci.yml si absent
# Crée Makefile si absent
```

## Patterns Détectés

Suggère pattern approprié selon contexte :
- Trop params → Functional Options
- Construction → Builder
- Algorithme → Strategy
- I/O parallèle → Worker Pool
- Pipeline → Channels

## Commandes

```bash
make lint   # Analyse complète
make test   # Tests + coverage
make build  # Compilation
```
