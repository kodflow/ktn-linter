package ktnconst

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzer002 checks that constants are grouped together and placed above var declarations
var Analyzer002 = &analysis.Analyzer{
	Name: "ktnconst002",
	Doc:  "KTN-CONST-002: Vérifie que les constantes sont groupées ensemble et placées au-dessus des déclarations var",
	Run:  runConst002,
}

// runConst002 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runConst002(pass *analysis.Pass) (any, error) {
	// Analyze each file independently
	for _, file := range pass.Files {
		tracker := &declTracker{
			constGroups: []declGroup{},
			varGroups:   []declGroup{},
		}

		// Collect const and var declarations
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Vérification de la condition
			if !ok {
				continue
			}

			// Sélection selon la valeur
			switch genDecl.Tok {
			// Traitement
			case token.CONST:
				tracker.constGroups = append(tracker.constGroups, declGroup{
					decl: genDecl,
					pos:  genDecl.Pos(),
				})
			// Traitement
			case token.VAR:
				tracker.varGroups = append(tracker.varGroups, declGroup{
					decl: genDecl,
					pos:  genDecl.Pos(),
				})
			}
		}

		// Check violations
		checkConstGrouping(pass, tracker)
	}

	// Retour de la fonction
	return nil, nil
}

type declTracker struct {
	constGroups []declGroup
	varGroups   []declGroup
}

// runConst002 exécute l'analyse KTN-CONST-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
type declGroup struct {
	decl *ast.GenDecl
	pos  token.Pos
}

// checkConstGrouping description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
func checkConstGrouping(pass *analysis.Pass, tracker *declTracker) {
	// If no const declarations, nothing to check
	if len(tracker.constGroups) == 0 {
		// Retour de la fonction
		return
	}

	// If no var declarations, only check if consts are scattered
	if len(tracker.varGroups) == 0 {
		checkScatteredConsts(pass, tracker.constGroups)
		// Retour de la fonction
		return
	}

	// Find the position of the first var declaration
	firstVarPos := tracker.varGroups[0].pos

	// Separate consts into those before and after first var
	constGroupsBeforeVar := []declGroup{}
	constGroupsAfterVar := []declGroup{}

	// Itération sur les éléments
	for _, constGroup := range tracker.constGroups {
		// Vérification de la condition
		if constGroup.pos < firstVarPos {
			constGroupsBeforeVar = append(constGroupsBeforeVar, constGroup)
			// Cas alternatif
		} else {
			constGroupsAfterVar = append(constGroupsAfterVar, constGroup)
		}
	}

	// Report consts that appear after var
	for _, constGroup := range constGroupsAfterVar {
		pass.Reportf(
			constGroup.pos,
			"KTN-CONST-002: les constantes doivent être groupées et placées au-dessus des déclarations var",
		)
	}

	// Check if consts before vars are scattered
	checkScatteredConsts(pass, constGroupsBeforeVar)
}

// checkScatteredConsts description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
func checkScatteredConsts(pass *analysis.Pass, constGroups []declGroup) {
	// If 0 or 1 const group, they're not scattered
	if len(constGroups) <= 1 {
		// Retour de la fonction
		return
	}

	// Report all const groups except the first as scattered
	for i := 1; i < len(constGroups); i++ {
		pass.Reportf(
			constGroups[i].pos,
			"KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc",
		)
	}
}
