# Go Best Practices - Bad Usage Examples (Non-idiomatique)

Ce dossier contient des exemples de **code qui compile mais viole les conventions Go officielles**.

## Objectif

Ces fichiers démontrent du code **non-idiomatique** qui :
- ✅ Compile sans erreur
- ❌ N'est pas idiomatique selon Effective Go
- ❌ Viole les Go Code Review Comments
- ❌ Ne suit pas les conventions de la communauté Go

## Différence avec bad_usage/ktn/

- **bad_usage/gospec/** : Viole les **conventions Go officielles** (Effective Go)
  - Exemple : `var x int = 42` au lieu de `x := 42`
  - Exemple : `func get_user_name()` au lieu de `func GetUserName()` (snake_case)
  - Exemple : Ignorer les erreurs silencieusement

- **bad_usage/ktn/** : Viole les **règles spécifiques KTN** (plus strictes)
  - Exemple : Fonction > 50 lignes (viole KTN-FUNC-006)
  - Exemple : > 5 paramètres (viole KTN-FUNC-005)
  - Exemple : Pas de commentaire sur return (viole KTN-FUNC-008)

## Structure des fichiers

### Fichiers actuels

1. **declarations_bad_practices.go** - Déclarations non-idiomatiques
   - Variables déclarées individuellement au lieu de groupées
   - Pas d'utilisation de short declarations (`:=`)
   - snake_case au lieu de MixedCaps

2. **control_flow_bad_practices.go** - Contrôle de flux non-idiomatique
   - Boucles C-style au lieu de range
   - else inutile après return
   - If imbriqués au lieu de early returns

3. **error_handling_bad_practices.go** - Gestion d'erreurs non-idiomatique
   - Erreurs ignorées
   - Pas de wrapping des erreurs
   - Panic pour erreurs normales
   - Comparaison de strings d'erreur

4. **naming_bad_practices.go** - Nommage non-idiomatique
   - snake_case au lieu de MixedCaps
   - ALL_CAPS pour variables
   - Acronymes mal formatés (Http au lieu de HTTP)
   - Noms trop génériques ou trop longs

5. **concurrency_bad_practices.go** - Concurrence non-idiomatique
   - Goroutines sans synchronisation
   - Channels pas fermés
   - Pas d'utilisation de context
   - Race conditions

6. **functions_bad_practices.go** - Fonctions non-idiomatiques
   - Trop de paramètres
   - Fonctions trop longues
   - Boolean parameters
   - Side effects cachés

7. **comments_bad_practices.go** - Documentation non-idiomatique
   - Pas de godoc sur exports
   - Commentaires ne commencent pas par le nom
   - Commentaires redondants
   - TODOs sans contexte

## Exemples de violations

### Déclarations
```go
// ❌ BAD: Explicit type when obvious
var x int = 42

// ❌ BAD: Individual declarations
const A = 1
const B = 2
const C = 3

// ❌ BAD: snake_case
func get_user_name() string { return "" }
```

### Gestion d'erreurs
```go
// ❌ BAD: Ignoring errors
result, _ := operation()

// ❌ BAD: No error wrapping
if err := load(); err != nil {
    return err  // No context
}

// ❌ BAD: Using panic for normal errors
if x < 0 {
    panic("negative value")
}
```

### Contrôle de flux
```go
// ❌ BAD: C-style loop instead of range
for i := 0; i < len(items); i++ {
    fmt.Println(items[i])
}

// ❌ BAD: Unnecessary else
if x > 0 {
    return x
} else {
    return 0
}

// ❌ BAD: Nested ifs instead of early return
if x > 0 {
    if x < 100 {
        if x%2 == 0 {
            return x
        }
    }
}
```

### Concurrence
```go
// ❌ BAD: No synchronization
go func() {
    fmt.Println("async work")
}()
// Program may exit before goroutine runs

// ❌ BAD: Channel not closed
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    // Should close(ch)
}()
```

### Nommage
```go
// ❌ BAD: snake_case
var user_count int

// ❌ BAD: ALL_CAPS for variables
var MAX_SIZE = 100

// ❌ BAD: Wrong acronym casing
var HttpClient int  // Should be HTTPClient
```

## Comment corriger ?

Pour chaque exemple non-idiomatique dans ce dossier, consultez le fichier correspondant dans `good_usage/gospec/` pour voir la version idiomatique.

### Ressources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Style Guide](https://google.github.io/styleguide/go/)

## Utilisation

Ces exemples servent à :
1. Identifier les patterns non-idiomatiques courants
2. Comprendre pourquoi certains patterns sont découragés
3. Apprendre à reconnaître le code non-idiomatique en code review
4. Former les développeurs aux conventions Go standard

## Note importante

⚠️ **Ce code COMPILE** mais n'est pas recommandé.

C'est la différence fondamentale avec les violations de la spec Go (erreurs de compilation).
Ici, le code fonctionne mais ne suit pas les best practices.

## Relation avec ktn/

- **gospec/** : Violations des conventions Go standard (Effective Go)
- **ktn/** : Violations des règles supplémentaires spécifiques au projet

Les règles KTN s'appuient sur gospec et ajoutent des contraintes plus strictes.
