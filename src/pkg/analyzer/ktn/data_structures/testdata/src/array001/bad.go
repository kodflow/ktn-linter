package array001

func BadArraySize() {
	// Cette fonction teste que l'analyzer détecte les arrays avec taille explicite
	// qui devraient utiliser ... à la place
	_ = [5]int{1, 2, 3, 4, 5}
}

func GoodArraySizeExplicit() {
	_ = [3]int{1, 2, 3}
}

func GoodArraySize() {
	arr := [3]int{1, 2, 3}
	_ = arr
}

func GoodArrayEllipsis() {
	arr := [...]int{1, 2, 3, 4, 5}
	_ = arr
}
