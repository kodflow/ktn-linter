# KTN-FUNC-011

**Sévérité**: INFO

## Description

La complexité cyclomatique ne doit pas dépasser 10.

## Exemple conforme

```go
func Process(x int) string {
    if x > 0 {
        return "positive"
    }
    return "non-positive"
}
```
