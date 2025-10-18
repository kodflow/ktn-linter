// Package rules_error_good contient du code conforme KTN-ERROR-001 dans contextes complexes.
package rules_error_good

import (
	"errors"
	"fmt"
)

// ✅ Erreurs dans goroutines avec wrapping correct

// GoodErrorInGoroutine lance une goroutine qui wrappe les erreurs correctement.
//
// Params:
//   - ch: channel pour communiquer l'erreur
func GoodErrorInGoroutine(ch chan error) {
	go func() {
		err := riskyOperationCtx()
		if err != nil {
			ch <- fmt.Errorf("goroutine failed: %w", err) // ✅ Correct
		}
	}()
}

// GoodMultipleGoroutinesWithErrors lance plusieurs goroutines avec erreurs wrappées.
//
// Params:
//   - ch: channel d'erreurs
//   - count: nombre de goroutines
func GoodMultipleGoroutinesWithErrors(ch chan error, count int) {
	for i := 0; i < count; i++ {
		go func(id int) {
			err := processTaskCtx(id)
			if err != nil {
				ch <- fmt.Errorf("goroutine %d: %w", id, err) // ✅ Correct
			}
		}(i)
	}
}

// ✅ Erreurs dans closures/fonctions anonymes avec wrapping

// GoodErrorInClosure utilise une closure qui wrappe les erreurs.
//
// Returns:
//   - error: erreur wrappée
func GoodErrorInClosure() error {
	processFunc := func() error {
		err := validateCtx("data")
		if err != nil {
			// Early return from function.
			return fmt.Errorf("closure validation: %w", err) // ✅ Correct
		}
		// Early return from function.
		return nil
	}

	// Early return from function.
	return processFunc()
}

// GoodErrorInNestedClosure a des closures imbriquées avec erreurs wrappées.
//
// Returns:
//   - error: erreur wrappée
func GoodErrorInNestedClosure() error {
	outer := func() error {
		inner := func() error {
			err := deepOperationCtx()
			if err != nil {
				// Early return from function.
				return fmt.Errorf("inner closure: %w", err) // ✅ Correct
			}
			// Early return from function.
			return nil
		}
		err := inner()
		if err != nil {
			// Early return from function.
			return fmt.Errorf("outer closure: %w", err) // ✅ Correct
		}
		// Early return from function.
		return nil
	}
	// Early return from function.
	return outer()
}

// ✅ Erreurs dans select statements avec wrapping

// GoodErrorInSelect retourne une erreur depuis select avec wrapping.
//
// Params:
//   - ch1: premier channel
//   - ch2: second channel
//
// Returns:
//   - error: erreur wrappée
func GoodErrorInSelect(ch1, ch2 chan error) error {
	select {
	case err := <-ch1:
		if err != nil {
			// Early return from function.
			return fmt.Errorf("channel 1: %w", err) // ✅ Correct
		}
	case err := <-ch2:
		if err != nil {
			// Early return from function.
			return fmt.Errorf("channel 2: %w", err) // ✅ Correct
		}
	}
	// Early return from function.
	return nil
}

// GoodErrorInSelectWithDefault retourne erreur depuis select avec wrapping.
//
// Params:
//   - errCh: channel d'erreurs
//
// Returns:
//   - error: erreur wrappée
func GoodErrorInSelectWithDefault(errCh chan error) error {
	select {
	case err := <-errCh:
		if err != nil {
			// Early return from function.
			return fmt.Errorf("received error: %w", err) // ✅ Correct
		}
	default:
		// Return error to caller.
		return errors.New("no error received") // ✅ OK car nouvelle erreur
	}
	// Early return from function.
	return nil
}

// ✅ Erreurs avec type assertions/conversions et wrapping

// GoodErrorFromTypeAssertion retourne erreur après type assertion avec wrapping.
//
// Params:
//   - v: valeur à asserter
//
// Returns:
//   - error: erreur wrappée
func GoodErrorFromTypeAssertion(v interface{}) error {
	if errProvider, ok := v.(ErrorProvider); ok {
		err := errProvider.GetError()
		if err != nil {
			// Early return from function.
			return fmt.Errorf("error provider: %w", err) // ✅ Correct
		}
	}
	// Early return from function.
	return nil
}

// GoodErrorFromInterfaceMethod retourne erreur depuis méthode d'interface avec wrapping.
//
// Params:
//   - provider: fournisseur d'erreurs
//
// Returns:
//   - error: erreur wrappée
func GoodErrorFromInterfaceMethod(provider ErrorProvider) error {
	err := provider.GetError()
	if err != nil {
		// Early return from function.
		return fmt.Errorf("interface method: %w", err) // ✅ Correct
	}
	// Early return from function.
	return nil
}

// ✅ Erreurs avec multiple named returns et wrapping

// GoodMultipleNamedReturns a plusieurs retours nommés avec erreurs wrappées.
//
// Returns:
//   - result: résultat
//   - count: compteur
//   - err: erreur wrappée
func GoodMultipleNamedReturns() (result string, count int, err error) {
	err = performOperationCtx()
	if err != nil {
		// Early return from function.
		return "", 0, fmt.Errorf("operation: %w", err) // ✅ Correct
	}

	data, err := fetchDataCtx(42)
	if err != nil {
		// Early return from function.
		return "", 0, fmt.Errorf("fetch data: %w", err) // ✅ Correct
	}

	// Early return from function.
	return string(data), len(data), nil
}

// ✅ Erreurs dans defer avec named returns et wrapping

// GoodDeferWithNamedReturn utilise defer pour modifier erreur nommée avec wrapping.
//
// Returns:
//   - err: erreur wrappée
func GoodDeferWithNamedReturn() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("defer wrapper: %w", err) // ✅ Correct
		}
	}()

	// Early return from function.
	return dangerousOperationCtx()
}

// ✅ Erreurs dans panic recovery avec wrapping

// GoodErrorInRecover capture panic et wrappe erreurs correctement.
//
// Returns:
//   - err: erreur wrappée
func GoodErrorInRecover() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = fmt.Errorf("recovered panic: %w", e) // ✅ Correct
			} else {
				err = fmt.Errorf("recovered panic: %v", r) // ✅ OK car nouveau contexte
			}
		}
	}()

	panickyFunctionCtx()
	// Early return from function.
	return nil
}

// ✅ Erreurs avec fonction variadique et wrapping

// GoodErrorInVariadic traite erreurs dans fonction variadique avec wrapping.
//
// Params:
//   - operations: opérations à effectuer
//
// Returns:
//   - error: erreur wrappée
func GoodErrorInVariadic(operations ...func() error) error {
	for i, op := range operations {
		err := op()
		if err != nil {
			// Early return from function.
			return fmt.Errorf("operation %d: %w", i, err) // ✅ Correct
		}
	}
	// Early return from function.
	return nil
}

// ✅ Erreurs avec higher-order functions et wrapping

// GoodErrorInMap applique fonction qui peut échouer avec wrapping.
//
// Params:
//   - items: éléments à traiter
//   - fn: fonction de transformation
//
// Returns:
//   - []string: résultats
//   - error: erreur wrappée
func GoodErrorInMap(items []string, fn func(string) (string, error)) ([]string, error) {
	results := make([]string, 0, len(items))
	for i, item := range items {
		result, err := fn(item)
		if err != nil {
			// Early return from function.
			return nil, fmt.Errorf("item %d (%s): %w", i, item, err) // ✅ Correct
		}
		results = append(results, result)
	}
	// Early return from function.
	return results, nil
}

// ✅ Erreurs dans init() function (mieux: éviter erreurs dans init)

// globalInitResult describes this variable.
var globalInitResult string

func init() {
	// ✅ Meilleure pratique: gérer erreurs dans fonction d'initialisation séparée
	// appelée depuis main(), pas directement dans init()
	result, err := initializeSystemCtx()
	if err != nil {
		// Dans init(), on ne peut pas retourner d'erreur
		// Donc: log ou panic, mais pas de return err
		panic(fmt.Errorf("init failed: %w", err)) // ✅ OK pour init()
	}
	globalInitResult = result
}

// ✅ Erreurs avec context cancellation et wrapping

// GoodErrorWithContextCancel retourne erreur de context avec wrapping.
//
// Params:
//   - done: channel de completion
//
// Returns:
//   - error: erreur wrappée
func GoodErrorWithContextCancel(done <-chan struct{}) error {
	select {
	case <-done:
		// Return error to caller.
		return errors.New("context cancelled") // ✅ OK car nouvelle erreur
	default:
		err := doWorkCtx()
		if err != nil {
			// Early return from function.
			return fmt.Errorf("work: %w", err) // ✅ Correct
		}
	}
	// Early return from function.
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

// riskyOperationCtx simule une opération risquée.
//
// Returns:
//   - error: erreur
func riskyOperationCtx() error {
	// Return error to caller.
	return errors.New("risky operation failed")
}

// processTaskCtx simule le traitement d'une tâche.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - error: erreur
func processTaskCtx(id int) error {
	// Early return from function.
	return fmt.Errorf("task %d failed", id)
}

// validateCtx simule une validation.
//
// Params:
//   - data: données à valider
//
// Returns:
//   - error: erreur
func validateCtx(data string) error {
	// Return error to caller.
	return errors.New("validation failed")
}

// deepOperationCtx simule une opération profonde.
//
// Returns:
//   - error: erreur
func deepOperationCtx() error {
	// Return error to caller.
	return errors.New("deep operation failed")
}

// performOperationCtx simule une opération.
//
// Returns:
//   - error: erreur
func performOperationCtx() error {
	// Return error to caller.
	return errors.New("operation failed")
}

// fetchDataCtx simule une récupération de données.
//
// Params:
//   - id: identifiant
//
// Returns:
//   - []byte: données
//   - error: erreur
func fetchDataCtx(id int) ([]byte, error) {
	// Return error if operation fails.
	return nil, errors.New("fetch failed")
}

// dangerousOperationCtx simule une opération dangereuse.
//
// Returns:
//   - error: erreur
func dangerousOperationCtx() error {
	// Return error to caller.
	return errors.New("danger")
}

// panickyFunctionCtx simule une fonction qui panique.
func panickyFunctionCtx() {
	panic("something went wrong")
}

// initializeSystemCtx simule l'initialisation du système.
//
// Returns:
//   - string: résultat
//   - error: erreur
func initializeSystemCtx() (string, error) {
	// Early return from function.
	return "initialized", nil
}

// doWorkCtx simule du travail.
//
// Returns:
//   - error: erreur
func doWorkCtx() error {
	// Return error to caller.
	return errors.New("work failed")
}
