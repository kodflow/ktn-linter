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
	tests := []struct {
		name            string
		aiMode          bool
		noColor         bool
		simpleMode      bool
		expectedMessage string
	}{
		{
			name:            "empty diagnostics in human mode",
			aiMode:          false,
			noColor:         false,
			simpleMode:      false,
			expectedMessage: "No issues found",
		},
		{
			name:            "empty diagnostics in simple mode",
			aiMode:          false,
			noColor:         true,
			simpleMode:      true,
			expectedMessage: "No issues found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, tt.aiMode, tt.noColor, tt.simpleMode, false)
			fset := token.NewFileSet()

			formatter.Format(fset, []analysis.Diagnostic{})

			output := buf.String()
			if !strings.Contains(output, tt.expectedMessage) {
				t.Errorf("Expected success message, got: %s", output)
			}
		})
	}
}

// TestFormatHumanMode tests the functionality of the corresponding implementation.
func TestFormatHumanMode(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "human mode output", contains: []string{"KTN-LINTER REPORT", "test.go"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, false, false, false, false)
			fset := token.NewFileSet()

			fset.AddFile("test.go", 1, 1000)
			diagnostics := createTestDiagnostics()

			formatter.Format(fset, diagnostics)

			output := buf.String()
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %q in output", expected)
				}
			}
		})
	}
}

// TestFormatAIMode tests the functionality of the corresponding implementation.
func TestFormatAIMode(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "AI mode output", contains: []string{"# KTN-Linter Report (AI Mode)", "## File:"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, true, false, false, false)
			fset := token.NewFileSet()

			fset.AddFile("test.go", 1, 1000)
			diagnostics := createTestDiagnostics()

			formatter.Format(fset, diagnostics)

			output := buf.String()
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %q in output", expected)
				}
			}
		})
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
	formatter := NewFormatter(buf, false, false, true, false)
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
	tests := []struct {
		name        string
		shouldError bool
	}{
		{name: "no ANSI color codes", shouldError: false},
		{name: "output without colors", shouldError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, false, true, false, false)
			fset := token.NewFileSet()

			fset.AddFile("test.go", 1, 1000)
			diagnostics := createTestDiagnostics()

			formatter.Format(fset, diagnostics)

			output := buf.String()

			// No ANSI color codes
			if strings.Contains(output, "\033[") {
				t.Error("Expected no ANSI color codes with noColor=true")
			}
		})
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
		{"KTN-VAR-001", Red},      // ERROR (camelCase pour var package)
		{"KTN-FUNC-002", Red},     // ERROR (context.Context en premier)
		{"KTN-TEST-003", Yellow},  // WARNING
		{"KTN-ALLOC-004", Yellow}, // WARNING (unknown defaults to WARNING)
		{"KTN-OTHER-999", Yellow}, // WARNING (unknown defaults to WARNING)
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
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{name: "no color for KTN-VAR-001", code: "KTN-VAR-001", expected: ""},
		{name: "no color for KTN-FUNC-002", code: "KTN-FUNC-002", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{noColor: true}
			got := formatter.getCodeColor(tt.code)
			if got != tt.expected {
				t.Errorf("Expected empty string with noColor=true, got %q", got)
			}
		})
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
	tests := []struct {
		name             string
		expectedGroups   int
		expectedFilename string
	}{
		{name: "filters only cache files, not tmp", expectedGroups: 2, expectedFilename: "/tmp/test.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{}
			fset := token.NewFileSet()

			file1 := fset.AddFile("normal.go", 1, 1000)
			file2 := fset.AddFile("/.cache/go-build/temp.go", 1002, 1000)
			file3 := fset.AddFile("/tmp/test.go", 2003, 1000)

			diagnostics := []analysis.Diagnostic{
				{Pos: file1.Pos(10), Message: "Issue 1", Category: "test"},
				{Pos: file2.Pos(20), Message: "Issue 2", Category: "test"},
				{Pos: file3.Pos(30), Message: "Issue 3", Category: "test"},
			}

			groups := formatter.groupByFile(fset, diagnostics)

			if len(groups) != tt.expectedGroups || groups[0].Filename != tt.expectedFilename {
				t.Errorf("Expected %d groups with filename %q, got %d groups",
					tt.expectedGroups, tt.expectedFilename, len(groups))
			}
		})
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
	const EXPECTED_LINE_COUNT int = 4

	tests := []struct {
		name  string
		check func(t *testing.T, output string, lines []string)
	}{
		{
			name: "correct line count after filtering",
			check: func(t *testing.T, output string, lines []string) {
				// Should have 4 lines (3 from normal.go + 1 from /tmp/test.go)
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
			name: "filters out only cache files, not tmp",
			check: func(t *testing.T, output string, lines []string) {
				// Vérification du filtrage - seuls les fichiers cache doivent être filtrés
				if strings.Contains(output, "temp.go") {
					t.Error("Expected cache files to be filtered out")
				}
				// /tmp/ files should NOT be filtered
				if !strings.Contains(output, "/tmp/") {
					t.Error("Expected /tmp/ files to be included")
				}
			},
		},
		{
			name: "sorts by filename then line - Issue 5 first (from /tmp)",
			check: func(t *testing.T, output string, lines []string) {
				// /tmp/test.go comes before normal.go alphabetically
				if !strings.Contains(lines[0], "Issue 5") {
					t.Error("First line should be Issue 5 (from /tmp/test.go)")
				}
			},
		},
		{
			name: "sorts by filename then line - Issue 1 second",
			check: func(t *testing.T, output string, lines []string) {
				// normal.go:10 - Issue 1
				if !strings.Contains(lines[1], "Issue 1") {
					t.Error("Second line should be Issue 1 (line 10)")
				}
			},
		},
		{
			name: "sorts by filename then line - Issue 4 third",
			check: func(t *testing.T, output string, lines []string) {
				// normal.go:50 - Issue 4
				if !strings.Contains(lines[2], "Issue 4") {
					t.Error("Third line should be Issue 4 (line 50)")
				}
			},
		},
	}

	// Préparation commune
	buf := &bytes.Buffer{}
	formatter := NewFormatter(buf, false, false, true, false)
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
	tests := []struct {
		name     string
		expected string
	}{
		{name: "empty diagnostics show success", expected: "No issues found"},
		{name: "formatForHuman with no errors", expected: "No issues found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected success message for empty diagnostics, got: %s", output)
			}
		})
	}
}
// TestFormatterImpl_formatForHuman tests formatForHuman private method
func TestFormatterImpl_formatForHuman(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{name: "shows test error", message: "KTN-VAR-001: test error", expected: "test error"},
		{name: "shows formatted output", message: "KTN-FUNC-002: another error", expected: "another error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, false, false, false, false).(*formatterImpl)
			fset := token.NewFileSet()
			fset.AddFile("test.go", 1, 100)

			diags := []analysis.Diagnostic{
				{Pos: token.Pos(1), Message: tt.message},
			}

			formatter.formatForHuman(fset, diags)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected error message in output")
			}
		})
	}
}

// TestFormatterImpl_formatForAI tests formatForAI private method
func TestFormatterImpl_formatForAI(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{name: "shows filename in output", filename: "test.go", expected: "test.go"},
		{name: "AI mode format", filename: "main.go", expected: "main.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, true, false, false, false).(*formatterImpl)
			fset := token.NewFileSet()
			fset.AddFile(tt.filename, 1, 100)

			diags := []analysis.Diagnostic{
				{Pos: token.Pos(1), Message: "KTN-VAR-001: test error"},
			}

			formatter.formatForAI(fset, diags)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected filename in AI output")
			}
		})
	}
}

// TestFormatterImpl_formatSimple tests formatSimple private method
func TestFormatterImpl_formatSimple(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{name: "simple format with filename", filename: "test.go", expected: "test.go:"},
		{name: "another file format", filename: "main.go", expected: "main.go:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := NewFormatter(buf, false, false, true, false).(*formatterImpl)
			fset := token.NewFileSet()
			fset.AddFile(tt.filename, 1, 100)

			diags := []analysis.Diagnostic{
				{Pos: token.Pos(1), Message: "KTN-VAR-001: test error"},
			}

			formatter.formatSimple(fset, diags)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected simple format with filename:line")
			}
		})
	}
}

// TestFormatterImpl_groupByFile tests groupByFile private method
func TestFormatterImpl_groupByFile(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
	}{
		{name: "groups by file correctly", expectedCount: 2},
		{name: "two files grouped", expectedCount: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{}
			fset := token.NewFileSet()
			file1 := fset.AddFile("test1.go", 1, 100)
			file2 := fset.AddFile("test2.go", 102, 100)

			diags := []analysis.Diagnostic{
				{Pos: file1.Pos(10), Message: "error1"},
				{Pos: file2.Pos(20), Message: "error2"},
				{Pos: file1.Pos(30), Message: "error3"},
			}

			grouped := formatter.groupByFile(fset, diags)

			if len(grouped) != tt.expectedCount {
				t.Errorf("Expected %d files, got %d", tt.expectedCount, len(grouped))
			}
		})
	}
}

// TestFormatterImpl_filterAndSortDiagnostics tests filterAndSortDiagnostics
func TestFormatterImpl_filterAndSortDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
	}{
		{name: "sorts and filters diagnostics", expectedCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{}
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", 1, 100)

			diags := []analysis.Diagnostic{
				{Pos: file.Pos(50), Message: "error2"},
				{Pos: file.Pos(10), Message: "error1"},
				{Pos: file.Pos(90), Message: "error3"},
			}

			sorted := formatter.filterAndSortDiagnostics(fset, diags)

			// Check count and sorting
			if len(sorted) != tt.expectedCount || (len(sorted) >= 2 && sorted[0].Pos > sorted[1].Pos) {
				t.Errorf("Expected %d sorted diagnostics", tt.expectedCount)
			}
		})
	}
}

// TestFormatterImpl_printHeader tests printHeader private method
func TestFormatterImpl_printHeader(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{name: "shows count 5", count: 5, expected: "5"},
		{name: "shows count 10", count: 10, expected: "10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{writer: buf, noColor: true}

			formatter.printHeader(tt.count)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected count in header")
			}
		})
	}
}

// TestFormatterImpl_printFileHeader tests printFileHeader private method
func TestFormatterImpl_printFileHeader(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		count    int
	}{
		{name: "test.go with 3 issues", filename: "test.go", count: 3},
		{name: "main.go with 5 issues", filename: "main.go", count: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{writer: buf, noColor: true}

			formatter.printFileHeader(tt.filename, tt.count)
			output := buf.String()

			if !strings.Contains(output, tt.filename) || !strings.Contains(output, string(rune(tt.count+'0'))) {
				t.Errorf("Expected filename and count in file header")
			}
		})
	}
}

// TestFormatterImpl_printDiagnostic tests printDiagnostic private method
func TestFormatterImpl_printDiagnostic(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{name: "prints test error", message: "KTN-VAR-001: test error\ndetails", expected: "test error"},
		{name: "prints another error", message: "KTN-FUNC-002: function error\nmore details", expected: "function error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{writer: buf, noColor: true}
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", 1, 100)

			diag := analysis.Diagnostic{
				Pos:     file.Pos(10),
				Message: tt.message,
			}

			pos := fset.Position(diag.Pos)
			formatter.printDiagnostic(1, pos, diag)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected error message in diagnostic output")
			}
		})
	}
}

// TestFormatterImpl_printSuccess tests printSuccess private method
func TestFormatterImpl_printSuccess(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{name: "shows success message", expected: "No issues found"},
		{name: "confirms no issues", expected: "No issues found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{writer: buf, noColor: true}

			formatter.printSuccess()
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected success message")
			}
		})
	}
}

// TestFormatterImpl_printSummary tests printSummary private method
func TestFormatterImpl_printSummary(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{name: "summary with 10 issues", count: 10},
		{name: "summary with 5 issues", count: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			formatter := &formatterImpl{writer: buf, noColor: true}

			formatter.printSummary(tt.count)
			output := buf.String()

			if !strings.Contains(output, "Total") {
				t.Errorf("Expected count in summary")
			}
		})
	}
}

// TestFormatterImpl_getCodeColor tests getCodeColor private method
func TestFormatterImpl_getCodeColor(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		noColor  bool
		wantNonEmpty bool
	}{
		{"with color", "KTN-VAR-001", false, true},
		{"no color", "KTN-VAR-001", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{noColor: tt.noColor}
			result := formatter.getCodeColor(tt.code)

			if tt.wantNonEmpty && result == "" {
				t.Errorf("Expected non-empty color code")
			}
			if !tt.wantNonEmpty && result != "" {
				t.Errorf("Expected empty color code with noColor=true")
			}
		})
	}
}

// TestFormatterImpl_getSymbol tests getSymbol private method
func TestFormatterImpl_getSymbol(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{name: "warning symbol", code: "KTN-VAR-001"},
		{name: "error symbol", code: "KTN-VAR-003"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &formatterImpl{}
			symbol := formatter.getSymbol(tt.code)
			if symbol == "" {
				t.Errorf("Expected non-empty symbol for %s", tt.code)
			}
		})
	}
}
