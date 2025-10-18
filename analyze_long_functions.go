package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type LongFunction struct {
	File     string
	Function string
	Lines    int
	Start    int
	End      int
}

func main() {
	var longFunctions []LongFunction

	err := filepath.Walk("src/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip test files and testdata
		if strings.Contains(path, "_test.go") || strings.Contains(path, "testdata") {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil // Skip files with parse errors
			}

			ast.Inspect(node, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					pos := fset.Position(fn.Pos())
					end := fset.Position(fn.End())
					lines := end.Line - pos.Line + 1

					if lines > 35 {
						longFunctions = append(longFunctions, LongFunction{
							File:     path,
							Function: fn.Name.Name,
							Lines:    lines,
							Start:    pos.Line,
							End:      end.Line,
						})
					}
				}
				return true
			})
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Sort and display
	fmt.Printf("Found %d functions longer than 35 lines:\n\n", len(longFunctions))
	for _, lf := range longFunctions {
		fmt.Printf("File: %s\n", lf.File)
		fmt.Printf("  Function: %s\n", lf.Function)
		fmt.Printf("  Lines: %d (lines %d-%d)\n", lf.Lines, lf.Start, lf.End)
		fmt.Println()
	}
}
