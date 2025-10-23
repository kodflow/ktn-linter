package var017

import (
	"fmt"
	"io"
	"os"
)

// badShadowingInIf démontre le shadowing d'erreur dans un if.
func badShadowingInIf(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		err := validateData(data) // want "KTN-VAR-017: shadowing de la variable 'err' avec ':=' au lieu de '='"
		if err != nil {
			return err
		}
	}

	return err
}

// badShadowingFmtErrorf démontre le shadowing dans fmt.Errorf.
func badShadowingFmtErrorf(url string) error {
	conn, err := dial(url)
	if err != nil {
		err := fmt.Errorf("failed to connect: %w", err) // want "KTN-VAR-017: shadowing de la variable 'err' avec ':=' au lieu de '='"
		_ = conn
		return err
	}
	return nil
}

// badShadowingInFor démontre le shadowing dans une boucle.
func badShadowingInFor(files []string) error {
	var err error
	for _, file := range files {
		err := processFile(file) // want "KTN-VAR-017: shadowing de la variable 'err' avec ':=' au lieu de '='"
		if err != nil {
			return err
		}
	}
	return err
}

// badMultipleShadowing démontre plusieurs shadowings.
func badMultipleShadowing() error {
	result, err := doSomething()
	if err != nil {
		return err
	}

	if result > 0 {
		err := doAnotherThing() // want "KTN-VAR-017: shadowing de la variable 'err' avec ':=' au lieu de '='"
		if err != nil {
			return err
		}
	}

	err := finalCheck() // want "KTN-VAR-017: shadowing de la variable 'err' avec ':=' au lieu de '='"
	return err
}

// badShadowingOtherVar démontre le shadowing d'autres variables.
func badShadowingOtherVar() {
	count := 0
	for i := 0; i < 10; i++ {
		count := i * 2 // want "KTN-VAR-017: shadowing de la variable 'count' avec ':=' au lieu de '='"
		_ = count
	}
	_ = count
}

// Fonctions helper pour les tests.
func validateData(data []byte) error { return nil }
func dial(url string) (interface{}, error) { return nil, nil }
func processFile(file string) error { return nil }
func doSomething() (int, error) { return 0, nil }
func doAnotherThing() error { return nil }
func finalCheck() error { return nil }
