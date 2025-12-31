# pkg/analyzer/ktn/ktnvar/ - Variable Rules

## Purpose
Analyze variable declarations for naming conventions, shadowing, and scope.

## Rules (18 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-VAR-001 | Variable name too short (min 2 chars) | Warning |
| KTN-VAR-002 | Variable name too long (max 30 chars) | Warning |
| KTN-VAR-003 | Non-idiomatic naming (use camelCase) | Error |
| KTN-VAR-004 | Shadowing built-in identifier | Error |
| KTN-VAR-005 | Shadowing package variable | Warning |
| KTN-VAR-006 | Unused variable | Error |
| KTN-VAR-007 | Redeclaration in same scope | Error |
| KTN-VAR-008 | Magic number (use named const) | Warning |
| KTN-VAR-009 | Unexported global variable | Info |
| KTN-VAR-010 | Global mutable state | Warning |
| KTN-VAR-011 | Blank identifier misuse | Warning |
| KTN-VAR-012 | Type inference preferred | Info |
| KTN-VAR-013 | Explicit type when obvious | Info |
| KTN-VAR-014 | Variable declared but immediately returned | Info |
| KTN-VAR-015 | Pointer to interface | Warning |
| KTN-VAR-016 | Empty interface{} usage | Info |
| KTN-VAR-017 | nil slice declaration | Info |
| KTN-VAR-018 | Inconsistent error naming | Warning |

## File Structure
```
ktnvar/
├── 001.go ... 018.go       # Rule implementations
├── *_external_test.go      # Tests
├── registry.go             # Analyzers()
└── testdata/src/var001...  # Test fixtures
```

## Common Patterns
```go
// Check naming
if len(ident.Name) < 2 && !utils.IsCommonAbbreviation(ident.Name) {
    pass.Reportf(ident.Pos(), "KTN-VAR-001: ...")
}

// Check shadowing
if utils.IsBuiltinIdentifier(ident.Name) {
    pass.Reportf(ident.Pos(), "KTN-VAR-004: ...")
}
```

## Allowed Short Names
- Loop: `i`, `j`, `k`, `n`
- Error: `err`
- Context: `ctx`
- Receiver: single letter matching type

## Dependencies
Uses `pkg/analyzer/utils` for naming checks and `pkg/analyzer/shared` for scope analysis.
