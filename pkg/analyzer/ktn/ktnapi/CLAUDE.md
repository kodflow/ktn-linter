# pkg/analyzer/ktn/ktnapi/ - API Rules

## Purpose
Analyze API-related patterns, particularly struct field access in HTTP handlers.

## Rules (1 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-API-001 | Mixed field and method access on struct | Warning |

## KTN-API-001 Details
Detects inconsistent access patterns where both field access and getter methods
are used on the same struct type within a function.

```go
// BAD - Mixed access
func handler(u User) {
    name := u.Name        // Direct field access
    email := u.GetEmail() // Method access
    // Inconsistent!
}

// GOOD - Consistent access
func handler(u User) {
    name := u.Name
    email := u.Email
}
```

## File Structure
```
ktnapi/
├── 001.go              # KTN-API-001 implementation
├── 001_external_test.go
├── registry.go         # Analyzers()
└── testdata/src/api001/
    ├── good.go
    └── bad.go
```

## Implementation
Tracks accessed fields and called methods per struct type per function.
Reports if both `obj.Field` and `obj.GetField()` patterns are detected.

## Use Cases
- HTTP handlers accessing request/response structs
- Service methods accessing domain entities
- Data transfer objects (DTOs)

## Dependencies
Uses `pkg/analyzer/shared` for method detection and type classification.
