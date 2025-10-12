package rules_const

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 1 : Constantes déclarées individuellement sans regroupement
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Les constantes sont déclarées une par une avec "const X = ..."
//    au lieu d'être regroupées dans un bloc const ().
//
// 📋 CE QU'ON A :
//    const EnableFeatureX bool = true
//    const EnableDebug bool = false
//    const isProduction bool = true
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Boolean configuration
//    // Ces constantes configurent les fonctionnalités
//    const (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//        // isProduction indique l'environnement production
//        isProduction bool = true
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur EnableFeatureX (pas de regroupement)
//    - KTN-CONST-001 sur EnableDebug (pas de regroupement)
//    - KTN-CONST-001 sur isProduction (pas de regroupement)
// ════════════════════════════════════════════════════════════════════════════
const EnableFeatureX bool = true
const EnableDebug bool = false
const isProduction bool = true

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 2 : Groupe sans commentaire ET constantes sans commentaires individuels
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le bloc const () existe mais :
//    1. Il n'y a pas de commentaire de groupe avant const ()
//    2. Aucune constante n'a de commentaire individuel
//
// 📋 CE QU'ON A :
//    const (
//        ApplicationName string = "MyApp"
//        Version         string = "1.0.0"
//        defaultEncoding string = "UTF-8"
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Application metadata
//    // Ces constantes définissent les métadonnées de l'application
//    const (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//        // defaultEncoding est l'encodage par défaut
//        defaultEncoding string = "UTF-8"
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-002 sur le groupe (pas de commentaire de groupe)
//    - KTN-CONST-003 sur ApplicationName (pas de commentaire individuel)
//    - KTN-CONST-003 sur Version (pas de commentaire individuel)
//    - KTN-CONST-003 sur defaultEncoding (pas de commentaire individuel)
// ════════════════════════════════════════════════════════════════════════════
const (
	ApplicationName string = "MyApp"
	Version         string = "1.0.0"
	defaultEncoding string = "UTF-8"
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 3 : Types non explicites (inférés par le compilateur)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Les constantes n'ont pas de type explicite.
//    Le compilateur infère "int" mais ce n'est pas clair.
//    Le groupe a un commentaire mais les constantes ont des commentaires individuels.
//
// 📋 CE QU'ON A :
//    // Ces constantes n'ont pas de type explicite
//    const (
//        // MaxConnections définit...
//        MaxConnections = 1000     // Type inféré : int
//        // DefaultPort est...
//        DefaultPort = 8080        // Type inféré : int
//        // maxRetries définit...
//        maxRetries = 3            // Type inféré : int
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Integer constants
//    // Ces constantes utilisent le type int explicite
//    const (
//        // MaxConnections définit le nombre maximum de connexions
//        MaxConnections int = 1000
//        // DefaultPort est le port par défaut
//        DefaultPort int = 8080
//        // maxRetries définit le nombre de tentatives
//        maxRetries int = 3
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-004 sur MaxConnections (type manquant)
//    - KTN-CONST-004 sur DefaultPort (type manquant)
//    - KTN-CONST-004 sur maxRetries (type manquant)
// ════════════════════════════════════════════════════════════════════════════
// Ces constantes n'ont pas de type explicite
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 4 : Commentaire de groupe OK mais pas de commentaires individuels
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le groupe a un commentaire, les types sont explicites, MAIS
//    aucune constante n'a de commentaire individuel.
//
// 📋 CE QU'ON A :
//    // Ces constantes utilisent des entiers 8 bits (-128 à 127)
//    const (
//        MinAge          int8 = 18
//        MaxAge          int8 = 120
//        defaultPriority int8 = 5
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Age limits
//    // Ces constantes définissent les limites d'âge
//    const (
//        // MinAge est l'âge minimum requis
//        MinAge int8 = 18
//        // MaxAge est l'âge maximum accepté
//        MaxAge int8 = 120
//        // defaultPriority est la priorité par défaut
//        defaultPriority int8 = 5
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-003 sur MinAge (pas de commentaire individuel)
//    - KTN-CONST-003 sur MaxAge (pas de commentaire individuel)
//    - KTN-CONST-003 sur defaultPriority (pas de commentaire individuel)
// ════════════════════════════════════════════════════════════════════════════
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 5 : Pas de regroupement ET pas de commentaire
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de plusieurs erreurs :
//    1. Constantes déclarées individuellement (pas de const ())
//    2. Aucun commentaire sur aucune constante
//
// 📋 CE QU'ON A :
//    const MaxQueueSize int16 = 10000
//    const DefaultBufferSize int16 = 4096
//    const minCacheSize int16 = 512
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Queue configuration
//    // Ces constantes configurent les tailles de queue
//    const (
//        // MaxQueueSize est la taille maximale de la queue
//        MaxQueueSize int16 = 10000
//        // DefaultBufferSize est la taille du buffer par défaut
//        DefaultBufferSize int16 = 4096
//        // minCacheSize est la taille minimale du cache
//        minCacheSize int16 = 512
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur MaxQueueSize (pas de regroupement)
//    - KTN-CONST-001 sur DefaultBufferSize (pas de regroupement)
//    - KTN-CONST-001 sur minCacheSize (pas de regroupement)
// ════════════════════════════════════════════════════════════════════════════
const MaxQueueSize int16 = 10000
const DefaultBufferSize int16 = 4096
const minCacheSize int16 = 512

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 6 : Type non explicite ET pas de commentaire individuel
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de deux erreurs :
//    1. Le groupe a un commentaire mais les constantes n'en ont pas
//    2. Les types sont inférés au lieu d'être explicites
//
// 📋 CE QU'ON A :
//    // Integer 32-bit constants
//    const (
//        MaxFileSize          = 104857600
//        DefaultTimeout       = 30000
//        maxRequestsPerMinute = 1000
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // File and timeout limits
//    // Ces constantes définissent les limites de fichier et timeout
//    const (
//        // MaxFileSize est la taille maximale d'un fichier
//        MaxFileSize int32 = 104857600
//        // DefaultTimeout est le timeout par défaut
//        DefaultTimeout int32 = 30000
//        // maxRequestsPerMinute limite les requêtes par minute
//        maxRequestsPerMinute int32 = 1000
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-003 sur MaxFileSize (pas de commentaire individuel)
//    - KTN-CONST-004 sur MaxFileSize (type manquant)
//    - KTN-CONST-003 sur DefaultTimeout (pas de commentaire individuel)
//    - KTN-CONST-004 sur DefaultTimeout (type manquant)
//    - KTN-CONST-003 sur maxRequestsPerMinute (pas de commentaire individuel)
//    - KTN-CONST-004 sur maxRequestsPerMinute (type manquant)
// ════════════════════════════════════════════════════════════════════════════
// Integer 32-bit constants
const (
	MaxFileSize          = 104857600
	DefaultTimeout       = 30000
	maxRequestsPerMinute = 1000
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 7 : Pas de commentaire du tout (ni groupe ni individuel)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le bloc const () existe et les types sont explicites, MAIS
//    il n'y a AUCUN commentaire (ni groupe, ni individuel).
//
// 📋 CE QU'ON A :
//    const (
//        MaxDiskSpace   int64 = 1099511627776
//        UnixEpoch      int64 = 0
//        nanosPerSecond int64 = 1000000000
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Disk and time constants
//    // Ces constantes utilisent int64 pour les grandes valeurs
//    const (
//        // MaxDiskSpace est l'espace disque maximum en octets
//        MaxDiskSpace int64 = 1099511627776
//        // UnixEpoch représente le timestamp Unix epoch
//        UnixEpoch int64 = 0
//        // nanosPerSecond est le nombre de nanosecondes par seconde
//        nanosPerSecond int64 = 1000000000
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-002 sur le groupe (pas de commentaire de groupe)
//    - KTN-CONST-003 sur MaxDiskSpace (pas de commentaire individuel)
//    - KTN-CONST-003 sur UnixEpoch (pas de commentaire individuel)
//    - KTN-CONST-003 sur nanosPerSecond (pas de commentaire individuel)
// ════════════════════════════════════════════════════════════════════════════
const (
	MaxDiskSpace   int64 = 1099511627776
	UnixEpoch      int64 = 0
	nanosPerSecond int64 = 1000000000
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 8 : Pas de regroupement ET type non explicite
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de deux erreurs :
//    1. Constantes déclarées individuellement
//    2. Types inférés au lieu d'être explicites
//
// 📋 CE QU'ON A :
//    const MaxUserID = 4294967295
//    const DefaultPoolSize = 100
//    const minWorkers = 4
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // User and pool limits
//    // Ces constantes définissent les limites utilisateur et pool
//    const (
//        // MaxUserID est l'ID utilisateur maximum
//        MaxUserID uint = 4294967295
//        // DefaultPoolSize est la taille par défaut du pool
//        DefaultPoolSize uint = 100
//        // minWorkers est le nombre minimum de workers
//        minWorkers uint = 4
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur MaxUserID (pas de regroupement)
//    - KTN-CONST-004 sur MaxUserID (type manquant)
//    - KTN-CONST-001 sur DefaultPoolSize (pas de regroupement)
//    - KTN-CONST-004 sur DefaultPoolSize (type manquant)
//    - KTN-CONST-001 sur minWorkers (pas de regroupement)
//    - KTN-CONST-004 sur minWorkers (type manquant)
// ════════════════════════════════════════════════════════════════════════════
const MaxUserID = 4294967295
const DefaultPoolSize = 100
const minWorkers = 4

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 9 : Commentaires OK mais types non explicites
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le bloc est bien structuré avec commentaires groupe + individuels,
//    MAIS les types sont inférés au lieu d'être explicites.
//
// 📋 CE QU'ON A :
//    // Unsigned 8-bit constants
//    const (
//        // MaxRetryAttempts définit...
//        MaxRetryAttempts = 10
//        // DefaultQuality est...
//        DefaultQuality = 85
//        // minCompressionLevel est...
//        minCompressionLevel = 1
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Unsigned 8-bit constants
//    // Ces constantes utilisent uint8 (0 à 255)
//    const (
//        // MaxRetryAttempts définit le nombre maximum de tentatives
//        MaxRetryAttempts uint8 = 10
//        // DefaultQuality est la qualité par défaut (0-100)
//        DefaultQuality uint8 = 85
//        // minCompressionLevel est le niveau de compression minimum
//        minCompressionLevel uint8 = 1
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-004 sur MaxRetryAttempts (type manquant)
//    - KTN-CONST-004 sur DefaultQuality (type manquant)
//    - KTN-CONST-004 sur minCompressionLevel (type manquant)
// ════════════════════════════════════════════════════════════════════════════
// Unsigned 8-bit constants
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel = 1
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 10 : Mélange - certaines avec commentaire, d'autres sans
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le groupe a un commentaire, les types sont explicites, MAIS
//    certaines constantes ont des commentaires individuels et d'autres non.
//    TOUTES les constantes doivent avoir leur commentaire individuel.
//
// 📋 CE QU'ON A :
//    // Unsigned 16-bit constants
//    const (
//        // HTTPPort est le port HTTP standard
//        HTTPPort   uint16 = 80
//        HTTPSPort  uint16 = 443          // Pas de commentaire !
//        customPort uint16 = 3000         // Pas de commentaire !
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // HTTP ports
//    // Ces constantes définissent les ports HTTP standards
//    const (
//        // HTTPPort est le port HTTP standard
//        HTTPPort uint16 = 80
//        // HTTPSPort est le port HTTPS standard
//        HTTPSPort uint16 = 443
//        // customPort est un port personnalisé
//        customPort uint16 = 3000
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-003 sur HTTPSPort (pas de commentaire individuel)
//    - KTN-CONST-003 sur customPort (pas de commentaire individuel)
// ════════════════════════════════════════════════════════════════════════════
// Unsigned 16-bit constants
const (
	// HTTPPort est le port HTTP standard
	HTTPPort   uint16 = 80
	HTTPSPort  uint16 = 443
	customPort uint16 = 3000
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 11 : Pas de regroupement + commentaires individuels + type manquant
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Les constantes ont des commentaires individuels (bien !) mais :
//    1. Elles ne sont pas regroupées dans const ()
//    2. Les types sont inférés au lieu d'être explicites
//
// 📋 CE QU'ON A :
//    // MaxRecordCount définit...
//    const MaxRecordCount = 1000000
//
//    // DefaultChunkSize est...
//    const DefaultChunkSize = 65536
//    const minBatchSize = 100
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Record processing
//    // Ces constantes configurent le traitement des enregistrements
//    const (
//        // MaxRecordCount définit le nombre maximum d'enregistrements
//        MaxRecordCount uint32 = 1000000
//        // DefaultChunkSize est la taille par défaut d'un chunk
//        DefaultChunkSize uint32 = 65536
//        // minBatchSize est la taille minimale d'un batch
//        minBatchSize uint32 = 100
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur MaxRecordCount (pas de regroupement)
//    - KTN-CONST-004 sur MaxRecordCount (type manquant)
//    - KTN-CONST-001 sur DefaultChunkSize (pas de regroupement)
//    - KTN-CONST-004 sur DefaultChunkSize (type manquant)
//    - KTN-CONST-001 sur minBatchSize (pas de regroupement)
//    - KTN-CONST-004 sur minBatchSize (type manquant)
// ════════════════════════════════════════════════════════════════════════════
// MaxRecordCount définit le nombre maximum d'enregistrements
const MaxRecordCount = 1000000

// DefaultChunkSize est la taille par défaut d'un chunk en octets
const DefaultChunkSize = 65536
const minBatchSize = 100

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 12 : Une seule constante sans type dans un groupe parfait
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le groupe est presque parfait : commentaire de groupe, commentaires
//    individuels, types explicites... SAUF une constante qui manque son type.
//
// 📋 CE QU'ON A :
//    // Unsigned 64-bit constants
//    const (
//        // MaxMemoryBytes définit...
//        MaxMemoryBytes uint64 = 17179869184
//        // MaxTransactionID est...
//        MaxTransactionID = 18446744073709551615    // ❌ Type manquant !
//        // defaultCacheExpiry est...
//        defaultCacheExpiry uint64 = 3600000000000
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Unsigned 64-bit constants
//    // Ces constantes utilisent uint64 pour de très grandes valeurs
//    const (
//        // MaxMemoryBytes définit la mémoire maximale en octets
//        MaxMemoryBytes uint64 = 17179869184
//        // MaxTransactionID est l'ID de transaction maximum
//        MaxTransactionID uint64 = 18446744073709551615
//        // defaultCacheExpiry est le délai d'expiration du cache
//        defaultCacheExpiry uint64 = 3600000000000
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-004 sur MaxTransactionID (type manquant)
// ════════════════════════════════════════════════════════════════════════════
// Unsigned 64-bit constants
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID = 18446744073709551615 // Type manquant !
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 13 : Pas de commentaire de groupe (mais commentaires individuels OK)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le bloc const () existe, les types sont explicites, les commentaires
//    individuels sont présents (ligne de fin), MAIS pas de commentaire de groupe.
//
// 📋 CE QU'ON A :
//    const (
//        NullByte    byte = 0x00  // commentaire de fin
//        NewlineByte byte = 0x0A  // commentaire de fin
//        tabByte     byte = 0x09  // commentaire de fin
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Byte constants
//    // Ces constantes représentent des octets individuels
//    const (
//        // NullByte représente l'octet null
//        NullByte byte = 0x00
//        // NewlineByte représente le caractère newline
//        NewlineByte byte = 0x0A
//        // tabByte représente le caractère tabulation
//        tabByte byte = 0x09
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-002 sur le groupe (pas de commentaire de groupe)
//    - KTN-CONST-003 sur NullByte (commentaire doit être au-dessus, pas à côté)
//    - KTN-CONST-003 sur NewlineByte (commentaire doit être au-dessus)
//    - KTN-CONST-003 sur tabByte (commentaire doit être au-dessus)
//
// 📝 NOTE : Les commentaires de fin de ligne ne sont pas acceptés, les
//          commentaires doivent être sur la ligne au-dessus de la constante.
// ════════════════════════════════════════════════════════════════════════════
const (
	NullByte    byte = 0x00
	NewlineByte byte = 0x0A
	tabByte     byte = 0x09
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 14 : Constantes individuelles sans type (rune)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de plusieurs erreurs :
//    1. Constantes déclarées individuellement (pas de const ())
//    2. Types inférés (rune) au lieu d'être explicites
//    3. Pas de commentaires
//
// 📋 CE QU'ON A :
//    const SpaceRune = ' '
//    const NewlineRune = '\n'
//    const heartEmoji = '❤'
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Rune constants
//    // Ces constantes représentent des caractères Unicode
//    const (
//        // SpaceRune représente le caractère espace
//        SpaceRune rune = ' '
//        // NewlineRune représente le caractère retour à la ligne
//        NewlineRune rune = '\n'
//        // heartEmoji représente l'emoji cœur
//        heartEmoji rune = '❤'
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur SpaceRune (pas de regroupement)
//    - KTN-CONST-004 sur SpaceRune (type manquant)
//    - KTN-CONST-001 sur NewlineRune (pas de regroupement)
//    - KTN-CONST-004 sur NewlineRune (type manquant)
//    - KTN-CONST-001 sur heartEmoji (pas de regroupement)
//    - KTN-CONST-004 sur heartEmoji (type manquant)
// ════════════════════════════════════════════════════════════════════════════
const SpaceRune = ' '
const NewlineRune = '\n'
const heartEmoji = '❤'

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 15 : Groupe OK, commentaire de groupe OK, mais pas de commentaires individuels
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Le bloc const () existe, le commentaire de groupe existe, les types
//    sont explicites, MAIS aucune constante n'a de commentaire individuel.
//
// 📋 CE QU'ON A :
//    // Float32 constants
//    const (
//        Pi32         float32 = 3.14159265
//        DefaultRate  float32 = 1.5
//        minThreshold float32 = 0.01
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Float32 constants
//    // Ces constantes utilisent float32 pour la précision simple
//    const (
//        // Pi32 est une approximation de Pi en float32
//        Pi32 float32 = 3.14159265
//        // DefaultRate est le taux par défaut
//        DefaultRate float32 = 1.5
//        // minThreshold est le seuil minimum
//        minThreshold float32 = 0.01
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-003 sur Pi32 (pas de commentaire individuel)
//    - KTN-CONST-003 sur DefaultRate (pas de commentaire individuel)
//    - KTN-CONST-003 sur minThreshold (pas de commentaire individuel)
// ════════════════════════════════════════════════════════════════════════════
// Float32 constants
const (
	Pi32         float32 = 3.14159265
	DefaultRate  float32 = 1.5
	minThreshold float32 = 0.01
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 16 : Pas de regroupement du tout (float64)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Les constantes sont déclarées individuellement au lieu d'être regroupées.
//    Elles sont liées thématiquement (constantes mathématiques) mais séparées.
//
// 📋 CE QU'ON A :
//    const Pi float64 = 3.14159265358979323846
//    const EulerNumber float64 = 2.71828182845904523536
//    const goldenRatio float64 = 1.618033988749894848204586
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Mathematical constants
//    // Ces constantes représentent des valeurs mathématiques importantes
//    const (
//        // Pi est une approximation de Pi en haute précision
//        Pi float64 = 3.14159265358979323846
//        // EulerNumber est le nombre d'Euler (e)
//        EulerNumber float64 = 2.71828182845904523536
//        // goldenRatio est le nombre d'or (phi)
//        goldenRatio float64 = 1.618033988749894848204586
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur Pi (pas de regroupement)
//    - KTN-CONST-001 sur EulerNumber (pas de regroupement)
//    - KTN-CONST-001 sur goldenRatio (pas de regroupement)
// ════════════════════════════════════════════════════════════════════════════
const Pi float64 = 3.14159265358979323846
const EulerNumber float64 = 2.71828182845904523536
const goldenRatio float64 = 1.618033988749894848204586

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 17 : Groupe sans commentaire de groupe ET types manquants
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de deux erreurs :
//    1. Pas de commentaire de groupe
//    2. Types inférés (complex64) au lieu d'être explicites
//    3. Pas de commentaires individuels
//
// 📋 CE QU'ON A :
//    const (
//        ImaginaryUnit64 = 0 + 1i
//        ComplexZero64   = 0 + 0i
//        sampleComplex64 = 3.5 + 2.8i
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Complex64 constants
//    // Ces constantes représentent des nombres complexes (float32+float32)
//    const (
//        // ImaginaryUnit64 représente l'unité imaginaire i
//        ImaginaryUnit64 complex64 = 0 + 1i
//        // ComplexZero64 représente zéro en complex64
//        ComplexZero64 complex64 = 0 + 0i
//        // sampleComplex64 est un exemple de nombre complexe
//        sampleComplex64 complex64 = 3.5 + 2.8i
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-002 sur le groupe (pas de commentaire de groupe)
//    - KTN-CONST-003 sur ImaginaryUnit64 (pas de commentaire individuel)
//    - KTN-CONST-004 sur ImaginaryUnit64 (type manquant)
//    - KTN-CONST-003 sur ComplexZero64 (pas de commentaire individuel)
//    - KTN-CONST-004 sur ComplexZero64 (type manquant)
//    - KTN-CONST-003 sur sampleComplex64 (pas de commentaire individuel)
//    - KTN-CONST-004 sur sampleComplex64 (type manquant)
// ════════════════════════════════════════════════════════════════════════════
const (
	ImaginaryUnit64 = 0 + 1i
	ComplexZero64   = 0 + 0i
	sampleComplex64 = 3.5 + 2.8i
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 18 : Un seul commentaire pour plusieurs constantes (pas individuel)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Il y a UN commentaire au début qui sert de commentaire de groupe,
//    mais les constantes individuelles n'ont pas chacune leur commentaire.
//    Le commentaire partagé est insuffisant pour documenter chaque constante.
//
// 📋 CE QU'ON A :
//    const (
//        // Ces constantes représentent des valeurs complexes
//        ImaginaryUnit     complex128 = 0 + 1i
//        ComplexZero       complex128 = 0 + 0i
//        eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
//    )
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Complex128 constants
//    // Ces constantes représentent des nombres complexes haute précision
//    const (
//        // ImaginaryUnit représente l'unité imaginaire i
//        ImaginaryUnit complex128 = 0 + 1i
//        // ComplexZero représente zéro en complex128
//        ComplexZero complex128 = 0 + 0i
//        // eulerIdentityBase est la base pour l'identité d'Euler
//        eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-003 sur ComplexZero (pas de commentaire individuel)
//    - KTN-CONST-003 sur eulerIdentityBase (pas de commentaire individuel)
//
// 📝 NOTE : ImaginaryUnit n'est pas détecté car le commentaire ligne 132
//          est attaché à la fois au groupe ET à la première constante.
//          Le linter le considère comme un commentaire partagé.
// ════════════════════════════════════════════════════════════════════════════
const (
	// Ces constantes représentent des valeurs complexes
	ImaginaryUnit     complex128 = 0 + 1i
	ComplexZero       complex128 = 0 + 0i
	eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
)

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 19 : Constante orpheline (aucune règle respectée)
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Cumul de TOUTES les erreurs possibles :
//    1. Déclarée individuellement (pas de const ())
//    2. Pas de commentaire
//    3. Type inféré au lieu d'être explicite
//
// 📋 CE QU'ON A :
//    const orphanConst = 42
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Answer constant
//    // Cette constante représente la réponse universelle
//    const (
//        // orphanConst est la réponse à la grande question
//        orphanConst int = 42
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur orphanConst (pas de regroupement)
//    - KTN-CONST-004 sur orphanConst (type manquant)
//
// 📝 NOTE : Pas d'erreur KTN-CONST-003 car pour les constantes individuelles,
//          on ne vérifie que le regroupement et le type.
// ════════════════════════════════════════════════════════════════════════════
const orphanConst = 42

// ════════════════════════════════════════════════════════════════════════════
// SCÉNARIO 20 : Mélange de const groupées et non groupées dans le même thème
// ════════════════════════════════════════════════════════════════════════════
// ❌ PROBLÈME :
//    Des constantes du même thème (theme configuration) sont séparées :
//    - Certaines dans un bloc const ()
//    - D'autres déclarées individuellement
//    Toutes devraient être dans le MÊME bloc const ().
//
// 📋 CE QU'ON A :
//    // Configuration theme - Partie 1
//    const (
//        // ThemeLight est le thème clair
//        ThemeLight string = "light"
//        // ThemeDark est le thème sombre
//        ThemeDark string = "dark"
//    )
//
//    // Configuration theme - Partie 2 (devrait être dans le même groupe)
//    const ThemeAuto string = "auto"
//    const ThemeCustom string = "custom"
//
// ✅ CE QU'ON DEVRAIT AVOIR :
//    // Theme configuration
//    // Ces constantes définissent les thèmes disponibles
//    const (
//        // ThemeLight est l'identifiant du thème clair
//        ThemeLight string = "light"
//        // ThemeDark est l'identifiant du thème sombre
//        ThemeDark string = "dark"
//        // ThemeAuto est l'identifiant du thème automatique
//        ThemeAuto string = "auto"
//        // ThemeCustom est l'identifiant du thème personnalisé
//        ThemeCustom string = "custom"
//    )
//
// 🔍 ERREURS DÉTECTÉES :
//    - KTN-CONST-001 sur ThemeAuto (pas de regroupement)
//    - KTN-CONST-001 sur ThemeCustom (pas de regroupement)
// ════════════════════════════════════════════════════════════════════════════
// Configuration theme - Partie 1
const (
	// ThemeLight est le thème clair
	ThemeLight string = "light"
	// ThemeDark est le thème sombre
	ThemeDark string = "dark"
)

// Configuration theme - Partie 2 (devrait être dans le même groupe)
const ThemeAuto string = "auto"
const ThemeCustom string = "custom"
