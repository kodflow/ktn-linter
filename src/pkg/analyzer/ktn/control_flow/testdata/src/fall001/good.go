package fall001

// Cas corrects - l'analyseur ne devrait pas rapporter d'erreur

func GoodExample() {
	// Code conforme aux règles
}

func GoodFunction(x int) int {
	return x * 2
}

// GoodFallthroughInSwitch - utilisation correcte de fallthrough
func GoodFallthroughInSwitch(x int) string {
	result := ""
	switch x {
	case 1:
		result += "one"
		fallthrough // OK: fallthrough dans un switch
	case 2:
		result += "two"
	case 3:
		result += "three"
		fallthrough // OK: fallthrough dans un switch
	default:
		result += "other"
	}
	return result
}

// GoodMultipleFallthrough - plusieurs fallthrough dans un switch
func GoodMultipleFallthrough(x int) int {
	result := 0
	switch x {
	case 1:
		result += 1
		fallthrough
	case 2:
		result += 2
		fallthrough
	case 3:
		result += 3
	default:
		result += 10
	}
	return result
}
