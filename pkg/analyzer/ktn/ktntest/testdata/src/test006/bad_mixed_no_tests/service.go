package test008 // want "KTN-TEST-006: le fichier 'service.go' contient des fonctions publiques ET privées. Il doit avoir DEUX fichiers de test : 'service_internal_test.go' \\(white-box\\) ET 'service_external_test.go' \\(black-box\\)"

// Execute exécute le service (fonction publique)
func Execute(input string) string {
	return process(input)
}

// process traite les données (fonction privée)
func process(data string) string {
	return data + "_processed"
}
