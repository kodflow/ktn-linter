# KTN-FUNC-005

**Sévérité**: WARNING

## Description

Les fonctions ne doivent pas dépasser 35 lignes de code pur (hors commentaires et lignes vides).

## Exemple conforme

```go
func Process(data string) string {
    result := transform(data)
    return result
}
```
