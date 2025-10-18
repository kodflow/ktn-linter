package ktn_package

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie l'absence de dot imports.
//
// KTN-PKG-001: Les dot imports (import . "pkg") sont interdits.
// Ils polluent le namespace et rendent le code confus car on ne sait plus
// d'où viennent les identifiants.
//
// Incorrect:
//
//	import . "fmt"
//	func main() {
//	    Println("hello")  // D'où vient Println?
//	}
//
// Correct:
//
//	import "fmt"
//	func main() {
//	    fmt.Println("hello")  // Clair: c'est fmt.Println
//	}
var Rule001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_PKG_001",
	Doc:  "Vérifie l'absence de dot imports",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-PKG-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// Vérifier les imports
		for _, imp := range file.Imports {
			checkDotImport(pass, imp)
		}
	}

	// Analysis completed successfully.
	return nil, nil
}

// checkDotImport vérifie l'utilisation de dot imports.
//
// Params:
//   - pass: la passe d'analyse
//   - imp: l'import spec
func checkDotImport(pass *analysis.Pass, imp *ast.ImportSpec) {
	if imp.Name == nil {
		// Early return from function.
		return
	}

	if imp.Name.Name == "." {
		reportDotImport(pass, imp)
	}
}

// reportDotImport rapporte une violation KTN-PKG-001.
//
// Params:
//   - pass: la passe d'analyse
//   - imp: l'import spec
func reportDotImport(pass *analysis.Pass, imp *ast.ImportSpec) {
	packagePath := imp.Path.Value

	pass.Reportf(imp.Pos(),
		"[KTN-PKG-001] Dot import de %s.\n"+
			"Les dot imports (import . \"pkg\") polluent le namespace et rendent le code confus.\n"+
			"On ne sait plus d'où viennent les identifiants.\n"+
			"Utilisez un alias normal ou pas d'alias.\n"+
			"Exception: tests uniquement dans certains cas.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - dot import\n"+
			"  import . \"fmt\"\n"+
			"  func main() {\n"+
			"      Println(\"hello\")  // D'où vient Println?\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - import normal\n"+
			"  import \"fmt\"\n"+
			"  func main() {\n"+
			"      fmt.Println(\"hello\")  // Clair: c'est fmt.Println\n"+
			"  }\n"+
			"\n"+
			"  // ✅ ACCEPTABLE - alias si nom long\n"+
			"  import pb \"google.golang.org/protobuf\"\n"+
			"  pb.Marshal(...)",
		packagePath)
}
