package fall001

func BadFallthroughOutsideSwitch() {
	// want `\[KTN-CONTROL-FALL-001\] fallthrough en dehors d'un switch`
	if true {
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
