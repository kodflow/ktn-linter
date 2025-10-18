package fall001

// TODO: L'analyseur FALL-001 détecte fallthrough hors switch
// MAIS ce code ne compile pas (erreur de syntaxe Go)
// On ne peut donc pas tester ce cas
//
// func BadFallthroughOutsideSwitch() {
// 	if true {
// 		fallthrough  // Erreur compilation: "fallthrough statement out of place"
// 	}
// }

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
