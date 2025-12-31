# pkg/analyzer/ktn/ktntest/ - Test Rules

## Purpose
Analyze test files for naming conventions, structure, and best practices.

## Rules (11 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-TEST-001 | Test file suffix convention (_internal/_external) | Error |
| KTN-TEST-002 | Orphan test file detection (no source file) | Error |
| KTN-TEST-003 | Function test coverage | Warning |
| KTN-TEST-004 | Table-driven test pattern required | Warning |
| KTN-TEST-005 | t.Skip() forbidden | Warning |
| KTN-TEST-006 | Test file coverage for source files | Warning |
| KTN-TEST-007 | Public function test in external file | Info |
| KTN-TEST-008 | Private function test in internal file | Info |
| KTN-TEST-009 | Package name convention | Error |
| KTN-TEST-010 | Mock file detection | Info |
| KTN-TEST-011 | Mock function detection | Info |

## File Naming Convention (KTN-TEST-001)
```
*_internal_test.go  → package same as source (white-box)
*_external_test.go  → package with _test suffix (black-box)
```

## Test Function Pattern
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test"
    expected := "TEST"

    // Act
    result := MyFunction(input)

    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

## Table-Driven Tests
```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"empty", "", ""},
        {"simple", "hello", "HELLO"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

## File Structure
```
ktntest/
├── 001.go ... 011.go       # Rule implementations
├── *_external_test.go      # Tests (dogfooding!)
├── registry.go             # Analyzers()
└── testdata/src/test001...
```

## Special: Applies Only to *_test.go Files
All ktntest rules skip non-test files automatically.
