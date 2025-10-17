package rules_var

// Violations intentionnelles pour tester les helpers d'interfaces
// Ce fichier teste les cas limites des variables utilisées avec interfaces

// bad_variable variable mal nommée (snake_case)
var bad_variable string = "test"

// UntypedVar variable sans type explicite dans l'initialisation
var UntypedVar = 42
