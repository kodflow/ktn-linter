// Good examples for the func011 test case.
package func011

// Constantes bien définies
const (
	DEFAULT_ITEM_COUNT  int     = 6
	MIN_LEGAL_AGE       int     = 18
	MAX_HUMAN_AGE       int     = 120
	DISCOUNT_RATE       float64 = 0.15
	MAX_RETRIES         int     = 3
	BUFFER_SIZE         int     = 1024
	HIGH_THRESHOLD      int     = 100
	DEFAULT_TIMEOUT_SEC int     = 30
	DEFAULT_PORT        int     = 8080
	ARRAY_SIZE          int     = 10
)

// processSixItems utilise une constante nommée
func processSixItems() {
	items := [DEFAULT_ITEM_COUNT]int{}
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
	return age >= MIN_LEGAL_AGE && age <= MAX_HUMAN_AGE
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
	return price * DISCOUNT_RATE
}

// processRetries utilise une constante nommée
func processRetries() {
	maxRetries := MAX_RETRIES
	_ = maxRetries
}

// setBufferSize utilise une constante nommée
func setBufferSize() {
	buffer := [BUFFER_SIZE]byte{}
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
	return value > HIGH_THRESHOLD
}

// waitTimeout utilise une constante nommée
func waitTimeout() {
	timeout := DEFAULT_TIMEOUT_SEC
	_ = timeout
}

// setPort utilise une constante nommée
func setPort() {
	port := DEFAULT_PORT
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

// arraySize tailles de tableaux doivent utiliser des constantes
func arraySize() {
	// arr tableau d'entiers de taille définie par constante
	var arr [ARRAY_SIZE]int
	_ = arr
}

// stringLiterals les littéraux string ne sont pas des magic numbers
func stringLiterals() {
	// Les strings ne déclenchent pas KTN-FUNC-011
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
	// Appel de arraySize
	arraySize()
	// Appel de stringLiterals
	stringLiterals()
}
