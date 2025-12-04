# KTN-VAR-017

**Sévérité**: INFO

## Description

Éviter les copies de mutex (`sync.Mutex`, `sync.RWMutex`).

## Exemple conforme

```go
type Safe struct {
    mu sync.Mutex  // Jamais copier cette struct
}

func (s *Safe) Do() {  // Utiliser pointeur
    s.mu.Lock()
    defer s.mu.Unlock()
}
```
