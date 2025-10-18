package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-001 : Variables non groupées dans var ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les variables package-level doivent être regroupées dans un bloc var ()
//    au lieu d'être déclarées individuellement avec "var X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité et l'organisation du code
//    - Facilite la maintenance (variables liées regroupées)
//    - Rend les variables mutables explicites et visibles
//    - Standard Go universel pour variables package-level
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces variables configurent les fonctionnalités (mutables)
//    var (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Bool non groupés
// ERREURS : KTN-VAR-001 sur EnableFeatureXV001, EnableDebugV001, isProductionV001
var EnableFeatureXV001 bool = true

// EnableDebugV001 describes this variable.
var EnableDebugV001 bool = false

// isProductionV001 describes this variable.
var isProductionV001 bool = true

// ❌ CAS INCORRECT 2 : String non groupés
// ERREURS : KTN-VAR-001 sur ThemeAutoV001, ThemeCustomV001
var ThemeAutoV001 string = "auto"

// ThemeCustomV001 describes this variable.
var ThemeCustomV001 string = "custom"

// ❌ CAS INCORRECT 3 : Int16 non groupés
// ERREURS : KTN-VAR-001 sur MaxQueueSizeV001, DefaultBufferSizeV001, minCacheSizeV001
var MaxQueueSizeV001 int16 = 10000

// DefaultBufferSizeV001 describes this variable.
var DefaultBufferSizeV001 int16 = 4096

// minCacheSizeV001 describes this variable.
var minCacheSizeV001 int16 = 512

// ❌ CAS INCORRECT 4 : Variables non groupées avec type manquant
// ERREURS : KTN-VAR-001 + KTN-VAR-004 sur MaxUserIDV001, DefaultPoolSizeV001, minWorkersV001
var MaxUserIDV001 = 4294967295

// DefaultPoolSizeV001 describes this variable.
var DefaultPoolSizeV001 = 100

// minWorkersV001 describes this variable.
var minWorkersV001 = 4

// ❌ CAS INCORRECT 5 : Variable orpheline (toutes les erreurs)
// ERREURS : KTN-VAR-001 + KTN-VAR-004
var orphanVarV001 = 42
