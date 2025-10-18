package range003

func BadCaptureInClosure() {
	items := []int{1, 2, 3}
	var funcs []func()
	for _, item := range items {
		funcs = append(funcs, func() { // want `\[KTN-CONTROL-RANGE-003\] Variable de range capturée dans une closure`
			process(item) // item est capturé
		})
	}
}

func BadCaptureInGoroutine() {
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() { // want `\[KTN-CONTROL-RANGE-003\] Variable de range capturée dans une closure`
			process(v) // v est capturé
		}()
	}
}

func GoodLocalCopy() {
	items := []int{1, 2, 3}
	var funcs []func()
	for _, item := range items {
		item := item // Copie locale
		funcs = append(funcs, func() {
			process(item)
		})
	}
}

func GoodParameter() {
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func(val string) {
			process(val)
		}(v)
	}
}

func process(v interface{}) {}
