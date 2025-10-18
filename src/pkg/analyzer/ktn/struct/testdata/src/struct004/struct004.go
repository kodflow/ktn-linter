package struct004

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
