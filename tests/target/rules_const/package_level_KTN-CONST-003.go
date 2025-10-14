package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-003 : Constante avec commentaire individuel
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    CHAQUE constante (publique ET privée) doit avoir son propre commentaire
//    individuel qui explique son rôle spécifique. Le commentaire doit être
//    sur la ligne juste au-dessus de la constante.
//
//    POURQUOI :
//    - Documente précisément le rôle de CETTE constante
//    - Obligatoire pour les constantes publiques (godoc)
//    - Recommandé aussi pour les privées (maintenabilité)
//    - Facilite la compréhension du code
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer 8-bit constants
//    // Ces constantes utilisent des entiers 8 bits (-128 à 127)
//    const (
//        // MinAge est l'âge minimum requis
//        MinAge int8 = 18
//        // MaxAge est l'âge maximum accepté
//        MaxAge int8 = 120
//        // defaultPriority est la priorité par défaut
//        defaultPriority int8 = 5
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Integer 8-bit constants
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	// MinAgeC003Good est l'âge minimum requis
	MinAgeC003Good int8 = 18
	// MaxAgeC003Good est l'âge maximum accepté
	MaxAgeC003Good int8 = 120
	// defaultPriorityC003Good est la priorité par défaut
	defaultPriorityC003Good int8 = 5
)

// Unsigned 16-bit constants
// Ces constantes utilisent des entiers non signés 16 bits (0 à 65535)
const (
	// HTTPPortC003Good est le port HTTP standard
	HTTPPortC003Good uint16 = 80
	// HTTPSPortC003Good est le port HTTPS standard
	HTTPSPortC003Good uint16 = 443
	// customPortC003Good est un port personnalisé
	customPortC003Good uint16 = 3000
)
