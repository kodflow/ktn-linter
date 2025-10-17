// Package rules_alloc_bad contient des violations KTN-ALLOC-001.
package rules_alloc_bad

// Viole KTN-ALLOC-001 : new() avec map

// BadNewMapString crée une map avec new() au lieu de make().
func BadNewMapString() {
	m := new(map[string]int) // Viole KTN-ALLOC-001
	(*m)["key"] = 42         // PANIC car m pointe vers nil map
}

// BadNewMapInt crée une map d'entiers avec new().
func BadNewMapInt() {
	numbers := new(map[int]string) // Viole KTN-ALLOC-001
	(*numbers)[1] = "one"
}

// Viole KTN-ALLOC-001 : new() avec slice

// BadNewSliceInt crée un slice avec new().
func BadNewSliceInt() {
	s := new([]int) // Viole KTN-ALLOC-001
	*s = append(*s, 1, 2, 3)
}

// BadNewSliceString crée un slice de strings avec new().
func BadNewSliceString() {
	items := new([]string) // Viole KTN-ALLOC-001
	*items = append(*items, "hello")
}

// BadNewSliceStruct crée un slice de structs avec new().
func BadNewSliceStruct() {
	type Item struct {
		ID   int
		Name string
	}
	records := new([]Item) // Viole KTN-ALLOC-001
	*records = append(*records, Item{ID: 1, Name: "test"})
}

// Viole KTN-ALLOC-001 : new() avec channel

// BadNewChanInt crée un channel avec new().
func BadNewChanInt() {
	ch := new(chan int) // Viole KTN-ALLOC-001
	*ch <- 42           // PANIC car ch pointe vers nil channel
}

// BadNewChanString crée un channel de strings avec new().
func BadNewChanString() {
	messages := new(chan string) // Viole KTN-ALLOC-001
	*messages <- "hello"
}

// BadNewChanBuffered crée un channel bufferisé avec new().
func BadNewChanBuffered() {
	buf := new(chan []byte) // Viole KTN-ALLOC-001
	go func() {
		*buf <- []byte("data")
	}()
}

// Cas farfelus supplémentaires

// BadNewMapNested crée une map imbriquée avec new().
func BadNewMapNested() {
	nested := new(map[string]map[int]string) // Viole KTN-ALLOC-001
	(*nested)["outer"] = make(map[int]string)
}

// BadNewSlice2D crée un slice 2D avec new().
func BadNewSlice2D() {
	matrix := new([][]int) // Viole KTN-ALLOC-001
	*matrix = append(*matrix, []int{1, 2, 3})
}

// BadNewChanChan crée un channel de channels avec new().
func BadNewChanChan() {
	chCh := new(chan chan int) // Viole KTN-ALLOC-001
	innerCh := make(chan int)
	*chCh <- innerCh
}

// BadMultipleNewMaps crée plusieurs maps problématiques.
func BadMultipleNewMaps() {
	m1 := new(map[string]int)    // Viole KTN-ALLOC-001
	m2 := new(map[int]bool)      // Viole KTN-ALLOC-001
	m3 := new(map[string]string) // Viole KTN-ALLOC-001
	_, _, _ = m1, m2, m3
}
