package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
)

// ConstAnalyzer vérifie que les constantes respectent les règles KTN
var ConstAnalyzer = &analysis.Analyzer{
	Name: "ktnconst",
	Doc:  "Vérifie que les constantes sont regroupées, documentées et typées explicitement",
	Run:  runConstAnalyzer,
}

// runConstAnalyzer vérifie que toutes les constantes respectent les règles KTN
// Retourne nil, nil car aucun résultat n'est nécessaire pour cet analyseur
func runConstAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkConstDeclarations(pass, file)
	}
	return nil, nil
}

// checkConstDeclarations parcourt et vérifie toutes les déclarations de constantes
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

// reportUngroupedConst signale une constante non groupée
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

// checkConstGroup vérifie un groupe de constantes
func checkConstGroup(pass *analysis.Pass, genDecl *ast.GenDecl) {
	hasGroupComment := false
	if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
		hasGroupComment = true
	}

	if !hasGroupComment {
		pass.Reportf(genDecl.Pos(),
			"[KTN-CONST-002] Groupe de constantes sans commentaire de groupe.\nAjoutez un commentaire avant le bloc const () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de constantes\n  const (...)")
	}

	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		isGroupCommentOnly := hasGroupComment &&
			valueSpec.Doc != nil &&
			genDecl.Doc != nil &&
			valueSpec.Doc.Pos() == genDecl.Doc.Pos()
		checkConstSpec(pass, valueSpec, isGroupCommentOnly)
	}
}

// checkConstSpec vérifie une spécification de constante individuelle
// Les paramètres pass, spec et isFirstWithGroupComment contrôlent la validation
// Le paramètre isFirstWithGroupComment indique si la constante partage le commentaire de groupe
func checkConstSpec(pass *analysis.Pass, spec *ast.ValueSpec, isFirstWithGroupComment bool) {
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

		// KTN-CONST-004: Vérifier le type explicite
		if spec.Type == nil {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-004] Constante '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, int8, uint, float64, etc.\nExemple:\n  %s int = ...",
				name.Name, name.Name)
		}
	}
}

// hasIndividualComment vérifie si une constante a un commentaire individuel
func hasIndividualComment(spec *ast.ValueSpec, isFirstWithGroupComment bool) bool {
	if spec.Doc != nil && len(spec.Doc.List) > 0 {
		if !isFirstWithGroupComment {
			return true
		}
	} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
		return true
	}
	return false
}
