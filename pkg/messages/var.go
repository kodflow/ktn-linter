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

	// VAR-020: Nil slice preferred
	Register(Message{
		Code:  "KTN-VAR-020",
		Short: "preferer nil slice a '%s'. Utiliser: var s []T",
		Verbose: `PROBLEME: Slice vide declaree avec []T{} ou make([]T, 0).

POURQUOI: Une nil slice est fonctionnellement equivalente a une
slice vide, mais plus efficace (pas d'allocation).
  - nil slice et slice vide ont len=0 et cap=0
  - Les deux supportent append, range, etc.

EXEMPLE INCORRECT:
  items := []string{}      // Allocation inutile
  data := make([]int, 0)   // Allocation inutile

EXEMPLE CORRECT:
  var items []string       // nil slice, pas d'allocation
  var data []int           // nil slice, pas d'allocation
  prealloc := make([]int, 0, 10)  // OK: capacite specifiee`,
	})

	// VAR-021: Receiver consistency
	Register(Message{
		Code:  "KTN-VAR-021",
		Short: "type '%s': receiver incohérent. Attendu: %s",
		Verbose: `PROBLEME: Le type '%s' a des receivers de types differents.

POURQUOI: Toutes les methodes d'un type doivent utiliser
le meme type de receiver (tous pointeur ou tous valeur).

RAISONS:
  - Coherence: comportement predictible
  - Semantique: un type est modifiable (pointeur) ou non (valeur)
  - API claire: pas de confusion pour l'utilisateur

EXEMPLE INCORRECT:
  func (s *Server) Start() {}  // pointeur
  func (s Server) Stop() {}    // valeur - incohérent!

EXEMPLE CORRECT:
  func (s *Server) Start() {}  // pointeur
  func (s *Server) Stop() {}   // pointeur - coherent`,
	})

	// VAR-022: Pointer to interface
	Register(Message{
		Code:  "KTN-VAR-022",
		Short: "pointeur vers interface '%s'. Utiliser l'interface directement",
		Verbose: `PROBLEME: Pointeur vers interface detecte (%s).

POURQUOI: Une interface est deja un fat pointer (type + data).
Un pointeur vers interface est rarement utile et souvent une erreur.
  - *io.Reader, *io.Writer, *interface{}, *any
  - Double indirection inutile
  - API plus complexe sans benefice

EXEMPLE INCORRECT:
  func process(r *io.Reader) { ... }
  var handler *interface{}

EXEMPLE CORRECT:
  func process(r io.Reader) { ... }
  var handler interface{}`,
	})

	// VAR-023: crypto/rand for secrets
	Register(Message{
		Code:  "KTN-VAR-023",
		Short: "math/rand en contexte securite. Utiliser crypto/rand",
		Verbose: `PROBLEME: math/rand utilise dans un contexte securite.

POURQUOI: math/rand utilise un PRNG previsible.
Pour la cryptographie, utiliser crypto/rand.

MOTS-CLES DETECTES:
  key, token, secret, password, salt, nonce, crypt, auth, credential

EXEMPLE INCORRECT:
  func generateToken() int64 {
      return rand.Int63()  // INSECURE!
  }

EXEMPLE CORRECT:
  func generateToken() *big.Int {
      token, _ := rand.Int(rand.Reader, big.NewInt(1000000))
      return token
  }`,
	})

	// VAR-024: any vs interface{}
	Register(Message{
		Code:  "KTN-VAR-024",
		Short: "preferer any a interface{}",
		Verbose: `PROBLEME: interface{} utilise au lieu de any.

POURQUOI: Depuis Go 1.18, any est l'alias de interface{}.
Utiliser any est plus lisible et idiomatique.

EXEMPLE INCORRECT:
  func process(data interface{}) {}
  var x interface{}
  type Container struct {
      value interface{}
  }

EXEMPLE CORRECT:
  func process(data any) {}
  var x any
  type Container struct {
      value any
  }`,
	})

	// VAR-025: clear() built-in (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-025",
		Short: "utiliser clear() au lieu de boucle range pour vider %s",
		Verbose: `PROBLEME: Boucle range utilisee pour vider une %s.

POURQUOI: Depuis Go 1.21, clear() est la fonction built-in pour:
  - Vider une map: clear(m) au lieu de for k := range m { delete(m, k) }
  - Remettre a zero une slice: clear(s) au lieu de for i := range s { s[i] = 0 }

AVANTAGES:
  - Plus lisible et idiomatique
  - Potentiellement plus performant (implementation optimisee)
  - Intention explicite

EXEMPLE INCORRECT:
  // Pour vider une map
  for k := range m {
      delete(m, k)
  }
  // Pour remettre a zero une slice
  for i := range s {
      s[i] = 0
  }

EXEMPLE CORRECT:
  clear(m)  // Vide la map
  clear(s)  // Remet tous les elements a leur valeur zero`,
	})

	// VAR-026: min()/max() built-in (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-026",
		Short: "utiliser %s() built-in au lieu de math.%s()",
		Verbose: `PROBLEME: Utilisation de math.%[2]s() au lieu de %[1]s() built-in.

POURQUOI: Depuis Go 1.21, min() et max() sont des fonctions built-in.
Elles sont:
  - Plus generiques (fonctionnent avec int, float64, string, etc.)
  - Supportent plusieurs arguments: min(a, b, c, d)
  - Pas besoin d'importer math

EXEMPLE INCORRECT:
  import "math"
  x := math.Min(a, b)    // float64 seulement
  y := math.Max(c, d)    // float64 seulement

EXEMPLE CORRECT:
  x := min(a, b)         // Fonctionne avec int, float64, etc.
  y := max(c, d, e)      // Supporte plusieurs arguments
  z := min(1, 2, 3, 4)   // Retourne 1`,
	})

	// VAR-027: range over integer (Go 1.22+)
	Register(Message{
		Code:  "KTN-VAR-027",
		Short: "utiliser 'for i := range n' au lieu de 'for i := 0; i < n; i++'",
		Verbose: `PROBLEME: Boucle for classique convertible en range over int.

POURQUOI: Depuis Go 1.22, on peut utiliser 'for i := range n'
au lieu de 'for i := 0; i < n; i++'. Cette syntaxe est:
  - Plus lisible et concise
  - Moins sujette aux erreurs (pas de i++, pas de condition)
  - Idiomatique en Go moderne

PATTERN DETECTE:
  for i := 0; i < n; i++ {
      // utilisation de i
  }

EXEMPLE INCORRECT:
  for i := 0; i < 10; i++ {
      process(i)
  }
  for i := 0; i < len(items); i++ {
      fmt.Println(i)
  }

EXEMPLE CORRECT:
  for i := range 10 {
      process(i)
  }
  for i := range len(items) {
      fmt.Println(i)
  }

NOTE: Ne s'applique pas si:
  - L'init n'est pas 0 (for i := 5; ...)
  - Le step n'est pas i++ (for i := 0; i < n; i += 2)
  - La condition n'est pas < (for i := 0; i <= n; i++)`,
	})

	// VAR-028: loop var copy obsolete (Go 1.22+)
	Register(Message{
		Code:  "KTN-VAR-028",
		Short: "pattern '%s := %s' obsolete depuis Go 1.22",
		Verbose: `PROBLEME: Copie de variable de boucle '%[1]s := %[1]s' obsolete.

POURQUOI: Depuis Go 1.22, les variables de boucle sont automatiquement
copiees a chaque iteration. Le pattern 'v := v' etait necessaire avant
pour eviter les captures de closure, mais il est maintenant inutile.

PATTERN DETECTE:
  for _, v := range items {
      v := v  // Obsolete en Go 1.22+
      go func() { use(v) }()
  }

EXEMPLE INCORRECT:
  for i, v := range items {
      i := i    // Obsolete
      v := v    // Obsolete
      go func() {
          process(i, v)
      }()
  }

EXEMPLE CORRECT (Go 1.22+):
  for i, v := range items {
      go func() {
          process(i, v)  // Safe: v est copie automatiquement
      }()
  }

NOTE: Ce pattern etait necessaire avant Go 1.22 pour eviter
que toutes les goroutines ne capturent la meme valeur.`,
	})

	// VAR-029: slices.Grow (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-029",
		Short: "utiliser slices.Grow() au lieu du pattern manuel de grow",
		Verbose: `PROBLEME: Pattern manuel de grow detecte (if cap-len < n, make+copy).

POURQUOI: Depuis Go 1.21, slices.Grow() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

PATTERN DETECTE:
  if cap(s)-len(s) < n {
      newSlice := make([]T, len(s), len(s)+n)
      copy(newSlice, s)
      s = newSlice
  }

EXEMPLE INCORRECT:
  func grow(s []int, n int) []int {
      if cap(s)-len(s) < n {
          newSlice := make([]int, len(s), len(s)+n)
          copy(newSlice, s)
          s = newSlice
      }
      return s
  }

EXEMPLE CORRECT:
  import "slices"
  func grow(s []int, n int) []int {
      return slices.Grow(s, n)
  }`,
	})

	// VAR-030: slices.Clone (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-030",
		Short: "utiliser slices.Clone() au lieu du pattern '%s'",
		Verbose: `PROBLEME: Clonage manuel d'une slice avec %s.

POURQUOI: Depuis Go 1.21, slices.Clone() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

PATTERNS DETECTES:
  1. make([]T, len(s)) + copy(clone, s)
  2. append([]T(nil), s...)

EXEMPLE INCORRECT:
  // Pattern 1: make + copy
  clone := make([]int, len(original))
  copy(clone, original)

  // Pattern 2: append nil
  clone := append([]int(nil), original...)

EXEMPLE CORRECT:
  import "slices"
  clone := slices.Clone(original)

NOTE: Ne s'applique pas aux copies partielles:
  partial := make([]int, len(s)/2)
  copy(partial, s[:len(s)/2])  // OK: copie partielle`,
	})

	// VAR-031: maps.Clone (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-031",
		Short: "utiliser maps.Clone() au lieu du pattern make+range",
		Verbose: `PROBLEME: Clonage manuel d'une map avec make+range.

POURQUOI: Depuis Go 1.21, maps.Clone() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

EXEMPLE INCORRECT:
  clone := make(map[K]V, len(m))
  for k, v := range m {
      clone[k] = v
  }

EXEMPLE CORRECT:
  import "maps"
  clone := maps.Clone(m)

NOTE: Ne s'applique pas aux transformations:
  for k, v := range m {
      result[k] = v * 2  // Transformation, pas un clone simple
  }`,
	})

	// VAR-033: cmp.Or (Go 1.22+)
	Register(Message{
		Code:  "KTN-VAR-033",
		Short: "utiliser cmp.Or() au lieu du pattern if x != zeroValue",
		Verbose: `PROBLEME: Pattern manuel de valeur par defaut avec if x != zeroValue.

POURQUOI: Depuis Go 1.22, cmp.Or() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

PATTERN DETECTE:
  if x != 0 { return x } return default     // int
  if x != "" { return x } return default    // string
  if x != nil { return x } return default   // pointer/slice/map

EXEMPLE INCORRECT:
  func getPort(port int) int {
      if port != 0 {
          return port
      }
      return 8080
  }

  func getHost(host string) string {
      if host != "" {
          return host
      }
      return "localhost"
  }

EXEMPLE CORRECT:
  import "cmp"
  func getPort(port int) int {
      return cmp.Or(port, 8080)
  }

  func getHost(host string) string {
      return cmp.Or(host, "localhost")
  }

NOTE: Ne s'applique pas aux conditions complexes:
  if x > 10 { return x }  // OK: pas une comparaison avec zero value`,
	})

	// VAR-034: WaitGroup.Go (Go 1.25+)
	Register(Message{
		Code:  "KTN-VAR-034",
		Short: "utiliser wg.Go() au lieu de wg.Add(1)+go func()+defer wg.Done()",
		Verbose: `PROBLEME: Pattern wg.Add(1) + go func() + defer wg.Done() detecte.

POURQUOI: Depuis Go 1.25, sync.WaitGroup a une methode Go() qui remplace
le pattern classique:
  wg.Add(1)
  go func() {
      defer wg.Done()
      // work
  }()

AVANTAGES de wg.Go():
  - Plus concis et lisible
  - Moins de code boilerplate
  - Impossible d'oublier wg.Add(1) ou defer wg.Done()
  - Pas de risque de mismatch entre Add et Done

EXEMPLE INCORRECT:
  var wg sync.WaitGroup
  for _, item := range items {
      wg.Add(1)
      go func(it int) {
          defer wg.Done()
          process(it)
      }(item)
  }
  wg.Wait()

EXEMPLE CORRECT:
  var wg sync.WaitGroup
  for _, item := range items {
      wg.Go(func() {
          process(item)
      })
  }
  wg.Wait()

NOTE: Ce pattern necessite Go 1.25 ou superieur.`,
	})

	// VAR-035: slices.Contains (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-035",
		Short: "utiliser slices.Contains() au lieu du pattern for-range manuel",
		Verbose: `PROBLEME: Recherche manuelle avec boucle for-range et if == return true/false.

POURQUOI: Depuis Go 1.21, slices.Contains() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

PATTERN DETECTE:
  for _, v := range slice {
      if v == target {
          return true
      }
  }
  return false

EXEMPLE INCORRECT:
  func contains(items []string, target string) bool {
      for _, v := range items {
          if v == target {
              return true
          }
      }
      return false
  }

EXEMPLE CORRECT:
  import "slices"
  func contains(items []string, target string) bool {
      return slices.Contains(items, target)
  }

NOTE: Ne s'applique pas aux comparaisons complexes:
  for _, v := range items {
      if v > threshold {  // OK: pas une simple egalite
          return true
      }
  }`,
	})

	// VAR-036: slices.Index (Go 1.21+)
	Register(Message{
		Code:  "KTN-VAR-036",
		Short: "utiliser slices.Index() au lieu du pattern de recherche manuel",
		Verbose: `PROBLEME: Recherche manuelle d'index avec boucle for-range.

POURQUOI: Depuis Go 1.21, slices.Index() est disponible.
Elle est:
  - Plus lisible et concise
  - Optimisee par le compilateur
  - Moins sujette aux erreurs

PATTERN DETECTE:
  for i, v := range slice {
      if v == target {
          return i
      }
  }
  return -1

EXEMPLE INCORRECT:
  func indexOf(items []int, target int) int {
      for i, v := range items {
          if v == target {
              return i
          }
      }
      return -1
  }

EXEMPLE CORRECT:
  import "slices"
  func indexOf(items []int, target int) int {
      return slices.Index(items, target)
  }

NOTE: Ne s'applique pas aux comparaisons complexes:
  for i, v := range items {
      if v > threshold {  // OK: pas une simple egalite
          return i
      }
  }`,
	})

	// VAR-037: maps.Keys/Values (Go 1.23+)
	Register(Message{
		Code:  "KTN-VAR-037",
		Short: "utiliser slices.Collect(maps.%s()) au lieu de boucle range manuelle",
		Verbose: `PROBLEME: Collection manuelle des %[1]s de map avec boucle range.

POURQUOI: Depuis Go 1.23, maps.Keys() et maps.Values() retournent des
iterateurs. Combines avec slices.Collect(), ils offrent une maniere
idiomatique de collecter les cles ou valeurs d'une map.

AVANTAGES:
  - Plus lisible et concis
  - Potentiellement plus performant (implementation optimisee)
  - Intention explicite

PATTERN DETECTE:
  // Pour collecter les cles
  var keys []K
  for k := range m {
      keys = append(keys, k)
  }

  // Pour collecter les valeurs
  var values []V
  for _, v := range m {
      values = append(values, v)
  }

EXEMPLE INCORRECT:
  func getKeys(m map[string]int) []string {
      var keys []string
      for k := range m {
          keys = append(keys, k)
      }
      return keys
  }

EXEMPLE CORRECT:
  import (
      "maps"
      "slices"
  )

  func getKeys(m map[string]int) []string {
      return slices.Collect(maps.Keys(m))
  }

  func getValues(m map[string]int) []int {
      return slices.Collect(maps.Values(m))
  }

NOTE: Ne s'applique pas aux transformations ou filtres:
  for k := range m {
      keys = append(keys, k+"_suffix")  // OK: transformation
  }`,
	})
}
