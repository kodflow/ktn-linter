package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-002 : Groupe sans commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc const () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ce groupe de constantes.
//
//    POURQUOI :
//    - Documente l'intention du regroupement (pourquoi ces constantes ensemble ?)
//    - Aide les développeurs à comprendre le contexte global
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité long terme
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces constantes définissent les métadonnées de l'application
//    const (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Pas de commentaire de groupe (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// ApplicationNameC002 est le nom de l'application
	ApplicationNameC002 string = "MyApp"
	// VersionC002 est la version actuelle
	VersionC002 string = "1.0.0"
	// defaultEncodingC002 est l'encodage par défaut
	defaultEncodingC002 string = "UTF-8"
)

// ❌ CAS INCORRECT 2 : Pas de commentaire de groupe avec int64 (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// MaxDiskSpaceC002 définit l'espace disque maximum en octets
	MaxDiskSpaceC002 int64 = 1099511627776
	// UnixEpochC002 représente le timestamp Unix de référence
	UnixEpochC002 int64 = 0
	// nanosPerSecondC002 est le nombre de nanosecondes dans une seconde
	nanosPerSecondC002 int64 = 1000000000
)

// ❌ CAS INCORRECT 3 : Pas de commentaire de groupe avec byte (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// NullByteC002 représente l'octet null
	NullByteC002 byte = 0x00
	// NewlineByteC002 représente l'octet newline
	NewlineByteC002 byte = 0x0A
	// tabByteC002 représente l'octet tabulation
	tabByteC002 byte = 0x09
)

// ❌ CAS INCORRECT 4 : Pas de commentaire de groupe avec complex64 (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// ImaginaryUnit64C002 est l'unité imaginaire en complex64
	ImaginaryUnit64C002 complex64 = 0 + 1i
	// ComplexZero64C002 est zéro en complex64
	ComplexZero64C002 complex64 = 0 + 0i
	// sampleComplex64C002 est un exemple de nombre complexe
	sampleComplex64C002 complex64 = 3.5 + 2.8i
)
