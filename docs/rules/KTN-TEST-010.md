# KTN-TEST-010

**Sévérité**: WARNING

## Description

Les tests de fonctions privées doivent être dans `_internal_test.go` uniquement (white-box testing).

## Exemple conforme

```go
// user_internal_test.go
package user

func Test_parseData(t *testing.T) {
    parseData() // Accès direct aux fonctions privées
}
```
