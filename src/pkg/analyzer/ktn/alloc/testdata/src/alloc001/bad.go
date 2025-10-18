package alloc001

func BadNewWithSlice() {
	// want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
	s := new([]int)
	_ = s
}

func BadNewWithMap() {
	// want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
	m := new(map[string]int)
	_ = m
}

func BadNewWithChan() {
	// want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
	ch := new(chan int)
	_ = ch
}

func GoodMakeSlice() {
	s := make([]int, 0, 10)
	_ = s
}

func GoodMakeMap() {
	m := make(map[string]int)
	_ = m
}

func GoodMakeChan() {
	ch := make(chan int)
	_ = ch
}
