# Go Spec - Good Usage Examples

Ce dossier contient des exemples de **code conforme à la spécification Go officielle**.

## Qu'est-ce que la spec Go ?

La spécification Go (https://go.dev/ref/spec) définit les **règles syntaxiques et sémantiques** du langage Go.

Elle couvre :
- La syntaxe valide (grammaire)
- Les règles de typage
- Les règles de scope et visibilité
- Le comportement des opérateurs
- Les règles de conversion de types

## Différence avec ktn/

- **gospec/** : Code **syntaxiquement valide** selon la spec Go
  - Exemple : `var x int = 42` est valide
  - Exemple : `var x string = 42` est invalide (erreur de type)

- **ktn/** : Code qui respecte les **conventions et best practices** au-delà de la spec
  - Exemple : `var x int = 42` est valide mais non idiomatique (devrait être `x := 42`)
  - Exemple : `const X = 42` est valide mais viole KTN-CONST-001 (devrait être groupé)

## Contenu actuel

🚧 **Ce dossier est actuellement vide.**

Il sera rempli avec des exemples de code qui :
1. Respecte strictement la spec Go
2. Compile sans erreur
3. A un comportement défini par la spec

Ces exemples serviront de référence pour comprendre :
- Ce qui est **permis** par Go (pas forcément recommandé)
- La différence entre "valide selon Go" vs "idiomatique selon KTN"

## Exemples à venir

### Déclarations valides mais non-idiomatiques
```go
// Valide selon gospec, mais non-idiomatique selon KTN
var x int = 42              // KTN préfère: x := 42
const A = 1                 // KTN préfère: const ( A = 1 )
func foo() {}               // KTN exige: documentation
```

### Types et conversions
```go
// Spec Go : conversions de types explicites requises
var i int = 42
var f float64 = float64(i)  // ✅ Valide selon spec
// var f float64 = i        // ❌ Invalide selon spec
```

### Scope et visibilité
```go
// Spec Go : majuscule = exporté, minuscule = privé
var PublicVar int = 1   // Exporté
var privateVar int = 2  // Privé au package
```

## Note

La majorité du code qui respecte la spec Go respecte aussi (ou peut être adapté pour respecter) les règles KTN.
Les règles KTN sont des **surcouches** de best practices, pas des contradictions avec la spec.
