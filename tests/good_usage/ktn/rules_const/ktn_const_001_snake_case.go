package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-001 : Constantes groupées dans const ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les constantes package-level doivent être regroupées dans un bloc const ()
//    au lieu d'être déclarées individuellement avec "const X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité et l'organisation du code
//    - Facilite la maintenance (constantes liées regroupées)
//    - Rend le code plus compact et structuré
//    - Standard Go universel pour constantes package-level
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces constantes représentent des valeurs booléennes
//    const (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Boolean configuration
// Ces constantes configurent les fonctionnalités de l'application
const (
	// EnableFeatureXC001Good active la fonctionnalité X
	EnableFeatureXC001Good bool = true
	// EnableDebugC001Good active le mode debug
	EnableDebugC001Good bool = false
	// isProductionC001Good indique si l'environnement est en production
	isProductionC001Good bool = true
)

// String configuration
// Ces constantes définissent les thèmes de l'application
const (
	// ThemeAutoC001Good est l'identifiant du thème automatique
	ThemeAutoC001Good string = "auto"
	// ThemeCustomC001Good est l'identifiant du thème personnalisé
	ThemeCustomC001Good string = "custom"
)

// Integer configuration
// Ces constantes configurent les limites entières
const (
	// MaxQueueSizeC001Good définit la taille maximale de la queue
	MaxQueueSizeC001Good int16 = 10000
	// DefaultBufferSizeC001Good est la taille du buffer par défaut
	DefaultBufferSizeC001Good int16 = 4096
	// minCacheSizeC001Good est la taille minimale du cache
	minCacheSizeC001Good int16 = 512
)
