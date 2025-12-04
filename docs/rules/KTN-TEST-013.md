# KTN-TEST-013

**Sévérité**: INFO

## Description

Les tests doivent couvrir les cas d'erreur (coverage des chemins d'erreur).

## Exemple conforme

```go
func TestProcess_Error(t *testing.T) {
    _, err := Process(invalidInput)
    if err == nil {
        t.Error("expected error")
    }
}
```
