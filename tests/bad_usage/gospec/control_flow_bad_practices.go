// Package gospec_bad_control_flow montre des patterns de contrôle non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_control_flow

import "fmt"

// ❌ BAD PRACTICE: Not using range when iterating over slice/array
func BadNoRange() {
	slice := []int{1, 2, 3, 4, 5}

	// Style C, pas idiomatique en Go
	for i := 0; i < len(slice); i++ {
		fmt.Println(slice[i])
	}

	// Devrait utiliser range
}

// ❌ BAD PRACTICE: Verbose range loop when only index/value needed
func BadVerboseRange() {
	slice := []int{1, 2, 3}

	// N'utilise pas l'index mais le déclare quand même
	for i, v := range slice {
		_ = i // Inutilisé
		fmt.Println(v)
	}
	// Devrait être: for _, v := range slice

	// N'utilise pas la valeur
	for i, v := range slice {
		_ = v // Inutilisé
		fmt.Println(i)
	}
	// Devrait être: for i := range slice
}

// ❌ BAD PRACTICE: Unnecessary else after return
func BadUnnecessaryElse(x int) int {
	if x > 0 {
		return x
	} else { // else inutile
		return -x
	}
	// Le else est redondant car if branch return
}

// ❌ BAD PRACTICE: Nested if instead of early return
func BadNestedIf(x int) error {
	if x > 0 {
		if x < 100 {
			if x%2 == 0 {
				fmt.Println("valid")
				return nil
			} else {
				return fmt.Errorf("odd")
			}
		} else {
			return fmt.Errorf("too large")
		}
	} else {
		return fmt.Errorf("negative")
	}
	// Devrait utiliser early returns et réduire nesting
}

// ❌ BAD PRACTICE: Not using early return for error handling
func BadNoEarlyReturn() error {
	err := step1()
	if err == nil {
		err = step2()
		if err == nil {
			err = step3()
			if err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
	// Devrait utiliser early returns
}

// ❌ BAD PRACTICE: Using break with label when unnecessary
func BadUnnecessaryLabel() {
OuterLoop:
	for i := 0; i < 10; i++ {
		if i == 5 {
			break OuterLoop // Label inutile ici
		}
	}
}

// ❌ BAD PRACTICE: Empty switch/select with default only
func BadEmptySwitch(x int) {
	switch x {
	default:
		fmt.Println("default")
	}
	// Inutilement complexe, devrait être simple if ou direct call
}

// ❌ BAD PRACTICE: Using switch when simple if would be clearer
func BadSwitchForSimpleIf(x int) {
	switch {
	case x > 0:
		fmt.Println("positive")
	}
	// Un simple if serait plus clair
}

// ❌ BAD PRACTICE: Infinite loop without clear exit strategy
func BadInfiniteLoop() {
	count := 0
	for {
		count++
		if count > 100 {
			break
		}
	}
	// Devrait être: for count := 0; count <= 100; count++
}

// ❌ BAD PRACTICE: Not using := in if statement initialization
func BadNoIfInit(m map[string]int) {
	var val int
	var ok bool
	val, ok = m["key"] // Devrait être dans le if
	if ok {
		fmt.Println(val)
	}
	// Devrait être: if val, ok := m["key"]; ok { ... }
}

// ❌ BAD PRACTICE: Verbose type switch
func BadVerboseTypeSwitch(v interface{}) {
	switch v.(type) {
	case int:
		i := v.(int) // Réassignation inutile
		fmt.Println(i)
	case string:
		s := v.(string) // Réassignation inutile
		fmt.Println(s)
	}
	// Devrait utiliser: switch x := v.(type)
}

// ❌ BAD PRACTICE: Using goto when structured control flow is clearer
func BadGotoUsage(x int) {
	if x < 0 {
		goto Error
	}
	if x > 100 {
		goto Error
	}
	fmt.Println("valid")
	return
Error:
	fmt.Println("error")
	// Devrait utiliser early returns
}

// Helper functions
func step1() error { return nil }
func step2() error { return nil }
func step3() error { return nil }
