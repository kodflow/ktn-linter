# KTN-ERROR - Règles pour la gestion des erreurs

Ce package contient toutes les règles de validation pour la gestion des erreurs en Go.

## Vue d'ensemble

La gestion d'erreurs est cruciale en Go. Les règles KTN-ERROR garantissent que les erreurs sont correctement wrappées et contextualisées pour faciliter le debugging en production.

## Règles disponibles

| Règle | Description | Exemples |
|-------|-------------|----------|
| [KTN-ERROR-001](./001.go) | Wrapping obligatoire des erreurs avec contexte | [✓ Bon](../../tests/good_usage/ktn/rules_error/ktn_error_001_wrapping.go) [✗ Mauvais](../../tests/bad_usage/ktn/rules_error/ktn_error_001_unwrapped_error.go) |

## Utilisation

```go
import "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/error"

// Obtenir toutes les règles
rules := ktn_error.GetRules()

// Utiliser la règle
analyzer := ktn_error.Rule001
```

## Exemples

### ✓ Gestion d'erreurs conforme

```go
func ProcessUser(id string) error {
    user, err := db.GetUser(id)
    if err != nil {
        // ✅ CORRECT: erreur wrappée avec contexte
        return fmt.Errorf("failed to get user %s: %w", id, err)
    }

    if err := user.Validate(); err != nil {
        // ✅ CORRECT: contexte ajouté
        return fmt.Errorf("user %s validation failed: %w", id, err)
    }

    // Succès
    return nil
}
```

### ✗ Gestion d'erreurs non-conforme

```go
func ProcessUser(id string) error {
    user, err := db.GetUser(id)
    if err != nil {
        // ❌ MAUVAIS: erreur retournée sans contexte
        return err  // Viole KTN-ERROR-001
    }

    if err := user.Validate(); err != nil {
        // ❌ MAUVAIS: perd le contexte
        return err  // Viole KTN-ERROR-001
    }

    return nil
}
```

## Pourquoi le wrapping d'erreurs ?

1. **Traçabilité** : Permet de suivre le chemin de l'erreur à travers les couches de l'application
2. **Debugging** : Facilite la compréhension des erreurs en production
3. **Context** : Ajoute des informations spécifiques au point d'échec
4. **Standards Go** : Suit les recommandations officielles de Go 1.13+

## Tests

Les tests sont organisés dans `/tests/`:
- `/tests/good_usage/ktn/rules_error/` - Exemples conformes
- `/tests/bad_usage/ktn/rules_error/` - Exemples non-conformes

Exécuter les tests :
```bash
go test ./analyzer/ktn/error/...
```
