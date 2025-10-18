package package004

import (
	// want `\[KTN-PKG-004\] Dot import de "fmt"`
	. "fmt"
	// want `\[KTN-PKG-004\] Dot import de "strings"`
	. "strings"
)

func BadDotImport() {
	// D'où viennent ces fonctions?
	Println("hello")
	ToUpper("test")
}
