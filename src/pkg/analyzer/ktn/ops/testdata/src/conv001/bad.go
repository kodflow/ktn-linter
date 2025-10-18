package conv002

func BadRedundantConversion(x int) int {
	// want `\[KTN-OPS-CONV-002\] Conversion redondante`
	return int(x)
}

func BadRedundantStringConv(s string) string {
	// want `\[KTN-OPS-CONV-002\] Conversion redondante`
	return string(s)
}

func GoodNeededConversion(x int32) int {
	return int(x)
}

func GoodBytesToString(b []byte) string {
	return string(b)
}
