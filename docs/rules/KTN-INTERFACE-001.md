# KTN-INTERFACE-001

**Sévérité**: WARNING

## Description

Les interfaces déclarées doivent être utilisées.

## Exemple conforme

```go
type Reader interface {
    Read() []byte
}

func Process(r Reader) {
    data := r.Read()
}
```
