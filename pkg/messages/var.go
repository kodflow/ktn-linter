// Package messages provides structured error messages for KTN rules.
// This file contains VAR rule messages.
package messages

// registerVarMessages enregistre les messages VAR.
func registerVarMessages() {
	// VAR-001: Types explicites (ex-VAR-002)
	Register(Message{
		Code:  "KTN-VAR-001",
		Short: "variable '%s' sans type explicite. Format: var name Type = value",
		Verbose: `PROBLEME: La variable de package '%s' n'a pas de type.

POURQUOI: Le type explicite:
  - Documente l'intention
  - Evite les conversions implicites

FORMAT ATTENDU:
  var nomVariable Type = valeur

EXEMPLE INCORRECT:
  var timeout = 30

EXEMPLE CORRECT:
  var timeout time.Duration = 30 * time.Second`,
	})

	// VAR-002: Ordre declaration (ex-VAR-014)
	Register(Message{
		Code:  "KTN-VAR-002",
		Short: "var avant const. Ordre: const -> var -> type -> func",
		Verbose: `PROBLEME: bloc var declare avant bloc const.

POURQUOI: L'ordre standard facilite la navigation:
  1. const
  2. var
  3. type
  4. func`,
	})

	// VAR-003: CamelCase - no underscores (fusion VAR-001+018)
	Register(Message{
		Code:  "KTN-VAR-003",
		Short: "variable '%s' contient underscore. Utiliser camelCase",
		Verbose: `PROBLEME: La variable '%s' contient un underscore.

POURQUOI: Go utilise camelCase pour toutes les variables.
Les underscores (snake_case ou SCREAMING_SNAKE_CASE) ne sont pas idiomatiques.

EXEMPLE INCORRECT:
  var MAX_SIZE = 1024        // SCREAMING_SNAKE_CASE
  var user_name string       // snake_case
  var Api_Key string         // Mixed_Case

EXEMPLE CORRECT:
  var maxSize = 1024         // camelCase prive
  var MaxSize = 1024         // PascalCase exporte
  var userName string
  var APIKey string`,
	})

	// VAR-004: Longueur min (NEW)
	Register(Message{
		Code:  "KTN-VAR-004",
		Short: "variable '%s' trop courte (min 2 caracteres)",
		Verbose: `PROBLEME: La variable '%s' a un nom trop court.

POURQUOI: Les noms a 1 caractere sont difficiles a comprendre
sauf dans des contextes specifiques (boucles, idiomes).

EXCEPTIONS AUTORISEES:
  - Boucles: i, j, k, n, x, y, z
  - Idiomes: ok

EXEMPLE INCORRECT:
  a := 42
  b := "hello"

EXEMPLE CORRECT:
  count := 42
  message := "hello"
  for i := 0; i < 10; i++ {}  // i autorise en boucle`,
	})

	// VAR-005: Longueur max (NEW)
	Register(Message{
		Code:  "KTN-VAR-005",
		Short: "variable '%s' trop longue (max 30 caracteres)",
		Verbose: `PROBLEME: La variable '%s' depasse 30 caracteres.

POURQUOI: Les noms trop longs nuisent a la lisibilite.

EXEMPLE INCORRECT:
  thisIsAVeryLongVariableNameThatExceedsLimit := 1

EXEMPLE CORRECT:
  maxConnPoolSize := 1`,
	})

	// VAR-006: Shadowing built-in identifiers (ISO avec CONST-006)
	Register(Message{
		Code:  "KTN-VAR-006",
		Short: "variable '%s' masque un identifiant built-in",
		Verbose: `PROBLEME: La variable '%s' masque un identifiant built-in Go.

POURQUOI: Masquer les built-ins cause des comportements inattendus:
  - Confusion pour les lecteurs du code
  - Erreurs subtiles difficiles a deboguer
  - Le built-in devient inaccessible dans ce scope

BUILT-INS PROTEGES:
  Types: bool, byte, int, string, error, any, ...
  Constants: true, false, iota, nil
  Functions: len, cap, append, make, new, panic, ...

EXEMPLE INCORRECT:
  var len int = 100  // Masque la fonction len()

EXEMPLE CORRECT:
  var maxLen int = 100`,
	})

	// VAR-007: := vs var (ex-VAR-003)
	Register(Message{
		Code:  "KTN-VAR-007",
		Short: "utiliser := au lieu de var pour variable locale",
		Verbose: `PROBLEME: 'var' utilise pour une variable locale.

POURQUOI: En Go, := est prefere pour les variables locales:
  - Plus concis
  - Idiomatique
  - Le type est infere

EXEMPLE INCORRECT:
  var x int = 42
  var err error = nil

EXEMPLE CORRECT:
  x := 42
  var err error  // OK si zero value voulue`,
	})

	// VAR-008: Slices prealloc (ex-VAR-004)
	Register(Message{
		Code:  "KTN-VAR-008",
		Short: "slice non preallouee. Utiliser make([]T, 0, %d)",
		Verbose: `PROBLEME: La slice n'est pas preallouee malgre capacite connue.

POURQUOI: Sans preallocation, append() realloue a chaque depassement,
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

	// VAR-009: make+append (ex-VAR-005)
	Register(Message{
		Code:  "KTN-VAR-009",
		Short: "make([]T, %d) avec append cause reallocation. Utiliser cap",
		Verbose: `PROBLEME: make([]T, n) cree n elements, append en ajoute apres.

POURQUOI: make([]T, n) initialise a zero, append ajoute EN PLUS.
  make([]T, 5) puis append(s, x) -> len=6, pas len=5!

EXEMPLE INCORRECT:
  s := make([]int, 10)
  s = append(s, 42)  // len=11!

EXEMPLE CORRECT:
  s := make([]int, 0, 10)
  s = append(s, 42)  // len=1, cap=10`,
	})

	// VAR-010: Buffer.Grow (ex-VAR-006)
	Register(Message{
		Code:  "KTN-VAR-010",
		Short: "Buffer/Builder sans Grow(). Preallouer avec Grow(%d)",
		Verbose: `PROBLEME: bytes.Buffer ou strings.Builder sans Grow().

POURQUOI: Sans Grow(), le buffer realloue a chaque depassement.

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

	// VAR-011: strings.Builder (ex-VAR-007)
	Register(Message{
		Code:  "KTN-VAR-011",
		Short: "%d concatenations string. Utiliser strings.Builder",
		Verbose: `PROBLEME: %d concatenations de string avec +.

POURQUOI: Chaque + cree une nouvelle string (immutable).
strings.Builder evite les allocations.

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

	// VAR-012: Alloc loops (ex-VAR-008)
	Register(Message{
		Code:  "KTN-VAR-012",
		Short: "allocation dans boucle chaude. Sortir de la boucle",
		Verbose: `PROBLEME: Allocation repetee dans une boucle.

POURQUOI: Allouer dans une boucle:
  - Cree de la pression GC
  - Ralentit l'execution
  - Peut etre evite

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

	// VAR-013: Struct size (ex-VAR-009)
	Register(Message{
		Code:  "KTN-VAR-013",
		Short: "struct de %d bytes passee par valeur (seuil: %d). Utiliser pointeur",
		Verbose: `PROBLEME: Struct de %d bytes passee par valeur (seuil: %d bytes).

POURQUOI: Passer par valeur copie toute la struct.
Un pointeur copie seulement 8 bytes (64-bit).

SEUIL: >%d bytes -> utiliser pointeur

EXEMPLE INCORRECT:
  func Process(data LargeStruct) { ... }

EXEMPLE CORRECT:
  func Process(data *LargeStruct) { ... }`,
	})

	// VAR-014: sync.Pool (ex-VAR-010)
	Register(Message{
		Code:  "KTN-VAR-014",
		Short: "buffer repete sans sync.Pool. Utiliser Pool",
		Verbose: `PROBLEME: Buffers alloues/liberes en boucle.

POURQUOI: sync.Pool reutilise les objets et reduit le GC.

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

	// VAR-015: string() (ex-VAR-012)
	Register(Message{
		Code:  "KTN-VAR-015",
		Short: "string() appele %d fois sur meme valeur. Stocker le resultat",
		Verbose: `PROBLEME: string() appele plusieurs fois sur meme []byte.

POURQUOI: Chaque string() alloue une nouvelle string.

EXEMPLE INCORRECT:
  if string(data) == "foo" || string(data) == "bar" { }

EXEMPLE CORRECT:
  s := string(data)
  if s == "foo" || s == "bar" { }`,
	})

	// VAR-016: Groupement (ex-VAR-013)
	Register(Message{
		Code:  "KTN-VAR-016",
		Short: "variables non groupees. Utiliser un seul bloc var()",
		Verbose: `PROBLEME: Plusieurs blocs var separes.

POURQUOI: Grouper ameliore la lisibilite et la coherence.

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

	// VAR-017: Map prealloc (ex-VAR-015)
	Register(Message{
		Code:  "KTN-VAR-017",
		Short: "map sans capacite. Utiliser make(map[K]V, %d)",
		Verbose: `PROBLEME: Map creee sans capacite initiale connue.

POURQUOI: Sans capacite, la map realloue en grandissant.

EXEMPLE INCORRECT:
  m := make(map[string]int)
  for _, item := range items {  // len(items) connu
      m[item.Key] = item.Val
  }

EXEMPLE CORRECT:
  m := make(map[string]int, len(items))`,
	})

	// VAR-018: Array vs slice (ex-VAR-016)
	Register(Message{
		Code:  "KTN-VAR-018",
		Short: "make([]T, %d) pour taille fixe. Utiliser [%d]T",
		Verbose: `PROBLEME: make() pour une taille constante connue a la compilation.

POURQUOI: Un array [N]T est sur la stack (pas d'allocation heap).

EXEMPLE INCORRECT:
  buf := make([]byte, 32)  // Heap allocation

EXEMPLE CORRECT:
  var buf [32]byte  // Stack allocation`,
	})

	// VAR-019: Mutex copies (ex-VAR-017)
	Register(Message{
		Code:  "KTN-VAR-019",
		Short: "copie de mutex '%s'. Utiliser pointeur ou embed",
		Verbose: `PROBLEME: sync.Mutex/RWMutex copie par valeur.

POURQUOI: Copier un mutex copie son etat interne,
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
}
