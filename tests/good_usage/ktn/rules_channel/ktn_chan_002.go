package rules_channel

// ✅ GOOD: sender ferme, receiver lit
func sender(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch) // ✅ sender ferme quand fini
}

func receiverGood(ch chan int) {
	for v := range ch { // ✅ détecte fermeture automatiquement
		processGood(v)
	}
}

func processGood(v int) {}
