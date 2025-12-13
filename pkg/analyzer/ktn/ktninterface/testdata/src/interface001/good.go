// Good examples for the interface001 test case.
package interface001

// goodUsedInterface is used in function parameter.
type goodUsedInterface interface {
	Method()
}

// goodProcess uses the interface as parameter.
//
// Params:
//   - g: interface implémentant Method
func goodProcess(g goodUsedInterface) {
	g.Method()
}

// goodReturnInterface is used as return type.
type goodReturnInterface interface {
	Get() string
}

// goodFactory returns the interface.
//
// Returns:
//   - goodReturnInterface: interface implementation ou nil
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
//
// Params:
//   - c: interface composite à utiliser
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
//
// Params:
//   - mr: interface avec méthode Process
func (g goodStruct) goodMethod(mr goodMethodReceiver) {
	mr.Process()
}

// UserServiceInterface interface suit le pattern XXXInterface pour struct UserService.
// Elle est légitime pour le mocking même si non utilisée directement.
type UserServiceInterface interface {
	Process() error
}

// UserService struct pour laquelle UserServiceInterface existe.
// Cette struct valide le pattern XXXInterface.
type UserService struct {
	db string
}

// goodCompileTimeCheck is used via compile-time interface check.
type goodCompileTimeCheck interface {
	Validate() error
}

// goodImpl implements goodCompileTimeCheck.
type goodImpl struct{}

// Validate implements goodCompileTimeCheck.
//
// Returns:
//   - error: validation error if any
func (g *goodImpl) Validate() error {
	// Implementation
	return nil
}

// Compile-time interface check - prouve que goodImpl implémente goodCompileTimeCheck
var _ goodCompileTimeCheck = (*goodImpl)(nil)

// init utilise les fonctions privées
func init() {
	// Appel de goodProcess
	_ = goodProcess(goodUsedInterface{})
	// Appel de goodFactory
	goodFactory()
	// Appel de goodUseComposite
	_ = goodUseComposite(goodCompositeInterface{})
}
