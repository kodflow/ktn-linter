// Package func005 provides good test cases.
package func005

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
	TWENTY      int = NINETEEN + 1
	TwentyOne   int = TWENTY + 1
	TwentyTwo   int = TwentyOne + 1
	TwentyThree int = TwentyTwo + 1
	TwentyFour  int = TwentyThree + 1
	TwentyFive  int = TwentyFour + 1
	TwentySix   int = TwentyFive + 1
	TwentySeven int = TwentySix + 1
	TwentyEight int = TwentySeven + 1
	TwentyNine  int = TwentyEight + 1
	THIRTY      int = TwentyNine + 1
	ThirtyOne   int = THIRTY + 1
	ThirtyTwo   int = ThirtyOne + 1
	ThirtyThree int = ThirtyTwo + 1
	ThirtyFour  int = ThirtyThree + 1
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
	u := TwentyOne                                                                                                                                    // 21
	v := TwentyTwo                                                                                                                                    // 22
	w := TwentyThree                                                                                                                                  // 23
	x := TwentyFour                                                                                                                                   // 24
	y := TwentyFive                                                                                                                                   // 25
	z := TwentySix                                                                                                                                    // 26
	aa := TwentySeven                                                                                                                                 // 27
	ab := TwentyEight                                                                                                                                 // 28
	ac := TwentyNine                                                                                                                                  // 29
	ad := THIRTY                                                                                                                                      // 30
	ae := ThirtyOne                                                                                                                                   // 31
	af := ThirtyTwo                                                                                                                                   // 32
	ag := ThirtyThree                                                                                                                                 // 33
	ah := ThirtyFour                                                                                                                                  // 34
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
	for i := range TEN {
		// Dans le for
		z := i * TWO
		// Utilisation de z
		_ = z
	}
	// Retour de la fonction
	return 0
}

// exemptTestFunction démontre qu'une fonction Test* est exemptée de FUNC-001
func exemptTestFunction() {
	// Les fonctions de test peuvent être aussi longues que nécessaire
	for i := range HUNDRED {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// exemptBenchmarkFunction démontre qu'une fonction Benchmark* est exemptée de FUNC-001
func exemptBenchmarkFunction() {
	// Les fonctions de benchmark peuvent être aussi longues que nécessaire
	for i := range HUNDRED {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// exemptExampleFunction démontre qu'une fonction Example* est exemptée de FUNC-001
func exemptExampleFunction() {
	// Les fonctions d'exemple peuvent être aussi longues que nécessaire
	for i := range HUNDRED {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// exemptFuzzFunction démontre qu'une fonction Fuzz* est exemptée de FUNC-001
func exemptFuzzFunction() {
	// Les fonctions de fuzzing peuvent être aussi longues que nécessaire
	for i := range HUNDRED {
		x := i * TWO
		y := i * THREE
		z := x + y
		_ = z
	}
}

// main est la fonction principale exemptée
func main() {
	// La fonction main peut être aussi longue que nécessaire
	for i := range HUNDRED {
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
	/* Block comment court */
	x := 1

	/* Autre block comment */
	y := TWO

	// Ligne de commentaire normale
	z := x + y

	// Retour de la fonction
	return z
}

// init utilise toutes les fonctions privées pour démontrer qu'elles sont correctement définies
func init() {
	// Utilisation de smallFunction
	_ = smallFunction()
	// Utilisation de exactlyThirtyFive
	exactlyThirtyFive()
	// Utilisation de manyCommentsButFewStatements
	manyCommentsButFewStatements()
	// Utilisation de withNestedBlocks
	_ = withNestedBlocks()
	// Utilisation de functionWithBlockComments
	_ = functionWithBlockComments()
	// Utilisation des fonctions exemptées
	exemptTestFunction()
	exemptBenchmarkFunction()
	exemptExampleFunction()
	exemptFuzzFunction()
}
