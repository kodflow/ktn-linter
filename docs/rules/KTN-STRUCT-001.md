# KTN-STRUCT-001

**Sévérité**: WARNING

## Description

Chaque struct doit avoir une interface reprenant 100% de ses méthodes publiques.

## Exemple conforme

```go
type UserInterface interface {
    Name() string
}

type User struct {
    name string
}

func (u *User) Name() string {
    return u.name
}
```
