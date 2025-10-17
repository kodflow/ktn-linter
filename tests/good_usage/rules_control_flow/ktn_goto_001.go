package rules_control_flow

import "fmt"

// ✅ GOOD: defer au lieu de goto cleanup
func processWithDefer(data string) error {
	defer doCleanupGood() // ✅ defer automatique

	if data == "" {
		return fmt.Errorf("empty data")
	}

	if err := validateGood(data); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	processGood(data)
	return nil
}

// ✅ GOOD: early return au lieu de goto error
func loadDataCorrectly(filename string) error {
	f := openFileGood(filename)
	if f == nil {
		fmt.Println("Error occurred")
		return fmt.Errorf("failed to open")
	}

	data := readDataGood(f)
	if data == nil {
		fmt.Println("Error occurred")
		return fmt.Errorf("failed to read")
	}

	return nil
}

// ✅ GOOD: return direct au lieu de goto found
func findCorrectly(items []int, target int) bool {
	for _, item := range items {
		if item == target {
			return true // ✅ return direct
		}
	}
	return false
}

// ✅ GOOD: switch/case au lieu de goto labels
func stateMachineCorrectly(state int) {
	switch state {
	case 0:
		handleAGood()
	case 1:
		handleBGood()
	default:
		handleCGood()
	}
	finalizeGood()
}

// ✅ GOOD: for loop au lieu de goto retry
func retryWithLoop(maxRetries int) {
	for retries := 0; retries < maxRetries; retries++ {
		if doWorkGood() {
			return
		}
	}
}

// ✅ GOOD: if/else au lieu de goto nested
func nestedIfCorrectly(a, b int) {
	if a > 0 {
		if b > 0 {
			handlePositiveGood()
		} else {
			handleMixedGood()
		}
	}
}

// ✅ GOOD: fonction helper au lieu de goto
func processSteps() error {
	if err := stepOne(); err != nil {
		return fmt.Errorf("step one failed: %w", err)
	}
	if err := stepTwo(); err != nil {
		return fmt.Errorf("step two failed: %w", err)
	}
	return stepThree()
}

func stepOne() error   { return nil }
func stepTwo() error   { return nil }
func stepThree() error { return nil }

// ✅ GOOD: pattern avec defer pour cleanup multiple
func complexProcessing() error {
	resource1 := acquireResource1()
	defer releaseResource1(resource1)

	resource2 := acquireResource2()
	defer releaseResource2(resource2)

	// Processing...
	return nil
}

// Fonctions helper
func validateGood(s string) error            { return nil }
func processGood(s string)                   {}
func doCleanupGood()                         {}
func openFileGood(name string) interface{}   { return &struct{}{} }
func readDataGood(f interface{}) interface{} { return &struct{}{} }
func handleAGood()                           {}
func handleBGood()                           {}
func handleCGood()                           {}
func finalizeGood()                          {}
func doWorkGood() bool                       { return true }
func handlePositiveGood()                    {}
func handleMixedGood()                       {}
func acquireResource1() interface{}          { return nil }
func releaseResource1(r interface{})         {}
func acquireResource2() interface{}          { return nil }
func releaseResource2(r interface{})         {}
