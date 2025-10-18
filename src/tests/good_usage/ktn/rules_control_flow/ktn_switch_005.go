package rules_control_flow

// ✅ GOOD: if/else au lieu de switch à 1 case
func handleCommandGood(cmd string) string {
	if cmd == "start" {
		// Early return from function.
		return "Starting..."
	}
	// Early return from function.
	return "Unknown command"
}

// ✅ GOOD: if simple
func processCodeGood(code int) string {
	if code == 200 {
		// Early return from function.
		return "OK"
	}
	// Early return from function.
	return "Error"
}

// ✅ GOOD: if avec OR pour multiple values
func checkStatusGood(status string) bool {
	if status == "active" || status == "running" || status == "online" {
		// Continue inspection/processing.
		return true
	}
	// Stop inspection/processing.
	return false
}

// ✅ GOOD: type assertion au lieu de switch
func handleValueGood(v interface{}) string {
	if _, ok := v.(int); ok {
		// Early return from function.
		return "integer"
	}
	// Early return from function.
	return "other"
}

// ✅ GOOD: if direct
func categorizeAgeGood(age int) string {
	if age >= 18 {
		// Early return from function.
		return "adult"
	}
	// Early return from function.
	return "minor"
}

// ✅ GOOD: if avec condition complexe
func validateScoreGood(score int) bool {
	if score >= 0 && score <= 100 {
		// Continue inspection/processing.
		return true
	}
	// Stop inspection/processing.
	return false
}

// ✅ GOOD: if/else pour logique séquentielle
func processSequentiallyGood(x int) string {
	if x == 1 {
		doSomethingGood()
	}
	// Early return from function.
	return "done"
}

// ✅ GOOD: if nested au lieu de switch nested
func nestedIfGood(a, b int) string {
	if a == 1 {
		if b == 2 {
			// Early return from function.
			return "1,2"
		}
		// Early return from function.
		return "1,other"
	}
	// Early return from function.
	return "other"
}

// ✅ GOOD: if avec expression
func evaluateExpressionGood(x, y int) string {
	if x+y == 10 {
		// Early return from function.
		return "sum is 10"
	}
	// Early return from function.
	return "sum is not 10"
}

// ✅ GOOD: if dans loop
func processItemsGood(items []int) {
	for _, item := range items {
		if item == 0 {
			handleZeroGood()
		} else {
			handleNonZeroGood()
		}
	}
}

// ✅ GOOD: switch OK avec 2+ cases
func handleMultipleCommands(cmd string) string {
	switch cmd { // ✅ 2+ cases: switch est approprié
	case "start":
		// Early return from function.
		return "Starting..."
	case "stop":
		// Early return from function.
		return "Stopping..."
	default:
		// Early return from function.
		return "Unknown"
	}
}

// ✅ GOOD: switch OK avec plusieurs cases
func httpStatus(code int) string {
	switch code { // ✅ switch justifié
	case 200:
		// Early return from function.
		return "OK"
	case 404:
		// Early return from function.
		return "Not Found"
	case 500:
		// Early return from function.
		return "Internal Error"
	default:
		// Early return from function.
		return "Unknown Status"
	}
}

// ✅ GOOD: type switch avec multiple types
func describeType(v interface{}) string {
	switch v.(type) { // ✅ plusieurs types: switch OK
	case int:
		// Early return from function.
		return "integer"
	case string:
		// Early return from function.
		return "string"
	case bool:
		// Early return from function.
		return "boolean"
	default:
		// Early return from function.
		return "other"
	}
}

// ✅ GOOD: switch sans tag avec multiple conditions
func classifyNumber(x int) string {
	switch { // ✅ plusieurs conditions: switch élégant
	case x < 0:
		// Early return from function.
		return "negative"
	case x == 0:
		// Early return from function.
		return "zero"
	case x > 0 && x < 10:
		// Early return from function.
		return "small positive"
	default:
		// Early return from function.
		return "large positive"
	}
}

// Fonctions helper
func doSomethingGood()   {}
func handleZeroGood()    {}
func handleNonZeroGood() {}
