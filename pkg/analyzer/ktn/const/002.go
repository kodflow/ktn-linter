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
			if !ok {
				continue
			}

			switch genDecl.Tok {
			case token.CONST:
				tracker.constGroups = append(tracker.constGroups, declGroup{
					decl: genDecl,
					pos:  genDecl.Pos(),
				})
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

	return nil, nil
}

type declTracker struct {
	constGroups []declGroup
	varGroups   []declGroup
}

type declGroup struct {
	decl *ast.GenDecl
	pos  token.Pos
}

func checkConstGrouping(pass *analysis.Pass, tracker *declTracker) {
	// If no const declarations, nothing to check
	if len(tracker.constGroups) == 0 {
		return
	}

	// If no var declarations, only check if consts are scattered
	if len(tracker.varGroups) == 0 {
		checkScatteredConsts(pass, tracker.constGroups)
		return
	}

	// Find the position of the first var declaration
	firstVarPos := tracker.varGroups[0].pos

	// Separate consts into those before and after first var
	constGroupsBeforeVar := []declGroup{}
	constGroupsAfterVar := []declGroup{}

	for _, constGroup := range tracker.constGroups {
		if constGroup.pos < firstVarPos {
			constGroupsBeforeVar = append(constGroupsBeforeVar, constGroup)
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

func checkScatteredConsts(pass *analysis.Pass, constGroups []declGroup) {
	// If 0 or 1 const group, they're not scattered
	if len(constGroups) <= 1 {
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
