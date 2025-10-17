# Go Best Practices - Good Usage Examples

Ce dossier contient des exemples de **code idiomatique suivant les conventions officielles de Go**.

## Qu'est-ce que les Go best practices ?

Les conventions Go sont définies dans plusieurs documents officiels :
- **Effective Go** (https://go.dev/doc/effective_go) : Le guide de style officiel
- **Go Code Review Comments** (https://github.com/golang/go/wiki/CodeReviewComments) : Conseils de code review
- **Go Proverbs** : Principes de conception Go

Ces pratiques couvrent :
- Les conventions de nommage (MixedCaps, interfaces en -er, etc.)
- La gestion des erreurs idiomatique (wrapping, early returns)
- Les patterns de concurrence (WaitGroups, channels, context)
- Le contrôle de flux (range, early returns, no else after return)
- La documentation (godoc, commentaires clairs)
- Les patterns de fonctions (options pattern, composition)

## Différence avec ktn/

- **gospec/** : Pratiques **idiomatiques officielles de Go**
  - Exemple : `x := 42` (court et clair)
  - Exemple : Utiliser `errors.Is()` pour vérifier les erreurs
  - Exemple : Interfaces avec suffixe -er pour méthode unique

- **ktn/** : Règles **spécifiques au projet KTN** (plus strictes)
  - Exemple : Fonctions max 50 lignes (KTN-FUNC-006)
  - Exemple : Max 5 paramètres par fonction (KTN-FUNC-005)
  - Exemple : Commentaires obligatoires sur tous les returns (KTN-FUNC-008)

## Structure des fichiers

### Fichiers actuels

1. **types_basic.go** - Types de base Go
   - Démontre l'utilisation correcte de tous les types Go

2. **declarations.go** - Déclarations idiomatiques
   - Constants groupées avec `const ()`
   - Variables avec inférence de type
   - Short declarations avec `:=`

3. **statements.go** - Statements idiomatiques
   - Range loops au lieu de for classiques
   - Early returns pour réduire la complexité
   - Select avec timeout/default

4. **naming_good_practices.go** - Conventions de nommage
   - MixedCaps pour exports, mixedCaps pour privés
   - Pas de snake_case
   - Acronymes en majuscules (HTTP, JSON, XML)
   - Interfaces avec suffixe -er

5. **error_handling_good_practices.go** - Gestion d'erreurs
   - Error wrapping avec `fmt.Errorf("...: %w", err)`
   - Sentinel errors avec `var ErrNotFound = ...`
   - Early returns au lieu de nested if
   - Utilisation de `errors.Is()` et `errors.As()`

6. **control_flow_good_practices.go** - Contrôle de flux
   - Range loops appropriés
   - No else after return
   - Switch au lieu de if/else chains
   - Labels pour break/continue en nested loops

7. **concurrency_good_practices.go** - Concurrence
   - WaitGroups pour synchronisation
   - Context pour cancellation
   - Worker pools
   - Sender closes channels

8. **functions_good_practices.go** - Fonctions
   - Functional options pattern
   - Focused, single-responsibility functions
   - Consistent receiver types
   - Defer for cleanup

9. **comments_good_practices.go** - Documentation
   - Godoc sur exports
   - Comments commencent par le nom
   - Package comment au début
   - TODOs avec contexte

## Exemples de patterns idiomatiques

### Déclarations
```go
// ✅ GOOD: Short declaration
x := 42

// ✅ GOOD: Grouped constants
const (
    StatusActive = "active"
    StatusInactive = "inactive"
)

// ✅ GOOD: Type inference
var count = 10
```

### Gestion d'erreurs
```go
// ✅ GOOD: Error wrapping
if err := operation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ✅ GOOD: Sentinel error
var ErrNotFound = errors.New("not found")

// ✅ GOOD: Early return
if x < 0 {
    return fmt.Errorf("invalid")
}
// Happy path continues
```

### Concurrence
```go
// ✅ GOOD: WaitGroup pattern
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        process(n)
    }(i)
}
wg.Wait()

// ✅ GOOD: Context for cancellation
func worker(ctx context.Context) {
    select {
    case <-ctx.Done():
        return
    default:
        // work
    }
}
```

### Nommage
```go
// ✅ GOOD: MixedCaps
func GetUserByID(id int) (*User, error)

// ✅ GOOD: Interface with -er suffix
type Reader interface {
    Read() ([]byte, error)
}

// ✅ GOOD: Acronyms in uppercase
var HTTPClient *http.Client
var JSONData []byte
```

## Utilisation

Ces exemples servent de référence pour :
1. Comprendre les patterns idiomatiques Go
2. Distinguer code idiomatique vs non-idiomatique (voir ../bad_usage/gospec/)
3. Apprendre les conventions avant d'appliquer les règles KTN plus strictes
4. Former les développeurs aux pratiques Go standard

## Note

Les pratiques KTN (dossier ktn/) s'appuient sur ces conventions Go et ajoutent des règles plus strictes adaptées aux besoins du projet.
