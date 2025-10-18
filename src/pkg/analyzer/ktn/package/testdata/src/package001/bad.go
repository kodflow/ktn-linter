package package004

import (
	. "fmt" // want `\[KTN-PKG-001\].*`
	. "strings" // want `\[KTN-PKG-001\].*`
)

func BadDotImport() {
	// D'où viennent ces fonctions?
	Println("hello")
	ToUpper("test")
}
