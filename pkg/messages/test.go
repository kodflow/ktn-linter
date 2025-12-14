// Package messages provides structured error messages for KTN rules.
// This file contains TEST rule messages.
package messages

// registerTestMessages enregistre les messages TEST.
func registerTestMessages() {
	Register(Message{
		Code:  "KTN-TEST-001",
		Short: "fichier '%s' doit être renommé en '%s', '%s', '%s' ou '%s'",
		Verbose: `PROBLÈME: Le fichier test '%s' n'a pas le bon suffixe.

POURQUOI: Les fichiers de test doivent suivre les conventions:
  - _internal_test.go pour white-box testing (même package)
  - _external_test.go pour black-box testing (package _test)
  - _bench_test.go pour les benchmarks
  - _integration_test.go pour les tests d'intégration

FICHIERS SUGGÉRÉS:
  - White-box: %s
  - Black-box: %s
  - Benchmark: %s
  - Integration: %s`,
	})

	Register(Message{
		Code:  "KTN-TEST-002",
		Short: "package test incorrect '%s'. Utiliser '%s_test'",
		Verbose: `PROBLÈME: Le fichier test utilise le package '%s'.

POURQUOI: Les tests doivent utiliser le package xxx_test
pour du black-box testing (tester l'API publique).

EXEMPLE INCORRECT:
  package mypackage  // Même package = white-box

EXEMPLE CORRECT:
  package mypackage_test  // Black-box testing`,
	})

	Register(Message{
		Code:  "KTN-TEST-003",
		Short: "fichier test '%s' sans fichier source '%s'",
		Verbose: `PROBLÈME: Le fichier test '%s' n'a pas de fichier source '%s'.

POURQUOI: Chaque fichier de test doit correspondre à un fichier source.

SOLUTIONS:
  1. Créer le fichier source correspondant
  2. Renommer le fichier de test
  3. Déplacer les tests vers le bon fichier`,
	})

	Register(Message{
		Code:  "KTN-TEST-004",
		Short: "fonction '%s' sans test correspondant",
		Verbose: `PROBLÈME: La fonction '%s' n'a pas de test.

POURQUOI: Chaque fonction doit avoir un test pour:
  - Documenter le comportement attendu
  - Prévenir les régressions
  - Faciliter le refactoring

FORMAT ATTENDU:
  - Publique: TestNomFonction dans *_external_test.go
  - Privée: Test_nomFonction dans *_internal_test.go`,
	})

	Register(Message{
		Code:  "KTN-TEST-005",
		Short: "test '%s' sans table-driven pattern",
		Verbose: `PROBLÈME: Le test '%s' n'utilise pas table-driven.

POURQUOI: Le pattern table-driven:
  - Facilite l'ajout de cas
  - Rend les tests lisibles
  - Évite la duplication

EXEMPLE:
  func TestAdd(t *testing.T) {
      tests := []struct {
          name     string
          a, b     int
          expected int
      }{
          {"positifs", 1, 2, 3},
          {"négatifs", -1, -1, -2},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              got := Add(tt.a, tt.b)
              if got != tt.expected {
                  t.Errorf("got %d, want %d", got, tt.expected)
              }
          })
      }
  }`,
	})

	Register(Message{
		Code:  "KTN-TEST-006",
		Short: "fichier test '%s' sans fichier source correspondant",
		Verbose: `PROBLÈME: Le fichier test '%s' n'a pas de fichier source 1:1.

POURQUOI: Chaque fichier de test doit correspondre exactement
à un fichier source (pattern 1:1).

CONVENTION:
  - foo.go → foo_internal_test.go et/ou foo_external_test.go
  - Pas de tests orphelins`,
	})

	Register(Message{
		Code:  "KTN-TEST-007",
		Short: "t.Skip() interdit dans '%s'. Les tests doivent passer",
		Verbose: `PROBLÈME: Le test '%s' utilise t.Skip().

POURQUOI: t.Skip() cache des tests cassés:
  - Les tests doivent toujours passer
  - Un test skipé est souvent oublié
  - Si le test n'est plus valide, le supprimer

ALTERNATIVES:
  - Fixer le test
  - Utiliser build tags pour environnements spécifiques
  - Supprimer si obsolète`,
	})

	Register(Message{
		Code:  "KTN-TEST-008",
		Short: "fichier '%s' sans tests. Créer %s",
		Verbose: `PROBLÈME: Le fichier '%s' n'a pas de fichier test approprié.

POURQUOI: Chaque fichier .go doit avoir ses tests.

FICHIER(S) À CRÉER: %s

CONVENTION:
  - Fonctions publiques → xxx_external_test.go (package xxx_test, black-box)
  - Fonctions privées → xxx_internal_test.go (package xxx, white-box)`,
	})

	Register(Message{
		Code:  "KTN-TEST-009",
		Short: "test '%s' de fonction publique doit être dans _external_test.go",
		Verbose: `PROBLÈME: Le test '%s' teste une fonction publique mais n'est pas dans _external_test.go.

POURQUOI: Les fonctions publiques (exportées) doivent être testées
en black-box testing uniquement:
  - Teste l'API publique sans connaître l'implémentation
  - Package xxx_test (séparé)
  - Fichier _external_test.go`,
	})

	Register(Message{
		Code:  "KTN-TEST-010",
		Short: "test '%s' de fonction privée doit être dans _internal_test.go",
		Verbose: `PROBLÈME: Le test '%s' teste une fonction privée mais n'est pas dans _internal_test.go.

POURQUOI: Les fonctions privées (non-exportées) doivent être testées
en white-box testing uniquement:
  - Accès aux fonctions privées
  - Package xxx (même package)
  - Fichier _internal_test.go`,
	})

	Register(Message{
		Code:  "KTN-TEST-011",
		Short: "package '%s' incorrect pour fichier %s",
		Verbose: `PROBLÈME: Le package '%s' ne correspond pas au type de fichier test.

CONVENTION:
  - _internal_test.go → package xxx (même package, white-box)
  - _external_test.go → package xxx_test (séparé, black-box)`,
	})

	Register(Message{
		Code:  "KTN-TEST-012",
		Short: "test '%s' sans assertions. Un test doit tester",
		Verbose: `PROBLÈME: Le test '%s' ne contient pas d'assertions.

POURQUOI: Un test sans assertion:
  - Ne vérifie rien
  - Passe toujours (faux positif)
  - Est inutile

UN TEST DOIT CONTENIR:
  - t.Error/t.Errorf
  - t.Fatal/t.Fatalf
  - Comparaisons avec expected`,
	})

	Register(Message{
		Code:  "KTN-TEST-013",
		Short: "test '%s' ne couvre pas les cas d'erreur",
		Verbose: `PROBLÈME: Le test '%s' ne teste pas les erreurs.

POURQUOI: Les cas d'erreur:
  - Sont souvent les plus critiques
  - Révèlent des bugs de gestion
  - Documentent le comportement d'erreur

À TESTER:
  - Paramètres invalides
  - Ressources indisponibles
  - Limites dépassées
  - nil/zero values`,
	})
}
