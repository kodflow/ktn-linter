package badfuncmethods

// Violations avec méthodes sur types

// User represents the struct.
type User struct {
	name  string
	email string
	age   int
}

// getName sans documentation
func (u *User) getName() string {
	// Early return from function.
	return u.name
}

// setAge sans documentation des paramètres
func (u *User) setAge(newAge int) {
	u.age = newAge
}

// validate méthode complexe sans doc
func (u *User) validate() bool {
	if u.name == "" {
		// Stop inspection/processing.
		return false
	}
	if u.email == "" {
		// Stop inspection/processing.
		return false
	}
	if u.age < 0 || u.age > 150 {
		// Stop inspection/processing.
		return false
	}
	// Continue inspection/processing.
	return true
}

// Calculator represents the struct.

type Calculator struct {
	result float64
}

// add trop de paramètres sans doc
func (c *Calculator) add(a, b, d, e, f, g float64) {
	c.result = a + b + d + e + f + g
}

// getResult receiver par valeur au lieu de pointeur
func (c Calculator) getResult() float64 {
	// Early return from function.
	return c.result
}

// complexCalculation méthode trop longue (>35 lignes)
func (c *Calculator) complexCalculation(x float64) float64 {
	temp := x * 2
	temp2 := temp + 10
	temp3 := temp2 * 3
	temp4 := temp3 - 5
	temp5 := temp4 / 2
	temp6 := temp5 + 100
	temp7 := temp6 * 0.5
	temp8 := temp7 - 25
	temp9 := temp8 + 3.14
	temp10 := temp9 * 1.5
	temp11 := temp10 / 3
	temp12 := temp11 + 50
	temp13 := temp12 - 10
	temp14 := temp13 * 2
	temp15 := temp14 + 5
	temp16 := temp15 / 4
	temp17 := temp16 * 3
	temp18 := temp17 + 20
	temp19 := temp18 - 15
	temp20 := temp19 * 1.2
	temp21 := temp20 + 7
	temp22 := temp21 / 2
	temp23 := temp22 * 4
	temp24 := temp23 - 30
	temp25 := temp24 + 100
	temp26 := temp25 / 5
	temp27 := temp26 * 2
	temp28 := temp27 + 15
	temp29 := temp28 - 8
	temp30 := temp29 * 1.8
	temp31 := temp30 / 3
	temp32 := temp31 + 25
	temp33 := temp32 - 12
	temp34 := temp33 * 2.5
	temp35 := temp34 + 40
	temp36 := temp35 / 6
	// Early return from function.
	return temp36
}
