package goroutine002

import (
	"context"
	"sync"
)

func BadNoSync() {
	data := "test"
	// want `\[KTN-GOROUTINE-002\] Goroutine lancée sans mécanisme de synchronisation`
	go func() {
		process(data)
	}()
	// La fonction peut se terminer avant la goroutine
}

func BadNoSyncMultiple() {
	for i := 0; i < 5; i++ {
		// want `\[KTN-GOROUTINE-002\] Goroutine lancée sans mécanisme de synchronisation`
		go func(n int) {
			process(n)
		}(i)
	}
}

func GoodWithWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		process("data")
	}()
	wg.Wait()
}

func GoodWithContext(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			process("data")
		}
	}()
}

func GoodWithChannel() {
	done := make(chan bool)
	go func() {
		process("data")
		done <- true
	}()
	<-done
}

func process(v interface{}) {}
