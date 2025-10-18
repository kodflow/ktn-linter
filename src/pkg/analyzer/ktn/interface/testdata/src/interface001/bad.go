package interface001 // want `\[KTN_INTERFACE_001\] Package 'interface001' sans fichier interfaces.go`

// Mauvais : package avec struct publique mais sans interfaces.go
// PublicService represents the struct.
// PublicService represents the struct.
type PublicService struct {
	name string
}

func (s *PublicService) DoSomething() string {
	return s.name
}
