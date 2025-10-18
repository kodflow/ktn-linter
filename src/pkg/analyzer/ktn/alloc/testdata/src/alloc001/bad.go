package alloc001

func BadNewWithSlice() {
	s := new([]int) // want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
	_ = s
}

func BadNewWithMap() {
	m := new(map[string]int) // want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
	_ = m
}

func BadNewWithChan() {
	ch := new(chan int) // want `\[KTN-ALLOC-001\] Utilisation de new\(\) avec un type référence`
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
