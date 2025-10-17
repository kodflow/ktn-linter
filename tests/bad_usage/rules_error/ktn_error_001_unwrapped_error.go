// Package rules_error_bad contient du code violant KTN-ERROR-001.
package rules_error_bad

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ❌ Code violant KTN-ERROR-001 : erreurs retournées sans wrapping

// BadReadFile lit un fichier sans wrapper l'erreur.
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - []byte: contenu du fichier
//   - error: erreur de lecture
func BadReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err // Viole KTN-ERROR-001
	}
	// Retourne les données lues
	return data, nil
}

// BadProcessData traite des données sans wrapper l'erreur.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
//   - error: erreur de traitement
func BadProcessData(data []byte) (string, error) {
	if len(data) == 0 {
		// Retourne une erreur
		return "", errors.New("empty data")
	}

	result, err := parseData(data)
	if err != nil {
		return "", err // Viole KTN-ERROR-001
	}

	// Retourne le résultat
	return result, nil
}

// BadMultipleReturns a plusieurs returns sans wrapping.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - string: résultat
//   - error: erreur
func BadMultipleReturns(id int) (string, error) {
	if id < 0 {
		// Retourne une erreur
		return "", errors.New("invalid id")
	}

	data, err := fetchData(id)
	if err != nil {
		return "", err // Viole KTN-ERROR-001
	}

	result, err := processResult(data)
	if err != nil {
		return "", err // Viole KTN-ERROR-001
	}

	// Retourne le résultat
	return result, nil
}

// BadNestedCalls a des appels imbriqués sans wrapping.
//
// Params:
//   - input: entrée à traiter
//
// Returns:
//   - error: erreur de traitement
func BadNestedCalls(input string) error {
	err := validate(input)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	err = process(input)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	err = save(input)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	// Retourne nil car succès
	return nil
}

// BadErrorInLoop retourne erreur sans wrapping dans boucle.
//
// Params:
//   - items: éléments à traiter
//
// Returns:
//   - error: erreur de traitement
func BadErrorInLoop(items []string) error {
	for _, item := range items {
		err := processItem(item)
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	}
	// Retourne nil car succès
	return nil
}

// BadErrorFromInterface retourne erreur depuis interface sans wrapping.
//
// Params:
//   - r: reader
//
// Returns:
//   - []byte: données lues
//   - error: erreur de lecture
func BadErrorFromInterface(r io.Reader) ([]byte, error) {
	buf := make([]byte, 100)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err // Viole KTN-ERROR-001
	}
	// Retourne les données lues
	return buf[:n], nil
}

// BadChainedErrors a des erreurs chaînées sans wrapping.
//
// Params:
//   - path: chemin
//
// Returns:
//   - error: erreur
func BadChainedErrors(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	if stat.Size() == 0 {
		// Retourne une erreur
		return errors.New("empty file")
	}

	// Retourne nil car succès
	return nil
}

// BadErrorWithAssignment retourne erreur assignée sans wrapping.
//
// Params:
//   - filename: nom du fichier
//
// Returns:
//   - error: erreur
func BadErrorWithAssignment(filename string) error {
	var err error
	err = validateFilename(filename)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	err = checkFileExists(filename)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	// Retourne nil car succès
	return nil
}

// BadDeferredError retourne erreur depuis defer sans wrapping.
//
// Params:
//   - path: chemin
//
// Returns:
//   - error: erreur
func BadDeferredError(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err // Viole KTN-ERROR-001
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = closeErr // Viole KTN-ERROR-001
		}
	}()

	// Retourne nil car succès
	return nil
}

// BadSwitchError retourne erreur dans switch sans wrapping.
//
// Params:
//   - operation: opération à effectuer
//
// Returns:
//   - error: erreur
func BadSwitchError(operation string) error {
	var err error
	switch operation {
	case "read":
		err = readOperation()
	case "write":
		err = writeOperation()
	case "delete":
		err = deleteOperation()
	}

	if err != nil {
		return err // Viole KTN-ERROR-001
	}

	// Retourne nil car succès
	return nil
}

// Fonctions helpers pour les tests

// parseData simule un parsing.
//
// Params:
//   - data: données à parser
//
// Returns:
//   - string: résultat
//   - error: erreur
func parseData(data []byte) (string, error) {
	// Retourne le résultat
	return string(data), nil
}

// fetchData simule une récupération de données.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - []byte: données
//   - error: erreur
func fetchData(id int) ([]byte, error) {
	// Retourne les données
	return []byte(fmt.Sprintf("data-%d", id)), nil
}

// processResult simule un traitement.
//
// Params:
//   - data: données
//
// Returns:
//   - string: résultat
//   - error: erreur
func processResult(data []byte) (string, error) {
	// Retourne le résultat
	return string(data), nil
}

// validate simule une validation.
//
// Params:
//   - input: entrée
//
// Returns:
//   - error: erreur
func validate(input string) error {
	// Retourne nil car succès
	return nil
}

// process simule un traitement.
//
// Params:
//   - input: entrée
//
// Returns:
//   - error: erreur
func process(input string) error {
	// Retourne nil car succès
	return nil
}

// save simule une sauvegarde.
//
// Params:
//   - input: entrée
//
// Returns:
//   - error: erreur
func save(input string) error {
	// Retourne nil car succès
	return nil
}

// processItem simule un traitement d'élément.
//
// Params:
//   - item: élément
//
// Returns:
//   - error: erreur
func processItem(item string) error {
	// Retourne nil car succès
	return nil
}

// validateFilename simule une validation.
//
// Params:
//   - filename: nom de fichier
//
// Returns:
//   - error: erreur
func validateFilename(filename string) error {
	// Retourne nil car succès
	return nil
}

// checkFileExists simule une vérification.
//
// Params:
//   - filename: nom de fichier
//
// Returns:
//   - error: erreur
func checkFileExists(filename string) error {
	// Retourne nil car succès
	return nil
}

// readOperation simule une lecture.
//
// Returns:
//   - error: erreur
func readOperation() error {
	// Retourne nil car succès
	return nil
}

// writeOperation simule une écriture.
//
// Returns:
//   - error: erreur
func writeOperation() error {
	// Retourne nil car succès
	return nil
}

// deleteOperation simule une suppression.
//
// Returns:
//   - error: erreur
func deleteOperation() error {
	// Retourne nil car succès
	return nil
}
