# pkg/severity/ - Severity Classification

## Purpose
Define severity levels for rules and their visual representation.

## Structure
```
severity/
├── severity.go   # Level enum + color codes
└── *_test.go     # Tests (100% coverage)
```

## Levels
```go
type Level int

const (
    Info    Level = iota  // Blue, informational
    Warning               // Yellow, potential issue
    Error                 // Red, definite problem
)
```

## Color Codes
```go
var Colors = map[Level]string{
    Info:    "\033[34m",  // Blue
    Warning: "\033[33m",  // Yellow
    Error:   "\033[31m",  // Red
}
var Reset = "\033[0m"
```

## Rule Mapping
```go
var RuleSeverity = map[string]Level{
    "KTN-FUNC-001": Error,    // Function too long
    "KTN-VAR-003":  Warning,  // Short variable name
    "KTN-COMMENT-001": Info,  // Missing comment
}
```

## Usage
```go
level := severity.GetLevel("KTN-FUNC-001")
color := severity.Colors[level]
fmt.Printf("%s%s%s\n", color, message, severity.Reset)
```
