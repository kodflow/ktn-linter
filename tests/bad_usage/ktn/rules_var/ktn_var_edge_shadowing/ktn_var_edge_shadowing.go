package badvarshadow

import "fmt"

// Violations de shadowing et redéclarations

var (
	global_value int = 100 // Snake_case violation
	result       string
)

func processData() {
	// Shadowing de global_value sans commentaire
	global_value := 50

	// Redéclaration avec :=
	result := "local" // Shadow de la variable globale

	if true {
		// Double shadowing
		result := "nested"
		fmt.Println(result)
	}

	fmt.Println(global_value, result)
}

// Shadowing dans boucle
var counter int = 0

func loopShadow() {
	for i := 0; i < 10; i++ {
		// Shadow de counter
		counter := i * 2
		fmt.Println(counter)
	}
}

// Shadowing avec short declaration problématique
var err error

func badErrorHandling() error {
	// err est shadowé au lieu d'être réutilisé
	data, err := readData()
	if err != nil {
		return err
	}

	// err shadowé à nouveau
	result, err := processResult(data)
	fmt.Println(result)
	return err
}

func readData() (string, error) {
	return "data", nil
}

func processResult(data string) (string, error) {
	return data, nil
}
