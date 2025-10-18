package slice001

// TODO: L'analyseur SLICE-001 a des limitations similaires à MAP-001
// Il ne détecte pas les accès sur des paramètres de fonction

func GoodSliceWithCheck(items []int) int {
	if len(items) > 0 {
		return items[0]
	}
	return 0
}

func GoodSliceIteration(data []string) {
	for i := range data {
		_ = data[i]
	}
}
