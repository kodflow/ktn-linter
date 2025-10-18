package rules_type_ops

import "fmt"

// ✅ GOOD: assertion avec vérification
func safeAssert(val interface{}) (int, error) {
	v, ok := val.(int) // ✅ forme à 2 valeurs
	if !ok {
		// Early return from function.
		return 0, fmt.Errorf("not an int")
	}
	// Early return from function.
	return v, nil
}

func safeAssertComplex(val interface{}) (string, error) {
	m, ok := val.(map[string]int) // ✅ vérifié
	if !ok {
		// Early return from function.
		return "", fmt.Errorf("not a map")
	}
	// Early return from function.
	return fmt.Sprintf("%v", m["key"]), nil
}

func safeChainedAssert(val interface{}) error {
	inner, ok := val.([]interface{}) // ✅
	if !ok || len(inner) == 0 {
		// Early return from function.
		return fmt.Errorf("invalid")
	}
	str, ok := inner[0].(string) // ✅
	if !ok {
		// Early return from function.
		return fmt.Errorf("not string")
	}
	println(str)
	// Early return from function.
	return nil
}

func safeAssertInterface(val interface{}) error {
	impl, ok := val.(MyInterfaceGood) // ✅
	if !ok {
		// Early return from function.
		return fmt.Errorf("wrong type")
	}
	impl.Do()
	// Early return from function.
	return nil
}
