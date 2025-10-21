package func012

// badCheckPositive has unnecessary else after return
func badCheckPositive(x int) string {
	if x > 0 {
		return "positive"
	} else {
		return "negative"
	}
}

// badProcessValue has unnecessary else after return
func badProcessValue(val int) int {
	if val < 0 {
		return 0
	} else {
		return val * 2
	}
}

// badFindMax has unnecessary else after return
func badFindMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// badLoopExample has unnecessary else after continue
func badLoopExample() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		} else {
			_ = i
		}
	}
}

// badSwitchExample has unnecessary else after break
func badSwitchExample(x int) {
	for {
		if x > 10 {
			break
		} else {
			x++
		}
	}
}

// badValidateInput has unnecessary else after return
func badValidateInput(input string) error {
	if input == "" {
		return nil
	} else {
		return nil
	}
}
