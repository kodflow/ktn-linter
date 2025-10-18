package interface006

// Mauvais : interface sans constructeur
type Service interface { // want `\[KTN_INTERFACE_006\] Interface 'Service' sans constructeur`
	DoSomething() string
}

// Mauvais : autre interface sans constructeur
type Processor interface { // want `\[KTN_INTERFACE_006\] Interface 'Processor' sans constructeur`
	Process(data string) error
}
