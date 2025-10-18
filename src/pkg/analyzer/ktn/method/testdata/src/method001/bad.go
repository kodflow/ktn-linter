package method001

type Counter struct {
	value int
}

type User struct {
	Name string
	Age  int
}

func (c Counter) Increment() { // want `\[KTN-METHOD-001\] Méthode 'Increment' avec receiver non-pointeur mais modifie le receiver`
	c.value++ // Modifie la copie, pas l'original
}

func (u User) SetName(name string) { // want `\[KTN-METHOD-001\] Méthode 'SetName' avec receiver non-pointeur mais modifie le receiver`
	u.Name = name // Modifie la copie
}

func (u User) SetAge(age int) { // want `\[KTN-METHOD-001\] Méthode 'SetAge' avec receiver non-pointeur mais modifie le receiver`
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
