// Package rules_error_good contient du code conforme à KTN-ERROR-001.
package rules_error_good

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ✅ Code conforme KTN-ERROR-001 : erreurs wrappées avec contexte

// GoodReadFile lit un fichier en wrappant l'erreur.
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - []byte: contenu du fichier
//   - error: erreur de lecture avec contexte
func GoodReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Retourne l'erreur wrappée avec contexte
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	// Retourne les données lues
	return data, nil
}

// GoodProcessData traite des données en wrappant l'erreur.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
//   - error: erreur de traitement avec contexte
func GoodProcessData(data []byte) (string, error) {
	if len(data) == 0 {
		// Retourne une erreur originale (pas besoin de wrapper)
		return "", errors.New("empty data")
	}

	result, err := parseData(data)
	if err != nil {
		// Retourne l'erreur wrappée avec contexte
		return "", fmt.Errorf("failed to parse data: %w", err)
	}

	// Retourne le résultat
	return result, nil
}

// GoodMultipleReturns a plusieurs returns avec wrapping.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - string: résultat
//   - error: erreur avec contexte
func GoodMultipleReturns(id int) (string, error) {
	if id < 0 {
		// Retourne une erreur originale
		return "", errors.New("invalid id")
	}

	data, err := fetchData(id)
	if err != nil {
		// Retourne l'erreur wrappée avec ID dans le contexte
		return "", fmt.Errorf("failed to fetch data for id %d: %w", id, err)
	}

	result, err := processResult(data)
	if err != nil {
		// Retourne l'erreur wrappée
		return "", fmt.Errorf("failed to process result: %w", err)
	}

	// Retourne le résultat
	return result, nil
}

// GoodNestedCalls a des appels imbriqués avec wrapping.
//
// Params:
//   - input: entrée à traiter
//
// Returns:
//   - error: erreur avec contexte
func GoodNestedCalls(input string) error {
	err := validate(input)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("validation failed for input %s: %w", input, err)
	}

	err = process(input)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("processing failed: %w", err)
	}

	err = save(input)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("save failed: %w", err)
	}

	// Retourne nil car succès
	return nil
}

// GoodErrorInLoop retourne erreur avec wrapping dans boucle.
//
// Params:
//   - items: éléments à traiter
//
// Returns:
//   - error: erreur avec contexte
func GoodErrorInLoop(items []string) error {
	for i, item := range items {
		err := processItem(item)
		if err != nil {
			// Retourne l'erreur wrappée avec index et item
			return fmt.Errorf("failed to process item %d (%s): %w", i, item, err)
		}
	}
	// Retourne nil car succès
	return nil
}

// GoodErrorFromInterface retourne erreur depuis interface avec wrapping.
//
// Params:
//   - r: reader
//
// Returns:
//   - []byte: données lues
//   - error: erreur avec contexte
func GoodErrorFromInterface(r io.Reader) ([]byte, error) {
	buf := make([]byte, 100)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		// Retourne l'erreur wrappée
		return nil, fmt.Errorf("failed to read from reader: %w", err)
	}
	// Retourne les données lues
	return buf[:n], nil
}

// GoodChainedErrors a des erreurs chaînées avec wrapping.
//
// Params:
//   - path: chemin
//
// Returns:
//   - error: erreur avec contexte
func GoodChainedErrors(path string) error {
	file, err := os.Open(path)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	if stat.Size() == 0 {
		// Retourne une erreur originale
		return errors.New("empty file")
	}

	// Retourne nil car succès
	return nil
}

// GoodErrorWithAssignment retourne erreur assignée avec wrapping.
//
// Params:
//   - filename: nom du fichier
//
// Returns:
//   - error: erreur avec contexte
func GoodErrorWithAssignment(filename string) error {
	var err error
	err = validateFilename(filename)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("invalid filename %s: %w", filename, err)
	}

	err = checkFileExists(filename)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("file check failed for %s: %w", filename, err)
	}

	// Retourne nil car succès
	return nil
}

// GoodDeferredError retourne erreur depuis defer avec wrapping.
//
// Params:
//   - path: chemin
//
// Returns:
//   - error: erreur avec contexte
func GoodDeferredError(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		// Retourne l'erreur wrappée
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			// Wrappe l'erreur de fermeture
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	// Retourne nil car succès
	return nil
}

// GoodSwitchError retourne erreur dans switch avec wrapping.
//
// Params:
//   - operation: opération à effectuer
//
// Returns:
//   - error: erreur avec contexte
func GoodSwitchError(operation string) error {
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
		// Retourne l'erreur wrappée avec l'opération
		return fmt.Errorf("operation %s failed: %w", operation, err)
	}

	// Retourne nil car succès
	return nil
}

// GoodReturnNil retourne nil explicitement (pas de wrapping nécessaire).
//
// Returns:
//   - error: toujours nil
func GoodReturnNil() error {
	// Retourne nil explicitement, pas de wrapping nécessaire
	return nil
}

// GoodOriginalError crée et retourne une erreur originale.
//
// Params:
//   - value: valeur à vérifier
//
// Returns:
//   - error: erreur originale
func GoodOriginalError(value int) error {
	if value < 0 {
		// Retourne une erreur originale (pas de wrapping car c'est la source)
		return errors.New("negative value not allowed")
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
