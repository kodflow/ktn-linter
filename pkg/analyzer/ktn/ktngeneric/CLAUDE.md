# pkg/analyzer/ktn/ktngeneric/ - Generic Rules

## Purpose
Analyze generic functions and types for proper constraint usage and Go 1.18+ best practices.

## Rules (5 total)
| Rule | Description | Severity | Go Version |
|------|-------------|----------|------------|
| KTN-GENERIC-001 | `==` or `!=` requires `comparable` constraint | Error | 1.18+ |
| KTN-GENERIC-002 | Unnecessary generics on interface types | Warning | 1.18+ |
| KTN-GENERIC-003 | `golang.org/x/exp/constraints` deprecated → `cmp` | Warning | 1.21+ |
| KTN-GENERIC-005 | Type params must not shadow predeclared identifiers | Warning | 1.18+ |
| KTN-GENERIC-006 | `<`, `>`, `+`, `-`, `*`, `/`, `%` require `cmp.Ordered` | Error | 1.21+ |

## File Structure
```
ktngeneric/
├── 001.go ... 006.go           # Rule implementations (note: 004 skipped)
├── *_internal_test.go          # White-box tests
├── *_external_test.go          # Black-box tests
├── registry.go                 # Analyzers()
└── testdata/src/
    └── generic001 ... generic006/  # Test fixtures
```

## KTN-GENERIC-001: Comparable Constraint
```go
// BAD - any constraint doesn't support ==
func Contains[T any](s []T, v T) bool {
    for _, x := range s {
        if x == v { return true }  // ERROR: == not allowed on any
    }
    return false
}

// GOOD - comparable constraint
func Contains[T comparable](s []T, v T) bool {
    for _, x := range s {
        if x == v { return true }  // OK
    }
    return false
}
```

## KTN-GENERIC-002: Unnecessary Generics
```go
// BAD - generic adds no value, just use io.Reader directly
func Process[T io.Reader](r T) { ... }

// GOOD - use interface directly
func Process(r io.Reader) { ... }
```

## KTN-GENERIC-003: Deprecated Package
```go
// BAD - x/exp/constraints is deprecated
import "golang.org/x/exp/constraints"
func Max[T constraints.Ordered](a, b T) T { ... }

// GOOD - use cmp package (Go 1.21+)
import "cmp"
func Max[T cmp.Ordered](a, b T) T { ... }
```

## KTN-GENERIC-005: Predeclared Shadowing
```go
// BAD - shadows built-in 'int' type
func Process[int any](x int) { ... }

// GOOD - use standard type parameter names
func Process[T any](x T) { ... }
```

## KTN-GENERIC-006: Ordered Constraint
```go
// BAD - any constraint doesn't support < > + - * / %
func Max[T any](a, b T) T {
    if a > b { return a }  // ERROR: > not allowed on any
    return b
}

// GOOD - cmp.Ordered constraint
func Max[T cmp.Ordered](a, b T) T {
    if a > b { return a }  // OK
    return b
}
```

## Dependencies
Uses `pkg/config`, `pkg/messages`
