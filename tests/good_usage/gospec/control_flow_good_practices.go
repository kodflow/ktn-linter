// Package gospec_good_control_flow démontre les patterns de contrôle idiomatiques.
// Référence: https://go.dev/doc/effective_go#control-structures
package gospec_good_control_flow

import "fmt"

// ✅ GOOD: Using range for slices
func GoodRangeSlice() {
	items := []int{1, 2, 3, 4, 5}
	for _, v := range items {
		fmt.Println(v)
	}
}

// ✅ GOOD: Using blank identifier when only value needed
func GoodRangeValue() {
	items := []string{"a", "b", "c"}
	for _, v := range items {
		fmt.Println(v)
	}
}

// ✅ GOOD: Using only index when value not needed
func GoodRangeIndex() {
	items := []int{1, 2, 3}
	for i := range items {
		fmt.Println("index:", i)
	}
}

// ✅ GOOD: Early return instead of nested if
func GoodEarlyReturn(x int) error {
	if x < 0 {
		return fmt.Errorf("negative")
	}
	if x == 0 {
		return fmt.Errorf("zero")
	}
	if x > 100 {
		return fmt.Errorf("too large")
	}

	// Happy path with minimal nesting
	fmt.Println("processing:", x)
	return nil
}

// ✅ GOOD: No unnecessary else after return
func GoodNoElse(x int) int {
	if x > 0 {
		return x * 2
	}
	return 0
}

// ✅ GOOD: Using short statement in if
func GoodIfInit(m map[string]int) {
	if val, ok := m["key"]; ok {
		fmt.Println("found:", val)
	}
}

// ✅ GOOD: Switch with multiple conditions
func GoodSwitch(x int) string {
	switch {
	case x < 0:
		return "negative"
	case x == 0:
		return "zero"
	case x > 0:
		return "positive"
	default:
		return "unknown"
	}
}

// ✅ GOOD: Switch with value
func GoodSwitchValue(s string) {
	switch s {
	case "a", "b", "c":
		fmt.Println("early letters")
	case "x", "y", "z":
		fmt.Println("late letters")
	default:
		fmt.Println("other")
	}
}

// ✅ GOOD: Type switch with assignment
func GoodTypeSwitch(v interface{}) {
	switch x := v.(type) {
	case int:
		fmt.Println("int:", x+1)
	case string:
		fmt.Println("string:", x+"!")
	case bool:
		fmt.Println("bool:", !x)
	default:
		fmt.Println("unknown type")
	}
}

// ✅ GOOD: For loop with clear structure
func GoodForLoop() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println(i)
	}
}

// ✅ GOOD: Range over map
func GoodRangeMap() {
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("%s: %d\n", k, v)
	}
}

// ✅ GOOD: Range over channel
func GoodRangeChannel() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
}

// ✅ GOOD: Using labels for nested loop break
func GoodLabeledBreak() {
	found := false
Outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i*j > 10 {
				found = true
				break Outer
			}
		}
	}
	_ = found
}

// ✅ GOOD: Select with timeout
func GoodSelectTimeout(ch chan int) {
	select {
	case v := <-ch:
		fmt.Println("received:", v)
	case <-timeoutChan():
		fmt.Println("timeout")
	}
}

// ✅ GOOD: Select with default (non-blocking)
func GoodSelectDefault(ch chan int) {
	select {
	case v := <-ch:
		fmt.Println("received:", v)
	default:
		fmt.Println("nothing available")
	}
}

// ✅ GOOD: Defer for cleanup in order
func GoodDeferOrder() {
	defer fmt.Println("first")
	defer fmt.Println("second")
	defer fmt.Println("third")
	// Prints: third, second, first
}

// ✅ GOOD: Using fallthrough when needed
func GoodFallthrough(x int) {
	switch x {
	case 1:
		fmt.Println("one")
		fallthrough
	case 2:
		fmt.Println("one or two")
	}
}

// ✅ GOOD: Simple condition without unnecessary complexity
func GoodSimpleCondition(enabled bool) {
	if enabled {
		activate()
	}
}

// ✅ GOOD: Using continue to reduce nesting
func GoodContinue(items []int) {
	for _, v := range items {
		if v < 0 {
			continue
		}
		if v > 100 {
			continue
		}
		// Process valid items
		fmt.Println("processing:", v)
	}
}

// ✅ GOOD: For with condition only
func GoodForCondition(n int) {
	for n > 0 {
		fmt.Println(n)
		n--
	}
}

// ✅ GOOD: Empty switch case for documentation
func GoodEmptyCase(x int) {
	switch x {
	case 0:
		// Zero is handled elsewhere
	case 1, 2, 3:
		process(x)
	default:
		processDefault(x)
	}
}

// Helper functions
func timeoutChan() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func activate()           {}
func process(int)         {}
func processDefault(int)  {}
