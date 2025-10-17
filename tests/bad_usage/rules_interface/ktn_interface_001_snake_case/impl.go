package KTN_INTERFACE_001

// implementation sans commentaire (violation)
type implementation struct{}

func (i *implementation) Process() error {
	return nil
}

// bad_constructor mauvais nommage (violation snake_case)
func bad_constructor() *implementation {
	return &implementation{}
}
