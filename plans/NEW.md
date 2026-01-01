# Futures Règles KTN-Linter (Go 1.18 - 1.25)

Generated: 2026-01-01
Go Version: 1.25.5
Sources: go.dev/doc/go1.XX release notes

---

## Résumé des Nouveautés Go par Version

| Version | Feature Majeure | Impact Variables |
|---------|-----------------|------------------|
| 1.18 | Generics, `any` | Nouvelles contraintes |
| 1.21 | `min`, `max`, `clear` built-ins | Nouveaux patterns |
| 1.22 | Loop var per-iteration, range int | Fix bug historique |
| 1.23 | Range-over-func iterators | Nouveaux patterns |
| 1.24 | Swiss Tables maps, generic aliases | Performance |
| 1.25 | `sync.WaitGroup.Go()`, Green Tea GC | Concurrence |

---

## Règles à Implémenter

### 1. VAR-020: Préférer nil slice (Confirmé)

**Go Version:** Toutes
**Source:** [Code Review Comments](https://go.dev/wiki/CodeReviewComments)

```go
// INCORRECT
items := []string{}           // want "KTN-VAR-020"
items := make([]string, 0)    // want "KTN-VAR-020"

// CORRECT
var items []string
```

**Exception:** JSON encoding où `[]` est préféré à `null`

---

### 2. VAR-021: Receiver type consistency (Confirmé)

**Go Version:** Toutes
**Source:** [Code Review Comments](https://go.dev/wiki/CodeReviewComments)

```go
// INCORRECT - mixing receivers
type Server struct{}
func (s *Server) Start() {}
func (s Server) Stop() {}    // want "KTN-VAR-021: mixed receiver types"

// CORRECT
func (s *Server) Start() {}
func (s *Server) Stop() {}
```

---

### 3. VAR-022: Pointer to interface (Confirmé)

**Go Version:** Toutes
**Source:** [Effective Go](https://go.dev/doc/effective_go)

```go
// INCORRECT
func Process(r *io.Reader) {}      // want "KTN-VAR-022"
var w *io.Writer                   // want "KTN-VAR-022"

// CORRECT
func Process(r io.Reader) {}
var w io.Writer
```

**Exception:** `*error` pour modifier une erreur existante

---

### 4. VAR-023: crypto/rand pour secrets (SÉCURITÉ)

**Go Version:** Toutes
**Source:** [Code Review Comments](https://go.dev/wiki/CodeReviewComments)

```go
// DANGEREUX - prédictible!
import "math/rand"
key := rand.Intn(1000000)          // want "KTN-VAR-023: use crypto/rand for secrets"
token := rand.Int63()              // want "KTN-VAR-023"

// CORRECT
import "crypto/rand"
key, _ := rand.Int(rand.Reader, big.NewInt(1000000))
```

**Détection:** Appels à `math/rand.*` dans contexte "key", "token", "secret", "password"

---

### 5. VAR-024: Utiliser `any` au lieu de `interface{}` (Go 1.18+)

**Go Version:** 1.18+
**Source:** [Go 1.18 Release Notes](https://go.dev/doc/go1.18)

```go
// OBSOLÈTE
func Process(data interface{}) {}  // want "KTN-VAR-024: use 'any' instead of 'interface{}'"
var x interface{}                  // want "KTN-VAR-024"

// MODERNE
func Process(data any) {}
var x any
```

**Note:** `any` est un alias de `interface{}` depuis Go 1.18

---

### 6. VAR-025: Utiliser `clear()` built-in (Go 1.21+)

**Go Version:** 1.21+
**Source:** [Go 1.21 Release Notes](https://go.dev/doc/go1.21)

```go
// OBSOLÈTE
for k := range m {
    delete(m, k)                   // want "KTN-VAR-025: use clear(m)"
}

for i := range s {
    s[i] = 0                       // want "KTN-VAR-025: use clear(s)"
}

// MODERNE
clear(m)  // Supprime tous éléments
clear(s)  // Met tous éléments à zero value
```

---

### 7. VAR-026: Utiliser `min()`/`max()` built-in (Go 1.21+)

**Go Version:** 1.21+
**Source:** [Go 1.21 Release Notes](https://go.dev/doc/go1.21)

```go
// OBSOLÈTE
func minInt(a, b int) int {
    if a < b {                     // want "KTN-VAR-026: use built-in min()"
        return a
    }
    return b
}

// OBSOLÈTE
import "math"
x := math.Min(float64(a), float64(b))  // want "KTN-VAR-026"

// MODERNE
x := min(a, b)
y := max(a, b, c)  // Supporte N arguments
```

---

### 8. VAR-027: Range over integer (Go 1.22+)

**Go Version:** 1.22+
**Source:** [Go 1.22 Release Notes](https://go.dev/doc/go1.22)

```go
// OBSOLÈTE
for i := 0; i < n; i++ {           // want "KTN-VAR-027: use 'range n'"
    process(i)
}

for i := 0; i < 10; i++ {          // want "KTN-VAR-027"
    fmt.Println(i)
}

// MODERNE
for i := range n {
    process(i)
}

for i := range 10 {
    fmt.Println(i)
}
```

**Condition:** Simple loop `for i := 0; i < N; i++` sans modification de `i`

---

### 9. VAR-028: Loop variable copy obsolète (Go 1.22+)

**Go Version:** 1.22+
**Source:** [Fixing For Loops](https://go.dev/blog/loopvar-preview)

```go
// OBSOLÈTE (workaround plus nécessaire)
for _, v := range items {
    v := v                         // want "KTN-VAR-028: unnecessary copy in Go 1.22+"
    go process(v)
}

// MODERNE (Go 1.22+ scope per-iteration)
for _, v := range items {
    go process(v)  // Safe - v has per-iteration scope
}
```

**Détection:** Pattern `x := x` immédiatement après range variable

---

### 10. VAR-029: Utiliser slices.Grow (Go 1.21+)

**Go Version:** 1.21+
**Source:** [slices package](https://pkg.go.dev/slices)

```go
// OBSOLÈTE
if cap(s)-len(s) < n {
    s = append(s, make([]T, n)...)[:len(s)]  // Complex
}

// ou
newCap := len(s) + n
if newCap > cap(s) {
    newSlice := make([]T, len(s), newCap)
    copy(newSlice, s)
    s = newSlice
}

// MODERNE
s = slices.Grow(s, n)
```

---

### 11. VAR-030: Utiliser slices.Clone (Go 1.21+)

**Go Version:** 1.21+
**Source:** [slices package](https://pkg.go.dev/slices)

```go
// OBSOLÈTE
clone := make([]T, len(original))
copy(clone, original)              // want "KTN-VAR-030: use slices.Clone"

// ou
clone := append([]T(nil), original...)

// MODERNE
clone := slices.Clone(original)
```

---

### 12. VAR-031: Utiliser maps.Clone (Go 1.21+)

**Go Version:** 1.21+
**Source:** [maps package](https://pkg.go.dev/maps)

```go
// OBSOLÈTE
clone := make(map[K]V, len(original))
for k, v := range original {
    clone[k] = v                   // want "KTN-VAR-031: use maps.Clone"
}

// MODERNE
clone := maps.Clone(original)
```

---

### 13. VAR-032: Utiliser slices.Delete correctement (Go 1.21+)

**Go Version:** 1.21+
**Source:** [Generic Slice Functions Blog](https://go.dev/blog/generic-slice-functions)

```go
// INCORRECT - ignore return value
slices.Delete(s, i, j)             // want "KTN-VAR-032: slices.Delete return must be used"

// CORRECT
s = slices.Delete(s, i, j)
```

**Note:** S'applique aussi à: `Compact`, `CompactFunc`, `DeleteFunc`, `Grow`, `Insert`, `Replace`

---

### 14. VAR-033: Utiliser cmp.Or (Go 1.22+)

**Go Version:** 1.22+
**Source:** [cmp package](https://pkg.go.dev/cmp)

```go
// OBSOLÈTE
func getPort(port int) int {
    if port != 0 {                 // want "KTN-VAR-033: use cmp.Or"
        return port
    }
    return 8080
}

// MODERNE
port := cmp.Or(port, 8080)
host := cmp.Or(configHost, envHost, "localhost")
```

---

### 15. VAR-034: Utiliser sync.WaitGroup.Go (Go 1.25+)

**Go Version:** 1.25+
**Source:** [Go 1.25 Release Notes](https://go.dev/doc/go1.25)

```go
// OBSOLÈTE
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(item T) {              // want "KTN-VAR-034: use wg.Go()"
        defer wg.Done()
        process(item)
    }(item)
}
wg.Wait()

// MODERNE
var wg sync.WaitGroup
for _, item := range items {
    wg.Go(func() {
        process(item)
    })
}
wg.Wait()
```

---

### 16. VAR-035: Préférer slices.Contains (Go 1.21+)

**Go Version:** 1.21+
**Source:** [slices package](https://pkg.go.dev/slices)

```go
// OBSOLÈTE
found := false
for _, v := range items {
    if v == target {               // want "KTN-VAR-035: use slices.Contains"
        found = true
        break
    }
}

// MODERNE
found := slices.Contains(items, target)
```

---

### 17. VAR-036: Préférer slices.Index (Go 1.21+)

**Go Version:** 1.21+
**Source:** [slices package](https://pkg.go.dev/slices)

```go
// OBSOLÈTE
index := -1
for i, v := range items {
    if v == target {               // want "KTN-VAR-036: use slices.Index"
        index = i
        break
    }
}

// MODERNE
index := slices.Index(items, target)
```

---

### 18. VAR-037: Préférer maps.Keys/Values (Go 1.23+)

**Go Version:** 1.23+
**Source:** [maps package iterators](https://pkg.go.dev/maps)

```go
// OBSOLÈTE
var keys []K
for k := range m {                 // want "KTN-VAR-037: use maps.Keys or slices.Collect"
    keys = append(keys, k)
}

// MODERNE
keys := slices.Collect(maps.Keys(m))
values := slices.Collect(maps.Values(m))
```

---

### 19. GENERIC-001: Contrainte comparable vs any

**Go Version:** 1.18+
**Source:** [Generics Intro](https://go.dev/blog/intro-generics)

```go
// TROP PERMISSIF
func Contains[T any](s []T, v T) bool {  // any ne garantit pas ==
    for _, x := range s {
        if x == v {  // Compile error si T ne supporte pas ==
            return true
        }
    }
}

// CORRECT
func Contains[T comparable](s []T, v T) bool {
    // ...
}
```

---

### 20. GENERIC-002: Utiliser slices/maps au lieu de ré-implémenter

**Go Version:** 1.21+
**Source:** [slices package](https://pkg.go.dev/slices)

```go
// OBSOLÈTE - réimplémentation
func Contains[T comparable](s []T, v T) bool {
    for _, x := range s {          // want "GENERIC-002: use slices.Contains"
        if x == v { return true }
    }
    return false
}

// MODERNE
import "slices"
found := slices.Contains(s, v)
```

---

## Priorités d'Implémentation

### Phase 1 - Haute priorité (Impact immédiat)

| Règle | Concept | Effort |
|-------|---------|--------|
| VAR-020 | nil slice | Low |
| VAR-021 | Receiver consistency | Medium |
| VAR-023 | crypto/rand (sécurité) | Low |
| VAR-024 | any vs interface{} | Low |
| VAR-028 | Loop var copy obsolète | Low |

### Phase 2 - Moyenne priorité (Modernisation)

| Règle | Concept | Effort |
|-------|---------|--------|
| VAR-025 | clear() | Medium |
| VAR-026 | min()/max() | Medium |
| VAR-027 | range over int | Medium |
| VAR-030 | slices.Clone | Low |
| VAR-031 | maps.Clone | Low |

### Phase 3 - Basse priorité (Optimisation)

| Règle | Concept | Effort |
|-------|---------|--------|
| VAR-029 | slices.Grow | High |
| VAR-032 | slices.Delete return | Medium |
| VAR-033 | cmp.Or | Medium |
| VAR-034 | WaitGroup.Go (1.25) | Low |
| VAR-035 | slices.Contains | Medium |
| VAR-036 | slices.Index | Medium |
| VAR-037 | maps.Keys/Values | Medium |

---

## Sources Officielles

| Document | URL |
|----------|-----|
| Go 1.18 Release Notes | https://go.dev/doc/go1.18 |
| Go 1.21 Release Notes | https://go.dev/doc/go1.21 |
| Go 1.22 Release Notes | https://go.dev/doc/go1.22 |
| Go 1.23 Release Notes | https://go.dev/doc/go1.23 |
| Go 1.24 Release Notes | https://go.dev/doc/go1.24 |
| Go 1.25 Release Notes | https://go.dev/doc/go1.25 |
| Generics Intro | https://go.dev/blog/intro-generics |
| Loop Var Fix | https://go.dev/blog/loopvar-preview |
| Generic Slice Functions | https://go.dev/blog/generic-slice-functions |
| slices package | https://pkg.go.dev/slices |
| maps package | https://pkg.go.dev/maps |
| cmp package | https://pkg.go.dev/cmp |

---

## Notes d'Implémentation

### Détection de version Go

Toutes les règles basées sur une version spécifique doivent:
1. Lire `go.mod` pour extraire la version minimum
2. Ne reporter que si `go X.Y` dans go.mod >= version requise

```go
func isGoVersionAtLeast(pass *analysis.Pass, major, minor int) bool {
    // Lire go.mod via pass.Module ou pkg.Module
}
```

### Catégories de règles

| Préfixe | Catégorie |
|---------|-----------|
| VAR-0XX | Variables générales |
| GENERIC-0XX | Génériques (types paramétrés) |
| ITER-0XX | Itérateurs (range-over-func) |
| SYNC-0XX | Concurrence (sync package) |

---

_Document de planification. Ne pas commiter._
