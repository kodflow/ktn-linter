package const002

// Edge case 1: only consts, no vars - single group is OK
const (
	OnlyConst1 string = "only1"
	OnlyConst2 string = "only2"
)

// Edge case 2: multiple const groups when no vars exist - scattered
const GroupA1 string = "a1" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"

const GroupB1 string = "b1" // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
