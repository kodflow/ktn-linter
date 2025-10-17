// Package gospec_good_concurrency démontre les patterns de concurrence idiomatiques.
// Référence: https://go.dev/doc/effective_go#concurrency
// Référence: https://go.dev/blog/pipelines
package gospec_good_concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ✅ GOOD: Using WaitGroup for goroutine synchronization
func GoodWaitGroup() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1) // Add before spawning goroutine
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait()
}

// ✅ GOOD: Sender closes channel
func GoodChannelClose() {
	ch := make(chan int)

	go func() {
		defer close(ch) // Sender closes
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// ✅ GOOD: Using context for cancellation
func GoodContext(ctx context.Context) {
	go worker(ctx)
	// Caller can cancel via context
}

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("cancelled")
			return
		default:
			// Do work
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ✅ GOOD: Worker pool pattern
func GoodWorkerPool(jobs []int) {
	const numWorkers = 3
	jobsCh := make(chan int, len(jobs))
	resultsCh := make(chan int, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				resultsCh <- processJob(job)
			}
		}()
	}

	// Send jobs
	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh)

	// Close results after all workers done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results
	for result := range resultsCh {
		fmt.Println(result)
	}
}

// ✅ GOOD: Using buffered channel when size known
func GoodBufferedChannel() {
	ch := make(chan int, 100)

	go func() {
		defer close(ch)
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// ✅ GOOD: Protecting shared state with mutex
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ✅ GOOD: Using RWMutex for read-heavy workloads
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func (m *SafeMap) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// ✅ GOOD: Using sync.Once for initialization
type Service struct {
	once     sync.Once
	instance *Connection
}

func (s *Service) Init() *Connection {
	s.once.Do(func() {
		s.instance = &Connection{}
	})
	return s.instance
}

// ✅ GOOD: Channel for signaling
func GoodSignaling() {
	done := make(chan struct{})

	go func() {
		defer close(done)
		time.Sleep(1 * time.Second)
		fmt.Println("work done")
	}()

	<-done
	fmt.Println("received signal")
}

// ✅ GOOD: Select with timeout pattern
func GoodTimeout() error {
	ch := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- 42
	}()

	select {
	case result := <-ch:
		fmt.Println("got result:", result)
		return nil
	case <-time.After(1 * time.Second):
		return fmt.Errorf("timeout")
	}
}

// ✅ GOOD: Goroutine with proper cleanup
func GoodCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("tick")
		}
	}
}

// ✅ GOOD: Fan-out, fan-in pattern
func GoodFanOutFanIn(inputs []int) []int {
	// Fan-out
	ch1 := processAsync(inputs[:len(inputs)/2])
	ch2 := processAsync(inputs[len(inputs)/2:])

	// Fan-in
	results := make([]int, 0, len(inputs))
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			results = append(results, v)
		case v := <-ch2:
			results = append(results, v)
		}
	}

	return results
}

func processAsync(inputs []int) chan int {
	ch := make(chan int, len(inputs))
	go func() {
		defer close(ch)
		for _, v := range inputs {
			ch <- v * 2
		}
	}()
	return ch
}

// ✅ GOOD: Pipeline pattern
func GoodPipeline() {
	// Stage 1: Generate numbers
	generate := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; i < 10; i++ {
				ch <- i
			}
		}()
		return ch
	}

	// Stage 2: Square numbers
	square := func(in <-chan int) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for v := range in {
				ch <- v * v
			}
		}()
		return ch
	}

	// Pipeline execution
	for v := range square(generate()) {
		fmt.Println(v)
	}
}

// ✅ GOOD: Error handling with errgroup
func GoodErrorGroup() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 3)

	tasks := []func() error{task1, task2, task3}

	for _, task := range tasks {
		wg.Add(1)
		go func(t func() error) {
			defer wg.Done()
			if err := t(); err != nil {
				errCh <- err
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Return first error if any
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

// ✅ GOOD: Non-blocking send with select
func GoodNonBlockingSend(ch chan int, value int) bool {
	select {
	case ch <- value:
		return true
	default:
		return false
	}
}

// ✅ GOOD: Non-blocking receive with select
func GoodNonBlockingReceive(ch chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

// ✅ GOOD: Using context with timeout
func GoodContextTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return doWork(ctx)
}

func doWork(ctx context.Context) error {
	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Helper types and functions
type Connection struct{}

func processJob(n int) int { return n * 2 }
func task1() error         { return nil }
func task2() error         { return nil }
func task3() error         { return nil }
