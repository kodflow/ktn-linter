package interface001

// goodUsedInterface is used in function parameter.
type goodUsedInterface interface {
	Method()
}

// goodProcess uses the interface as parameter.
func goodProcess(g goodUsedInterface) {
	g.Method()
}

// goodReturnInterface is used as return type.
type goodReturnInterface interface {
	Get() string
}

// goodFactory returns the interface.
func goodFactory() goodReturnInterface {
	return nil
}

// goodFieldInterface is used in struct field.
type goodFieldInterface interface {
	Execute()
}

// goodContainer contains the interface as field.
type goodContainer struct {
	handler goodFieldInterface
}

// goodEmbeddedInterface is embedded in another interface.
type goodEmbeddedInterface interface {
	Base()
}

// goodCompositeInterface embeds goodEmbeddedInterface.
type goodCompositeInterface interface {
	goodEmbeddedInterface
	Extended()
}

// goodUseComposite uses the composite interface.
func goodUseComposite(c goodCompositeInterface) {
	c.Extended()
}

// goodMethodReceiver is used as method parameter.
type goodMethodReceiver interface {
	Process()
}

// goodStruct uses the interface in method.
type goodStruct struct{}

// goodMethod accepts the interface.
func (g goodStruct) goodMethod(mr goodMethodReceiver) {
	mr.Process()
}
