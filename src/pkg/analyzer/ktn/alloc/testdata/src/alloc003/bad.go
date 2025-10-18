package alloc004

type User struct {
	Name string
	Age  int
}

type Config struct {
	Host string
	Port int
}

func BadNewStruct() {
	// want `\[KTN-ALLOC-004\] Utilisez le composite literal &User\{\} au lieu de new\(User\)`
	u := new(User)
	_ = u
}

func BadNewStructConfig() {
	// want `\[KTN-ALLOC-004\] Utilisez le composite literal &Config\{\} au lieu de new\(Config\)`
	c := new(Config)
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
