# KTN-VAR-006

**Sévérité**: WARNING

## Description

Préallocation `bytes.Buffer`/`strings.Builder` avec `Grow`.

## Exemple conforme

```go
var buf bytes.Buffer
buf.Grow(1024)
```
