package var017

import (
	"fmt"
	"io"
	"os"
)

const (
	// LOOP_MAX_ITERATIONS est le nombre maximum d'itérations
	LOOP_MAX_ITERATIONS int = 10
	// MULTIPLIER_VALUE est le multiplicateur utilisé
	MULTIPLIER_VALUE int = 2
)

// badShadowingInIf démontre le shadowing d'erreur dans un if.
//
// Params:
//   - path: chemin du fichier à ouvrir
//
// Returns:
//   - error: erreur éventuelle
func badShadowingInIf(path string) error {
	file, err := os.Open(path)
	// Vérification d'erreur
	if err != nil {
		// Retour avec erreur
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	// Vérification d'erreur
	if err != nil {
		// Retour avec erreur
		return err
	}

	// Vérification de la longueur des données
	if len(data) > 0 {
		err := validateData(data)
		// Vérification d'erreur
		if err != nil {
			// Retour avec erreur
			return err
		}
	}

	// Retour avec dernière erreur
	return err
}

// badShadowingFmtErrorf démontre le shadowing dans fmt.Errorf.
//
// Params:
//   - url: URL de connexion
//
// Returns:
//   - error: erreur éventuelle
func badShadowingFmtErrorf(url string) error {
	conn, err := dial(url)
	// Vérification d'erreur
	if err != nil {
		err := fmt.Errorf("failed to connect: %w", err) // SHADOWING: err redéclaré
		_ = conn
		// Retour avec erreur wrappée
		return err
	}
	// Retour sans erreur
	return nil
}

// badShadowingInFor démontre le shadowing dans une boucle.
//
// Params:
//   - files: liste des fichiers à traiter
//
// Returns:
//   - error: erreur éventuelle
func badShadowingInFor(files []string) error {
	var err error
	// Traitement de chaque fichier
	for _, file := range files {
		err := processFile(file)
		// Vérification d'erreur
		if err != nil {
			// Retour avec erreur
			return err
		}
	}
	// Retour avec dernière erreur
	return err
}

// badMultipleShadowing démontre plusieurs shadowings.
//
// Returns:
//   - error: erreur éventuelle
func badMultipleShadowing() error {
	result, err := doSomething()
	// Vérification d'erreur
	if err != nil {
		// Retour avec erreur
		return err
	}

	// Vérification du résultat
	if result > 0 {
		err := doAnotherThing()
		// Vérification d'erreur
		if err != nil {
			// Retour avec erreur
			return err
		}
	}

	err = finalCheck()
	// Retour avec dernière erreur
	return err
}

// badShadowingOtherVar démontre le shadowing d'autres variables.
func badShadowingOtherVar() {
	count := 0
	// Boucle sur les itérations
	for i := 0; i < LOOP_MAX_ITERATIONS; i++ {
		count := i * MULTIPLIER_VALUE
		_ = count
	}
	_ = count
}

// validateData valide les données.
//
// Params:
//   - data: données à valider
//
// Returns:
//   - error: erreur éventuelle
func validateData(data []byte) error {
	// Retour sans erreur
	return nil
}

// dial établit une connexion.
//
// Params:
//   - url: URL de connexion
//
// Returns:
//   - interface{}: connexion établie
//   - error: erreur éventuelle
func dial(url string) (interface{}, error) {
	// Retour sans erreur
	return nil, nil
}

// processFile traite un fichier.
//
// Params:
//   - file: chemin du fichier
//
// Returns:
//   - error: erreur éventuelle
func processFile(file string) error {
	// Retour sans erreur
	return nil
}

// doSomething effectue une opération.
//
// Returns:
//   - int: résultat de l'opération
//   - error: erreur éventuelle
func doSomething() (int, error) {
	// Retour avec résultat
	return 0, nil
}

// doAnotherThing effectue une autre opération.
//
// Returns:
//   - error: erreur éventuelle
func doAnotherThing() error {
	// Retour sans erreur
	return nil
}

// finalCheck effectue une vérification finale.
//
// Returns:
//   - error: erreur éventuelle
func finalCheck() error {
	// Retour sans erreur
	return nil
}
