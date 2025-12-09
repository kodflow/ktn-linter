# KTN-CONST-002

**Sévérité**: INFO

## Description

Les constantes doivent être :
1. Groupées ensemble dans un seul bloc `const`
2. Placées en haut du fichier, avant toutes les autres déclarations

L'ordre attendu dans un fichier Go est : `const` → `var` → `type` → `func`

## Exemple non conforme

```go
package example

var globalVar = "var"

// ❌ Constante après var
const MaxSize int = 100

type MyType struct{}

// ❌ Constante après type
const Timeout int = 30

func myFunc() {}

// ❌ Constante après func
const ApiKey string = "key"
```

## Exemple conforme

```go
package example

// ✅ Toutes les constantes groupées en haut
const (
    MaxSize int = 100
    Timeout int = 30
    ApiKey string = "key"
)

var globalVar = "var"

type MyType struct{}

func myFunc() {}
```

## Vérifications effectuées

- Plusieurs blocs `const` dispersés sont signalés
- Constantes déclarées après `var` sont signalées
- Constantes déclarées après `type` sont signalées
- Constantes déclarées après `func` sont signalées
