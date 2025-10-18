package for001

// TODO: L'analyseur FOR-001 nécessite des améliorations
// Il détecte uniquement les cas avec := pas avec =
// Les vrais cas bad seraient:
// - for _ = range items (mais n'est pas détecté)
// - for _, _ = range items (mais n'est pas détecté)

// Ce sont des exemples valides pour tester la détection
func GoodUseIndexOnly() {
	items := []int{1, 2, 3}
	for i := range items {
		_ = i
	}
}

func GoodUseNoVars() {
	items := []int{1, 2, 3}
	for range items {
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
