package goroutine002

import (
	"context"
	"sync"
)

// Cas corrects - goroutines avec synchronisation

func GoodWithChannelSend() {
	ch := make(chan int)
	go func() {
		ch <- 42 // SendStmt - synchronisation via channel
	}()
	<-ch
}

func GoodWithChannelReceive() {
	ch := make(chan string)
	go func() {
		msg := <-ch // UnaryExpr avec <- - synchronisation via channel
		_ = msg
	}()
	ch <- "hello"
}

func GoodWithSelect() {
	ch := make(chan int)
	go func() {
		select { // SelectStmt - synchronisation via select
		case val := <-ch:
			_ = val
		}
	}()
	ch <- 10
}

func GoodWithContextParam(ctx context.Context) {
	go func(c context.Context) { // hasContextParam
		<-c.Done()
	}(ctx)
}

func GoodWithWaitGroupParam(wg *sync.WaitGroup) {
	wg.Add(1)
	go func(w *sync.WaitGroup) { // isWaitGroupTypeCheck param
		defer w.Done()
	}(wg)
}

func GoodMultipleSync() {
	var wg sync.WaitGroup
	ch := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ch <- true
	}()

	<-ch
	wg.Wait()
}
