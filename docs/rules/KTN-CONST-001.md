# KTN-CONST-001

**Sévérité**: WARNING

## Description

Les constantes doivent avoir un type explicite pour éviter les conversions inattendues et clarifier l'intention du développeur.

## Exemple non conforme

```go
const (
    MaxSize = 100        // Type implicite
    ApiKey = "secret"    // Type implicite
    Timeout = 5.0        // Type implicite
)
```

## Exemple conforme

```go
const (
    MaxSize int = 100
    ApiKey string = "secret"
    Timeout float64 = 5.0
)
```

## Cas particulier : iota

Pour les constantes utilisant `iota`, seule la première constante du bloc nécessite un type explicite. Les constantes suivantes héritent automatiquement du type.

```go
const (
    StatusPending Status = iota  // Type explicite requis
    StatusRunning                // Hérite de Status
    StatusCompleted              // Hérite de Status
)
```
