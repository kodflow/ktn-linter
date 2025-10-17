// Package gospec_statements démontre les statements selon la spec Go.
// Référence: https://go.dev/ref/spec#Statements
package gospec_statements

import "fmt"

// Spec Go: If statements
// https://go.dev/ref/spec#If_statements

func ValidIfStatement(x int) {
	if x > 0 {
		fmt.Println("positive")
	}
}

func ValidIfElse(x int) {
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	}
}

func ValidIfElseIf(x int) {
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}
}

func ValidIfWithInit(x int) {
	if y := x * 2; y > 10 {
		fmt.Println("large")
	}
}

// Spec Go: For statements
// https://go.dev/ref/spec#For_statements

func ValidInfiniteLoop() {
	for {
		break
	}
}

func ValidForCondition(n int) {
	for n > 0 {
		n--
	}
}

func ValidForClause() {
	for i := 0; i < 10; i++ {
		_ = i
	}
}

func ValidForRange() {
	slice := []int{1, 2, 3}
	for i, v := range slice {
		_, _ = i, v
	}
}

func ValidForRangeMap() {
	m := map[string]int{"a": 1}
	for k, v := range m {
		_, _ = k, v
	}
}

// Spec Go: Switch statements
// https://go.dev/ref/spec#Switch_statements

func ValidSwitch(x int) {
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}
}

func ValidSwitchWithInit(x int) {
	switch y := x * 2; y {
	case 2:
		fmt.Println("two")
	case 4:
		fmt.Println("four")
	}
}

func ValidSwitchNoExpression() {
	x := 10
	switch {
	case x > 0:
		fmt.Println("positive")
	case x < 0:
		fmt.Println("negative")
	default:
		fmt.Println("zero")
	}
}

// Spec Go: Type switch
// https://go.dev/ref/spec#Type_switches

func ValidTypeSwitch(v interface{}) {
	switch v.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	default:
		fmt.Println("other")
	}
}

func ValidTypeSwitchWithAssignment(v interface{}) {
	switch x := v.(type) {
	case int:
		_ = x + 1
	case string:
		_ = x + "!"
	}
}

// Spec Go: Select statements
// https://go.dev/ref/spec#Select_statements

func ValidSelect() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	select {
	case v := <-ch1:
		_ = v
	case v := <-ch2:
		_ = v
	}
}

func ValidSelectWithDefault() {
	ch := make(chan int)

	select {
	case v := <-ch:
		_ = v
	default:
		fmt.Println("no value")
	}
}

// Spec Go: Return statements
// https://go.dev/ref/spec#Return_statements

func ValidReturnVoid() {
	return
}

func ValidReturnValue() int {
	return 42
}

func ValidReturnMultiple() (int, string) {
	return 42, "hello"
}

func ValidReturnNamed() (x int, y string) {
	x = 42
	y = "hello"
	return
}

// Spec Go: Break and Continue
// https://go.dev/ref/spec#Break_statements
// https://go.dev/ref/spec#Continue_statements

func ValidBreakContinue() {
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		if i%2 == 0 {
			continue
		}
	}
}

// Spec Go: Goto statements
// https://go.dev/ref/spec#Goto_statements

func ValidGoto() {
	goto End
	fmt.Println("skipped")
End:
	fmt.Println("end")
}

// Spec Go: Fallthrough statements
// https://go.dev/ref/spec#Fallthrough_statements

func ValidFallthrough(x int) {
	switch x {
	case 1:
		fmt.Println("one")
		fallthrough
	case 2:
		fmt.Println("one or two")
	}
}

// Spec Go: Defer statements
// https://go.dev/ref/spec#Defer_statements

func ValidDefer() {
	defer fmt.Println("deferred")
	fmt.Println("normal")
}

func ValidMultipleDefers() {
	defer fmt.Println("first")
	defer fmt.Println("second")
	defer fmt.Println("third")
}

// Spec Go: Go statements
// https://go.dev/ref/spec#Go_statements

func ValidGoroutine() {
	go func() {
		fmt.Println("goroutine")
	}()
}

func ValidGoroutineWithParams() {
	go fmt.Println("hello")
}

// Spec Go: Assignment statements
// https://go.dev/ref/spec#Assignment_statements

func ValidAssignments() {
	var x int
	x = 42

	var y, z int
	y, z = 1, 2

	x += 10
	x -= 5
	x *= 2
	x /= 2

	_, _, _ = x, y, z
}

// Spec Go: Inc and Dec statements
// https://go.dev/ref/spec#IncDec_statements

func ValidIncDec() {
	x := 0
	x++
	x--
	_ = x
}

// Spec Go: Labeled statements
// https://go.dev/ref/spec#Labeled_statements

func ValidLabeled() {
Loop:
	for i := 0; i < 10; i++ {
		if i == 5 {
			break Loop
		}
	}
}
