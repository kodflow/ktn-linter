# Tasks - KTN-COMMENT

## Vue d'ensemble

Le module `ktncomment` contient 7 règles de commentaires (001-007) visant à garantir une documentation complète et cohérente du code.

---

## KTN-COMMENT-001 : Longueur maximale des commentaires inline (80 caractères)

### Points positifs
- Encourage des commentaires concis et lisibles
- Ignore correctement les commentaires de documentation
- Complémentaire avec KTN-COMMENT-007

### Points négatifs / Problèmes identifiés
1. **Commentaires multilignes mal gérés** : Un bloc `/* ... */` est traité comme une seule chaîne. La règle pourrait signaler un bloc multi-ligne dont la somme dépasse 80, même si chaque ligne est courte
2. **Constante non configurable** : `MAX_COMMENT_LENGTH = 80` est fixe, devrait être paramétrable
3. **URLs longues** : Les URLs dans les commentaires dépassent souvent 80 caractères mais sont légitimes

### Actions à mener
- [ ] **AMÉLIORATION** : Vérifier la longueur de chaque ligne séparément dans les blocs `/* ... */`
- [ ] **AMÉLIORATION** : Ajouter une exception pour les lignes contenant des URLs (regex `https?://`)
- [ ] **REFACTORING** : Rendre `MAX_COMMENT_LENGTH` configurable via constante documentée
- [ ] **VÉRIFICATION** : S'assurer que les fichiers `_test.go` sont exclus de cette règle

### Scénarios non couverts
- Commentaires de documentation au début de fichier (avant `package`)
- URLs ou messages d'erreur longs légitimes

---

## KTN-COMMENT-002 : Commentaire de package requis (min 3 caractères)

### Points positifs
- Assure une documentation minimale de chaque package
- Aligne avec les conventions Go pour les packages exportés

### Points négatifs / Problèmes identifiés
1. **Portée trop stricte** : Exige un commentaire par fichier au lieu de par package
2. **Fichiers tests inclus** : Les fichiers `_test.go` n'ont généralement pas besoin de commentaire de package
3. **Package main** : Une description n'est pas toujours pertinente pour `main`
4. **Seuil trop bas** : 3 caractères minimum permet des commentaires inutiles comme `// X`

### Actions à mener
- [ ] **REFACTORING MAJEUR** : Changer la logique pour exiger UN commentaire de package par package (pas par fichier)
- [ ] **AMÉLIORATION** : Ignorer les fichiers `_test.go`
- [ ] **AMÉLIORATION** : Ignorer ou assouplir pour `package main`
- [ ] **AMÉLIORATION** : Relever `MIN_PACKAGE_COMMENT_LENGTH` à 10+ caractères
- [ ] **VÉRIFICATION** : Gérer le cas licence/copyright en en-tête (détecter le commentaire juste avant `package`)

### Scénarios non couverts
- Fichiers avec licence en en-tête mais sans description du package
- Packages avec plusieurs fichiers (duplication de commentaires)

---

## KTN-COMMENT-003 : Commentaire sur chaque constante

### Points positifs
- Encourage la documentation des symboles globaux
- Cohérent avec les autres règles de documentation (004, 005)

### Points négatifs / Problèmes identifiés
1. **Constantes groupées** : Dans un bloc `const (...)`, faut-il un commentaire par constante ou un commentaire de groupe suffit-il ?
2. **Constantes non exportées** : Documenter chaque constante interne peut être lourd
3. **Constantes évidentes** : `const Pi = 3.14` n'a pas besoin de commentaire détaillé

### Actions à mener
- [ ] **CLARIFICATION** : Définir la politique pour les blocs `const (...)` - accepter un commentaire de groupe ?
- [ ] **AMÉLIORATION** : Considérer limiter l'obligation aux constantes exportées uniquement
- [ ] **AMÉLIORATION** : Gérer les énumérations `iota` - un commentaire de bloc suffit pour le groupe
- [ ] **DOCUMENTATION** : Documenter clairement le comportement attendu dans le README

### Scénarios non couverts
- Constantes "évidentes" qui génèrent des commentaires redondants
- Groupes `iota` qui représentent un ensemble logique

---

## KTN-COMMENT-004 : Commentaire sur chaque variable de package

### Points positifs
- Cohérent avec KTN-COMMENT-003 pour les constantes
- Assure la documentation des éléments globaux

### Points négatifs / Problèmes identifiés
1. **Variables groupées** : Même problème que pour les constantes avec les blocs `var (...)`
2. **Variables internes triviales** : `var maxRetries = 3` parle d'elle-même
3. **Surcharge potentielle** : Tout commenter peut créer du bruit

### Actions à mener
- [ ] **ALIGNEMENT** : Harmoniser le comportement avec KTN-COMMENT-003 pour les blocs groupés
- [ ] **AMÉLIORATION** : Considérer limiter aux variables exportées
- [ ] **VÉRIFICATION** : S'assurer que les variables d'erreur globales (`var ErrX = ...`) sont bien couvertes

### Scénarios non couverts
- Variables de métrique/compteurs auto-explicatives
- Variables groupées dans un bloc `var (...)`

---

## KTN-COMMENT-005 : Documentation des structs exportées (2 lignes)

### Points positifs
- Suit les conventions Go officielles
- Garantit une documentation minimale de l'API publique
- Complémentaire avec KTN-STRUCT-002

### Points négatifs / Problèmes identifiés
1. **Structs non exportées** : Aucune vérification pour les structs internes complexes
2. **Qualité du contenu** : Vérifie juste la présence, pas si le commentaire commence par le nom du type

### Actions à mener
- [ ] **AMÉLIORATION** : Vérifier que le commentaire commence par le nom de la struct (style Go officiel)
- [ ] **OPTIONNEL** : Ajouter un warning de moindre sévérité pour les structs internes complexes
- [ ] **VÉRIFICATION** : Gérer les structs générées automatiquement (protobuf, etc.)

### Scénarios non couverts
- Structs générées par des outils (protobuf, go generate)
- Structs embedded anonymement

---

## KTN-COMMENT-006 : Documentation de toutes les fonctions

### Points positifs
- Assure une documentation intégrale
- Cohérent avec KTN-TEST-004 (tester toutes les fonctions)

### Points négatifs / Problèmes identifiés
1. **Fonctions internes** : Documenter chaque fonction privée peut être excessif
2. **Fonctions de test** : Exiger un commentaire sur chaque `TestXxx` alourdit les tests
3. **Fonctions triviales** : Getters simples génèrent des commentaires paraphrasant le nom

### Actions à mener
- [ ] **AMÉLIORATION MAJEURE** : Ignorer les fichiers `_test.go` pour cette règle
- [ ] **AMÉLIORATION** : Considérer limiter aux fonctions/méthodes exportées (niveau ERROR pour exportées, WARNING pour internes)
- [ ] **VÉRIFICATION** : S'assurer que les méthodes sont aussi vérifiées
- [ ] **DOCUMENTATION** : Clarifier les attentes pour les fonctions triviales

### Scénarios non couverts
- Fonctions de test (TestXxx, BenchmarkXxx)
- Fonctions triviales où le nom est suffisamment explicite

---

## KTN-COMMENT-007 : Commenter les blocs de contrôle et logique significative

### Points positifs
- Encourage la documentation du code complexe
- Complémentaire avec KTN-COMMENT-001 (commentaires concis)

### Points négatifs / Problèmes identifiés
1. **Trop zélée** : Exige un commentaire même pour `if err != nil { return err }` évident
2. **Définition floue** : Qu'est-ce que "logique significative" ?
3. **Position du commentaire** : Tolère-t-on un commentaire à l'intérieur du bloc plutôt qu'avant ?
4. **Risque de bruit** : Trop de commentaires peut nuire à la lisibilité

### Actions à mener
- [ ] **AMÉLIORATION MAJEURE** : Ignorer les blocs très simples (1 ligne, pattern `if err != nil { return err }`)
- [ ] **AMÉLIORATION** : Définir des exceptions pour les patterns évidents
- [ ] **VÉRIFICATION** : S'assurer que les commentaires à l'intérieur du bloc sont détectés
- [ ] **DOCUMENTATION** : Clarifier la position souhaitée des commentaires (avant le bloc)

### Scénarios non couverts
- Commentaires positionnés à l'intérieur du bloc plutôt qu'avant
- Blocs `default:` dans un switch trivial
- Boucles très courtes

---

## Résumé des priorités

### Priorité HAUTE
1. Ignorer les fichiers `_test.go` pour les règles 002, 006, 007
2. Améliorer la gestion des blocs multi-lignes dans 001
3. Ignorer les blocs simples `if err != nil { return err }` dans 007

### Priorité MOYENNE
1. Harmoniser la gestion des blocs groupés `const/var (...)` dans 003/004
2. Vérifier la qualité du commentaire (commence par le nom) dans 005
3. Ajouter exception pour URLs longues dans 001

### Priorité BASSE
1. Rendre les constantes configurables
2. Différencier les niveaux de sévérité (exporté vs interne)
