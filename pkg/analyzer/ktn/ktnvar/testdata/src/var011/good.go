// Package var007 provides good test cases.
package var007

import "strings"

const (
	// AvgItemLength is the average item length estimate
	AvgItemLength int = 10
	// LoopCount is the number of iterations
	LoopCount int = 10
)

// init demonstrates proper string concatenation
func init() {
	items := []string{"a", "b", "c"}

	// Good: using strings.Builder
	var sb strings.Builder
	// Iteration over items to build string
	for _, item := range items {
		sb.WriteString(item)
	}
	_ = sb.String()

	// Good: using strings.Builder with preallocated size
	var sb2 strings.Builder
	sb2.Grow(len(items) * AvgItemLength)
	// Iteration over items to build string with capacity
	for _, item := range items {
		sb2.WriteString(item)
	}
	_ = sb2.String()

	// Good: using strings.Join for simple cases
	joined := strings.Join(items, "")
	_ = joined

	// Good: single concatenation, not in a loop
	result := "a" + "b"
	_ = result

	// Good: not string concatenation
	sum := 0
	// Iteration to compute sum
	for i := range LoopCount {
		sum += i
	}
	_ = sum
}
