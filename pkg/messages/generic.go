// Package messages provides structured error messages for KTN rules.
// This file contains GENERIC rule messages.
package messages

// registerGenericMessages enregistre les messages GENERIC.
func registerGenericMessages() {
	// Enregistrer KTN-GENERIC-001
	registerGeneric001()
	// Enregistrer KTN-GENERIC-002
	registerGeneric002()
	// Enregistrer KTN-GENERIC-003
	registerGeneric003()
	// Enregistrer KTN-GENERIC-005
	registerGeneric005()
	// Enregistrer KTN-GENERIC-006
	registerGeneric006()
}

// registerGeneric001 enregistre le message KTN-GENERIC-001.
func registerGeneric001() {
	Register(Message{
		Code:  "KTN-GENERIC-001",
		Short: "fonction generique '%s' utilise == ou != sans contrainte comparable",
		Verbose: `PROBLEME: La fonction generique '%s' utilise == ou != sur un type parameter avec contrainte 'any'.

POURQUOI: Les operateurs == et != necessitent que le type soit comparable.
La contrainte 'any' accepte des types non-comparables (slices, maps, fonctions).

EXEMPLE INCORRECT:
  func Contains[T any](s []T, v T) bool {
      for _, x := range s {
          if x == v { return true }  // ERREUR: T peut etre non-comparable
      }
      return false
  }

EXEMPLE CORRECT:
  func Contains[T comparable](s []T, v T) bool {
      for _, x := range s {
          if x == v { return true }  // OK: T est comparable
      }
      return false
  }

ALTERNATIVE: Utiliser une fonction de comparaison:
  func Contains[T any](s []T, v T, eq func(T, T) bool) bool {
      for _, x := range s {
          if eq(x, v) { return true }
      }
      return false
  }`,
	})
}

// registerGeneric002 enregistre le message KTN-GENERIC-002.
func registerGeneric002() {
	Register(Message{
		Code:  "KTN-GENERIC-002",
		Short: "fonction '%s' utilise un generique [%s %s] inutilement",
		Verbose: `PROBLEME: La fonction '%s' utilise un type parameter [%s %s] alors que l'interface pourrait etre utilisee directement.

POURQUOI: Les generiques avec une contrainte interface simple n'apportent aucun benefice si le type n'est pas retourne ou utilise pour la preservation du type.

EXEMPLE INCORRECT:
  func ReadSome[T io.Reader](r T) ([]byte, error) {
      buf := make([]byte, 1024)
      n, _ := r.Read(buf)
      return buf[:n], nil
  }

EXEMPLE CORRECT:
  func ReadSome(r io.Reader) ([]byte, error) {
      buf := make([]byte, 1024)
      n, _ := r.Read(buf)
      return buf[:n], nil
  }

EXCEPTION: Les generiques sont justifies si le type est preserve:
  func PassThrough[T io.Reader](r T) T {
      return r  // Le type T est retourne, justifiant le generique
  }`,
	})
}

// registerGeneric003 enregistre le message KTN-GENERIC-003.
func registerGeneric003() {
	Register(Message{
		Code:  "KTN-GENERIC-003",
		Short: "import obsolete golang.org/x/exp/constraints, utiliser cmp",
		Verbose: `PROBLEME: Le package golang.org/x/exp/constraints est obsolete.

POURQUOI: Depuis Go 1.21, le package standard 'cmp' fournit cmp.Ordered
et cmp.Compare qui remplacent constraints.Ordered de x/exp/constraints.

EXEMPLE INCORRECT:
  import "golang.org/x/exp/constraints"
  func Max[T constraints.Ordered](a, b T) T { ... }

EXEMPLE CORRECT:
  import "cmp"
  func Max[T cmp.Ordered](a, b T) T { ... }

AVANTAGE: Le package 'cmp' est dans la bibliotheque standard, pas de dependance externe.`,
	})
}

// registerGeneric005 enregistre le message KTN-GENERIC-005.
func registerGeneric005() {
	Register(Message{
		Code:  "KTN-GENERIC-005",
		Short: "type parameter '%s' masque un identifiant predeclare",
		Verbose: `PROBLEME: Le type parameter '%s' masque un identifiant predeclare de Go.

POURQUOI: Utiliser un identifiant predeclare comme nom de type parameter:
- Cree de la confusion dans le code
- Peut generer des messages d'erreur cryptiques
- Rend le code difficile a maintenir

IDENTIFIANTS PREDECLARES:
  Types: bool, byte, complex64, complex128, error, float32, float64,
         int, int8, int16, int32, int64, rune, string, uint, uint8,
         uint16, uint32, uint64, uintptr, any, comparable
  Constantes: true, false, iota, nil
  Fonctions: append, cap, clear, close, complex, copy, delete, imag,
             len, make, max, min, new, panic, print, println, real, recover

EXEMPLE INCORRECT:
  func Process[string any](s string) { ... }  // "string" masque le type predeclare
  func Handle[error any](e error) { ... }     // "error" masque le type predeclare

EXEMPLE CORRECT:
  func Process[T any](s T) { ... }            // "T" est conventionnel
  func Handle[E any](e E) { ... }             // "E" est conventionnel

CONVENTIONS:
  - T, U, V pour les types generiques
  - K, V pour les cles/valeurs de maps
  - E pour les elements de collections
  - Noms descriptifs: Element, Item, Key, Value`,
	})
}

// registerGeneric006 enregistre le message KTN-GENERIC-006.
func registerGeneric006() {
	Register(Message{
		Code:  "KTN-GENERIC-006",
		Short: "fonction generique '%s' utilise des operateurs ordered/arithmetiques sans contrainte cmp.Ordered",
		Verbose: `PROBLEME: La fonction generique '%s' utilise des operateurs (<, >, <=, >=, +, -, *, /, %%) sur un type parameter avec contrainte 'any'.

POURQUOI: Les operateurs de comparaison ordered et arithmetiques necessitent que le type supporte ces operations.
La contrainte 'any' accepte des types qui ne supportent pas ces operateurs (structs, slices, maps, fonctions).

EXEMPLE INCORRECT:
  func Min[T any](a, b T) T {
      if a < b { return a }  // ERREUR: < non defini sur any
      return b
  }

  func Sum[T any](values ...T) T {
      var sum T
      for _, v := range values {
          sum = sum + v  // ERREUR: + non defini sur any
      }
      return sum
  }

EXEMPLE CORRECT:
  func Min[T cmp.Ordered](a, b T) T {
      if a < b { return a }  // OK: T est ordonne
      return b
  }

  func Sum[T cmp.Ordered](values ...T) T {
      var sum T
      for _, v := range values {
          sum = sum + v  // OK: T supporte +
      }
      return sum
  }

ALTERNATIVE: Utiliser une fonction de comparaison/operation:
  func MinFunc[T any](a, b T, less func(T, T) bool) T {
      if less(a, b) { return a }
      return b
  }`,
	})
}
