package func012

// checkPositive uses early return (no else)
func checkPositive(x int) string {
	if x > 0 {
		return "positive"
	}
	return "negative"
}

// processValue uses early return (no else)
func processValue(val int) int {
	if val < 0 {
		return 0
	}
	return val * 2
}

// findMax uses early return (no else)
func findMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// loopExample uses continue without else
func loopExample() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		_ = i
	}
}

// switchExample uses break without else
func switchExample(x int) {
	for {
		if x > 10 {
			break
		}
		x++
	}
}

// validateInput uses early return (no else)
func validateInput(input string) error {
	if input == "" {
		return nil
	}
	return nil
}

// complexLogic else is acceptable when if doesn't have return/break/continue
func complexLogic(x int) int {
	if x > 0 {
		x = x * 2
	} else {
		x = x * 3
	}
	return x
}

// nestedConditions multiple conditions without early exit
func nestedConditions(a, b int) int {
	if a > 0 {
		if b > 0 {
			return a + b
		}
		return a
	}
	return b
}
