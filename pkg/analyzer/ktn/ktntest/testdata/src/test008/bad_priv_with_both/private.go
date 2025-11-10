package test008 // want "KTN-TEST-008: le fichier 'private.go' contient UNIQUEMENT des fonctions privées. Le fichier 'private_external_test.go' est inutile et doit être supprimé \\(utilisez 'private_internal_test.go' pour tester les fonctions privées\\)"

// computeValue calcule une valeur (fonction privée)
func computeValue(x int) int {
	return x * 2
}

// processData traite des données (fonction privée)
func processData(data string) string {
	return data + "_processed"
}
