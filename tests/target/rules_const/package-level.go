package rules_const

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : BOOL
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Regroupement dans const ()
//    ✓ Commentaire de groupe (2 lignes : titre + description)
//    ✓ Chaque constante a son commentaire individuel
//    ✓ Type bool explicite pour toutes
//    ✓ Naming MixedCaps (publiques et privées)
//
// 💡 PATTERN :
//    // Titre du groupe
//    // Description détaillée du groupe
//    const (
//        // Commentaire constante publique
//        PublicConst bool = true
//        // Commentaire constante privée
//        privateConst bool = false
//    )
// ════════════════════════════════════════════════════════════════════════════

// Boolean configuration
// Ces constantes configurent les fonctionnalités de l'application
const (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : STRING
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type string explicite (même si Go pourrait l'inférer)
//    ✓ Commentaires décrivent le rôle, pas juste le nom
//    ✓ Cohésion thématique (métadonnées ensemble, thèmes ensemble)
// ════════════════════════════════════════════════════════════════════════════

// Application metadata
// Ces constantes contiennent les métadonnées de l'application
const (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle de l'application
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut utilisé
	defaultEncoding string = "UTF-8"
)

// Theme configuration
// Ces constantes définissent les thèmes disponibles dans l'interface
const (
	// ThemeLight est l'identifiant du thème clair
	ThemeLight string = "light"
	// ThemeDark est l'identifiant du thème sombre
	ThemeDark string = "dark"
	// ThemeAuto est l'identifiant du thème automatique
	ThemeAuto string = "auto"
	// ThemeCustom est l'identifiant du thème personnalisé
	ThemeCustom string = "custom"
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : INT (taille dépend de l'architecture: 32 ou 64 bits)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int explicite
//    ✓ Mix de constantes publiques et privées
//    ✓ Commentaires précis sur le rôle de chaque constante
//
// ⚠️  IMPORTANT :
//    Le type "int" peut être 32 ou 64 bits selon l'architecture.
//    Si la taille est critique, utiliser int32 ou int64 explicitement.
// ════════════════════════════════════════════════════════════════════════════

// Integer constants (int)
// Ces constantes utilisent le type int (taille dépend de l'architecture)
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : INT8 (valeurs -128 à 127)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int8 explicite
//    ✓ Choix justifié : valeurs garanties dans la plage -128 à 127
//    ✓ Approprié pour âges, priorités, petites valeurs
// ════════════════════════════════════════════════════════════════════════════

// Integer 8-bit constants
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	// MinAge est l'âge minimum requis
	MinAge int8 = 18
	// MaxAge est l'âge maximum accepté
	MaxAge int8 = 120
	// defaultPriority est la priorité par défaut
	defaultPriority int8 = 5
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : INT16 (valeurs -32768 à 32767)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int16 explicite
//    ✓ Approprié pour tailles de queue, buffers, caches
//    ✓ Commentaires mentionnent l'unité ou le contexte
// ════════════════════════════════════════════════════════════════════════════

// Integer 16-bit constants
// Ces constantes utilisent des entiers 16 bits (-32768 à 32767)
const (
	// MaxQueueSize définit la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille du buffer par défaut
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : INT32 (valeurs -2147483648 à 2147483647)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int32 explicite
//    ✓ Approprié pour tailles de fichiers, timeouts, compteurs
//    ✓ Commentaires incluent les unités (octets, millisecondes)
// ════════════════════════════════════════════════════════════════════════════

// Integer 32-bit constants
// Ces constantes utilisent des entiers 32 bits (-2147483648 à 2147483647)
const (
	// MaxFileSize définit la taille maximale d'un fichier en octets
	MaxFileSize int32 = 104857600 // 100 MB
	// DefaultTimeout est le timeout par défaut en millisecondes
	DefaultTimeout int32 = 30000
	// maxRequestsPerMinute limite le nombre de requêtes par minute
	maxRequestsPerMinute int32 = 1000
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : INT64 (très grandes valeurs)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int64 explicite
//    ✓ Approprié pour espace disque, timestamps, nanosecondes
//    ✓ Valeur 0 documentée comme "intentionnelle" (UnixEpoch)
//    ✓ Commentaires explicites avec unités
// ════════════════════════════════════════════════════════════════════════════

// Integer 64-bit constants
// Ces constantes utilisent des entiers 64 bits (très grandes valeurs)
const (
	// MaxDiskSpace définit l'espace disque maximum en octets
	MaxDiskSpace int64 = 1099511627776 // 1 TB
	// UnixEpoch représente le timestamp Unix de référence
	UnixEpoch int64 = 0
	// nanosPerSecond est le nombre de nanosecondes dans une seconde
	nanosPerSecond int64 = 1000000000
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : UINT (taille dépend de l'architecture: 32 ou 64 bits)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint explicite
//    ✓ Approprié pour valeurs garanties positives (IDs, compteurs, tailles)
//    ✓ Double la plage positive vs int de même taille
// ════════════════════════════════════════════════════════════════════════════

// Unsigned integer constants (uint)
// Ces constantes utilisent des entiers non signés (taille dépend de l'architecture)
const (
	// MaxUserID est l'ID utilisateur maximum
	MaxUserID uint = 4294967295
	// DefaultPoolSize est la taille par défaut du pool
	DefaultPoolSize uint = 100
	// minWorkers est le nombre minimum de workers
	minWorkers uint = 4
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : UINT8 (valeurs 0 à 255)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint8 explicite
//    ✓ Approprié pour qualité (0-100), niveaux (0-10), pourcentages
//    ✓ Commentaires mentionnent les plages valides quand pertinent
//
// 📝 NOTE : uint8 et byte sont équivalents. Utiliser byte pour données binaires,
//           uint8 pour nombres.
// ════════════════════════════════════════════════════════════════════════════

// Unsigned 8-bit constants
// Ces constantes utilisent des entiers non signés 8 bits (0 à 255)
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts uint8 = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality uint8 = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel uint8 = 1
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : UINT16 (valeurs 0 à 65535)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint16 explicite
//    ✓ Parfait pour ports réseau (plage 0-65535)
//    ✓ Naming avec initialismes corrects (HTTPPort, HTTPSPort)
//
// 📝 NAMING CONVENTION :
//    - HTTPPort (pas HttpPort, ni HTTP_PORT)
//    - URLMaxLength (pas UrlMaxLength)
//    - Les initialismes restent en majuscules dans MixedCaps
// ════════════════════════════════════════════════════════════════════════════

// Unsigned 16-bit constants
// Ces constantes utilisent des entiers non signés 16 bits (0 à 65535)
const (
	// HTTPPort est le port HTTP standard
	HTTPPort uint16 = 80
	// HTTPSPort est le port HTTPS standard
	HTTPSPort uint16 = 443
	// customPort est un port personnalisé
	customPort uint16 = 3000
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : UINT32 (valeurs 0 à 4294967295)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint32 explicite
//    ✓ Approprié pour compteurs d'enregistrements, chunks, IDs
//    ✓ Commentaires avec unités (octets) quand applicable
// ════════════════════════════════════════════════════════════════════════════

// Unsigned 32-bit constants
// Ces constantes utilisent des entiers non signés 32 bits (0 à 4294967295)
const (
	// MaxRecordCount définit le nombre maximum d'enregistrements
	MaxRecordCount uint32 = 1000000
	// DefaultChunkSize est la taille par défaut d'un chunk en octets
	DefaultChunkSize uint32 = 65536
	// minBatchSize est la taille minimale d'un batch
	minBatchSize uint32 = 100
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : UINT64 (valeurs 0 à 18446744073709551615)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint64 explicite sur TOUTES les constantes
//    ✓ Approprié pour mémoire, IDs de transaction, très grandes valeurs positives
//    ✓ Commentaires en ligne pour clarifier les grandes valeurs (16 GB)
//
// ⚠️  IMPORTANT : Chaque constante DOIT avoir son type explicite, même dans un
//                 groupe où toutes ont le même type. Cohérence > concision.
// ════════════════════════════════════════════════════════════════════════════

// Unsigned 64-bit constants
// Ces constantes utilisent des entiers non signés 64 bits (très grandes valeurs positives)
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184 // 16 GB
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID uint64 = 18446744073709551615
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : BYTE (alias de uint8)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type byte explicite (pour données binaires/octets)
//    ✓ Notation hexadécimale appropriée (0x00, 0x0A, 0x09)
//    ✓ Commentaires expliquent la signification de chaque octet
//
// 📝 BYTE vs UINT8 :
//    - byte : Pour données binaires, protocoles, encodages
//    - uint8 : Pour valeurs numériques de 0 à 255
// ════════════════════════════════════════════════════════════════════════════

// Byte constants
// Ces constantes représentent des octets individuels pour encodages et protocoles
const (
	// NullByte représente l'octet null
	NullByte byte = 0x00
	// NewlineByte représente le caractère newline
	NewlineByte byte = 0x0A
	// tabByte représente le caractère tabulation
	tabByte byte = 0x09
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : RUNE (alias de int32)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type rune explicite (pour caractères Unicode)
//    ✓ Mix de caractères ASCII (' ', '\n') et Unicode (❤)
//    ✓ Commentaires expliquent chaque caractère
//
// 📝 RUNE vs INT32 :
//    - rune : Pour représenter des code points Unicode
//    - int32 : Pour valeurs numériques signées 32 bits
//    rune et int32 sont équivalents mais rune est plus expressif pour Unicode
// ════════════════════════════════════════════════════════════════════════════

// Rune constants
// Ces constantes représentent des caractères Unicode (code points)
const (
	// SpaceRune représente le caractère espace
	SpaceRune rune = ' '
	// NewlineRune représente le caractère retour à la ligne
	NewlineRune rune = '\n'
	// heartEmoji représente l'emoji cœur
	heartEmoji rune = '❤'
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : FLOAT32 (précision simple, ~7 décimales)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type float32 explicite
//    ✓ Commentaires mentionnent "approximation" pour valeurs mathématiques
//    ✓ Approprié quand 32 bits de précision suffisent
//
// 📝 FLOAT32 vs FLOAT64 :
//    - float32 : ~7 décimales de précision, économise mémoire
//    - float64 : ~15 décimales de précision, standard en Go
//    Utiliser float32 seulement si économie mémoire critique
// ════════════════════════════════════════════════════════════════════════════

// Float32 constants
// Ces constantes utilisent des nombres à virgule flottante 32 bits (précision simple)
const (
	// Pi32 est une approximation de Pi en float32
	Pi32 float32 = 3.14159265
	// DefaultRate est le taux par défaut
	DefaultRate float32 = 1.5
	// minThreshold est le seuil minimum
	minThreshold float32 = 0.01
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : FLOAT64 (double précision, ~15 décimales)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type float64 explicite (précision double)
//    ✓ Haute précision pour constantes mathématiques
//    ✓ Commentaires mentionnent les noms mathématiques (Euler, nombre d'or)
//
// 📝 BEST PRACTICE :
//    Préférer float64 par défaut en Go (c'est le type par défaut pour les littéraux).
//    Utiliser float32 seulement si besoin spécifique d'économie mémoire.
// ════════════════════════════════════════════════════════════════════════════

// Float64 constants
// Ces constantes utilisent des nombres à virgule flottante 64 bits (double précision)
const (
	// Pi est une approximation de Pi en haute précision
	Pi float64 = 3.14159265358979323846
	// EulerNumber est le nombre d'Euler (e)
	EulerNumber float64 = 2.71828182845904523536
	// goldenRatio est le nombre d'or (phi)
	goldenRatio float64 = 1.618033988749894848204586
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : COMPLEX64 (float32 + float32)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type complex64 explicite (composé de deux float32)
//    ✓ Notation mathématique i pour unité imaginaire
//    ✓ Commentaires expliquent chaque valeur complexe
//
// 📝 COMPLEX64 :
//    - Composé de deux float32 (partie réelle + partie imaginaire)
//    - Utilisé pour calculs scientifiques, traitement signal
//    - Notation : a + bi où i est l'unité imaginaire
// ════════════════════════════════════════════════════════════════════════════

// Complex64 constants
// Ces constantes représentent des nombres complexes en précision simple (float32 + float32)
const (
	// ImaginaryUnit64 représente l'unité imaginaire i en complex64
	ImaginaryUnit64 complex64 = 0 + 1i
	// ComplexZero64 représente zéro en complex64
	ComplexZero64 complex64 = 0 + 0i
	// sampleComplex64 est un exemple de nombre complexe
	sampleComplex64 complex64 = 3.5 + 2.8i
)

// ════════════════════════════════════════════════════════════════════════════
// CONSTANTES TYPE : COMPLEX128 (float64 + float64)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type complex128 explicite (composé de deux float64)
//    ✓ Commentaire de groupe explique "haute précision"
//    ✓ Chaque constante documentée avec son rôle mathématique
//
// 📝 COMPLEX128 :
//    - Composé de deux float64 (haute précision)
//    - Type complexe par défaut en Go
//    - Pour calculs scientifiques nécessitant précision
// ════════════════════════════════════════════════════════════════════════════

// Complex128 constants
// Ces constantes représentent des nombres complexes haute précision (float64 + float64)
const (
	// ImaginaryUnit représente l'unité imaginaire i
	ImaginaryUnit complex128 = 0 + 1i
	// ComplexZero représente zéro en complex128
	ComplexZero complex128 = 0 + 0i
	// eulerIdentityBase est la base pour l'identité d'Euler
	eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
)

// ════════════════════════════════════════════════════════════════════════════
// 📚 RÉSUMÉ DES BONNES PRATIQUES
// ════════════════════════════════════════════════════════════════════════════
//
// 1. REGROUPEMENT :
//    ✓ Toujours utiliser const () pour regrouper
//    ✓ Grouper les constantes par thème/domaine fonctionnel
//    ✓ Ne jamais déclarer const X = ... individuellement
//
// 2. COMMENTAIRES :
//    ✓ Commentaire de groupe : 2 lignes (titre + description)
//    ✓ Commentaire individuel : 1 ligne par constante
//    ✓ Commentaires au-dessus de la constante (pas à côté)
//
// 3. TYPES :
//    ✓ TOUJOURS spécifier le type explicitement
//    ✓ Choisir le bon type selon la plage de valeurs
//    ✓ Ne jamais se fier à l'inférence de type
//
// 4. NAMING :
//    ✓ MixedCaps : MaxConnections, defaultPort
//    ✓ Jamais underscore : max_connections ❌
//    ✓ Jamais ALL_CAPS : MAX_CONNECTIONS ❌
//    ✓ Initialismes en majuscules : HTTPPort, URLMaxLength
//
// 5. DOCUMENTATION :
//    ✓ Mentionner les unités (octets, millisecondes, etc.)
//    ✓ Mentionner les plages valides si pertinent (0-100)
//    ✓ Expliquer le rôle, pas juste répéter le nom
//
// 6. ORGANISATION :
//    ✓ Constantes du même domaine ensemble
//    ✓ Ordre logique par type (bool → int → float → complex)
//    ✓ Séparation visuelle avec commentaires de section
//
// ════════════════════════════════════════════════════════════════════════════
