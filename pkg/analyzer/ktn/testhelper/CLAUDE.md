# pkg/analyzer/ktn/testhelper/ - Test Utilities

## Purpose
Shared test utilities for analyzer tests. NOT a rule category.

## Functions
```go
// TestGoodBad runs analysistest with good.go and bad.go
func TestGoodBad(t *testing.T, a *analysis.Analyzer, dir string, expectedErrors int)

// TestGood runs analysistest expecting 0 errors
func TestGood(t *testing.T, a *analysis.Analyzer, dir string)

// TestBad runs analysistest expecting specific errors
func TestBad(t *testing.T, a *analysis.Analyzer, dir string, expectedErrors int)

// GetTestdataPath returns absolute path to testdata directory
func GetTestdataPath(category string) string
```

## Usage in Tests
```go
package ktnfunc_test

import (
    "testing"
    "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
    "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
)

func TestFunc001(t *testing.T) {
    testhelper.TestGoodBad(t, ktnfunc.Analyzer001, "func001", 3)
}
```

## Testdata Structure
```
testdata/src/<category><NNN>/
├── good.go   # Must produce 0 diagnostics
└── bad.go    # Must produce expected diagnostics with // want comments
```

## Want Comment Format
```go
// bad.go
func badExample() { // want "KTN-FUNC-001"
    // 40 lines of code...
}

// For multiple errors on same line:
func badMultiple() { // want "KTN-FUNC-001" "KTN-FUNC-002"
```

## Avoiding Redeclaration
Functions in bad.go and good.go are in same package. Use prefixes:
```go
// bad.go
func badTooLong() {}

// good.go
func shortEnough() {}
```
