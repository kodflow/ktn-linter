# pkg/analyzer/ktn/ktninterface/ - Interface Rules

## Purpose
Analyze interface declarations for size, design, and Go idioms.

## Rules (1 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-INTERFACE-001 | Interface too large (max 5 methods) | Warning |

## Go Interface Philosophy
"The bigger the interface, the weaker the abstraction." - Rob Pike

Small interfaces are preferred in Go:
- `io.Reader` has 1 method
- `io.Writer` has 1 method
- `io.Closer` has 1 method

## File Structure
```
ktninterface/
├── 001.go              # KTN-INTERFACE-001 implementation
├── 001_external_test.go
├── registry.go         # Analyzers()
└── testdata/src/interface001/
    ├── good.go
    └── bad.go
```

## Implementation
```go
func runInterface001(pass *analysis.Pass) (any, error) {
    inspect.Preorder([]ast.Node{(*ast.InterfaceType)(nil)}, func(n ast.Node) {
        iface := n.(*ast.InterfaceType)
        methodCount := len(iface.Methods.List)

        if methodCount > 5 {
            pass.Reportf(n.Pos(), "KTN-INTERFACE-001: ...")
        }
    })
    return nil, nil
}
```

## Testdata Example
```go
// bad.go
type badLargeInterface interface { // want "KTN-INTERFACE-001"
    Method1()
    Method2()
    Method3()
    Method4()
    Method5()
    Method6()
}

// good.go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

## Composition Over Large Interfaces
```go
// Instead of one large interface, compose small ones
type ReadWriter interface {
    Reader
    Writer
}
```
