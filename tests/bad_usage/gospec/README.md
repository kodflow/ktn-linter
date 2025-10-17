# Go Spec - Bad Usage Examples (Violations de la spec)

Ce dossier contient des exemples de **code qui viole la spécification Go officielle**.

## Objectif

Ces fichiers démontrent du code **invalide selon la spec Go** qui :
- Ne compile pas
- Viole les règles de typage
- Viole les règles syntaxiques
- A un comportement non-défini

## Différence avec bad_usage/ktn/

- **bad_usage/gospec/** : Code qui **ne respecte PAS la spec Go**
  - ❌ Erreurs de compilation
  - ❌ Erreurs de typage
  - ❌ Syntaxe invalide
  - Exemple : `var x string = 42` (erreur de type)

- **bad_usage/ktn/** : Code qui **compile** mais **viole les conventions KTN**
  - ✅ Compile sans erreur
  - ❌ Non-idiomatique
  - ❌ Mauvaises pratiques
  - Exemple : `const X = 42` (pas groupé, viole KTN-CONST-001)

## Contenu actuel

🚧 **Ce dossier est actuellement vide.**

Il sera rempli avec des exemples de violations de la spec Go tels que :

### Erreurs de typage
```go
// ❌ Type mismatch
var x string = 42

// ❌ Cannot assign to constant
const PI = 3.14
PI = 3.14159

// ❌ Invalid operation
var result = "hello" + 42
```

### Erreurs syntaxiques
```go
// ❌ Missing package declaration
func main() {}

// ❌ Invalid identifier
var 123abc int

// ❌ Incorrect function signature
func foo() int {
    return  // Missing return value
}
```

### Violations de scope
```go
// ❌ Undeclared variable
func foo() {
    x = 42  // x not declared
}

// ❌ Redeclaration in same block
func bar() {
    var x int
    var x string  // Redeclared
}
```

### Violations de visibilité
```go
// Dans package A
package a
var privateVar int = 1

// Dans package B
package b
import "a"
func foo() {
    a.privateVar = 2  // ❌ Cannot access private var
}
```

## Usage

Ces fichiers servent à :
1. Documenter les erreurs courantes par rapport à la spec Go
2. Tester que le compilateur Go détecte bien les violations
3. Éduquer sur ce qui est **strictement interdit** par Go

## Note importante

⚠️ **Ce code NE COMPILE PAS** car il viole la spec Go.

C'est la différence fondamentale avec `bad_usage/ktn/` qui compile mais n'est pas idiomatique.

## Relation avec ktn/

Tous les fichiers dans `bad_usage/ktn/` :
- ✅ Respectent la spec Go (compilent)
- ❌ Violent les conventions KTN (non-idiomatiques)

Tous les fichiers dans `bad_usage/gospec/` :
- ❌ Violent la spec Go (ne compilent pas)
- N/A pour KTN (impossible d'évaluer les conventions si le code ne compile pas)
