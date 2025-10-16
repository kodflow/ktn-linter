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

	// objectPool est un pool d'objets DataObject.
	objectPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			return &DataObject{}
		},
	}
)

// DataObject est une struct de test.
type DataObject struct {
	// ID est l'identifiant de l'objet.
	ID int
	// Data contient les données de l'objet.
	Data []byte
}

// ✅ Code conforme KTN-POOL-001 : pool.Get() avec defer Put()

// GoodDeferPut utilise defer correctement.
func GoodDeferPut() {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf) // ✅ Correct : defer Put()

	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}
}

// GoodDeferPutObject utilise defer avec objet.
func GoodDeferPutObject() {
	obj := objectPool.Get().(*DataObject)
	defer objectPool.Put(obj) // ✅ Correct : defer Put()

	obj.ID = 42
	obj.Data = make([]byte, 100)
}

// GoodDeferMultiple utilise defer pour plusieurs pools.
func GoodDeferMultiple() {
	buf1 := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf1) // ✅ Correct

	buf2 := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf2) // ✅ Correct

	copy(buf1, buf2)
}

// GoodDeferInGoroutine utilise defer dans goroutine avec synchronisation.
func GoodDeferInGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf) // ✅ Correct

		buf[0] = 1
	}()
	wg.Wait()
}

// GoodDeferWithReturn utilise defer avant return.
//
// Params:
//   - condition: condition de retour anticipé
func GoodDeferWithReturn(condition bool) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf) // ✅ Correct : protège tous les returns

	if condition {
		// Retourne tôt si condition vraie, defer s'exécutera quand même
		return
	}

	processBuffer(buf)
}

// GoodDeferWithPanic utilise defer qui protège contre panic.
//
// Params:
//   - index: index à accéder dans le buffer
func GoodDeferWithPanic(index int) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf) // ✅ Correct : protège même si panic

	_ = buf[index] // Si panic, defer s'exécute quand même
}

// GoodLoopWithDefer utilise defer dans boucle (attention : coût).
//
// Params:
//   - n: nombre d'itérations
func GoodLoopWithDefer(n int) {
	for i := 0; i < n; i++ {
		func() {
			buf := bufferPool.Get().([]byte)
			defer bufferPool.Put(buf) // ✅ Correct dans closure

			buf[0] = byte(i)
		}()
	}
}

// GoodDeferInCondition utilise defer dans condition.
//
// Params:
//   - condition: condition pour traitement
func GoodDeferInCondition(condition bool) {
	if condition {
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf) // ✅ Correct

		processBuffer(buf)
	}
}

// GoodDeferInSwitch utilise defer dans switch.
//
// Params:
//   - choice: choix de traitement
func GoodDeferInSwitch(choice int) {
	switch choice {
	case 1:
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf) // ✅ Correct
		_ = buf
	case 2:
		buf := bufferPool.Get().([]byte)
		defer bufferPool.Put(buf) // ✅ Correct
		_ = buf
	}
}

// GoodDeferWithError utilise defer avant erreur potentielle.
//
// Returns:
//   - error: erreur de traitement
func GoodDeferWithError() error {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf) // ✅ Correct : protège même si erreur

	_ = buf
	// Retourne l'erreur de traitement
	return processWithError()
}

// processBuffer traite un buffer en remplissant avec des valeurs.
//
// Params:
//   - buf: buffer à traiter
func processBuffer(buf []byte) {
	for i := 0; i < len(buf); i++ {
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

// LargeConfig est une struct de 200 bytes.
type LargeConfig struct {
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

// HugeData est une struct de 512 bytes.
type HugeData struct {
	// Buffer contient les données brutes.
	Buffer [512]byte
}

// MediumStruct est une struct de 160 bytes.
type MediumStruct struct {
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

// SmallStruct est une struct de 24 bytes (OK par valeur).
type SmallStruct struct {
	// ID est l'identifiant.
	ID int // 8 bytes
	// Count est le compteur.
	Count int // 8 bytes
	// Active indique si actif.
	Active bool // 1 byte
}

// GoodProcessLargeConfig passe LargeConfig par pointeur.
//
// Params:
//   - cfg: configuration par pointeur (8 bytes)
func GoodProcessLargeConfig(cfg *LargeConfig) {
	_ = cfg.Host
}

// GoodProcessHugeData passe HugeData par pointeur.
//
// Params:
//   - data: données par pointeur (8 bytes)
func GoodProcessHugeData(data *HugeData) {
	_ = data.Buffer[0]
}

// GoodProcessMediumStruct passe MediumStruct par pointeur.
//
// Params:
//   - m: struct moyenne par pointeur (8 bytes)
func GoodProcessMediumStruct(m *MediumStruct) {
	_ = m.Flags
}

// GoodProcessSmallStruct passe SmallStruct par valeur (OK car petite).
//
// Params:
//   - s: petite struct par valeur (24 bytes)
func GoodProcessSmallStruct(s SmallStruct) {
	_ = s.ID
}

// GoodMultipleLargeParams a plusieurs paramètres par pointeur.
//
// Params:
//   - cfg: configuration par pointeur
//   - data: données par pointeur
func GoodMultipleLargeParams(cfg *LargeConfig, data *HugeData) {
	_ = cfg.Host
	_ = data.Buffer
}

// GoodMixedParams mélange types correctement.
//
// Params:
//   - id: identifiant
//   - cfg: configuration par pointeur
//   - name: nom
func GoodMixedParams(id int, cfg *LargeConfig, name string) {
	_ = id
	_ = cfg
	_ = name
}

// GoodReturnLargeStruct retourne pointeur vers large struct.
//
// Returns:
//   - *LargeConfig: configuration par pointeur
func GoodReturnLargeStruct() *LargeConfig {
	// Retourne une nouvelle configuration avec valeurs par défaut
	return &LargeConfig{
		Host: "localhost",
		Port: 8080,
	}
}

// GoodVariadicLargeStruct utilise variadic avec pointeurs.
//
// Params:
//   - configs: configurations par pointeurs
func GoodVariadicLargeStruct(configs ...*LargeConfig) {
	for _, cfg := range configs {
		_ = cfg
	}
}

// GoodAnonymousParam a un paramètre anonyme par pointeur.
//
// Params:
//   - (unnamed): configuration anonyme par pointeur
func GoodAnonymousParam(*LargeConfig) {
	// Paramètre non utilisé
}

// GoodMethodReceiver a un receiver par pointeur.
func (cfg *LargeConfig) GoodMethodReceiver() {
	_ = cfg.Host
}

// NestedStruct a des structs imbriquées larges.
type NestedStruct struct {
	// Config est la configuration.
	Config LargeConfig // 200 bytes
	// Data contient les données.
	Data HugeData // 512 bytes
}

// GoodProcessNestedStructs passe struct imbriquée par pointeur.
//
// Params:
//   - n: struct imbriquée par pointeur
func GoodProcessNestedStructs(n *NestedStruct) {
	_ = n.Config.Host
}

// GoodClosureWithLargeStruct crée closure avec pointeur.
//
// Returns:
//   - func(*LargeConfig): closure avec pointeur
func GoodClosureWithLargeStruct() func(*LargeConfig) {
	// Retourne une closure qui accepte un pointeur vers LargeConfig
	return func(cfg *LargeConfig) {
		_ = cfg.Host
	}
}

// GoodDeferWithLargeStruct utilise defer avec pointeur.
func GoodDeferWithLargeStruct() {
	defer func(cfg *LargeConfig) {
		_ = cfg.Host
	}(&LargeConfig{})
}

// GoodGoRoutineWithLargeStruct lance goroutine avec pointeur et synchronisation.
//
// Params:
//   - cfg: configuration par pointeur
func GoodGoRoutineWithLargeStruct(cfg *LargeConfig) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(c *LargeConfig) {
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
	s := SmallStruct{ID: 1, Count: 10} // ✅ Correct : petite struct
	GoodProcessSmallStruct(s)
}
