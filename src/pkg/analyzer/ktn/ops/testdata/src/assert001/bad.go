package assert001

func BadTypeAssertionNoOk(i interface{}) {
	s := i.(string) // want `\[KTN-ASSERT-001\] Assertion de type sans vérification`
	_ = s
}

func BadTypeAssertionMultiple(data interface{}) {
	num := data.(int) // want `\[KTN-ASSERT-001\] Assertion de type sans vérification`
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
