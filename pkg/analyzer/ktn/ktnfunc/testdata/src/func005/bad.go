// Bad examples for the func005 test case.
package func005

// Constants used in complexity tests
const (
	BAD_THRESHOLD_TWO     int = 2
	BAD_THRESHOLD_THREE   int = 3
	BAD_THRESHOLD_FIVE    int = 5
	BAD_THRESHOLD_TEN     int = 10
	BAD_THRESHOLD_TWENTY  int = 20
	BAD_THRESHOLD_FIFTY   int = 50
	BAD_THRESHOLD_HUNDRED int = 100
	BAD_MIN_TEN           int = -10
	BAD_MIN_TWENTY        int = -20
	BAD_CASE_TWO          int = 2
	BAD_CASE_THREE        int = 3
	BAD_CASE_FOUR         int = 4
	BAD_CASE_FIVE         int = 5
	BAD_CASE_SIX          int = 6
	BAD_CASE_SEVEN        int = 7
	BAD_CASE_EIGHT        int = 8
	BAD_CASE_NINE         int = 9
	BAD_CASE_TEN          int = 10
)

// ComplexityEleven demonstrates a function with cyclomatic complexity = 11 (exceeds limit of 10). This function intentionally violates KTN-FUNC-005 to test the complexity analyzer.
//
// Params:
//   - x: input value to process
func ComplexityEleven(x int) {
	// Check if x is positive
	if x > 0 {
		x++
	}
	// Check if x exceeds threshold five
	if x > BAD_THRESHOLD_FIVE {
		x++
	}
	// Check if x exceeds threshold ten
	if x > BAD_THRESHOLD_TEN {
		x++
	}
	// First loop iteration
	for i := 0; i < BAD_THRESHOLD_THREE; i++ {
		x++
	}
	// Second loop iteration
	for i := 0; i < BAD_THRESHOLD_THREE; i++ {
		x++
	}
	// Third loop iteration
	for i := 0; i < BAD_THRESHOLD_THREE; i++ {
		x++
	}
	// Switch based on x value
	switch x {
	// Case for value 1
	case 1:
		x++
	// Case for value 2
	case BAD_CASE_TWO:
		x++
	// Case for value 3
	case BAD_CASE_THREE:
		x++
	// Case for value 4
	case BAD_CASE_FOUR:
		x++
	}
}

// VeryComplex is a very complex function with nested conditions and multiple paths. This function intentionally violates KTN-FUNC-005 to test high complexity detection.
//
// Params:
//   - x: input value to evaluate
//
// Returns:
//   - int: processed result based on input
func VeryComplex(x int) int {
	// Check positive range
	if x > 0 {
		// Check threshold ten
		if x > BAD_THRESHOLD_TEN {
			// Check threshold twenty
			if x > BAD_THRESHOLD_TWENTY {
				// Return one for very high positive values
				return 1
			}
		}
	}

	// Check negative range
	if x < 0 {
		// Check threshold minus ten
		if x < BAD_MIN_TEN {
			// Check threshold minus twenty
			if x < BAD_MIN_TWENTY {
				// Return negative one for very low values
				return -1
			}
		}
	}

	// Process value in loop
	for i := 0; i < x; i++ {
		// Check for even numbers
		if i%BAD_THRESHOLD_TWO == 0 {
			x++
		} else {
			// Decrement for odd numbers
			x--
		}
	}

	// Return result based on final value
	switch x {
	// Case for value 1
	case 1:
		// Return one for case one
		return 1
	// Case for value 2
	case BAD_CASE_TWO:
		// Return two for case two
		return BAD_CASE_TWO
	// Case for value 3
	case BAD_CASE_THREE:
		// Return three for case three
		return BAD_CASE_THREE
	// Default case
	default:
		// Return zero for all other cases
		return 0
	}
}

// ManyLogicalOps contains many logical operators increasing complexity. This function intentionally violates KTN-FUNC-005 through multiple conditions.
//
// Params:
//   - a: first boolean condition
//   - b: second boolean condition
//   - c: third boolean condition
//   - d: fourth boolean condition
//
// Returns:
//   - bool: result of complex logical evaluation
func ManyLogicalOps(a, b, c, d bool) bool {
	// Check a and b combination
	if a && b {
		// Return true when both a and b are true
		return true
	}
	// Check c or d combination
	if c || d {
		// Return true when either c or d is true
		return true
	}
	// Check a and c combination
	if a && c {
		// Return true when both a and c are true
		return true
	}
	// Check b or d combination
	if b || d {
		// Return true when either b or d is true
		return true
	}
	// Check a, b, and c combination
	if a && b && c {
		// Return true when all three are true
		return true
	}
	// Check c, d, or a combination
	if c || d || a {
		// Return true when any of these is true
		return true
	}
	// Check complex combination
	if a && b || c && d {
		// Return true for this complex condition
		return true
	}
	// Return false if no conditions matched
	return false
}

// ComplexSelectRange demonstrates high complexity with select and range operations. This function intentionally violates KTN-FUNC-005 through multiple control structures.
//
// Params:
//   - ch: channel to receive or send integers
//   - items: slice of items to process
//
// Returns:
//   - int: computed result from channel and range operations
func ComplexSelectRange(ch chan int, items []int) int {
	result := 0

	// Process each item in the slice
	for _, item := range items {
		// Add positive items to result
		if item > 0 {
			result += item
		}
	}

	// Handle channel communication
	select {
	// Case for receiving from channel
	case x := <-ch:
		// Process received value if above threshold
		if x > BAD_THRESHOLD_TEN {
			result += x
		}
	// Case for sending to channel
	case ch <- result:
		// Increment result if above threshold before sending
		if result > BAD_THRESHOLD_FIVE {
			result++
		}
	// Default case
	default:
		// Reset result for default case
		result = 0
	}

	// Cap result at maximum
	if result > BAD_THRESHOLD_HUNDRED {
		// Limit to hundred
		result = BAD_THRESHOLD_HUNDRED
	}
	// Ensure non-negative result
	if result < 0 {
		// Reset to zero for negative values
		result = 0
	}
	// Adjust even results
	if result%BAD_THRESHOLD_TWO == 0 {
		// Increment even results
		result++
	}
	// Cap result at fifty
	if result > BAD_THRESHOLD_FIFTY {
		// Limit to fifty
		result = BAD_THRESHOLD_FIFTY
	}

	// Return final computed result
	return result
}

// SwitchManyWithDefault has a switch with many cases causing high complexity. Complexity = 12 (1 base + 10 cases + 1 if). This intentionally violates KTN-FUNC-005.
//
// Params:
//   - x: input value to match against cases
//
// Returns:
//   - int: result based on switch case matching
func SwitchManyWithDefault(x int) int {
	result := 0
	// Switch on input value
	switch x {
	// Case for 1
	case 1:
		result = 1
	// Case for 2
	case BAD_CASE_TWO:
		result = BAD_CASE_TWO
	// Case for 3
	case BAD_CASE_THREE:
		result = BAD_CASE_THREE
	// Case for 4
	case BAD_CASE_FOUR:
		result = BAD_CASE_FOUR
	// Case for 5
	case BAD_CASE_FIVE:
		result = BAD_CASE_FIVE
	// Case for 6
	case BAD_CASE_SIX:
		result = BAD_CASE_SIX
	// Case for 7
	case BAD_CASE_SEVEN:
		result = BAD_CASE_SEVEN
	// Case for 8
	case BAD_CASE_EIGHT:
		result = BAD_CASE_EIGHT
	// Case for 9
	case BAD_CASE_NINE:
		result = BAD_CASE_NINE
	// Case for 10
	case BAD_CASE_TEN:
		result = BAD_CASE_TEN
	// Default case
	default:
		result = 0
	}
	// Increment result if above threshold
	if result > BAD_THRESHOLD_FIVE {
		result++
	}
	// Return computed result
	return result
}

// SelectManyWithDefault has a select with many communication cases causing high complexity. Complexity = 11 (1 base + 5 comm cases + 5 ifs). This intentionally violates KTN-FUNC-005.
//
// Params:
//   - ch1: first channel for receiving
//   - ch2: second channel for receiving
//   - ch3: third channel for receiving
//   - ch4: fourth channel for receiving
//   - ch5: fifth channel for receiving
//
// Returns:
//   - int: result from channel operations
func SelectManyWithDefault(ch1, ch2, ch3, ch4, ch5 chan int) int {
	result := 0
	// Select from multiple channels
	select {
	// Case from channel one
	case x := <-ch1:
		// Process value from channel one
		if x > 0 {
			result = x
		}
	// Case from channel two
	case y := <-ch2:
		// Process value from channel two
		if y > 0 {
			result = y
		}
	// Case from channel three
	case z := <-ch3:
		// Process value from channel three
		if z > 0 {
			result = z
		}
	// Case from channel four
	case a := <-ch4:
		// Process value from channel four
		if a > 0 {
			result = a
		}
	// Case from channel five
	case b := <-ch5:
		// Process value from channel five
		if b > 0 {
			result = b
		}
	// Default case
	default:
		// Reset for default case
		result = 0
	}
	// Cap result at maximum
	if result > BAD_THRESHOLD_TEN {
		// Limit to ten
		result = BAD_THRESHOLD_TEN
	}
	// Return final result
	return result
}
