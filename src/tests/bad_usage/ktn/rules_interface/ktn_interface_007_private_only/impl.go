package emptyinterfaces

// implementation implémente privateHelper.
type implementation struct{}

// help implémente l'interface privateHelper.
func (i *implementation) help() {}

// NewImplementation crée une nouvelle instance.
//
// Returns:
//   - privateHelper: une nouvelle instance
func NewImplementation() privateHelper {
	return &implementation{}
}
