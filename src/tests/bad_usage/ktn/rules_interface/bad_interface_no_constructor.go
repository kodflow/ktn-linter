package rules_interface

// ANTI-PATTERN: Interfaces sans constructeur New*
// Viole KTN-INTERFACE-004

// Repository interface SANS constructeur
type Repository interface {
	Save(data string) error
	Load(id string) (string, error)
	Delete(id string) error
}

// Logger interface sans New*
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

// Validator interface sans constructeur
type Validator interface {
	Validate(data interface{}) error
	ValidateField(field string, value interface{}) error
}

// Transformer interface sans New*
type Transformer interface {
	Transform(input string) string
	Reverse(output string) string
}

// HTTPClient interface sans constructeur
type HTTPClient interface {
	Get(url string) ([]byte, error)
	Post(url string, data []byte) ([]byte, error)
	Put(url string, data []byte) error
	Delete(url string) error
}

// Serializer interface sans New*
type Serializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(data []byte, obj interface{}) error
}

// EventBus interface sans constructeur
type EventBus interface {
	Publish(event string, data interface{}) error
	Subscribe(event string, handler func(interface{})) error
}
