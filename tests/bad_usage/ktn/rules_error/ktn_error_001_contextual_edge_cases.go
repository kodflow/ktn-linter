// Package rules_error_bad contient des violations KTN-ERROR-001 dans contextes complexes.
package rules_error_bad

import (
	"errors"
	"fmt"
)

// ❌ Erreurs dans goroutines sans wrapping

// BadErrorInGoroutine lance une goroutine qui retourne une erreur sans wrapping.
//
// Params:
//   - ch: channel pour communiquer l'erreur
func BadErrorInGoroutine(ch chan error) {
	go func() {
		err := riskyOperation()
		if err != nil {
			ch <- err // Viole KTN-ERROR-001
		}
	}()
}

// BadMultipleGoroutinesWithErrors lance plusieurs goroutines avec erreurs non wrappées.
//
// Params:
//   - ch: channel d'erreurs
//   - count: nombre de goroutines
func BadMultipleGoroutinesWithErrors(ch chan error, count int) {
	for i := 0; i < count; i++ {
		go func(id int) {
			err := processTask(id)
			if err != nil {
				ch <- err // Viole KTN-ERROR-001
			}
		}(i)
	}
}

// ❌ Erreurs dans closures/fonctions anonymes

// BadErrorInClosure utilise une closure qui retourne une erreur sans wrapping.
//
// Returns:
//   - error: erreur
func BadErrorInClosure() error {
	processFunc := func() error {
		err := validate("data")
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
		return nil
	}

	return processFunc()
}

// BadErrorInNestedClosure a des closures imbriquées avec erreurs non wrappées.
//
// Returns:
//   - error: erreur
func BadErrorInNestedClosure() error {
	outer := func() error {
		inner := func() error {
			err := deepOperation()
			if err != nil {
				return err // Viole KTN-ERROR-001
			}
			return nil
		}
		return inner()
	}
	return outer()
}

// ❌ Erreurs dans select statements

// BadErrorInSelect retourne une erreur depuis select sans wrapping.
//
// Params:
//   - ch1: premier channel
//   - ch2: second channel
//
// Returns:
//   - error: erreur
func BadErrorInSelect(ch1, ch2 chan error) error {
	select {
	case err := <-ch1:
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	case err := <-ch2:
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	}
	return nil
}

// BadErrorInSelectWithDefault retourne erreur depuis select avec default.
//
// Params:
//   - errCh: channel d'erreurs
//
// Returns:
//   - error: erreur
func BadErrorInSelectWithDefault(errCh chan error) error {
	select {
	case err := <-errCh:
		return err // Viole KTN-ERROR-001
	default:
		return errors.New("no error received")
	}
}

// ❌ Erreurs avec type assertions/conversions

// BadErrorFromTypeAssertion retourne erreur après type assertion sans wrapping.
//
// Params:
//   - v: valeur à asserter
//
// Returns:
//   - error: erreur
func BadErrorFromTypeAssertion(v interface{}) error {
	if errProvider, ok := v.(ErrorProvider); ok {
		err := errProvider.GetError()
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	}
	return nil
}

// BadErrorFromInterfaceMethod retourne erreur depuis méthode d'interface sans wrapping.
//
// Params:
//   - provider: fournisseur d'erreurs
//
// Returns:
//   - error: erreur
func BadErrorFromInterfaceMethod(provider ErrorProvider) error {
	err := provider.GetError()
	if err != nil {
		return err // Viole KTN-ERROR-001
	}
	return nil
}

// ❌ Erreurs avec multiple named returns

// BadMultipleNamedReturns a plusieurs retours nommés avec erreurs non wrappées.
//
// Returns:
//   - result: résultat
//   - count: compteur
//   - err: erreur
func BadMultipleNamedReturns() (result string, count int, err error) {
	err = performOperation()
	if err != nil {
		return "", 0, err // Viole KTN-ERROR-001
	}

	data, err := fetchData(42)
	if err != nil {
		return "", 0, err // Viole KTN-ERROR-001
	}

	return string(data), len(data), nil
}

// ❌ Erreurs dans defer avec named returns

// BadDeferWithNamedReturn utilise defer pour modifier erreur nommée sans wrapping.
//
// Returns:
//   - err: erreur
func BadDeferWithNamedReturn() (err error) {
	defer func() {
		if err != nil {
			// Modifier l'erreur dans defer sans wrapping
			err = errors.New("defer: " + err.Error()) // Viole KTN-ERROR-001 (pas de %w)
		}
	}()

	return dangerousOperation()
}

// ❌ Erreurs dans panic recovery

// BadErrorInRecover capture panic mais retourne erreur sans wrapping.
//
// Returns:
//   - err: erreur
func BadErrorInRecover() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e // Viole KTN-ERROR-001
			} else {
				err = fmt.Errorf("panic: %v", r) // OK car c'est un nouveau contexte
			}
		}
	}()

	panickyFunction()
	return nil
}

// ❌ Erreurs avec fonction variadique

// BadErrorInVariadic traite erreurs dans fonction variadique sans wrapping.
//
// Params:
//   - operations: opérations à effectuer
//
// Returns:
//   - error: erreur
func BadErrorInVariadic(operations ...func() error) error {
	for _, op := range operations {
		err := op()
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	}
	return nil
}

// ❌ Erreurs avec higher-order functions

// BadErrorInMap applique fonction qui peut échouer sans wrapping.
//
// Params:
//   - items: éléments à traiter
//   - fn: fonction de transformation
//
// Returns:
//   - []string: résultats
//   - error: erreur
func BadErrorInMap(items []string, fn func(string) (string, error)) ([]string, error) {
	results := make([]string, 0, len(items))
	for _, item := range items {
		result, err := fn(item)
		if err != nil {
			return nil, err // Viole KTN-ERROR-001
		}
		results = append(results, result)
	}
	return results, nil
}

// ❌ Erreurs dans init() function

// BadErrorInInit simule une erreur dans init (anti-pattern).
var globalInitError error

func init() {
	err := initializeSystem()
	if err != nil {
		globalInitError = err // Viole KTN-ERROR-001
	}
}

// ❌ Erreurs avec context cancellation

// BadErrorWithContextCancel retourne erreur de context sans wrapping.
//
// Params:
//   - done: channel de completion
//
// Returns:
//   - error: erreur
func BadErrorWithContextCancel(done <-chan struct{}) error {
	select {
	case <-done:
		return errors.New("cancelled") // Pourrait être wrappé si vient d'ailleurs
	default:
		err := doWork()
		if err != nil {
			return err // Viole KTN-ERROR-001
		}
	}
	return nil
}

// Helpers

// ErrorProvider est une interface pour fournir des erreurs.
type ErrorProvider interface {
	// GetError retourne une erreur.
	//
	// Returns:
	//   - error: erreur
	GetError() error
}

// riskyOperation simule une opération risquée.
//
// Returns:
//   - error: erreur
func riskyOperation() error {
	return errors.New("risky operation failed")
}

// processTask simule le traitement d'une tâche.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - error: erreur
func processTask(id int) error {
	return fmt.Errorf("task %d failed", id)
}

// validate simule une validation.
//
// Params:
//   - data: données à valider
//
// Returns:
//   - error: erreur
func validate(data string) error {
	return errors.New("validation failed")
}

// deepOperation simule une opération profonde.
//
// Returns:
//   - error: erreur
func deepOperation() error {
	return errors.New("deep operation failed")
}

// performOperation simule une opération.
//
// Returns:
//   - error: erreur
func performOperation() error {
	return errors.New("operation failed")
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
	return nil, errors.New("fetch failed")
}

// dangerousOperation simule une opération dangereuse.
//
// Returns:
//   - error: erreur
func dangerousOperation() error {
	return errors.New("danger")
}

// panickyFunction simule une fonction qui panique.
func panickyFunction() {
	panic("something went wrong")
}

// initializeSystem simule l'initialisation du système.
//
// Returns:
//   - error: erreur
func initializeSystem() error {
	return errors.New("init failed")
}

// doWork simule du travail.
//
// Returns:
//   - error: erreur
func doWork() error {
	return errors.New("work failed")
}
