# Tasks - KTN-CONST

## Vue d'ensemble

Le module `ktnconst` contient 3 règles pour la déclaration et le nommage des constantes.

---

## KTN-CONST-001 : Constantes à type explicite

### Points positifs
- Évite les conversions inattendues avec les constantes non typées
- Clarifie l'intention du développeur
- Complémentaire avec KTN-CONST-003 (nommage)

### Points négatifs / Problèmes identifiés
1. **Cas idiomatiques Go** : `const Timeout = 5 * time.Second` fonctionne naturellement sans type explicite
2. **Verbosité** : Forcer `const Timeout time.Duration = 5 * time.Second` est plus verbeux
3. **Iota** : Les constantes `iota` héritent du type de la première, mais la règle exige le type sur chacune ?

### Actions à mener
- [ ] **VÉRIFICATION** : Confirmer que seule la première constante d'un bloc `iota` nécessite le type
- [ ] **AMÉLIORATION** : Tolérer les cas où l'expression donne déjà un type évident (`time.Second`)
- [ ] **DOCUMENTATION** : Documenter clairement le comportement avec `iota`
- [ ] **VÉRIFICATION** : S'assurer que les constantes locales (dans les fonctions) sont aussi vérifiées

### Scénarios non couverts
- Constantes non exportées (le besoin de type explicite est moins crucial)
- Expressions qui ont déjà un type implicite clair

---

## KTN-CONST-002 : Regrouper les constantes en haut de fichier (avant les variables)

### Points positifs
- Améliore la lisibilité en groupant les définitions immuables
- Cohérent avec KTN-VAR-014 (variables après constantes)
- Facilite la navigation dans le code

### Points négatifs / Problèmes identifiés
1. **Constantes dispersées** : La règle détecte-t-elle les `const` définis après des fonctions ?
2. **Plusieurs blocs const** : Faut-il fusionner en un seul bloc ou juste les avoir avant `var` ?
3. **Ordre avec types/fonctions** : La règle ne vérifie pas la position relative aux `type` ou `func`

### Actions à mener
- [ ] **VÉRIFICATION** : Confirmer que les `const` définis après des fonctions sont signalés
- [ ] **AMÉLIORATION** : Signaler si plusieurs blocs `const` séparés existent (recommander fusion)
- [ ] **DOCUMENTATION** : Clarifier que l'ordre `import` → `const` → `var` → `type` → `func` est souhaité
- [ ] **VÉRIFICATION** : Gérer le cas où il n'y a pas de `var` (pas de faux positif)

### Scénarios non couverts
- Fichiers avec plusieurs blocs `const` séparés (pour différentes catégories)
- Position relative aux déclarations de types

---

## KTN-CONST-003 : Format de nom des constantes (SCREAMING_SNAKE_CASE)

### Points positifs
- Distinction visuelle immédiate entre constantes et variables
- Convention uniforme dans le projet
- Facilite la recherche de constantes dans le code

### Points négatifs / Problèmes identifiés
1. **Non-standard Go** : La convention Go idiomatique utilise CamelCase pour les constantes
2. **Export forcé** : SCREAMING_SNAKE_CASE commence par majuscule → constante toujours exportée
3. **Pas de constante privée possible** : Impossible d'avoir une constante interne avec cette convention
4. **Conflits potentiels** : Les développeurs Go habitués à CamelCase seront surpris

### Actions à mener
- [ ] **DOCUMENTATION CRITIQUE** : Documenter clairement ce choix non-standard et ses implications
- [ ] **RÉFLEXION** : Décider si les constantes internes doivent suivre un autre pattern (ex: `_INTERNAL_CONST`)
- [ ] **VÉRIFICATION** : S'assurer que les acronymes sont bien gérés (`MAX_HTTP_CONN` OK)
- [ ] **COMMUNICATION** : Préparer une explication pour les nouveaux contributeurs

### Scénarios non couverts
- Constantes d'un seul mot (`ERROR`) - techniquement conforme
- Constantes qui doivent rester internes au package (impossible avec cette convention)

---

## Interactions entre règles CONST

### Synergies
- **001 + 003** : Type explicite + nom SCREAMING_SNAKE_CASE = déclaration complète et lisible
- **002 + VAR-014** : Ordre cohérent const → var dans tout le projet

### Conflits potentiels
- **003 vs conventions Go** : Le format SCREAMING_SNAKE_CASE est contraire aux conventions officielles
- **003 export forcé** : Toute constante sera exportée, pas de constante privée possible

---

## Résumé des priorités

### Priorité HAUTE
1. **Documenter** le choix non-standard de SCREAMING_SNAKE_CASE (003)
2. **Clarifier** le comportement avec `iota` (001)
3. **Vérifier** la détection des `const` dispersés après les fonctions (002)

### Priorité MOYENNE
1. Signaler les blocs `const` multiples à fusionner (002)
2. Gérer les expressions avec type implicite évident (001)

### Priorité BASSE
1. Réflexion sur les constantes internes (003)
2. Rendre les seuils/patterns configurables

---

## Décisions architecturales à prendre

### Question : Constantes internes
Avec la convention SCREAMING_SNAKE_CASE, toutes les constantes sont exportées. Options :
1. **Accepter** : Toutes les constantes sont publiques (choix actuel)
2. **Pattern alternatif** : `_INTERNAL_CONST` pour les internes (underscore initial)
3. **Assouplir** : Autoriser camelCase pour les constantes non exportées

### Recommandation
Option 1 ou 2 selon la philosophie du projet. Documenter clairement le choix.
