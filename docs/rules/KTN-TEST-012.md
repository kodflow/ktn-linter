# KTN-TEST-012

**Sévérité**: WARNING

## Description

Les tests doivent contenir des assertions et vraiment tester quelque chose.

## Exemple conforme

```go
func TestAdd(t *testing.T) {
    result := Add(1, 2)
    if result != 3 {
        t.Errorf("got %d, want 3", result)
    }
}
```
