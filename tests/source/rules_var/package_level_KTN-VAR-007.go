package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-007 : Channel sans buffer size explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les channels doivent avoir le buffer size explicite dans le commentaire
//    ou préciser "unbuffered" si intentionnel.
//
//    POURQUOI :
//    - Clarté sur la sémantique (synchrone vs asynchrone)
//    - Aide à comprendre les performances attendues
//    - Évite les deadlocks non intentionnels
//    - Important pour la concurrence
//
// ✅ CAS PARFAIT (buffer size explicite) :
//
//    // Message channels
//    // Ces variables gèrent les messages inter-goroutines
//    var (
//        // MessageQueue canal de messages (buffer=100)
//        MessageQueue chan string = make(chan string, 100)
//        // DoneSignal signale la fin (unbuffered intentionnel)
//        DoneSignal chan bool = make(chan bool)
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Channels sans buffer info dans commentaire
// NOTE : Tout est parfait (groupe + commentaire groupe + commentaires individuels + types) SAUF buffer size manquant
// ERREURS ATTENDUES : KTN-VAR-007 UNIQUEMENT sur MessageQueueV007, doneSignalV007
// Channel variables
// Ces variables gèrent les channels de communication
var (
	// MessageQueueV007 est la file de messages
	MessageQueueV007 chan string = make(chan string)
	// ErrorQueueV007 est la file d'erreurs
	ErrorQueueV007 chan error
	// doneSignalV007 signale la fin
	doneSignalV007 chan bool = make(chan bool)
)
