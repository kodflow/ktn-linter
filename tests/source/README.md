# Tests Source - Anti-Patterns Catalog

Ce répertoire contient **l'archétype du code FOIREUX** - tous les anti-patterns que le linter KTN doit détecter.

## 📊 Statistiques

- **405 violations** détectées par le linter
- Couvre **TOUTES** les règles KTN implémentées
- Exemples de code réel montrant ce qu'il NE FAUT PAS faire

## 🔴 Catégories d'Anti-Patterns

### CONST - Constantes mal déclarées
- ❌ Constantes individuelles (KTN-CONST-001)
- ❌ Groupes sans commentaire (KTN-CONST-002)
- ❌ Constantes sans documentation (KTN-CONST-003)

### FUNC - Fonctions problématiques
- ❌ **Trop de paramètres** (8-10 params au lieu de ≤5)
- ❌ **Profondeur d'imbrication extrême** (6-8 niveaux au lieu de ≤3)
- ❌ **Fonctions trop longues** (60+ lignes au lieu de ≤35)
- ❌ **Godoc absent ou incomplet** (sans Params/Returns)
- ❌ Commentaires internes manquants
- ❌ Documentation inadéquate

### INTERFACE - Design anti-pattern
- ❌ **Structs publics exposés** (UserService, PaymentGateway, etc.)
- ❌ **Interfaces sans constructeurs** New*()
- ❌ Packages sans fichier interfaces.go
- ❌ Types publics qui devraient être des interfaces

### VAR - Variables anarchiques
- ❌ **Variables déclarées individuellement** (au lieu de var())
- ❌ **Variables qui devraient être const** (valeurs littérales jamais modifiées)
- ❌ **Shadowing extrême** (même nom à 4-5 niveaux d'imbrication)
- ❌ Variables mal groupées
- ❌ Déclarations non optimales

### TEST - Tests inadéquats
- ❌ Fichiers sans tests correspondants
- ❌ Tests mal packagés
- ❌ Godoc manquant sur les tests

## 📁 Fichiers Clés d'Anti-Patterns

### Violations FUNC
```
rules_func/bad_func_too_many_params.go       # 8-10 paramètres
rules_func/bad_func_extreme_nesting.go       # Profondeur 6-8
rules_func/bad_func_too_long.go              # 60+ lignes
rules_func/bad_func_no_godoc.go              # Sans documentation
```

### Violations INTERFACE
```
rules_interface/bad_interface_public_structs.go    # 5 structs publics
rules_interface/bad_interface_no_constructor.go   # 7 interfaces sans New*
```

### Violations VAR
```
rules_var/bad_var_individual.go           # 10 vars individuelles
rules_var/bad_var_should_be_const.go      # 10 vars → const
rules_var/bad_var_extreme_shadowing.go    # Shadowing profond
```

### Violations TEST
```
rules_test/bad_no_test_file.go           # Fichier sans *_test.go
```

## 🎯 Utilisation

```bash
# Lancer le linter sur les anti-patterns
./builds/ktn-linter ./tests/source/...

# Résultat attendu: 405 violations détectées
# ✅ Prouve que le linter fonctionne correctement
```

## ⚠️ AVERTISSEMENT

**CE CODE EST VOLONTAIREMENT MAUVAIS !**

- Ne JAMAIS copier/coller ce code dans un projet réel
- Utiliser comme référence "à ne pas faire"
- Comparer avec tests/target/ pour voir le bon code

## 📖 Exemples de Violations

### Fonction avec trop de paramètres
```go
// ❌ MAUVAIS - 8 paramètres
func ProcessUserDataBad(name string, email string, age int,
    country string, city string, zipcode string,
    phone string, newsletter bool) error

// ✅ BON - Utiliser un struct
type UserData struct {
    Name, Email, Country, City, Zipcode, Phone string
    Age int
    Newsletter bool
}
func ProcessUserData(data UserData) error
```

### Profondeur d'imbrication excessive
```go
// ❌ MAUVAIS - 8 niveaux
func ProcessDataWithExtremeNesting(data map[string]interface{}) error {
    if data != nil {              // 1
        if val, ok := data["config"]; ok {  // 2
            if cfg, ok := val.(map[string]interface{}); ok {  // 3
                if enabled, ok := cfg["enabled"]; ok {  // 4
                    if enabled == true {  // 5
                        if settings, ok := cfg["settings"]; ok {  // 6
                            if s, ok := settings.([]interface{}); ok {  // 7
                                for _, item := range s {  // 8
                                    _ = item
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// ✅ BON - Extraire des fonctions, early returns
```

### Struct public au lieu d'interface
```go
// ❌ MAUVAIS - Struct exposé
type UserService struct {
    db      string
    cache   map[string]interface{}
}

// ✅ BON - Interface publique, impl privée
type UserService interface {
    GetUser(id int) (User, error)
}

type userService struct {
    db    string
    cache map[string]interface{}
}

func NewUserService() UserService {
    return &userService{}
}
```

## 🔗 Voir Aussi

- [tests/target/](../target/) - Code PARFAIT conforme à toutes les règles
- [Documentation KTN-Linter](../../README.md)
