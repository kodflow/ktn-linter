// Package gospec_bad_concurrency montre des patterns de concurrence non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_concurrency

import (
	"fmt"
	"sync"
	"time"
)

// ❌ BAD PRACTICE: Starting goroutine without any synchronization
func BadNoSync() {
	go func() {
		fmt.Println("async work")
	}()
	// Pas de WaitGroup, channel, ou autre mécanisme de sync
	// Le programme peut se terminer avant l'exécution
}

// ❌ BAD PRACTICE: Not closing channels
func BadNoChannelClose() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		// Devrait fermer: close(ch)
	}()

	// Range va bloquer indéfiniment
	count := 0
	for v := range ch {
		fmt.Println(v)
		count++
		if count == 5 {
			break // Forcé de break manuellement
		}
	}
}

// ❌ BAD PRACTICE: Closing channel from receiver side
func BadReceiverClose() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		time.Sleep(time.Second) // Attente arbitraire
	}()

	for v := range ch {
		fmt.Println(v)
		if v == 4 {
			close(ch) // Receiver ne devrait pas closer
			break
		}
	}
}

// ❌ BAD PRACTICE: Not using context for cancellation
func BadNoContext() {
	done := make(chan bool)
	go worker(done)

	time.Sleep(5 * time.Second)
	done <- true
	// Devrait utiliser context.Context
}

// ❌ BAD PRACTICE: Using sleep for synchronization
func BadSleepSync() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("done")
	}()

	time.Sleep(200 * time.Millisecond) // Timing fragile
	fmt.Println("continuing")
}

// ❌ BAD PRACTICE: Creating goroutine in loop without limit
func BadUnlimitedGoroutines(items []int) {
	for _, item := range items {
		go process(item) // Peut créer des milliers de goroutines
	}
	// Devrait utiliser worker pool
}

// ❌ BAD PRACTICE: Not handling channel send blocking
func BadBlockingSend() {
	ch := make(chan int) // Unbuffered
	ch <- 1              // Bloque forever, pas de receiver
	fmt.Println("never reached")
}

// ❌ BAD PRACTICE: Mutex locked but not unlocked
func BadNoUnlock() {
	var mu sync.Mutex
	mu.Lock()
	// Oubli de mu.Unlock()
	fmt.Println("critical section")
	// Devrait utiliser defer mu.Unlock()
}

// ❌ BAD PRACTICE: Copying mutex value
func BadCopyMutex() {
	var mu sync.Mutex
	mu.Lock()
	mu2 := mu // Copie le mutex (incorrect)
	mu.Unlock()
	_ = mu2
}

// ❌ BAD PRACTICE: Channel closed multiple times
func BadDoubleClose() {
	ch := make(chan int)
	close(ch)
	// Plus tard...
	close(ch) // Panic: close of closed channel
}

// ❌ BAD PRACTICE: Sending on closed channel
func BadSendOnClosed() {
	ch := make(chan int)
	close(ch)
	ch <- 1 // Panic: send on closed channel
}

// ❌ BAD PRACTICE: WaitGroup Add called in wrong place
func BadWaitGroupTiming() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		go func(n int) {
			wg.Add(1) // Devrait être avant go
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait()
}

// ❌ BAD PRACTICE: Range over channel without checking if closed
func BadNoClosedCheck() {
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
		// Oublie de fermer
	}()

	// Bloquera indéfiniment après les 2 valeurs
	timeout := time.After(1 * time.Second)
	for {
		select {
		case v := <-ch:
			fmt.Println(v)
		case <-timeout:
			return
		}
	}
}

// ❌ BAD PRACTICE: Not using buffered channel when appropriate
func BadUnbufferedForKnownSize() {
	ch := make(chan int) // Unbuffered
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i // Bloque à chaque send
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
	// Devrait utiliser: make(chan int, 100) ou taille raisonnable
}

// ❌ BAD PRACTICE: Goroutine leak - no way to stop
var cache = make(map[string]string)
var cacheMu sync.RWMutex

func BadGoroutineLeak() {
	go func() {
		for { // Infinite loop sans mécanisme d'arrêt
			cacheMu.Lock()
			// Simulate cache cleanup
			cache = make(map[string]string)
			cacheMu.Unlock()
			time.Sleep(1 * time.Hour)
		}
	}()
	// Goroutine ne peut jamais être arrêtée
}

// ❌ BAD PRACTICE: Race condition - shared variable without protection
var counter int

func BadRaceCondition() {
	for i := 0; i < 10; i++ {
		go func() {
			counter++ // Race condition
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(counter)
}

// ❌ BAD PRACTICE: Defer in loop (performance)
func BadDeferInLoop() {
	var mu sync.Mutex
	for i := 0; i < 1000; i++ {
		mu.Lock()
		defer mu.Unlock() // Defer s'accumule, unlock seulement à la fin
		// Work...
	}
	// Devrait unlock manuellement dans loop
}

// ❌ BAD PRACTICE: Using sync.Map when regular map with mutex is clearer
func BadSyncMapOveruse() {
	// Pour des use cases simples, map + mutex est plus clair
	var sm sync.Map
	sm.Store("key", "value")
	val, _ := sm.Load("key")
	fmt.Println(val)
	// sync.Map utile seulement pour patterns spécifiques
}

// ❌ BAD PRACTICE: Not using select with timeout
func BadNoTimeout() {
	ch := make(chan int)
	result := <-ch // Peut bloquer forever
	fmt.Println(result)
	// Devrait utiliser select avec time.After
}

// Helper functions
func worker(done chan bool) {
	<-done
}

func process(int) {}
