package var013

// Constantes pour les tests
const (
	LOOP_COUNT_TEN    int = 10
	LOOP_COUNT_FIVE   int = 5
	SLICE_VALUE_ONE   int = 1
	SLICE_VALUE_TWO   int = 2
	SLICE_VALUE_THREE int = 3
	NESTED_SLICE_SIZE int = 10
)

// badLoopSliceAlloc crée un slice à l'intérieur d'une boucle for.
func badLoopSliceAlloc() {
	// Allocation dans boucle for
	for index := 0; index < LOOP_COUNT_TEN; index++ {
		dataSlice := make([]int, 0, LOOP_COUNT_TEN)
		_ = dataSlice
	}
}

// badLoopMapAlloc crée une map à l'intérieur d'une boucle for.
func badLoopMapAlloc() {
	// Allocation dans boucle for
	for loopIndex := 0; loopIndex < LOOP_COUNT_TEN; loopIndex++ {
		cacheMap := make(map[string]int, LOOP_COUNT_TEN)
		_ = cacheMap
	}
}

// badRangeSliceAlloc crée un slice à l'intérieur d'une boucle range.
func badRangeSliceAlloc() {
	items := []int{SLICE_VALUE_ONE, SLICE_VALUE_TWO, SLICE_VALUE_THREE}
	// Allocation dans boucle range
	for _, itemValue := range items {
		stringBuffer := make([]string, 0, LOOP_COUNT_TEN)
		_ = itemValue
		_ = stringBuffer
	}
}

// badRangeMapAlloc crée une map à l'intérieur d'une boucle range.
func badRangeMapAlloc() {
	items := []string{"a", "b", "c"}
	// Allocation dans boucle range
	for _, stringItem := range items {
		boolMap := make(map[string]bool, LOOP_COUNT_TEN)
		_ = stringItem
		_ = boolMap
	}
}

// badNestedLoopAlloc crée un slice dans une boucle imbriquée.
func badNestedLoopAlloc() {
	// Allocation dans boucle imbriquée
	for outerIndex := 0; outerIndex < LOOP_COUNT_FIVE; outerIndex++ {
		// Boucle interne avec allocation
		for innerIndex := 0; innerIndex < LOOP_COUNT_FIVE; innerIndex++ {
			tempSlice := make([]int, 0, NESTED_SLICE_SIZE)
			_ = tempSlice
			_ = innerIndex
		}
		_ = outerIndex
	}
}

// badVarDeclInLoop crée un slice avec var dans une boucle.
func badVarDeclInLoop() {
	// Allocation avec var dans boucle
	for varIndex := 0; varIndex < LOOP_COUNT_TEN; varIndex++ {
		dataArray := make([]int, 0, SLICE_VALUE_THREE)
		dataArray = append(dataArray, SLICE_VALUE_ONE, SLICE_VALUE_TWO, SLICE_VALUE_THREE)
		_ = dataArray
		_ = varIndex
	}
}

// badVarMapDeclInLoop crée une map avec var dans une boucle.
func badVarMapDeclInLoop() {
	// Allocation avec var dans boucle
	for mapIndex := 0; mapIndex < LOOP_COUNT_TEN; mapIndex++ {
		stringCache := make(map[string]int, LOOP_COUNT_TEN)
		_ = stringCache
		_ = mapIndex
	}
}
