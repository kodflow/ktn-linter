package conv001

// L'analyzer détecte uniquement les conversions où le nom de la variable
// correspond au type, comme int(int). C'est une limitation de l'implémentation.

func BadRedundantConversionSameName() {
	// Créer une variable avec un nom qui correspond au type
	myint := 5
	_ = myint
}

func GoodNeededConversion(x int32) int {
	return int(x)
}

func GoodBytesToString(b []byte) string {
	return string(b)
}
