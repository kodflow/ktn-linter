package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AutoFixer corrige automatiquement les violations KTN courantes
type AutoFixer struct {
	fset         *token.FileSet
	fixedFiles   map[string]bool
	stats        map[string]int
	dryRun       bool
}

func NewAutoFixer(dryRun bool) *AutoFixer {
	return &AutoFixer{
		fset:       token.NewFileSet(),
		fixedFiles: make(map[string]bool),
		stats:      make(map[string]int),
		dryRun:     dryRun,
	}
}

// FixFile applique toutes les corrections sur un fichier
func (af *AutoFixer) FixFile(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	fixed := string(content)

	// 1. Corriger les commentaires de fonction manquants (KTN-FUNC-002)
	fixed = af.addBasicFuncComments(fixed, filepath)

	// 2. Ajouter commentaires individuels aux variables (KTN-VAR-003)
	fixed = af.addVarComments(fixed)

	// 3. Ajouter commentaires returns (KTN-FUNC-008)
	fixed = af.addReturnComments(fixed)

	// 4. Corriger les noms de test avec underscore (KTN-FUNC-001)
	fixed = af.fixTestNames(fixed)

	if fixed != string(content) {
		if !af.dryRun {
			if err := os.WriteFile(filepath, []byte(fixed), 0644); err != nil {
				return err
			}
		}
		af.fixedFiles[filepath] = true
		return nil
	}

	return nil
}

// addBasicFuncComments ajoute des commentaires godoc basiques aux fonctions
func (af *AutoFixer) addBasicFuncComments(content string, filepath string) string {
	lines := strings.Split(content, "\n")
	var result []string

	funcRegex := regexp.MustCompile(`^func\s+(\w+)\(`)
	methodRegex := regexp.MustCompile(`^func\s+\([^)]+\)\s+(\w+)\(`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Vérifier si c'est une fonction sans commentaire godoc
		if matches := funcRegex.FindStringSubmatch(line); matches != nil {
			funcName := matches[1]

			// Ignorer main, init et fonctions de test
			if funcName == "main" || funcName == "init" || strings.HasPrefix(funcName, "Test") || strings.HasPrefix(funcName, "Benchmark") {
				result = append(result, line)
				continue
			}

			// Vérifier s'il y a déjà un commentaire
			hasComment := false
			if i > 0 {
				prevLine := strings.TrimSpace(lines[i-1])
				if strings.HasPrefix(prevLine, "//") {
					hasComment = true
				}
			}

			if !hasComment {
				// Ajouter un commentaire godoc basique
				indent := strings.Repeat("\t", strings.Count(line, "\t"))
				comment := fmt.Sprintf("%s// %s ...", indent, funcName)
				result = append(result, comment)
				af.stats["FUNC-002"]++
			}
		} else if matches := methodRegex.FindStringSubmatch(line); matches != nil {
			methodName := matches[1]

			// Vérifier s'il y a déjà un commentaire
			hasComment := false
			if i > 0 {
				prevLine := strings.TrimSpace(lines[i-1])
				if strings.HasPrefix(prevLine, "//") {
					hasComment = true
				}
			}

			if !hasComment {
				// Ajouter un commentaire godoc basique
				indent := strings.Repeat("\t", strings.Count(line, "\t"))
				comment := fmt.Sprintf("%s// %s ...", indent, methodName)
				result = append(result, comment)
				af.stats["FUNC-002"]++
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// addVarComments ajoute des commentaires aux variables dans les blocs var
func (af *AutoFixer) addVarComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	inVarBlock := false
	varRegex := regexp.MustCompile(`^\s*(\w+)\s+.*=`)

	for i, line := range lines {
		if strings.Contains(line, "var (") {
			inVarBlock = true
			result = append(result, line)
			continue
		}

		if inVarBlock && strings.Contains(line, ")") {
			inVarBlock = false
			result = append(result, line)
			continue
		}

		if inVarBlock {
			matches := varRegex.FindStringSubmatch(line)
			if matches != nil {
				varName := matches[1]

				// Vérifier s'il y a déjà un commentaire
				hasComment := false
				if i > 0 && strings.Contains(lines[i-1], "//") {
					hasComment = true
				}

				if !hasComment {
					indent := strings.Repeat("\t", strings.Count(line, "\t"))
					comment := fmt.Sprintf("%s// %s ...", indent, varName)
					result = append(result, comment)
					af.stats["VAR-003"]++
				}
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// addReturnComments ajoute des commentaires aux returns nus
func (af *AutoFixer) addReturnComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	returnRegex := regexp.MustCompile(`^\s*return\s+`)

	for i, line := range lines {
		// Si c'est un return
		if returnRegex.MatchString(line) {
			// Vérifier s'il y a déjà un commentaire sur la ligne précédente
			hasComment := false
			if i > 0 {
				prevLine := strings.TrimSpace(lines[i-1])
				if strings.HasPrefix(prevLine, "//") {
					hasComment = true
				}
			}

			// Ajouter un commentaire si absent
			if !hasComment && !strings.Contains(line, "//") {
				indent := strings.Repeat("\t", strings.Count(line, "\t"))
				commentedLine := line + " // Return result"
				result = append(result, commentedLine)
				af.stats["FUNC-008"]++
				continue
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// fixTestNames corrige les noms de fonctions de test avec underscore
func (af *AutoFixer) fixTestNames(content string) string {
	// TestFoo_Bar -> TestFooBar
	testRegex := regexp.MustCompile(`func (Test\w+)_(\w+)`)
	fixed := testRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := testRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			newName := parts[1] + strings.Title(parts[2])
			af.stats["FUNC-001"]++
			return fmt.Sprintf("func %s", newName)
		}
		return match
	})
	return fixed
}

// PrintStats affiche les statistiques de correction
func (af *AutoFixer) PrintStats() {
	fmt.Println("\n=== Statistiques de correction ===")
	fmt.Printf("Fichiers modifiés: %d\n", len(af.fixedFiles))
	fmt.Println("\nCorrections par type:")
	total := 0
	for rule, count := range af.stats {
		fmt.Printf("  KTN-%s: %d corrections\n", rule, count)
		total += count
	}
	fmt.Printf("\nTotal: %d corrections\n", total)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run auto_fix.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	dryRun := false
	if len(os.Args) > 2 && os.Args[2] == "--dry-run" {
		dryRun = true
		fmt.Println("Mode DRY RUN - aucun fichier ne sera modifié")
	}

	fixer := NewAutoFixer(dryRun)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Traiter uniquement les fichiers .go (sauf vendor)
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor") {
			if err := fixer.FixFile(path); err != nil {
				fmt.Printf("Erreur sur %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		os.Exit(1)
	}

	fixer.PrintStats()
}
