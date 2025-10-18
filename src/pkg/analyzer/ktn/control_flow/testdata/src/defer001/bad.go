package defer001

import "os"

func BadDeferInLoop() {
	files := []string{"a.txt", "b.txt", "c.txt"}
	for _, f := range files {
		file, _ := os.Open(f)
		// want `\[KTN-CONTROL-DEFER-001\] defer dans une boucle`
		defer file.Close()
	}
}

func BadDeferInForLoop() {
	for i := 0; i < 10; i++ {
		file, _ := os.Open("test.txt")
		// want `\[KTN-CONTROL-DEFER-001\] defer dans une boucle`
		defer file.Close()
	}
}

func GoodDeferOutsideLoop() {
	file, _ := os.Open("test.txt")
	defer file.Close()
}

func GoodDeferWithFunc() {
	files := []string{"a.txt", "b.txt"}
	for _, f := range files {
		func() {
			file, _ := os.Open(f)
			defer file.Close()
		}()
	}
}
