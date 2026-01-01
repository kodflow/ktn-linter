# pkg/analyzer/ktn/ktnvar/ - Variable Rules

## Purpose
Analyze variable declarations for naming, patterns, and modern Go idioms.

## Rules (36 total)

### Core Rules (001-019)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-VAR-001 | Explicit types | Warning |
| KTN-VAR-002 | Declaration order | Warning |
| KTN-VAR-003 | CamelCase naming | Error |
| KTN-VAR-004 | Min length (scope-aware) | Warning |
| KTN-VAR-005 | Max length (30 chars) | Warning |
| KTN-VAR-006 | Shadowing detection | Error |
| KTN-VAR-007 | := vs var (zero-value aware) | Info |
| KTN-VAR-008 | Slice preallocation | Info |
| KTN-VAR-009 | make+append pattern | Info |
| KTN-VAR-010 | Buffer.Grow | Info |
| KTN-VAR-011 | strings.Builder | Info |
| KTN-VAR-012 | Alloc in loops | Warning |
| KTN-VAR-013 | Struct size (64 bytes) | Info |
| KTN-VAR-014 | sync.Pool usage | Info |
| KTN-VAR-015 | string() conversion | Info |
| KTN-VAR-016 | Grouping | Info |
| KTN-VAR-017 | Map preallocation | Info |
| KTN-VAR-018 | Array vs slice (≤64 bytes) | Info |
| KTN-VAR-019 | Mutex copies | Error |

### New Core (020-024)
| Rule | Description | Go Version |
|------|-------------|------------|
| KTN-VAR-020 | Nil slice preferred | - |
| KTN-VAR-021 | Receiver consistency | - |
| KTN-VAR-022 | Pointer to interface | - |
| KTN-VAR-023 | crypto/rand for secrets | - |
| KTN-VAR-024 | any vs interface{} | 1.18+ |

### Go 1.21+ Rules
| Rule | Description |
|------|-------------|
| KTN-VAR-025 | clear() built-in |
| KTN-VAR-026 | min()/max() built-in |
| KTN-VAR-029 | slices.Grow |
| KTN-VAR-030 | slices.Clone |
| KTN-VAR-031 | maps.Clone |
| KTN-VAR-035 | slices.Contains |
| KTN-VAR-036 | slices.Index |

### Go 1.22+ Rules
| Rule | Description |
|------|-------------|
| KTN-VAR-027 | range over integer |
| KTN-VAR-028 | loop var copy obsolete |
| KTN-VAR-033 | cmp.Or |

### Go 1.23+/1.25+ Rules
| Rule | Description | Go Version |
|------|-------------|------------|
| KTN-VAR-037 | maps.Keys/Values | 1.23+ |
| KTN-VAR-034 | WaitGroup.Go | 1.25+ |

## File Structure
```
ktnvar/
├── 001.go ... 037.go       # Rule implementations
├── *_external_test.go      # Tests
├── registry.go             # Analyzers()
└── testdata/src/var001...  # Test fixtures
```

## Idiomatic 1-char Names (VAR-004)
Function scope allows: `i`, `j`, `k`, `n`, `b`, `c`, `f`, `m`, `r`, `s`, `t`, `w`, `_`

## Dependencies
Uses `pkg/analyzer/utils`, `pkg/analyzer/shared`, `pkg/config`, `pkg/messages`
