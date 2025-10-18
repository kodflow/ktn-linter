package error001

import (
	"errors"
	"fmt"
)

// CORRECT: Erreur wrappée avec contexte
func goodWrappedError() error {
	err := errors.New("base error")
	if err != nil {
		return fmt.Errorf("failed to process: %w", err)
	}
	return nil
}

// CORRECT: Retourne nil (pas d'erreur)
func goodNilReturn() error {
	return nil
}

// CORRECT: Erreur créée directement
func goodNewError() error {
	return errors.New("new error")
}

// CORRECT: fmt.Errorf avec wrapping
func goodContextWrapping(userID string) error {
	err := errors.New("database error")
	if err != nil {
		return fmt.Errorf("failed to fetch user %s: %w", userID, err)
	}
	return nil
}

// BAD: Erreur retournée sans wrapping
func badUnwrappedError() error {
	err := errors.New("base error")
	if err != nil {
		return err // want "KTN-ERROR-001.*sans contexte"
	}
	return nil
}

// BAD: Erreur retournée directement
func badDirectReturn() error {
	err := doSomething()
	return err // want "KTN-ERROR-001.*sans contexte"
}

// BAD: Erreur dans une condition
func badConditionalReturn() error {
	err := doSomething()
	if err != nil {
		return err // want "KTN-ERROR-001.*sans contexte"
	}
	return nil
}

// BAD: Multiple return, erreur non wrappée
func badMultipleReturn() (string, error) {
	err := errors.New("error")
	if err != nil {
		return "", err // want "KTN-ERROR-001.*sans contexte"
	}
	return "success", nil
}

func doSomething() error {
	return errors.New("something went wrong")
}
