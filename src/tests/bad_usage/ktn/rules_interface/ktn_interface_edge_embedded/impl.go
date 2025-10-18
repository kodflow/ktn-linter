package badembedded

// implementation sans commentaire
type implementation struct {
	data []byte
}

func (i *implementation) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (i *implementation) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (i *implementation) Close() error {
	return nil
}

func (i *implementation) CustomMethod() error {
	return nil
}

// newImplementation constructeur qui retourne un type privé
func newImplementation() *implementation {
	return &implementation{}
}
