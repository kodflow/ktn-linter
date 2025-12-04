# KTN-COMMENT-007

**Sévérité**: WARNING

## Description

Tous les blocs de contrôle (`if`/`else`/`switch`/`for`), returns et logique significative doivent être commentés.

## Exemple conforme

```go
func Process(x int) string {
    // Check if x is positive
    if x > 0 {
        // Return positive result
        return "positive"
    }
    // Return negative result
    return "negative"
}
```
