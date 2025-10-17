package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-002 : Groupe sans commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc var () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ces variables mutables.
//
//    POURQUOI :
//    - Documente l'intention du regroupement
//    - Aide à comprendre pourquoi ces variables sont mutables
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces variables contiennent les métadonnées (mutables en production)
//    var (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Pas de commentaire de groupe avec strings
// NOTE : Tout est parfait (commentaires individuels + types) SAUF pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-VAR-002 UNIQUEMENT sur le groupe

var (
	// ApplicationNameV002 est le nom de l'application
	ApplicationNameV002 string = "MyApp"
	// VersionV002 est la version actuelle
	VersionV002 string = "1.0.0"
	// defaultEncodingV002 est l'encodage par défaut
	defaultEncodingV002 string = "UTF-8"
)

// ❌ CAS INCORRECT 2 : Pas de commentaire de groupe avec int64
// NOTE : Tout est parfait (commentaires individuels + types) SAUF pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-VAR-002 UNIQUEMENT sur le groupe

var (
	// MaxDiskSpaceV002 définit l'espace disque maximum en octets
	MaxDiskSpaceV002 int64 = 1099511627776
	// UnixEpochV002 représente le timestamp Unix de référence
	UnixEpochV002 int64 = 0
	// nanosPerSecondV002 est le nombre de nanosecondes dans une seconde
	nanosPerSecondV002 int64 = 1000000000
)
