package rules_const

// Violations intentionnelles pour tester les helpers d'interfaces
// Ce fichier teste les cas limites des constantes utilisées avec interfaces

// bad_constant constante mal nommée (snake_case)
const bad_constant string = "test"

// UntypedConst constante sans type explicite
const UntypedConst = 42
