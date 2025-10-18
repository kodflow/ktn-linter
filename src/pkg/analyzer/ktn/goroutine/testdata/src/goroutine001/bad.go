package goroutine001

import "sync"

type Request struct {
	ID int
}

func BadGoroutineInLoop() {
	requests := []Request{{1}, {2}, {3}}
	for _, req := range requests {
		go handleRequest(req) // want `\[KTN-GOROUTINE-001\].*`
	}
}

func BadGoroutineInForLoop() {
	for i := 0; i < 100; i++ {
		go processItem(i) // want `\[KTN-GOROUTINE-001\].*`
	}
}

func GoodGoroutineWithWaitGroup() {
	var wg sync.WaitGroup
	requests := []Request{{1}, {2}, {3}}
	for _, req := range requests {
		wg.Add(1)
		go func(r Request) {
			defer wg.Done()
			handleRequest(r)
		}(req)
	}
	wg.Wait()
}

func GoodWorkerPool() {
	jobs := make(chan Request, 100)
	var wg sync.WaitGroup
	// Worker pool - nombre limité de goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}
	wg.Wait()
}

func handleRequest(req Request) {}
func processItem(i int)         {}
func worker(jobs chan Request, wg *sync.WaitGroup) {
	defer wg.Done()
}
