package rules_control_flow

// ✅ GOOD: if/else au lieu de switch à 1 case
func handleCommandGood(cmd string) string {
	if cmd == "start" {
		return "Starting..."
	}
	return "Unknown command"
}

// ✅ GOOD: if simple
func processCodeGood(code int) string {
	if code == 200 {
		return "OK"
	}
	return "Error"
}

// ✅ GOOD: if avec OR pour multiple values
func checkStatusGood(status string) bool {
	if status == "active" || status == "running" || status == "online" {
		return true
	}
	return false
}

// ✅ GOOD: type assertion au lieu de switch
func handleValueGood(v interface{}) string {
	if _, ok := v.(int); ok {
		return "integer"
	}
	return "other"
}

// ✅ GOOD: if direct
func categorizeAgeGood(age int) string {
	if age >= 18 {
		return "adult"
	}
	return "minor"
}

// ✅ GOOD: if avec condition complexe
func validateScoreGood(score int) bool {
	if score >= 0 && score <= 100 {
		return true
	}
	return false
}

// ✅ GOOD: if/else pour logique séquentielle
func processSequentiallyGood(x int) string {
	if x == 1 {
		doSomethingGood()
	}
	return "done"
}

// ✅ GOOD: if nested au lieu de switch nested
func nestedIfGood(a, b int) string {
	if a == 1 {
		if b == 2 {
			return "1,2"
		}
		return "1,other"
	}
	return "other"
}

// ✅ GOOD: if avec expression
func evaluateExpressionGood(x, y int) string {
	if x+y == 10 {
		return "sum is 10"
	}
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
		return "Starting..."
	case "stop":
		return "Stopping..."
	default:
		return "Unknown"
	}
}

// ✅ GOOD: switch OK avec plusieurs cases
func httpStatus(code int) string {
	switch code { // ✅ switch justifié
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Error"
	default:
		return "Unknown Status"
	}
}

// ✅ GOOD: type switch avec multiple types
func describeType(v interface{}) string {
	switch v.(type) { // ✅ plusieurs types: switch OK
	case int:
		return "integer"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "other"
	}
}

// ✅ GOOD: switch sans tag avec multiple conditions
func classifyNumber(x int) string {
	switch { // ✅ plusieurs conditions: switch élégant
	case x < 0:
		return "negative"
	case x == 0:
		return "zero"
	case x > 0 && x < 10:
		return "small positive"
	default:
		return "large positive"
	}
}

// Fonctions helper
func doSomethingGood()   {}
func handleZeroGood()    {}
func handleNonZeroGood() {}
