// Package gospec_bad_functions montre des patterns de fonctions non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_functions

import "fmt"

// ❌ BAD PRACTICE: Long parameter list (should use struct)
func BadLongParams(name string, age int, email string, phone string,
	address string, city string, zip string, country string) {
	_, _, _, _, _, _, _, _ = name, age, email, phone, address, city, zip, country
}

// ❌ BAD PRACTICE: Using output parameters instead of return values
func BadOutputParams(x int, result *int) {
	*result = x * 2
	// En Go, préférer: func(x int) int { return x * 2 }
}

// ❌ BAD PRACTICE: Inconsistent parameter order across similar functions
func BadInconsistentOrder1(name string, id int) {}
func BadInconsistentOrder2(id int, name string) {} // Ordre différent

// ❌ BAD PRACTICE: Not using variadic parameters when appropriate
func BadNoVariadic(items []int) {
	for _, item := range items {
		fmt.Println(item)
	}
}
// Devrait être: func(items ...int)

// ❌ BAD PRACTICE: Function does too many things (violates SRP)
func BadDoEverything(data []byte) error {
	// Valide les données
	if len(data) == 0 {
		return fmt.Errorf("empty")
	}

	// Parse les données
	parsed := string(data)

	// Sauvegarde dans DB
	saveToDatabase(parsed)

	// Envoie notification
	sendNotification(parsed)

	// Log l'opération
	logOperation(parsed)

	// Update cache
	updateCache(parsed)

	return nil
	// Devrait être divisé en plusieurs fonctions
}

// ❌ BAD PRACTICE: Function modifies global state
var globalCounter int

func BadGlobalMutation() {
	globalCounter++ // Modifie état global, difficile à tester
}

// ❌ BAD PRACTICE: Function has side effects not obvious from signature
func BadHiddenSideEffect(name string) string {
	// Signature suggère fonction pure, mais...
	logToFile(name)        // Side effect: I/O
	updateMetrics(name)    // Side effect: modifie state
	sendAnalytics(name)    // Side effect: network call
	return "processed: " + name
}

// ❌ BAD PRACTICE: Boolean parameter instead of two functions
func BadBooleanParam(enable bool) {
	if enable {
		startService()
	} else {
		stopService()
	}
}
// Devrait être deux fonctions: EnableService() et DisableService()

// ❌ BAD PRACTICE: Returning interface{} when type can be known
func BadInterfaceReturn() interface{} {
	return 42 // Type connu, pas besoin d'interface{}
}
// Devrait être: func() int

// ❌ BAD PRACTICE: init() function doing too much
func init() {
	// Setup DB connection
	connectDB()
	// Load config
	loadConfig()
	// Start background workers
	startWorkers()
	// Register handlers
	registerHandlers()
	// init devrait être minimal
}

// ❌ BAD PRACTICE: Not using functional options pattern for complex config
type ServerConfig struct {
	Host    string
	Port    int
	Timeout int
	// ... many fields
}

func BadNewServer(host string, port int, timeout int, /* many params */) *Server {
	return &Server{host, port, timeout}
}
// Devrait utiliser functional options pattern

// ❌ BAD PRACTICE: Method receiver inconsistency (pointer vs value)
type Handler struct {
	count int
}

func (h Handler) Method1() {
	h.count++ // Modifie copie, pas l'original
}

func (h *Handler) Method2() {
	h.count++ // Modifie l'original
}
// Devrait être cohérent

// ❌ BAD PRACTICE: Using naked return without clear benefit
func BadNakedReturn(x int) (result int) {
	if x > 0 {
		result = x * 2
		return // Naked return peu clair ici
	}
	result = 0
	return
}

// ❌ BAD PRACTICE: Ignoring function return values
func BadIgnoringReturns() {
	getValue()        // Ignore return value
	getError()        // Ignore error
	getMultiple()     // Ignore multiple returns
}

// ❌ BAD PRACTICE: Function name doesn't describe what it does
func DoIt() {}
func Process() {}
func Handle() {}
func Execute() {}
// Noms trop vagues

// ❌ BAD PRACTICE: Exporting functions that should be private
func HelperFunction() {} // Exported mais devrait être private
func InternalUtil() {}   // "Internal" mais exporté

// ❌ BAD PRACTICE: Function with complex nested logic
func BadComplexNesting(x int) int {
	if x > 0 {
		if x < 100 {
			if x%2 == 0 {
				if x%3 == 0 {
					return x * 2
				} else {
					return x * 3
				}
			} else {
				return x + 1
			}
		} else {
			return x - 1
		}
	} else {
		return 0
	}
	// Devrait utiliser early returns
}

// ❌ BAD PRACTICE: Function does I/O but name doesn't indicate it
func GetUser(id int) string {
	// Nom suggère simple getter, mais...
	return fetchFromDatabase(id) // I/O hidden
}

// ❌ BAD PRACTICE: Function signature changes behavior based on nil
func BadNilBehavior(data []byte) {
	if data == nil {
		// Comportement complètement différent avec nil
		initializeDefaults()
		return
	}
	processData(data)
}
// Devrait être deux fonctions distinctes

// ❌ BAD PRACTICE: Method doesn't use receiver
func (h Handler) BadUnusedReceiver() {
	// Ne utilise pas 'h' - devrait être fonction standalone
	fmt.Println("hello")
}

// ❌ BAD PRACTICE: Constructor returns error instead of panicking for programmer errors
func NewBadValidator(pattern string) (*Validator, error) {
	if pattern == "" {
		return nil, fmt.Errorf("empty pattern") // Erreur programmeur
	}
	return &Validator{pattern}, nil
}
// Pour erreurs programmeur, panic est acceptable dans constructor

// Helper types and functions
type Server struct {
	host    string
	port    int
	timeout int
}

type Validator struct {
	pattern string
}

func saveToDatabase(string)   {}
func sendNotification(string) {}
func logOperation(string)     {}
func updateCache(string)      {}
func logToFile(string)        {}
func updateMetrics(string)    {}
func sendAnalytics(string)    {}
func startService()           {}
func stopService()            {}
func connectDB()              {}
func loadConfig()             {}
func startWorkers()           {}
func registerHandlers()       {}
func getValue() int           { return 0 }
func getError() error         { return nil }
func getMultiple() (int, error) { return 0, nil }
func fetchFromDatabase(int) string { return "" }
func initializeDefaults()          {}
func processData([]byte)           {}
