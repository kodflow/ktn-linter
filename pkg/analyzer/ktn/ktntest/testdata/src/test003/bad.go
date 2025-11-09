package test003

// BadResource représente une ressource sans tests.
// Les méthodes publiques n'ont PAS de tests correspondants.
type BadResource struct {
	name string
}

// NewBadResource crée une nouvelle instance. // want "KTN-TEST-003: fonction publique 'NewBadResource' n'a pas de test correspondant. Créer un test nommé 'TestNewBadResource'"
//
// Returns:
//   - *BadResource: nouvelle instance
func NewBadResource() *BadResource {
	// Retour de la nouvelle instance
	return &BadResource{}
}

// GetData retourne des données. // want "KTN-TEST-003: fonction publique 'GetData' n'a pas de test correspondant. Créer un test nommé 'TestBadResource_GetData'"
//
// Returns:
//   - string: données
func (r *BadResource) GetData() string {
	// Retour des données
	return "data"
}

// Process traite les données. // want "KTN-TEST-003: fonction publique 'Process' n'a pas de test correspondant. Créer un test nommé 'TestBadResource_Process'"
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
