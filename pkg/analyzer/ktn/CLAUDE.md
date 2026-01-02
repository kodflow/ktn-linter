# pkg/analyzer/ktn/ - KTN Rules Registry

## Purpose
Central registry aggregating all KTN rule categories. Provides:
- `GetAllRules()` - All 62+ analyzers
- `GetRulesByCategory(cat)` - Filter by category
- `GetRuleByCode(code)` - Single rule lookup

## Categories (11 total)
| Package | Prefix | Count | Focus |
|---------|--------|-------|-------|
| ktnfunc | KTN-FUNC | 12 | Function length, params, docs |
| ktnvar | KTN-VAR | 36 | Variable naming, patterns, modern idioms |
| ktnstruct | KTN-STRUCT | 6 | Struct fields, embedding |
| ktnconst | KTN-CONST | 6 | Explicit types, grouping, naming |
| ktngeneric | KTN-GENERIC | 1 | Generic type constraints |
| ktncomment | KTN-COMMENT | 7 | Comment format, placement |
| ktntest | KTN-TEST | 11 | Test naming, file conventions |
| ktnapi | KTN-API | 1 | API field access patterns |
| ktninterface | KTN-INTERFACE | 1 | Interface size limits |
| ktnreturn | KTN-RETURN | 1 | Naked return detection |
| testhelper | - | 0 | Test utilities (not rules) |

## Registry Pattern
```go
// Each category provides:
func Analyzers() []*analysis.Analyzer { ... }

// Master registry aggregates:
func GetAllRules() []*analysis.Analyzer {
    return append(
        ktnfunc.GetAnalyzers(),
        ktnvar.Analyzers(),
        ktnstruct.Analyzers(),
        // ...
    )
}
```

## Rule Naming Convention
- File: `<NNN>.go` (e.g., `001.go`, `012.go`)
- Analyzer: `Analyzer<NNN>` (e.g., `Analyzer001`)
- Test: `<NNN>_external_test.go` or `<NNN>_internal_test.go`
- Testdata: `testdata/src/<cat><NNN>/` (e.g., `testdata/src/func001/`)

## Testdata Structure
```
testdata/src/<cat><NNN>/
├── good.go   # Must produce 0 errors
└── bad.go    # Must produce ONLY KTN-<CAT>-<NNN> errors
```

## Critical Rules
- NO path-based exclusions (IsTestdataPath, etc.)
- Testdata must be REALLY compliant, not artificially excluded
- Function names in bad.go/good.go must differ to avoid redeclaration
