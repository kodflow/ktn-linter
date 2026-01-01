// Package var002 provides good test cases.
package var002

// Good: Variables with explicit type AND value (format: var name type = value)

const (
	// DefaultTimeout is the default timeout
	DefaultTimeout int = 30
	// DefaultPort is the default port
	DefaultPort int = 8080
	// MaxConnections is the maximum connections
	MaxConnections int = 100
	// TheAnswer is the answer to everything
	TheAnswer int = 42
	// BufferSize is the buffer size
	BufferSize int = 1024
	// ByteH is the byte value for H
	ByteH byte = 72
	// ByteE is the byte value for e
	ByteE byte = 101
	// ByteL is the byte value for l
	ByteL byte = 108
	// ByteO is the byte value for o
	ByteO byte = 111
	// DefaultRetries is the default retry count
	DefaultRetries int = 3
	// DefaultCacheSize is the default cache size
	DefaultCacheSize int = 10
)

// Style requis: var name type (= value optionnel)
var (
	// ===== Avec type explicite ET valeur =====

	// defaultRetries has explicit type and value
	defaultRetries int = DefaultTimeout
	// configuration has explicit type and value
	configuration string = "default"
	// isEnabled has explicit type and value
	isEnabled bool = false
	// serverPort has explicit type and value
	serverPort int = DefaultPort
	// serverHost has explicit type and value
	serverHost string = "localhost"
	// maxConnections has explicit type and value
	maxConnections int = MaxConnections
	// endpoints has explicit type and value
	endpoints []string = []string{"http://localhost:8080", "http://localhost:9090"}
	// configMap has explicit type and value
	configMap map[string]int = map[string]int{"timeout": DefaultTimeout, "retries": DefaultRetries}
	// buffer has explicit type and value
	buffer []byte = make([]byte, 0, BufferSize)
	// cache has explicit type and value
	cache map[string]string = make(map[string]string, DefaultCacheSize)
	// convertedInt has explicit type and value
	convertedInt int = int(TheAnswer)
	// convertedStr has explicit type and value
	convertedStr string = string([]byte{ByteH, ByteE, ByteL, ByteL, ByteO})
	// convertedFloat has explicit type and value
	convertedFloat float64 = float64(TheAnswer)

	// ===== Zéro-values (type explicite, pas de valeur) =====

	// zeroInt uses zero-value (idiomatic Go)
	zeroInt int
	// zeroString uses zero-value
	zeroString string
	// zeroBool uses zero-value
	zeroBool bool
	// zeroSlice uses zero-value (nil slice)
	zeroSlice []string
	// zeroMap uses zero-value (nil map)
	zeroMap map[string]int
)

// init utilise les variables pour éviter les erreurs de compilation
func init() {
	// Utilisation des variables privées
	_ = defaultRetries
	_ = configuration
	_ = isEnabled
	_ = serverPort
	_ = serverHost
	_ = maxConnections
	_ = endpoints
	_ = configMap
	_ = buffer
	_ = cache
	_ = convertedInt
	_ = convertedStr
	_ = convertedFloat
	// Utilisation des zéro-values
	_ = zeroInt
	_ = zeroString
	_ = zeroBool
	_ = zeroSlice
	_ = zeroMap
}
