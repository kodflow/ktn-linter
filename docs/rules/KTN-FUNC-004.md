# KTN-FUNC-004

**Sévérité**: ERROR

## Description

Les fonctions privées non utilisées dans le code de production sont du code mort et doivent être supprimées.

## Exemple conforme

```go
func process() string {
    return helper() // helper est utilisée
}

func helper() string {
    return "data"
}
```
