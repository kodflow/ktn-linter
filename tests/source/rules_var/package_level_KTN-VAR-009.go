package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-009 : Nom en ALL_CAPS (utiliser MixedCaps)
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les noms en ALL_CAPS sont réservés aux constantes dans d'autres langages
//    (C, Java) mais Go utilise MixedCaps pour tout.
//
//    POURQUOI :
//    - Convention Go standard (pas ALL_CAPS)
//    - Évite confusion avec conventions d'autres langages
//    - MixedCaps est le style unifié Go
//
// ✅ CAS PARFAIT (MixedCaps) :
//
//    // Buffer configuration
//    // Cette variable configure la taille du buffer
//    var (
//        // MaxBufferSize est la taille maximale du buffer
//        MaxBufferSize int = 1024
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Noms en ALL_CAPS (sans underscore)
// ERREURS : KTN-VAR-009 UNIQUEMENT sur MAXSIZEV009, TIMEOUTV009, BUFFERSIZEV009
// Size configuration
var (
	// MAXSIZEV009 est la taille maximale (ALL_CAPS sans underscore)
	MAXSIZEV009 int = 1024
	// TIMEOUTV009 est le timeout par défaut (ALL_CAPS sans underscore)
	TIMEOUTV009 int = 30
	// BUFFERSIZEV009 est la taille du buffer (ALL_CAPS sans underscore)
	BUFFERSIZEV009 int = 512
)

// ❌ ERREUR : Cumul VAR-008 + VAR-009 (underscore ET ALL_CAPS)
// ERREURS : KTN-VAR-008 + KTN-VAR-009 sur MAX_BUFFER_SIZEV009, DEFAULT_TIMEOUTV009
// Buffer configuration
var (
	// MAX_BUFFER_SIZEV009 viole les deux règles (underscore + ALL_CAPS)
	MAX_BUFFER_SIZEV009 int = 2048
	// DEFAULT_TIMEOUTV009 viole les deux règles (underscore + ALL_CAPS)
	DEFAULT_TIMEOUTV009 int = 60
)
