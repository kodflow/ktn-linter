// Package messages provides structured error messages for KTN rules.
// This file contains API rule messages for dependency/coupling rules.
package messages

// registerAPIMessages enregistre les messages API.
func registerAPIMessages() {
	Register(Message{
		Code:  "KTN-API-001",
		Short: "paramètre '%s' dépend du type concret externe '%s'; utiliser une interface minimale avec {%s}",
		Verbose: `PROBLÈME: Le paramètre '%s' est typé par '%s' (type concret externe).

POURQUOI C'EST UN PROBLÈME:
  - Couplage fort avec une implémentation externe
  - Difficile à tester (impossible de mocker)
  - Viole le principe d'inversion des dépendances

SOLUTION: Définir une interface minimale côté consumer:

  // Interface minimale contenant uniquement les méthodes utilisées
  type %s interface {
      %s
  }

  // Modifier la signature pour accepter l'interface
  func %s(%s %s) { ... }

AVANTAGES:
  - Code testable (injection de mocks)
  - Découplage (ne dépend que du comportement nécessaire)
  - Principe ISP (Interface Segregation Principle)

EXCEPTIONS AUTOMATIQUES:
  - Paramètres déjà typés par interface
  - Types du même package
  - Types allowlist (time.Time, etc.)
  - Paramètres utilisés uniquement comme données (sans appel de méthodes)`,
	})
}
