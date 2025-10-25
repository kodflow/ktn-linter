package func001

// Constantes pour éviter les magic numbers
const (
	TWO          int = 1 + 1
	THREE        int = TWO + 1
	FOUR         int = THREE + 1
	FIVE         int = FOUR + 1
	SIX          int = FIVE + 1
	SEVEN        int = SIX + 1
	EIGHT        int = SEVEN + 1
	NINE         int = EIGHT + 1
	TEN          int = NINE + 1
	HUNDRED      int = TEN * TEN
	ELEVEN       int = TEN + 1
	TWELVE       int = ELEVEN + 1
	THIRTEEN     int = TWELVE + 1
	FOURTEEN     int = THIRTEEN + 1
	FIFTEEN      int = FOURTEEN + 1
	SIXTEEN      int = FIFTEEN + 1
	SEVENTEEN    int = SIXTEEN + 1
	EIGHTEEN     int = SEVENTEEN + 1
	NINETEEN     int = EIGHTEEN + 1
	TWENTY       int = NINETEEN + 1
	TWENTY_ONE   int = TWENTY + 1
	TWENTY_TWO   int = TWENTY_ONE + 1
	TWENTY_THREE int = TWENTY_TWO + 1
	TWENTY_FOUR  int = TWENTY_THREE + 1
	TWENTY_FIVE  int = TWENTY_FOUR + 1
	TWENTY_SIX   int = TWENTY_FIVE + 1
	TWENTY_SEVEN int = TWENTY_SIX + 1
	TWENTY_EIGHT int = TWENTY_SEVEN + 1
	TWENTY_NINE  int = TWENTY_EIGHT + 1
	THIRTY       int = TWENTY_NINE + 1
	THIRTY_ONE   int = THIRTY + 1
	THIRTY_TWO   int = THIRTY_ONE + 1
	THIRTY_THREE int = THIRTY_TWO + 1
	THIRTY_FOUR  int = THIRTY_THREE + 1
)

// smallFunction est une petite fonction avec peu de code
//
// Returns:
//   - string: résultat du calcul
func smallFunction() string {
	// Déclaration des variables
	x := 1
	y := 1 + 1
	z := x + y
	// Retour de la fonction
	_ = z
	// Retour de la fonction
	return "result"
}

// exactlyThirtyFive a exactement 35 lignes de code pur (limite)
func exactlyThirtyFive() {
	a := 1                                                                                                                                            // 1
	b := TWO                                                                                                                                          // 2
	c := THREE                                                                                                                                        // 3
	d := FOUR                                                                                                                                         // 4
	e := FIVE                                                                                                                                         // 5
	f := SIX                                                                                                                                          // 6
	g := SEVEN                                                                                                                                        // 7
	h := EIGHT                                                                                                                                        // 8
	i := NINE                                                                                                                                         // 9
	j := TEN                                                                                                                                          // 10
	k := ELEVEN                                                                                                                                       // 11
	l := TWELVE                                                                                                                                       // 12
	m := THIRTEEN                                                                                                                                     // 13
	n := FOURTEEN                                                                                                                                     // 14
	o := FIFTEEN                                                                                                                                      // 15
	p := SIXTEEN                                                                                                                                      // 16
	q := SEVENTEEN                                                                                                                                    // 17
	r := EIGHTEEN                                                                                                                                     // 18
	s := NINETEEN                                                                                                                                     // 19
	t := TWENTY                                                                                                                                       // 20
	u := TWENTY_ONE                                                                                                                                   // 21
	v := TWENTY_TWO                                                                                                                                   // 22
	w := TWENTY_THREE                                                                                                                                 // 23
	x := TWENTY_FOUR                                                                                                                                  // 24
	y := TWENTY_FIVE                                                                                                                                  // 25
	z := TWENTY_SIX                                                                                                                                   // 26
	aa := TWENTY_SEVEN                                                                                                                                // 27
	ab := TWENTY_EIGHT                                                                                                                                // 28
	ac := TWENTY_NINE                                                                                                                                 // 29
	ad := THIRTY                                                                                                                                      // 30
	ae := THIRTY_ONE                                                                                                                                  // 31
	af := THIRTY_TWO                                                                                                                                  // 32
	ag := THIRTY_THREE                                                                                                                                // 33
	ah := THIRTY_FOUR                                                                                                                                 // 34
	_ = a + b + c + d + e + f + g + h + i + j + k + l + m + n + o + p + q + r + s + t + u + v + w + x + y + z + aa + ab + ac + ad + ae + af + ag + ah // 35
}

// manyCommentsButFewStatements démontre que les commentaires ne comptent pas
func manyCommentsButFewStatements() {
	// Ce commentaire ne compte pas
	// Ni celui-ci
	// Ni celui-là
	// Les commentaires sont ignorés
	// Encore un commentaire
	// Toujours des commentaires
	// Plus de commentaires
	// Commentaires partout
	x := 1 // Statement 1
	// Commentaire entre statements
	y := 1 + 1 // Statement 2
	// Encore un commentaire
	z := x + y // Statement 3
	// Commentaire final
	_ = z // Statement 4
}

// withNestedBlocks démontre que les blocs imbriqués sont comptés correctement
//
// Returns:
//   - int: résultat
func withNestedBlocks() int {
	// Déclaration de la variable
	x := 1
	// Condition if
	if x > 0 {
		// Dans le if
		y := 1 + 1
		// Retour de la fonction
		return y
	}
	// Boucle for
	for i := 0; i < TEN; i++ {
		// Dans le for
		z := i * TWO
		// Utilisation de z
		_ = z
	}
	// Retour de la fonction
	return 0
}

// TestSomething est une fonction de test exemptée
func TestSomething() {
	// Les fonctions de test peuvent être aussi longues que nécessaire
	for i := 0; i < HUNDRED; i++ {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// BenchmarkSomething est une fonction de benchmark exemptée
func BenchmarkSomething() {
	// Les fonctions de benchmark peuvent être aussi longues que nécessaire
	for i := 0; i < HUNDRED; i++ {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// ExampleSomething est une fonction d'exemple exemptée
func ExampleSomething() {
	// Les fonctions d'exemple peuvent être aussi longues que nécessaire
	for i := 0; i < HUNDRED; i++ {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// FuzzSomething est une fonction de fuzzing exemptée
func FuzzSomething() {
	// Les fonctions de fuzzing peuvent être aussi longues que nécessaire
	for i := 0; i < HUNDRED; i++ {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// main est la fonction principale exemptée
func main() {
	// La fonction main peut être aussi longue que nécessaire
	for i := 0; i < HUNDRED; i++ {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// functionWithBlockComments teste les commentaires de bloc
//
// Returns:
//   - int: résultat
func functionWithBlockComments() int {
	/* Début du block comment
	   Ligne dans le block comment
	   Fin du block comment */
	x := 1

	/* Un autre block comment */
	y := TWO

	// Ligne de commentaire normale
	z := x + y

	// Retour de la fonction
	return z
}
