# KTN-VAR-007

**Sévérité**: WARNING

## Description

Utiliser `strings.Builder` pour plus de 2 concaténations.

## Exemple conforme

```go
var b strings.Builder
b.WriteString("Hello")
b.WriteString(" ")
b.WriteString("World")
result := b.String()
```
