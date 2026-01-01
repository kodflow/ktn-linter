# Plan d'Implémentation - Règles VAR

Generated: 2026-01-01
Total: 22 règles (4 MODIFY + 18 ADD)

---

## Phase 1: Règles à MODIFIER (4)

### TASK-001: VAR-004 - Scope-aware length check

**Type:** MODIFY
**File:** `pkg/analyzer/ktn/ktnvar/004.go`
**Effort:** Medium

**Current behavior:** Min 2 chars pour toutes variables
**New behavior:**
- Package-level: min 2 chars
- Function-level: 1 char OK si idiomatique

**Noms 1-char autorisés en local:**
```go
// Loop counters
i, j, k, n

// Type hints
b, c, f, m, r, s, t, w

// Results
ok, _
```

**Testdata bad.go:**
```go
package var004

var a int = 1           // Package scope - want error
var b string = "x"      // Package scope - want error

func example() {
    var x int = 42      // 1 char non-idiomatique - want error
    var q string = "q"  // q pas idiomatique - want error
}
```

**Testdata good.go:**
```go
package var004

var count int = 1
var name string = "x"

func example() {
    for i := 0; i < 10; i++ {}  // i OK
    r := strings.NewReader("")   // r OK
    ok := condition              // ok OK
    m := make(map[string]int)    // m OK
}
```

---

### TASK-002: VAR-007 - Zero value exception

**Type:** MODIFY
**File:** `pkg/analyzer/ktn/ktnvar/007.go`
**Effort:** Low

**Current behavior:** Reporte tout `var x Type`
**New behavior:** Ne reporte que `var x Type = value` (avec init explicite)

**Testdata bad.go:**
```go
package var007

func example() {
    var x int = 42        // want error - use :=
    var s string = "test" // want error - use :=
    var err error = nil   // want error - redundant nil
}
```

**Testdata good.go:**
```go
package var007

func example() {
    var err error         // OK - zero value intentionnelle
    var buf bytes.Buffer  // OK - zero value usable
    var wg sync.WaitGroup // OK - zero value
    x := 42               // OK - short syntax
}
```

---

### TASK-003: VAR-013 - 64 bytes threshold

**Type:** MODIFY
**File:** `pkg/analyzer/ktn/ktnvar/013.go`
**Effort:** Low

**Current behavior:** Seuil variable (64 par défaut)
**New behavior:** Seuil fixe 64 bytes (L1 cache line)

**Testdata:** Ajuster les structs pour être précisément autour de 64 bytes

---

### TASK-004: VAR-018 - Constant ≤64 bytes only

**Type:** MODIFY
**File:** `pkg/analyzer/ktn/ktnvar/018.go`
**Effort:** Medium

**Current behavior:** Suggère array pour toute taille constante
**New behavior:** Seulement si constant ET ≤64 bytes

**Testdata bad.go:**
```go
package var018

func example() {
    buf := make([]byte, 32)  // want error - use [32]byte
    arr := make([]int, 8)    // want error - use [8]int (64 bytes)
}
```

**Testdata good.go:**
```go
package var018

func example() {
    buf := make([]byte, 128)  // OK - >64 bytes, heap acceptable
    var small [32]byte        // OK - array

    n := 32
    dynamic := make([]byte, n) // OK - dynamic size
}
```

---

## Phase 2: Règles à AJOUTER - Core (4)

### TASK-005: VAR-020 - Nil slice preferred

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/020.go`
**Effort:** Low

**Logic:** Détecte `[]T{}` et `make([]T, 0)` sans capacité

**Testdata bad.go:**
```go
package var020

func example() {
    items := []string{}           // want "KTN-VAR-020"
    data := make([]int, 0)        // want "KTN-VAR-020"
    list := []Item{}              // want "KTN-VAR-020"
}
```

**Testdata good.go:**
```go
package var020

func example() {
    var items []string           // OK - nil slice
    data := make([]int, 0, 10)   // OK - has capacity
    var list []Item              // OK - nil slice
}
```

---

### TASK-006: VAR-021 - Receiver consistency

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/021.go`
**Effort:** Medium

**Logic:** Tous les methods d'un type doivent avoir même receiver type

**Testdata bad.go:**
```go
package var021

type Server struct{ data int }

func (s *Server) Start() {}
func (s Server) Stop() {}     // want "KTN-VAR-021: mixed receiver"
func (s *Server) Restart() {}
func (s Server) Status() int { return 0 }  // want "KTN-VAR-021"
```

**Testdata good.go:**
```go
package var021

type Server struct{ data int }

func (s *Server) Start() {}
func (s *Server) Stop() {}
func (s *Server) Restart() {}
func (s *Server) Status() int { return s.data }

type Point struct{ X, Y int }

func (p Point) Add(other Point) Point { return Point{p.X+other.X, p.Y+other.Y} }
func (p Point) String() string { return fmt.Sprintf("(%d,%d)", p.X, p.Y) }
```

---

### TASK-007: VAR-022 - Pointer to interface

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/022.go`
**Effort:** Low

**Logic:** Détecte `*io.Reader`, `*io.Writer`, `*interface{}`, etc.

**Testdata bad.go:**
```go
package var022

import "io"

func Process(r *io.Reader) {}      // want "KTN-VAR-022"
func Handle(w *io.Writer) {}       // want "KTN-VAR-022"
var handler *interface{}           // want "KTN-VAR-022"

type Service struct {
    logger *Logger                 // OK if Logger is concrete
    reader *io.Reader              // want "KTN-VAR-022"
}
```

**Testdata good.go:**
```go
package var022

import "io"

func Process(r io.Reader) {}
func Handle(w io.Writer) {}
var handler interface{}

type Service struct {
    reader io.Reader
}
```

---

### TASK-008: VAR-023 - crypto/rand for secrets

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/023.go`
**Effort:** Medium

**Logic:** Détecte `math/rand` dans contexte "key", "token", "secret", "password", "salt", "nonce"

**Testdata bad.go:**
```go
package var023

import "math/rand"

func generateKey() int {
    return rand.Intn(1000000)      // want "KTN-VAR-023"
}

func createToken() string {
    token := rand.Int63()          // want "KTN-VAR-023"
    return fmt.Sprintf("%d", token)
}

var secretKey = rand.Uint64()      // want "KTN-VAR-023"
```

**Testdata good.go:**
```go
package var023

import (
    "crypto/rand"
    "math/big"
    mathrand "math/rand"
)

func generateKey() *big.Int {
    key, _ := rand.Int(rand.Reader, big.NewInt(1000000))
    return key
}

func shuffleItems(items []int) {
    mathrand.Shuffle(len(items), func(i, j int) {
        items[i], items[j] = items[j], items[i]
    })  // OK - not security context
}
```

---

## Phase 3: Règles Go 1.18+ (2)

### TASK-009: VAR-024 - any vs interface{}

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/024.go`
**Effort:** Low
**Go Version:** 1.18+

**Logic:** Remplacer `interface{}` par `any`

**Testdata bad.go:**
```go
package var024

func Process(data interface{}) {}  // want "KTN-VAR-024"
var x interface{}                  // want "KTN-VAR-024"

type Container struct {
    value interface{}              // want "KTN-VAR-024"
}

func returns() interface{} {       // want "KTN-VAR-024"
    return nil
}
```

**Testdata good.go:**
```go
package var024

func Process(data any) {}
var x any

type Container struct {
    value any
}

func returns() any {
    return nil
}
```

---

### TASK-010: GENERIC-001 - comparable constraint

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktngeneric/001.go`
**Effort:** Medium
**Go Version:** 1.18+

**Logic:** Si generic func utilise `==` ou `!=`, contrainte doit être `comparable`

**Testdata bad.go:**
```go
package generic001

func Contains[T any](s []T, v T) bool {  // want "GENERIC-001"
    for _, x := range s {
        if x == v { return true }  // == requires comparable
    }
    return false
}

func Index[T any](s []T, v T) int {      // want "GENERIC-001"
    for i, x := range s {
        if x == v { return i }
    }
    return -1
}
```

**Testdata good.go:**
```go
package generic001

func Contains[T comparable](s []T, v T) bool {
    for _, x := range s {
        if x == v { return true }
    }
    return false
}

func Map[T, U any](s []T, f func(T) U) []U {  // OK - no == used
    result := make([]U, len(s))
    for i, x := range s {
        result[i] = f(x)
    }
    return result
}
```

---

## Phase 4: Règles Go 1.21+ (7)

### TASK-011: VAR-025 - clear() built-in

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/025.go`
**Effort:** Medium
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var025

func clearMap(m map[string]int) {
    for k := range m {
        delete(m, k)               // want "KTN-VAR-025: use clear(m)"
    }
}

func zeroSlice(s []int) {
    for i := range s {
        s[i] = 0                   // want "KTN-VAR-025: use clear(s)"
    }
}
```

**Testdata good.go:**
```go
package var025

func clearMap(m map[string]int) {
    clear(m)
}

func zeroSlice(s []int) {
    clear(s)
}
```

---

### TASK-012: VAR-026 - min()/max() built-in

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/026.go`
**Effort:** Medium
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var026

func minInt(a, b int) int {
    if a < b {                     // want "KTN-VAR-026"
        return a
    }
    return b
}

import "math"

func maxFloat(a, b float64) float64 {
    return math.Max(a, b)          // want "KTN-VAR-026"
}
```

**Testdata good.go:**
```go
package var026

func example() {
    x := min(a, b)
    y := max(a, b, c)
    z := min(1, 2, 3, 4, 5)
}
```

---

### TASK-013: VAR-029 - slices.Grow

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/029.go`
**Effort:** High
**Go Version:** 1.21+

**Note:** Complexe - détecte pattern de grow manuel

---

### TASK-014: VAR-030 - slices.Clone

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/030.go`
**Effort:** Low
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var030

func example(original []int) []int {
    clone := make([]int, len(original))
    copy(clone, original)          // want "KTN-VAR-030"
    return clone
}

func example2(original []string) []string {
    return append([]string(nil), original...)  // want "KTN-VAR-030"
}
```

**Testdata good.go:**
```go
package var030

import "slices"

func example(original []int) []int {
    return slices.Clone(original)
}
```

---

### TASK-015: VAR-031 - maps.Clone

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/031.go`
**Effort:** Low
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var031

func example(original map[string]int) map[string]int {
    clone := make(map[string]int, len(original))
    for k, v := range original {
        clone[k] = v               // want "KTN-VAR-031"
    }
    return clone
}
```

**Testdata good.go:**
```go
package var031

import "maps"

func example(original map[string]int) map[string]int {
    return maps.Clone(original)
}
```

---

### TASK-016: VAR-035 - slices.Contains

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/035.go`
**Effort:** Medium
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var035

func contains(items []string, target string) bool {
    for _, v := range items {
        if v == target {           // want "KTN-VAR-035"
            return true
        }
    }
    return false
}
```

**Testdata good.go:**
```go
package var035

import "slices"

func contains(items []string, target string) bool {
    return slices.Contains(items, target)
}
```

---

### TASK-017: VAR-036 - slices.Index

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/036.go`
**Effort:** Medium
**Go Version:** 1.21+

**Testdata bad.go:**
```go
package var036

func indexOf(items []int, target int) int {
    for i, v := range items {
        if v == target {           // want "KTN-VAR-036"
            return i
        }
    }
    return -1
}
```

**Testdata good.go:**
```go
package var036

import "slices"

func indexOf(items []int, target int) int {
    return slices.Index(items, target)
}
```

---

## Phase 5: Règles Go 1.22+ (3)

### TASK-018: VAR-027 - range over integer

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/027.go`
**Effort:** Medium
**Go Version:** 1.22+

**Testdata bad.go:**
```go
package var027

func example(n int) {
    for i := 0; i < n; i++ {       // want "KTN-VAR-027"
        process(i)
    }

    for i := 0; i < 10; i++ {      // want "KTN-VAR-027"
        fmt.Println(i)
    }
}
```

**Testdata good.go:**
```go
package var027

func example(n int) {
    for i := range n {
        process(i)
    }

    for i := range 10 {
        fmt.Println(i)
    }

    // OK - modifies i
    for i := 0; i < n; i += 2 {
        process(i)
    }
}
```

---

### TASK-019: VAR-028 - loop var copy obsolete

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/028.go`
**Effort:** Low
**Go Version:** 1.22+

**Testdata bad.go:**
```go
package var028

func example(items []int) {
    for _, v := range items {
        v := v                     // want "KTN-VAR-028"
        go process(v)
    }

    for i, item := range items {
        i := i                     // want "KTN-VAR-028"
        item := item               // want "KTN-VAR-028"
        go func() {
            fmt.Println(i, item)
        }()
    }
}
```

**Testdata good.go:**
```go
package var028

func example(items []int) {
    for _, v := range items {
        go process(v)  // Safe in Go 1.22+
    }

    for i, item := range items {
        go func() {
            fmt.Println(i, item)  // Safe in Go 1.22+
        }()
    }
}
```

---

### TASK-020: VAR-033 - cmp.Or

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/033.go`
**Effort:** Medium
**Go Version:** 1.22+

**Testdata bad.go:**
```go
package var033

func getPort(port int) int {
    if port != 0 {                 // want "KTN-VAR-033"
        return port
    }
    return 8080
}

func getHost(host string) string {
    if host != "" {                // want "KTN-VAR-033"
        return host
    }
    return "localhost"
}
```

**Testdata good.go:**
```go
package var033

import "cmp"

func getPort(port int) int {
    return cmp.Or(port, 8080)
}

func getHost(host string) string {
    return cmp.Or(host, "localhost")
}
```

---

## Phase 6: Règles Go 1.23+ et 1.25+ (2)

### TASK-021: VAR-037 - maps.Keys/Values

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/037.go`
**Effort:** Medium
**Go Version:** 1.23+

**Testdata bad.go:**
```go
package var037

func getKeys(m map[string]int) []string {
    var keys []string
    for k := range m {
        keys = append(keys, k)     // want "KTN-VAR-037"
    }
    return keys
}
```

**Testdata good.go:**
```go
package var037

import (
    "maps"
    "slices"
)

func getKeys(m map[string]int) []string {
    return slices.Collect(maps.Keys(m))
}
```

---

### TASK-022: VAR-034 - WaitGroup.Go

**Type:** ADD
**File:** `pkg/analyzer/ktn/ktnvar/034.go`
**Effort:** Low
**Go Version:** 1.25+

**Testdata bad.go:**
```go
package var034

import "sync"

func example(items []int) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func(item int) {        // want "KTN-VAR-034"
            defer wg.Done()
            process(item)
        }(item)
    }
    wg.Wait()
}
```

**Testdata good.go:**
```go
package var034

import "sync"

func example(items []int) {
    var wg sync.WaitGroup
    for _, item := range items {
        wg.Go(func() {
            process(item)
        })
    }
    wg.Wait()
}
```

---

## Récapitulatif Exécution

### Batch 1: MODIFY (4 agents parallèles)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-001 | VAR-004 | Medium |
| TASK-002 | VAR-007 | Low |
| TASK-003 | VAR-013 | Low |
| TASK-004 | VAR-018 | Medium |

### Batch 2: ADD Core (4 agents parallèles)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-005 | VAR-020 | Low |
| TASK-006 | VAR-021 | Medium |
| TASK-007 | VAR-022 | Low |
| TASK-008 | VAR-023 | Medium |

### Batch 3: Go 1.18+ (2 agents parallèles)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-009 | VAR-024 | Low |
| TASK-010 | GENERIC-001 | Medium |

### Batch 4: Go 1.21+ (7 agents, par groupes de 4)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-011 | VAR-025 | Medium |
| TASK-012 | VAR-026 | Medium |
| TASK-013 | VAR-029 | High |
| TASK-014 | VAR-030 | Low |
| TASK-015 | VAR-031 | Low |
| TASK-016 | VAR-035 | Medium |
| TASK-017 | VAR-036 | Medium |

### Batch 5: Go 1.22+ (3 agents parallèles)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-018 | VAR-027 | Medium |
| TASK-019 | VAR-028 | Low |
| TASK-020 | VAR-033 | Medium |

### Batch 6: Go 1.23+/1.25+ (2 agents parallèles)
| Task | Règle | Effort |
|------|-------|--------|
| TASK-021 | VAR-037 | Medium |
| TASK-022 | VAR-034 | Low |

---

## Template Agent Prompt

Pour chaque agent, utiliser ce template:

```
Implémenter la règle KTN-{RULE_CODE}

FICHIERS À CRÉER:
1. pkg/analyzer/ktn/ktnvar/{NNN}.go - Implémentation
2. pkg/analyzer/ktn/ktnvar/{NNN}_external_test.go - Tests
3. pkg/analyzer/ktn/ktnvar/testdata/src/var{NNN}/good.go
4. pkg/analyzer/ktn/ktnvar/testdata/src/var{NNN}/bad.go

SPÉCIFICATION:
{SPEC_FROM_PLAN}

TESTDATA BAD.GO:
{BAD_GO_CONTENT}

TESTDATA GOOD.GO:
{GOOD_GO_CONTENT}

CONTRAINTES:
- Pas de verbosité
- Suivre le pattern des règles existantes
- Ajouter au registry.go
- Ajouter message dans pkg/messages/var.go
- Vérifier avec: go test ./pkg/analyzer/ktn/ktnvar/...
```

---

_Plan prêt pour /apply_
