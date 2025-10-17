package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-001 : Constantes non groupées dans const ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les constantes doivent être regroupées dans un bloc const () au lieu
//    d'être déclarées individuellement avec "const X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité en regroupant les constantes liées
//    - Facilite la maintenance (une section = un thème)
//    - Évite la pollution du namespace package-level
//    - Standard Go universellement accepté
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces constantes configurent les fonctionnalités
//    const (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Bool non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur EnableFeatureX

// EnableFeatureXC001 active la fon¨ctionnalité X
const EnableFeatureXC001 bool = true

// ❌ CAS INCORRECT 2 : String non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur ThemeAutoC001

// ThemeAutoC001 est l'identifiant du thème automatique
const ThemeAutoC001 string = "auto"

// ❌ CAS INCORRECT 3 : Int non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur MaxUserIDC001

// MaxUserIDC001 définit l'ID utilisateur maximum
const MaxUserIDC001 uint32 = 4294967295

// ❌ CAS INCORRECT 4 : Int16 non groupés (SEULE ERREUR : KTN-CONST-001 x3)
// NOTE : Tout est parfait (commentaires + types) SAUF pas de ()
// ERREURS ATTENDUES : KTN-CONST-001 sur MaxQueueSizeC001, DefaultBufferSizeC001, minCacheSizeC001

// MaxQueueSizeC001 définit la taille maximale de la queue
const MaxQueueSizeC001 int16 = 10000

// DefaultBufferSizeC001 est la taille du buffer par défaut
const DefaultBufferSizeC001 int16 = 4096

// minCacheSizeC001 est la taille minimale du cache
const minCacheSizeC001 int16 = 512

// ❌ CAS INCORRECT 5 : Float64 non groupés (SEULE ERREUR : KTN-CONST-001 x3)
// NOTE : Tout est parfait (commentaires + types) SAUF pas de ()
// ERREURS ATTENDUES : KTN-CONST-001 sur PiC001, EulerNumberC001, goldenRatioC001

// PiC001 est une approximation de Pi en haute précision
const PiC001 float64 = 3.14159265358979323846

// EulerNumberC001 est le nombre d'Euler (e)
const EulerNumberC001 float64 = 2.71828182845904523536

// goldenRatioC001 est le nombre d'or (phi)
const goldenRatioC001 float64 = 1.618033988749894848204586
