# KTN-FUNC-001

**Sévérité**: ERROR

## Description

L'erreur doit toujours être en dernière position dans les valeurs de retour.

## Exemple conforme

```go
func Process(data string) (string, error) {
    return data, nil
}
```
