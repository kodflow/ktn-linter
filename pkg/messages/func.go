// Package messages provides structured error messages for KTN rules.
// This file contains FUNC rule messages.
package messages

// registerFuncMessages enregistre les messages FUNC.
func registerFuncMessages() {
	Register(Message{
		Code:  "KTN-FUNC-001",
		Short: "error doit être le dernier retour, pas en position %d",
		Verbose: `PROBLÈME: Le type 'error' est en position %d au lieu de dernier.

POURQUOI: Convention Go universelle - error toujours en dernier:
  - Cohérence avec la stdlib
  - Pattern if err != nil { } naturel
  - Facilite la lecture

EXEMPLE INCORRECT:
  func Process() (error, *Result) { ... }

EXEMPLE CORRECT:
  func Process() (*Result, error) { ... }`,
	})

	Register(Message{
		Code:  "KTN-FUNC-002",
		Short: "context.Context doit être le 1er paramètre, pas position %d",
		Verbose: `PROBLÈME: context.Context est en position %d au lieu de 1er.

POURQUOI: Convention Go standard (context package doc):
  "Context should be the first parameter, named ctx"

EXEMPLE INCORRECT:
  func GetUser(id int, ctx context.Context) (*User, error)

EXEMPLE CORRECT:
  func GetUser(ctx context.Context, id int) (*User, error)

NOTE: Pour méthodes, ctx vient après le receiver.`,
	})

	Register(Message{
		Code:  "KTN-FUNC-003",
		Short: "else inutile après return/break/continue. Early return préféré",
		Verbose: `PROBLÈME: Un bloc else suit un return/break/continue.

POURQUOI: L'early return réduit l'indentation et améliore la lisibilité.

EXEMPLE INCORRECT:
  if err != nil {
      return err
  } else {
      processData()
      return nil
  }

EXEMPLE CORRECT:
  if err != nil {
      return err
  }
  processData()
  return nil`,
	})

	Register(Message{
		Code:  "KTN-FUNC-004",
		Short: "fonction privée '%s' jamais appelée (code mort)",
		Verbose: `PROBLÈME: La fonction privée '%s' n'est appelée nulle part.

POURQUOI: Le code mort:
  - Alourdit la maintenance
  - Induit en erreur
  - Peut contenir des bugs cachés

ACTIONS:
  1. Utilisée dans tests → vérifier si vraiment utile
  2. "Pour plus tard" → Supprimer, git la garde
  3. Oubliée après refactoring → Supprimer`,
	})

	Register(Message{
		Code:  "KTN-FUNC-005",
		Short: "fonction '%s' trop complexe (%d statements > %d). Extraire en sous-fonctions",
		Verbose: `PROBLÈME: La fonction '%s' contient %d statements (max %d).

POURQUOI: Les fonctions courtes:
  - Sont plus faciles à tester
  - Ont un nom qui documente leur rôle
  - Sont réutilisables

CALCUL: Compte les instructions logiques (1 appel multi-ligne = 1 statement).
Les blocs if/for/switch ajoutent leurs statements internes.

SOLUTION: Extraire en sous-fonctions nommées.

EXEMPLE:
  // Avant: trop de statements
  func ProcessOrder(order Order) error { ... }

  // Après: fonctions courtes
  func ProcessOrder(order Order) error {
      if err := validateOrder(order); err != nil {
          return err
      }
      return saveOrder(order)
  }`,
	})

	Register(Message{
		Code:  "KTN-FUNC-006",
		Short: "fonction '%s' a %d paramètres (max 5). Grouper dans struct",
		Verbose: `PROBLÈME: La fonction '%s' a %d paramètres (max 5).

POURQUOI: Trop de paramètres:
  - Rendent l'appel difficile à lire
  - Augmentent le risque d'inversion
  - Signalent un problème de conception

SOLUTION: Grouper dans une struct Options/Config.

EXEMPLE INCORRECT:
  func CreateUser(name, email, phone, address, city string) error

EXEMPLE CORRECT:
  type CreateUserRequest struct {
      Name, Email, Phone, Address, City string
  }
  func CreateUser(req CreateUserRequest) error`,
	})

	Register(Message{
		Code:  "KTN-FUNC-007",
		Short: "getter '%s' a des side effects. Un getter doit être pur",
		Verbose: `PROBLÈME: Le getter '%s' (Get*/Is*/Has*) modifie l'état.

POURQUOI: Un getter doit être "pur":
  - Retourner sans modifier l'état
  - Appelable plusieurs fois sans conséquence
  - Pas d'effets observables

EXEMPLE INCORRECT:
  func (c *Counter) GetCount() int {
      c.accessCount++  // Side effect!
      return c.count
  }

EXEMPLE CORRECT:
  func (c *Counter) GetCount() int {
      return c.count
  }`,
	})

	Register(Message{
		Code:  "KTN-FUNC-008",
		Short: "paramètre '%s' non utilisé. Préfixer _ ou supprimer",
		Verbose: `PROBLÈME: Le paramètre '%s' n'est pas utilisé.

SOLUTIONS:
  1. Supprimer si pas requis
  2. Préfixer _ si imposé par interface: _param
  3. Si "_ = param" → c'est du contournement, nettoyer

EXEMPLE INCORRECT:
  func Process(ctx context.Context, data []byte) error {
      // ctx jamais utilisé
      return nil
  }

EXEMPLE CORRECT (si interface l'impose):
  func Process(_ctx context.Context, data []byte) error`,
	})

	Register(Message{
		Code:  "KTN-FUNC-009",
		Short: "nombre magique %v. Extraire en constante nommée",
		Verbose: `PROBLÈME: Le nombre %v est utilisé sans explication.

POURQUOI: Les magic numbers:
  - N'expliquent pas leur signification
  - Sont difficiles à modifier
  - Rendent le code incompréhensible

EXCEPTIONS: 0, 1, -1, 2 (souvent évidents)

EXEMPLE INCORRECT:
  if age < 18 { ... }
  time.Sleep(30 * time.Second)

EXEMPLE CORRECT:
  const MajorityAge = 18
  const DefaultTimeout = 30 * time.Second
  if age < MajorityAge { ... }`,
	})

	Register(Message{
		Code:  "KTN-FUNC-010",
		Short: "naked return interdit (fonction > 5 lignes)",
		Verbose: `PROBLÈME: Return sans valeur dans fonction longue.

POURQUOI: Les naked returns:
  - Cachent ce qui est retourné
  - Rendent le debugging difficile
  - Sont source de bugs subtils

EXCEPTION: Autorisé pour fonctions < 5 lignes.

EXEMPLE INCORRECT:
  func GetUser(id int) (user *User, err error) {
      user, err = db.Find(id)
      return  // Que retourne-t-on?
  }

EXEMPLE CORRECT:
  func GetUser(id int) (*User, error) {
      user, err := db.Find(id)
      return user, err
  }`,
	})

	Register(Message{
		Code:  "KTN-FUNC-011",
		Short: "complexité cyclomatique %d (max 15). Simplifier",
		Verbose: `PROBLÈME: Complexité cyclomatique de %d (max 15).

POURQUOI: Complexité élevée signifie:
  - Trop de chemins d'exécution
  - Difficile à tester
  - Difficile à maintenir

CALCUL: +1 pour if, else, case, for, &&, ||

SOLUTIONS:
  - Extraire des sous-fonctions
  - Utiliser early returns
  - Remplacer switch par map`,
	})

	Register(Message{
		Code:  "KTN-FUNC-012",
		Short: "%d retours sans noms. Named returns requis pour >3 valeurs",
		Verbose: `PROBLÈME: La fonction retourne %d valeurs sans noms.

POURQUOI: Pour >3 retours, les noms:
  - Documentent chaque position
  - Facilitent la lecture
  - Servent de documentation

EXEMPLE INCORRECT:
  func Parse() (string, int, bool, []string, error)

EXEMPLE CORRECT:
  func Parse() (path string, port int, debug bool, hosts []string, err error)`,
	})
}
