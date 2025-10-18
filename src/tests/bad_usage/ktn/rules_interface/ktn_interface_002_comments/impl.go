package KTN_INTERFACE_002

// implementation sans commentaire propre
type implementation struct{}

func (i *implementation) Execute() error {
	// Early return from function.
	return nil
}

func (i *implementation) Validate() bool {
	// Continue inspection/processing.
	return true
}

func newImpl() *implementation {
	// Early return from function.
	return &implementation{}
}
