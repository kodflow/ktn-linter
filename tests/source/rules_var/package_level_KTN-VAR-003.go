package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-003 : Variable sans commentaire individuel
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    CHAQUE variable (publique ET privée) doit avoir son propre commentaire
//    individuel qui explique son rôle spécifique. Le commentaire doit être
//    sur la ligne juste au-dessus de la variable.
//
//    POURQUOI :
//    - Documente précisément le rôle de CETTE variable
//    - Obligatoire pour les variables publiques (godoc)
//    - Recommandé aussi pour les privées (maintenabilité)
//    - Variables mutables nécessitent plus de documentation
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // HTTP configuration
//    // Ces variables configurent le serveur HTTP (mutables)
//    var (
//        // HTTPPort est le port HTTP à utiliser
//        HTTPPort uint16 = 80
//        // HTTPSPort est le port HTTPS à utiliser
//        HTTPSPort uint16 = 443
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int8 avec commentaire de groupe mais pas individuels
// ERREURS : KTN-VAR-003 sur MinAgeV003, MaxAgeV003, defaultPriorityV003
// Ces variables utilisent des entiers 8 bits (-128 à 127)
var (
	MinAgeV003          int8 = 18
	MaxAgeV003          int8 = 120
	defaultPriorityV003 int8 = 5
)
