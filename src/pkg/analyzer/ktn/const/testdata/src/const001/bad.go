package const001

// Mauvais : constante non groupée
const UngroupedConst = "value" // want `\[KTN_CONST_001\] Constante 'UngroupedConst' déclarée individuellement`

// Mauvais : plusieurs constantes non groupées
const AnotherConst int = 42 // want `\[KTN_CONST_001\] Constante 'AnotherConst' déclarée individuellement`

const ThirdConst bool = true // want `\[KTN_CONST_001\] Constante 'ThirdConst' déclarée individuellement`
