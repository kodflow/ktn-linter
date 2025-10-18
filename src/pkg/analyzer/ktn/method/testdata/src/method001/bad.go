package method001

type Counter struct {
	value int
}

type User struct {
	Name string
	Age  int
}

// want `\[KTN-METHOD-001\] Méthode 'Increment' avec receiver non-pointeur mais modifie le receiver`
func (c Counter) Increment() {
	c.value++ // Modifie la copie, pas l'original
}

// want `\[KTN-METHOD-001\] Méthode 'SetName' avec receiver non-pointeur mais modifie le receiver`
func (u User) SetName(name string) {
	u.Name = name // Modifie la copie
}

// want `\[KTN-METHOD-001\] Méthode 'SetAge' avec receiver non-pointeur mais modifie le receiver`
func (u User) SetAge(age int) {
	u.Age = age // Modifie la copie
}

// Correct - receiver pointeur
func (c *Counter) IncrementCorrect() {
	c.value++
}

// Correct - ne modifie pas le receiver
func (c Counter) GetValue() int {
	return c.value
}

// Correct - receiver pointeur pour modification
func (u *User) UpdateUser(name string, age int) {
	u.Name = name
	u.Age = age
}
