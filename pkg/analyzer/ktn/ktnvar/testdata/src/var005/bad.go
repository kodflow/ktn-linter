package var005

// Bad examples - names too long (over 30 chars)

var thisIsAVeryLongVariableNameThatExceedsLimit int = 1                    // want "KTN-VAR-005"
var anotherExtremelyLongVariableNameForSomeConfigurationValue string = "x" // want "KTN-VAR-005"

func badShortDecl() {
	thisIsAVeryLongLocalVariableNameOverThirty := 42 // want "KTN-VAR-005"
	_ = thisIsAVeryLongLocalVariableNameOverThirty
}

func badBlockVar() {
	var thisIsAnotherVeryLongVariableNameOverMax int = 10 // want "KTN-VAR-005"
	_ = thisIsAnotherVeryLongVariableNameOverMax
}

func badRangeStmt() {
	items := []int{1, 2, 3}
	for thisIsAVeryLongRangeKeyVariableNameX := range items { // want "KTN-VAR-005"
		_ = thisIsAVeryLongRangeKeyVariableNameX
	}
}

func badRangeWithValue() {
	items := []int{1, 2, 3}
	for _, thisIsAVeryLongRangeValueVariableName := range items { // want "KTN-VAR-005"
		_ = thisIsAVeryLongRangeValueVariableName
	}
}
