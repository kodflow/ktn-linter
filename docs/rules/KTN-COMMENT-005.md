# KTN-COMMENT-005

**Sévérité**: WARNING

## Description

Toute struct exportée doit avoir une documentation complète (au moins 2 lignes décrivant son rôle).

## Exemple conforme

```go
// User représente un utilisateur du système.
// Stocke les informations de base d'un utilisateur.
type User struct {
    Name string
    Age  int
}
```
