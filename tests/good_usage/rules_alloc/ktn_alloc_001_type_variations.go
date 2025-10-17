// Package rules_alloc_good contient du code conforme KTN-ALLOC-001 avec variations de types.
package rules_alloc_good

// ✅ Code conforme KTN-ALLOC-001 : make() pour types complexes et custom types

// GoodMakeMapComplex128 crée une map avec clés complex128 via make().
//
// Returns:
//   - map[complex128]string: map de nombres complexes vers strings
func GoodMakeMapComplex128() map[complex128]string {
	m := make(map[complex128]string) // ✅ Correct : make() pour map
	m[complex(1.0, 2.0)] = "complex"
	return m
}

// GoodMakeMapComplex64 crée une map avec valeurs complex64 via make().
//
// Returns:
//   - map[string]complex64: map de strings vers nombres complexes
func GoodMakeMapComplex64() map[string]complex64 {
	coords := make(map[string]complex64, 10) // ✅ Correct : make() avec capacité
	coords["point"] = complex64(complex(3.0, 4.0))
	return coords
}

// GoodMakeSliceComplex crée un slice de complex128 via make().
//
// Returns:
//   - []complex128: slice de nombres complexes
func GoodMakeSliceComplex() []complex128 {
	numbers := make([]complex128, 0, 5) // ✅ Correct : make() avec capacité
	numbers = append(numbers, complex(1, 2), complex(3, 4))
	return numbers
}

// GoodMakeChanComplex crée un channel de complex64 via make().
//
// Returns:
//   - chan complex64: channel de nombres complexes
func GoodMakeChanComplex() chan complex64 {
	ch := make(chan complex64, 3) // ✅ Correct : make() pour channel
	go func() {
		ch <- complex64(complex(5, 6))
	}()
	return ch
}

// Custom types

// MyMap est un type personnalisé basé sur map.
type MyMap map[string]int

// MySlice est un type personnalisé basé sur slice.
type MySlice []int

// MyChan est un type personnalisé basé sur channel.
type MyChan chan string

// GoodMakeCustomMap crée un type map personnalisé via make().
//
// Returns:
//   - MyMap: type map personnalisé
func GoodMakeCustomMap() MyMap {
	m := make(MyMap) // ✅ Correct : make() pour type custom
	m["key"] = 42
	return m
}

// GoodMakeCustomSlice crée un type slice personnalisé via make().
//
// Returns:
//   - MySlice: type slice personnalisé
func GoodMakeCustomSlice() MySlice {
	s := make(MySlice, 0, 10) // ✅ Correct : make() avec capacité
	s = append(s, 1, 2, 3)
	return s
}

// GoodMakeCustomChan crée un type channel personnalisé via make().
//
// Returns:
//   - MyChan: type channel personnalisé
func GoodMakeCustomChan() MyChan {
	ch := make(MyChan, 5) // ✅ Correct : make() pour channel custom
	go func() {
		ch <- "message"
	}()
	return ch
}

// Function types

// GoodMakeMapFunc crée une map de fonctions via make().
//
// Returns:
//   - map[string]func(): map de fonctions
func GoodMakeMapFunc() map[string]func() {
	handlers := make(map[string]func(), 5) // ✅ Correct : make() pour map
	handlers["start"] = func() {}
	return handlers
}

// GoodMakeSliceFunc crée un slice de fonctions via make().
//
// Returns:
//   - []func(int) string: slice de fonctions
func GoodMakeSliceFunc() []func(int) string {
	callbacks := make([]func(int) string, 0, 3) // ✅ Correct : make() avec capacité
	callbacks = append(callbacks, func(n int) string { return "" })
	return callbacks
}

// GoodMakeChanFunc crée un channel de fonctions via make().
//
// Returns:
//   - chan func(): channel de fonctions
func GoodMakeChanFunc() chan func() {
	ch := make(chan func(), 2) // ✅ Correct : make() pour channel
	go func() {
		ch <- func() {}
	}()
	return ch
}

// Interface types

// Processor est une interface de traitement.
type Processor interface {
	// Process traite des données.
	//
	// Params:
	//   - data: données à traiter
	//
	// Returns:
	//   - string: résultat du traitement
	Process(data string) string
}

// GoodMakeMapInterface crée une map avec interface via make().
//
// Returns:
//   - map[string]Processor: map d'interfaces
func GoodMakeMapInterface() map[string]Processor {
	processors := make(map[string]Processor) // ✅ Correct : make() pour map
	processors["main"] = nil
	return processors
}

// GoodMakeSliceInterface crée un slice d'interfaces via make().
//
// Returns:
//   - []Processor: slice d'interfaces
func GoodMakeSliceInterface() []Processor {
	items := make([]Processor, 0, 10) // ✅ Correct : make() avec capacité
	items = append(items, nil)
	return items
}

// GoodMakeChanInterface crée un channel d'interfaces via make().
//
// Returns:
//   - chan Processor: channel d'interfaces
func GoodMakeChanInterface() chan Processor {
	ch := make(chan Processor, 3) // ✅ Correct : make() pour channel
	go func() {
		ch <- nil
	}()
	return ch
}

// Nested custom types

// GoodMakeMapOfCustomSlices crée une map de MySlice via make().
//
// Returns:
//   - map[string]MySlice: map de slices personnalisés
func GoodMakeMapOfCustomSlices() map[string]MySlice {
	data := make(map[string]MySlice, 5) // ✅ Correct : make() avec capacité
	data["numbers"] = MySlice{1, 2, 3}
	return data
}

// GoodMakeSliceOfCustomMaps crée un slice de MyMap via make().
//
// Returns:
//   - []MyMap: slice de maps personnalisées
func GoodMakeSliceOfCustomMaps() []MyMap {
	configs := make([]MyMap, 0, 3) // ✅ Correct : make() avec capacité
	configs = append(configs, MyMap{"key": 1})
	return configs
}

// GoodMakeChanOfCustomTypes crée un channel de MySlice via make().
//
// Returns:
//   - chan MySlice: channel de slices personnalisés
func GoodMakeChanOfCustomTypes() chan MySlice {
	ch := make(chan MySlice, 2) // ✅ Correct : make() pour channel
	go func() {
		ch <- MySlice{1, 2, 3}
	}()
	return ch
}

// Rune and byte variations

// GoodMakeMapRune crée une map avec rune via make().
//
// Returns:
//   - map[rune]int: map de runes vers int
func GoodMakeMapRune() map[rune]int {
	chars := make(map[rune]int, 26) // ✅ Correct : make() avec capacité
	chars['a'] = 1
	return chars
}

// GoodMakeSliceRune crée un slice de runes via make().
//
// Returns:
//   - []rune: slice de runes
func GoodMakeSliceRune() []rune {
	text := make([]rune, 0, 10) // ✅ Correct : make() avec capacité
	text = append(text, 'h', 'e', 'l', 'l', 'o')
	return text
}

// GoodMakeChanByte crée un channel de bytes via make().
//
// Returns:
//   - chan byte: channel de bytes
func GoodMakeChanByte() chan byte {
	ch := make(chan byte, 5) // ✅ Correct : make() pour channel
	go func() {
		ch <- byte('x')
	}()
	return ch
}

// Pointer types in collections

// GoodMakeMapPointer crée une map de pointeurs via make().
//
// Returns:
//   - map[string]*int: map de pointeurs
func GoodMakeMapPointer() map[string]*int {
	ptrs := make(map[string]*int, 5) // ✅ Correct : make() avec capacité
	val := 42
	ptrs["answer"] = &val
	return ptrs
}

// GoodMakeSlicePointer crée un slice de pointeurs via make().
//
// Returns:
//   - []*int: slice de pointeurs
func GoodMakeSlicePointer() []*int {
	nums := make([]*int, 0, 10) // ✅ Correct : make() avec capacité
	val := 42
	nums = append(nums, &val)
	return nums
}

// GoodMakeChanPointer crée un channel de pointeurs via make().
//
// Returns:
//   - chan *string: channel de pointeurs
func GoodMakeChanPointer() chan *string {
	ch := make(chan *string, 3) // ✅ Correct : make() pour channel
	go func() {
		msg := "hello"
		ch <- &msg
	}()
	return ch
}
