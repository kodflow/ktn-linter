// Package badnomockfile démontre une violation de KTN-MOCK-001.
package badnomockfile

// Service est une interface qui nécessite un mock.go.
type Service interface {
	Process() error
	GetStatus() string
}

// Repository autre interface.
type Repository interface {
	Save(data string) error
}
