# Tasks - KTN-VAR

## Vue d'ensemble

Le module `ktnvar` contient 17 règles (001-017) pour la déclaration, le nommage et l'optimisation des variables.

---

## KTN-VAR-001 : Variables de package en camelCase (pas SCREAMING_SNAKE_CASE)

### Points positifs
- Suit la convention Go idiomatique
- Distinction claire entre variables et constantes
- Cohérent avec Effective Go

### Points négatifs / Problèmes identifiés
1. **Conflit avec CONST-003** : Les constantes utilisent SCREAMING_SNAKE_CASE
2. **Confusion possible** : Variable vs constante visuellement similaires si même casse

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Documenter la distinction variable (camelCase) vs constante (SCREAMING_SNAKE_CASE)
- [ ] **VÉRIFICATION** : S'assurer que les acronymes sont bien gérés (`httpClient` vs `HTTPClient`)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Règles pour les variables exportées (`MaxRetries` vs `maxRetries`)

### Scénarios non couverts
1. Variables avec acronymes (`xmlParser`, `httpServer`)
2. Variables exportées commençant par acronyme

---

## KTN-VAR-002 : Format 'var name type = value' pour les variables de package

### Points positifs
- Type explicite améliore la lisibilité
- Évite les inférences de type inattendues
- Cohérent avec CONST-001

### Points négatifs / Problèmes identifiés
1. **Verbosité** : `var timeout time.Duration = 5 * time.Second` vs `var timeout = 5 * time.Second`
2. **Types évidents** : Le type est parfois évident de l'expression
3. **Zéro-values** : `var count int` vs `var count int = 0`

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Faut-il `= value` même pour les zéro-values ?
- [ ] **AMÉLIORATION** : Tolérer les expressions avec type évident (`time.Second`, `http.StatusOK`)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Exemples de bon format
- [ ] **VÉRIFICATION** : Gérer les déclarations multiples `var a, b int = 1, 2`

### Scénarios non couverts
1. Variables sans valeur initiale (`var mu sync.Mutex`)
2. Variables avec expression typée (`var t = time.Now()`)

---

## KTN-VAR-003 : Variables locales avec ':=' au lieu de 'var'

### Points positifs
- Plus concis et idiomatique Go
- Réduit la verbosité dans les fonctions
- Convention standard

### Points négatifs / Problèmes identifiés
1. **Shadowing** : `:=` peut masquer une variable extérieure (voir VAR-011)
2. **Cas où var est nécessaire** : Type différent de l'expression
3. **Nil explicite** : `var slice []int` vs impossible avec `:=`

### Actions à mener

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Ignorer les cas où `var` est nécessaire (nil, type différent)
- [ ] **COHÉRENCE** : Aligner avec VAR-011 (shadowing)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Cas légitimes d'utilisation de `var`
- [ ] **VÉRIFICATION** : Gérer les déclarations dans les blocs `if`, `for`, etc.

### Scénarios non couverts
1. `var err error` avant un bloc `if` pour scope
2. `var x T` quand on veut un nil slice/map/pointer

---

## KTN-VAR-004 : Préallouer slices avec capacité si connue

### Points positifs
- Évite les réallocations multiples
- Meilleure performance
- Pattern recommandé

### Points négatifs / Problèmes identifiés
1. **Détection complexe** : Comment savoir si la capacité est "connue" ?
2. **Faux positifs** : Capacité approximative vs exacte
3. **Over-allocation** : Capacité trop grande gaspille la mémoire

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Définir "capacité connue" (constante, len(autre), etc.)
- [ ] **AMÉLIORATION** : Limiter aux cas évidents (boucle avec range sur collection de taille connue)

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message avec la capacité suggérée
- [ ] **VÉRIFICATION** : Détecter les patterns `for _, x := range items { result = append(result, ...) }`

### Scénarios non couverts
1. Capacité dépendant d'une condition
2. Slice agrandi progressivement avec critère de sélection

---

## KTN-VAR-005 : Éviter make([]T, length) si utilisation avec append

### Points positifs
- Évite les éléments zéro-value non désirés au début du slice
- Détecte un bug courant

### Points négatifs / Problèmes identifiés
1. **Intention** : Parfois on VEUT les zéro-values initiales
2. **Distinction length vs capacity** : `make([]T, 0, cap)` est correct

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Distinguer `make([]T, n)` (problème) vs `make([]T, 0, n)` (OK)
- [ ] **AMÉLIORATION** : Suggérer le bon pattern `make([]T, 0, len)`

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Expliquer le piège de `make([]T, n)` + `append`

### Scénarios non couverts
1. `make([]T, n)` rempli avec index direct (OK)
2. Slice partiellement rempli intentionnellement

---

## KTN-VAR-006 : Préallouer bytes.Buffer/strings.Builder avec Grow

### Points positifs
- Évite les réallocations dans le buffer
- Meilleure performance pour les grandes chaînes

### Points négatifs / Problèmes identifiés
1. **Taille inconnue** : Souvent on ne connaît pas la taille finale
2. **Over-estimation** : Grow trop grand gaspille la mémoire
3. **Petites chaînes** : Overhead inutile pour petites concaténations

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Quand la règle s'applique-t-elle ? (taille connue uniquement)
- [ ] **AMÉLIORATION** : Ignorer si la taille ne peut pas être estimée

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Suggérer la formule de calcul de taille
- [ ] **CONFIGURATION** : Seuil minimal avant d'exiger Grow

### Scénarios non couverts
1. Buffer utilisé dans une boucle avec taille variable
2. Buffer passé à une fonction qui écrit dedans

---

## KTN-VAR-007 : Utiliser strings.Builder pour >2 concaténations

### Points positifs
- Évite les allocations multiples
- Plus performant que `+` pour plusieurs concaténations
- Recommandation standard Go

### Points négatifs / Problèmes identifiés
1. **Seuil arbitraire** : Pourquoi 2 ? (overhead pour 3 concats peut être négligeable)
2. **fmt.Sprintf** : Parfois plus lisible que Builder
3. **Contexte** : Dans un path chaud vs code rarement exécuté

### Actions à mener

#### Priorité HAUTE
- [ ] **CONFIGURATION** : Rendre le seuil configurable
- [ ] **AMÉLIORATION** : Ignorer les concaténations de constantes (compilateur optimise)

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Suggérer aussi `fmt.Sprintf` comme alternative
- [ ] **DOCUMENTATION** : Benchmarks pour justifier le seuil

### Scénarios non couverts
1. Concaténations de constantes uniquement
2. Code non-performance-critical

---

## KTN-VAR-008 : Éviter les allocations de slices/maps dans les boucles chaudes

### Points positifs
- Détecte un problème de performance courant
- Encourage la réutilisation des allocations

### Points négatifs / Problèmes identifiés
1. **Définition "boucle chaude"** : Comment détecter ?
2. **Faux positifs** : Boucle exécutée rarement
3. **Clear vs réalloc** : Parfois réallouer est plus simple

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Définir "boucle chaude" (toute boucle ? ou heuristique)
- [ ] **AMÉLIORATION** : Suggérer le pattern de clear et réutilisation

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Proposer `clear(slice)` (Go 1.21+) comme alternative
- [ ] **DOCUMENTATION** : Patterns de réutilisation de mémoire

### Scénarios non couverts
1. Boucle avec break/continue fréquent
2. Allocation conditionnelle dans la boucle

---

## KTN-VAR-009 : Utiliser des pointeurs pour les structs >64 bytes

### Points positifs
- Évite les copies coûteuses
- Améliore les performances pour les grandes structs

### Points négatifs / Problèmes identifiés
1. **Calcul de taille** : Comment calculer la taille réelle avec padding ?
2. **Seuil arbitraire** : 64 bytes est-il optimal ?
3. **Sémantique** : Pointeur vs valeur change la sémantique

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Algorithme de calcul de taille (avec padding/alignment)
- [ ] **CONFIGURATION** : Rendre le seuil configurable

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message avec la taille calculée
- [ ] **DOCUMENTATION** : Expliquer quand utiliser pointeur vs valeur

### Scénarios non couverts
1. Structs avec champs dynamiques (slice, map, string)
2. Interfaces (toujours 16 bytes en interne)

---

## KTN-VAR-010 : Utiliser sync.Pool pour les buffers répétés

### Points positifs
- Réduit la pression sur le GC
- Pattern recommandé pour les objets temporaires réutilisables

### Points négatifs / Problèmes identifiés
1. **Détection complexe** : Comment savoir qu'un buffer est "répété" ?
2. **Overhead** : sync.Pool a un coût
3. **Correctness** : Risque de data race si mal utilisé

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Critères de détection (même allocation dans une boucle ?)
- [ ] **AMÉLIORATION** : Template de code sync.Pool dans le message

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Bonnes pratiques sync.Pool
- [ ] **VÉRIFICATION** : Détecter les mauvais usages de Pool

### Scénarios non couverts
1. Buffers de tailles différentes
2. Objets avec état qui doit être reset

---

## KTN-VAR-011 : Shadowing de variables avec := au lieu de =

### Points positifs
- Détecte les bugs subtils de shadowing
- Évite les modifications accidentelles

### Points négatifs / Problèmes identifiés
1. **Shadowing intentionnel** : Parfois voulu pour limiter le scope
2. **err notamment** : `err` shadowed dans les if est très courant
3. **Faux positifs** : Shadowing dans les closures

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Quand le shadowing est-il acceptable ?
- [ ] **AMÉLIORATION** : Exception pour `err` dans les blocs if ?

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message précisant quelle variable est shadowée où
- [ ] **CONFIGURATION** : Liste de variables OK à shadow (`err`, `ctx`, etc.)

### Scénarios non couverts
1. Shadowing dans les goroutines (souvent voulu)
2. Shadowing de paramètres de fonction

---

## KTN-VAR-012 : Conversions string() répétées

### Points positifs
- Évite les allocations multiples
- Encourage la mise en cache

### Points négatifs / Problèmes identifiés
1. **Détection** : Comment détecter "répétées" ?
2. **Scope** : Conversions dans des fonctions différentes
3. **Mutabilité** : Le []byte source peut changer

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Définir "répétées" (même variable, même scope ?)
- [ ] **AMÉLIORATION** : Suggérer la mise en variable

#### Priorité MOYENNE
- [ ] **VÉRIFICATION** : Gérer les conversions dans les boucles
- [ ] **DOCUMENTATION** : Impact performance de string()

### Scénarios non couverts
1. Conversion dans une closure capturée
2. []byte modifié entre les conversions

---

## KTN-VAR-013 : Variables de package groupées dans un seul bloc var

### Points positifs
- Meilleure organisation du code
- Facilite la lecture
- Cohérent avec CONST-002

### Points négatifs / Problèmes identifiés
1. **Catégories** : Parfois plusieurs blocs pour différentes catégories
2. **Génération** : Code généré peut avoir des blocs séparés
3. **Ordre** : Dépendances entre variables

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Un bloc ou plusieurs par catégorie ?
- [ ] **COHÉRENCE** : Aligner avec CONST-002

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Suggérer la fusion avec ordre préservé
- [ ] **DOCUMENTATION** : Convention de groupement

### Scénarios non couverts
1. Variables avec dépendances d'initialisation
2. Variables générées automatiquement

---

## KTN-VAR-014 : Variables de package déclarées après les constantes

### Points positifs
- Ordre cohérent dans les fichiers
- Facilite la navigation
- Complémentaire avec CONST-002

### Points négatifs / Problèmes identifiés
1. **Génération** : Code généré peut ne pas respecter l'ordre
2. **Dépendances** : Variable utilisant une constante (OK car const avant)

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Gérer les fichiers sans constantes (pas de faux positif)
- [ ] **COHÉRENCE** : Aligner avec CONST-002 (const avant var)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Ordre recommandé complet (import → const → var → type → func)

### Scénarios non couverts
1. Fichiers avec types entre const et var
2. Fichiers générés

---

## KTN-VAR-015 : Préallouer maps avec capacité si connue

### Points positifs
- Évite les réallocations de map
- Meilleure performance
- Symétrique avec VAR-004 (slices)

### Points négatifs / Problèmes identifiés
1. **Détection** : Comment savoir si la capacité est "connue" ?
2. **Estimation** : Capacité approximative peut être pire
3. **Maps petites** : Overhead pour petites maps

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Même logique que VAR-004
- [ ] **AMÉLIORATION** : Détecter les patterns évidents (range sur collection)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Différence entre capacité slice et map

### Scénarios non couverts
1. Maps avec clés conditionnellement ajoutées
2. Maps utilisées comme set

---

## KTN-VAR-016 : Utiliser [N]T au lieu de make([]T, N) si N est constant

### Points positifs
- Array sur la stack vs slice sur le heap
- Meilleure performance pour tailles fixes connues

### Points négatifs / Problèmes identifiés
1. **Sémantique différente** : Array vs slice ne sont pas interchangeables
2. **Taille limite** : Gros arrays sur la stack peuvent causer stack overflow
3. **Passage en paramètre** : Array copié, slice non

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Seuil de taille au-delà duquel slice est préférable
- [ ] **AMÉLIORATION** : Ignorer si le slice est passé à une fonction

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Différences array vs slice
- [ ] **VÉRIFICATION** : Détecter les cas où la conversion [:]  est nécessaire

### Scénarios non couverts
1. Slice passé à une fonction acceptant []T
2. Taille > 1MB (stack overflow risk)

---

## KTN-VAR-017 : Copies de mutex (sync.Mutex, sync.RWMutex, atomic.Value)

### Points positifs
- Détecte un bug grave et subtil
- Les mutex ne doivent jamais être copiés
- go vet le détecte mais autant le signaler aussi

### Points négatifs / Problèmes identifiés
1. **Redondant avec go vet** : `go vet` détecte déjà ça
2. **Embedded mutex** : Struct embedant un mutex copié
3. **Atomic types** : Tous les atomic.* ont ce problème

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Couvrir tous les types atomic.* (Value, Int32, etc.)
- [ ] **AMÉLIORATION** : Détecter les copies via assignment et passage de paramètre

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Pourquoi les mutex ne doivent pas être copiés
- [ ] **COHÉRENCE** : Peut-être retirer si go vet suffit, ou améliorer le message

### Scénarios non couverts
1. Mutex dans une struct passée par valeur
2. sync.Cond, sync.Once, sync.WaitGroup (mêmes problèmes)

---

## Interactions entre règles VAR

### Synergies
- **004 + 005 + 015** : Préallocation cohérente (slice et map)
- **006 + 007** : Optimisation des buffers string
- **001 + 002 + 003** : Convention de déclaration cohérente
- **013 + 014** : Organisation du fichier (avec CONST-002)

### Cohérence requise
- **001 vs CONST-003** : Distinction visuelle variable (camelCase) vs constante (SCREAMING_SNAKE_CASE)
- **013 vs CONST-002** : Même logique de groupement
- **003 vs 011** : := peut causer du shadowing

### Potentiels faux positifs
- **008** : Définition de "boucle chaude" floue
- **009** : Seuil 64 bytes arbitraire
- **016** : Array vs slice ont des sémantiques différentes

---

## Résumé des priorités

### Priorité HAUTE
1. **Clarifier** les critères de détection pour VAR-004/008/015 (capacité "connue")
2. **Configurer** les seuils pour VAR-007 (concaténations) et VAR-009 (taille struct)
3. **Améliorer** VAR-011 pour les cas de shadowing légitimes (err)
4. **Vérifier** VAR-017 couvre tous les types sync/atomic

### Priorité MOYENNE
1. Aligner VAR-013/014 avec CONST-002
2. Documenter les patterns de préallocation
3. Gérer les cas array vs slice pour VAR-016

### Priorité BASSE
1. Améliorer les messages avec suggestions de code
2. Rendre les seuils configurables

---

## Tests à vérifier

1. Variable `SCREAMING_CASE` → erreur VAR-001
2. Variable sans type `var x = 5` → erreur VAR-002
3. `var x int` dans une fonction → erreur VAR-003 (suggérer `:=`)
4. `make([]T, len(items))` suivi de `append` → erreur VAR-005
5. `x := y` où `y` est une variable externe → erreur VAR-011 si shadow
6. `sync.Mutex` passé par valeur → erreur VAR-017
