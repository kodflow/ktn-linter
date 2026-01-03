# pkg/analyzer/ktn/ktnstruct/ - Struct Rules

## Purpose
Analyze struct declarations for field count, embedding, organization, and method receivers.

## Rules (9 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-STRUCT-001 | Too many fields (max 10) | Warning |
| KTN-STRUCT-002 | Missing struct documentation | Warning |
| KTN-STRUCT-003 | Getters should not have "Get" prefix | Warning |
| KTN-STRUCT-004 | Embedded type not first | Info |
| KTN-STRUCT-005 | Multiple embedding levels | Warning |
| KTN-STRUCT-006 | Private field with serialization tag in DTO | Info |
| KTN-STRUCT-007 | DTO exported fields missing json/xml tags | Info |
| KTN-STRUCT-008 | Inconsistent receiver types (pointer vs value) | Warning |
| KTN-STRUCT-009 | Inconsistent receiver names | Warning |

## File Structure
```
ktnstruct/
├── 001.go ... 006.go       # Rule implementations
├── *_external_test.go      # Tests
├── registry.go             # Analyzers()
└── testdata/src/struct001...
```

## Struct Analysis Pattern
```go
func runStruct001(pass *analysis.Pass) (any, error) {
    inspect.Preorder([]ast.Node{(*ast.TypeSpec)(nil)}, func(n ast.Node) {
        ts := n.(*ast.TypeSpec)
        st, ok := ts.Type.(*ast.StructType)
        if !ok {
            return
        }

        fieldCount := len(st.Fields.List)
        if fieldCount > 10 {
            pass.Reportf(ts.Pos(), "KTN-STRUCT-001: ...")
        }
    })
    return nil, nil
}
```

## Field Ordering Convention
1. Embedded types first
2. Exported fields grouped
3. Unexported fields last
4. Related fields adjacent

## Testdata Example
```go
// bad.go
type badTooManyFields struct { // want "KTN-STRUCT-001"
    A, B, C, D, E, F, G, H, I, J, K int
}

// good.go
type User struct {
    ID    int
    Name  string
    Email string
}
```
