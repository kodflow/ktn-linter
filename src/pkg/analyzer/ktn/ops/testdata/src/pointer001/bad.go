package pointer001

type User struct {
	Name string
}

func BadNilDereference() {
	var u *User
	// want `\[KTN-OPS-POINTER-001\] Déréférence d'un pointeur potentiellement nil`
	_ = u.Name
}

func BadNilPointerAccess(u *User) {
	// want `\[KTN-OPS-POINTER-001\] Déréférence d'un pointeur potentiellement nil`
	name := u.Name
	_ = name
}

func GoodWithNilCheck(u *User) {
	if u != nil {
		_ = u.Name
	}
}

func GoodNonNilPointer() {
	u := &User{Name: "test"}
	_ = u.Name
}
