// Package rules_alloc_bad contient des violations KTN-ALLOC-004.
package rules_alloc_bad

// Viole KTN-ALLOC-004 : new(struct) au lieu de &struct{}

// User est une struct de test.
type User struct {
	ID   int
	Name string
	Age  int
}

// Config est une struct de configuration.
type Config struct {
	Host    string
	Port    int
	Timeout int
}

// BadNewUser crée un User avec new().
func BadNewUser() *User {
	u := new(User) // Viole KTN-ALLOC-004
	u.ID = 1
	u.Name = "Alice"
	// Retourne le pointeur vers User
	return u
}

// BadNewConfig crée une Config avec new().
func BadNewConfig() *Config {
	cfg := new(Config) // Viole KTN-ALLOC-004
	cfg.Host = "localhost"
	cfg.Port = 8080
	// Retourne le pointeur vers Config
	return cfg
}

// BadNewStructInline crée une struct inline avec new().
func BadNewStructInline() {
	type Point struct {
		X, Y int
	}
	p := new(Point) // Viole KTN-ALLOC-004
	p.X = 10
	p.Y = 20
}

// BadNewMultipleStructs crée plusieurs structs avec new().
func BadNewMultipleStructs() {
	u1 := new(User)   // Viole KTN-ALLOC-004
	u2 := new(User)   // Viole KTN-ALLOC-004
	cfg := new(Config) // Viole KTN-ALLOC-004
	_, _, _ = u1, u2, cfg
}

// BadNewNestedStruct crée une struct imbriquée avec new().
func BadNewNestedStruct() {
	type Address struct {
		Street string
		City   string
	}
	type Person struct {
		Name    string
		Address *Address
	}

	addr := new(Address) // Viole KTN-ALLOC-004
	addr.Street = "123 Main St"
	addr.City = "Paris"

	person := new(Person) // Viole KTN-ALLOC-004
	person.Name = "Bob"
	person.Address = addr
}

// Cas farfelus

// BadNewEmptyStruct crée une struct vide avec new().
func BadNewEmptyStruct() {
	type Empty struct{}
	e := new(Empty) // Viole KTN-ALLOC-004
	_ = e
}

// BadNewStructWithMethods crée une struct avec méthodes via new().
func BadNewStructWithMethods() {
	type Counter struct {
		count int
	}
	c := new(Counter) // Viole KTN-ALLOC-004
	c.count = 10
}

// BadNewStructAsReturn retourne un new(struct).
func BadNewStructAsReturn() *User {
	// Retourne un pointeur créé avec new() (non-idiomatique)
	return new(User) // Viole KTN-ALLOC-004
}

// BadNewStructInLoop crée des structs avec new() dans une boucle.
func BadNewStructInLoop() {
	var users []*User
	for i := 0; i < 5; i++ {
		u := new(User) // Viole KTN-ALLOC-004
		u.ID = i
		users = append(users, u)
	}
}

// BadNewStructComplex crée une struct complexe avec new().
func BadNewStructComplex() {
	type Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Schema   string
		Pool     int
	}
	db := new(Database) // Viole KTN-ALLOC-004
	db.Host = "localhost"
	db.Port = 5432
}

// BadNewStructPointerField crée une struct avec champ pointeur via new().
func BadNewStructPointerField() {
	type Node struct {
		Value int
		Next  *Node
	}
	node := new(Node) // Viole KTN-ALLOC-004
	node.Value = 42
	node.Next = nil
}

// BadNewGenericStruct crée une struct générique avec new().
func BadNewGenericStruct() {
	type Container struct {
		Data interface{}
	}
	cont := new(Container) // Viole KTN-ALLOC-004
	cont.Data = "test"
}
