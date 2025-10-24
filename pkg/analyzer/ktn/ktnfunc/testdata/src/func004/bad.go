package func004

const (
	// INCREMENT_ONE représente l'incrément de 1
	INCREMENT_ONE int = 1
	// INCREMENT_TWO représente l'incrément de 2
	INCREMENT_TWO int = 2
	// INCREMENT_THREE représente l'incrément de 3
	INCREMENT_THREE int = 3
	// INCREMENT_FOUR représente l'incrément de 4
	INCREMENT_FOUR int = 4
	// INCREMENT_FIVE représente l'incrément de 5
	INCREMENT_FIVE int = 5
	// INCREMENT_SIX représente l'incrément de 6
	INCREMENT_SIX int = 6
	// INCREMENT_TEN représente l'incrément de 10
	INCREMENT_TEN int = 10
	// INCREMENT_TWENTY représente l'incrément de 20
	INCREMENT_TWENTY int = 20
	// INCREMENT_THIRTY représente l'incrément de 30
	INCREMENT_THIRTY int = 30
	// INCREMENT_FORTY représente l'incrément de 40
	INCREMENT_FORTY int = 40
	// INCREMENT_FIFTY représente l'incrément de 50
	INCREMENT_FIFTY int = 50
)

// badLongWithNakedReturn utilise naked return dans fonction de 5 lignes
//
// Returns:
//   - result: résultat calculé
func badLongWithNakedReturn() (result int) {
	result = INCREMENT_ONE
	result += INCREMENT_TWO
	result += INCREMENT_THREE
	result += INCREMENT_FOUR
	result += INCREMENT_FIVE
	// Retour naked interdit car >= 5 lignes
	return
}

// badVeryLongNakedReturn utilise naked return dans fonction longue
//
// Returns:
//   - a: premier entier
//   - b: chaîne de caractères
func badVeryLongNakedReturn() (a int, b string) {
	a = INCREMENT_ONE
	a += INCREMENT_TWO
	a += INCREMENT_THREE
	a += INCREMENT_FOUR
	a += INCREMENT_FIVE
	a += INCREMENT_SIX
	b = "test"
	// Retour naked interdit car >= 5 lignes
	return
}

// badMultipleNakedReturns utilise plusieurs naked returns dans fonction longue
//
// Returns:
//   - result: résultat calculé
func badMultipleNakedReturns() (result int) {
	// Vérification de la condition
	if true {
		result = INCREMENT_ONE
		result += INCREMENT_TWO
		result += INCREMENT_THREE
		result += INCREMENT_FOUR
		result += INCREMENT_FIVE
		// Retour naked interdit car >= 5 lignes
		return
	}
	result = INCREMENT_TEN
	result += INCREMENT_TWENTY
	result += INCREMENT_THIRTY
	result += INCREMENT_FORTY
	result += INCREMENT_FIFTY
	// Retour naked interdit car >= 5 lignes
	return
}
