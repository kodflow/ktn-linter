// Package rules_pool_good contient du code conforme aux règles KTN-POOL et KTN-STRUCT.
package rules_pool_good

import "sync"

// Pools de ressources
//
// Ces pools permettent de réutiliser les objets pour réduire les allocations.
var (
	// bufferPool est un pool de buffers de 1024 bytes.
	bufferPool sync.Pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}

	// objectPool est un pool d'objets dataObject.
	objectPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			return &dataObject{}
		},
	}
)

// dataObject est une struct de test.
type dataObject struct {
	// ID est l'identifiant de l'objet.
	ID int
	// Data contient les données de l'objet.
	Data []byte
}

// ✅ Code conforme KTN-POOL-001 : pool.Get() avec defer Put()

// GoodDeferPut utilise defer correctement.
func GoodDeferPut() {
	if buf, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf) // ✅ Correct : defer Put()

		for i := range buf {
			buf[i] = byte(i)
		}
	} else {
		panic("type assertion failed: expected []byte")
	}
}

// GoodDeferPutObject utilise defer avec objet.
func GoodDeferPutObject() {
	if obj, ok := objectPool.Get().(*dataObject); ok {
		defer objectPool.Put(obj) // ✅ Correct : defer Put()

		obj.ID = 42
		obj.Data = make([]byte, 100)
	} else {
		panic("type assertion failed: expected *dataObject")
	}
}

// GoodDeferMultiple utilise defer pour plusieurs pools.
func GoodDeferMultiple() {
	if buf1, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf1) // ✅ Correct

		if buf2, ok2 := bufferPool.Get().([]byte); ok2 {
			defer bufferPool.Put(buf2) // ✅ Correct

			copy(buf1, buf2)
		} else {
			panic("type assertion failed: expected []byte for buf2")
		}
	} else {
		panic("type assertion failed: expected []byte for buf1")
	}
}

// GoodDeferInGoroutine utilise defer dans goroutine avec synchronisation.
func GoodDeferInGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if buf, ok := bufferPool.Get().([]byte); ok {
			defer bufferPool.Put(buf) // ✅ Correct

			buf[0] = 1
		} else {
			panic("type assertion failed: expected []byte")
		}
	}()
	wg.Wait()
}

// GoodDeferWithReturn utilise defer avant return.
//
// Params:
//   - condition: condition de retour anticipé
func GoodDeferWithReturn(condition bool) {
	if buf, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf) // ✅ Correct : protège tous les returns

		if condition {
			// Retourne tôt si condition vraie, defer s'exécutera quand même
			return
		}

		processBuffer(buf)
	} else {
		panic("type assertion failed: expected []byte")
	}
}

// GoodDeferWithPanic utilise defer qui protège contre panic.
//
// Params:
//   - index: index à accéder dans le buffer
func GoodDeferWithPanic(index int) {
	if buf, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf) // ✅ Correct : protège même si panic

		// ✅ Utiliser range loop pour éviter indexation directe
		for i := range buf {
			if i == index {
				_ = buf[i]
				break
			}
		}
	} else {
		panic("type assertion failed: expected []byte")
	}
}

// GoodLoopWithDefer utilise defer dans boucle (attention : coût).
//
// Params:
//   - n: nombre d'itérations
func GoodLoopWithDefer(n int) {
	for i := 0; i < n; i++ {
		processBufferWithDefer(i)
	}
}

// processBufferWithDefer traite un buffer dans une closure avec defer.
//
// Params:
//   - i: index de l'itération
func processBufferWithDefer(i int) {
	if buf, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf) // ✅ Correct dans fonction séparée

		buf[0] = byte(i)
	} else {
		panic("type assertion failed: expected []byte")
	}
}

// GoodDeferInCondition utilise defer dans condition.
//
// Params:
//   - condition: condition pour traitement
func GoodDeferInCondition(condition bool) {
	if condition {
		if buf, ok := bufferPool.Get().([]byte); ok {
			defer bufferPool.Put(buf) // ✅ Correct

			processBuffer(buf)
		} else {
			panic("type assertion failed: expected []byte")
		}
	}
}

// GoodDeferInSwitch utilise defer dans switch.
//
// Params:
//   - choice: choix de traitement
func GoodDeferInSwitch(choice int) {
	switch choice {
	case 1:
		if buf, ok := bufferPool.Get().([]byte); ok {
			defer bufferPool.Put(buf) // ✅ Correct
			_ = buf
		} else {
			panic("type assertion failed: expected []byte")
		}
	case 2:
		if buf, ok := bufferPool.Get().([]byte); ok {
			defer bufferPool.Put(buf) // ✅ Correct
			_ = buf
		} else {
			panic("type assertion failed: expected []byte")
		}
	}
}

// GoodDeferWithError utilise defer avant erreur potentielle.
//
// Returns:
//   - error: erreur de traitement
func GoodDeferWithError() error {
	if buf, ok := bufferPool.Get().([]byte); ok {
		defer bufferPool.Put(buf) // ✅ Correct : protège même si erreur

		_ = buf
		// Retourne l'erreur de traitement
		return processWithError()
	}
	panic("type assertion failed: expected []byte")
}

// processBuffer traite un buffer en remplissant avec des valeurs.
//
// Params:
//   - buf: buffer à traiter
func processBuffer(buf []byte) {
	for i := range buf {
		buf[i] = byte(i % 256)
	}
}

// processWithError simule un traitement avec erreur potentielle.
//
// Returns:
//   - error: toujours nil dans ce cas
func processWithError() error {
	// Retourne nil car pas d'erreur
	return nil
}

// ✅ Code conforme KTN-STRUCT-004 : grandes structs par pointeur

// largeConfig est une struct de 200 bytes.
type largeConfig struct {
	// Host est l'adresse du serveur.
	Host string // 16 bytes
	// Port est le port du serveur.
	Port int // 8 bytes
	// Timeout est le délai d'attente.
	Timeout int // 8 bytes
	// Buffer est le tampon de données.
	Buffer [128]byte // 128 bytes
	// MaxConns est le nombre max de connexions.
	MaxConns int // 8 bytes
	// Retries est le nombre de tentatives.
	Retries int // 8 bytes
}

// hugeData est une struct de 512 bytes.
type hugeData struct {
	// Buffer contient les données brutes.
	Buffer [512]byte
}

// mediumStruct est une struct de 160 bytes.
type mediumStruct struct {
	// Data1 est le premier bloc de données.
	Data1 [64]byte
	// Data2 est le second bloc de données.
	Data2 [64]byte
	// Flags contient les drapeaux.
	Flags int
	// Count est le compteur.
	Count int
	// Size est la taille.
	Size int
}

// smallStruct est une struct de 24 bytes (OK par valeur).
type smallStruct struct {
	// ID est l'identifiant.
	ID int // 8 bytes
	// Count est le compteur.
	Count int // 8 bytes
	// Active indique si actif.
	Active bool // 1 byte
}

// GoodProcesslargeConfig passe largeConfig par pointeur.
//
// Params:
//   - cfg: configuration par pointeur (8 bytes)
func GoodProcesslargeConfig(cfg *largeConfig) {
	_ = cfg.Host
}

// GoodProcesshugeData passe hugeData par pointeur.
//
// Params:
//   - data: données par pointeur (8 bytes)
func GoodProcesshugeData(data *hugeData) {
	_ = data.Buffer[0]
}

// GoodProcessmediumStruct passe mediumStruct par pointeur.
//
// Params:
//   - m: struct moyenne par pointeur (8 bytes)
func GoodProcessmediumStruct(m *mediumStruct) {
	_ = m.Flags
}

// GoodProcesssmallStruct passe smallStruct par valeur (OK car petite).
//
// Params:
//   - s: petite struct par valeur (24 bytes)
func GoodProcesssmallStruct(s smallStruct) {
	_ = s.ID
}

// GoodMultipleLargeParams a plusieurs paramètres par pointeur.
//
// Params:
//   - cfg: configuration par pointeur
//   - data: données par pointeur
func GoodMultipleLargeParams(cfg *largeConfig, data *hugeData) {
	_ = cfg.Host
	_ = data.Buffer
}

// GoodMixedParams mélange types correctement.
//
// Params:
//   - id: identifiant
//   - cfg: configuration par pointeur
//   - name: nom
func GoodMixedParams(id int, cfg *largeConfig, name string) {
	_ = id
	_ = cfg
	_ = name
}

// GoodReturnLargeStruct retourne pointeur vers large struct.
//
// Returns:
//   - *largeConfig: configuration par pointeur
func GoodReturnLargeStruct() *largeConfig {
	// Retourne une nouvelle configuration avec valeurs par défaut
	return &largeConfig{
		Host: "localhost",
		Port: 8080,
	}
}

// GoodVariadicLargeStruct utilise variadic avec pointeurs.
//
// Params:
//   - configs: configurations par pointeurs
func GoodVariadicLargeStruct(configs ...*largeConfig) {
	for _, cfg := range configs {
		_ = cfg
	}
}

// GoodAnonymousParam a un paramètre anonyme par pointeur.
//
// Params:
//   - (unnamed): configuration anonyme par pointeur
func GoodAnonymousParam(*largeConfig) {
	// Paramètre non utilisé
}

// GoodMethodReceiver a un receiver par pointeur.
func (cfg *largeConfig) GoodMethodReceiver() {
	_ = cfg.Host
}

// nestedStruct a des structs imbriquées larges.
type nestedStruct struct {
	// Config est la configuration.
	Config largeConfig // 200 bytes
	// Data contient les données.
	Data hugeData // 512 bytes
}

// GoodProcessnestedStructs passe struct imbriquée par pointeur.
//
// Params:
//   - n: struct imbriquée par pointeur
func GoodProcessnestedStructs(n *nestedStruct) {
	_ = n.Config.Host
}

// GoodClosureWithLargeStruct crée closure avec pointeur.
//
// Returns:
//   - func(*largeConfig): closure avec pointeur
func GoodClosureWithLargeStruct() func(*largeConfig) {
	// Retourne une closure qui accepte un pointeur vers largeConfig
	return func(cfg *largeConfig) {
		_ = cfg.Host
	}
}

// GoodDeferWithLargeStruct utilise defer avec pointeur.
func GoodDeferWithLargeStruct() {
	defer func(cfg *largeConfig) {
		_ = cfg.Host
	}(&largeConfig{})
}

// GoodGoRoutineWithLargeStruct lance goroutine avec pointeur et synchronisation.
//
// Params:
//   - cfg: configuration par pointeur
func GoodGoRoutineWithLargeStruct(cfg *largeConfig) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(c *largeConfig) {
		defer wg.Done()
		_ = c.Host
	}(cfg)
	wg.Wait()
}

// GoodNoPoolUsage ne fait pas d'allocation pool (pas de règle).
func GoodNoPoolUsage() {
	buf := make([]byte, 1024) // ✅ Correct : pas de pool
	processBuffer(buf)
}

// GoodStackAllocation utilise petite struct sur stack.
func GoodStackAllocation() {
	s := smallStruct{ID: 1, Count: 10} // ✅ Correct : petite struct
	GoodProcesssmallStruct(s)
}
