# KTN-VAR-008

**Sévérité**: WARNING

## Description

Éviter les allocations dans les boucles chaudes.

## Exemple conforme

```go
buffer := make([]byte, 1024)
for i := 0; i < n; i++ {
    // Réutiliser buffer
    process(buffer)
}
```
