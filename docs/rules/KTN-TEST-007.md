# KTN-TEST-007

**Sévérité**: WARNING

## Description

Interdiction d'utiliser `t.Skip()` dans les tests.

## Exemple conforme

```go
func TestFeature(t *testing.T) {
    // Implémenter le test au lieu de le skip
    result := Feature()
    if result != expected {
        t.Error("failed")
    }
}
```
