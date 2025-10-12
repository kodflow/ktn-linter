package rules_const

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 1 : Constantes correctement regroupées
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Les constantes sont regroupées dans const ()
//    ✓ Le groupe a un commentaire de groupe (2 lignes)
//    ✓ Chaque constante a son commentaire individuel
//    ✓ Les types sont explicites (bool)
//    ✓ Naming MixedCaps respecté
//
// 📝 STRUCTURE :
//    // Ligne 1 : Titre court du groupe
//    // Ligne 2 : Description détaillée du groupe
//    const (
//        // Commentaire individuel constante 1
//        Constante1 Type = valeur
//        // Commentaire individuel constante 2
//        Constante2 Type = valeur
//    )
//
// 💡 POURQUOI :
//    - Regroupement facilite la navigation et la compréhension
//    - Commentaires permettent génération automatique documentation
//    - Types explicites évitent ambiguïté et erreurs de conversion
// ════════════════════════════════════════════════════════════════════════════

// Boolean configuration
// Ces constantes représentent des valeurs booléennes pour la configuration
const (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 2 : String constants avec métadonnées
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Groupe thématique cohérent (métadonnées d'application)
//    ✓ Commentaire de groupe explicite
//    ✓ Chaque constante documentée individuellement
//    ✓ Type string explicite pour toutes
//
// 📝 PATTERN :
//    Pour les constantes string, toujours être explicite sur le type
//    même si Go pourrait l'inférer. Cela améliore la lisibilité.
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

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 3 : Integer constants avec type int explicite
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type int explicite (important : taille dépend de l'architecture)
//    ✓ Commentaires précis sur le rôle de chaque constante
//    ✓ Mix de constantes publiques (MajConnections) et privées (maxRetries)
//
// 📝 IMPORTANT :
//    Le type "int" peut être 32 ou 64 bits selon l'architecture.
//    Si la taille est critique, utiliser int32 ou int64 explicitement.
// ════════════════════════════════════════════════════════════════════════════

// Integer constants (int)
// Ces constantes utilisent le type int (taille dépend de l'architecture: 32 ou 64 bits)
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 4 : Integer 8-bit (int8) pour petites valeurs
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type int8 explicite (valeurs -128 à 127)
//    ✓ Commentaire de groupe mentionne la plage de valeurs
//    ✓ Choix de int8 justifié (âges, priorités : petites valeurs)
//
// 📝 QUAND UTILISER int8 :
//    - Valeurs garanties dans la plage -128 à 127
//    - Optimisation mémoire importante (tableaux, structures)
//    - Besoin d'être explicite sur la taille
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
// ✅ EXEMPLE 5 : Integer 16-bit (int16) pour valeurs moyennes
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type int16 explicite (valeurs -32768 à 32767)
//    ✓ Approprié pour tailles de queue, buffers
//    ✓ Commentaires mentionnent l'unité ou le contexte
//
// 📝 QUAND UTILISER int16 :
//    - Valeurs dans la plage -32768 à 32767
//    - Compteurs, tailles de buffers modérés
//    - Besoin d'économiser mémoire vs int/int32
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
// ✅ EXEMPLE 6 : Integer 32-bit (int32) pour grandes valeurs
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type int32 explicite (valeurs -2147483648 à 2147483647)
//    ✓ Approprié pour tailles de fichiers, timeouts en ms
//    ✓ Commentaires incluent les unités (octets, millisecondes)
//
// 📝 QUAND UTILISER int32 :
//    - Valeurs dans la plage ~-2 milliards à ~2 milliards
//    - Garantir 32 bits même sur architecture 64 bits
//    - Compatibilité avec APIs qui attendent int32
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
// ✅ EXEMPLE 7 : Integer 64-bit (int64) pour très grandes valeurs
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type int64 explicite (très grandes valeurs)
//    ✓ Approprié pour espace disque, timestamps, nanosecondes
//    ✓ Commentaires explicites avec unités (octets, nanosecondes)
//    ✓ Valeur 0 documentée comme "intentionnelle" (UnixEpoch)
//
// 📝 QUAND UTILISER int64 :
//    - Très grandes valeurs (espace disque, timestamps)
//    - Timestamps Unix (secondes depuis 1970)
//    - Nanosecondes, microsecondes
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
// ✅ EXEMPLE 8 : Unsigned integer (uint) pour valeurs positives
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type uint explicite (taille architecture: 32 ou 64 bits)
//    ✓ Approprié pour IDs, tailles qui ne sont jamais négatives
//    ✓ Double la plage positive vs int de même taille
//
// 📝 QUAND UTILISER uint :
//    - Valeurs garanties positives (IDs, compteurs)
//    - Besoin de doubler la plage positive vs int
//    - APIs qui requièrent unsigned (ex: longueur d'array)
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
// ✅ EXEMPLE 9 : Unsigned 8-bit (uint8) alias de byte
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type uint8 explicite (valeurs 0 à 255)
//    ✓ Approprié pour qualité (0-100), niveaux (0-10)
//    ✓ Commentaires mentionnent les plages valides
//
// 📝 NOTE :
//    uint8 et byte sont équivalents en Go.
//    Utiliser byte pour données binaires, uint8 pour nombres.
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
// ✅ EXEMPLE 10 : Unsigned 16-bit (uint16) pour ports réseau
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type uint16 explicite (valeurs 0 à 65535)
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
// ✅ EXEMPLE 11 : Unsigned 32-bit (uint32) pour compteurs
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type uint32 explicite (valeurs 0 à 4294967295)
//    ✓ Approprié pour compteurs d'enregistrements, chunks
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
// ✅ EXEMPLE 12 : Unsigned 64-bit (uint64) pour très grandes valeurs positives
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type uint64 explicite (0 à 18446744073709551615)
//    ✓ TOUTES les constantes ont leur type (même MaxTransactionID!)
//    ✓ Approprié pour mémoire, IDs de transaction
//    ✓ Commentaires en ligne pour clarifier les valeurs (16 GB)
//
// 📝 IMPORTANT :
//    Chaque constante DOIT avoir son type explicite, même dans un groupe
//    où toutes ont le même type. Cohérence > concision.
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
// ✅ EXEMPLE 13 : Byte constants (alias de uint8)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type byte explicite (pour données binaires/octets)
//    ✓ Commentaire de groupe ET commentaires individuels
//    ✓ Notation hexadécimale appropriée (0x00, 0x0A)
//
// 📝 BYTE vs UINT8 :
//    - byte : Pour données binaires, protocoles, encodages
//    - uint8 : Pour valeurs numériques de 0 à 255
// ════════════════════════════════════════════════════════════════════════════

// Byte constants (alias de uint8)
// Ces constantes représentent des octets individuels (0 à 255)
const (
	// NullByte représente l'octet null
	NullByte byte = 0x00
	// NewlineByte représente le caractère newline
	NewlineByte byte = 0x0A
	// tabByte représente le caractère tabulation
	tabByte byte = 0x09
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 14 : Rune constants (alias de int32) pour Unicode
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type rune explicite (pour caractères Unicode)
//    ✓ Mix de caractères ASCII (' ', '\n') et Unicode (❤)
//    ✓ Commentaires expliquent chaque caractère
//
// 📝 RUNE vs INT32 :
//    - rune : Pour représenter des code points Unicode
//    - int32 : Pour valeurs numériques signées 32 bits
//    rune et int32 sont équivalents mais rune est plus expressif
// ════════════════════════════════════════════════════════════════════════════

// Rune constants (alias de int32)
// Ces constantes représentent des caractères Unicode
const (
	// SpaceRune représente le caractère espace
	SpaceRune rune = ' '
	// NewlineRune représente le caractère retour à la ligne
	NewlineRune rune = '\n'
	// heartEmoji représente l'emoji cœur
	heartEmoji rune = '❤'
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 15 : Float32 constants (précision simple)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type float32 explicite
//    ✓ Commentaires mentionnent "approximation" pour Pi
//    ✓ Approprié quand 32 bits de précision suffisent
//
// 📝 FLOAT32 vs FLOAT64 :
//    - float32 : ~7 décimales de précision, économise mémoire
//    - float64 : ~15 décimales de précision, standard en Go
//    Utiliser float32 seulement si économie mémoire critique
// ════════════════════════════════════════════════════════════════════════════

// Float32 constants
// Ces constantes utilisent des nombres à virgule flottante 32 bits
const (
	// Pi32 est une approximation de Pi en float32
	Pi32 float32 = 3.14159265
	// DefaultRate est le taux par défaut
	DefaultRate float32 = 1.5
	// minThreshold est le seuil minimum
	minThreshold float32 = 0.01
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 16 : Float64 constants (double précision)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type float64 explicite (précision double)
//    ✓ Haute précision pour constantes mathématiques
//    ✓ Commentaires mentionnent les noms mathématiques (Euler, nombre d'or)
//
// 📝 BEST PRACTICE :
//    Préférer float64 par défaut en Go (c'est le type par défaut).
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
// ✅ EXEMPLE 17 : Complex64 constants (nombres complexes simple précision)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type complex64 explicite (float32 + float32)
//    ✓ Notation mathématique i pour unité imaginaire
//    ✓ Commentaires expliquent chaque valeur complexe
//
// 📝 COMPLEX64 :
//    - Composé de deux float32 (partie réelle + partie imaginaire)
//    - Utilisé pour calculs scientifiques, traitement signal
//    - Notation : a + bi où i est l'unité imaginaire
// ════════════════════════════════════════════════════════════════════════════

// Complex64 constants
// Ces constantes représentent des nombres complexes (float32 + float32)
const (
	// ImaginaryUnit64 représente l'unité imaginaire i en complex64
	ImaginaryUnit64 complex64 = 0 + 1i
	// ComplexZero64 représente zéro en complex64
	ComplexZero64 complex64 = 0 + 0i
	// sampleComplex64 est un exemple de nombre complexe
	sampleComplex64 complex64 = 3.5 + 2.8i
)

// ════════════════════════════════════════════════════════════════════════════
// ✅ EXEMPLE 18 : Complex128 constants (nombres complexes haute précision)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ Type complex128 explicite (float64 + float64)
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
// ✅ EXEMPLE 19 : Theme configuration (cohésion thématique)
// ════════════════════════════════════════════════════════════════════════════
// 🎯 CE QUI EST CORRECT ICI :
//    ✓ TOUTES les constantes du même thème dans UN seul groupe
//    ✓ Pas de constantes séparées du même thème ailleurs
//    ✓ Cohérence : 4 thèmes ensemble, pas dispersés
//
// 📝 PRINCIPE DE COHÉSION :
//    Les constantes du même domaine fonctionnel doivent être groupées
//    ensemble, même si ajoutées à des moments différents.
//    Ne pas créer de "Partie 2" ailleurs dans le fichier.
// ════════════════════════════════════════════════════════════════════════════

// Theme configuration
// Ces constantes définissent les thèmes disponibles
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
//    ✓ Commentaires au-dessus (pas à côté)
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
//    ✓ Ordre logique (types simples → complexes)
//    ✓ Séparation visuelle avec commentaires de section
//
// ════════════════════════════════════════════════════════════════════════════
