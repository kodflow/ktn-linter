# KTN-FUNC-007

**Sévérité**: WARNING

## Description

Les getters (`Get*`/`Is*`/`Has*`) ne doivent pas avoir de side effects.

## Exemple conforme

```go
func (u *User) GetName() string {
    return u.name
}
```
