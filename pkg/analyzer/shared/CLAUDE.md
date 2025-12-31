# pkg/analyzer/shared/ - Classification Utilities

## Purpose
Higher-level classification and categorization of AST nodes for rule logic.

## Files
| File | Purpose |
|------|---------|
| `classify.go` | Node classification (function type, variable scope) |
| `comments.go` | Comment extraction and association |
| `methods.go` | Method detection, receiver analysis |
| `types.go` | Type classification (struct, interface, etc.) |
| `files.go` | File type detection (test, generated) |
| `testtarget.go` | Test target detection for ktntest rules |

## Key Functions
```go
// Classification
func ClassifyFunction(fn *ast.FuncDecl) FunctionType
func ClassifyVariable(decl *ast.ValueSpec) VariableScope

// Comments
func GetAssociatedComment(node ast.Node, cmap ast.CommentMap) string
func HasDocComment(decl ast.GenDecl) bool

// Methods
func IsMethod(fn *ast.FuncDecl) bool
func GetReceiverName(fn *ast.FuncDecl) string

// Test detection
func IsTestFunction(fn *ast.FuncDecl) bool
func IsBenchmarkFunction(fn *ast.FuncDecl) bool
```

## Usage
```go
import "github.com/kodflow/ktn-linter/pkg/analyzer/shared"

func runRule(pass *analysis.Pass) (any, error) {
    fn := n.(*ast.FuncDecl)
    if shared.IsTestFunction(fn) {
        // Apply test-specific rules
    }
    ftype := shared.ClassifyFunction(fn)
    switch ftype {
    case shared.FunctionTypeHandler:
        // HTTP handler specific checks
    }
}
```

## Difference from utils/
- `utils/`: Low-level primitives (naming, counting)
- `shared/`: High-level classification (what kind of code is this?)
