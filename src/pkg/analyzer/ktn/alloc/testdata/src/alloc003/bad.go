package alloc003

type User struct {
	Name string
	Age  int
}

type Config struct {
	Host string
	Port int
}

func BadNewStruct() {
	u := new(User) // want `\[KTN-ALLOC-003\] Utilisez le composite literal &User\{\} au lieu de new\(User\)`
	_ = u
}

func BadNewStructConfig() {
	c := new(Config) // want `\[KTN-ALLOC-003\] Utilisez le composite literal &Config\{\} au lieu de new\(Config\)`
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
