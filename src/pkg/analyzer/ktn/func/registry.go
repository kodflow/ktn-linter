package ktn_func

import "golang.org/x/tools/go/analysis"

// AllRules contient toutes les règles KTN-FUNC.
//
// Cette liste doit être mise à jour lorsqu'une nouvelle règle est ajoutée.
var AllRules = []*analysis.Analyzer{
	Rule001, // KTN-FUNC-001: Nommage MixedCaps
	Rule002, // KTN-FUNC-002: Documentation godoc
	Rule003, // KTN-FUNC-003: Format section Params
	Rule004, // KTN-FUNC-004: Format section Returns
	Rule005, // KTN-FUNC-005: Nombre de paramètres (max 5)
	Rule006, // KTN-FUNC-006: Longueur de fonction (max 35/100)
	Rule007, // KTN-FUNC-007: Complexité cyclomatique (max 10/50)
	Rule008, // KTN-FUNC-008: Commentaires sur return
	Rule009, // KTN-FUNC-010: Profondeur d'imbrication (max 3)
}

// GetRules retourne toutes les règles FUNC.
//
// Returns:
//   - []*analysis.Analyzer: la liste de toutes les règles KTN-FUNC
func GetRules() []*analysis.Analyzer {
	return AllRules
}
