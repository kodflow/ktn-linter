# KTN Rules - Good Usage Examples

Ce dossier contient des exemples de **code conforme aux règles KTN**.

## Qu'est-ce que KTN ?

KTN (Kodflow Technical Norms) est un ensemble de règles de style et de best practices **plus strictes que la spec Go officielle**.

Ces règles visent à :
- Améliorer la lisibilité du code
- Faciliter la maintenance
- Standardiser les patterns
- Réduire les bugs potentiels
- Améliorer les performances

## Différence avec gospec/

- **gospec/** : Code qui respecte uniquement la spécification Go officielle (syntaxe valide)
- **ktn/** : Code qui respecte les conventions KTN (idiomatique, optimisé, documenté)

## Catégories de règles

### Allocation (rules_alloc/)
- `KTN-ALLOC-001` : Interdire new() avec types référence (utiliser make())
- `KTN-ALLOC-002` : make([]T, 0) suivi d'append doit spécifier une capacité
- `KTN-ALLOC-004` : Préférer &struct{} à new(struct)

### Constantes (rules_const/)
- `KTN-CONST-001` : Constantes groupées dans const ()
- `KTN-CONST-002` : Commentaire de groupe pour const
- `KTN-CONST-003` : Commentaire individuel pour chaque const
- `KTN-CONST-004` : Types explicites pour const publiques

### Variables (rules_var/)
- `KTN-VAR-001` à `KTN-VAR-009` : Conventions de déclaration, nommage, grouping

### Fonctions (rules_func/)
- `KTN-FUNC-001` à `KTN-FUNC-010` : Nommage, documentation, complexité, longueur

### Erreurs (rules_error/)
- `KTN-ERROR-001` : Wrapping obligatoire des erreurs avec %w

### Goroutines (rules_goroutine/)
- `KTN-GOROUTINE-001` : Goroutines dans boucles avec synchronisation
- `KTN-GOROUTINE-002` : Goroutines avec mécanisme de sync obligatoire

### Et autres...
- rules_pool/ : sync.Pool patterns
- rules_struct/ : Structures et fields
- rules_interface/ : Architecture avec interfaces
- rules_mock/ : Génération de mocks
- rules_test/ : Organisation des tests

## Usage

Ces fichiers servent de référence pour :
1. Comprendre les règles KTN
2. Tester le linter ktn-linter
3. Documenter les best practices

Chaque fichier de test démontre un pattern correct selon KTN.
