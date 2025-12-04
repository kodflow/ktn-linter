# KTN-TEST-008

**Sévérité**: WARNING

## Description

Chaque fichier `.go` doit avoir les fichiers de test appropriés (`_internal_test.go` si fonctions privées, `_external_test.go` si fonctions publiques).

## Exemple conforme

```
user.go                    // Contient GetUser() et parseData()
user_external_test.go      // Teste GetUser()
user_internal_test.go      // Teste parseData()
```
