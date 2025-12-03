// Entry point for the ktn-linter CLI tool.
package main

import (
	"github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd"
)

// Version est injectée au moment du build via -ldflags
var Version string = "dev"

// main est le point d'entrée du linter KTN.
// Returns: aucun
//
// Params: aucun
func main() {
	cmd.SetVersion(Version)
	cmd.Execute()
}
