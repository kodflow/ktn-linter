package var013

// badLoopSliceAlloc crée un slice à l'intérieur d'une boucle for.
func badLoopSliceAlloc() {
	// Allocation dans boucle for
	for i := 0; i < 10; i++ {
		data := []int{} // want "KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle"
		_ = data
	}
}

// badLoopMapAlloc crée une map à l'intérieur d'une boucle for.
func badLoopMapAlloc() {
	// Allocation dans boucle for
	for i := 0; i < 10; i++ {
		cache := make(map[string]int) // want "KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle"
		_ = cache
	}
}

// badRangeSliceAlloc crée un slice à l'intérieur d'une boucle range.
func badRangeSliceAlloc() {
	items := []int{1, 2, 3}
	// Allocation dans boucle range
	for _, item := range items {
		buffer := []byte{} // want "KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle"
		_ = item
		_ = buffer
	}
}

// badRangeMapAlloc crée une map à l'intérieur d'une boucle range.
func badRangeMapAlloc() {
	items := []string{"a", "b", "c"}
	// Allocation dans boucle range
	for _, item := range items {
		m := map[string]bool{} // want "KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle"
		_ = item
		_ = m
	}
}

// badNestedLoopAlloc crée un slice dans une boucle imbriquée.
func badNestedLoopAlloc() {
	// Allocation dans boucle imbriquée
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			temp := make([]int, 10) // want "KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle"
			_ = temp
			_ = j
		}
		_ = i
	}
}
