package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-004 : Constante sans type explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    TOUTES les constantes doivent avoir un type explicite spécifié.
//    Ne jamais laisser le compilateur inférer le type, même si c'est évident.
//
//    POURQUOI :
//    - Élimine l'ambiguïté (int ? int32 ? int64 ?)
//    - Rend le contrat explicite (importante pour APIs)
//    - Évite les surprises de conversion de types
//    - Facilite la relecture et la maintenance
//    - Standard pour code production
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer constants
//    // Ces constantes utilisent le type int explicite
//    const (
//        // MaxConnections définit le nombre maximum de connexions
//        MaxConnections int = 1000
//        // DefaultPort est le port par défaut
//        DefaultPort int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxConnectionsC004, DefaultPortC004, maxRetriesC004
// Connection limits
// Ces constantes définissent les limites de connexion
const (
	// MaxConnectionsC004 définit le nombre maximum de connexions simultanées
	MaxConnectionsC004 = 1000
	// DefaultPortC004 est le port par défaut de l'application
	DefaultPortC004 = 8080
	// maxRetriesC004 définit le nombre maximum de tentatives
	maxRetriesC004 = 3
)

// ❌ CAS INCORRECT 2 : Int32 sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxFileSizeC004, DefaultTimeoutC004, maxRequestsPerMinuteC004
// File and timeout settings
// Ces constantes définissent les limites de fichiers et timeouts
const (
	// MaxFileSizeC004 définit la taille maximale d'un fichier en octets
	MaxFileSizeC004 = 104857600
	// DefaultTimeoutC004 est le timeout par défaut en millisecondes
	DefaultTimeoutC004 = 30000
	// maxRequestsPerMinuteC004 définit le nombre maximum de requêtes par minute
	maxRequestsPerMinuteC004 = 1000
)

// ❌ CAS INCORRECT 3 : Uint8 sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxRetryAttemptsC004, DefaultQualityC004, minCompressionLevelC004
// Quality settings
// Ces constantes définissent les paramètres de qualité
const (
	// MaxRetryAttemptsC004 définit le nombre maximum de tentatives
	MaxRetryAttemptsC004 = 10
	// DefaultQualityC004 est la qualité par défaut (0-100)
	DefaultQualityC004 = 85
	// minCompressionLevelC004 est le niveau de compression minimum
	minCompressionLevelC004 = 1
)

// ❌ CAS INCORRECT 4 : Une seule constante sans type dans un groupe presque parfait (SEULE ERREUR : KTN-CONST-004 x1)
// NOTE : Groupe OK, commentaires OK, 2 constantes avec types, MAIS MaxTransactionIDC004 sans type
// ERREUR ATTENDUE : KTN-CONST-004 sur MaxTransactionIDC004 uniquement
// Transaction limits
// Ces constantes définissent les limites de transactions
const (
	// MaxMemoryBytesC004 définit la mémoire maximale en octets
	MaxMemoryBytesC004 uint64 = 17179869184
	// MaxTransactionIDC004 est l'ID de transaction maximum
	MaxTransactionIDC004 = 18446744073709551615
	// defaultCacheExpiryC004 est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiryC004 uint64 = 3600000000000
)
