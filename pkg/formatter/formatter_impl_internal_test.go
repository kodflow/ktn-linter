package formatter

import (
	"bytes"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func createTestDiagnostics() []analysis.Diagnostic {
	return []analysis.Diagnostic{
		{
			Pos:      token.Pos(1),
			Message:  "KTN-VAR-001: Variable naming issue\nThis is a test diagnostic.\nExample: var myVar int",
			Category: "naming",
		},
		{
			Pos:      token.Pos(10),
			Message:  "KTN-FUNC-002: Function complexity too high\nSplit into smaller functions",
			Category: "complexity",
		},
	}
}

// TestFormatEmpty tests the functionality of the corresponding implementation.
func TestFormatEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, false, false)
	fset := token.NewFileSet()

	formatter.Format(fset, []analysis.Diagnostic{})

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

// TestFormatHumanMode tests the functionality of the corresponding implementation.
func TestFormatHumanMode(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, false, false)
	fset := token.NewFileSet()

	// Add a test file
	fset.AddFile("test.go", 1, 1000)
	diagnostics := createTestDiagnostics()

	formatter.Format(fset, diagnostics)

	output := buf.String()
	if !strings.Contains(output, "KTN-LINTER REPORT") {
		t.Error("Expected human-readable report header")
	}
	if !strings.Contains(output, "test.go") {
		t.Error("Expected filename in output")
	}
}

// TestFormatAIMode tests the functionality of the corresponding implementation.
func TestFormatAIMode(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, true, false, false)
	fset := token.NewFileSet()

	fset.AddFile("test.go", 1, 1000)
	diagnostics := createTestDiagnostics()

	formatter.Format(fset, diagnostics)

	output := buf.String()
	if !strings.Contains(output, "# KTN-Linter Report (AI Mode)") {
		t.Error("Expected AI mode header")
	}
	if !strings.Contains(output, "## File:") {
		t.Error("Expected AI format file markers")
	}
}

// TestFormatSimpleMode tests the functionality of the corresponding implementation.
func TestFormatSimpleMode(t *testing.T) {
	const MIN_LINE_COUNT int = 2

	tests := []struct {
		name  string
		check func(t *testing.T, output string, lines []string)
	}{
		{
			name: "has minimum line count",
			check: func(t *testing.T, output string, lines []string) {
				// Simple mode: one line per diagnostic
				if len(lines) < MIN_LINE_COUNT {
					t.Errorf("Expected at least %d lines in simple mode, got %d", MIN_LINE_COUNT, len(lines))
				}
			},
		},
		{
			name: "has filename:line:col format",
			check: func(t *testing.T, output string, lines []string) {
				// Check format: filename:line:col: [CODE] message
				if !strings.Contains(lines[0], "test.go:") {
					t.Error("Expected filename:line:col format")
				}
			},
		},
		{
			name: "has error code in brackets",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification du code d'erreur
				if !strings.Contains(lines[0], "[KTN-") {
					t.Errorf("Expected error code in brackets, got: %s", lines[0])
				}
			},
		},
	}

	// Préparation commune
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, false, true)
	fset := token.NewFileSet()
	fset.AddFile("test.go", 1, 1000)
	diagnostics := createTestDiagnostics()
	formatter.Format(fset, diagnostics)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, output, lines)
		})
	}
}

// TestFormatNoColor tests the functionality of the corresponding implementation.
func TestFormatNoColor(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, true, false)
	fset := token.NewFileSet()

	fset.AddFile("test.go", 1, 1000)
	diagnostics := createTestDiagnostics()

	formatter.Format(fset, diagnostics)

	output := buf.String()

	// No ANSI color codes
	if strings.Contains(output, "\033[") {
		t.Error("Expected no ANSI color codes with noColor=true")
	}
}

// TestExtractCode tests the functionality of the corresponding implementation.
func TestExtractCode(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			"valid code",
			"KTN-VAR-001: Variable issue",
			"KTN-VAR-001",
		},
		{
			"no code",
			"Plain message without code",
			"UNKNOWN",
		},
		{
			"incomplete code",
			"[KTN-VAR- incomplete",
			"UNKNOWN",
		},
		{
			"code in middle",
			"Some text [KTN-FUNC-002] and more text",
			"KTN-FUNC-002",
		},
		{
			"prefix format with colon",
			"KTN-TEST-005: some error message",
			"KTN-TEST-005",
		},
		{
			"prefix format without colon",
			"KTN-TEST-006 missing colon",
			"UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCode(tt.message)
			if got != tt.expected {
				t.Errorf("extractCode(%q) = %q, want %q", tt.message, got, tt.expected)
			}
		})
	}
}

// TestExtractMessage tests the functionality of the corresponding implementation.
func TestExtractMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			"with code and newline",
			"KTN-VAR-001: Variable issue\nDetails here",
			"Variable issue",
		},
		{
			"with code only",
			"KTN-VAR-001: Variable issue",
			"Variable issue",
		},
		{
			"no code",
			"Plain message\nWith newline",
			"Plain message",
		},
		{
			"empty",
			"",
			"",
		},
		{
			"bracket at end",
			"KTN-TEST-001:",
			"KTN-TEST-001:",
		},
		{
			"prefix format with colon",
			"KTN-FUNC-002: Function issue here",
			"Function issue here",
		},
		{
			"prefix format with colon at end",
			"KTN-FUNC-003:",
			"KTN-FUNC-003:",
		},
		{
			"prefix format no colon",
			"KTN-FUNC-004 no colon",
			"KTN-FUNC-004 no colon",
		},
		{
			"only newline",
			"\nJust newline",
			"",
		},
		{
			"bracket exactly at end",
			"message]",
			"message]",
		},
		{
			"colon exactly at end",
			"KTN-TEST:",
			"KTN-TEST:",
		},
		{
			"bracket format basic",
			"[KTN-VAR-001] issue here",
			"issue here",
		},
		{
			"KTN prefix without any separator",
			"KTN-ERROR",
			"KTN-ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMessage(tt.message)
			if got != tt.expected {
				t.Errorf("extractMessage(%q) = %q, want %q", tt.message, got, tt.expected)
			}
		})
	}
}

// TestGetCodeColor tests the functionality of the corresponding implementation.
func TestGetCodeColor(t *testing.T) {
	formatter := &formatterImpl{noColor: false}

	tests := []struct {
		code     string
		expected string
	}{
		{"KTN-VAR-001", YELLOW},   // Toutes les règles = WARNING
		{"KTN-FUNC-002", YELLOW},  // Toutes les règles = WARNING
		{"KTN-TEST-003", YELLOW},  // Toutes les règles = WARNING
		{"KTN-ALLOC-004", YELLOW}, // Toutes les règles = WARNING
		{"KTN-OTHER-999", YELLOW}, // Toutes les règles = WARNING
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := formatter.getCodeColor(tt.code)
			if got != tt.expected {
				t.Errorf("getCodeColor(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

// TestGetCodeColorNoColor tests the functionality of the corresponding implementation.
func TestGetCodeColorNoColor(t *testing.T) {
	formatter := &formatterImpl{noColor: true}
	got := formatter.getCodeColor("KTN-VAR-001")
	if got != "" {
		t.Errorf("Expected empty string with noColor=true, got %q", got)
	}
}

// TestGroupByFile tests the functionality of the corresponding implementation.
func TestGroupByFile(t *testing.T) {
	const EXPECTED_GROUP_COUNT int = 2

	tests := []struct {
		name  string
		check func(t *testing.T, groups []DiagnosticGroupData, fset *token.FileSet)
	}{
		{
			name: "correct group count",
			check: func(t *testing.T, groups []DiagnosticGroupData, fset *token.FileSet) {
				// Vérification du nombre de groupes
				if len(groups) != EXPECTED_GROUP_COUNT {
					t.Errorf("Expected %d groups, got %d", EXPECTED_GROUP_COUNT, len(groups))
				}
			},
		},
		{
			name: "groups sorted by filename",
			check: func(t *testing.T, groups []DiagnosticGroupData, fset *token.FileSet) {
				// Check sorting (by filename)
				if groups[0].Filename > groups[1].Filename {
					t.Error("Groups should be sorted by filename")
				}
			},
		},
		{
			name: "diagnostics sorted by line within groups",
			check: func(t *testing.T, groups []DiagnosticGroupData, fset *token.FileSet) {
				// Check that diagnostics are sorted by line within each group
				for _, group := range groups {
					// Itération sur les diagnostics
					for i := 1; i < len(group.Diagnostics); i++ {
						posI := fset.Position(group.Diagnostics[i-1].Pos)
						posJ := fset.Position(group.Diagnostics[i].Pos)
						// Vérification de l'ordre
						if posI.Line > posJ.Line {
							t.Error("Diagnostics should be sorted by line number")
						}
					}
				}
			},
		},
	}

	// Préparation commune
	formatter := &formatterImpl{}
	fset := token.NewFileSet()
	file1 := fset.AddFile("file1.go", 1, 1000)
	file2 := fset.AddFile("file2.go", 1002, 1000)
	diagnostics := []analysis.Diagnostic{
		{Pos: file1.Pos(10), Message: "Issue 1", Category: "test"},
		{Pos: file2.Pos(20), Message: "Issue 2", Category: "test"},
		{Pos: file1.Pos(5), Message: "Issue 3", Category: "test"},
	}
	groups := formatter.groupByFile(fset, diagnostics)

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, groups, fset)
		})
	}
}

// TestGroupByFileFiltering tests the functionality of the corresponding implementation.
func TestGroupByFileFiltering(t *testing.T) {
	formatter := &formatterImpl{}
	fset := token.NewFileSet()

	// Add files including temp/cache files
	file1 := fset.AddFile("normal.go", 1, 1000)
	file2 := fset.AddFile("/.cache/go-build/temp.go", 1002, 1000)
	file3 := fset.AddFile("/tmp/test.go", 2003, 1000)

	diagnostics := []analysis.Diagnostic{
		{Pos: file1.Pos(10), Message: "Issue 1", Category: "test"},
		{Pos: file2.Pos(20), Message: "Issue 2", Category: "test"}, // Should be filtered
		{Pos: file3.Pos(30), Message: "Issue 3", Category: "test"}, // Should be filtered
	}

	groups := formatter.groupByFile(fset, diagnostics)

	if len(groups) != 1 {
		t.Errorf("Expected 1 group after filtering, got %d", len(groups))
	}
	if groups[0].Filename != "normal.go" {
		t.Errorf("Expected normal.go, got %s", groups[0].Filename)
	}
}

// TestFilterAndSortDiagnostics tests the functionality of the corresponding implementation.
func TestFilterAndSortDiagnostics(t *testing.T) {
	const EXPECTED_DIAG_COUNT int = 3

	tests := []struct {
		name  string
		check func(t *testing.T, filtered []analysis.Diagnostic, fset *token.FileSet)
	}{
		{
			name: "correct diagnostic count",
			check: func(t *testing.T, filtered []analysis.Diagnostic, fset *token.FileSet) {
				// Vérification du nombre de diagnostics
				if len(filtered) != EXPECTED_DIAG_COUNT {
					t.Errorf("Expected %d diagnostics, got %d", EXPECTED_DIAG_COUNT, len(filtered))
				}
			},
		},
		{
			name: "first diagnostic from a.go",
			check: func(t *testing.T, filtered []analysis.Diagnostic, fset *token.FileSet) {
				// Check sorting: by filename, then line, then column
				positionStrings := make([]string, len(filtered))
				// Itération sur les diagnostics
				for i, diag := range filtered {
					pos := fset.Position(diag.Pos)
					positionStrings[i] = pos.String()
				}

				// Should be: a.go line 10, a.go line 20, b.go line 10
				if !strings.Contains(positionStrings[0], "a.go") {
					t.Error("First diagnostic should be from a.go")
				}
			},
		},
		{
			name: "last diagnostic from b.go",
			check: func(t *testing.T, filtered []analysis.Diagnostic, fset *token.FileSet) {
				// Calcul des positions
				positionStrings := make([]string, len(filtered))
				// Itération sur les diagnostics
				for i, diag := range filtered {
					pos := fset.Position(diag.Pos)
					positionStrings[i] = pos.String()
				}

				// Vérification du dernier diagnostic
				if !strings.Contains(positionStrings[2], "b.go") {
					t.Error("Last diagnostic should be from b.go")
				}
			},
		},
	}

	// Préparation commune
	formatter := &formatterImpl{}
	fset := token.NewFileSet()
	file1 := fset.AddFile("a.go", 1, 1000)
	file2 := fset.AddFile("b.go", 1002, 1000)
	diagnostics := []analysis.Diagnostic{
		{Pos: file2.Pos(10), Message: "B1", Category: "test"},
		{Pos: file1.Pos(20), Message: "A2", Category: "test"},
		{Pos: file1.Pos(10), Message: "A1", Category: "test"},
	}
	filtered := formatter.filterAndSortDiagnostics(fset, diagnostics)

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, filtered, fset)
		})
	}
}

// TestPrintFunctions tests the functionality of the corresponding implementation.
func TestPrintFunctions(t *testing.T) {
	tests := []struct {
		name    string
		noColor bool
	}{
		{"with color", false},
		{"without color", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{
				writer:  buf,
				noColor: tt.noColor,
			}

			// Test printSuccess
			formatter.printSuccess()
			if !strings.Contains(buf.String(), "No issues found") {
				t.Error("printSuccess should contain success message")
			}

			// Test printHeader
			buf.Reset()
			formatter.printHeader(5)
			if !strings.Contains(buf.String(), "5 issue(s) found") {
				t.Error("printHeader should contain issue count")
			}

			// Test printFileHeader
			buf.Reset()
			formatter.printFileHeader("test.go", 3)
			if !strings.Contains(buf.String(), "test.go") {
				t.Error("printFileHeader should contain filename")
			}

			// Test printSummary
			buf.Reset()
			formatter.printSummary(10)
			if !strings.Contains(buf.String(), "issue(s) to fix") {
				t.Error("printSummary should contain 'issue(s) to fix'")
			}

			// Test printDiagnostic
			buf.Reset()
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", 1, 1000)
			pos := fset.Position(file.Pos(10))
			diag := analysis.Diagnostic{
				Message:  "KTN-TEST-001: Test issue\nDetails",
				Category: "test",
			}
			formatter.printDiagnostic(1, pos, diag)
			output := buf.String()
			if !strings.Contains(output, "test.go") {
				t.Error("printDiagnostic should contain filename")
			}
			if !strings.Contains(output, "KTN-TEST-001") {
				t.Error("printDiagnostic should contain error code")
			}
		})
	}
}

// TestFormatSimpleModeWithFiltering tests the functionality of the corresponding implementation.
func TestFormatSimpleModeWithFiltering(t *testing.T) {
	const EXPECTED_LINE_COUNT int = 3

	tests := []struct {
		name  string
		check func(t *testing.T, output string, lines []string)
	}{
		{
			name: "correct line count after filtering",
			check: func(t *testing.T, output string, lines []string) {
				// Should have 3 lines (from normal.go on different lines)
				if len(lines) != EXPECTED_LINE_COUNT {
					t.Errorf("Expected %d lines after filtering, got %d", EXPECTED_LINE_COUNT, len(lines))
				}
			},
		},
		{
			name: "contains normal.go in output",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification du fichier normal
				if !strings.Contains(output, "normal.go") {
					t.Error("Expected normal.go in output")
				}
			},
		},
		{
			name: "filters out cache and tmp files",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification du filtrage
				if strings.Contains(output, "temp.go") || strings.Contains(output, "/tmp/") {
					t.Error("Expected cache/tmp files to be filtered out")
				}
			},
		},
		{
			name: "sorts by line number - Issue 1 first",
			check: func(t *testing.T, output string, lines []string) {
				// Check sorting: line 2, then line 3, then line 4
				if !strings.Contains(lines[0], "Issue 1") {
					t.Error("First line should be Issue 1 (line 2)")
				}
			},
		},
		{
			name: "sorts by line number - Issue 4 second",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification de la deuxième ligne
				if !strings.Contains(lines[1], "Issue 4") {
					t.Error("Second line should be Issue 4 (line 3)")
				}
			},
		},
		{
			name: "sorts by line number - Issue 3 third",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification de la troisième ligne
				if !strings.Contains(lines[2], "Issue 3") {
					t.Error("Third line should be Issue 3 (line 4)")
				}
			},
		},
	}

	// Préparation commune
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, false, true)
	fset := token.NewFileSet()
	file1 := fset.AddFile("normal.go", 1, 1000)
	file1.AddLine(10)
	file1.AddLine(50)
	file1.AddLine(100)
	file2 := fset.AddFile("/.cache/go-build/temp.go", 1002, 1000)
	file3 := fset.AddFile("/tmp/test.go", 2003, 1000)
	diagnostics := []analysis.Diagnostic{
		{Pos: file1.Pos(100), Message: "KTN-TEST-003: Issue 3", Category: "test"},
		{Pos: file1.Pos(10), Message: "KTN-TEST-001: Issue 1", Category: "test"},
		{Pos: file2.Pos(20), Message: "KTN-TEST-002: Issue 2", Category: "test"},
		{Pos: file1.Pos(50), Message: "KTN-TEST-004: Issue 4", Category: "test"},
		{Pos: file3.Pos(30), Message: "KTN-TEST-005: Issue 5", Category: "test"},
	}
	formatter.Format(fset, diagnostics)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, output, lines)
		})
	}
}

// TestFormatHumanModeEmpty tests the functionality of the corresponding implementation.
func TestFormatHumanModeEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := &formatterImpl{
		writer:     buf,
		noColor:    false,
		aiMode:     false,
		simpleMode: false,
	}
	fset := token.NewFileSet()

	// Call formatForHuman directly with empty diagnostics
	formatter.formatForHuman(fset, []analysis.Diagnostic{})

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected success message for empty diagnostics, got: %s", output)
	}
}
