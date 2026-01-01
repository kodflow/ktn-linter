# pkg/analyzer/ktn/ktncomment/ - Comment Rules

## Purpose
Analyze comments for format, placement, and documentation standards.

## Rules (7 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-COMMENT-001 | Missing package comment | Warning |
| KTN-COMMENT-002 | Comment not starting with name | Warning |
| KTN-COMMENT-003 | TODO without owner/issue | Info |
| KTN-COMMENT-004 | FIXME without priority | Info |
| KTN-COMMENT-005 | Comment line too long | Info |
| KTN-COMMENT-006 | Trailing comment on statement | Info |
| KTN-COMMENT-007 | Empty comment block | Warning |

## Go Doc Convention
```go
// Package mypackage provides utilities for...
package mypackage

// MyFunction does something important.
// It takes an input and returns an output.
func MyFunction(input string) string { ... }

// MyType represents a thing.
type MyType struct { ... }
```

## Comment Format Rules
- Start with the name being documented
- End with a period
- No blank line between comment and declaration
- Use `//` not `/* */` for doc comments

## File Structure
```
ktncomment/
├── 001.go ... 007.go       # Rule implementations
├── *_external_test.go      # Tests
├── registry.go             # Analyzers()
└── testdata/src/comment001...
```

## TODO/FIXME Format
```go
// TODO(username): Implement caching - issue #123
// FIXME(P1): Memory leak in connection pool
```

## Testdata Example
```go
// bad.go
// does something  // want "KTN-COMMENT-002" (doesn't start with name)
func badExample() {}

// good.go
// GoodExample demonstrates proper documentation.
func GoodExample() {}
```
