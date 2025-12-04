# KTN-VAR-010

**Sévérité**: WARNING

## Description

Utiliser `sync.Pool` pour les buffers répétés.

## Exemple conforme

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}
```
