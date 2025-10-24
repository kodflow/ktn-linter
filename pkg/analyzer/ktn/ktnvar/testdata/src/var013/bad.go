package var013

// badLoopSliceAlloc crée un slice à l'intérieur d'une boucle for.
func badLoopSliceAlloc() {
	// Allocation dans boucle for
	for i := 0; i < 10; i++ {
		data := []int{}
		_ = data
	}
}

// badLoopMapAlloc crée une map à l'intérieur d'une boucle for.
func badLoopMapAlloc() {
	// Allocation dans boucle for
	for i := 0; i < 10; i++ {
		cache := make(map[string]int)
		_ = cache
	}
}

// badRangeSliceAlloc crée un slice à l'intérieur d'une boucle range.
func badRangeSliceAlloc() {
	items := []int{1, 2, 3}
	// Allocation dans boucle range
	for _, item := range items {
		buffer := []byte{}
		_ = item
		_ = buffer
	}
}

// badRangeMapAlloc crée une map à l'intérieur d'une boucle range.
func badRangeMapAlloc() {
	items := []string{"a", "b", "c"}
	// Allocation dans boucle range
	for _, item := range items {
		m := map[string]bool{}
		_ = item
		_ = m
	}
}

// badNestedLoopAlloc crée un slice dans une boucle imbriquée.
func badNestedLoopAlloc() {
	// Allocation dans boucle imbriquée
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			temp := make([]int, 10)
			_ = temp
			_ = j
		}
		_ = i
	}
}

// badVarDeclInLoop crée un slice avec var dans une boucle.
func badVarDeclInLoop() {
	// Allocation avec var dans boucle
	for i := 0; i < 10; i++ {
		var data = []int{1, 2, 3}
		_ = data
		_ = i
	}
}

// badVarMapDeclInLoop crée une map avec var dans une boucle.
func badVarMapDeclInLoop() {
	// Allocation avec var dans boucle
	for i := 0; i < 10; i++ {
		var cache = make(map[string]int)
		_ = cache
		_ = i
	}
}
