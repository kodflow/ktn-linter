# KTN-COMMENT-006

**Sévérité**: WARNING

## Description

Toutes les fonctions doivent avoir une documentation au format strict avec description, `Params:` et `Returns:`.

## Exemple conforme

```go
// ProcessData traite les données entrantes.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - string: le résultat formaté
//   - error: erreur éventuelle
func ProcessData(data string) (string, error) {
    return data, nil
}
```
