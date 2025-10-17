// Package badinterfaceswithstructs démontre une violation de KTN-INTERFACE-008.
package badinterfaceswithstructs

// Service est une interface valide.
type Service interface {
	Process() error
}

// serviceImpl est un struct PUBLIC dans interfaces.go - VIOLATION !
type ServiceImpl struct {
	data string
}

// Process implémente Service.
func (s *ServiceImpl) Process() error {
	return nil
}

// Repository autre interface valide.
type Repository interface {
	Save(data string) error
}

// repositoryImpl struct dans interfaces.go - VIOLATION !
type repositoryImpl struct {
	db string
}

// Save implémente Repository.
func (r *repositoryImpl) Save(data string) error {
	return nil
}
