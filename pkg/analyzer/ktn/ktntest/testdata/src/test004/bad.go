// Package test004 provides test cases for untested functions.
package test004

// BadResourceInterface defines the public API for BadResource.
type BadResourceInterface interface {
	GetData() string
	Process(input string) string
}

// Ensure BadResource implements BadResourceInterface.
var _ BadResourceInterface = (*BadResource)(nil)

// BadResource représente une ressource sans tests.
// Les méthodes publiques n'ont PAS de tests correspondants.
type BadResource struct {
	name string
}

// NewBadResource crée une nouvelle instance. // want "KTN-TEST-004: fonction publique 'NewBadResource' n'a pas de test correspondant. Créer un test nommé 'TestNewBadResource' dans le fichier 'bad_external_test.go' \\(black-box testing avec package xxx_test\\)"
//
// Returns:
//   - *BadResource: nouvelle instance
func NewBadResource() *BadResource {
	r := &BadResource{}
	// Use private functions
	r.name = r.formatOutput("init")
	_ = validateInput(r.name)
	// Retour de la nouvelle instance
	return r
}

// GetData retourne des données. // want "KTN-TEST-004: fonction publique 'GetData' n'a pas de test correspondant. Créer un test nommé 'TestBadResource_GetData' dans le fichier 'bad_external_test.go' \\(black-box testing avec package xxx_test\\)"
//
// Returns:
//   - string: données
func (r *BadResource) GetData() string {
	// Retour des données
	return "data"
}

// Process traite les données. // want "KTN-TEST-004: fonction publique 'Process' n'a pas de test correspondant. Créer un test nommé 'TestBadResource_Process' dans le fichier 'bad_external_test.go' \\(black-box testing avec package xxx_test\\)"
//
// Params:
//   - input: données à traiter
//
// Returns:
//   - string: résultat
func (r *BadResource) Process(input string) string {
	// Traitement
	return input + "_processed"
}

// validateInput valide l'entrée (fonction privée). // want "KTN-TEST-004: fonction privée 'validateInput' n'a pas de test correspondant. Créer un test nommé 'TestvalidateInput' dans le fichier 'bad_internal_test.go' \\(white-box testing avec package xxx\\)"
//
// Params:
//   - input: données à valider
//
// Returns:
//   - bool: true si valide
func validateInput(input string) bool {
	// Validation
	return len(input) > 0
}

// formatOutput formate la sortie (fonction privée). // want "KTN-TEST-004: fonction privée 'formatOutput' n'a pas de test correspondant. Créer un test nommé 'TestBadResource_formatOutput' dans le fichier 'bad_internal_test.go' \\(white-box testing avec package xxx\\)"
//
// Params:
//   - data: données à formater
//
// Returns:
//   - string: données formatées
func (r *BadResource) formatOutput(data string) string {
	// Formatage
	return "[" + data + "]"
}
