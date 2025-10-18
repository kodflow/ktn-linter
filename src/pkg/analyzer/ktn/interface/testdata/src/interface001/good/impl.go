package good

// serviceImpl is the private implementation
type serviceImpl struct {
	name string
}

func (s *serviceImpl) DoSomething() string {
	return s.name
}

func NewService(name string) Service {
	return &serviceImpl{name: name}
}
