package func008

// goodWithComments a des commentaires sur tous les return.
func goodWithComments(x int) int {
	if x > 0 {
		// Retourne la valeur positive
		return x
	}
	// Retourne zéro par défaut
	return 0
}

// goodWithError a des commentaires sur les returns d'erreur.
func goodWithError(x int) error {
	if x < 0 {
		// Erreur de validation
		return nil
	}
	// Succès
	return nil
}

// badMissingComment n'a pas de commentaires.
func badMissingComment(x int) int {
	if x > 0 {
		return x // want "KTN-FUNC-008.*Return statement sans commentaire"
	}
	return 0 // want "KTN-FUNC-008.*Return statement sans commentaire"
}

// badMissingErrorComment manque commentaires erreur.
func badMissingErrorComment(x int) error {
	if x < 0 {
		return nil // want "KTN-FUNC-008.*Return statement sans commentaire"
	}
	return nil // want "KTN-FUNC-008.*Return statement sans commentaire"
}
