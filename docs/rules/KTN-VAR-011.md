# KTN-VAR-011

**Sévérité**: WARNING

## Description

Éviter le shadowing de variables.

## Exemple conforme

```go
var err error
result, err = process()  // Pas de := pour éviter le shadowing
```
