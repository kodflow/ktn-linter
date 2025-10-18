package slice001

func BadSliceIndexWithoutCheck(items []int, idx int) int {
	return items[idx] // want `\[KTN-SLICE-001\] Indexation du slice`
}

func BadSliceIndexMultiple(data []string) {
	i := 0
	j := 1
	first := data[i] // want `\[KTN-SLICE-001\] Indexation du slice`
	second := data[j] // want `\[KTN-SLICE-001\] Indexation du slice`
	_, _ = first, second
}

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
