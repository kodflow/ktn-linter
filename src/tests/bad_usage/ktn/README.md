# KTN Rules - Bad Usage Examples (Violations)

Ce dossier contient des exemples de **code qui viole les règles KTN**.

## Objectif

Ces fichiers démontrent des **anti-patterns** et des violations des règles KTN. Ils servent à :
1. Tester que le linter ktn-linter détecte correctement les violations
2. Documenter ce qu'il NE faut PAS faire
3. Fournir des exemples concrets de code problématique

## Structure

Chaque fichier de violation correspond à une ou plusieurs règles KTN :

### Allocation (rules_alloc/)
❌ Violations détectées :
- `new()` avec maps, slices, channels au lieu de `make()`
- `make([]T, 0)` suivi d'append sans capacité
- `new(struct)` au lieu de `&struct{}`

### Constantes (rules_const/)
❌ Violations détectées :
- Constantes déclarées individuellement au lieu d'être groupées
- Absence de commentaires
- Types implicites pour constantes publiques

### Variables (rules_var/)
❌ Violations détectées :
- Variables non groupées
- Nommage incorrect (snake_case, ALL_CAPS)
- Absence de commentaires
- Variables qui devraient être const

### Fonctions (rules_func/)
❌ Violations détectées :
- Fonctions trop longues (>50 lignes)
- Trop de paramètres (>5)
- Complexité cyclomatique élevée
- Absence de documentation
- Nesting trop profond

### Erreurs (rules_error/)
❌ Violations détectées :
- `return err` sans wrapping avec `fmt.Errorf("...: %w", err)`
- Perte du contexte d'erreur

### Goroutines (rules_goroutine/)
❌ Violations détectées :
- Goroutines dans boucles sans sync.WaitGroup
- Goroutines sans mécanisme de synchronisation

## Usage par le linter

Le linter ktn-linter analyse ces fichiers et doit détecter **toutes les violations**.

Exemple :
```bash
ktn-linter tests/bad_usage/ktn/rules_alloc/ktn_alloc_001_new_with_ref_types.go
# Doit détecter : KTN-ALLOC-001 sur chaque new(map), new([]T), new(chan T)
```

## Relation avec good_usage/

Pour chaque fichier dans `bad_usage/ktn/`, il existe un fichier correspondant dans `good_usage/ktn/` qui montre la **version correcte** du même code.

## Note importante

⚠️ **Ce code compile** mais **viole les conventions KTN**.
La différence avec `bad_usage/gospec/` (à venir) est que celui-ci viole uniquement les best practices KTN, pas la spec Go elle-même.
