package func001

// smallFunction est une petite fonction avec peu de code
//
// Returns:
//   - string: résultat du calcul
func smallFunction() string {
	// Déclaration des variables
	x := 1
	y := 2
	z := x + y
	// Retour de la fonction
	_ = z
	// Retour de la fonction
	return "result"
}

// exactlyThirtyFive a exactement 35 lignes de code pur (limite)
func exactlyThirtyFive() {
	a := 1                                                                                                                                            // 1
	b := 2                                                                                                                                            // 2
	c := 3                                                                                                                                            // 3
	d := 4                                                                                                                                            // 4
	e := 5                                                                                                                                            // 5
	f := 6                                                                                                                                            // 6
	g := 7                                                                                                                                            // 7
	h := 8                                                                                                                                            // 8
	i := 9                                                                                                                                            // 9
	j := 10                                                                                                                                           // 10
	k := 11                                                                                                                                           // 11
	l := 12                                                                                                                                           // 12
	m := 13                                                                                                                                           // 13
	n := 14                                                                                                                                           // 14
	o := 15                                                                                                                                           // 15
	p := 16                                                                                                                                           // 16
	q := 17                                                                                                                                           // 17
	r := 18                                                                                                                                           // 18
	s := 19                                                                                                                                           // 19
	t := 20                                                                                                                                           // 20
	u := 21                                                                                                                                           // 21
	v := 22                                                                                                                                           // 22
	w := 23                                                                                                                                           // 23
	x := 24                                                                                                                                           // 24
	y := 25                                                                                                                                           // 25
	z := 26                                                                                                                                           // 26
	aa := 27                                                                                                                                          // 27
	ab := 28                                                                                                                                          // 28
	ac := 29                                                                                                                                          // 29
	ad := 30                                                                                                                                          // 30
	ae := 31                                                                                                                                          // 31
	af := 32                                                                                                                                          // 32
	ag := 33                                                                                                                                          // 33
	ah := 34                                                                                                                                          // 34
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
	y := 2 // Statement 2
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
		y := 2
		// Retour de la fonction
		return y
	}
	// Boucle for
	for i := 0; i < 10; i++ {
		// Dans le for
		z := i * 2
		// Utilisation de z
		_ = z
	}
	// Retour de la fonction
	return 0
}

// TestSomething est une fonction de test exemptée
func TestSomething() {
	// Les fonctions de test peuvent être aussi longues que nécessaire
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// BenchmarkSomething est une fonction de benchmark exemptée
func BenchmarkSomething() {
	// Les fonctions de benchmark peuvent être aussi longues que nécessaire
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// ExampleSomething est une fonction d'exemple exemptée
func ExampleSomething() {
	// Les fonctions d'exemple peuvent être aussi longues que nécessaire
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// FuzzSomething est une fonction de fuzzing exemptée
func FuzzSomething() {
	// Les fonctions de fuzzing peuvent être aussi longues que nécessaire
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// main est la fonction principale exemptée
func main() {
	// La fonction main peut être aussi longue que nécessaire
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}
