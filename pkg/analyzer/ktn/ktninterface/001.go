package ktninterface

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Analyzer001 detects unused interface declarations.
var Analyzer001 = &analysis.Analyzer{
	Name:     "ktninterface001",
	Doc:      "KTN-INTERFACE-001: interface non utilisée",
	Run:      runInterface001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runInterface001 analyzes interfaces to detect unused ones.
// Params:
//   - pass: Analysis pass
// Returns: TODO
func runInterface001(pass *analysis.Pass) (any, error) {
	// Collect all interface declarations
	interfaces := make(map[string]*ast.TypeSpec)
	usedInterfaces := make(map[string]bool)
	structNames := make(map[string]bool)

	// First pass: collect all interface and struct declarations
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			genDecl, ok := node.(*ast.GenDecl)
			// Continue if not general declaration
			if !ok {
				return true
			}

   // Verification de la condition
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

				_, isStruct := typeSpec.Type.(*ast.StructType)
				// Store struct names for KTN-STRUCT-002 pattern detection
				if isStruct {
					structNames[typeSpec.Name.Name] = true
				}
			}
			return true
		})
	}

	// Second pass: find interface usages
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
   // Verification de la condition
			switch n := node.(type) {
   // Verification de la condition
			case *ast.FuncDecl:
				// Check parameters
				if n.Type.Params != nil {
					checkFieldList(n.Type.Params, usedInterfaces)
				}
				// Check results
				if n.Type.Results != nil {
					checkFieldList(n.Type.Results, usedInterfaces)
				}
   // Verification de la condition
			case *ast.Field:
				// Check field types
				checkType(n.Type, usedInterfaces)
   // Verification de la condition
			case *ast.InterfaceType:
				// Check embedded interfaces
				if n.Methods != nil {
     // Verification de la condition
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
		// Skip if interface is used
		if usedInterfaces[name] {
			continue
		}

		// Skip if interface follows XXXInterface pattern for struct XXX (KTN-STRUCT-002)
		if isStructInterfacePattern(name, structNames) {
			continue
		}

		// Report unused interface
		pass.Reportf(
			typeSpec.Pos(),
			"KTN-INTERFACE-001: interface non utilisée",
		)
	}

	return nil, nil
}

// isStructInterfacePattern checks if interface follows XXXInterface pattern.
// Params:
//   - interfaceName: Interface name to check
//   - structs: Map of struct names
func isStructInterfacePattern(interfaceName string, structs map[string]bool) bool {
	// Check if interface name ends with "Interface"
	if len(interfaceName) <= 9 || interfaceName[len(interfaceName)-9:] != "Interface" {
		return false
	}

	// Extract potential struct name (remove "Interface" suffix)
	structName := interfaceName[:len(interfaceName)-9]

	// Check if corresponding struct exists
	return structs[structName]
}

// checkFieldList checks field list for interface usage.
// Params:
//   - fields: Field list to check
//   - used: Map to track used interfaces
func checkFieldList(fields *ast.FieldList, used map[string]bool) {
 // Verification de la condition
	for _, field := range fields.List {
		checkType(field.Type, used)
	}
}

// checkType checks type for interface usage.
// Params:
//   - expr: Expression to check
//   - used: Map to track used interfaces
func checkType(expr ast.Expr, used map[string]bool) {
 // Verification de la condition
	switch t := expr.(type) {
 // Verification de la condition
	case *ast.Ident:
		// Mark identifier as used
		used[t.Name] = true
 // Verification de la condition
	case *ast.StarExpr:
		// Check pointer type
		checkType(t.X, used)
 // Verification de la condition
	case *ast.ArrayType:
		// Check array element type
		checkType(t.Elt, used)
 // Verification de la condition
	case *ast.MapType:
		// Check map key and value types
		checkType(t.Key, used)
		checkType(t.Value, used)
 // Verification de la condition
	case *ast.ChanType:
		// Check channel element type
		checkType(t.Value, used)
 // Verification de la condition
	case *ast.SelectorExpr:
		// Check selector expression
		checkType(t.X, used)
	}
}
