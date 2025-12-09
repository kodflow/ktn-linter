// Bad examples for the var011 test case.
package var010

// Bad: Creating []byte buffers repeatedly in loops without sync.Pool

const (
	// LoopIterations est le nombre d'itérations de boucle
	LoopIterations int = 100
	// BufferSize est la taille du buffer
	BufferSize int = 1024
	// SmallBufferSize est la taille du petit buffer
	SmallBufferSize int = 512
	// TinyBufferSize est la taille du très petit buffer
	TinyBufferSize int = 128
	// MediumBufferSize est la taille moyenne du buffer
	MediumBufferSize int = 256
	// ConstBufferSize est la taille constante du buffer
	ConstBufferSize int = 2048
	// OuterIterations est le nombre d'itérations externes
	OuterIterations int = 10
	// InnerIterations est le nombre d'itérations internes
	InnerIterations int = 10
	// NestedBufferSize est la taille du buffer imbriqué
	NestedBufferSize int = 64
	// LargeBufferSize est la taille du grand buffer
	LargeBufferSize int = 512
	// LargeBufferCapacity est la capacité du grand buffer
	LargeBufferCapacity int = 1024
	// ItemsCount est le nombre d'éléments
	ItemsCount int = 3
	// CounterLimit est la limite du compteur
	CounterLimit int = 100
	// RangeIterations est le nombre d'itérations de range
	RangeIterations int = 50
)

// badProcessInLoop creates buffer in loop without sync.Pool
func badProcessInLoop() {
	items := make([]int, 0, LoopIterations)
	// Loop processes items
	for range LoopIterations {
		// Buffer created repeatedly in loop
		buffer := make([]byte, 0, BufferSize)
		items = append(items, len(buffer))
	}
}

// badRangeLoop creates buffer in range loop without sync.Pool
func badRangeLoop() {
	items := [ItemsCount]string{"a", "b", "c"}
	results := make([]int, 0, ItemsCount)
	// Loop processes items
	for _, item := range items {
		// Buffer created in each iteration
		buf := make([]byte, 0, SmallBufferSize)
		_ = item
		results = append(results, len(buf))
	}
}

// badInfiniteLoop creates buffer in infinite loop
func badInfiniteLoop() {
	// Infinite loop processing
	for {
		// Buffer allocated every iteration
		b := make([]byte, 0, TinyBufferSize)
		_ = b
		break
	}
}

// badWhileStyle creates buffer in while-style loop
func badWhileStyle() {
	counter := 0
	results := make([]int, 0, CounterLimit)
	// While-style loop
	for counter < CounterLimit {
		// Buffer created each iteration
		data := make([]byte, 0, MediumBufferSize)
		results = append(results, len(data))
		counter++
	}
}

// badConstSizeBuffer creates buffer with const size in loop
func badConstSizeBuffer() {
	results := make([]int, 0, RangeIterations)
	// Loop processes items
	for range RangeIterations {
		// Buffer with const size
		buf := make([]byte, 0, ConstBufferSize)
		results = append(results, len(buf))
	}
}

// badNestedLoop creates buffer in nested loop
func badNestedLoop() {
	// Outer loop
	for range OuterIterations {
		// Inner loop
		for range InnerIterations {
			// Buffer allocated in nested loop
			temp := make([]byte, 0, NestedBufferSize)
			_ = temp
		}
	}
}

// badMakeWithCapacity creates buffer with capacity in loop
func badMakeWithCapacity() {
	results := make([]int, 0, LoopIterations)
	// Loop processes items
	for range LoopIterations {
		// Buffer with capacity allocated in loop
		buffer := make([]byte, 0, LargeBufferCapacity)
		results = append(results, len(buffer))
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badProcessInLoop
	badProcessInLoop()
	// Appel de badRangeLoop
	badRangeLoop()
	// Appel de badInfiniteLoop
	badInfiniteLoop()
	// Appel de badWhileStyle
	badWhileStyle()
	// Appel de badConstSizeBuffer
	badConstSizeBuffer()
	// Appel de badNestedLoop
	badNestedLoop()
	// Appel de badMakeWithCapacity
	badMakeWithCapacity()
}
