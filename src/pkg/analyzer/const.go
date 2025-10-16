package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
)

// Analyzers
var (
	// ConstAnalyzer vérifie que les constantes respectent les règles KTN
	ConstAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnconst",
		Doc:  "Vérifie que les constantes sont regroupées, documentées et typées explicitement",
		Run:  runConstAnalyzer,
	}
)

// runConstAnalyzer vérifie que toutes les constantes respectent les règles KTN.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runConstAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkConstDeclarations(pass, file)
	}
	// Retourne nil car l'analyseur rapporte via pass.Reportf
	return nil, nil
}

// checkConstDeclarations parcourt et vérifie toutes les déclarations de constantes.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - file: le fichier AST à analyser
func checkConstDeclarations(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		if genDecl.Lparen == token.NoPos {
			reportUngroupedConst(pass, genDecl)
			continue
		}

		checkConstGroup(pass, genDecl)
	}
}

// reportUngroupedConst signale une constante non groupée.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - genDecl: la déclaration générale contenant la constante non groupée
func reportUngroupedConst(pass *analysis.Pass, genDecl *ast.GenDecl) {
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		for _, name := range valueSpec.Names {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-001] Constante '%s' déclarée individuellement. Regroupez les constantes dans un bloc const ().\nExemple:\n  const (\n      %s %s = ...\n  )",
				name.Name, name.Name, astutil.GetTypeString(valueSpec))
		}
	}
}

// groupContainsIota vérifie si un groupe de constantes utilise iota.
//
// Params:
//   - genDecl: la déclaration de groupe const() à analyser
//
// Returns:
//   - bool: true si au moins une constante du groupe utilise iota
func groupContainsIota(genDecl *ast.GenDecl) bool {
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		if usesIota(valueSpec) {
			// Retourne true car au moins une constante du groupe utilise iota
			return true
		}
	}
	// Retourne false car aucune constante du groupe n'utilise iota
	return false
}

// checkConstGroup vérifie un groupe de constantes.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - genDecl: la déclaration de groupe const() à vérifier
func checkConstGroup(pass *analysis.Pass, genDecl *ast.GenDecl) {
	hasGroupComment := false
	if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
		hasGroupComment = true
	}

	if !hasGroupComment {
		pass.Reportf(genDecl.Pos(),
			"[KTN-CONST-002] Groupe de constantes sans commentaire de groupe.\nAjoutez un commentaire avant le bloc const () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de constantes\n  const (...)")
	}

	// Détecter si le groupe utilise iota
	groupUsesIota := groupContainsIota(genDecl)

	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		isGroupCommentOnly := hasGroupComment &&
			valueSpec.Doc != nil &&
			genDecl.Doc != nil &&
			valueSpec.Doc.Pos() == genDecl.Doc.Pos()
		checkConstSpec(pass, valueSpec, isGroupCommentOnly, groupUsesIota)
	}
}

// checkConstSpec vérifie une spécification de constante individuelle.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - spec: la spécification de constante à vérifier
//   - isFirstWithGroupComment: true si la constante partage le commentaire de groupe
//   - groupUsesIota: true si le groupe de constantes utilise iota
func checkConstSpec(pass *analysis.Pass, spec *ast.ValueSpec, isFirstWithGroupComment bool, groupUsesIota bool) {
	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		// KTN-CONST-003: Vérifier le commentaire individuel
		if !hasIndividualComment(spec, isFirstWithGroupComment) {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-003] Constante '%s' sans commentaire individuel.\nChaque constante doit avoir un commentaire explicatif.\nExemple:\n  // %s décrit son rôle\n  %s %s = ...",
				name.Name, name.Name, name.Name, astutil.GetTypeString(spec))
		}

		// KTN-CONST-004: Vérifier le type explicite (exception pour iota)
		if spec.Type == nil && !groupUsesIota {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-004] Constante '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, int8, uint, float64, etc.\nExemple:\n  %s int = ...",
				name.Name, name.Name)
		}
	}
}

// usesIota vérifie si une constante utilise iota directement.
//
// Params:
//   - spec: la spécification de constante à vérifier
//
// Returns:
//   - bool: true si la constante utilise iota explicitement
func usesIota(spec *ast.ValueSpec) bool {
	// Vérifier si l'expression contient iota
	for _, value := range spec.Values {
		if containsIota(value) {
			// Retourne true car l'expression contient iota
			return true
		}
	}

	// Retourne false car aucune valeur ne contient iota
	return false
}

// containsIota vérifie récursivement si une expression contient iota.
//
// Params:
//   - expr: l'expression à analyser
//
// Returns:
//   - bool: true si iota est trouvé dans l'expression
func containsIota(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Retourne true si l'identifiant est iota, false sinon
		return e.Name == "iota"
	case *ast.BinaryExpr:
		// Retourne true si iota est trouvé dans l'opérande gauche ou droite
		return containsIota(e.X) || containsIota(e.Y)
	case *ast.UnaryExpr:
		// Retourne true si iota est trouvé dans l'opérande
		return containsIota(e.X)
	case *ast.ParenExpr:
		// Retourne true si iota est trouvé dans l'expression entre parenthèses
		return containsIota(e.X)
	case *ast.CallExpr:
		for _, arg := range e.Args {
			if containsIota(arg) {
				// Retourne true car iota est trouvé dans un argument de la fonction
				return true
			}
		}
		// Retourne false car iota n'est trouvé dans aucun argument
		return false
	}
	// Retourne false car le type d'expression ne contient pas iota
	return false
}

// hasIndividualComment vérifie si une constante a un commentaire individuel.
//
// Params:
//   - spec: la spécification de constante à vérifier
//   - isFirstWithGroupComment: true si la constante partage le commentaire de groupe
//
// Returns:
//   - bool: true si un commentaire individuel existe
func hasIndividualComment(spec *ast.ValueSpec, isFirstWithGroupComment bool) bool {
	if spec.Doc != nil && len(spec.Doc.List) > 0 {
		if !isFirstWithGroupComment {
			// Retourne true car un commentaire individuel existe
			return true
		}
	} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
		// Retourne true car un commentaire en ligne existe
		return true
	}
	// Retourne false car aucun commentaire individuel n'a été trouvé
	return false
}
