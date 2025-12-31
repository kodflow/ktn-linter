# pkg/rules/ - Rule Information Extraction

## Purpose
Extract rule metadata (code, description, examples) from KTN analyzers for
documentation and AI-assisted development.

## Files
| File | Responsibility |
|------|----------------|
| `info.go` | RuleInfo struct, extraction from analyzer Doc strings |
| `examples.go` | Load good.go examples from testdata directories |

## Data Structures
```go
type RuleInfo struct {
    Code        string  // KTN-FUNC-001
    Category    string  // func
    Name        string  // ktnfunc001
    Description string  // Short description
    GoodExample string  // Content from testdata/good.go
}

type RulesOutput struct {
    TotalCount int
    Categories []string
    Rules      []RuleInfo
}
```

## Key Functions
- `GetAllRuleInfos()` - All rules from registry
- `GetRuleInfosByCategory(cat)` - Filter by category
- `GetRuleInfoByCode(code)` - Single rule lookup
- `LoadGoodExamples(infos)` - Enrich with examples

## Testdata Resolution
```
KTN-FUNC-001 → pkg/analyzer/ktn/ktnfunc/testdata/src/func001/good.go
KTN-VAR-002  → pkg/analyzer/ktn/ktnvar/testdata/src/var002/good.go
```

## Dependencies
- `pkg/analyzer/ktn` (registry access)
