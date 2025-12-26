package test008 // want "KTN-TEST-006: le fichier 'internal.go' contient des fonctions privées. Il doit avoir un fichier 'internal_internal_test.go' \\(white-box\\)"

// validateData valide les données (fonction privée)
func validateData(data string) bool {
	return len(data) > 0
}

// cleanupData nettoie les données (fonction privée)
func cleanupData(data string) string {
	return data + "_clean"
}
