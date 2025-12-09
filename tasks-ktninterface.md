# Tasks - KTN-INTERFACE

## Vue d'ensemble

Le module `ktninterface` contient 1 règle (001) pour la gestion des interfaces.

---

## KTN-INTERFACE-001 : Interface non utilisée

### Points positifs
- Élimine le code mort (interfaces fantômes)
- Maintient la base de code propre
- Cohérent avec FUNC-004 (fonctions inutilisées)

### Points négatifs / Problèmes identifiés

#### 1. Conflit majeur avec STRUCT-001
**C'est le problème principal de cette règle.**

- **STRUCT-001** exige : "Chaque struct doit avoir une interface associée"
- **INTERFACE-001** signale : "Interface non utilisée"

**Scénario de conflit :**
```go
// Pour satisfaire STRUCT-001, on crée l'interface
type UserServiceInterface interface {
    GetUser(id int) (*User, error)
}

type UserService struct{}

func (s *UserService) GetUser(id int) (*User, error) { ... }
```

Si `UserServiceInterface` n'est jamais utilisée en paramètre ou variable, INTERFACE-001 la signale comme inutilisée, alors qu'on l'a créée POUR satisfaire STRUCT-001.

#### 2. Usage uniquement dans les tests
Une interface créée pour le mocking en tests mais pas utilisée en prod sera signalée.

#### 3. Interfaces API publiques
Dans une bibliothèque, on peut exposer une interface destinée aux utilisateurs externes sans l'utiliser en interne.

### Actions à mener

#### Priorité CRITIQUE
- [ ] **RÉSOUDRE LE CONFLIT STRUCT-001** : Options :
  1. **Option A** : Ignorer les interfaces qui suivent le pattern `XxxInterface` pour une struct `Xxx` du même package
  2. **Option B** : Considérer qu'une interface est "utilisée" si une struct du package l'implémente
  3. **Option C** : Exiger que le constructeur `NewXxx()` retourne l'interface (pas la struct concrète)

#### Priorité HAUTE
- [ ] **AMÉLIORATION** : Détecter si une interface est implémentée implicitement par une struct du package
- [ ] **AMÉLIORATION** : Ignorer les interfaces utilisées uniquement dans les fichiers `_test.go` (ou les considérer comme utilisées)

#### Priorité MOYENNE
- [ ] **AMÉLIORATION** : Ignorer les interfaces exportées dans les packages de type "library" (détection par convention de nommage ?)
- [ ] **DOCUMENTATION** : Expliquer clairement l'interaction avec STRUCT-001

### Scénarios non couverts

1. **Interface créée pour STRUCT-001 mais jamais utilisée**
   - Statut : CONFLIT → à résoudre impérativement

2. **Interface utilisée uniquement dans les tests (mocking)**
   - L'analyse sans les tests verra l'interface comme inutilisée
   - Décision : l'interface devrait-elle être dans un fichier `_test.go` ?

3. **Interface d'API publique pour utilisateurs externes**
   - L'analyse intra-package ne voit pas l'usage externe
   - Rare dans un code applicatif, plus fréquent dans les libs

4. **Interface vide exportée intentionnellement**
   - Ex: `type Serializable interface{}` pour marquer des types
   - Sera signalée comme inutilisée

---

## Résolution proposée du conflit STRUCT-001 / INTERFACE-001

### Recommandation : Option C (Constructeur retourne l'interface)

```go
// ✅ Cette approche satisfait les deux règles

type UserServiceInterface interface {
    GetUser(id int) (*User, error)
}

type UserService struct{}

// Le constructeur retourne l'interface → elle est "utilisée"
func NewUserService() UserServiceInterface {
    return &UserService{}
}
```

**Avantages :**
- L'interface est utilisée (retour de constructeur)
- Encourage le pattern d'injection de dépendances
- Facilite le mocking

**Implémentation requise :**
- INTERFACE-001 considère une interface comme utilisée si elle est type de retour d'une fonction
- Documenter cette best practice dans le README

---

## Interactions avec autres règles

### Conflit direct
- **vs STRUCT-001** : Résolution requise (voir ci-dessus)

### Complémentarité
- **avec FUNC-004** : Même philosophie d'élimination du code mort
- **avec STRUCT-002** : Si le constructeur retourne l'interface, les deux règles sont satisfaites

---

## Résumé des priorités

### Priorité CRITIQUE
1. **RÉSOUDRE** le conflit avec STRUCT-001

### Priorité HAUTE
1. Définir ce qui constitue un "usage" d'interface (retour de fonction, paramètre, variable)
2. Gérer le cas des interfaces utilisées uniquement en tests

### Priorité MOYENNE
1. Documenter les best practices
2. Améliorer les messages d'erreur avec suggestions

---

## Tests à ajouter

1. Interface créée pour une struct (pattern XxxInterface) → pas d'erreur si constructeur la retourne
2. Interface utilisée uniquement en test → comportement défini
3. Interface implémentée mais jamais référencée → comportement défini
4. Interface avec méthodes identiques à celles d'une struct → comportement défini
