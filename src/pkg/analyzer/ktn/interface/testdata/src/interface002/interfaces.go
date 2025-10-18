package interface002

// Bon : interface publique dans interfaces.go
// Service defines the interface.
type Service interface {
	DoSomething() string
}

// Processor defines the interface.
// Bon : autre interface publique
type Processor interface {
	Process(data string) error
}
