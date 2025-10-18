package struct004

import "fmt"

// CORRECT: Struct avec moins de 15 champs

// GoodConfig a un nombre acceptable de champs.
type GoodConfig struct {
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
}

// EmptyStruct est une struct vide (pas d'erreur).
type EmptyStruct struct {
}

// SmallStruct a peu de champs.
type SmallStruct struct {
	Field1 string
	Field2 int
}

// BAD: Struct avec plus de 15 champs

// BadTooBig a trop de champs (16 > 15).
type BadTooBig struct { // want "KTN-STRUCT-004.*trop de champs"
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
	Field16 string // Champ en trop
}

// AnotherBadStruct a beaucoup trop de champs.
type AnotherBadStruct struct { // want "KTN-STRUCT-004.*trop de champs"
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
	Field16 string
	Field17 int
	Field18 bool
	Field19 float64
	Field20 string
}

// EmbeddedBase est une struct de base.
type EmbeddedBase struct {
	BaseField string
}

// MultipleEmbedded a plusieurs champs embedded.
type MultipleEmbedded struct {
	EmbeddedBase
	Field1 string
	Field2, Field3 int
}

// EmbeddedWithTooMany a trop de champs (embedded + nommés).
type EmbeddedWithTooMany struct { // want "KTN-STRUCT-004.*trop de champs"
	EmbeddedBase
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
}

// TooManyWithMultiplePerLine a trop de champs déclarés sur plusieurs lignes.
type TooManyWithMultiplePerLine struct { // want "KTN-STRUCT-004.*trop de champs"
	Field1, Field2, Field3, Field4 string
	Field5, Field6, Field7, Field8 int
	Field9, Field10, Field11, Field12 bool
	Field13, Field14, Field15, Field16 float64
}

// Types non-struct (pas d'erreur)
// MyInt is a custom type.
type MyInt int

// MyFunc est une fonction (pas d'erreur).
func MyFunc() {
	fmt.Println("test")
}
