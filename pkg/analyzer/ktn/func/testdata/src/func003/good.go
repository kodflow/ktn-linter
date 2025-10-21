package func003

// Constantes bien définies
const (
	DEFAULT_ITEM_COUNT  = 6
	MIN_LEGAL_AGE       = 18
	MAX_HUMAN_AGE       = 120
	DISCOUNT_RATE       = 0.15
	MAX_RETRIES         = 3
	BUFFER_SIZE         = 1024
	HIGH_THRESHOLD      = 100
	DEFAULT_TIMEOUT_SEC = 30
	DEFAULT_PORT        = 8080
)

// processSixItems utilise une constante nommée
func processSixItems() {
	items := make([]int, DEFAULT_ITEM_COUNT)
	_ = items
}

// validateAge utilise des constantes nommées
func validateAge(age int) bool {
	return age >= MIN_LEGAL_AGE && age <= MAX_HUMAN_AGE
}

// calculateDiscount utilise une constante nommée
func calculateDiscount(price float64) float64 {
	return price * DISCOUNT_RATE
}

// processRetries utilise une constante nommée
func processRetries() {
	maxRetries := MAX_RETRIES
	_ = maxRetries
}

// setBufferSize utilise une constante nommée
func setBufferSize() {
	buffer := make([]byte, BUFFER_SIZE)
	_ = buffer
}

// checkThreshold utilise une constante nommée
func checkThreshold(value int) bool {
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

// arraySize tailles de tableaux sont acceptables
func arraySize() {
	var arr [10]int
	_ = arr
}
