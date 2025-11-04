package return002

// badReturnNilSlice returns nil for a slice type.
func badReturnNilSlice() []string {
	return nil // want "KTN-RETURN-002"
}

// badReturnNilMap returns nil for a map type.
func badReturnNilMap() map[string]int {
	return nil // want "KTN-RETURN-002"
}

// badReturnNilSliceConditional returns nil conditionally.
func badReturnNilSliceConditional(x int) []int {
	if x > 0 {
		return nil // want "KTN-RETURN-002"
	}
	return []int{x}
}

// badReturnNilMapConditional returns nil conditionally.
func badReturnNilMapConditional(key string) map[string]string {
	if key == "" {
		return nil // want "KTN-RETURN-002"
	}
	return map[string]string{key: "value"}
}

// badMultipleReturnsWithNil has multiple return statements with nil.
func badMultipleReturnsWithNil(flag bool) []byte {
	if flag {
		return nil // want "KTN-RETURN-002"
	}
	return nil // want "KTN-RETURN-002"
}
