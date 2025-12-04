# KTN-FUNC-012

**Sévérité**: INFO

## Description

Les fonctions avec plus de 3 valeurs de retour doivent utiliser des named returns.

## Exemple conforme

```go
func Process() (result string, count int, valid bool, err error) {
    return "ok", 1, true, nil
}
```
