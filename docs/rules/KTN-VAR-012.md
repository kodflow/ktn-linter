# KTN-VAR-012

**Sévérité**: WARNING

## Description

Éviter les conversions `string()` répétées.

## Exemple conforme

```go
s := string(bytes)
use(s)
use(s)  // Réutiliser la string convertie
```
