// Package rules_pool_bad contient des violations de KTN-POOL-001.
package rules_pool_bad

import "sync"

// bufferPool manages buffer allocation and reuse.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// Early return from function.
		return make([]byte, 1024)
	},
}

// objectPool manages object allocation and reuse.
var objectPool = &sync.Pool{
	New: func() interface{} {
		// Early return from function.
		return &DataObject{}
	},
}

// DataObject est une struct de test.
type DataObject struct {
	ID   int
	Data []byte
}

// ❌ KTN-POOL-001 : pool.Get() sans defer Put()

// BadNoDefer obtient un buffer sans le retourner.
func BadNoDefer() {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	// Traitement
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}
	// ❌ Oubli de bufferPool.Put(buf)
}

// BadNoDeferObject obtient un objet sans le retourner.
func BadNoDeferObject() {
	obj := objectPool.Get().(*DataObject) // Viole KTN-POOL-001
	obj.ID = 42
	obj.Data = make([]byte, 100)
	// ❌ Oubli de objectPool.Put(obj)
}

// BadConditionalReturn retourne avant le Put().
func BadConditionalReturn(condition bool) {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	if condition {
		// ❌ Return sans Put()
		return
	}
	// Even with Put here, early return causes violation
	bufferPool.Put(buf)
}

// BadPanicBeforePut peut paniquer sans Put().
func BadPanicBeforePut(index int) {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	// ❌ Si panic ici, buf n'est jamais retourné
	_ = buf[index] // Peut paniquer si index out of bounds
	bufferPool.Put(buf)
}

// BadMultipleGets obtient plusieurs buffers sans les retourner.
func BadMultipleGets() {
	buf1 := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	buf2 := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	copy(buf1, buf2)
	// ❌ Oubli de retourner buf1 et buf2
}

// BadLoopGet obtient des buffers en boucle sans les retourner.
func BadLoopGet(n int) {
	for i := 0; i < n; i++ {
		buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
		buf[0] = byte(i)
		// ❌ Oubli de Put(), fuite à chaque itération
	}
}

// BadAssignmentChain assigne depuis pool sans defer.
func BadAssignmentChain() {
	var buf []byte
	buf = bufferPool.Get().([]byte) // Viole KTN-POOL-001
	_ = buf
	// ❌ Oubli de Put()
}

// BadNoDeferInGoroutine utilise pool dans goroutine sans defer.
func BadNoDeferInGoroutine() {
	go func() {
		buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
		buf[0] = 1
		// ❌ Oubli de Put()
	}()
}

// BadPutWithoutDefer fait Put() mais sans defer.
func BadPutWithoutDefer() {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001 (pas de defer)
	processBuffer(buf)
	bufferPool.Put(buf) // Put existe mais sans defer → risque si processBuffer panic
}

func processBuffer(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 256)
	}
}

// BadNestedGet obtient un buffer dans une fonction imbriquée.
func BadNestedGet() {
	helper := func() []byte {
		// Early return from function.
		return bufferPool.Get().([]byte) // Viole KTN-POOL-001
	}
	buf := helper()
	_ = buf
	// ❌ Oubli de Put()
}

// BadGetInCondition obtient buffer dans condition.
func BadGetInCondition(condition bool) {
	if condition {
		buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
		_ = buf
		// ❌ Oubli de Put()
	}
}

// BadGetInSwitch obtient buffer dans switch.
func BadGetInSwitch(choice int) {
	switch choice {
	case 1:
		buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
		_ = buf
		// ❌ Oubli de Put()
	case 2:
		buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
		_ = buf
		// ❌ Oubli de Put()
	}
}

// BadGetBeforeError obtient buffer avant erreur potentielle.
func BadGetBeforeError() error {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	_ = buf
	// ❌ Si erreur retournée, buf n'est jamais retourné
	return processWithError()
}

func processWithError() error {
	// Early return from function.
	return nil
}

// BadDeferWrongVariable défère Put sur mauvaise variable.
func BadDeferWrongVariable() {
	buf1 := bufferPool.Get().([]byte) // Viole KTN-POOL-001
	buf2 := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf2) // Défère seulement buf2
	_ = buf1
	// ❌ buf1 n'est jamais retourné
}

// BadPutOutsideDefer met Put() en dehors d'un defer.
func BadPutOutsideDefer() {
	buf := bufferPool.Get().([]byte) // Viole KTN-POOL-001 (Put sans defer)
	for i := 0; i < 10; i++ {
		buf[i] = byte(i)
	}
	bufferPool.Put(buf) // ❌ Pas un defer, risque si panic dans la boucle
}
