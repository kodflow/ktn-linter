package fall001

func BadFallthroughOutsideSwitch() {
	if true { // want `\[KTN-FALL-001\].*`
		fallthrough
	}
}

func GoodFallthroughInSwitch(x int) {
	switch x {
	case 1:
		process()
		fallthrough
	case 2:
		process()
	}
}

func GoodNoFallthrough(x int) {
	switch x {
	case 1:
		process()
	case 2:
		process()
	}
}

func process() {}
