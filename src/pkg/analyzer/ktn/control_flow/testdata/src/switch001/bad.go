package switch005

func BadSingleCase(x int) {
	switch x { // want `\[KTN-SWITCH-005\].*`
	case 1:
		process()
	}
}

func BadSingleCaseString(s string) {
	switch s { // want `\[KTN-SWITCH-005\].*`
	case "test":
		process()
	}
}

func GoodMultipleCases(x int) {
	switch x {
	case 1:
		process()
	case 2:
		process()
	default:
		process()
	}
}

func GoodIfInstead(x int) {
	if x == 1 {
		process()
	}
}

func process() {}
