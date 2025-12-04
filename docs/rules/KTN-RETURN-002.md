# KTN-RETURN-002

**Sévérité**: WARNING

## Description

Préférer retourner une slice/map vide plutôt que `nil`.

## Exemple conforme

```go
func GetItems() []string {
    if empty {
        return []string{}  // Pas nil
    }
    return items
}
```
