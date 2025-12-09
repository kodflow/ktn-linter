# Tasks - KTN-STRUCT

## Vue d'ensemble

Le module `ktnstruct` contient 7 règles (001-007) pour la structure et l'organisation des structs.

---

## KTN-STRUCT-001 : Interface pour chaque struct (100% méthodes)

### Points positifs
- Encourage la programmation contre des interfaces
- Facilite le mocking en tests
- Prépare la future interchangeabilité des implémentations
- Suit le principe d'inversion de dépendances

### Points négatifs / Problèmes identifiés

#### 1. Conflit avec INTERFACE-001
**Problème majeur déjà identifié dans tasks-ktninterface.md**
- STRUCT-001 exige une interface
- INTERFACE-001 la signale comme inutilisée si pas référencée

#### 2. Structs sans méthodes
La règle devrait ignorer les structs purement "data" (DTOs).

#### 3. Structs implémentant des interfaces standard
Ex: une struct qui implémente `fmt.Stringer` → exiger une interface spécifique est redondant

#### 4. Boilerplate excessif
Chaque struct avec méthodes = 1 interface à créer et maintenir

### Actions à mener

#### Priorité CRITIQUE
- [ ] **RÉSOUDRE** le conflit avec INTERFACE-001 (voir tasks-ktninterface.md)
- [ ] **VÉRIFICATION** : Confirmer que les structs sans méthodes sont ignorées

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Tolérer les structs implémentant seulement des interfaces standard (Stringer, error, etc.)
- [ ] **AMÉLIORATION** : Vérifier par contenu (mêmes méthodes) pas seulement par nom

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Convention de nommage recommandée (XxxInterface ou IXxx ?)
- [ ] **AMÉLIORATION** : Autoriser l'interface dans un ordre de méthodes différent

### Scénarios non couverts
1. Struct implémentant seulement `fmt.Stringer` → interface redondante
2. Deux structs avec mêmes méthodes → interfaces dupliquées vs interface partagée
3. Struct embedant une autre struct avec méthodes → interface requise ?

---

## KTN-STRUCT-002 : Struct exportée avec méthodes => constructeur obligatoire

### Points positifs
- Garantit une initialisation correcte
- Encapsule les invariants
- Pattern factory clair

### Points négatifs / Problèmes identifiés
1. **Nommage strict** : Le constructeur doit-il s'appeler exactement `NewXxx` ?
2. **Structs simples** : `type Point struct { X, Y int }` avec une méthode triviale → constructeur vraiment nécessaire ?
3. **Constructeurs dans autre package** : Rare mais possible

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Documenter le pattern de détection (nom `NewXxx`, retour `*Xxx` ou `Xxx`)
- [ ] **AMÉLIORATION** : Tolérer plusieurs constructeurs (`NewXxxWithOption`, etc.)

#### Priorité MOYENNE
- [ ] **COHÉRENCE** : Le constructeur devrait retourner l'interface (lien avec INTERFACE-001)
- [ ] **DOCUMENTATION** : Best practice = constructeur retourne l'interface, pas la struct concrète

#### Priorité BASSE
- [ ] **RÉFLEXION** : Exceptions pour structs très simples avec 1-2 champs publics ?

### Scénarios non couverts
1. Constructeur nommé `CreateXxx` ou `MakeXxx` → non détecté
2. Constructeur dans un package séparé → non détecté
3. Struct simple initialisable littéralement → constructeur semble superflu

---

## KTN-STRUCT-003 : Ne pas préfixer les getters par "Get"

### Points positifs
- Suit la convention Go idiomatique (Effective Go)
- Code plus concis
- Cohérent avec la stdlib Go

### Points négatifs / Problèmes identifiés
1. **Interaction avec FUNC-007** : FUNC-007 cherche Get*/Is*/Has* pour vérifier les effets de bord
2. **Fonctions libres** : La règle ne couvre que les méthodes, pas les fonctions globales

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Aligner avec FUNC-007 qui utilise Get* pour détecter les getters
- [ ] **VÉRIFICATION** : Couvre-t-on aussi `get` minuscule (méthode non exportée) ?

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Étendre aux fonctions globales `GetXxx()` si pertinent
- [ ] **DOCUMENTATION** : Rappeler la convention Go (Name() au lieu de GetName())

#### Priorité BASSE
- [ ] **EXCEPTION** : Interfaces externes imposant `GetXxx()` (gRPC, etc.)

### Scénarios non couverts
1. Méthode `GetDataAndClear()` qui fait plus qu'un get → signalée mais c'est correct
2. Interfaces externes imposant `Get*` → devra être ignoré

---

## KTN-STRUCT-004 : Un seul struct par fichier

### Points positifs
- Meilleure organisation du code
- Facilite la navigation
- Un fichier = un concept

### Points négatifs / Problèmes identifiés
1. **Structs privées liées** : Parfois logique de garder une struct helper privée avec sa struct principale
2. **Fichiers de test** : Les tests peuvent définir plusieurs structs de mock
3. **Fragmentation excessive** : Package avec 10 structs = 10+ fichiers

### Actions à mener

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Ignorer les fichiers `_test.go` (structs de mock acceptables)
- [ ] **CLARIFICATION** : L'interface associée (STRUCT-001) peut-elle être dans le même fichier ? (Oui car interface ≠ struct)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Convention fichier = struct + interface + constructeur
- [ ] **RÉFLEXION** : Tolérer 1 struct privée helper avec la struct principale ?

### Scénarios non couverts
1. Deux structs très liées (parent/enfant) → fichiers séparés obligatoires
2. Fichiers de test avec structs de mock → à ignorer

---

## KTN-STRUCT-005 : Champs exportés avant champs non exportés

### Points positifs
- API publique visible immédiatement
- Convention visuelle claire
- Facilite la documentation

### Points négatifs / Problèmes identifiés
1. **Champs embedded anonymes** : Un champ embedded de type exporté est-il considéré "exporté" ?
2. **Impact sur les tags** : Réordonner les champs peut affecter l'ordre d'encodage (JSON, etc.)

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Traiter les champs embedded anonymes correctement
- [ ] **AMÉLIORATION** : Message précisant quels champs sont mal ordonnés

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Mentionner l'impact potentiel sur l'ordre d'encodage (négligeable en pratique)

### Scénarios non couverts
1. Champs embedded anonymes de type exporté
2. Struct sans champs exportés ou sans champs privés → rien à vérifier

---

## KTN-STRUCT-006 : Pas de tags sur les champs privés des DTOs

### Points positifs
- Les champs privés ne sont pas sérialisés, donc les tags sont inutiles
- Évite la confusion
- Code plus propre

### Points négatifs / Problèmes identifiés
1. **Définition de DTO** : Comment détecter qu'une struct est un DTO ?
   - Par suffixe (DTO, Request, Response) ?
   - Par présence de tags sur des champs ?

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Documenter les critères de détection d'un DTO
- [ ] **AMÉLIORATION** : S'assurer que la détection est cohérente avec STRUCT-007

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Expliquer pourquoi les tags sur champs privés sont inutiles

### Scénarios non couverts
1. Struct avec tag `-` (ignore) sur champ privé → techniquement inutile aussi
2. Struct non-DTO avec tags sur champs publics ET champs privés → seulement privés signalés

---

## KTN-STRUCT-007 : Getters obligatoires pour champs privés (non-DTO)

### Points positifs
- Encapsulation propre
- Permet l'évolution future (logique dans getter)
- Interface complète pour la struct

### Points négatifs / Problèmes identifiés
1. **Interaction avec STRUCT-003** : Le getter doit s'appeler `Field()` pas `GetField()`
2. **Définition non-DTO** : Inverse de STRUCT-006, doit être cohérent
3. **Getters triviaux** : Beaucoup de code boilerplate pour juste `return s.field`

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Aligner la détection DTO avec STRUCT-006
- [ ] **VÉRIFICATION** : Le getter attendu est `Field()` pas `GetField()` (STRUCT-003)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Convention de nommage claire (champ `name` → getter `Name()`)
- [ ] **AMÉLIORATION** : Message indiquant le nom de getter attendu

### Scénarios non couverts
1. Champ dont on ne veut PAS de getter (vraiment privé) → doit ignorer la règle
2. Struct avec méthodes mais pas de champs privés → rien à vérifier

---

## Interactions entre règles STRUCT

### Synergies
- **001 + 002** : Interface + Constructeur = API propre
- **003 + 007** : Getters sans "Get" + getters obligatoires = encapsulation idiomatique
- **006 + 007** : Détection DTO cohérente entre les deux

### Conflit majeur
- **001 vs INTERFACE-001** : Interface requise mais potentiellement signalée comme inutilisée

### Cohérence requise
- **006 et 007** : Même détection de DTO (suffixes, tags, etc.)
- **003 et 007** : Même convention de nommage des getters

---

## Résumé des priorités

### Priorité CRITIQUE
1. **RÉSOUDRE** le conflit STRUCT-001 / INTERFACE-001

### Priorité HAUTE
1. Ignorer les tests pour STRUCT-004
2. Aligner la détection DTO entre 006 et 007
3. Traiter les champs embedded pour 005

### Priorité MOYENNE
1. Documenter les conventions de nommage
2. Gérer les structs implémentant des interfaces standard pour 001
3. Constructeur retournant l'interface pour 002

### Priorité BASSE
1. Exceptions pour structs très simples
2. Amélioration des messages d'erreur
