# KTN-FUNC - Règles pour les fonctions

Ce package contient toutes les règles de validation pour les fonctions en Go.

## Vue d'ensemble

Les règles KTN-FUNC garantissent que les fonctions sont bien nommées, documentées et maintiennent une complexité raisonnable.

## Règles disponibles

| Règle | Description | Exemples |
|-------|-------------|----------|
| [KTN-FUNC-001](./001.go) | Nommage MixedCaps (pas de snake_case) | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_001_naming.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_001_naming.go) |
| [KTN-FUNC-002](./002.go) | Documentation godoc obligatoire | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_002_godoc.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_002_godoc.go) |
| [KTN-FUNC-003](./003.go) | Format strict section Params: | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_003_params.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_003_params.go) |
| [KTN-FUNC-004](./004.go) | Format strict section Returns: | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_004_returns.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_004_returns.go) |
| [KTN-FUNC-005](./005.go) | Maximum 5 paramètres par fonction | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_005_max_params.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_005_max_params.go) |
| [KTN-FUNC-006](./006.go) | Longueur max: 35 lignes (100 pour tests) | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_006_length.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_006_length.go) |
| [KTN-FUNC-007](./007.go) | Complexité cyclomatique max: 10 (50 tests) | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_007_complexity.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_007_complexity.go) |
| [KTN-FUNC-008](./008.go) | Commentaires obligatoires sur return | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_008_return_comments.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_008_return_comments.go) |
| [KTN-FUNC-010](./010.go) | Profondeur d'imbrication max: 3 niveaux | [✓ Bon](../../tests/good_usage/ktn/rules_func/ktn_func_010_nesting.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_func/ktn_func_010_nesting.go) |

## Utilisation

```go
import "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"

// Obtenir toutes les règles
rules := ktn_func.GetRules()

// Utiliser une règle spécifique
analyzer := ktn_func.Rule001 // KTN-FUNC-001
```

## Exemples

### ✓ Fonction conforme

```go
// CalculateTotal calcule le total des items avec la taxe.
//
// Params:
//   - items: liste des montants à additionner
//   - taxRate: taux de taxe à appliquer (ex: 0.2 pour 20%)
//
// Returns:
//   - float64: le total calculé avec taxe
//   - error: erreur si taxRate invalide
func CalculateTotal(items []float64, taxRate float64) (float64, error) {
    if taxRate < 0 || taxRate > 1 {
        // Taux de taxe invalide
        return 0, errors.New("invalid tax rate")
    }

    total := 0.0
    for _, item := range items {
        total += item
    }

    // Retourne le total avec taxe appliquée
    return total * (1 + taxRate), nil
}
```

### ✗ Fonction non-conforme

```go
// Violations multiples :
// - FUNC-001: snake_case au lieu de MixedCaps
// - FUNC-002: pas de description
// - FUNC-003: section Params: manquante
// - FUNC-004: section Returns: manquante
// - FUNC-008: pas de commentaire sur return
func calculate_total(items []float64, tax_rate float64) (float64, error) {
    total := 0.0
    for _, item := range items {
        total += item
    }
    return total * (1 + tax_rate), nil
}
```

## Tests

Les tests sont organisés dans `/tests/`:
- `/tests/good_usage/ktn/rules_func/` - Exemples conformes
- `/tests/bad_usage/ktn/rules_func/` - Exemples non-conformes

Exécuter les tests :
```bash
go test ./analyzer/ktn/func/...
```
