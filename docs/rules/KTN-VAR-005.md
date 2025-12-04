# KTN-VAR-005

**Sévérité**: WARNING

## Description

Éviter `make([]T, length)` avec `append` (cause réallocation).

## Exemple conforme

```go
items := make([]string, 0, length)
for _, v := range source {
    items = append(items, v)
}
```
