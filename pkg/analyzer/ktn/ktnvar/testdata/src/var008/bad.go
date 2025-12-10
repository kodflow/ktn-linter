// Bad examples for the var009 test case.
package var008

// Constantes pour les tests
const (
	LoopCountTen    int = 10
	LoopCountFive   int = 5
	SliceValueOne   int = 1
	SliceValueTwo   int = 2
	SliceValueThree int = 3
	NestedSliceSize int = 10
)

// badLoopSliceAlloc crée un slice à l'intérieur d'une boucle for.
func badLoopSliceAlloc() {
	// Allocation dans boucle for
	for index := 0; index < LoopCountTen; index++ {
		dataSlice := make([]int, 0, LoopCountTen)
		_ = dataSlice
	}
}

// badLoopMapAlloc crée une map à l'intérieur d'une boucle for.
func badLoopMapAlloc() {
	// Allocation dans boucle for
	for loopIndex := 0; loopIndex < LoopCountTen; loopIndex++ {
		cacheMap := make(map[string]int, LoopCountTen)
		_ = cacheMap
	}
}

// badRangeSliceAlloc crée un slice à l'intérieur d'une boucle range.
func badRangeSliceAlloc() {
	items := []int{SliceValueOne, SliceValueTwo, SliceValueThree}
	// Allocation dans boucle range
	for _, itemValue := range items {
		stringBuffer := make([]string, 0, LoopCountTen)
		_ = itemValue
		_ = stringBuffer
	}
}

// badRangeMapAlloc crée une map à l'intérieur d'une boucle range.
func badRangeMapAlloc() {
	items := []string{"a", "b", "c"}
	// Allocation dans boucle range
	for _, stringItem := range items {
		boolMap := make(map[string]bool, LoopCountTen)
		_ = stringItem
		_ = boolMap
	}
}

// badNestedLoopAlloc crée un slice dans une boucle imbriquée.
func badNestedLoopAlloc() {
	// Allocation dans boucle imbriquée
	for outerIndex := 0; outerIndex < LoopCountFive; outerIndex++ {
		// Boucle interne avec allocation
		for innerIndex := 0; innerIndex < LoopCountFive; innerIndex++ {
			tempSlice := make([]int, 0, NestedSliceSize)
			_ = tempSlice
			_ = innerIndex
		}
		_ = outerIndex
	}
}

// badVarDeclInLoop crée un slice avec var dans une boucle.
func badVarDeclInLoop() {
	// Allocation avec var dans boucle
	for varIndex := 0; varIndex < LoopCountTen; varIndex++ {
		dataArray := make([]int, 0, SliceValueThree)
		dataArray = append(dataArray, SliceValueOne, SliceValueTwo, SliceValueThree)
		_ = dataArray
		_ = varIndex
	}
}

// badVarMapDeclInLoop crée une map avec var dans une boucle.
func badVarMapDeclInLoop() {
	// Allocation avec var dans boucle
	for mapIndex := 0; mapIndex < LoopCountTen; mapIndex++ {
		stringCache := make(map[string]int, LoopCountTen)
		_ = stringCache
		_ = mapIndex
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badLoopSliceAlloc
	badLoopSliceAlloc()
	// Appel de badLoopMapAlloc
	badLoopMapAlloc()
	// Appel de badRangeSliceAlloc
	badRangeSliceAlloc()
	// Appel de badRangeMapAlloc
	badRangeMapAlloc()
	// Appel de badNestedLoopAlloc
	badNestedLoopAlloc()
	// Appel de badVarDeclInLoop
	badVarDeclInLoop()
	// Appel de badVarMapDeclInLoop
	badVarMapDeclInLoop()
}
