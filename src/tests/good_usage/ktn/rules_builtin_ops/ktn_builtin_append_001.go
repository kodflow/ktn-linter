package rules_builtin_ops

// ✅ GOOD: assigner le résultat d'append
func correctAppend() {
	slice := []int{1, 2, 3}
	slice = append(slice, 4) // ✅ assigner!
	println(len(slice))      // 4
}

func appendInLoopGood() {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i) // ✅ assigner
	}
	println(len(s)) // 10
}
