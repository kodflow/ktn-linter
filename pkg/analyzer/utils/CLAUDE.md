# pkg/analyzer/utils/ - AST Utilities

## Purpose
Reusable helper functions for AST analysis. 100% test coverage.

## Files
| File | Functions |
|------|-----------|
| `ast.go` | AST node traversal, position helpers |
| `ident.go` | Identifier analysis (exported, naming) |
| `types.go` | Type checking, interface detection |
| `file.go` | File path utilities, test file detection |
| `naming.go` | Naming conventions (camelCase, SCREAMING_SNAKE) |
| `make.go` | Make/new detection for allocations |
| `slice.go` | Slice utilities |

## Key Functions
```go
// Naming
func IsCamelCase(name string) bool
func IsScreamingSnakeCase(name string) bool
func ToScreamingSnake(name string) string

// Types
func IsExported(name string) bool
func IsTestFile(filename string) bool
func IsInternalTestFile(filename string) bool
func IsExternalTestFile(filename string) bool

// AST
func GetFunctionLength(fn *ast.FuncDecl) int
func CountParameters(fn *ast.FuncDecl) int
func GetReceiverType(fn *ast.FuncDecl) string
```

## Usage in Analyzers
```go
import "github.com/kodflow/ktn-linter/pkg/analyzer/utils"

func runMyRule(pass *analysis.Pass) (any, error) {
    if utils.IsTestFile(pass.Fset.File(node.Pos()).Name()) {
        return nil, nil // Skip test files
    }
    if !utils.IsCamelCase(ident.Name) {
        pass.Reportf(ident.Pos(), "...")
    }
}
```

## Testing
All functions have unit tests in `*_test.go` files. Coverage: 100%.
