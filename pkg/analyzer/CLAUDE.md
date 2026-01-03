# pkg/analyzer/ - Static Analysis Engine

## Structure
```
analyzer/
├── ktn/              # KTN rules (80 rules across 10 categories)
│   ├── registry.go   # Master registry: GetAllRules(), GetRuleByCode()
│   ├── ktnfunc/      # Function rules (13): length, params, docs, nil returns
│   ├── ktnvar/       # Variable rules (36): naming, shadowing
│   ├── ktnstruct/    # Struct rules (6): field count, embedding
│   ├── ktnconst/     # Const rules (6): naming, grouping
│   ├── ktngeneric/   # Generic rules (5): type constraints
│   ├── ktncomment/   # Comment rules (7): format, placement
│   ├── ktntest/      # Test rules (11): naming, assertions
│   ├── ktnapi/       # API rules (1): field access
│   ├── ktninterface/ # Interface rules (1): size limits
│   └── testhelper/   # Test utilities for analyzers
├── modernize/        # Wrapper for golang.org/x/tools/go/analysis/passes/modernize
├── utils/            # AST utilities (100% coverage)
└── shared/           # Classification helpers
```

## Adding a New Rule
1. Create `pkg/analyzer/ktn/ktn<cat>/<NNN>.go`
2. Create `pkg/analyzer/ktn/ktn<cat>/<NNN>_external_test.go`
3. Create `testdata/src/<cat><NNN>/good.go` + `bad.go`
4. Add to `pkg/analyzer/ktn/ktn<cat>/registry.go`
5. Add messages in `pkg/messages/<cat>.go`
6. Add severity in `pkg/severity/severity.go`

## Analyzer Template
```go
var Analyzer<NNN> = &analysis.Analyzer{
    Name:     "ktn<cat><NNN>",
    Doc:      "KTN-<CAT>-<NNN>: Description",
    Run:      run<Cat><NNN>,
    Requires: []*analysis.Analyzer{inspect.Analyzer},
}
```

## Testing Convention
- `analysistest.Run()` with testdata/src/<cat><NNN>/
- `good.go`: 0 errors expected
- `bad.go`: Only KTN-<CAT>-<NNN> errors expected (use `// want` comments)

## Dependencies
- `golang.org/x/tools/go/analysis` (analyzer framework)
- `golang.org/x/tools/go/ast/inspector` (AST traversal)
