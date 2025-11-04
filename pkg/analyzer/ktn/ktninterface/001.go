package ktninterface

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Analyzer001 detects unused interface declarations.
// Params:
//   - N/A
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktninterface001",
	Doc:      "KTN-INTERFACE-001: interface non utilisée",
	Run:      runInterface001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface001 analyzes interfaces to detect unused ones.
// Params:
//   - pass: Analysis pass
func runInterface001(pass *analysis.Pass) (any, error) {
	// Collect all interface declarations
	interfaces := make(map[string]*ast.TypeSpec)
	usedInterfaces := make(map[string]bool)

	// First pass: collect all interface declarations
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			genDecl, ok := node.(*ast.GenDecl)
			// Continue if not general declaration
			if !ok {
				return true
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				// Continue if not type spec
				if !ok {
					continue
				}

				_, isInterface := typeSpec.Type.(*ast.InterfaceType)
				// Store if interface type
				if isInterface {
					interfaces[typeSpec.Name.Name] = typeSpec
				}
			}
			return true
		})
	}

	// Second pass: find interface usages
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				// Check parameters
				if n.Type.Params != nil {
					checkFieldList(n.Type.Params, usedInterfaces)
				}
				// Check results
				if n.Type.Results != nil {
					checkFieldList(n.Type.Results, usedInterfaces)
				}
			case *ast.Field:
				// Check field types
				checkType(n.Type, usedInterfaces)
			case *ast.InterfaceType:
				// Check embedded interfaces
				if n.Methods != nil {
					for _, method := range n.Methods.List {
						// Embedded interface has no function type
						if method.Type != nil {
							checkType(method.Type, usedInterfaces)
						}
					}
				}
			}
			return true
		})
	}

	// Report unused interfaces
	for name, typeSpec := range interfaces {
		// Report if not used
		if !usedInterfaces[name] {
			pass.Reportf(
				typeSpec.Pos(),
				"KTN-INTERFACE-001: interface non utilisée",
			)
		}
	}

	return nil, nil
}

// checkFieldList checks field list for interface usage.
// Params:
//   - fields: Field list to check
//   - used: Map to track used interfaces
func checkFieldList(fields *ast.FieldList, used map[string]bool) {
	for _, field := range fields.List {
		checkType(field.Type, used)
	}
}

// checkType checks type for interface usage.
// Params:
//   - expr: Expression to check
//   - used: Map to track used interfaces
func checkType(expr ast.Expr, used map[string]bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		// Mark identifier as used
		used[t.Name] = true
	case *ast.StarExpr:
		// Check pointer type
		checkType(t.X, used)
	case *ast.ArrayType:
		// Check array element type
		checkType(t.Elt, used)
	case *ast.MapType:
		// Check map key and value types
		checkType(t.Key, used)
		checkType(t.Value, used)
	case *ast.ChanType:
		// Check channel element type
		checkType(t.Value, used)
	case *ast.SelectorExpr:
		// Check selector expression
		checkType(t.X, used)
	}
}
