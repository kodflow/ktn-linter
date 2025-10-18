package chan001

func BadReceiverCloses(ch chan int) {
	for v := range ch {
		process(v)
	}
	close(ch) // want `\[KTN-CHAN-001\] close\(\) appelé par le receiver`
}

func BadReceiverClosesAfterReceive(ch chan string) {
	val := <-ch
	_ = val
	close(ch) // want `\[KTN-CHAN-001\] close\(\) appelé par le receiver`
}

func GoodSenderOnlyCloses(ch chan int) {
	ch <- 42
	close(ch)
}

func process(v int) {}
