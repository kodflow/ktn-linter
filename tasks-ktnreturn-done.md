# Tasks - KTN-RETURN

## Vue d'ensemble

Le module `ktnreturn` contient 1 règle (002) pour la gestion des valeurs de retour.

> **Note** : Il n'y a pas de KTN-RETURN-001. La numérotation commence à 002.

---

## KTN-RETURN-002 : Préférer slice/map vide à nil

### Points positifs
- Évite aux appelants de gérer le cas `nil`
- Permet d'itérer directement sur le résultat sans vérification
- Simplifie le code appelant
- Suit les recommandations de style Go

### Points négatifs / Problèmes identifiés

#### 1. Légère allocation vs nil
- `[]T{}` crée un objet slice (même si minime)
- `nil` ne crée rien
- Impact négligeable mais existe

#### 2. Distinction sémantique nil vs vide
Certaines APIs utilisent intentionnellement cette distinction :
- `nil` = "pas de résultat" ou "non applicable"
- `[]T{}` = "résultat vide mais valide"

Cette pratique est déconseillée mais existe.

#### 3. Cohérence slice vs map
Les deux devraient être traités de la même manière :
- `return nil` pour slice → suggérer `return []T{}`
- `return nil` pour map → suggérer `return make(map[K]V)`

### Actions à mener

#### Priorité HAUTE
- [ ] **VÉRIFICATION** : Confirmer que les maps sont couvertes en plus des slices
- [ ] **AMÉLIORATION** : Ajouter un quickfix suggérant la bonne syntaxe (`[]T{}` ou `make(map[K]V)`)

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Limiter la règle aux fonctions exportées (API publique) pour réduire le bruit
- [ ] **DOCUMENTATION** : Expliquer pourquoi nil vs vide n'est pas une bonne distinction sémantique

#### Priorité BASSE
- [ ] **RÉFLEXION** : Considérer une exception si la fonction a une documentation explicite sur nil

### Scénarios non couverts

1. **Distinction intentionnelle nil vs vide**
   - L'API utilise volontairement nil comme signal
   - Devra ignorer la règle ou repenser l'API

2. **Performance critique**
   - Cas très rare où l'allocation de `[]T{}` importerait
   - En pratique négligeable, mais signaler au développeur

3. **Maps**
   - Vérifier que la règle couvre `return nil` pour les maps
   - Comportement identique à nil pour itération mais différent pour insertion

---

## Détails techniques

### Comportement Go avec nil vs vide

```go
// Les deux fonctionnent pour l'itération
var nilSlice []int          // nil
emptySlice := []int{}       // non-nil, len=0

for _, v := range nilSlice {    // OK, 0 itérations
    // ...
}
for _, v := range emptySlice {  // OK, 0 itérations
    // ...
}

// Les deux fonctionnent pour len/cap
len(nilSlice)  // 0
len(emptySlice) // 0

// Maps aussi
var nilMap map[string]int
emptyMap := make(map[string]int)

for k, v := range nilMap {     // OK, 0 itérations
    // ...
}
len(nilMap) // 0

// MAIS pour l'insertion, nil map panique
nilMap["key"] = 1  // PANIC!
emptyMap["key"] = 1  // OK
```

### Recommandation
La règle est bien fondée : retourner une collection vide évite des bugs potentiels et simplifie le code appelant.

---

## Interactions avec autres règles

### Complémentarité
- **avec FUNC-001** : Erreur en dernière position + slice vide = pattern propre
- **avec VAR-008** : Éviter les allocations en boucle, mais une allocation de slice vide au return est acceptable

### Aucun conflit identifié

---

## Résumé des priorités

### Priorité HAUTE
1. **Confirmer** que les maps sont couvertes
2. **Ajouter** des quickfix avec la bonne syntaxe

### Priorité MOYENNE
1. Considérer de limiter aux fonctions exportées
2. Documenter les cas edge

### Priorité BASSE
1. Optimisation du message d'erreur

---

## Tests à ajouter/vérifier

1. `return nil` pour slice → erreur
2. `return nil` pour map → erreur
3. `return []T{}` pour slice → OK
4. `return make(map[K]V)` pour map → OK
5. Fonction interne vs exportée → même traitement actuellement
6. Fonction avec plusieurs returns, certains nil → tous signalés
