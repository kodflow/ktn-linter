package interface002

// Bon : interface publique dans interfaces.go
type Service interface {
	DoSomething() string
}

// Bon : autre interface publique
type Processor interface {
	Process(data string) error
}
