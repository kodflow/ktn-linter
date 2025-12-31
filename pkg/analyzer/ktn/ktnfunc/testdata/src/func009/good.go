// Package func009 provides good test cases.
package func009

// Constantes bien définies
const (
	defaultItemCount  int     = 6
	minLegalAge       int     = 18
	maxHumanAge       int     = 120
	discountRate      float64 = 0.15
	maxRetries        int     = 3
	bufferSize        int     = 1024
	highThreshold     int     = 100
	defaultTimeoutSec int     = 30
	defaultPort       int     = 8080
	arraySize         int     = 10

	// Constantes avec iota (bitflags) - ne doivent pas déclencher KTN-FUNC-009
	flagNone  int = 0
	flagRead  int = 1 << iota
	flagWrite     // = 1 << 1
	flagExec      // = 1 << 2

	// Constantes avec iota simple
	levelDebug int = iota
	levelInfo
	levelWarn
	levelError
)

// processSixItems utilise une constante nommée
func processSixItems() {
	items := [defaultItemCount]int{}
	_ = items
}

// validateAge utilise des constantes nommées
//
// Params:
//   - age: âge à valider
//
// Returns:
//   - bool: true si l'âge est valide
func validateAge(age int) bool {
	// Retourne la validation de l'âge
	return age >= minLegalAge && age <= maxHumanAge
}

// calculateDiscount utilise une constante nommée
//
// Params:
//   - price: prix d'origine
//
// Returns:
//   - float64: montant de la réduction
func calculateDiscount(price float64) float64 {
	// Retourne le montant de la réduction
	return price * discountRate
}

// processRetries utilise une constante nommée
func processRetries() {
	retries := maxRetries
	_ = retries
}

// setBufferSize utilise une constante nommée
func setBufferSize() {
	buffer := [bufferSize]byte{}
	_ = buffer
}

// checkThreshold utilise une constante nommée
//
// Params:
//   - value: valeur à vérifier
//
// Returns:
//   - bool: true si la valeur dépasse le seuil
func checkThreshold(value int) bool {
	// Retourne true si la valeur dépasse le seuil
	return value > highThreshold
}

// waitTimeout utilise une constante nommée
func waitTimeout() {
	timeout := defaultTimeoutSec
	_ = timeout
}

// setPort utilise une constante nommée
func setPort() {
	port := defaultPort
	_ = port
}

// allowedNumbers utilise des nombres autorisés (0, 1, -1)
func allowedNumbers() {
	zero := 0
	one := 1
	minusOne := -1
	_ = zero
	_ = one
	_ = minusOne
}

// arraySizeFunc tailles de tableaux doivent utiliser des constantes
func arraySizeFunc() {
	// arr tableau d'entiers de taille définie par constante
	var arr [arraySize]int
	_ = arr
}

// stringLiterals les littéraux string ne sont pas des magic numbers
func stringLiterals() {
	// Les strings ne déclenchent pas KTN-FUNC-009
	message := "hello"
	code := "CODE123"
	_ = message
	_ = code
}

// init utilise les fonctions privées
func init() {
	// Appel de processSixItems
	processSixItems()
	// Appel de validateAge
	_ = validateAge(0)
	// Appel de calculateDiscount
	_ = calculateDiscount(0)
	// Appel de processRetries
	processRetries()
	// Appel de setBufferSize
	setBufferSize()
	// Appel de checkThreshold
	_ = checkThreshold(0)
	// Appel de waitTimeout
	waitTimeout()
	// Appel de setPort
	setPort()
	// Appel de allowedNumbers
	allowedNumbers()
	// Appel de arraySizeFunc
	arraySizeFunc()
	// Appel de stringLiterals
	stringLiterals()
}
