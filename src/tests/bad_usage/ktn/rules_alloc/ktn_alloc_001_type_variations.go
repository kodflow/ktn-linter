// Package rules_alloc_bad contient des violations KTN-ALLOC-001 avec variations de types.
package rules_alloc_bad

// Viole KTN-ALLOC-001 : new() avec types complexes et custom types

// BadNewMapComplex128 crée une map avec clés complex128 via new().
func BadNewMapComplex128() {
	m := new(map[complex128]string) // Viole KTN-ALLOC-001
	(*m)[complex(1.0, 2.0)] = "complex"
}

// BadNewMapComplex64 crée une map avec valeurs complex64 via new().
func BadNewMapComplex64() {
	coords := new(map[string]complex64) // Viole KTN-ALLOC-001
	(*coords)["point"] = complex64(complex(3.0, 4.0))
}

// BadNewSliceComplex crée un slice de complex128 via new().
func BadNewSliceComplex() {
	numbers := new([]complex128) // Viole KTN-ALLOC-001
	*numbers = append(*numbers, complex(1, 2), complex(3, 4))
}

// BadNewChanComplex crée un channel de complex64 via new().
func BadNewChanComplex() {
	ch := new(chan complex64) // Viole KTN-ALLOC-001
	*ch <- complex64(complex(5, 6))
}

// Custom types

// MyMap est un type personnalisé basé sur map.
type MyMap map[string]int

// MySlice est un type personnalisé basé sur slice.
type MySlice []int

// MyChan est un type personnalisé basé sur channel.
type MyChan chan string

// BadNewCustomMap crée un type map personnalisé via new().
func BadNewCustomMap() {
	m := new(MyMap) // Viole KTN-ALLOC-001
	(*m)["key"] = 42
}

// BadNewCustomSlice crée un type slice personnalisé via new().
func BadNewCustomSlice() {
	s := new(MySlice) // Viole KTN-ALLOC-001
	*s = append(*s, 1, 2, 3)
}

// BadNewCustomChan crée un type channel personnalisé via new().
func BadNewCustomChan() {
	ch := new(MyChan) // Viole KTN-ALLOC-001
	*ch <- "message"
}

// Function types

// BadNewMapFunc crée une map de fonctions via new().
func BadNewMapFunc() {
	handlers := new(map[string]func()) // Viole KTN-ALLOC-001
	(*handlers)["start"] = func() {}
}

// BadNewSliceFunc crée un slice de fonctions via new().
func BadNewSliceFunc() {
	callbacks := new([]func(int) string) // Viole KTN-ALLOC-001
	*callbacks = append(*callbacks, func(n int) string { return "" })
}

// BadNewChanFunc crée un channel de fonctions via new().
func BadNewChanFunc() {
	ch := new(chan func()) // Viole KTN-ALLOC-001
	*ch <- func() {}
}

// Interface types

// Processor est une interface de traitement.
type Processor interface {
	Process(data string) string
}

// BadNewMapInterface crée une map avec interface via new().
func BadNewMapInterface() {
	processors := new(map[string]Processor) // Viole KTN-ALLOC-001
	(*processors)["main"] = nil
}

// BadNewSliceInterface crée un slice d'interfaces via new().
func BadNewSliceInterface() {
	items := new([]Processor) // Viole KTN-ALLOC-001
	*items = append(*items, nil)
}

// BadNewChanInterface crée un channel d'interfaces via new().
func BadNewChanInterface() {
	ch := new(chan Processor) // Viole KTN-ALLOC-001
	*ch <- nil
}

// Nested custom types

// BadNewMapOfCustomSlices crée une map de MySlice via new().
func BadNewMapOfCustomSlices() {
	data := new(map[string]MySlice) // Viole KTN-ALLOC-001
	(*data)["numbers"] = MySlice{1, 2, 3}
}

// BadNewSliceOfCustomMaps crée un slice de MyMap via new().
func BadNewSliceOfCustomMaps() {
	configs := new([]MyMap) // Viole KTN-ALLOC-001
	*configs = append(*configs, MyMap{"key": 1})
}

// BadNewChanOfCustomTypes crée un channel de MySlice via new().
func BadNewChanOfCustomTypes() {
	ch := new(chan MySlice) // Viole KTN-ALLOC-001
	*ch <- MySlice{1, 2, 3}
}

// Rune and byte variations

// BadNewMapRune crée une map avec rune via new().
func BadNewMapRune() {
	chars := new(map[rune]int) // Viole KTN-ALLOC-001
	(*chars)['a'] = 1
}

// BadNewSliceRune crée un slice de runes via new().
func BadNewSliceRune() {
	text := new([]rune) // Viole KTN-ALLOC-001
	*text = append(*text, 'h', 'e', 'l', 'l', 'o')
}

// BadNewChanByte crée un channel de bytes via new().
func BadNewChanByte() {
	ch := new(chan byte) // Viole KTN-ALLOC-001
	*ch <- byte('x')
}

// Pointer types in collections

// BadNewMapPointer crée une map de pointeurs via new().
func BadNewMapPointer() {
	ptrs := new(map[string]*int) // Viole KTN-ALLOC-001
	val := 42
	(*ptrs)["answer"] = &val
}

// BadNewSlicePointer crée un slice de pointeurs via new().
func BadNewSlicePointer() {
	nums := new([]*int) // Viole KTN-ALLOC-001
	val := 42
	*nums = append(*nums, &val)
}

// BadNewChanPointer crée un channel de pointeurs via new().
func BadNewChanPointer() {
	ch := new(chan *string) // Viole KTN-ALLOC-001
	msg := "hello"
	*ch <- &msg
}
