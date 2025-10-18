package goto001

func BadGotoUsage(x int) {
	// want `\[KTN-CONTROL-GOTO-001\] Utilisation de goto non idiomatique`
	if x > 0 {
		goto skip
	}
	process()
skip:
	return
}

func BadGotoInLoop() {
	for i := 0; i < 10; i++ {
		// want `\[KTN-CONTROL-GOTO-001\] Utilisation de goto non idiomatique`
		goto end
	}
end:
}

func GoodWithoutGoto(x int) {
	if x <= 0 {
		return
	}
	process()
}

func process() {}
