package KTN_INTERFACE_002

// implementation sans commentaire propre
type implementation struct{}

func (i *implementation) Execute() error {
	return nil
}

func (i *implementation) Validate() bool {
	return true
}

func newImpl() *implementation {
	return &implementation{}
}
