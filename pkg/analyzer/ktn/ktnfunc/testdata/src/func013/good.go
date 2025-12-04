// Good examples for the func013 test case.
package func013

const (
	// THRESHOLD_FIVE represents a threshold value of 5
	THRESHOLD_FIVE int = 5
	// THRESHOLD_TEN represents a threshold value of 10
	THRESHOLD_TEN int = 10
	// LOOP_ITERATIONS represents the number of iterations in loops
	LOOP_ITERATIONS int = 3
	// CASE_VALUE_TWO represents case value 2
	CASE_VALUE_TWO int = 2
	// CASE_VALUE_THREE represents case value 3
	CASE_VALUE_THREE int = 3
	// CASE_VALUE_FOUR represents case value 4
	CASE_VALUE_FOUR int = 4
	// CASE_VALUE_FIVE represents case value 5
	CASE_VALUE_FIVE int = 5
)

// Simple is a simple function with minimal complexity.
// Returns: nothing
func Simple() {
	// Initialize x with value 1
	x := 1
	// Use x to avoid unused variable
	_ = x
}

// OneIf demonstrates a function with a single if statement.
// Params:
//   - x: input value to check
//
// Returns: nothing
func OneIf(x int) {
	// Check if x is positive
	if x > 0 {
		// Use x to avoid unused variable
		_ = x
	}
}

// ComplexityTen demonstrates a function with complexity exactly at the limit of 10 (1 base + 3 if + 3 for + 3 case).
// Params:
//   - x: input value to process
//
// Returns: nothing
func ComplexityTen(x int) {
	// Check if x is positive
	if x > 0 {
		// Increment x
		x++
	}
	// Check if x exceeds threshold of 5
	if x > THRESHOLD_FIVE {
		// Increment x
		x++
	}
	// Check if x exceeds threshold of 10
	if x > THRESHOLD_TEN {
		// Increment x
		x++
	}

	// First loop to increment x
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}
	// Second loop to increment x
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}
	// Third loop to increment x
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}

	// Switch to handle different cases of x
	switch x {
	// Handle case 1
	case 1:
		x++
	// Handle case 2
	case CASE_VALUE_TWO:
		x++
	// Handle case 3
	case CASE_VALUE_THREE:
		x++
	}
}

// MultipleIfs demonstrates a function with multiple if statements.
// Params:
//   - x: input value to evaluate
//
// Returns:
//   - int: normalized value based on input
func MultipleIfs(x int) int {
	// Check if x is positive
	if x > 0 {
		// Return 1 for positive values
		return 1
	}
	// Check if x is negative
	if x < 0 {
		// Return -1 for negative values
		return -1
	}
	// Check if x is zero
	if x == 0 {
		// Return 0 for zero value
		return 0
	}
	// Check if x exceeds upper threshold
	if x > THRESHOLD_TEN {
		// Return capped value
		return THRESHOLD_TEN
	}
	// Check if x is below lower threshold
	if x < -THRESHOLD_TEN {
		// Return capped negative value
		return -THRESHOLD_TEN
	}
	// Return original value
	return x
}

// SwitchWithDefault demonstrates a switch statement with default case (complexity 4: 1 base + 3 cases).
// Params:
//   - x: input value to match
//
// Returns:
//   - int: matched case value or 0
func SwitchWithDefault(x int) int {
	// Switch on x value
	switch x {
	// Return 1 for case 1
	case 1:
		// Return 1 for case 1
		return 1
	// Return 2 for case 2
	case CASE_VALUE_TWO:
		// Return 2 for case 2
		return CASE_VALUE_TWO
	// Return 3 for case 3
	case CASE_VALUE_THREE:
		// Return 3 for case 3
		return CASE_VALUE_THREE
	// Return 0 for default case
	default:
		// Return 0 for default case
		return 0
	}
}

// SelectWithDefault demonstrates a select statement with channels (complexity 3: 1 base + 2 comm cases).
// Params:
//   - ch1: first channel to receive from
//   - ch2: second channel to receive from
//
// Returns:
//   - int: received value or 0
func SelectWithDefault(ch1, ch2 chan int) int {
	// Select from channels
	select {
	case x := <-ch1:
		// Return value from first channel
		return x
	case y := <-ch2:
		// Return value from second channel
		return y
	default:
		// Return 0 if no channel ready
		return 0
	}
}

// RangeLoop demonstrates a range loop over a slice.
// Params:
//   - items: slice of integers to sum
//
// Returns:
//   - int: sum of all items
func RangeLoop(items []int) int {
	// Initialize sum
	sum := 0
	// Iterate over all items
	for _, item := range items {
		// Add item to sum
		sum += item
	}
	// Return total sum
	return sum
}

// TestComplexFunction is a test helper function that should be skipped by linting.
// Params:
//   - t: test parameter
//
// Returns: nothing
func TestComplexFunction(t int) {
	// Check if t is positive
	if t > 0 {
		// Increment t
		t++
	}
	// Check if t exceeds threshold of 5
	if t > THRESHOLD_FIVE {
		// Increment t
		t++
	}
	// Check if t exceeds threshold of 10
	if t > THRESHOLD_TEN {
		// Increment t
		t++
	}
	// First loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment t in loop
		t++
	}
	// Second loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment t in loop
		t++
	}
	// Third loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment t in loop
		t++
	}
	// Switch on t value
	switch t {
	case 1:
		// Handle case 1
		t++
	case CASE_VALUE_TWO:
		// Handle case 2
		t++
	case CASE_VALUE_THREE:
		// Handle case 3
		t++
	case CASE_VALUE_FOUR:
		// Handle case 4
		t++
	case CASE_VALUE_FIVE:
		// Handle case 5
		t++
	}
}

// BenchmarkComplexFunction is a benchmark helper function that should be skipped by linting.
// Params:
//   - b: benchmark parameter
//
// Returns: nothing
func BenchmarkComplexFunction(b int) {
	// Check if b is positive
	if b > 0 {
		// Increment b
		b++
	}
	// Check if b exceeds threshold of 5
	if b > THRESHOLD_FIVE {
		// Increment b
		b++
	}
	// Check if b exceeds threshold of 10
	if b > THRESHOLD_TEN {
		// Increment b
		b++
	}
	// First loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment b in loop
		b++
	}
	// Second loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment b in loop
		b++
	}
	// Third loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment b in loop
		b++
	}
	// Switch on b value
	switch b {
	case 1:
		// Handle case 1
		b++
	case CASE_VALUE_TWO:
		// Handle case 2
		b++
	case CASE_VALUE_THREE:
		// Handle case 3
		b++
	case CASE_VALUE_FOUR:
		// Handle case 4
		b++
	case CASE_VALUE_FIVE:
		// Handle case 5
		b++
	}
}

// ExampleComplexFunction is an example function that should be skipped by linting.
// Returns: nothing
func ExampleComplexFunction() {
	// Initialize x
	x := 0
	// Check if x is positive
	if x > 0 {
		// Increment x
		x++
	}
	// Check if x exceeds threshold of 5
	if x > THRESHOLD_FIVE {
		// Increment x
		x++
	}
	// Check if x exceeds threshold of 10
	if x > THRESHOLD_TEN {
		// Increment x
		x++
	}
	// First loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}
	// Second loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}
	// Third loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment x in loop
		x++
	}
	// Switch on x value
	switch x {
	case 1:
		// Handle case 1
		x++
	case CASE_VALUE_TWO:
		// Handle case 2
		x++
	case CASE_VALUE_THREE:
		// Handle case 3
		x++
	case CASE_VALUE_FOUR:
		// Handle case 4
		x++
	case CASE_VALUE_FIVE:
		// Handle case 5
		x++
	}
}

// FuzzComplexFunction is a fuzz test function that should be skipped by linting.
// Params:
//   - f: fuzz parameter
//
// Returns: nothing
func FuzzComplexFunction(f int) {
	// Check if f is positive
	if f > 0 {
		// Increment f
		f++
	}
	// Check if f exceeds threshold of 5
	if f > THRESHOLD_FIVE {
		// Increment f
		f++
	}
	// Check if f exceeds threshold of 10
	if f > THRESHOLD_TEN {
		// Increment f
		f++
	}
	// First loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment f in loop
		f++
	}
	// Second loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment f in loop
		f++
	}
	// Third loop iteration
	for i := 0; i < LOOP_ITERATIONS; i++ {
		// Increment f in loop
		f++
	}
	// Switch on f value
	switch f {
	case 1:
		// Handle case 1
		f++
	case CASE_VALUE_TWO:
		// Handle case 2
		f++
	case CASE_VALUE_THREE:
		// Handle case 3
		f++
	case CASE_VALUE_FOUR:
		// Handle case 4
		f++
	case CASE_VALUE_FIVE:
		// Handle case 5
		f++
	}
}
