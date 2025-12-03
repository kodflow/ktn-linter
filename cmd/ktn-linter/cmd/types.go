// Types for the cmd package.
package cmd

// textEdit représente une modification de texte avec position.
type textEdit struct {
	start   int
	end     int
	newText []byte
}
