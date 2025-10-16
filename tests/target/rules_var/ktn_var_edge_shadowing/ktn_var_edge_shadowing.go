package rules_var

import "fmt"

// Variables correctement nommées et utilisées sans shadowing problématique.

// Constantes du package.
const (
	// GlobalValue valeur globale accessible partout.
	GlobalValue int = 100
	// Counter compteur global pour les boucles.
	Counter int = 0
)

// Variables du package.
var (
	// Result résultat du traitement.
	Result string
	// Err erreur globale si nécessaire.
	Err error
)

// ProcessData traite les données sans shadowing.
//
// Returns:
//   - string: le résultat du traitement
func ProcessData() string {
	// Utilisation de noms différents pour éviter le shadowing
	localValue := 50

	// Variable locale avec nom distinct
	localResult := "local"

	if true {
		// Variable de scope limité avec nom explicite
		nestedResult := "nested"
		fmt.Println(nestedResult)
	}

	fmt.Println(localValue, localResult)
	// Retourne le résultat local
	return localResult
}

// LoopWithoutShadow traite une boucle sans shadowing.
func LoopWithoutShadow() {
	for i := 0; i < 10; i++ {
		// Variable locale avec nom distinct
		doubledValue := i * 2
		fmt.Println(doubledValue)
	}
}

// GoodErrorHandling gère les erreurs correctement sans shadowing.
//
// Returns:
//   - error: erreur si le traitement échoue
func GoodErrorHandling() error {
	var localErr error

	// Utilisation de variable locale distincte
	data, localErr := readData()
	if localErr != nil {
		// Retourne l'erreur wrappée avec contexte
		return fmt.Errorf("failed to read data: %w", localErr)
	}

	// Réutilisation de la même variable locale
	result, localErr := processResult(data)
	if localErr != nil {
		// Retourne l'erreur wrappée avec contexte
		return fmt.Errorf("failed to process result: %w", localErr)
	}

	fmt.Println(result)
	// Retourne nil car le traitement est terminé avec succès
	return nil
}

// readData lit les données.
//
// Returns:
//   - string: les données lues
//   - error: erreur si la lecture échoue
func readData() (string, error) {
	// Retourne les données lues et nil pour l'erreur
	return "data", nil
}

// processResult traite le résultat.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
//   - error: erreur si le traitement échoue
func processResult(data string) (string, error) {
	// Retourne les données traitées et nil pour l'erreur
	return data, nil
}
