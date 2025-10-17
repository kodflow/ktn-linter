// Package gospec_bad_errors montre des patterns d'erreurs non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_errors

import (
	"errors"
	"fmt"
)

// ❌ BAD PRACTICE: Ignoring errors silently
func BadIgnoringErrors() {
	result, _ := riskyOperation() // Erreur ignorée sans raison
	fmt.Println(result)
}

// ❌ BAD PRACTICE: Using panic for normal error handling
func BadUsingPanic(x int) int {
	if x < 0 {
		panic("negative number") // Devrait retourner error
	}
	return x * 2
}

// ❌ BAD PRACTICE: Not wrapping errors to add context
func BadNoErrorWrapping() error {
	err := loadConfig()
	if err != nil {
		return err // Perd le contexte
	}
	return nil
	// Devrait être: return fmt.Errorf("failed to initialize: %w", err)
}

// ❌ BAD PRACTICE: Creating errors with fmt.Errorf without context
func BadGenericError() error {
	return fmt.Errorf("error") // Trop générique, pas de contexte
}

// ❌ BAD PRACTICE: Using errors.New with formatted string
func BadErrorsNewWithFormat(id int) error {
	// Devrait utiliser fmt.Errorf
	return errors.New(fmt.Sprintf("item %d not found", id))
}

// ❌ BAD PRACTICE: Not checking error before using result
func BadNoErrorCheck() int {
	result, err := calculate()
	// Utilise result avant de vérifier err
	if result > 0 {
		return result
	}
	if err != nil {
		return 0
	}
	return result
}

// ❌ BAD PRACTICE: Checking error with err.Error() string comparison
func BadErrorStringComparison() error {
	err := operation()
	if err != nil {
		// Ne devrait pas comparer les strings d'erreur
		if err.Error() == "not found" {
			return fmt.Errorf("item not found")
		}
		return err
	}
	return nil
	// Devrait utiliser errors.Is() ou sentinel errors
}

// ❌ BAD PRACTICE: Returning both value and error
func BadReturningBoth() (int, error) {
	err := validate()
	if err != nil {
		return 42, err // Retourne une valeur même en cas d'erreur
	}
	return 42, nil
	// En cas d'erreur, devrait retourner zero value: return 0, err
}

// ❌ BAD PRACTICE: Verbose error checking
func BadVerboseErrorCheck() error {
	err := step1()
	if err != nil {
		return err
	}

	err = step2()
	if err != nil {
		return err
	}

	err = step3()
	if err != nil {
		return err
	}

	return nil
	// Réutilise la variable err au lieu de déclarer localement
}

// ❌ BAD PRACTICE: Not using defer for cleanup on error
func BadNoDefer() error {
	f := openFile()

	err := processFile(f)
	if err != nil {
		closeFile(f)
		return err
	}

	err = validateFile(f)
	if err != nil {
		closeFile(f) // Duplication
		return err
	}

	closeFile(f)
	return nil
	// Devrait utiliser: defer closeFile(f)
}

// ❌ BAD PRACTICE: Swallowing errors with log only
func BadSwallowingError() {
	err := importantOperation()
	if err != nil {
		fmt.Println("error:", err) // Log seulement, ne propage pas
	}
	// L'erreur devrait être propagée ou gérée
}

// ❌ BAD PRACTICE: Using named return to modify error in defer (unclear)
func BadNamedReturnErrorModification() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("wrapped: %w", err) // Pattern valide mais peut être confus
		}
	}()

	return operation()
	// Ce pattern est techniquement correct mais peut être source de confusion
}

// ❌ BAD PRACTICE: Multiple return points with different error messages for same error
func BadInconsistentErrorMessages(x int) error {
	if x < 0 {
		return fmt.Errorf("invalid input")
	}
	if x == 0 {
		return fmt.Errorf("bad input")
	}
	if x > 100 {
		return fmt.Errorf("wrong input")
	}
	// Messages d'erreur incohérents pour le même type d'erreur
	return nil
}

// ❌ BAD PRACTICE: Not using sentinel errors for common error conditions
func BadNoSentinelError() error {
	// Crée une nouvelle erreur à chaque fois
	return errors.New("not found")
	// Devrait définir: var ErrNotFound = errors.New("not found")
}

// ❌ BAD PRACTICE: Returning nil for pointer types on error
func BadReturningNilPointer() (*Result, error) {
	if !isValid() {
		return nil, fmt.Errorf("invalid") // OK, mais...
	}
	return &Result{}, nil
	// Pattern accepté mais l'appelant doit vérifier error avant d'utiliser result
}

// Helper types and functions
type Result struct{}

func riskyOperation() (int, error) { return 0, nil }
func loadConfig() error            { return nil }
func calculate() (int, error)      { return 0, nil }
func operation() error             { return nil }
func validate() error              { return nil }
func step1() error                 { return nil }
func step2() error                 { return nil }
func step3() error                 { return nil }
func openFile() int                { return 0 }
func closeFile(int)                {}
func processFile(int) error        { return nil }
func validateFile(int) error       { return nil }
func importantOperation() error    { return nil }
func isValid() bool                { return true }
