package for001

// Cas corrects - pas de détection FOR-001

// GoodRangeNoVars - range sans variables
func GoodRangeNoVars() {
	items := []int{1, 2, 3}
	for range items {
		process()
	}
}

// GoodRangeIndexOnly - range avec index seulement
func GoodRangeIndexOnly() {
	items := []int{1, 2, 3}
	for i := range items {
		_ = i
	}
}

// GoodRangeBothVars - range avec index et valeur
func GoodRangeBothVars() {
	items := []int{1, 2, 3}
	for i, v := range items {
		_, _ = i, v
	}
}

// GoodRangeValueUnderscore - range avec _ pour index (CORRECT)
func GoodRangeValueUnderscore() {
	items := []int{1, 2, 3}
	for _, v := range items {
		_ = v
	}
}

// GoodMapRange - range sur map
func GoodMapRange() {
	m := map[string]int{"a": 1}
	for k, v := range m {
		_, _ = k, v
	}
}

// GoodMapRangeKeyOnly - range sur map (clés seulement)
func GoodMapRangeKeyOnly() {
	m := map[string]int{"a": 1}
	for k := range m {
		_ = k
	}
}

// GoodChannelRange - range sur channel
func GoodChannelRange() {
	ch := make(chan int)
	go func() {
		close(ch)
	}()
	for v := range ch {
		_ = v
	}
}
