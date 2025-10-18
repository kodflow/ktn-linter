package switch005

func BadSingleCase(x int) {
	// want `\[KTN-CONTROL-SWITCH-005\] Switch avec un seul case`
	switch x {
	case 1:
		process()
	}
}

func BadSingleCaseString(s string) {
	// want `\[KTN-CONTROL-SWITCH-005\] Switch avec un seul case`
	switch s {
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
