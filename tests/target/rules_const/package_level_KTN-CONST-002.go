package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-002 : Groupe avec commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc const () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ces constantes.
//
//    POURQUOI :
//    - Documente l'intention du regroupement
//    - Aide à comprendre le rôle global des constantes
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces constantes contiennent les métadonnées de l'application
//    const (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Application metadata
// Ces constantes contiennent les métadonnées de l'application
const (
	// ApplicationNameC002Good est le nom de l'application
	ApplicationNameC002Good string = "MyApp"
	// VersionC002Good est la version actuelle de l'application
	VersionC002Good string = "1.0.0"
	// defaultEncodingC002Good est l'encodage par défaut utilisé
	defaultEncodingC002Good string = "UTF-8"
)

// Integer 64-bit constants
// Ces constantes utilisent des entiers 64 bits (très grandes valeurs)
const (
	// MaxDiskSpaceC002Good définit l'espace disque maximum en octets
	MaxDiskSpaceC002Good int64 = 1099511627776 // 1 TB
	// UnixEpochC002Good représente le timestamp Unix de référence
	UnixEpochC002Good int64 = 0
	// nanosPerSecondC002Good est le nombre de nanosecondes dans une seconde
	nanosPerSecondC002Good int64 = 1000000000
)
