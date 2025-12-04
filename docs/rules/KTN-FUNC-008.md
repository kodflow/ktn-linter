# KTN-FUNC-008

**Sévérité**: WARNING

## Description

Les paramètres non utilisés doivent être préfixés par `_` ou assignés à `_`.

## Exemple conforme

```go
func Process(_ctx context.Context, data string) string {
    return data
}
```
