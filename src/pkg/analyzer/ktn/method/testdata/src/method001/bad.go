package method001

// Counter represents the struct.
type Counter struct {
	value int
}
// User represents the struct.

type User struct {
	Name string
	Age  int
}

func (c Counter) Increment() { // want `\[KTN-METHOD-001\].*`
	c.value++ // Modifie la copie, pas l'original
}

func (u User) SetName(name string) { // want `\[KTN-METHOD-001\].*`
	u.Name = name // Modifie la copie
}

func (u User) SetAge(age int) { // want `\[KTN-METHOD-001\].*`
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
// Data represents the struct.

// Test IndexExpr avec slice
type Data struct {
	items []int
}

func (d Data) UpdateItem(index int, value int) {
	d.items[index] = value // OK: slices sont des références
// Cache represents the struct.
}

// Test IndexExpr avec map
type Cache struct {
	data map[string]int
}

func (c Cache) Set(key string, value int) {
	c.data[key] = value // OK: maps sont des références
}

// Correct - receiver pointeur avec IndexExpr
func (d *Data) UpdateItemCorrect(index int, value int) {
	d.items[index] = value
}

func (c *Cache) SetCorrect(key string, value int) {
// Value represents the struct.
	c.data[key] = value
}

// Test assignation directe du receiver
type Value struct {
	n int
}

func (v Value) Reset() { // want `\[KTN-METHOD-001\].*`
	v = Value{n: 0} // Assigne une nouvelle valeur au receiver
}

// Correct - receiver pointeur pour réassignation
func (v *Value) ResetCorrect() {
	*v = Value{n: 0}
}
