# KTN-FUNC-010

**Sévérité**: INFO

## Description

Les naked returns sont interdits sauf pour les fonctions très courtes (<5 lignes).

## Exemple conforme

```go
func Process(x int) (result int, err error) {
    result = x * 2
    return result, nil
}
```
