# KTN-FUNC-009

**Sévérité**: INFO

## Description

Les nombres littéraux doivent être des constantes nommées (pas de magic numbers). Exception: 0, 1, -1.

## Exemple conforme

```go
const MAX_RETRIES int = 3

func Retry() int {
    return MAX_RETRIES
}
```
