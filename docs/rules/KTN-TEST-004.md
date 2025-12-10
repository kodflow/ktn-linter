# KTN-TEST-004

**Sévérité**: WARNING

## Description

Toutes les fonctions publiques doivent avoir des tests.

## Exemple conforme

```go
// user.go
func GetUser() {}

// user_external_test.go
func TestGetUser(t *testing.T) {}
```
