package interface006

// Bon : interface avec constructeur
type Repository interface {
	Save(data string) error
	Load() (string, error)
}

type repositoryImpl struct {
	storage map[string]string
}

func (r *repositoryImpl) Save(data string) error {
	r.storage["key"] = data
	return nil
}

func (r *repositoryImpl) Load() (string, error) {
	return r.storage["key"], nil
}

// Bon : constructeur présent
func NewRepository() Repository {
	return &repositoryImpl{
		storage: make(map[string]string),
	}
}
