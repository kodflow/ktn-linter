package interface006

// Mauvais : interface sans constructeur
// Service defines the interface.
type Service interface { // want `\[KTN_INTERFACE_004\] Interface 'Service' sans constructeur`
	DoSomething() string
}

// Processor defines the interface.
// Mauvais : autre interface sans constructeur
type Processor interface { // want `\[KTN_INTERFACE_004\] Interface 'Processor' sans constructeur`
	Process(data string) error
}
