# KTN-TEST-005

**Sévérité**: WARNING

## Description

Les tests avec plusieurs cas doivent utiliser le pattern table-driven.

## Exemple conforme

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 2, 3},
        {0, 0, 0},
    }
    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    }
}
```
