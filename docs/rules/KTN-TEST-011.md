# KTN-TEST-011

**Sévérité**: WARNING

## Description

Les fichiers `_internal_test.go` doivent utiliser `package xxx`, les fichiers `_external_test.go` doivent utiliser `package xxx_test`.

## Exemple conforme

```go
// user_internal_test.go
package user  // Même package

// user_external_test.go
package user_test  // Package _test
```
