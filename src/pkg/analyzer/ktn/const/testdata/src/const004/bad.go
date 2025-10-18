package const004

const BadNoComment = "value" // want "KTN-CONST-004: la constante 'BadNoComment' doit avoir un commentaire associé"

const (
	BadMultiple1 = "value1" // want "KTN-CONST-004: la constante 'BadMultiple1' doit avoir un commentaire associé"
	BadMultiple2 = "value2" // want "KTN-CONST-004: la constante 'BadMultiple2' doit avoir un commentaire associé"
)

const BadTyped int = 42 // want "KTN-CONST-004: la constante 'BadTyped' doit avoir un commentaire associé"

const (
	BadName1, BadName2 = 1, 2 // want "KTN-CONST-004: la constante 'BadName1' doit avoir un commentaire associé" "KTN-CONST-004: la constante 'BadName2' doit avoir un commentaire associé"
)

const (
	BadIota1 = iota // want "KTN-CONST-004: la constante 'BadIota1' doit avoir un commentaire associé"
	BadIota2        // want "KTN-CONST-004: la constante 'BadIota2' doit avoir un commentaire associé"
	BadIota3        // want "KTN-CONST-004: la constante 'BadIota3' doit avoir un commentaire associé"
)

const BadString = "hello world" // want "KTN-CONST-004: la constante 'BadString' doit avoir un commentaire associé"

const BadBool = true // want "KTN-CONST-004: la constante 'BadBool' doit avoir un commentaire associé"

const BadFloat = 3.14 // want "KTN-CONST-004: la constante 'BadFloat' doit avoir un commentaire associé"

const BadComplex = 1 + 2i // want "KTN-CONST-004: la constante 'BadComplex' doit avoir un commentaire associé"

const BadExpr = 1 + 2*3 // want "KTN-CONST-004: la constante 'BadExpr' doit avoir un commentaire associé"
