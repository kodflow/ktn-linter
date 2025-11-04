package comment001

// badRedundantReturn has obvious comment.
func badRedundantReturn() int {
	// Return 42 // want "KTN-COMMENT-001"
	return 42
}

// badRedundantAssignment states the obvious.
func badRedundantAssignment() {
	// Set x to 10 // want "KTN-COMMENT-001"
	x := 10
	_ = x
}

// badRedundantNil comments on nil return.
func badRedundantNil() error {
	// Return nil // want "KTN-COMMENT-001"
	return nil
}

// badRedundantIncrement comments obvious operation.
func badRedundantIncrement() {
	i := 0
	// Increment i // want "KTN-COMMENT-001"
	i++
	_ = i
}

// badRedundantLoop comments obvious loop.
func badRedundantLoop() {
	// Loop from 0 to 10 // want "KTN-COMMENT-001"
	for i := 0; i < 10; i++ {
		_ = i
	}
}
