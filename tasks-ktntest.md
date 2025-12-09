# Tasks - KTN-TEST

## Vue d'ensemble

Le module `ktntest` contient 13 règles (001-013) pour la structure et la qualité des tests.

---

## KTN-TEST-001 : Fichier de test doit se terminer par _internal_test.go ou _external_test.go

### Points positifs
- Distinction claire entre tests white-box et black-box
- Facilite la compréhension de la portée des tests
- Convention cohérente dans tout le projet

### Points négatifs / Problèmes identifiés
1. **Contrainte forte** : Beaucoup de projets utilisent simplement `_test.go`
2. **Migration** : Projets existants devront renommer beaucoup de fichiers
3. **Outils externes** : Certains outils ne reconnaissent que `*_test.go`

### Actions à mener

#### Priorité HAUTE
- [ ] **DOCUMENTATION** : Expliquer clairement la différence internal/external
- [ ] **VÉRIFICATION** : S'assurer que `go test` fonctionne avec ces suffixes

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message d'erreur avec suggestion de renommage
- [ ] **MIGRATION** : Script de migration pour renommer les fichiers existants

### Scénarios non couverts
1. Fichiers de benchmark (`*_bench_test.go`)
2. Fichiers de test d'intégration (`*_integration_test.go`)

---

## KTN-TEST-002 : Les fichiers de test doivent utiliser le package xxx_test

### Points positifs
- Encourage les tests black-box
- Teste l'API publique comme un utilisateur externe
- Détecte les problèmes d'encapsulation

### Points négatifs / Problèmes identifiés
1. **Conflit avec TEST-011** : TEST-011 distingue internal (package xxx) vs external (package xxx_test)
2. **Tests internes nécessaires** : Parfois il faut tester des fonctions privées

### Actions à mener

#### Priorité CRITIQUE
- [ ] **COHÉRENCE** : Aligner avec TEST-001 et TEST-011 (internal vs external)
- [ ] **RÉFLEXION** : Cette règle semble redondante avec TEST-011

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Clarifier quand utiliser package xxx vs xxx_test

### Scénarios non couverts
1. Tests internes légitimes de fonctions privées
2. Fichiers `_internal_test.go` qui DOIVENT utiliser `package xxx`

---

## KTN-TEST-003 : Chaque fichier _test.go doit avoir un fichier .go correspondant

### Points positifs
- Évite les tests orphelins
- Structure 1:1 entre code et tests
- Facilite la navigation

### Points négatifs / Problèmes identifiés
1. **Tests d'intégration** : Un test peut tester plusieurs fichiers
2. **Fichiers de helpers** : `helpers_test.go` ou `utils_test.go` n'ont pas de source
3. **Conflit avec TEST-006** : Semble être la même règle

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Différence entre TEST-003 et TEST-006 ?
- [ ] **AMÉLIORATION** : Ignorer les fichiers helpers de test (`testutils_test.go`, etc.)

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Exceptions acceptables (integration tests, etc.)

### Scénarios non couverts
1. Fichiers `testdata_test.go` ou `fixtures_test.go`
2. Tests d'intégration globaux

---

## KTN-TEST-004 : Toutes les fonctions (publiques et privées) doivent avoir des tests

### Points positifs
- Couverture complète du code
- Élimine le code non testé
- Qualité garantie

### Points négatifs / Problèmes identifiés
1. **Fonctions triviales** : Getters/setters simples méritent-ils un test dédié ?
2. **Code généré** : Fonctions générées automatiquement
3. **Fonctions main** : `func main()` dans les packages main
4. **Détection complexe** : Comment matcher fonction ↔ test ?

### Actions à mener

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Ignorer `func main()` et `func init()`
- [ ] **AMÉLIORATION** : Ignorer les fichiers avec commentaire `// Code generated`
- [ ] **VÉRIFICATION** : Détecter les tests par convention de nommage (`TestXxx` pour `Xxx`)

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Tolérer les getters/setters triviaux sans test dédié
- [ ] **CONFIGURATION** : Seuil configurable de couverture minimale

#### Priorité BASSE
- [ ] **STATISTIQUES** : Afficher le % de fonctions testées

### Scénarios non couverts
1. Fonction testée indirectement via une autre fonction
2. Fonction dans un fichier avec build tags
3. Méthodes d'interface implémentées

---

## KTN-TEST-005 : Les tests avec plusieurs cas doivent utiliser table-driven tests

### Points positifs
- Pattern idiomatique Go
- Facile d'ajouter de nouveaux cas
- Réduit la duplication de code de test

### Points négatifs / Problèmes identifiés
1. **Seuil de détection** : À partir de combien de cas ?
2. **Tests séquentiels** : Certains tests doivent être séquentiels (état partagé)
3. **Faux positifs** : Tests avec plusieurs assertions mais cas unique

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Définir le seuil (2+ assertions similaires ?)
- [ ] **VÉRIFICATION** : Distinguer assertions sur le même cas vs cas multiples

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Template de table-driven test recommandé
- [ ] **AMÉLIORATION** : Suggestion de refactoring dans le message

### Scénarios non couverts
1. Tests avec sous-tests (`t.Run`) déjà structurés
2. Tests de benchmark avec plusieurs configurations

---

## KTN-TEST-006 : Chaque fichier _test.go doit correspondre à un fichier source (pattern 1:1)

### Points positifs
- Organisation claire
- Navigation facile
- Responsabilité unique par fichier

### Points négatifs / Problèmes identifiés
1. **Doublon avec TEST-003** : Semble identique
2. **Fichiers helpers** : `testhelper_test.go` légitime
3. **Tests d'intégration** : Testent plusieurs sources

### Actions à mener

#### Priorité CRITIQUE
- [ ] **CLARIFICATION** : Fusionner ou différencier de TEST-003

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Liste d'exceptions (testhelper, testutil, etc.)

### Scénarios non couverts
1. Voir TEST-003 (mêmes scénarios)

---

## KTN-TEST-007 : Interdiction d'utiliser t.Skip() dans les tests

### Points positifs
- Évite les tests "oubliés" indéfiniment
- Force à corriger ou supprimer les tests cassés
- Détecte la dette technique de tests

### Points négatifs / Problèmes identifiés
1. **Cas légitimes** : Tests dépendant de l'environnement (CI, OS, etc.)
2. **Tests lents** : Skip en développement local
3. **Dépendances externes** : Tests nécessitant une DB, API externe

### Actions à mener

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Autoriser `t.Skip` avec commentaire explicatif
- [ ] **AMÉLIORATION** : Ignorer `testing.Short()` qui est idiomatique

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Alternatives à t.Skip (build tags, etc.)
- [ ] **VÉRIFICATION** : Couvrir aussi `t.SkipNow()`

### Scénarios non couverts
1. `t.Skip` conditionnel (`if testing.Short()`)
2. Tests avec dépendances d'environnement

---

## KTN-TEST-008 : Chaque fichier .go doit avoir les fichiers de test appropriés

### Points positifs
- Garantit la couverture de tests
- Lié à TEST-001 (suffixes appropriés)
- Structure complète

### Points négatifs / Problèmes identifiés
1. **Fichiers sans fonctions exportées** : Ne nécessitent pas de `_external_test.go`
2. **Fichiers générés** : Code généré sans tests
3. **Fichiers de types** : `types.go` avec seulement des types

### Actions à mener

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Ignorer les fichiers avec seulement des types/constantes
- [ ] **AMÉLIORATION** : Ignorer les fichiers générés
- [ ] **VÉRIFICATION** : Ne demander `_internal_test.go` que si fonctions privées

#### Priorité MOYENNE
- [ ] **CONFIGURATION** : Liste de fichiers à ignorer configurable

### Scénarios non couverts
1. Fichiers `doc.go` (documentation package)
2. Fichiers `version.go` avec seulement des constantes

---

## KTN-TEST-009 : Les tests de fonctions publiques doivent être dans _external_test.go

### Points positifs
- Tests black-box pour l'API publique
- Teste comme un utilisateur externe
- Vérifie l'encapsulation

### Points négatifs / Problèmes identifiés
1. **Détection des tests** : Comment lier `TestFoo` à `Foo` ?
2. **Tests combinés** : Un test peut tester plusieurs fonctions
3. **Overlap avec TEST-010** : Complémentaires mais peuvent être confus

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Algorithme de matching test ↔ fonction
- [ ] **DOCUMENTATION** : Convention de nommage `TestNomFonction`

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message indiquant quelle fonction devrait être dans quel fichier

### Scénarios non couverts
1. Test `TestFooAndBar` testant deux fonctions
2. Tests de méthodes (`TestType_Method`)

---

## KTN-TEST-010 : Les tests de fonctions privées doivent être dans _internal_test.go

### Points positifs
- Tests white-box séparés
- Accès aux fonctions privées
- Structure claire

### Points négatifs / Problèmes identifiés
1. **Complément de TEST-009** : Peuvent être confus ensemble
2. **Fonctions privées légitimement testées via l'API publique**

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Aligner avec TEST-009
- [ ] **VÉRIFICATION** : Convention de nommage pour tests de fonctions privées

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Quand tester directement vs via l'API publique

### Scénarios non couverts
1. Fonction privée testée indirectement (coverage suffisante)

---

## KTN-TEST-011 : Les _internal_test.go doivent utiliser package xxx, les _external_test.go doivent utiliser package xxx_test

### Points positifs
- Cohérence avec TEST-001/009/010
- Enforce le pattern internal/external
- Convention claire

### Points négatifs / Problèmes identifiés
1. **Redondance** : TEST-002 existe déjà pour les packages
2. **Stricte** : Erreur si mauvaise combinaison

### Actions à mener

#### Priorité HAUTE
- [ ] **COHÉRENCE** : Vérifier que TEST-002 et TEST-011 ne sont pas en conflit
- [ ] **RÉFLEXION** : TEST-002 devrait peut-être être supprimée en faveur de TEST-011

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Message d'erreur clair avec la combinaison attendue

### Scénarios non couverts
- Voir TEST-002

---

## KTN-TEST-012 : Les tests doivent contenir des assertions et vraiment tester quelque chose

### Points positifs
- Évite les tests vides ou triviaux
- Garantit que les tests vérifient quelque chose
- Qualité des tests

### Points négatifs / Problèmes identifiés
1. **Détection** : Comment détecter une "vraie" assertion ?
2. **Tests d'effets de bord** : Certains tests vérifient juste l'absence d'erreur
3. **Faux positifs** : `if err != nil { t.Fatal(err) }` compte comme assertion ?

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Définir ce qui compte comme assertion
  - `t.Error`, `t.Errorf`, `t.Fatal`, `t.Fatalf`
  - `require.*`, `assert.*` (testify)
  - Comparaisons avec `if x != expected`
- [ ] **AMÉLIORATION** : Reconnaître les patterns d'assertion courants

#### Priorité MOYENNE
- [ ] **CONFIGURATION** : Liste de packages d'assertion reconnus

### Scénarios non couverts
1. Tests utilisant des frameworks d'assertion non reconnus
2. Tests de panic (`defer func() { recover() }`)

---

## KTN-TEST-013 : Les tests doivent couvrir les cas d'erreur et exceptions

### Points positifs
- Robustesse des tests
- Vérifie le happy path ET les erreurs
- Meilleure couverture

### Points négatifs / Problèmes identifiés
1. **Détection complexe** : Comment savoir si les erreurs sont testées ?
2. **Fonctions sans erreur** : Pas toutes les fonctions retournent des erreurs
3. **Faux positifs** : Peut être trop strict

### Actions à mener

#### Priorité HAUTE
- [ ] **CLARIFICATION** : Comment détecter qu'un cas d'erreur est testé ?
- [ ] **AMÉLIORATION** : Limiter aux fonctions qui retournent `error`

#### Priorité MOYENNE
- [ ] **DOCUMENTATION** : Bonnes pratiques pour tester les erreurs
- [ ] **STATISTIQUES** : Pourcentage de chemins d'erreur testés

### Scénarios non couverts
1. Erreurs internes (panics récupérés)
2. Erreurs dans des callbacks/handlers

---

## Interactions entre règles TEST

### Synergies
- **001 + 009 + 010 + 011** : Structure complète internal/external
- **003 + 006** : Pattern 1:1 (mais potentiellement redondantes)
- **004 + 008** : Couverture complète

### Redondances à clarifier
- **TEST-002 vs TEST-011** : Même sujet (package de test)
- **TEST-003 vs TEST-006** : Même sujet (fichier correspondant)

### Tensions
- **TEST-007** : Interdiction stricte de t.Skip peut être trop contraignante

---

## Résumé des priorités

### Priorité CRITIQUE
1. **Clarifier** les redondances TEST-002/011 et TEST-003/006
2. **Définir** les algorithmes de matching fonction ↔ test

### Priorité HAUTE
1. Ignorer `func main()` et `func init()` pour TEST-004
2. Autoriser `t.Skip` avec conditions légitimes pour TEST-007
3. Définir ce qui compte comme "assertion" pour TEST-012

### Priorité MOYENNE
1. Améliorer les messages avec suggestions concrètes
2. Documenter les conventions de nommage
3. Gérer les fichiers helpers de test

### Priorité BASSE
1. Scripts de migration
2. Configuration des seuils

---

## Tests à vérifier

1. Fichier `foo.go` → doit avoir `foo_internal_test.go` et/ou `foo_external_test.go`
2. Fichier `_internal_test.go` avec `package xxx_test` → erreur TEST-011
3. Test sans assertion → erreur TEST-012
4. Fonction publique testée dans `_internal_test.go` → erreur TEST-009
