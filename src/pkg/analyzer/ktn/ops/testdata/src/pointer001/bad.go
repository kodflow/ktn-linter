package pointer001

type User struct {
	Name string
}

// L'analyzer détecte uniquement les déréférencements explicites avec *
// pas l'accès aux champs via pointeur

func BadNilDereference() {
	var u *int
	u = nil
	_ = *u // want `\[KTN-POINTER-001\] Déréférencement potentiel d'un pointeur nil`
}

func BadNilPointerAccess() {
	var p *int
	p = nil
	value := *p // want `\[KTN-POINTER-001\] Déréférencement potentiel d'un pointeur nil`
	_ = value
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
