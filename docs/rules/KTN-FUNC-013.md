# KTN-FUNC-013

**Sévérité**: WARNING

## Description

Préférer retourner une slice/map vide plutôt que nil pour éviter les nil pointer dereferences.

## Exemple non-conforme

```go
func GetUsers() []string {
    return nil
}

func GetConfig() map[string]int {
    return nil
}
```

## Exemple conforme

```go
func GetUsers() []string {
    return []string{}
}

func GetConfig() map[string]int {
    return map[string]int{}
}
```
