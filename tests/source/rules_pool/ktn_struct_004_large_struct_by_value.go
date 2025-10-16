// Package rules_pool_bad contient des violations de KTN-STRUCT-004.
package rules_pool_bad

// ❌ KTN-STRUCT-004 : Grandes structs passées par valeur

// LargeConfig est une struct de 200 bytes.
type LargeConfig struct {
	Host     string    // 16 bytes
	Port     int       // 8 bytes
	Timeout  int       // 8 bytes
	Buffer   [128]byte // 128 bytes
	MaxConns int       // 8 bytes
	Retries  int       // 8 bytes
	// Total: ~176 bytes (mais avec padding > 128)
}

// HugeData est une struct de 512 bytes.
type HugeData struct {
	Buffer [512]byte // 512 bytes
}

// MediumStruct est une struct de 160 bytes.
type MediumStruct struct {
	Data1 [64]byte // 64 bytes
	Data2 [64]byte // 64 bytes
	Flags int      // 8 bytes
	Count int      // 8 bytes
	Size  int      // 8 bytes
	// Total: ~152 bytes
}

// ComplexStruct est une struct avec plusieurs champs.
type ComplexStruct struct {
	ID        int64     // 8 bytes
	Name      string    // 16 bytes
	Data      [100]byte // 100 bytes
	Timestamp int64     // 8 bytes
	Active    bool      // 1 byte
	Priority  int       // 8 bytes
	// Total: ~141 bytes
}

// BadProcessLargeConfig passe LargeConfig par valeur.
//
// Params:
//   - cfg: configuration (copie 200 bytes) // Viole KTN-STRUCT-004
func BadProcessLargeConfig(cfg LargeConfig) {
	_ = cfg.Host
}

// BadProcessHugeData passe HugeData par valeur.
//
// Params:
//   - data: données (copie 512 bytes) // Viole KTN-STRUCT-004
func BadProcessHugeData(data HugeData) {
	_ = data.Buffer[0]
}

// BadProcessMediumStruct passe MediumStruct par valeur.
//
// Params:
//   - m: struct moyenne (copie 160 bytes) // Viole KTN-STRUCT-004
func BadProcessMediumStruct(m MediumStruct) {
	_ = m.Flags
}

// BadProcessComplexStruct passe ComplexStruct par valeur.
//
// Params:
//   - c: struct complexe (copie 141 bytes) // Viole KTN-STRUCT-004
func BadProcessComplexStruct(c ComplexStruct) {
	_ = c.ID
}

// BadMultipleLargeParams a plusieurs paramètres larges.
//
// Params:
//   - cfg: configuration // Viole KTN-STRUCT-004
//   - data: données // Viole KTN-STRUCT-004
func BadMultipleLargeParams(cfg LargeConfig, data HugeData) {
	_ = cfg.Host
	_ = data.Buffer
}

// BadMixedParams mélange petits et grands paramètres.
//
// Params:
//   - id: identifiant
//   - cfg: configuration // Viole KTN-STRUCT-004
//   - name: nom
func BadMixedParams(id int, cfg LargeConfig, name string) {
	_ = id
	_ = cfg
	_ = name
}

// BadReturnLargeStruct retourne large struct.
//
// Returns:
//   - LargeConfig: configuration (copie au retour)
func BadReturnLargeStruct() LargeConfig {
	// Note: return type n'est pas détecté par notre analyseur (seulement params)
	return LargeConfig{
		Host: "localhost",
		Port: 8080,
	}
}

// BadVariadicLargeStruct utilise variadic avec large struct.
//
// Params:
//   - configs: configurations // Viole KTN-STRUCT-004 pour chaque élément
func BadVariadicLargeStruct(configs ...LargeConfig) {
	for _, cfg := range configs {
		_ = cfg
	}
}

// BadAnonymousParam a un paramètre anonyme large.
//
// Params:
//   - (unnamed): configuration anonyme // Viole KTN-STRUCT-004
func BadAnonymousParam(LargeConfig) {
	// Paramètre non utilisé
}

// BadMethodReceiver a un receiver large (method, pas function).
//
// Note: les méthodes ne sont généralement pas vérifiées par notre analyseur
func (cfg LargeConfig) BadMethodReceiver() {
	_ = cfg.Host
}

// BadNestedStructs a des structs imbriquées larges.
type NestedStruct struct {
	Config LargeConfig // 200 bytes
	Data   HugeData    // 512 bytes
	// Total: 712 bytes
}

// BadProcessNestedStructs passe struct imbriquée par valeur.
//
// Params:
//   - n: struct imbriquée (copie 712 bytes) // Viole KTN-STRUCT-004
func BadProcessNestedStructs(n NestedStruct) {
	_ = n.Config.Host
}

// BadArrayOfStructs a un array de structs larges.
type ArrayOfStructs struct {
	Items [4]LargeConfig // 4 * 200 = 800 bytes
}

// BadProcessArrayOfStructs passe array de structs par valeur.
//
// Params:
//   - a: array de structs (copie 800 bytes) // Viole KTN-STRUCT-004
func BadProcessArrayOfStructs(a ArrayOfStructs) {
	_ = a.Items[0]
}

// BadInlineStruct utilise struct inline large.
//
// Params:
//   - data: struct inline // Viole KTN-STRUCT-004
func BadInlineStruct(data struct {
	Buffer [200]byte
	Count  int
}) {
	_ = data.Count
}

// BadClosureWithLargeStruct crée closure avec large param.
func BadClosureWithLargeStruct() func(LargeConfig) {
	// Viole KTN-STRUCT-004 dans la closure
	return func(cfg LargeConfig) {
		_ = cfg.Host
	}
}

// BadDeferWithLargeStruct utilise defer avec large struct.
func BadDeferWithLargeStruct() {
	defer func(cfg LargeConfig) { // Viole KTN-STRUCT-004
		_ = cfg.Host
	}(LargeConfig{})
}

// BadGoRoutineWithLargeStruct lance goroutine avec large struct.
func BadGoRoutineWithLargeStruct(cfg LargeConfig) {
	// Param cfg viole KTN-STRUCT-004
	go func(c LargeConfig) { // Viole KTN-STRUCT-004 aussi
		_ = c.Host
	}(cfg)
}
