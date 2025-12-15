// Package messages provides structured error messages for KTN rules.
// This file contains API rule messages for dependency/coupling rules.
package messages

// registerAPIMessages enregistre les messages API.
func registerAPIMessages() {
	Register(Message{
		Code:  "KTN-API-001",
		Short: "paramètre '%s' utilise type concret '%s'; suggérer interface '%s' avec: %s",
		Verbose: `PROBLÈME: Le paramètre '%s' est typé par '%s' (type concret externe).

POURQUOI C'EST UN PROBLÈME:
  - Couplage fort avec une implémentation externe
  - Difficile à tester (impossible de mocker)
  - Viole le principe d'inversion des dépendances

SOLUTION: Définir une interface minimale côté consumer:

  type %s interface {
      %s
  }

  func %s(%s %s) { ... }

AVANTAGES:
  - Code testable (injection de mocks)
  - Découplage (ne dépend que du comportement nécessaire)
  - Principe ISP (Interface Segregation Principle)

LIMITES V1:
  - y := x; y.Method() non détecté (variable intermédiaire)
  - T.Method(x) non détecté (expression de méthode)`,
	})
}
