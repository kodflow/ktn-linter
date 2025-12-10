# KTN-VAR-009

**Sévérité**: WARNING

## Description

Utiliser des pointeurs pour les structs de plus de 64 bytes.

## Exemple conforme

```go
func Process(data *LargeStruct) {
    // Évite la copie de la struct
}
```
