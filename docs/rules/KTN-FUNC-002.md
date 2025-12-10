# KTN-FUNC-002

**Sévérité**: ERROR

## Description

`context.Context` doit toujours être le premier paramètre.

## Exemple conforme

```go
func Process(ctx context.Context, data string) error {
    return nil
}
```
