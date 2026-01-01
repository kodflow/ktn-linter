// Package const001 contains test cases for KTN-CONST-001.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-001 errors.
package const001

// =============================================================================
// SECTION 1: Basic types with explicit type (int, string, bool, float64)
// =============================================================================

const (
	// GoodInt is an integer with explicit type
	GoodInt int = 42

	// GoodNegativeInt is a negative integer with explicit type
	GoodNegativeInt int = -100

	// GoodString is a string with explicit type
	GoodString string = "hello"

	// GoodEmptyString is an empty string with explicit type
	GoodEmptyString string = ""

	// GoodBoolTrue is a boolean true with explicit type
	GoodBoolTrue bool = true

	// GoodBoolFalse is a boolean false with explicit type
	GoodBoolFalse bool = false

	// GoodFloat is a float with explicit type
	GoodFloat float64 = 3.14159

	// GoodNegativeFloat is a negative float with explicit type
	GoodNegativeFloat float64 = -2.5
)

// =============================================================================
// SECTION 2: All integer types with explicit type
// =============================================================================

const (
	// GoodInt8 is an int8 with explicit type
	GoodInt8 int8 = -128

	// GoodInt8Max is max int8 value
	GoodInt8Max int8 = 127

	// GoodInt16 is an int16 with explicit type
	GoodInt16 int16 = -32768

	// GoodInt16Max is max int16 value
	GoodInt16Max int16 = 32767

	// GoodInt32 is an int32 with explicit type
	GoodInt32 int32 = -2147483648

	// GoodInt32Max is max int32 value
	GoodInt32Max int32 = 2147483647

	// GoodInt64 is an int64 with explicit type
	GoodInt64 int64 = -9223372036854775808

	// GoodInt64Max is max int64 value
	GoodInt64Max int64 = 9223372036854775807
)

// =============================================================================
// SECTION 3: All unsigned integer types with explicit type
// =============================================================================

const (
	// GoodUint is a uint with explicit type
	GoodUint uint = 42

	// GoodUint8 is a uint8 with explicit type
	GoodUint8 uint8 = 255

	// GoodUint16 is a uint16 with explicit type
	GoodUint16 uint16 = 65535

	// GoodUint32 is a uint32 with explicit type
	GoodUint32 uint32 = 4294967295

	// GoodUint64 is a uint64 with explicit type
	GoodUint64 uint64 = 18446744073709551615

	// GoodUintptr is a uintptr with explicit type
	GoodUintptr uintptr = 0x1234

	// GoodByte is a byte (alias for uint8) with explicit type
	GoodByte byte = 0xFF
)

// =============================================================================
// SECTION 4: Float types with explicit type
// =============================================================================

const (
	// GoodFloat32 is a float32 with explicit type
	GoodFloat32 float32 = 3.14

	// GoodFloat32Neg is a negative float32
	GoodFloat32Neg float32 = -2.5

	// GoodFloat64 is a float64 with explicit type
	GoodFloat64 float64 = 3.14159265358979

	// GoodFloat64Scientific is scientific notation float64
	GoodFloat64Scientific float64 = 1e10

	// GoodFloat64ScientificNeg is negative exponent float64
	GoodFloat64ScientificNeg float64 = 1e-10
)

// =============================================================================
// SECTION 5: Complex types with explicit type
// =============================================================================

const (
	// GoodComplex64 is a complex64 with explicit type
	GoodComplex64 complex64 = 1 + 2i

	// GoodComplex128 is a complex128 with explicit type
	GoodComplex128 complex128 = 1 + 2i

	// GoodComplexPure is a pure imaginary with explicit type
	GoodComplexPure complex128 = 3i

	// GoodComplexNeg is a negative complex with explicit type
	GoodComplexNeg complex128 = -1 - 2i
)

// =============================================================================
// SECTION 6: Rune type with explicit type
// =============================================================================

const (
	// GoodRuneA is a rune with explicit type
	GoodRuneA rune = 'a'

	// GoodRuneNewline is a newline rune with explicit type
	GoodRuneNewline rune = '\n'

	// GoodRuneTab is a tab rune with explicit type
	GoodRuneTab rune = '\t'

	// GoodRuneUnicode is a unicode rune with explicit type
	GoodRuneUnicode rune = '世'

	// GoodRuneHex is a hex escape rune with explicit type
	GoodRuneHex rune = '\x00'

	// GoodRuneOctal is an octal escape rune with explicit type
	GoodRuneOctal rune = '\000'

	// GoodRuneUnicodeEsc is a unicode escape with explicit type
	GoodRuneUnicodeEsc rune = '\u4e16'

	// GoodRuneAsInt32 uses int32 (rune alias) with explicit type
	GoodRuneAsInt32 int32 = 'A'
)

// =============================================================================
// SECTION 7: Numeric literal formats with explicit type
// =============================================================================

const (
	// GoodHex is a hexadecimal literal with explicit type
	GoodHex int = 0xFF

	// GoodHexUint is a hex uint with explicit type
	GoodHexUint uint = 0xABCDEF

	// GoodOctal is an octal literal with explicit type
	GoodOctal int = 0o755

	// GoodOctalOld is an old-style octal with explicit type
	GoodOctalOld int = 0644

	// GoodBinary is a binary literal with explicit type
	GoodBinary int = 0b1010

	// GoodUnderscored is a number with underscores with explicit type
	GoodUnderscored int = 1_000_000

	// GoodHexUnderscored is hex with underscores with explicit type
	GoodHexUnderscored int64 = 0xFF_FF_FF_FF
)

// =============================================================================
// SECTION 8: String variants with explicit type
// =============================================================================

const (
	// GoodRawString is a raw string with explicit type
	GoodRawString string = `raw string`

	// GoodMultilineRaw is a multiline raw string with explicit type
	GoodMultilineRaw string = `line1
line2`

	// GoodUnicodeString is a unicode string with explicit type
	GoodUnicodeString string = "Hello, 世界"

	// GoodEscapedString is an escaped string with explicit type
	GoodEscapedString string = "tab:\there"

	// GoodEmptyRaw is an empty raw string with explicit type
	GoodEmptyRaw string = ``
)

// =============================================================================
// SECTION 9: Iota with explicit type on first line (inheritance OK)
// =============================================================================

const (
	// GoodIotaFirst has explicit type, others inherit
	GoodIotaFirst int = iota
	// GoodIotaSecond inherits type (no value = no error)
	GoodIotaSecond
	// GoodIotaThird inherits type (no value = no error)
	GoodIotaThird
)

const (
	// GoodIotaExpr uses iota in expression with explicit type
	GoodIotaExpr int = 1 << iota
	// GoodIotaExprTwo inherits (no value = no error)
	GoodIotaExprTwo
	// GoodIotaExprThree inherits (no value = no error)
	GoodIotaExprThree
)

const (
	// GoodIotaUint64 is iota with uint64 type
	GoodIotaUint64 uint64 = iota
	// GoodIotaUint64Two inherits
	GoodIotaUint64Two
)

// =============================================================================
// SECTION 10: Expressions with explicit type
// =============================================================================

const (
	// GoodExprAdd is an addition with explicit type
	GoodExprAdd int = 10 + 5

	// GoodExprMul is a multiplication with explicit type
	GoodExprMul int = 3 * 4

	// GoodExprDiv is a division with explicit type
	GoodExprDiv int = 100 / 4

	// GoodExprShift is a shift with explicit type
	GoodExprShift int = 1 << 10

	// GoodExprBitOr is a bitwise or with explicit type
	GoodExprBitOr int = 0x0F | 0xF0

	// GoodExprLen uses len with explicit type
	GoodExprLen int = len("hello")

	// GoodExprFloat is a float expression with explicit type
	GoodExprFloat float64 = 3.14 * 2.0
)

// =============================================================================
// SECTION 11: Multi-name declarations with explicit type
// =============================================================================

const (
	// GoodMultiInt declares multiple ints with explicit type
	GoodMultiIntA, GoodMultiIntB int = 1, 2

	// GoodMultiString declares multiple strings with explicit type
	GoodMultiStringA, GoodMultiStringB string = "a", "b"

	// GoodMultiFloat declares multiple floats with explicit type
	GoodMultiFloatA, GoodMultiFloatB float64 = 1.0, 2.0
)

// =============================================================================
// SECTION 12: Edge cases with explicit type
// =============================================================================

const (
	// GoodZero is zero with explicit type
	GoodZero int = 0

	// GoodOne is one with explicit type
	GoodOne int = 1

	// GoodMinusOne is minus one with explicit type
	GoodMinusOne int = -1
)

// =============================================================================
// SECTION 13: Blank identifier with value (should be skipped by rule)
// =============================================================================

const (
	// Blank identifier is intentionally ignored
	_ int = 0
	// Another blank with string
	_ string = "ignored"
)
