# KTN-INTERFACE-001

**Sévérité**: WARNING

## Description

Les interfaces privées (minuscule) déclarées doivent être utilisées dans le package.

Les interfaces exportées (majuscule) sont ignorées car elles font partie de l'API publique.

## Exemple non-conforme

```go
type unusedInterface interface { // want "KTN-INTERFACE-001"
    Method()
}
// Interface privée jamais utilisée
```

## Exemple conforme

```go
type usedInterface interface {
    Method()
}

func Process(i usedInterface) {
    i.Method()
}
```
