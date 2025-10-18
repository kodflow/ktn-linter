package goto001

func BadGotoUsage(x int) {
	if x > 0 {
		goto skip // want `\[KTN-GOTO-001\].*`
	}
	process()
skip:
	return
}

func BadGotoInLoop() {
	for i := 0; i < 10; i++ {
		goto end // want `\[KTN-GOTO-001\].*`
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
