package return001

// TODO: L'analyzer return nécessite des améliorations
// Il détecte uniquement les naked returns dans les fonctions TRÈS longues

func ShortFunction(x int) (result int) {
	result = x * 2
	return result
}

func GoodExplicitReturn(x int) (result int) {
	result = x * 2
	return result
}

func GoodSimpleReturn(x int) int {
	return x * 2
}
