package pointer001

type User struct {
	Name string
}

// L'analyzer détecte uniquement certains patterns de déréférence nil
// TODO: Besoin d'amélioration pour détecter plus de cas

func TestPointer() {
	var p *int
	if p != nil {
		_ = *p
	}
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
