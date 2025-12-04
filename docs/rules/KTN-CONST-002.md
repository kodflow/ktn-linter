# KTN-CONST-002

**Sévérité**: INFO

## Description

Les constantes doivent être groupées ensemble dans un seul bloc et placées au-dessus des déclarations `var`.

## Exemple conforme

```go
const (
    MAX_SIZE int = 100
    MIN_SIZE int = 10
)

var (
    currentSize int = MAX_SIZE
)
```
