package slice001

func BadSliceIndexWithoutCheck(items []int) int {
	// want `\[KTN-DS-SLICE-001\] Index de slice sans vérification de bounds`
	return items[0]
}

func BadSliceIndexMultiple(data []string) {
	// want `\[KTN-DS-SLICE-001\] Index de slice sans vérification de bounds`
	first := data[0]
	// want `\[KTN-DS-SLICE-001\] Index de slice sans vérification de bounds`
	second := data[1]
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
