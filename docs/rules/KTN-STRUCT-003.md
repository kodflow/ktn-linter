# KTN-STRUCT-003

**Sévérité**: WARNING

## Description

Les getters ne doivent pas avoir le préfixe `Get` (convention Go idiomatique).

## Exemple conforme

```go
type User struct {
    name string
}

func (u *User) Name() string {
    return u.name
}
```
