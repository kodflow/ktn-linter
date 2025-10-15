# Tests Target - Code Parfait

Ce répertoire contient **l'archétype du code PARFAIT** - exemples de code conforme à TOUTES les règles KTN-Linter.

## 📊 Statistiques

- **0 violation** détectée par le linter ✅
- Couvre **TOUTES** les règles KTN implémentées
- Exemples de code réel montrant les bonnes pratiques

## 🎯 Objectif

`tests/target/` sert de **référence de qualité** :
- Code production-ready
- Documentation complète
- Best practices Go appliquées
- Design patterns recommandés

## ✅ Catégories de Code Parfait

### CONST - Constantes Optimales
✅ Groupées dans des blocs `const ()`
✅ Commentaire de groupe obligatoire
✅ Commentaire individuel pour chaque constante
✅ Type explicite (sauf iota)
✅ Nommage en MixedCaps

### FUNC - Fonctions Exemplaires
✅ Godoc complet avec Params/Returns
✅ Profondeur d'imbrication ≤ 3
✅ Longueur ≤ 35 lignes
✅ Complexité cyclomatique < 10
✅ Commentaires internes pour logique complexe
✅ Commentaires sur returns multiples

### INTERFACE - Design Interface-First
✅ Fichier `interfaces.go` dans chaque package
✅ Types publics exposés comme interfaces
✅ Implémentations privées
✅ Constructeurs `New*()` obligatoires
✅ Godoc complet sur interfaces et méthodes

### VAR - Variables Structurées
✅ Groupées dans des blocs `var ()`
✅ Commentaire de groupe obligatoire
✅ Type explicite
✅ Pas de shadowing
✅ Conversion en const si jamais réassignée

### TEST - Tests Complets
✅ Package naming `package_test`
✅ Fichier `*_test.go` pour chaque fichier source
✅ Godoc avec section Params sur tests
✅ Couverture complète

## 📁 Structure

```
tests/target/
├── README.md                          # Ce fichier
├── rules_const/                       # Constantes parfaites
│   ├── ktn_const_001_grouped.go       # Constantes groupées
│   ├── ktn_const_002_comments.go      # Avec commentaires
│   └── ktn_const_edge_iota/           # Edge cases (iota)
├── rules_func/                        # Fonctions exemplaires
│   ├── ktn_func_002_godoc.go          # Godoc complet
│   ├── ktn_func_006_length.go         # Longueur optimale
│   ├── ktn_func_010_nesting_depth.go  # Profondeur contrôlée
│   └── ktn_func_edge_*/               # Edge cases
├── rules_interface/                   # Interface-first design
│   ├── ktn_interface_001_file/        # Avec interfaces.go
│   ├── ktn_interface_002_public/      # Interfaces publiques
│   ├── ktn_interface_006_constructor/ # Avec New*()
│   └── ktn_interface_edge_*/          # Edge cases
├── rules_var/                         # Variables optimales
│   ├── ktn_var_001_grouped.go         # Variables groupées
│   ├── ktn_var_002_comments.go        # Avec commentaires
│   └── ktn_var_edge_*/                # Edge cases
└── rules_test/                        # Tests conformes
    ├── KTN-TEST-001-package-naming/   # Package correct
    └── KTN-TEST-002-*/                # Avec tests
```

## 🎓 Exemples de Bonnes Pratiques

### Constantes Groupées et Documentées
```go
// Configuration constants.
// Define application limits and defaults.
const (
    // MaxConnections nombre maximum de connexions simultanées.
    MaxConnections int = 100

    // DefaultTimeout timeout par défaut en secondes.
    DefaultTimeout int = 30
)
```

### Fonction avec Godoc Complet
```go
// ProcessUser traite les données utilisateur et les valide.
//
// Params:
//   - user: les données utilisateur à traiter
//   - options: options de traitement
//
// Returns:
//   - *Result: résultat du traitement
//   - error: erreur si la validation échoue
func ProcessUser(user *User, options ProcessOptions) (*Result, error) {
    // Validation des données
    if err := validateUser(user); err != nil {
        return nil, err
    }

    // Traitement
    result := &Result{
        Status: "processed",
        User:   user,
    }

    return result, nil
}
```

### Interface-First Design
```go
// Package myservice fournit les services métier.
package myservice

// Service définit l'interface du service.
type Service interface {
    Process(data string) error
    GetStatus() string
}

// service implémentation privée.
type service struct {
    db Database
}

// NewService crée une nouvelle instance de Service.
//
// Params:
//   - db: base de données à utiliser
//
// Returns:
//   - Service: instance du service
func NewService(db Database) Service {
    return &service{db: db}
}

// Process implémente Service.Process.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func (s *service) Process(data string) error {
    return s.db.Save(data)
}
```

### Profondeur d'Imbrication Contrôlée
```go
// ProcessData traite les données avec early returns.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func ProcessData(data map[string]interface{}) error {
    // Early return si pas de données
    if data == nil {
        return errors.New("no data")
    }

    // Extraction avec helper
    cfg, err := extractConfig(data)
    if err != nil {
        return err
    }

    // Traitement avec helper
    return processConfig(cfg)
}

// extractConfig extrait la configuration.
func extractConfig(data map[string]interface{}) (*Config, error) {
    // Profondeur limitée grâce à la fonction dédiée
    if val, ok := data["config"]; ok {
        if cfg, ok := val.(*Config); ok {
            return cfg, nil
        }
    }
    return nil, errors.New("invalid config")
}
```

## 🔗 Voir Aussi

- [tests/source/](../source/) - Code FOIREUX avec 405 violations
- [Documentation KTN-Linter](../../README.md)
- [FUNCTION_BEST_PRACTICES.md](../../FUNCTION_BEST_PRACTICES.md)

## ✨ Utilisation

```bash
# Vérifier que le code est parfait
./builds/ktn-linter ./tests/target/...

# Résultat attendu: ✅ No issues found! Code is compliant.
```

**Ce code est la référence** - Utilisez-le comme modèle pour vos propres projets !
