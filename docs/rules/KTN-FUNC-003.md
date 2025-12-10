# KTN-FUNC-003

**Sévérité**: ERROR

## Description

Éviter `else` après `return`/`continue`/`break`. Préférer le pattern early return.

## Exemple conforme

```go
func Process(x int) string {
    if x < 0 {
        return "negative"
    }
    return "positive"
}
```
