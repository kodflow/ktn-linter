package badembedded

// implementation sans commentaire
type implementation struct {
	data []byte
}

func (i *implementation) Read(p []byte) (n int, err error) {
	// Early return from function.
	return 0, nil
}

func (i *implementation) Write(p []byte) (n int, err error) {
	// Early return from function.
	return len(p), nil
}

func (i *implementation) Close() error {
	// Early return from function.
	return nil
}

func (i *implementation) CustomMethod() error {
	// Early return from function.
	return nil
}

// newImplementation constructeur qui retourne un type privé
func newImplementation() *implementation {
	// Early return from function.
	return &implementation{}
}
