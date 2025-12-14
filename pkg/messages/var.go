// Var messages for KTN-VAR rules.
package messages

// registerVarMessages enregistre les messages VAR.
func registerVarMessages() {
	Register(Message{
		Code:  "KTN-VAR-001",
		Short: "variable '%s' utilise SCREAMING_SNAKE. Utiliser camelCase",
		Verbose: `PROBLÈME: La variable '%s' utilise SCREAMING_SNAKE_CASE.

POURQUOI: Go utilise camelCase pour tout, y compris variables.

EXEMPLE INCORRECT:
  var MAX_SIZE = 1024
  var DEFAULT_TIMEOUT = 30

EXEMPLE CORRECT:
  var maxSize = 1024       // Privée
  var DefaultTimeout = 30  // Exportée`,
	})

	Register(Message{
		Code:  "KTN-VAR-002",
		Short: "variable '%s' sans type explicite. Format: var name Type = value",
		Verbose: `PROBLÈME: La variable de package '%s' n'a pas de type.

POURQUOI: Le type explicite:
  - Documente l'intention
  - Évite les conversions implicites

FORMAT ATTENDU:
  var nomVariable Type = valeur

EXEMPLE INCORRECT:
  var timeout = 30

EXEMPLE CORRECT:
  var timeout time.Duration = 30 * time.Second`,
	})

	Register(Message{
		Code:  "KTN-VAR-003",
		Short: "utiliser := au lieu de var pour variable locale",
		Verbose: `PROBLÈME: 'var' utilisé pour une variable locale.

POURQUOI: En Go, := est préféré pour les variables locales:
  - Plus concis
  - Idiomatique
  - Le type est inféré

EXEMPLE INCORRECT:
  var x int = 42
  var err error = nil

EXEMPLE CORRECT:
  x := 42
  var err error  // OK si zero value voulue`,
	})

	Register(Message{
		Code:  "KTN-VAR-004",
		Short: "slice non préallouée. Utiliser make([]T, 0, %d)",
		Verbose: `PROBLÈME: La slice n'est pas préallouée malgré capacité connue.

POURQUOI: Sans préallocation, append() réalloue à chaque dépassement,
causant des copies et de la pression GC.

EXEMPLE INCORRECT:
  var items []Item
  for _, x := range data {
      items = append(items, x)
  }

EXEMPLE CORRECT:
  items := make([]Item, 0, len(data))
  for _, x := range data {
      items = append(items, x)
  }`,
	})

	Register(Message{
		Code:  "KTN-VAR-005",
		Short: "make([]T, %d) avec append cause réallocation. Utiliser cap",
		Verbose: `PROBLÈME: make([]T, n) crée n éléments, append en ajoute après.

POURQUOI: make([]T, n) initialise à zéro, append ajoute EN PLUS.
  make([]T, 5) puis append(s, x) → len=6, pas len=5!

EXEMPLE INCORRECT:
  s := make([]int, 10)
  s = append(s, 42)  // len=11!

EXEMPLE CORRECT:
  s := make([]int, 0, 10)
  s = append(s, 42)  // len=1, cap=10`,
	})

	Register(Message{
		Code:  "KTN-VAR-006",
		Short: "Buffer/Builder sans Grow(). Préallouer avec Grow(%d)",
		Verbose: `PROBLÈME: bytes.Buffer ou strings.Builder sans Grow().

POURQUOI: Sans Grow(), le buffer réalloue à chaque dépassement.

EXEMPLE INCORRECT:
  var buf bytes.Buffer
  for _, s := range items {
      buf.WriteString(s)
  }

EXEMPLE CORRECT:
  var buf bytes.Buffer
  buf.Grow(estimatedSize)
  for _, s := range items {
      buf.WriteString(s)
  }`,
	})

	Register(Message{
		Code:  "KTN-VAR-007",
		Short: "%d concaténations string. Utiliser strings.Builder",
		Verbose: `PROBLÈME: %d concaténations de string avec +.

POURQUOI: Chaque + crée une nouvelle string (immutable).
strings.Builder évite les allocations.

EXEMPLE INCORRECT:
  s := ""
  for _, x := range items {
      s += x + ","
  }

EXEMPLE CORRECT:
  var b strings.Builder
  for _, x := range items {
      b.WriteString(x)
      b.WriteString(",")
  }
  s := b.String()`,
	})

	Register(Message{
		Code:  "KTN-VAR-008",
		Short: "allocation dans boucle chaude. Sortir de la boucle",
		Verbose: `PROBLÈME: Allocation répétée dans une boucle.

POURQUOI: Allouer dans une boucle:
  - Crée de la pression GC
  - Ralentit l'exécution
  - Peut être évité

EXEMPLE INCORRECT:
  for i := 0; i < 1000; i++ {
      buf := make([]byte, 1024)
      process(buf)
  }

EXEMPLE CORRECT:
  buf := make([]byte, 1024)
  for i := 0; i < 1000; i++ {
      process(buf)
      clear(buf)
  }`,
	})

	Register(Message{
		Code:  "KTN-VAR-009",
		Short: "struct de %d bytes passée par valeur. Utiliser pointeur",
		Verbose: `PROBLÈME: Struct de %d bytes passée par valeur (>64 bytes).

POURQUOI: Passer par valeur copie toute la struct.
Un pointeur copie seulement 8 bytes (64-bit).

SEUIL: >64 bytes → utiliser pointeur

EXEMPLE INCORRECT:
  func Process(data LargeStruct) { ... }

EXEMPLE CORRECT:
  func Process(data *LargeStruct) { ... }`,
	})

	Register(Message{
		Code:  "KTN-VAR-010",
		Short: "buffer répété sans sync.Pool. Utiliser Pool",
		Verbose: `PROBLÈME: Buffers alloués/libérés en boucle.

POURQUOI: sync.Pool réutilise les objets et réduit le GC.

EXEMPLE INCORRECT:
  for req := range requests {
      buf := make([]byte, 4096)
      process(buf)
  }

EXEMPLE CORRECT:
  var bufPool = sync.Pool{
      New: func() any { return make([]byte, 4096) },
  }
  for req := range requests {
      buf := bufPool.Get().([]byte)
      process(buf)
      bufPool.Put(buf)
  }`,
	})

	Register(Message{
		Code:  "KTN-VAR-011",
		Short: "shadowing de '%s' avec :=. Utiliser = pour réassigner",
		Verbose: `PROBLÈME: La variable '%s' est shadowée avec :=.

POURQUOI: := crée une NOUVELLE variable qui cache l'ancienne.
C'est souvent un bug subtil.

EXEMPLE INCORRECT:
  x := 10
  if cond {
      x := 20  // Nouvelle variable, shadow!
  }
  // x est toujours 10

EXEMPLE CORRECT:
  x := 10
  if cond {
      x = 20  // Réassignation
  }
  // x est 20`,
	})

	Register(Message{
		Code:  "KTN-VAR-012",
		Short: "string() appelé %d fois sur même valeur. Stocker le résultat",
		Verbose: `PROBLÈME: string() appelé plusieurs fois sur même []byte.

POURQUOI: Chaque string() alloue une nouvelle string.

EXEMPLE INCORRECT:
  if string(data) == "foo" || string(data) == "bar" { }

EXEMPLE CORRECT:
  s := string(data)
  if s == "foo" || s == "bar" { }`,
	})

	Register(Message{
		Code:  "KTN-VAR-013",
		Short: "variables non groupées. Utiliser un seul bloc var()",
		Verbose: `PROBLÈME: Plusieurs blocs var séparés.

POURQUOI: Grouper améliore la lisibilité et la cohérence.

EXEMPLE INCORRECT:
  var x int = 1
  var y int = 2
  var z int = 3

EXEMPLE CORRECT:
  var (
      x int = 1
      y int = 2
      z int = 3
  )`,
	})

	Register(Message{
		Code:  "KTN-VAR-014",
		Short: "var avant const. Ordre: const → var → type → func",
		Verbose: `PROBLÈME: bloc var déclaré avant bloc const.

POURQUOI: L'ordre standard facilite la navigation:
  1. const
  2. var
  3. type
  4. func`,
	})

	Register(Message{
		Code:  "KTN-VAR-015",
		Short: "map sans capacité. Utiliser make(map[K]V, %d)",
		Verbose: `PROBLÈME: Map créée sans capacité initiale connue.

POURQUOI: Sans capacité, la map réalloue en grandissant.

EXEMPLE INCORRECT:
  m := make(map[string]int)
  for _, item := range items {  // len(items) connu
      m[item.Key] = item.Val
  }

EXEMPLE CORRECT:
  m := make(map[string]int, len(items))`,
	})

	Register(Message{
		Code:  "KTN-VAR-016",
		Short: "make([]T, %d) pour taille fixe. Utiliser [%d]T",
		Verbose: `PROBLÈME: make() pour une taille constante connue à la compilation.

POURQUOI: Un array [N]T est sur la stack (pas d'allocation heap).

EXEMPLE INCORRECT:
  buf := make([]byte, 32)  // Heap allocation

EXEMPLE CORRECT:
  var buf [32]byte  // Stack allocation`,
	})

	Register(Message{
		Code:  "KTN-VAR-017",
		Short: "copie de mutex '%s'. Utiliser pointeur ou embed",
		Verbose: `PROBLÈME: sync.Mutex/RWMutex copié par valeur.

POURQUOI: Copier un mutex copie son état interne,
causant des deadlocks ou data races.

EXEMPLE INCORRECT:
  func process(m sync.Mutex) { ... }  // Copie!

EXEMPLE CORRECT:
  func process(m *sync.Mutex) { ... }
  // Ou embed dans struct:
  type Safe struct {
      mu sync.Mutex
      data int
  }`,
	})

	Register(Message{
		Code:  "KTN-VAR-018",
		Short: "variable '%s' utilise snake_case. Utiliser camelCase",
		Verbose: `PROBLÈME: La variable '%s' utilise snake_case.

POURQUOI: Go utilise camelCase pour toutes les variables.

EXEMPLE INCORRECT:
  var user_name string
  var max_size int

EXEMPLE CORRECT:
  var userName string
  var maxSize int`,
	})
}
