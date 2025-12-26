package test008 // want "KTN-TEST-006: le fichier 'manager.go' contient des fonctions publiques. Il doit avoir un fichier 'manager_external_test.go' \\(black-box\\)"

// Start démarre le manager (fonction publique)
func Start() string {
	return initialize()
}

// initialize initialise le manager (fonction privée)
func initialize() string {
	return "initialized"
}
