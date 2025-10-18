package for001

func BadUnderscoreForIndex() {
	items := []int{1, 2, 3}
	for _ = range items { // want `\[KTN-CONTROL-FOR-001\] Utilisation de _ inutile dans for range`
		process()
	}
}

func BadUnderscoreForValue() {
	items := []string{"a", "b", "c"}
	for _, _ = range items { // want `\[KTN-CONTROL-FOR-001\] Utilisation de _ inutile dans for range`
		process()
	}
}

func GoodNoUnderscore() {
	items := []int{1, 2, 3}
	for range items {
		process()
	}
}

func GoodWithIndex() {
	items := []int{1, 2, 3}
	for i := range items {
		_ = i
	}
}

func process() {}
