# pkg/analyzer/modernize/ - Go Modernize Wrapper

## Purpose
Wrapper around `golang.org/x/tools/go/analysis/passes/modernize` suite.
Provides access to official Go modernization analyzers.

## Structure
```
modernize/
├── registry.go   # Analyzers() - filters/exposes modernize suite
└── (no tests - wraps external package)
```

## Active Analyzers (17 of 18)
The modernize suite includes analyzers for updating code to use newer Go features:

| Analyzer | Purpose |
|----------|---------|
| appendassign | Use `x = append(x, ...)` patterns |
| bloop | Use range-over-int loops (Go 1.22) |
| efaceany | Use `any` instead of `interface{}` |
| fmtappendf | Use `fmt.Appendf` (Go 1.19) |
| minmax | Use `min()`/`max()` built-ins (Go 1.21) |
| omitzero | Use `omitzero` struct tag (Go 1.24) |
| rangeint | Use range over integer |
| slicesclone | Use `slices.Clone` |
| slicesdelete | Use `slices.Delete` |
| slicesinsert | Use `slices.Insert` |
| sortslice | Use `slices.Sort` |
| stringbuilder | Use `strings.Builder` |
| testingcontext | Use `testing.Context` (Go 1.24) |
| waitgroup | Use structured WaitGroup |

## Integration
```go
// registry.go
func Analyzers() []*analysis.Analyzer {
    suite := modernize.Suite
    // Filter or customize as needed
    return suite
}
```

## Usage in Main Registry
```go
// pkg/analyzer/ktn/registry.go
func GetAllRules() []*analysis.Analyzer {
    rules := append(ktnfunc.GetAnalyzers(), ...)
    rules = append(rules, modernize.Analyzers()...)
    return rules
}
```

## Output Prefix
Modernize diagnostics are prefixed with `MODERNIZE-` in output for distinction:
```
file.go:10:5  info  MODERNIZE-efaceany  Use 'any' instead of 'interface{}'
```

## Version Requirements
Some modernize checks only apply when targeting specific Go versions.
The project targets Go 1.25.5, so all modernize features are available.
