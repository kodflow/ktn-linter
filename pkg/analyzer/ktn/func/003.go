package ktnfunc

import (
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks that exported functions start with a verb
var Analyzer003 = &analysis.Analyzer{
	Name:     "ktnfunc003",
	Doc:      "KTN-FUNC-003: Les fonctions publiques doivent commencer par un verbe (Get, Set, Create, Update, Delete, etc.)",
	Run:      runFunc003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// Common verb prefixes for function names
var commonVerbs = map[string]bool{
	"Get": true, "Set": true, "Is": true, "Has": true, "Can": true, "Should": true,
	"Create": true, "New": true, "Make": true, "Build": true, "Generate": true,
	"Add": true, "Insert": true, "Append": true, "Push": true, "Put": true,
	"Update": true, "Modify": true, "Change": true, "Edit": true, "Replace": true,
	"Delete": true, "Remove": true, "Clear": true, "Reset": true, "Drop": true,
	"Find": true, "Search": true, "Query": true, "Filter": true, "Select": true,
	"Load": true, "Read": true, "Fetch": true, "Retrieve": true, "Parse": true,
	"Save": true, "Write": true, "Store": true, "Persist": true, "Flush": true,
	"Send": true, "Receive": true, "Publish": true, "Subscribe": true, "Emit": true,
	"Start": true, "Stop": true, "Run": true, "Execute": true, "Perform": true,
	"Open": true, "Close": true, "Connect": true, "Disconnect": true, "Listen": true,
	"Validate": true, "Verify": true, "Check": true, "Test": true, "Ensure": true,
	"Convert": true, "Transform": true, "Format": true, "Encode": true, "Decode": true,
	"Register": true, "Unregister": true, "Unsubscribe": true,
	"Enable": true, "Disable": true, "Activate": true, "Deactivate": true,
	"Handle": true, "Process": true, "Apply": true, "Calculate": true, "Compute": true,
	"Render": true, "Draw": true, "Display": true, "Show": true, "Hide": true,
	"Print": true, "Log": true, "Debug": true, "Trace": true, "Warn": true,
	"Count": true, "Average": true,
	"Sort": true, "Order": true, "Group": true, "Merge": true, "Split": true,
	"Copy": true, "Clone": true, "Duplicate": true, "Move": true, "Swap": true,
	"Compare": true, "Match": true, "Contains": true, "Equals": true, "Diff": true,
	"Wait": true, "Sleep": true, "Pause": true, "Resume": true, "Continue": true,
}

func runFunc003(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		funcName := funcDecl.Name.Name

		// Skip unexported functions
		if !ast.IsExported(funcName) {
   // Retour de la fonction
			return
		}

		// Skip test functions
		if isTestFunction(funcName) {
   // Retour de la fonction
			return
		}

		// Skip main function
		if funcName == "main" {
   // Retour de la fonction
			return
		}

		// Skip init function
		if funcName == "init" {
   // Retour de la fonction
			return
		}

		// Check if function name starts with a verb
		if !startsWithVerb(funcName) {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-003: la fonction publique '%s' doit commencer par un verbe (Get, Set, Create, etc.)",
				funcName,
			)
		}
	})

 // Retour de la fonction
	return nil, nil
}

// startsWithVerb checks if a function name starts with a known verb
func startsWithVerb(name string) bool {
	// Extract the first word (before the first uppercase letter after position 0)
	firstWord := extractFirstWord(name)

	// Check if it's in our verb list
	return commonVerbs[firstWord]
}

// extractFirstWord extracts the first word from a PascalCase/camelCase name
func extractFirstWord(name string) string {
 // Vérification de la condition
	if len(name) == 0 {
  // Retour de la fonction
		return ""
	}

	// Find the first uppercase letter after the start
	for i := 1; i < len(name); i++ {
  // Vérification de la condition
		if unicode.IsUpper(rune(name[i])) {
   // Retour de la fonction
			return name[:i]
		}
	}

	// No uppercase found, return the whole name
	return name
}
