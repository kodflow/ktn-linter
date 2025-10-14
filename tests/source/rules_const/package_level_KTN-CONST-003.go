package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-003 : Constante sans commentaire individuel
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
//    - Facilite la compréhension sans avoir à lire le code
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // HTTP configuration
//    // Ces constantes définissent les ports HTTP standards
//    const (
//        // HTTPPort est le port HTTP standard
//        HTTPPort uint16 = 80
//        // HTTPSPort est le port HTTPS standard
//        HTTPSPort uint16 = 443
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int8 sans commentaires individuels (SEULE ERREUR : KTN-CONST-003 x3)
// NOTE : Groupe OK, commentaire de groupe OK, types OK, MAIS pas de commentaires individuels
// ERREURS ATTENDUES : KTN-CONST-003 sur MinAge, MaxAge, defaultPriority
// Age configuration
// Ces constantes utilisent des entiers 8 bits pour les âges
const (
	MinAgeC003          int8 = 18
	MaxAgeC003          int8 = 120
	defaultPriorityC003 int8 = 5
)

// ❌ CAS INCORRECT 2 : Uint16 avec partiellement commentées (SEULE ERREUR : KTN-CONST-003 x2)
// NOTE : Groupe OK, commentaire groupe OK, types OK, HTTPPort commenté, MAIS HTTPSPort et customPort non commentés
// ERREURS ATTENDUES : KTN-CONST-003 sur HTTPSPortC003, customPortC003
// Port configuration
// Ces constantes définissent les ports réseau standards
const (
	// HTTPPortC003 est le port HTTP standard
	HTTPPortC003   uint16 = 80
	HTTPSPortC003  uint16 = 443
	customPortC003 uint16 = 3000
)

// ❌ CAS INCORRECT 3 : Float32 sans commentaires individuels (SEULE ERREUR : KTN-CONST-003 x3)
// NOTE : Groupe OK, commentaire groupe OK, types OK, MAIS pas de commentaires individuels
// ERREURS ATTENDUES : KTN-CONST-003 sur Pi32C003, DefaultRateC003, minThresholdC003
// Mathematical constants
// Ces constantes représentent des valeurs mathématiques en float32
const (
	Pi32C003         float32 = 3.14159265
	DefaultRateC003  float32 = 1.5
	minThresholdC003 float32 = 0.01
)

// ❌ CAS INCORRECT 4 : Complex128 avec première constante non commentée (SEULE ERREUR : KTN-CONST-003 x1)
// NOTE : Groupe OK, commentaire groupe OK, types OK, MAIS ImaginaryUnit sans commentaire individuel
// ERREUR ATTENDUE : KTN-CONST-003 sur ImaginaryUnitC003
// Complex number constants
// Ces constantes représentent des nombres complexes en complex128
const (
	ImaginaryUnitC003 complex128 = 0 + 1i
	// ComplexZeroC003 est zéro en complex128
	ComplexZeroC003 complex128 = 0 + 0i
	// eulerIdentityBaseC003 est la base de l'identité d'Euler
	eulerIdentityBaseC003 complex128 = 2.71828182845904523536 + 0i
)
