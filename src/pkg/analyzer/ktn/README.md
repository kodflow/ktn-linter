# KTN - Kodflow Technical Norms

Toutes les règles de style et de bonnes pratiques pour Go.

## Vue d'ensemble

Le linter KTN contient **~64 règles** organisées en **16 catégories** pour garantir un code Go propre, maintenable et performant.

## Catégories

| Catégorie | Règles | Description |
|-----------|--------|-------------|
| [func](./func/) | 9 | Fonctions et méthodes |
| [var](./var/) | 9 | Variables |
| [struct](./struct/) | 4 | Structures |
| [interface](./interface/) | 6 | Interfaces |
| [const](./const/) | 4 | Constantes |
| [error](./error/) | 1 | Gestion d'erreurs |
| [test](./test/) | 4 | Tests unitaires |
| [alloc](./alloc/) | 3 | Allocations mémoire |
| [goroutine](./goroutine/) | 2 | Concurrence |
| [pool](./pool/) | 1 | Object pooling |
| [mock](./mock/) | 2 | Mocks |
| [method](./method/) | 1 | Méthodes |
| [package](./package/) | 1 | Packages |
| [control_flow](./control_flow/) | 7 | Structures de contrôle |
| [data_structures](./data_structures/) | 3 | Structures de données |
| [ops](./ops/) | 8 | Opérations diverses |

## Utilisation

```go
import "github.com/kodflow/ktn-linter/pkg/analyzer/ktn"

// Obtenir toutes les règles
all := ktn.GetAllRules()

// Obtenir les règles par catégorie
funcRules := ktn.GetRulesByCategory("func")
errorRules := ktn.GetRulesByCategory("error")
```

## Architecture

```
analyzer/ktn/
├── func/           # 9 règles de fonctions
│   ├── 001.go      # Nommage MixedCaps
│   ├── 002.go      # Documentation godoc
│   ├── ...
│   ├── registry.go
│   └── README.md
├── error/          # 1 règle cruciale
│   ├── 001.go      # Wrapping d'erreurs
│   ├── registry.go
│   └── README.md
├── var/            # 9 règles de variables
├── ...
├── registry.go     # Registre global
└── README.md       # Ce fichier
```

## Règles prioritaires

### 🔥 Critique
- **KTN-ERROR-001**: Wrapping d'erreurs obligatoire
- **KTN-FUNC-002**: Documentation godoc
- **KTN-VAR-005**: Pas de globales mutables

### ⚠️ Important
- **KTN-FUNC-006**: Longueur de fonction (max 35 lignes)
- **KTN-FUNC-007**: Complexité cyclomatique (max 10)
- **KTN-GOROUTINE-001**: Context pour cancellation

## Compilation

```bash
# Compiler toutes les règles
go build ./analyzer/ktn/...

# Compiler une catégorie spécifique
go build ./analyzer/ktn/func/...
go build ./analyzer/ktn/error/...
```

## Tests

```bash
# Tester toutes les catégories
go test ./analyzer/ktn/...

# Tester une catégorie
go test ./analyzer/ktn/func/...
```
