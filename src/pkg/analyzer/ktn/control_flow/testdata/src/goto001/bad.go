package goto001

func BadGotoUsage(x int) {
	if x > 0 { // want `\[KTN-CONTROL-GOTO-001\] Utilisation de goto non idiomatique`
		goto skip
	}
	process()
skip:
	return
}

func BadGotoInLoop() {
	for i := 0; i < 10; i++ {
		goto end // want `\[KTN-CONTROL-GOTO-001\] Utilisation de goto non idiomatique`
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
