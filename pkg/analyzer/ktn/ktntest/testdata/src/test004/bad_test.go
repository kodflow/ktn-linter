package test004_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test004"
)

// TestCheckPositiveOnlyValid teste SEULEMENT le cas valide (PAS BIEN).
// Ne teste pas les cas d'erreur (valeurs négatives, zéro).
//
// Params:
//   - t: contexte de test
func TestCheckPositiveOnlyValid(t *testing.T) {
	// Appel uniquement avec valeur valide
	test004.CheckPositive(10)
	// Manque: test avec valeur négative ou zéro (cas d'erreur)
}

// TestFormatStringJustHappyPath teste seulement le chemin heureux (PAS BIEN).
// Aucune vérification des cas limites ou exceptionnels.
//
// Params:
//   - t: contexte de test
func TestFormatStringJustHappyPath(t *testing.T) {
	const EXPECTED_OUTPUT string = "[hello]"
	// Appel avec donnée valide uniquement
	output := formatAndIgnore("hello")
	// Vérification basique
	if output != EXPECTED_OUTPUT {
		t.Logf("unexpected output")
	}
	// Manque: vérification des cas limites (vide, spéciaux, etc.)
}

// formatAndIgnore formatte une string sans gérer les cas spéciaux.
//
// Params:
//   - s: string à formater
//
// Returns:
//   - string: résultat formaté
func formatAndIgnore(s string) string {
	const EMPTY_RESULT string = ""
	// Appel de la fonction
	output, _ := test004.FormatString(s)
	// Vérification résultat
	if output == EMPTY_RESULT {
		// Retour vide
		return EMPTY_RESULT
	}
	// Retour résultat
	return output
}
