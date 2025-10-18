package interface001_with_interface // want `\[KTN_INTERFACE_001\] Package 'interface001_with_interface' sans fichier interfaces.go`

// Package avec interface publique mais sans interfaces.go
type MyInterface interface {
	DoSomething() error
}
