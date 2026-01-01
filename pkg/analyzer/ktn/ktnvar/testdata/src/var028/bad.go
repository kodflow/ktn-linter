package var028

func badLoopCopy(items []int) {
	for _, v := range items {
		v := v // want "KTN-VAR-028"
		go process(v)
	}
}

func badIndexCopy(items []int) {
	for i, item := range items {
		i := i       // want "KTN-VAR-028"
		item := item // want "KTN-VAR-028"
		go func() {
			use(i, item)
		}()
	}
}

func process(v int) {}
func use(i, item int) {}
