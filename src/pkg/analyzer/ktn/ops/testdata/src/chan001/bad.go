package chan002

func BadReceiverCloses(ch chan int) {
	// want `\[KTN-OPS-CHAN-002\] close\(\) appelé par le receiver`
	for v := range ch {
		process(v)
	}
	close(ch)
}

func BadReceiverClosesAfterReceive(ch chan string) {
	val := <-ch
	_ = val
	// want `\[KTN-OPS-CHAN-002\] close\(\) appelé par le receiver`
	close(ch)
}

func GoodSenderCloses() {
	ch := make(chan int)
	go func() {
		ch <- 42
		close(ch)
	}()
	<-ch
}

func process(v int) {}
