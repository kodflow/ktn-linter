# KTN-TEST-009

**Sévérité**: WARNING

## Description

Les tests de fonctions publiques doivent être dans `_external_test.go` uniquement (black-box testing).

## Exemple conforme

```go
// user_external_test.go
package user_test

func TestGetUser(t *testing.T) {
    user.GetUser() // Test via l'API publique
}
```
