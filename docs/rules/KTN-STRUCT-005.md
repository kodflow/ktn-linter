# KTN-STRUCT-005

**Sévérité**: INFO

## Description

Les champs exportés doivent être placés avant les champs privés dans une struct.

## Exemple conforme

```go
type User struct {
    Name  string // exporté
    Age   int    // exporté
    id    int    // privé
    email string // privé
}
```
