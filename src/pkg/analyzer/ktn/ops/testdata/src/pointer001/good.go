package pointer001

type Data struct {
	Value int
}

// Cas corrects - déréférencement sécurisé

func GoodNonNilPointer() {
	var p *int
	x := 5
	p = &x
	_ = *p
}

func GoodNilAssignedButChecked() {
	var p *int
	p = nil
	if p != nil {
		_ = *p // OK - vérifié
	}
}

func GoodStructPointer() {
	d := &Data{Value: 42}
	_ = d.Value
}

func GoodCheckedBeforeDeref(p *int) {
	if p != nil {
		v := *p
		_ = v
	}
}
