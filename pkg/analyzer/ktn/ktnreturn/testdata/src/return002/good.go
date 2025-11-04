package return002

// goodReturnEmptySlice returns an empty slice instead of nil.
func goodReturnEmptySlice() []string {
	return []string{}
}

// goodReturnEmptyMap returns an empty map instead of nil.
func goodReturnEmptyMap() map[string]int {
	return map[string]int{}
}

// goodReturnEmptySliceConditional returns empty slice conditionally.
func goodReturnEmptySliceConditional(x int) []int {
	if x > 0 {
		return []int{}
	}
	return []int{x}
}

// goodReturnPointerNil can return nil for pointer types (allowed).
func goodReturnPointerNil() *string {
	return nil
}

// goodReturnInterfaceNil can return nil for interface types (allowed).
func goodReturnInterfaceNil() error {
	return nil
}

// goodReturnError returns nil for error type (allowed).
func goodReturnError() error {
	return nil
}

// goodReturnMakeSlice returns a slice created with make.
func goodReturnMakeSlice(size int) []int {
	return make([]int, size)
}

// goodReturnMakeMap returns a map created with make.
func goodReturnMakeMap() map[string]bool {
	return make(map[string]bool)
}
