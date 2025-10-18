package interface005

// Bon : interface publique dans interfaces.go
type Repository interface {
	Save(data string) error
	Load() (string, error)
}

// Bon : autre interface publique dans interfaces.go
type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string)
}
