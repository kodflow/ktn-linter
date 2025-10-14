package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-007 : Channel avec buffer size explicite
// ════════════════════════════════════════════════════════════════════════════

// Channel variables
// Ces variables sont des channels pour la communication inter-goroutines
var (
	// MessageQueueV007Good est le channel pour les messages (buffer=100)
	MessageQueueV007Good chan string = make(chan string, 100)
	// ErrorQueueV007Good est le channel pour les erreurs (buffer=50)
	ErrorQueueV007Good chan error = make(chan error, 50)
	// doneSignalV007Good signale la fin d'exécution (unbuffered intentionnel)
	doneSignalV007Good chan bool = make(chan bool)
)
