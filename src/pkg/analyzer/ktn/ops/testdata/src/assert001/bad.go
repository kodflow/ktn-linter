package assert001

func BadTypeAssertionNoOk(i interface{}) {
	// want `\[KTN-OPS-ASSERT-001\] Type assertion sans vérification ok`
	s := i.(string)
	_ = s
}

func BadTypeAssertionMultiple(data interface{}) {
	// want `\[KTN-OPS-ASSERT-001\] Type assertion sans vérification ok`
	num := data.(int)
	_ = num
}

func GoodTypeAssertionWithOk(i interface{}) {
	s, ok := i.(string)
	if ok {
		_ = s
	}
}

func GoodTypeSwitch(i interface{}) {
	switch v := i.(type) {
	case string:
		_ = v
	case int:
		_ = v
	}
}
