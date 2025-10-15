package rules_var

import "fmt"

// Variables correctement nommées et utilisées sans shadowing problématique.

var (
	// GlobalValue valeur globale accessible partout.
	GlobalValue int = 100
	// Result résultat du traitement.
	Result string
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
	return localResult
}

// Counter compteur global pour les boucles.
var Counter int = 0

// LoopWithoutShadow traite une boucle sans shadowing.
func LoopWithoutShadow() {
	for i := 0; i < 10; i++ {
		// Variable locale avec nom distinct
		doubledValue := i * 2
		fmt.Println(doubledValue)
	}
}

// Err erreur globale si nécessaire.
var Err error

// GoodErrorHandling gère les erreurs correctement sans shadowing.
//
// Returns:
//   - error: erreur si le traitement échoue
func GoodErrorHandling() error {
	var localErr error

	// Utilisation de variable locale distincte
	data, localErr := readData()
	if localErr != nil {
		return localErr
	}

	// Réutilisation de la même variable locale
	result, localErr := processResult(data)
	if localErr != nil {
		return localErr
	}

	fmt.Println(result)
	return nil
}

func readData() (string, error) {
	return "data", nil
}

func processResult(data string) (string, error) {
	return data, nil
}
