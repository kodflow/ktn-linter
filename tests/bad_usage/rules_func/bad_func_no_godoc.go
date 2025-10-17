package rules_func

// ANTI-PATTERN: Fonctions sans godoc ou godoc incomplet
// Viole KTN-FUNC-002, KTN-FUNC-003, KTN-FUNC-004

// Fonction SANS AUCUN GODOC - GRAVE !
func ProcessWithoutDoc(data string) error {
	return nil
}

// Fonction sans godoc du tout
func AnotherBadFunction(x int, y int) int {
	return x + y
}

// IncompleteDoc manque la section Params et Returns
func IncompleteDoc(name string, age int) (string, error) {
	return name, nil
}

// NoParams cette fonction a des paramètres mais pas de section Params:
//
// Returns:
//   - string: le résultat
func NoParams(input string, count int) string {
	return input
}

// NoReturns cette fonction a des valeurs de retour mais pas de section Returns:
//
// Params:
//   - data: les données
func NoReturns(data []byte) error {
	return nil
}

func totallyUndocumented(a, b, c int) (int, int, error) {
	return a + b, c, nil
}

// privateWithoutDoc fonction privée sans documentation
func privateWithoutDoc(x float64) float64 {
	return x * 2
}
