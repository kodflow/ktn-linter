package struct004

// BadUser documentation insuffisante (1 seule ligne) - VIOLATION
type BadUser struct { // want "KTN-STRUCT-004"
	Name string
	Age  int
}

// NoDoc struct exportée sans documentation - VIOLATION
type NoDoc struct { // want "KTN-STRUCT-004"
	Value string
}

type MissingDoc struct { // want "KTN-STRUCT-004"
	Data int
}

// Ceci est une mauvaise documentation - VIOLATION
// Car elle ne commence pas par le nom de la struct
type BadConfig struct { // want "KTN-STRUCT-004"
	Host string
	Port int
}
