package const001

// Mauvais : constante non groupée

// UngroupedConst is an ungrouped constant.
const UngroupedConst string = "value" // want `\[KTN_CONST_001\] Constante 'UngroupedConst' déclarée individuellement`

// Mauvais : plusieurs constantes non groupées

// AnotherConst is another ungrouped constant.
const AnotherConst int = 42 // want `\[KTN_CONST_001\] Constante 'AnotherConst' déclarée individuellement`

// ThirdConst is a third ungrouped constant.
const ThirdConst bool = true // want `\[KTN_CONST_001\] Constante 'ThirdConst' déclarée individuellement`
