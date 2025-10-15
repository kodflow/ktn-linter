package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-008 : Nom avec underscore (utiliser MixedCaps)
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les noms de variables doivent utiliser MixedCaps, pas underscore.
//    Go style : HTTPPort, maxRetries (pas HTTP_PORT, max_retries)
//
//    POURQUOI :
//    - Convention Go standard (Effective Go)
//    - Cohérence avec la stdlib Go
//    - Facilite la lecture (style uniforme)
//
// ✅ CAS PARFAIT (MixedCaps) :
//
//    // HTTP codes
//    // Ces variables contiennent les codes HTTP standards
//    var (
//        // HTTPOK représente le code 200
//        HTTPOK int = 200
//        // NotFound représente le code 404
//        NotFound int = 404
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Noms avec underscore (mais PAS en ALL_CAPS)
// ERREURS : KTN-VAR-008 UNIQUEMENT sur max_sizeV008, buffer_SizeV008, http_clientV008
// Buffer settings
var (
	// max_sizeV008 définit la taille maximale (snake_case)
	max_sizeV008 int = 1024
	// buffer_SizeV008 définit la taille du buffer (mixte avec underscore)
	buffer_SizeV008 int = 512
	// http_clientV008 est le client HTTP (snake_case)
	http_clientV008 string = "default"
)
