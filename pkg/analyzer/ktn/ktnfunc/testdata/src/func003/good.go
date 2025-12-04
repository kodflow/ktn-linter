// Good examples for the func003 test case.
package func003

const (
	// MULTIPLIER_DOUBLE constante pour doubler une valeur
	MULTIPLIER_DOUBLE int = 2
	// MULTIPLIER_TRIPLE constante pour tripler une valeur
	MULTIPLIER_TRIPLE int = 3
	// MAX_LOOP_ITERATIONS nombre maximum d'itérations de boucle
	MAX_LOOP_ITERATIONS int = 10
	// MODULO_EVEN modulo pour vérifier les nombres pairs
	MODULO_EVEN int = 2
)

// checkPositive vérifie si un nombre est positif
//
// Params:
//   - x: nombre à vérifier
//
// Returns:
//   - string: "positive" ou "negative"
func checkPositive(x int) string {
	// Vérification si positif
	if x > 0 {
		// Retour cas positif
		return "positive"
	}
	// Retour cas négatif ou zéro
	return "negative"
}

// processValue traite une valeur en la doublant si positive
//
// Params:
//   - val: valeur à traiter
//
// Returns:
//   - int: 0 si négatif, sinon valeur doublée
func processValue(val int) int {
	// Vérification si négatif
	if val < 0 {
		// Retour zéro pour valeur négative
		return 0
	}
	// Retour valeur doublée
	return val * MULTIPLIER_DOUBLE
}

// findMax trouve le maximum entre deux nombres
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: le maximum des deux nombres
func findMax(a, b int) int {
	// Comparaison a > b
	if a > b {
		// Retour a si supérieur
		return a
	}
	// Retour b sinon
	return b
}

// loopExample illustre l'utilisation de continue sans else
func loopExample() {
	// Boucle sur les 10 premières valeurs
	for i := 0; i < MAX_LOOP_ITERATIONS; i++ {
		// Vérification si pair
		if i%MODULO_EVEN == 0 {
			// Continue si pair
			continue
		}
		// Utilisation de la valeur impaire
		_ = i
	}
}

// switchExample illustre l'utilisation de break sans else
//
// Params:
//   - x: valeur initiale
func switchExample(x int) {
	// Boucle infinie avec condition de sortie
	for {
		// Vérification condition de sortie
		if x > MAX_LOOP_ITERATIONS {
			// Sortie de boucle
			break
		}
		// Incrémentation
		x++
	}
}

// validateInput valide une entrée
//
// Params:
//   - input: chaîne à valider
//
// Returns:
//   - error: erreur de validation ou nil
func validateInput(input string) error {
	// Vérification si vide
	if input == "" {
		// Retour nil si vide
		return nil
	}
	// Retour nil pour entrée valide
	return nil
}

// complexLogic applique une logique avec else acceptable
//
// Params:
//   - x: nombre à traiter
//
// Returns:
//   - int: résultat après transformation
func complexLogic(x int) int {
	// Vérification si positif
	if x > 0 {
		// Doublement si positif
		x = x * MULTIPLIER_DOUBLE
	} else {
		// Triplement si négatif ou nul
		x = x * MULTIPLIER_TRIPLE
	}
	// Retour résultat transformé
	return x
}

// nestedConditions gère des conditions imbriquées avec early returns
//
// Params:
//   - a: premier nombre
//   - b: deuxième nombre
//
// Returns:
//   - int: résultat selon les conditions
func nestedConditions(a, b int) int {
	// Vérification si a positif
	if a > 0 {
		// Vérification si b positif
		if b > 0 {
			// Retour somme si les deux positifs
			return a + b
		}
		// Retour a si seul a positif
		return a
	}
	// Retour b si a non positif
	return b
}

// emptyIfBody teste les blocs if vides (edge case)
//
// Params:
//   - x: valeur à tester
func emptyIfBody(x int) {
	// Bloc if vide (ne déclenche pas KTN-FUNC-003)
	if x > 0 {
	}
	// Retour de la fonction
}

// init utilise les fonctions privées
func init() {
	// Appel de checkPositive
	_ = checkPositive(0)
	// Appel de processValue
	_ = processValue(0)
	// Appel de findMax
	_ = findMax(1, 0)
	// Appel de loopExample
	loopExample()
	// Appel de switchExample
	switchExample(0)
	// Appel de validateInput
	_ = validateInput("")
	// Appel de complexLogic
	_ = complexLogic(0)
	// Appel de nestedConditions
	_ = nestedConditions(1, 0)
	// Appel de emptyIfBody
	emptyIfBody(0)
}
