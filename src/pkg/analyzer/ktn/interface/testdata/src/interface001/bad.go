package interface001 // want `\[KTN_INTERFACE_001\] Package 'interface001' sans fichier interfaces.go`

// Mauvais : package avec struct publique mais sans interfaces.go
type PublicService struct {
	name string
}

func (s *PublicService) DoSomething() string {
	return s.name
}
