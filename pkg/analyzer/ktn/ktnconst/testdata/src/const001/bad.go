// Package const001 contains test cases for KTN-CONST-001.
// This file contains ALL cases that MUST trigger KTN-CONST-001 errors.
package const001

// =============================================================================
// SECTION 1: Basic types without explicit type
// =============================================================================

const (
	// BadInt is an integer without explicit type
	BadInt = 42 // want "KTN-CONST-001"

	// BadNegativeInt is a negative integer without explicit type
	BadNegativeInt = -100 // want "KTN-CONST-001"

	// BadString is a string without explicit type
	BadString = "hello" // want "KTN-CONST-001"

	// BadEmptyString is an empty string without explicit type
	BadEmptyString = "" // want "KTN-CONST-001"

	// BadBoolTrue is a boolean true without explicit type
	BadBoolTrue = true // want "KTN-CONST-001"

	// BadBoolFalse is a boolean false without explicit type
	BadBoolFalse = false // want "KTN-CONST-001"

	// BadFloat is a float without explicit type
	BadFloat = 3.14159 // want "KTN-CONST-001"

	// BadNegativeFloat is a negative float without explicit type
	BadNegativeFloat = -2.5 // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 2: Numeric literal formats without explicit type
// =============================================================================

const (
	// BadHex is a hexadecimal literal without explicit type
	BadHex = 0xFF // want "KTN-CONST-001"

	// BadHexLower is a lowercase hex literal without explicit type
	BadHexLower = 0xabcdef // want "KTN-CONST-001"

	// BadOctal is an octal literal without explicit type
	BadOctal = 0o755 // want "KTN-CONST-001"

	// BadOctalOld is an old-style octal literal without explicit type
	BadOctalOld = 0644 // want "KTN-CONST-001"

	// BadBinary is a binary literal without explicit type
	BadBinary = 0b1010 // want "KTN-CONST-001"

	// BadScientific is a scientific notation without explicit type
	BadScientific = 1e10 // want "KTN-CONST-001"

	// BadScientificNeg is a negative exponent without explicit type
	BadScientificNeg = 1e-5 // want "KTN-CONST-001"

	// BadUnderscored is a number with underscores without explicit type
	BadUnderscored = 1_000_000 // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 3: Rune/character literals without explicit type
// =============================================================================

const (
	// BadRuneA is a rune without explicit type
	BadRuneA = 'a' // want "KTN-CONST-001"

	// BadRuneNewline is a newline rune without explicit type
	BadRuneNewline = '\n' // want "KTN-CONST-001"

	// BadRuneTab is a tab rune without explicit type
	BadRuneTab = '\t' // want "KTN-CONST-001"

	// BadRuneUnicode is a unicode rune without explicit type
	BadRuneUnicode = '世' // want "KTN-CONST-001"

	// BadRuneHex is a hex escape rune without explicit type
	BadRuneHex = '\x00' // want "KTN-CONST-001"

	// BadRuneOctal is an octal escape rune without explicit type
	BadRuneOctal = '\000' // want "KTN-CONST-001"

	// BadRuneUnicodeEsc is a unicode escape without explicit type
	BadRuneUnicodeEsc = '\u4e16' // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 4: Complex numbers without explicit type
// =============================================================================

const (
	// BadComplex is a complex number without explicit type
	BadComplex = 1 + 2i // want "KTN-CONST-001"

	// BadComplexPure is a pure imaginary without explicit type
	BadComplexPure = 3i // want "KTN-CONST-001"

	// BadComplexNeg is a negative complex without explicit type
	BadComplexNeg = -1 - 2i // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 5: String variants without explicit type
// =============================================================================

const (
	// BadRawString is a raw string without explicit type
	BadRawString = `raw string` // want "KTN-CONST-001"

	// BadMultilineRaw is a multiline raw string without explicit type
	BadMultilineRaw = `line1
line2` // want "KTN-CONST-001"

	// BadUnicodeString is a unicode string without explicit type
	BadUnicodeString = "Hello, 世界" // want "KTN-CONST-001"

	// BadEscapedString is an escaped string without explicit type
	BadEscapedString = "tab:\there" // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 6: Iota without explicit type on first line
// =============================================================================

const (
	// BadIotaFirst without explicit type starts the sequence
	BadIotaFirst = iota // want "KTN-CONST-001"
	// BadIotaSecond inherits (no error - no value)
	BadIotaSecond
	// BadIotaThird inherits (no error - no value)
	BadIotaThird
)

const (
	// BadIotaExpr uses iota in expression without type
	BadIotaExpr = 1 << iota // want "KTN-CONST-001"
	// BadIotaExprTwo inherits (no error)
	BadIotaExprTwo
)

// =============================================================================
// SECTION 7: Expressions without explicit type
// =============================================================================

const (
	// BadExprAdd is an addition without explicit type
	BadExprAdd = 10 + 5 // want "KTN-CONST-001"

	// BadExprMul is a multiplication without explicit type
	BadExprMul = 3 * 4 // want "KTN-CONST-001"

	// BadExprDiv is a division without explicit type
	BadExprDiv = 100 / 4 // want "KTN-CONST-001"

	// BadExprShift is a shift without explicit type
	BadExprShift = 1 << 10 // want "KTN-CONST-001"

	// BadExprBitOr is a bitwise or without explicit type
	BadExprBitOr = 0x0F | 0xF0 // want "KTN-CONST-001"

	// BadExprLen uses len without explicit type
	BadExprLen = len("hello") // want "KTN-CONST-001"
)

// =============================================================================
// SECTION 8: Multi-name declarations without explicit type
// =============================================================================

const (
	// BadMultiInt declares multiple ints without type
	BadMultiIntA, BadMultiIntB = 1, 2 // want "KTN-CONST-001" "KTN-CONST-001"

	// BadMultiMixed declares mixed values without type
	BadMultiMixedA, BadMultiMixedB = "a", "b" // want "KTN-CONST-001" "KTN-CONST-001"
)

// =============================================================================
// SECTION 9: Edge cases
// =============================================================================

const (
	// BadZero is zero without explicit type
	BadZero = 0 // want "KTN-CONST-001"

	// BadOne is one without explicit type
	BadOne = 1 // want "KTN-CONST-001"

	// BadMinusOne is minus one without explicit type
	BadMinusOne = -1 // want "KTN-CONST-001"

	// BadMaxInt64 is max int64 value without explicit type
	BadMaxInt64 = 9223372036854775807 // want "KTN-CONST-001"

	// BadMinInt64 is min int64 value without explicit type
	BadMinInt64 = -9223372036854775808 // want "KTN-CONST-001"
)
