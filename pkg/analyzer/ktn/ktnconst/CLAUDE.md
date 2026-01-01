# pkg/analyzer/ktn/ktnconst/ - Constant Rules

## Purpose
Analyze constant declarations for explicit types, organization, and naming conventions.

## Rules (3 total)
| Rule | Description | Severity | Errors |
|------|-------------|----------|--------|
| KTN-CONST-001 | Constants must have explicit types | Error | 47 |
| KTN-CONST-002 | Constants grouped at top (before var/type/func) | Info | 7 |
| KTN-CONST-003 | Constants must use CamelCase (no underscores) | Info | 44 |

## File Structure
```
ktnconst/
├── 001.go, 002.go, 003.go      # Rule implementations
├── 001_internal_test.go         # White-box tests
├── 001_external_test.go         # Black-box tests
├── registry.go                  # Analyzers()
└── testdata/src/
    ├── const001/                # 47 error cases
    ├── const002/                # 7 error cases
    └── const003/                # 44 error cases
```

## KTN-CONST-001: Explicit Types

**Rule**: All constants must have explicit type declarations.

```go
// CORRECT
const MaxSize int = 100
const APIKey string = "secret"

// WRONG - implicit type
const maxSize = 100  // want "KTN-CONST-001"
```

**Exception**: iota inheritance (subsequent lines inherit type from first)
```go
const (
    StatusOK int = iota  // explicit type
    StatusErr            // inherits int - OK
)
```

## KTN-CONST-002: Declaration Order

**Rule**: Constants must be placed at the top of the file (before var/type/func).
Multiple const blocks at the top are allowed.

```go
// CORRECT - all const at top
const A int = 1
const B int = 2
var x = A  // var after const

// WRONG - const after var
var x = 1
const A int = 1  // want "KTN-CONST-002"
```

**Exception**: const with iota using custom type after type declaration
```go
type Status int
const (
    StatusOK Status = iota  // OK - uses custom type
    StatusErr
)
```

## KTN-CONST-003: CamelCase Naming

**Rule**: Constants must use CamelCase (MixedCaps), not underscores.

```go
// CORRECT
const MaxSize int = 100
const APIEndpoint string = "/api"
const maxInternalSize int = 50  // unexported

// WRONG
const MAX_SIZE int = 100      // want "KTN-CONST-003"
const max_size int = 100      // want "KTN-CONST-003"
const Max_Size int = 100      // want "KTN-CONST-003"
```

**Valid patterns**:
- PascalCase: `MaxSize`, `APIKey`, `HTTPStatus`
- camelCase: `maxSize`, `apiKey` (unexported)
- Single letters: `A`, `i`, `n`
- With numbers: `Http2`, `Version100`
- All caps (no underscore): `MAXSIZE` (valid but unusual)

## Sources
- [Effective Go - MixedCaps](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
