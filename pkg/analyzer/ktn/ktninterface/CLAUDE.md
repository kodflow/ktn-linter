# pkg/analyzer/ktn/ktninterface/ - Interface Rules

## Purpose
Analyze interface declarations for usage, design, naming, and Go idioms.

## Rules (3 total)
| Rule | Description | Severity |
|------|-------------|----------|
| KTN-INTERFACE-001 | Interface privée non utilisée | Warning |
| KTN-INTERFACE-003 | Single-method interfaces should follow -er naming | Info |
| KTN-INTERFACE-004 | Overuse of empty interface (interface{}/any) | Info |

## Go Interface Philosophy
"The bigger the interface, the weaker the abstraction." - Rob Pike

Interfaces belong in the consumer package, not the producer.

## File Structure
```
ktninterface/
├── 001.go              # KTN-INTERFACE-001 implementation
├── 001_external_test.go
├── registry.go         # Analyzers()
└── testdata/src/interface001/
    ├── good.go
    └── bad.go
```

## KTN-INTERFACE-001: Unused Private Interface

Détecte les interfaces privées (minuscule) qui ne sont jamais utilisées dans le package.

**Comportement:**
- Interfaces exportées (majuscule): ignorées (API publique)
- Interfaces privées (minuscule): reportées si non utilisées

**Usage considéré:**
- Champ de struct
- Paramètre de fonction/méthode
- Retour de fonction/méthode
- Déclaration de variable
- Compile-time check (`var _ MyInterface = (*S)(nil)`)

## Testdata Example
```go
// bad.go
type unusedInterface interface { // want "KTN-INTERFACE-001"
    Method()
}
// Interface privée jamais utilisée

// good.go
type usedInterface interface {
    Method()
}
func process(i usedInterface) { ... }
```

## Composition Over Large Interfaces
```go
// Instead of one large interface, compose small ones
type ReadWriter interface {
    Reader
    Writer
}
```
