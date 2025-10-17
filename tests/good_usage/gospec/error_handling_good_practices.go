// Package gospec_good_errors démontre les patterns d'erreurs idiomatiques de Go.
// Référence: https://go.dev/doc/effective_go#errors
// Référence: https://go.dev/blog/error-handling-and-go
package gospec_good_errors

import (
	"errors"
	"fmt"
	"os"
)

// ✅ GOOD: Defining sentinel errors
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrTimeout      = errors.New("timeout")
)

// ✅ GOOD: Custom error type with context
type ValidationError struct {
	Field string
	Value string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %s: %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ✅ GOOD: Checking errors immediately
func GoodErrorCheck() error {
	data, err := loadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if err := processData(data); err != nil {
		return fmt.Errorf("failed to process data: %w", err)
	}

	return nil
}

// ✅ GOOD: Using fmt.Errorf with %w for error wrapping
func GoodErrorWrapping() error {
	err := readConfig()
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	return nil
}

// ✅ GOOD: Early return pattern
func GoodEarlyReturn(x int) error {
	if x < 0 {
		return fmt.Errorf("negative value: %d", x)
	}
	if x > 100 {
		return fmt.Errorf("value too large: %d", x)
	}

	// Happy path with minimal nesting
	fmt.Println("processing:", x)
	return nil
}

// ✅ GOOD: Using defer for cleanup
func GoodDeferCleanup() error {
	f, err := openFile("data.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closeFile(f)

	if err := processFile(f); err != nil {
		return fmt.Errorf("failed to process: %w", err)
	}

	return nil
}

// ✅ GOOD: Using errors.Is for sentinel error checking
func GoodErrorIs() error {
	err := operation()
	if errors.Is(err, ErrNotFound) {
		return fmt.Errorf("item not found: %w", err)
	}
	return err
}

// ✅ GOOD: Using errors.As for type assertion
func GoodErrorAs() error {
	err := validate()
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return fmt.Errorf("validation error on field %s: %w",
			validationErr.Field, err)
	}
	return err
}

// ✅ GOOD: Return zero value on error
func GoodReturnZeroValue() (int, error) {
	err := checkCondition()
	if err != nil {
		return 0, fmt.Errorf("condition check failed: %w", err)
	}
	return 42, nil
}

// ✅ GOOD: Multiple error checks with clear flow
func GoodMultipleChecks() error {
	if err := step1(); err != nil {
		return fmt.Errorf("step 1 failed: %w", err)
	}

	if err := step2(); err != nil {
		return fmt.Errorf("step 2 failed: %w", err)
	}

	if err := step3(); err != nil {
		return fmt.Errorf("step 3 failed: %w", err)
	}

	return nil
}

// ✅ GOOD: Error messages are lowercase, no punctuation
func GoodErrorMessage() error {
	return fmt.Errorf("failed to connect to database")
}

// ✅ GOOD: Adding context to errors
func GoodErrorContext(userID int) error {
	err := fetchUser(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user %d: %w", userID, err)
	}
	return nil
}

// ✅ GOOD: Using panic only for unrecoverable errors
func GoodPanicUsage() {
	if !isInitialized() {
		panic("system not initialized - programmer error")
	}
}

// ✅ GOOD: Recovering from panic when appropriate
func GoodRecover() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	riskyOperation()
	return nil
}

// ✅ GOOD: Error handling with cleanup on multiple returns
func GoodMultipleReturnsWithCleanup() (err error) {
	f, err := openFile("config.txt")
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer func() {
		if closeErr := closeFile(f); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close: %w", closeErr)
		}
	}()

	if err := processFile(f); err != nil {
		return fmt.Errorf("failed to process: %w", err)
	}

	return nil
}

// ✅ GOOD: Using named return for defer error wrapping
func GoodNamedReturnForDefer() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("operation failed: %w", err)
		}
	}()

	return operation()
}

// ✅ GOOD: Checking error before using result
func GoodCheckBeforeUse() error {
	result, err := calculate()
	if err != nil {
		return fmt.Errorf("calculation failed: %w", err)
	}

	// Only use result after checking error
	if result > 0 {
		fmt.Println("positive result:", result)
	}

	return nil
}

// ✅ GOOD: Error handling in goroutines
func GoodGoroutineError() error {
	errCh := make(chan error, 1)

	go func() {
		if err := asyncOperation(); err != nil {
			errCh <- fmt.Errorf("async operation failed: %w", err)
			return
		}
		errCh <- nil
	}()

	return <-errCh
}

// ✅ GOOD: Creating descriptive error messages
func GoodDescriptiveError(filename string) error {
	_, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open configuration file %q: %w", filename, err)
	}
	return nil
}

// Helper functions
func loadData() ([]byte, error)     { return nil, nil }
func processData([]byte) error      { return nil }
func readConfig() error             { return nil }
func openFile(string) (int, error)  { return 0, nil }
func closeFile(int) error           { return nil }
func processFile(int) error         { return nil }
func operation() error              { return nil }
func validate() error               { return nil }
func checkCondition() error         { return nil }
func step1() error                  { return nil }
func step2() error                  { return nil }
func step3() error                  { return nil }
func fetchUser(int) error           { return nil }
func isInitialized() bool           { return true }
func riskyOperation()               {}
func calculate() (int, error)       { return 0, nil }
func asyncOperation() error         { return nil }
