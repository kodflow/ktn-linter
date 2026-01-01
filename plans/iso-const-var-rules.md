# Plan: Réorganisation ISO CONST/VAR

## Résumé Exécutif

Alignement des règles CONST et VAR pour avoir les mêmes numéros pour les mêmes concepts.

| # | Concept | CONST | VAR |
|---|---------|-------|-----|
| 001 | Types explicites | ✓ Existe | Rename VAR-002 |
| 002 | Ordre déclaration | ✓ Existe | Rename VAR-014 |
| 003 | CamelCase | ✓ Existe | Fusion VAR-001+018 |
| 004 | Longueur min | **CRÉER** | **CRÉER** |
| 005 | Longueur max | **CRÉER** | **CRÉER** |
| 006 | Shadowing | **CRÉER** | Rename VAR-011 |

---

## EPIC 1: CONST - Nouvelles règles (004, 005, 006)

### Task 1.1: CONST-004 - Longueur minimale

**Fichiers à créer:**
- `pkg/analyzer/ktn/ktnconst/004.go`
- `pkg/analyzer/ktn/ktnconst/004_external_test.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const004/good.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const004/bad.go`

**Algorithme:**
```go
func isConstNameTooShort(name string) bool {
    if name == "_" { return false } // blank identifier OK
    return len(name) < 2
}
```

**Cas good.go (0 erreurs):**
```go
const OK int = 0           // 2 chars - OK
const Pi float64 = 3.14    // 2 chars - OK
const MaxSize int = 100    // Normal
const _ int = 0            // Blank identifier - toujours OK
```

**Cas bad.go:**
```go
const A int = 1            // want "KTN-CONST-004"
const B = "hello"          // want "KTN-CONST-004"
const X int = 100          // want "KTN-CONST-004"
const (
    C = iota               // want "KTN-CONST-004"
    D                      // want "KTN-CONST-004"
)
```

**Edge cases:**
- `_` (blank) → Ignoré
- Unicode single rune (é) → Compte comme 1 char → Erreur
- Iota blocks → Chaque nom vérifié individuellement

---

### Task 1.2: CONST-005 - Longueur maximale

**Fichiers à créer:**
- `pkg/analyzer/ktn/ktnconst/005.go`
- `pkg/analyzer/ktn/ktnconst/005_external_test.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const005/good.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const005/bad.go`

**Algorithme:**
```go
const maxConstNameLength = 30

func isConstNameTooLong(name string) bool {
    return len(name) > maxConstNameLength
}
```

**Cas good.go:**
```go
const MaxConnectionPoolSize int = 100      // 21 chars - OK
const DefaultTimeout int = 30              // 14 chars - OK
const ThisIsExactlyThirtyCharsX int = 1    // 30 chars - OK (limite)
```

**Cas bad.go:**
```go
const ThisIsAVeryLongConstantNameThatExceedsLimit int = 1
// want "KTN-CONST-005"

const VeryLongConstantNameForSomeConfigurationValue string = "x"
// want "KTN-CONST-005"
```

---

### Task 1.3: CONST-006 - Shadowing built-in

**Fichiers à créer:**
- `pkg/analyzer/ktn/ktnconst/006.go`
- `pkg/analyzer/ktn/ktnconst/006_external_test.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const006/good.go`
- `pkg/analyzer/ktn/ktnconst/testdata/src/const006/bad.go`

**Liste des 44 built-ins Go:**
```go
var builtins = map[string]bool{
    // Types (25)
    "bool": true, "byte": true, "complex64": true, "complex128": true,
    "error": true, "float32": true, "float64": true, "int": true,
    "int8": true, "int16": true, "int32": true, "int64": true,
    "rune": true, "string": true, "uint": true, "uint8": true,
    "uint16": true, "uint32": true, "uint64": true, "uintptr": true,
    "any": true,
    // Constantes (3)
    "true": true, "false": true, "iota": true,
    // Zero-value (1)
    "nil": true,
    // Fonctions (15)
    "append": true, "cap": true, "close": true, "complex": true,
    "copy": true, "delete": true, "imag": true, "len": true,
    "make": true, "new": true, "panic": true, "print": true,
    "println": true, "real": true, "recover": true,
}
```

**Cas good.go:**
```go
const MaxSize int = 100        // Ne shadow rien
const ApiKey string = "secret" // Ne shadow rien
const ErrorCode int = 1        // "ErrorCode" != "error"
const NewValue int = 42        // "NewValue" != "new"
```

**Cas bad.go:**
```go
const int int = 32             // want "KTN-CONST-006"
const bool = true              // want "KTN-CONST-006"
const error = 1                // want "KTN-CONST-006"
const true = 1                 // want "KTN-CONST-006"
const false = 0                // want "KTN-CONST-006"
const nil = 0                  // want "KTN-CONST-006"
const append = "log"           // want "KTN-CONST-006"
const make = 42                // want "KTN-CONST-006"
const len = 100                // want "KTN-CONST-006"
const new = 1                  // want "KTN-CONST-006"
const panic = 3                // want "KTN-CONST-006"
```

---

### Task 1.4: Mise à jour registry CONST

**Fichier:** `pkg/analyzer/ktn/ktnconst/registry.go`

```go
func Analyzers() []*analysis.Analyzer {
    return []*analysis.Analyzer{
        Analyzer001,
        Analyzer002,
        Analyzer003,
        Analyzer004, // NEW
        Analyzer005, // NEW
        Analyzer006, // NEW
    }
}
```

---

## EPIC 2: VAR - Réorganisation majeure

### Task 2.1: Fusion VAR-001 + VAR-018 → VAR-003 (CamelCase)

**Action:** Combiner la détection SCREAMING_SNAKE_CASE + snake_case

**Nouveau VAR-003:**
```go
func hasUnderscoreInName(name string) bool {
    if name == "_" { return false } // blank identifier OK
    return strings.Contains(name, "_")
}
```

**Cas bad.go (fusionné):**
```go
var MAX_SIZE int = 100       // SCREAMING_SNAKE_CASE
var my_var int = 200         // snake_case
var Api_Key string = "x"     // Mixed_Case
```

**Cas good.go:**
```go
var maxSize int = 100        // camelCase
var MaxSize int = 100        // PascalCase
var HTTPStatus int = 200     // Acronymes OK
```

**Fichiers à supprimer après fusion:**
- `001.go` (absorbé)
- `018.go` (absorbé)

---

### Task 2.2: Renommage VAR-002 → VAR-001 (Types explicites)

**Fichiers à renommer:**
```
002.go → 001.go
002_external_test.go → 001_external_test.go
002_internal_test.go → 001_internal_test.go
testdata/src/var002/ → testdata/src/var001/
```

**Mise à jour dans le code:**
- `ruleCodeVar002` → `ruleCodeVar001`
- `Analyzer002` → `Analyzer001`
- Doc: "KTN-VAR-002" → "KTN-VAR-001"

---

### Task 2.3: Renommage VAR-014 → VAR-002 (Ordre déclaration)

**Fichiers à renommer:**
```
014.go → 002.go
014_*_test.go → 002_*_test.go
testdata/src/var014/ → testdata/src/var002/
```

---

### Task 2.4: VAR-004 - Longueur minimale (CRÉER)

**Algorithme avec exceptions:**
```go
var loopVars = map[string]bool{
    "i": true, "j": true, "k": true, "n": true,
    "x": true, "y": true, "z": true,
}

var idiomaticShort = map[string]bool{
    "err": true, "ok": true, "ctx": true,
    "id": true, "db": true, "tx": true,
}

func isVarNameTooShort(name string, isLoopVar bool) bool {
    if name == "_" { return false }
    if len(name) >= 2 { return false }
    if loopVars[name] && isLoopVar { return false }
    if idiomaticShort[name] { return false }
    return true
}
```

**Cas good.go:**
```go
for i := 0; i < 10; i++ {}     // loop var OK
for j, v := range slice {}     // loop vars OK
err := doSomething()           // idiomatique OK
val, ok := m[key]              // idiomatique OK
ctx := context.Background()    // idiomatique OK
```

**Cas bad.go:**
```go
func BadExample() {
    a := 42                    // want "KTN-VAR-004"
    b := "hello"               // want "KTN-VAR-004"
    x := compute()             // want "KTN-VAR-004" (pas dans loop)
}
```

---

### Task 2.5: VAR-005 - Longueur maximale (CRÉER)

**Algorithme:**
```go
const maxVarNameLength = 30

func isVarNameTooLong(name string) bool {
    return len(name) > maxVarNameLength
}
```

---

### Task 2.6: Renommage VAR-011 → VAR-006 (Shadowing)

**Fichiers à renommer:**
```
011.go → 006.go
011_*_test.go → 006_*_test.go
testdata/src/var011/ → testdata/src/var006/
```

---

### Task 2.7: Renommage cascade VAR-003→007 jusqu'à VAR-017→019

**Table de renommage complète:**

| Ancien | Nouveau | Règle |
|--------|---------|-------|
| VAR-002 | VAR-001 | Types explicites |
| VAR-014 | VAR-002 | Ordre déclaration |
| FUSION | VAR-003 | CamelCase |
| CRÉER | VAR-004 | Longueur min |
| CRÉER | VAR-005 | Longueur max |
| VAR-011 | VAR-006 | Shadowing |
| VAR-003 | VAR-007 | := vs var |
| VAR-004 | VAR-008 | Slices préalloc |
| VAR-005 | VAR-009 | make+append |
| VAR-006 | VAR-010 | Buffer.Grow |
| VAR-007 | VAR-011 | strings.Builder |
| VAR-008 | VAR-012 | Alloc loops |
| VAR-009 | VAR-013 | Struct size |
| VAR-010 | VAR-014 | sync.Pool |
| VAR-012 | VAR-015 | string() |
| VAR-013 | VAR-016 | Groupement |
| VAR-015 | VAR-017 | Map prealloc |
| VAR-016 | VAR-018 | Array vs slice |
| VAR-017 | VAR-019 | Mutex copies |

**Fichiers à supprimer:**
- VAR-001 (absorbé dans VAR-003)
- VAR-018 (absorbé dans VAR-003)

---

### Task 2.8: Mise à jour registry VAR

```go
func Analyzers() []*analysis.Analyzer {
    return []*analysis.Analyzer{
        Analyzer001, // Types explicites (ex-002)
        Analyzer002, // Ordre déclaration (ex-014)
        Analyzer003, // CamelCase (fusion 001+018)
        Analyzer004, // Longueur min (NEW)
        Analyzer005, // Longueur max (NEW)
        Analyzer006, // Shadowing (ex-011)
        Analyzer007, // := vs var (ex-003)
        Analyzer008, // Slices préalloc (ex-004)
        Analyzer009, // make+append (ex-005)
        Analyzer010, // Buffer.Grow (ex-006)
        Analyzer011, // strings.Builder (ex-007)
        Analyzer012, // Alloc loops (ex-008)
        Analyzer013, // Struct size (ex-009)
        Analyzer014, // sync.Pool (ex-010)
        Analyzer015, // string() (ex-012)
        Analyzer016, // Groupement (ex-013)
        Analyzer017, // Map prealloc (ex-015)
        Analyzer018, // Array vs slice (ex-016)
        Analyzer019, // Mutex copies (ex-017)
    }
}
```

---

## EPIC 3: Messages et Sévérités

### Task 3.1: Ajouter messages CONST

**Fichier:** `pkg/messages/const.go`

```go
"KTN-CONST-004": {
    Short: "constant name too short",
    Verbose: "constant name '%s' is too short (min 2 characters)",
},
"KTN-CONST-005": {
    Short: "constant name too long",
    Verbose: "constant name '%s' exceeds maximum length of 30 characters",
},
"KTN-CONST-006": {
    Short: "constant shadows built-in",
    Verbose: "constant '%s' shadows Go built-in identifier",
},
```

### Task 3.2: Mettre à jour messages VAR

Renommer tous les codes de règles dans `pkg/messages/var.go`

### Task 3.3: Mettre à jour sévérités

**Fichier:** `pkg/severity/severity.go`

---

## EPIC 4: Tests et Validation

### Task 4.1: Exécuter tests après chaque étape

```bash
go test ./pkg/analyzer/ktn/ktnconst/... -v
go test ./pkg/analyzer/ktn/ktnvar/... -v
```

### Task 4.2: Valider testdata avec linter direct

```bash
./builds/ktn-linter lint pkg/analyzer/ktn/ktnconst/testdata/src/const004/bad.go
./builds/ktn-linter lint pkg/analyzer/ktn/ktnconst/testdata/src/const004/good.go
```

### Task 4.3: Mise à jour CLAUDE.md

- `pkg/analyzer/ktn/ktnconst/CLAUDE.md`
- `pkg/analyzer/ktn/ktnvar/CLAUDE.md`

---

## Ordre d'exécution recommandé

1. **CONST d'abord** (plus simple, pas de réorg)
   - Task 1.1 → 1.2 → 1.3 → 1.4
   - Tests après chaque task

2. **VAR ensuite** (réorg complexe)
   - Task 2.1 (fusion) → 2.2 → 2.3 (renommages critiques)
   - Task 2.4 → 2.5 (nouvelles règles)
   - Task 2.6 → 2.7 (cascade renommages)
   - Task 2.8 (registry)

3. **Messages et sévérités**
   - Task 3.1 → 3.2 → 3.3

4. **Validation finale**
   - Task 4.1 → 4.2 → 4.3

---

## Estimation

| Epic | Tasks | Complexité |
|------|-------|------------|
| CONST | 4 | Moyenne |
| VAR | 8 | Élevée |
| Messages | 3 | Faible |
| Tests | 3 | Moyenne |

**Total:** 18 tasks
