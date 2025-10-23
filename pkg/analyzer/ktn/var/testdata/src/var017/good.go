package var017

import (
	"fmt"
	"io"
	"os"
)

// goodNoShadowingInIf utilise la réassignation correcte.
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - error: erreur éventuelle
func goodNoShadowingInIf(path string) error {
	file, err := os.Open(path)
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}

	// Vérification de la condition
	if len(data) > 0 {
		err = goodValidateData(data) // OK: réassignation avec '='
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Retour de la fonction
	return err
}

// goodFmtErrorf utilise la réassignation correcte.
//
// Params:
//   - url: URL de connexion
//
// Returns:
//   - error: erreur éventuelle
func goodFmtErrorf(url string) error {
	conn, err := goodDial(url)
	// Vérification de la condition
	if err != nil {
		err = fmt.Errorf("failed to connect: %w", err) // OK: réassignation
		_ = conn
		// Retour de la fonction
		return err
	}
	// Retour de la fonction
	return nil
}

// goodInFor utilise la réassignation correcte dans une boucle.
//
// Params:
//   - files: liste de fichiers
//
// Returns:
//   - error: erreur éventuelle
func goodInFor(files []string) error {
	var err error
	// Parcours des éléments
	for _, file := range files {
		err = goodProcessFile(file) // OK: réassignation
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}
	// Retour de la fonction
	return err
}

// goodNewVariable déclare une nouvelle variable avec un nom différent.
//
// Returns:
//   - error: erreur éventuelle
func goodNewVariable() error {
	result, err := goodDoSomething()
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}

	// Vérification de la condition
	if result > 0 {
		err2 := goodDoAnotherThing() // OK: nouvelle variable avec nom différent
		// Vérification de la condition
		if err2 != nil {
			// Retour de la fonction
			return err2
		}
	}

	err = goodFinalCheck() // OK: réassignation
	// Retour de la fonction
	return err
}

// goodLocalScopeErr déclare err dans un scope différent (OK).
//
// Returns:
//   - error: erreur éventuelle
func goodLocalScopeErr() error {
	// Vérification de la condition
	if true {
		err := goodDoSomething2() // OK: première déclaration dans ce scope
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Vérification de la condition
	if false {
		err := goodDoAnotherThing() // OK: première déclaration dans ce scope
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Retour de la fonction
	return nil
}

// goodValidateData valide les données.
//
// Params:
//   - data: données à valider
//
// Returns:
//   - error: erreur éventuelle
func goodValidateData(data []byte) error {
	// Retour de la fonction
	return nil
}

// goodDial établit une connexion.
//
// Params:
//   - url: URL de connexion
//
// Returns:
//   - interface{}: connexion
//   - error: erreur éventuelle
func goodDial(url string) (interface{}, error) {
	// Retour de la fonction
	return nil, nil
}

// goodProcessFile traite un fichier.
//
// Params:
//   - file: fichier à traiter
//
// Returns:
//   - error: erreur éventuelle
func goodProcessFile(file string) error {
	// Retour de la fonction
	return nil
}

// goodDoSomething effectue une opération.
//
// Returns:
//   - int: résultat
//   - error: erreur éventuelle
func goodDoSomething() (int, error) {
	// Retour de la fonction
	return 0, nil
}

// goodDoSomething2 effectue une autre opération.
//
// Returns:
//   - error: erreur éventuelle
func goodDoSomething2() error {
	// Retour de la fonction
	return nil
}

// goodDoAnotherThing effectue encore une autre opération.
//
// Returns:
//   - error: erreur éventuelle
func goodDoAnotherThing() error {
	// Retour de la fonction
	return nil
}

// goodFinalCheck effectue une vérification finale.
//
// Returns:
//   - error: erreur éventuelle
func goodFinalCheck() error {
	// Retour de la fonction
	return nil
}
