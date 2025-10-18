package const002

// Bad example 1: scattered consts before var
const (
	FirstGroup1 string = "first"
	FirstGroup2 string = "group"
)

const SecondGroup1 string = "scattered" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"

var GlobalVar string = "test"

const BadConst1 string = "bad" // want "KTN-CONST-002: les constantes doivent être groupées et placées au-dessus des déclarations var"
