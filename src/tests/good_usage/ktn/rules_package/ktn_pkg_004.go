package rules_package

import (
	"fmt"         // ✅ import normal
	"strings"     // ✅
	str "strings" // ✅ alias OK
)

func useNormalImports() {
	fmt.Println("hello")    // ✅ clair: vient de fmt
	strings.ToUpper("test") // ✅ clair: vient de strings
}

// ✅ GOOD: alias OK
func useAlias() {
	_ = str.ToLower("TEST") // ✅ alias acceptable
}
