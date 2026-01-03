# pkg/analyzer/ktn/ktnfunc/ - Function Rules

## Purpose
Analyze function declarations for style, complexity, and documentation.

## Rules (13 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-FUNC-001 | Function too long (max 35 lines) | Error |
| KTN-FUNC-002 | Too many parameters (max 5) | Error |
| KTN-FUNC-003 | Missing function documentation | Warning |
| KTN-FUNC-004 | Function name too long | Warning |
| KTN-FUNC-005 | Unexported function missing doc | Info |
| KTN-FUNC-006 | Too many return values | Warning |
| KTN-FUNC-007 | Doc missing Params/Returns sections | Warning |
| KTN-FUNC-008 | Cognitive complexity too high | Error |
| KTN-FUNC-009 | Nested function depth exceeded | Warning |
| KTN-FUNC-010 | Empty function body | Warning |
| KTN-FUNC-011 | Missing comments on if/switch/return | Info |
| KTN-FUNC-012 | Else after return (simplify) | Warning |
| KTN-FUNC-013 | Prefer empty slice/map over nil | Warning |

## File Structure
```
ktnfunc/
├── 001.go ... 013.go       # Rule implementations
├── 001_external_test.go    # Black-box tests
├── registry.go             # GetAnalyzers()
└── testdata/src/
    ├── func001/            # Testdata for each rule
    │   ├── good.go
    │   └── bad.go
    └── ...
```

## Adding a New Rule
1. Create `<NNN>.go` with `Analyzer<NNN>`
2. Create `<NNN>_external_test.go`
3. Create `testdata/src/func<NNN>/good.go` + `bad.go`
4. Add `Analyzer<NNN>` to `registry.go`
5. Add messages in `pkg/messages/func.go`

## Testdata Convention
```go
// bad.go - Use prefix to avoid redeclaration
func badTooLongFunction() { ... }  // want "KTN-FUNC-001"

// good.go - Clean names
func shortFunction() { ... }
```

## Key Imports
```go
import (
    "go/ast"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "github.com/kodflow/ktn-linter/pkg/analyzer/utils"
)
```
