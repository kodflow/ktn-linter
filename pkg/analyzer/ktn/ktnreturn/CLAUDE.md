# pkg/analyzer/ktn/ktnreturn/ - Return Rules

## Purpose
Analyze return statements for clarity and maintainability.

## Rules (1 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-RETURN-001 | Naked return in function with named results | Warning |

## Naked Returns Problem
```go
// BAD - Unclear what's being returned
func compute(x int) (result int, err error) {
    result = x * 2
    if x < 0 {
        err = errors.New("negative")
        return // Naked return - what values?
    }
    return // Naked return
}

// GOOD - Explicit returns
func compute(x int) (int, error) {
    result := x * 2
    if x < 0 {
        return 0, errors.New("negative")
    }
    return result, nil
}
```

## When Naked Returns Are Acceptable
- Very short functions (< 5 lines)
- defer statements that modify named returns
- Generated code

## File Structure
```
ktnreturn/
├── 001.go              # KTN-RETURN-001 implementation
├── 001_external_test.go
├── registry.go         # Analyzers()
└── testdata/src/return001/
    ├── good.go
    └── bad.go
```

## Implementation
```go
func runReturn001(pass *analysis.Pass) (any, error) {
    inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
        fn := n.(*ast.FuncDecl)
        if !hasNamedResults(fn) {
            return
        }

        ast.Inspect(fn.Body, func(n ast.Node) bool {
            ret, ok := n.(*ast.ReturnStmt)
            if ok && len(ret.Results) == 0 {
                pass.Reportf(ret.Pos(), "KTN-RETURN-001: ...")
            }
            return true
        })
    })
    return nil, nil
}
```

## Related Rules
- KTN-FUNC-006: Too many return values
- KTN-FUNC-012: Else after return
