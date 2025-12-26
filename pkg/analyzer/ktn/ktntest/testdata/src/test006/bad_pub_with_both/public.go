package test008 // want "KTN-TEST-006: le fichier 'public.go' contient UNIQUEMENT des fonctions publiques. Le fichier 'public_internal_test.go' est inutile et doit être supprimé \\(utilisez 'public_external_test.go' pour tester l'API publique\\)"

// GetValue retourne une valeur (fonction publique)
func GetValue() string {
	return "value"
}

// SetValue définit une valeur (fonction publique)
func SetValue(v string) {
	// Set value
}
