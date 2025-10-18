package for001

// Cas détectés par FOR-001: utilisation inutile de _ dans range

func BadForIndexUnderscore() {
	items := []int{1, 2, 3}
	for i, _ := range items { // want `\[KTN-FOR-001\].*`
		_ = i
	}
}

func BadForIndexUnderscoreInLoop() {
	data := []string{"a", "b", "c"}
	for idx, _ := range data { // want `\[KTN-FOR-001\].*`
		println(idx)
	}
}

// Cas corrects pour référence

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

func GoodUseBothVars() {
	items := []int{1, 2, 3}
	for i, v := range items {
		_, _ = i, v
	}
}

func GoodUseValueOnly() {
	items := []int{1, 2, 3}
	for _, v := range items {
		_ = v
	}
}

func process() {}
