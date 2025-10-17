package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// PackageOpsAnalyzer vérifie les opérations sur les packages.
	PackageOpsAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnpackageops",
		Doc:  "Vérifie les opérations sur les packages (dot imports)",
		Run:  runPackageOpsAnalyzer,
	}
)

// runPackageOpsAnalyzer exécute l'analyseur d'opérations sur packages.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runPackageOpsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Vérifier les imports
		for _, imp := range file.Imports {
			checkDotImport(pass, imp)
		}
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkDotImport vérifie l'utilisation de dot imports.
//
// Params:
//   - pass: la passe d'analyse
//   - imp: l'import spec
func checkDotImport(pass *analysis.Pass, imp *ast.ImportSpec) {
	if imp.Name == nil {
		// Pas d'alias, c'est ok
		// Retourne
		return
	}

	if imp.Name.Name == "." {
		// C'est un dot import
		reportDotImport(pass, imp)
	}
}

// reportDotImport rapporte une violation KTN-PKG-004.
//
// Params:
//   - pass: la passe d'analyse
//   - imp: l'import spec
func reportDotImport(pass *analysis.Pass, imp *ast.ImportSpec) {
	packagePath := imp.Path.Value

	pass.Reportf(imp.Pos(),
		"[KTN-PKG-004] Dot import de %s.\n"+
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
