# Tasks - KTN-FUNC

## Vue d'ensemble

Le module `ktnfunc` contient 12 règles (001-012) couvrant la structure, les paramètres, les retours et la complexité des fonctions.

---

## KTN-FUNC-001 : Position de retour de l'erreur (toujours dernière)

### Points positifs
- Suit la convention idiomatique Go
- Facilite le pattern `if err != nil`
- Améliore la cohérence du code

### Points négatifs / Problèmes identifiés
1. **Multiple errors** : Rare, mais une fonction avec 2 erreurs en retour n'est pas gérée
2. **Interfaces externes** : Si une interface externe impose un ordre différent, pas de solution

### Actions à mener
- [ ] **AMÉLIORATION** : Détecter et signaler les fonctions avec plus d'une `error` en retour
- [ ] **DOCUMENTATION** : Mentionner la possibilité d'ignorer via annotation pour interfaces externes
- [ ] **VÉRIFICATION** : Confirmer que les méthodes sont aussi vérifiées

### Scénarios non couverts
- Fonctions sans erreur (ignorées, OK)
- Interfaces externes mal conçues imposant un ordre différent

---

## KTN-FUNC-002 : context.Context en premier paramètre

### Points positifs
- Suit la convention idiomatique Go
- Facilite le passage de contexte dans les appels chaînés

### Points négatifs / Problèmes identifiés
1. **Alias d'import** : Si `context` est importé avec un alias, la détection fonctionne-t-elle ?
2. **Plusieurs contextes** : Très improbable mais non géré

### Actions à mener
- [ ] **VÉRIFICATION CRITIQUE** : S'assurer que la détection utilise l'info de type, pas juste le nom littéral
- [ ] **AMÉLIORATION** : Signaler s'il y a plus d'un `context.Context` dans les paramètres
- [ ] **DOCUMENTATION** : Mentionner `//nolint` pour les rares cas d'interfaces externes

### Scénarios non couverts
- Import avec alias (`import ctx "context"`)
- Paramètres de méthodes (receveur exclu du compte)

---

## KTN-FUNC-003 : Pas de else après un return (early return)

### Points positifs
- Réduit l'indentation et la complexité visuelle
- Encourage les early returns
- Cohérent avec les bonnes pratiques Go

### Points négatifs / Problèmes identifiés
1. **Else-if chaînés** : `if { return } else if {...}` devrait aussi être signalé
2. **Tous les jump statements** : Couvre-t-on `break`, `continue` aussi ?
3. **Cas simples** : Parfois un if/else court est plus lisible

### Actions à mener
- [ ] **VÉRIFICATION** : Confirmer que `break` et `continue` sont aussi couverts
- [ ] **AMÉLIORATION** : Détecter les chaînes `else if` après return
- [ ] **DOCUMENTATION** : Clarifier que la règle est stricte (pas d'exception pour les cas simples)

### Scénarios non couverts
- Chaînes `if { return } else if {...}` (devraient être plates)

---

## KTN-FUNC-004 : Fonctions privées non utilisées (code mort)

### Points positifs
- Élimine le code mort
- Maintient la base de code propre

### Points négatifs / Problèmes identifiés
1. **Usage dans les tests** : Une fonction privée utilisée UNIQUEMENT dans les tests est-elle "utilisée" ?
2. **Réflexion** : Fonctions appelées via reflection ou pointeurs de fonction
3. **Build tags** : Fonctions utilisées sous certains build tags seulement

### Actions à mener
- [ ] **CLARIFICATION CRITIQUE** : Définir si l'usage dans `_test.go` compte comme "utilisée en production"
- [ ] **AMÉLIORATION** : Ignorer explicitement les appels trouvés uniquement dans les fichiers `_test.go`
- [ ] **EXTENSION** : Considérer étendre aux méthodes privées et variables globales non utilisées
- [ ] **DOCUMENTATION** : Mentionner les cas de réflexion/build tags

### Scénarios non couverts
- Fonctions avec usage conditionnel (build tags)
- Fonctions utilisées via réflexion ou passage en paramètre
- Hooks ou callbacks optionnels

---

## KTN-FUNC-005 : Taille maximale d'une fonction (35 lignes)

### Points positifs
- Encourage le découpage en fonctions plus petites
- Améliore la lisibilité et la maintenabilité
- Complémentaire avec KTN-FUNC-011 (complexité)

### Points négatifs / Problèmes identifiés
1. **Méthode de comptage** : Compte-t-on les lignes vides et commentaires ?
2. **Fichiers de test** : Les tests peuvent légitimement être plus longs (setup)
3. **Code généré** : Les fichiers générés peuvent dépasser le seuil

### Actions à mener
- [ ] **CLARIFICATION** : Documenter la méthode de comptage (lignes effectives vs totales)
- [ ] **AMÉLIORATION** : Envisager d'ignorer ou relever le seuil pour les fichiers `_test.go`
- [ ] **AMÉLIORATION** : Ignorer les fichiers avec commentaire `// Code generated`
- [ ] **CONFIGURATION** : Rendre `MAX_FUNCTION_LINES = 35` configurable

### Scénarios non couverts
- Switch avec beaucoup de cases mais logique simple
- Code généré automatiquement

---

## KTN-FUNC-006 : Nombre maximal de paramètres (5)

### Points positifs
- Évite les signatures trop complexes
- Encourage l'utilisation de structs de configuration

### Points négatifs / Problèmes identifiés
1. **Interfaces externes** : Si une interface impose >5 paramètres
2. **Constructeurs** : Les constructeurs ont parfois beaucoup de paramètres

### Actions à mener
- [ ] **VÉRIFICATION** : Confirmer que les méthodes sont vérifiées (receveur exclu)
- [ ] **VÉRIFICATION** : Confirmer que les variadics comptent pour 1 paramètre
- [ ] **DOCUMENTATION** : Recommander l'usage de struct de config quand la limite est atteinte
- [ ] **ANNOTATION** : Permettre `//nolint` pour les cas d'interfaces externes

### Scénarios non couverts
- Interfaces externes imposant >5 paramètres
- Pattern variadic `...opts` comptant pour 1

---

## KTN-FUNC-007 : Pas d'effets de bord dans les getters (Get*/Is*/Has*)

### Points positifs
- Garantit que les getters sont purs
- Prévient les surprises lors de l'appel de getters

### Points négatifs / Problèmes identifiés
1. **Détection des effets de bord** : Comment détecter statiquement ?
2. **Cache/lazy init** : Un getter avec cache modifie l'état interne
3. **Interaction avec STRUCT-003** : Si on n'utilise pas Get*, cette règle ne s'applique pas

### Actions à mener
- [ ] **VÉRIFICATION** : Documenter la méthode de détection des effets de bord
- [ ] **AMÉLIORATION** : Améliorer les messages pour préciser l'effet de bord détecté
- [ ] **RÉFLEXION** : Comment gérer le pattern cache/lazy initialization ?
- [ ] **COHÉRENCE** : Aligner avec STRUCT-003 qui décourage Get*

### Scénarios non couverts
- Pattern cache/lazy init (techniquement un effet de bord mais souvent acceptable)
- Getters appelant des APIs externes ou DB
- Fonctions nommées sans Get* mais qui sont des getters

---

## KTN-FUNC-008 : Paramètres non utilisés (préfixer par _ ou supprimer)

### Points positifs
- Élimine les paramètres inutiles
- Rend explicite l'intention de ne pas utiliser un paramètre

### Points négatifs / Problèmes identifiés
1. **Couvert par le compilateur** : Un paramètre non utilisé ne compile pas (sauf si nommé `_`)
2. **Pattern `_ = param`** : Contournement à détecter
3. **Interfaces** : Impossible de supprimer un paramètre si l'interface l'impose

### Actions à mener
- [ ] **VÉRIFICATION** : Détecter le pattern `_ = param` qui contourne le compilateur
- [ ] **AMÉLIORATION** : Distinguer les méthodes implémentant une interface (ne pas suggérer suppression)
- [ ] **DOCUMENTATION** : Expliquer les deux approches (préfixe `_` vs suppression)

### Scénarios non couverts
- Paramètres utilisés conditionnellement (build tags)
- Paramètres nommés pour la documentation godoc

---

## KTN-FUNC-009 : Pas de nombres magiques (utiliser des constantes)

### Points positifs
- Améliore la lisibilité
- Facilite les modifications
- Renforce KTN-CONST-001/003

### Points négatifs / Problèmes identifiés
1. **Exceptions 0, 1** : Ces valeurs sont souvent légitimes sans constante
2. **Tests** : Les littéraux dans les tests sont souvent acceptables
3. **Contexte** : `for i := 0; i < 10; i++` - le 10 est-il magique ?

### Actions à mener
- [ ] **AMÉLIORATION CRITIQUE** : Ignorer 0 et 1 (valeurs universelles)
- [ ] **AMÉLIORATION** : Ignorer ou assouplir dans les fichiers `_test.go`
- [ ] **AMÉLIORATION** : Ignorer les littéraux dans les boucles locales simples
- [ ] **VÉRIFICATION** : Ne pas signaler les littéraux dans les chaînes de caractères

### Scénarios non couverts
- Littéraux dans les logs ou fmt.Printf
- Valeurs de test
- Compteurs de boucle locaux

---

## KTN-FUNC-010 : Interdiction des naked returns (sauf fonctions < 5 lignes)

### Points positifs
- Améliore la lisibilité des returns explicites
- Réduit la charge cognitive

### Points négatifs / Problèmes identifiés
1. **Tension avec FUNC-012** : On encourage les named returns (012) mais on interdit naked returns
2. **Seuil < 5** : 5 lignes pile est interdit ? (clarifier)

### Actions à mener
- [ ] **CLARIFICATION** : Préciser si < 5 signifie ≤ 4 ou ≤ 5
- [ ] **DOCUMENTATION** : Expliquer l'interaction avec FUNC-012 (nommer pour documenter, mais return explicite)
- [ ] **CONFIGURATION** : Rendre le seuil configurable

### Scénarios non couverts
- Fonctions exactement à la limite (5 lignes)
- Fonctions `init()` courtes

---

## KTN-FUNC-011 : Complexité cyclomatique maximale (10)

### Points positifs
- Limite la complexité des fonctions
- Complémentaire avec FUNC-005 (taille)

### Points négatifs / Problèmes identifiés
1. **Calcul précis** : Tous les branchements sont-ils comptés ? (&&, ||, case, etc.)
2. **Tests** : Les tests avec beaucoup de cas peuvent être complexes
3. **Switch avec beaucoup de cases** : Peut être clair malgré haute complexité

### Actions à mener
- [ ] **VÉRIFICATION** : Documenter précisément la méthode de calcul
- [ ] **AMÉLIORATION** : Considérer d'ignorer ou relever le seuil pour `_test.go`
- [ ] **AMÉLIORATION** : Ajouter des suggestions de refactoring dans le message
- [ ] **CONFIGURATION** : Rendre le seuil configurable

### Scénarios non couverts
- Switch avec 11+ cases mais logique simple
- Tests table-driven avec beaucoup de cas

---

## KTN-FUNC-012 : Plus de 3 retours => utiliser des retours nommés

### Points positifs
- Améliore la documentation des retours multiples
- Clarifie l'intention de chaque valeur retournée

### Points négatifs / Problèmes identifiés
1. **Alternative struct** : >3 retours suggère plutôt un struct de résultat
2. **Tension avec FUNC-010** : Named returns permettent naked returns qu'on interdit
3. **Interfaces externes** : Si une interface impose >3 retours sans noms

### Actions à mener
- [ ] **DOCUMENTATION** : Expliquer l'interaction avec FUNC-010 (nommer mais return explicite)
- [ ] **AMÉLIORATION** : Suggérer dans le message de considérer un struct de résultat
- [ ] **VÉRIFICATION** : Gérer le cas où l'erreur est déjà nommée `err`

### Scénarios non couverts
- Interfaces externes imposant >3 retours sans noms
- Cas où un struct serait plus approprié

---

## Interactions entre règles FUNC

### Synergies
- **005 + 011** : Taille + complexité = fonctions simples et courtes
- **001 + 002** : Ordre cohérent des paramètres (context) et retours (error)
- **003 + 011** : Early returns réduisent la complexité

### Tensions à gérer
- **010 vs 012** : Named returns encouragés mais naked returns interdits → documenter que les deux sont compatibles
- **004 vs TEST-004** : Fonction privée inutilisée en prod mais testée → clarifier le comportement

---

## Résumé des priorités

### Priorité HAUTE
1. **Ignorer 0 et 1** dans FUNC-009 (magic numbers)
2. **Clarifier** l'usage dans les tests pour FUNC-004 (code mort)
3. **Ignorer les tests** pour FUNC-005 (taille) et FUNC-011 (complexité)

### Priorité MOYENNE
1. **Détecter le pattern `_ = param`** dans FUNC-008
2. **Documenter la méthode de calcul** pour FUNC-011
3. **Clarifier l'interaction** FUNC-010 vs FUNC-012

### Priorité BASSE
1. Rendre les seuils configurables
2. Améliorer les messages avec suggestions de refactoring
