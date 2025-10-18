package defer001

import "os"

func BadDeferInLoop() {
	files := []string{"a.txt", "b.txt", "c.txt"}
	for _, f := range files {
		file, _ := os.Open(f)
		defer file.Close() // want `\[KTN-DEFER-001\].*`
	}
}

func BadDeferInForLoop() {
	for i := 0; i < 10; i++ {
		file, _ := os.Open("test.txt")
		defer file.Close() // want `\[KTN-DEFER-001\].*`
	}
}

func GoodDeferOutsideLoop() {
	file, _ := os.Open("test.txt")
	defer file.Close()
}

// TODO: L'analyseur ne détecte pas les fonctions anonymes
// func GoodDeferWithFunc() {
// 	files := []string{"a.txt", "b.txt"}
// 	for _, f := range files {
// 		func() {
// 			file, _ := os.Open(f)
// 			defer file.Close()  // Faux positif - dans func anonyme
// 		}()
// 	}
// }
