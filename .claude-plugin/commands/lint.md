# Lint - KTN-Linter Analysis

Execute le linter KTN avec 94 règles strictes sur le projet Go.

## Action

1. Build le linter si nécessaire (`make build`)
2. Exécute `./builds/ktn-linter lint ./...`
3. Affiche tous les problèmes détectés (ERROR/WARNING/INFO)
4. **STOP** si > 0 issues détectées

## Règle Absolue

**0 issues = 0 issues**. Même les INFO doivent être corrigées avant de continuer.

## Severity Levels

- **ERROR** (✖ rouge) : Blocant, bug potentiel, violation majeure
- **WARNING** (⚠ orange) : Maintenabilité, qualité de code
- **INFO** (ℹ bleu) : Style, conventions, best practices

## Commandes Disponibles

```bash
make build  # Compile ktn-linter
make lint   # Exécute ktn-linter lint ./...
make test   # Tests avec coverage
```

## Workflow

Après correction de tous les issues :
1. Re-linter pour vérifier
2. Lancer les tests (`make test`)
3. Vérifier la coverage
