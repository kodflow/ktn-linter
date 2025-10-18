package package004

import (
	. "fmt" // want `\[KTN-PKG-001\] Dot import de "fmt"`
	. "strings" // want `\[KTN-PKG-001\] Dot import de "strings"`
)

func BadDotImport() {
	// D'où viennent ces fonctions?
	Println("hello")
	ToUpper("test")
}
