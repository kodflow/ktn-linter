# pkg/formatter/ - Output Formatting

## Purpose
Format diagnostic output for terminal, JSON, or SARIF display.

## Structure
```
formatter/
├── formatter.go       # Formatter interface
├── formatter_impl.go  # Text formatter implementation
├── factory.go         # NewFormatterByFormat factory
├── format.go          # OutputFormat type + validation
├── json.go            # JSON formatter
├── json_*.go          # JSON DTOs (report, result, location, etc.)
├── sarif.go           # SARIF formatter
└── *_test.go          # Tests (100% coverage)
```

## Interface
```go
type Formatter interface {
    Format(fset *token.FileSet, diagnostics []analysis.Diagnostic)
}

type OutputFormat string
const (
    FormatText  OutputFormat = "text"
    FormatJSON  OutputFormat = "json"
    FormatSARIF OutputFormat = "sarif"
)
```

## Factory
```go
func NewFormatterByFormat(format OutputFormat, w io.Writer, opts FormatterOptions) Formatter
```

## Output Formats

### Text (default)
```
path/to/file.go
  12:5  error    KTN-FUNC-001  Function too long
  25:1  warning  KTN-VAR-003   Variable name too short
```

### JSON
```json
{"tool": {...}, "summary": {...}, "results": [...]}
```

### SARIF
Standard SARIF 2.1.0 format for IDE integration.

## Color Codes (via severity package)
- ERROR: Red
- WARNING: Yellow
- INFO: Blue
