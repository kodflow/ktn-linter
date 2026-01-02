package var027

import "fmt"

// goodRangeInt demonstrates the correct Go 1.22+ range over int.
func goodRangeInt(n int) {
	for i := range n {
		processGood(i)
	}
}

// goodRangeConst demonstrates range over constant.
func goodRangeConst() {
	for i := range 10 {
		fmt.Println(i)
	}
}

// goodModifiedStep is OK because step != 1.
func goodModifiedStep() {
	for i := 0; i < 10; i += 2 {
		fmt.Println(i)
	}
}

// goodStartNonZero is OK because start != 0.
func goodStartNonZero() {
	for i := 5; i < 10; i++ {
		fmt.Println(i)
	}
}

// goodDecrement is OK because it decrements.
func goodDecrement() {
	for i := 10; i > 0; i-- {
		fmt.Println(i)
	}
}

// goodLessOrEqual is OK because condition is <=, not <.
func goodLessOrEqual() {
	for i := 0; i <= 10; i++ {
		fmt.Println(i)
	}
}

// goodComplexCondition is OK because condition is more complex.
func goodComplexCondition(limit int) {
	for i := 0; i < limit && i < 100; i++ {
		fmt.Println(i)
	}
}

// goodInfiniteLoop is OK because no condition.
func goodInfiniteLoop() {
	count := 0
	for {
		count++
		// Break after first iteration
		if count > 0 {
			break
		}
	}
}

// goodRangeSlice demonstrates range over slice (not convertible).
func goodRangeSlice(items []int) {
	for i, v := range items {
		fmt.Println(i, v)
	}
}

// processGood is a helper function for testing.
func processGood(i int) {
	fmt.Println(i)
}
