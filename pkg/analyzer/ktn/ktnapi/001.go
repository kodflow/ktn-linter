// Package ktnapi provides analyzers for API design lint rules.
//
// KTN-API-001: Consumer-side minimal interfaces for external dependencies.
package ktnapi

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeAPI001 is the rule code for this analyzer.
	ruleCodeAPI001 string = "KTN-API-001"
	// defaultMapCapacity initial map capacity for paramInfo maps.
	defaultMapCapacity int = 8
)

var (
	// defaultAllowedPackages contains packages exempt from KTN-API-001.
	//
	// RATIONALE:
	// - time, context: Standard library types universally used as values (time.Time, context.Context)
	// - strings, bytes: Standard library builder types (strings.Builder, bytes.Buffer) - common utilities
	// - go/ast, go/token, go/types: Go toolchain AST packages - fundamental for any Go analyzer
	// - golang.org/x/tools/go/analysis: Official Go analysis framework - all analyzers depend on *analysis.Pass
	// - golang.org/x/tools/go/ast/inspector: Official AST traversal - performance-critical for linters
	// - ktn-linter/pkg/config: Internal config - self-reference would be circular
	//
	// These packages are exempt because:
	// 1. They are part of Go's official toolchain for building analyzers
	// 2. Creating interfaces for them would add complexity without benefit
	// 3. They are not "external dependencies" in the typical sense
	//
	// TODO: Make configurable via .ktn-linter.yaml if needed
	defaultAllowedPackages map[string]bool = map[string]bool{
		"time":                                     true,
		"context":                                  true,
		"strings":                                  true,
		"bytes":                                    true,
		"go/ast":                                   true,
		"go/token":                                 true,
		"go/types":                                 true,
		"golang.org/x/tools/go/analysis":           true,
		"golang.org/x/tools/go/ast/inspector":      true,
		"github.com/kodflow/ktn-linter/pkg/config": true,
	}

	// defaultAllowedTypes contains fully qualified types that are allowed by default.
	defaultAllowedTypes map[string]bool = map[string]bool{
		"time.Time":       true,
		"time.Duration":   true,
		"context.Context": true,
	}

	// Analyzer001 checks that external concrete types used for method calls
	// should be replaced by minimal consumer-side interfaces.
	Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnapi001",
		Doc:      "KTN-API-001: utiliser des interfaces minimales côté consumer pour les dépendances externes",
		Run:      runAPI001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
)

// paramInfo holds information about a parameter being analyzed.
type paramInfo struct {
	ident          *ast.Ident
	paramType      types.Type
	namedType      *types.Named
	methodsCalled  map[string]bool
	hasFieldAccess bool // true if any field is accessed (interface wouldn't help)
}

// runAPI001 exécute l'analyse KTN-API-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runAPI001(pass *analysis.Pass) (any, error) {
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeAPI001) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir les déclarations de fonction
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeAPI001, filename) {
			// Fichier exclu
			return
		}

		// Analyze function
		analyzeFunction(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeFunction analyzes a function declaration for external concrete dependencies.
//
// Params:
//   - pass: analysis pass
//   - funcDecl: function declaration to analyze
func analyzeFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Skip if no parameters
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		// Pas de paramètres
		return
	}

	// Skip if no body (interface method, external declaration)
	if funcDecl.Body == nil {
		// Pas de corps
		return
	}

	// Collect parameter info
	params := collectParams(pass, funcDecl)
	// Vérifier si on a des paramètres
	if len(params) == 0 {
		// Pas de paramètres externes
		return
	}

	// Scan body for method calls on parameters
	scanBodyForMethodCalls(pass, funcDecl.Body, params)

	// Scan body for field accesses on parameters
	scanBodyForFieldAccess(pass, funcDecl.Body, params)

	// Report diagnostics for external concrete types with method calls
	reportDiagnostics(pass, funcDecl, params)
}

// collectParams collects information about external concrete type parameters.
//
// Params:
//   - pass: analysis pass
//   - funcDecl: function declaration
//
// Returns:
//   - map[*ast.Ident]*paramInfo: map of parameter idents to their info
func collectParams(pass *analysis.Pass, funcDecl *ast.FuncDecl) map[*ast.Ident]*paramInfo {
	params := make(map[*ast.Ident]*paramInfo, defaultMapCapacity)

	// Parcourir les champs de paramètres
	for _, field := range funcDecl.Type.Params.List {
		// Get the type
		paramType := pass.TypesInfo.TypeOf(field.Type)
		// Vérifier si le type est nil
		if paramType == nil {
			// Type non résolu
			continue
		}

		// Check if it's an external concrete type
		namedType := getExternalConcreteType(pass, paramType)
		// Vérifier si c'est un type externe concret
		if namedType == nil {
			// Pas un type externe concret
			continue
		}

		// Add each named parameter
		for _, ident := range field.Names {
			// Vérifier si le paramètre est ignoré
			if ident.Name == "_" {
				// Paramètre ignoré
				continue
			}
			params[ident] = &paramInfo{
				ident:         ident,
				paramType:     paramType,
				namedType:     namedType,
				methodsCalled: make(map[string]bool, defaultMapCapacity),
			}
		}
	}

	// Retour de la map
	return params
}

// getExternalConcreteType checks if a type is an external concrete type.
// Returns the named type if it is, nil otherwise.
// Handles type aliases via types.Unalias.
//
// Params:
//   - pass: analysis pass
//   - t: type to check
//
// Returns:
//   - *types.Named: the named type if external concrete, nil otherwise
func getExternalConcreteType(pass *analysis.Pass, t types.Type) *types.Named {
	// Unwrap alias first
	t = types.Unalias(t)

	// Unwrap pointer
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}

	// Must be a named type
	named, ok := t.(*types.Named)
	// Vérifier si c'est un type nommé
	if !ok {
		// Pas un type nommé
		return nil
	}

	// Check if it's an interface (skip interfaces)
	if _, isIface := named.Underlying().(*types.Interface); isIface {
		// C'est une interface
		return nil
	}

	// Check if it's from the same package (skip local types)
	obj := named.Obj()
	// Vérifier si l'objet ou le package est nil
	if obj == nil || obj.Pkg() == nil {
		// Objet ou package nil
		return nil
	}

	// Vérifier si c'est du même package
	if obj.Pkg() == pass.Pkg {
		// Même package
		return nil
	}

	// Check allowlist
	pkgPath := obj.Pkg().Path()
	typeName := pkgPath + "." + obj.Name()

	// Vérifier si le type est dans l'allowlist
	if isAllowedType(pkgPath, typeName) {
		// Type autorisé
		return nil
	}

	// Retour du type nommé
	return named
}

// isAllowedType checks if a type is in the allowlist.
//
// Params:
//   - pkgPath: package path
//   - typeName: fully qualified type name
//
// Returns:
//   - bool: true if allowed
func isAllowedType(pkgPath, typeName string) bool {
	// Check allowed types
	if defaultAllowedTypes[typeName] {
		// Type explicitement autorisé
		return true
	}

	// Check allowed packages
	if defaultAllowedPackages[pkgPath] {
		// Package autorisé
		return true
	}

	// Check if any allowed package is a prefix
	for pkg := range defaultAllowedPackages {
		// Vérifier si le chemin commence par le package
		if strings.HasPrefix(pkgPath, pkg+"/") {
			// Sous-package autorisé
			return true
		}
	}

	// Non autorisé
	return false
}

// scanBodyForMethodCalls scans a function body for method calls on parameters.
//
// Params:
//   - pass: analysis pass
//   - body: function body
//   - params: parameter info map
func scanBodyForMethodCalls(pass *analysis.Pass, body *ast.BlockStmt, params map[*ast.Ident]*paramInfo) {
	// Parcourir le corps de la fonction
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		// Vérifier si c'est un appel
		if !ok {
			// Continuer la traversée
			return true
		}

		// Check if it's a method call (selector expression)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		// Vérifier si c'est un sélecteur
		if !ok {
			// Continuer la traversée
			return true
		}

		// Get the base identifier of the receiver
		baseIdent := getBaseIdent(sel.X)
		// Vérifier si on a trouvé un identifiant
		if baseIdent == nil {
			// Continuer la traversée
			return true
		}

		// Check if the receiver matches any parameter
		for paramIdent, info := range params {
			// Vérifier si le récepteur correspond au paramètre
			if matchesParam(pass, baseIdent, paramIdent) {
				info.methodsCalled[sel.Sel.Name] = true
			}
		}

		// Continuer la traversée
		return true
	})
}

// scanBodyForFieldAccess scans a function body for field accesses on parameters.
// This detects param.Field patterns that are NOT method calls.
// If any field access is found, interface suggestion would be inapplicable.
//
// Params:
//   - pass: analysis pass
//   - body: function body
//   - params: parameter info map
func scanBodyForFieldAccess(pass *analysis.Pass, body *ast.BlockStmt, params map[*ast.Ident]*paramInfo) {
	// First, collect all selector expressions that are method calls
	methodCallSelectors := make(map[*ast.SelectorExpr]bool, defaultMapCapacity)
	// Collecter les sélecteurs d'appels de méthode
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		// Vérifier si c'est un appel
		if !ok {
			// Continuer
			return true
		}
		// Vérifier si c'est un sélecteur
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			methodCallSelectors[sel] = true
		}
		// Continuer
		return true
	})

	// Now find selector expressions that are NOT method calls (field accesses)
	ast.Inspect(body, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		// Vérifier si c'est un sélecteur
		if !ok {
			// Continuer
			return true
		}

		// Skip if this is a method call
		if methodCallSelectors[sel] {
			// C'est un appel de méthode
			return true
		}

		// Get the base identifier
		baseIdent := getBaseIdent(sel.X)
		// Vérifier si on a trouvé un identifiant
		if baseIdent == nil {
			// Continuer
			return true
		}

		// Check if this matches any parameter
		for paramIdent, info := range params {
			// Vérifier si le récepteur correspond au paramètre
			if matchesParam(pass, baseIdent, paramIdent) {
				info.hasFieldAccess = true
			}
		}

		// Continuer
		return true
	})
}

// getMethodSignature retrieves the signature of a method on a named type.
//
// Params:
//   - namedType: the named type
//   - methodName: name of the method
//
// Returns:
//   - string: the method signature or empty string if not found
func getMethodSignature(namedType *types.Named, methodName string) string {
	// Lookup method on the type
	obj, _, _ := types.LookupFieldOrMethod(namedType, true, nil, methodName)
	// Vérifier si la méthode existe
	if obj == nil {
		// Méthode non trouvée
		return ""
	}

	// Cast to func
	fn, ok := obj.(*types.Func)
	// Vérifier si c'est une fonction
	if !ok {
		// Pas une fonction
		return ""
	}

	// Get signature
	sig := fn.Type().(*types.Signature)
	// Retour de la signature formatée
	return formatMethodSignature(methodName, sig)
}

// formatMethodSignature formats a method signature for display.
//
// Params:
//   - name: method name
//   - sig: the signature
//
// Returns:
//   - string: formatted signature line
func formatMethodSignature(name string, sig *types.Signature) string {
	// Build params
	params := formatTupleTypes(sig.Params())
	// Build results
	results := formatTupleTypes(sig.Results())

	// Format the signature
	if results == "" {
		// Pas de retour
		return name + "(" + params + ")"
	}

	// Vérifier si plusieurs résultats
	if sig.Results().Len() > 1 {
		// Multiple results
		return name + "(" + params + ") (" + results + ")"
	}

	// Single result
	return name + "(" + params + ") " + results
}

// formatTupleTypes formats a tuple of types for display.
//
// Params:
//   - tuple: the tuple
//
// Returns:
//   - string: formatted types
func formatTupleTypes(tuple *types.Tuple) string {
	// Vérifier si tuple nil ou vide
	if tuple == nil || tuple.Len() == 0 {
		// Retour vide
		return ""
	}

	tupleLen := tuple.Len()
	parts := make([]string, 0, tupleLen)
	// Parcourir les éléments (range over int - Go 1.22+)
	for i := range tupleLen {
		v := tuple.At(i)
		typeName := types.TypeString(v.Type(), shortQualifier)
		// Vérifier si le paramètre a un nom
		if v.Name() != "" {
			parts = append(parts, v.Name()+" "+typeName)
		} else {
			// Paramètre anonyme
			parts = append(parts, typeName)
		}
	}
	// Retour des parties jointes
	return strings.Join(parts, ", ")
}

// shortQualifier returns package name for qualification.
//
// Params:
//   - pkg: the package
//
// Returns:
//   - string: package name or empty
func shortQualifier(pkg *types.Package) string {
	// Vérifier si le package est nil
	if pkg == nil {
		// Retour vide
		return ""
	}
	// Retour du nom du package
	return pkg.Name()
}

// getBaseIdent extracts the base identifier from an expression.
// Handles parentheses and dereferencing.
//
// Params:
//   - expr: expression to analyze
//
// Returns:
//   - *ast.Ident: base identifier or nil
func getBaseIdent(expr ast.Expr) *ast.Ident {
	// Boucle d'extraction
	for {
		// Switch sur le type d'expression
		switch e := expr.(type) {
		// Identifiant simple
		case *ast.Ident:
			// Retour de l'identifiant
			return e
		// Expression parenthésée
		case *ast.ParenExpr:
			expr = e.X
		// Déréférencement
		case *ast.StarExpr:
			expr = e.X
		// Sélecteur (x.Field.Method())
		case *ast.SelectorExpr:
			// x.Field.Method() - we don't count this
			return nil
		// Type non supporté
		default:
			// Retour nil
			return nil
		}
	}
}

// matchesParam checks if an identifier refers to a parameter.
//
// Params:
//   - pass: analysis pass
//   - ident: identifier to check
//   - paramIdent: parameter identifier
//
// Returns:
//   - bool: true if they refer to the same object
func matchesParam(pass *analysis.Pass, ident, paramIdent *ast.Ident) bool {
	// Use type info to check if they refer to the same object
	identObj := pass.TypesInfo.Uses[ident]
	paramObj := pass.TypesInfo.Defs[paramIdent]

	// Vérifier si les objets sont nil
	if identObj == nil || paramObj == nil {
		// Fallback to name comparison
		return ident.Name == paramIdent.Name
	}

	// Retour de la comparaison
	return identObj == paramObj
}

// reportDiagnostics reports diagnostics for parameters with method calls.
//
// Params:
//   - pass: analysis pass
//   - funcDecl: function declaration
//   - params: parameter info map
func reportDiagnostics(pass *analysis.Pass, funcDecl *ast.FuncDecl, params map[*ast.Ident]*paramInfo) {
	// Parcourir les paramètres
	for _, info := range params {
		// Vérifier si des méthodes ont été appelées
		if len(info.methodsCalled) == 0 {
			// Pas de méthodes appelées
			continue
		}

		// Skip if there's field access - interface wouldn't help
		if info.hasFieldAccess {
			// Accès aux champs - suggestion inapplicable
			continue
		}

		// Generate interface name suggestion
		ifaceName := suggestInterfaceName(info.ident.Name, info.namedType.Obj().Name())

		// Get sorted method names and signatures
		methods := getSortedMethods(info.methodsCalled)
		sigLines := buildInterfaceSignatures(info.namedType, methods)

		// Format the message
		msg, _ := messages.Get(ruleCodeAPI001)
		typeName := formatTypeName(info.paramType)

		pass.Reportf(
			info.ident.Pos(),
			"%s: %s",
			ruleCodeAPI001,
			msg.Format(
				config.Get().Verbose,
				info.ident.Name,    // param name
				typeName,           // type name
				ifaceName,          // suggested interface name
				sigLines,           // signatures for interface
				funcDecl.Name.Name, // function name (verbose only)
				info.ident.Name,    // param name (verbose only)
				ifaceName,          // interface name (verbose only)
			),
		)
	}
}

// buildInterfaceSignatures builds the interface method signatures.
// Returns a multi-line format suitable for copy/paste into code.
//
// Params:
//   - namedType: the named type
//   - methods: sorted method names
//
// Returns:
//   - string: formatted interface body with signatures (multi-line)
func buildInterfaceSignatures(namedType *types.Named, methods []string) string {
	lines := make([]string, 0, len(methods))
	// Parcourir les méthodes
	for _, method := range methods {
		sig := getMethodSignature(namedType, method)
		// Vérifier si la signature existe
		if sig != "" {
			lines = append(lines, "\t"+sig)
		} else {
			// Fallback sans signature
			lines = append(lines, "\t"+method+"(...)")
		}
	}
	// Retour des signatures jointes (multi-ligne)
	return strings.Join(lines, "\n")
}

// suggestInterfaceName suggests an interface name based on param and type names.
//
// Params:
//   - paramName: name of the parameter
//   - typeName: name of the type
//
// Returns:
//   - string: suggested interface name (unexported)
func suggestInterfaceName(paramName, typeName string) string {
	// If param name is meaningful, use it
	if paramName != "" && paramName != "_" {
		// Check if it ends with common suffixes
		suffixes := []string{"repo", "store", "client", "svc", "service", "publisher", "reader", "writer", "clock", "handler"}
		lowerParam := strings.ToLower(paramName)
		// Parcourir les suffixes
		for _, suffix := range suffixes {
			// Vérifier si le nom se termine par le suffixe
			if strings.HasSuffix(lowerParam, suffix) {
				// Retour du nom du paramètre
				return paramName
			}
		}
	}

	// Use type name, make it lowercase (unexported)
	return lowerFirst(typeName)
}

// lowerFirst lowercases the first letter of a string.
//
// Params:
//   - s: string to process
//
// Returns:
//   - string: string with lowercase first letter
func lowerFirst(s string) string {
	// Vérifier si la chaîne est vide
	if s == "" {
		// Retour de la chaîne vide
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	// Retour de la chaîne modifiée
	return string(r)
}

// getSortedMethods returns sorted method names from a map.
//
// Params:
//   - methods: map of method names
//
// Returns:
//   - []string: sorted method names
func getSortedMethods(methods map[string]bool) []string {
	result := make([]string, 0, len(methods))
	// Parcourir les méthodes
	for method := range methods {
		result = append(result, method)
	}
	sort.Strings(result)
	// Retour du résultat trié
	return result
}

// formatTypeName formats a type for display.
// Handles aliases by showing the underlying type.
//
// Params:
//   - t: type to format
//
// Returns:
//   - string: formatted type name
func formatTypeName(t types.Type) string {
	// Handle pointer to alias
	isPtr := false
	// Vérifier si c'est un pointeur
	if ptr, ok := t.(*types.Pointer); ok {
		isPtr = true
		t = ptr.Elem()
	}

	// Unwrap alias
	t = types.Unalias(t)

	// Rebuild pointer if needed
	if isPtr {
		t = types.NewPointer(t)
	}

	// Retour du type formaté
	return types.TypeString(t, shortQualifier)
}
