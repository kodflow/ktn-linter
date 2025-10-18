package interface005

// Mauvais : interface publique dans un fichier autre que interfaces.go
type Service interface { // want `\[KTN_INTERFACE_005\] Interface 'Service' définie dans bad.go`
	DoSomething() string
}

// Mauvais : autre interface publique hors interfaces.go
type Processor interface { // want `\[KTN_INTERFACE_005\] Interface 'Processor' définie dans bad.go`
	Process(data string) error
}
