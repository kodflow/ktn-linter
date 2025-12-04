# KTN-STRUCT-002

**Sévérité**: WARNING

## Description

Les structs exportées avec méthodes doivent avoir un constructeur `NewX()`.

## Exemple conforme

```go
type User struct {
    name string
}

func NewUser(name string) *User {
    return &User{name: name}
}
```
