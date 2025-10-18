package slice001

// TODO: L'analyseur SLICE-001 a des limitations
// Il ne détecte actuellement que des cas très spécifiques
// Les tests ci-dessous correspondent à ce qu'il peut détecter

// Cas corrects - pas de détection nécessaire

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
