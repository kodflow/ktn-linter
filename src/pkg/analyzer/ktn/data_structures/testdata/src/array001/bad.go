package array002

func BadArraySize() {
	// want `\[KTN-DS-ARRAY-002\] Taille d'array incohérente`
	arr := [3]int{1, 2, 3, 4, 5}
}

func BadArraySizeString() {
	// want `\[KTN-DS-ARRAY-002\] Taille d'array incohérente`
	arr := [2]string{"a", "b", "c"}
	_ = arr
}

func GoodArraySize() {
	arr := [3]int{1, 2, 3}
	_ = arr
}

func GoodArrayEllipsis() {
	arr := [...]int{1, 2, 3, 4, 5}
	_ = arr
}
