package func010

// Good: 1 return value (unnamed OK)
func OneReturn() int {
	return 1
}

// Good: 2 return values (unnamed OK)
func TwoReturns() (int, error) {
	return 1, nil
}

// Good: 3 return values (unnamed OK - at limit)
func ThreeReturns() (int, string, error) {
	return 1, "test", nil
}

// Good: 4 return values with names
func FourNamedReturns() (count int, name string, valid bool, err error) {
	return 1, "test", true, nil
}

// Good: 5 return values with names
func FiveNamedReturns() (a int, b int, c string, d bool, e error) {
	return 1, 2, "test", true, nil
}

// Good: More than 3 returns but all named
func ManyNamedReturns() (id int, name string, age int, active bool, score float64) {
	return 1, "test", 30, true, 95.5
}
