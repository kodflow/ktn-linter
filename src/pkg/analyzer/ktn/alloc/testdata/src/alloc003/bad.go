package alloc003

// User represents the struct.
type User struct {
	Name string
	Age  int
}
// Config represents the struct.

type Config struct {
	Host string
	Port int
}

func BadNewStruct() {
	u := new(User) // want `\[KTN-ALLOC-003\].*`
	_ = u
}

func BadNewStructConfig() {
	c := new(Config) // want `\[KTN-ALLOC-003\].*`
	_ = c
}

func GoodCompositeLiteral() {
	u := &User{}
	_ = u
}

func GoodCompositeLiteralWithFields() {
	u := &User{
		Name: "Alice",
		Age:  30,
	}
	_ = u
}
