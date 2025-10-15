package naming_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// TestIsAllCaps teste IsAllCaps.
//
// Params:
//   - t: instance de test
func TestIsAllCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Cas ALL_CAPS
		{name: "simple caps", input: "MAXSIZE", expected: true},
		{name: "caps with underscore", input: "MAX_SIZE", expected: true},
		{name: "caps with multiple underscores", input: "MAX_BUFFER_SIZE", expected: true},
		{name: "single letter", input: "A", expected: true},
		{name: "caps with numbers", input: "HTTP2", expected: true},
		{name: "caps with numbers and underscore", input: "HTTP_2_SERVER", expected: true},

		// Cas MixedCaps
		{name: "mixed caps", input: "MaxSize", expected: false},
		{name: "camelCase", input: "maxSize", expected: false},
		{name: "lowercase", input: "maxsize", expected: false},
		{name: "mixed with underscore", input: "Max_Size", expected: false},

		// Cas spéciaux
		{name: "empty string", input: "", expected: false},
		{name: "only underscore", input: "_", expected: false},
		{name: "only numbers", input: "123", expected: false},
		{name: "underscore only", input: "___", expected: false},
		{name: "numbers and underscore", input: "123_456", expected: false},

		// Initialismes Go
		{name: "HTTP", input: "HTTP", expected: true},
		{name: "URL", input: "URL", expected: true},
		{name: "API", input: "API", expected: true},
		{name: "HTTPSPort", input: "HTTPSPORT", expected: true}, // Tout en majuscules
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.IsAllCaps(tt.input)
			if result != tt.expected {
				t.Errorf("naming.IsAllCaps(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsValidInitialism teste IsValidInitialism.
//
// Params:
//   - t: instance de test
func TestIsValidInitialism(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Initialismes valides simples
		{name: "HTTP", input: "HTTP", expected: true},
		{name: "HTTPS", input: "HTTPS", expected: true},
		{name: "URL", input: "URL", expected: true},
		{name: "URI", input: "URI", expected: true},
		{name: "ID", input: "ID", expected: true},
		{name: "API", input: "API", expected: true},
		{name: "JSON", input: "JSON", expected: true},
		{name: "XML", input: "XML", expected: true},
		{name: "HTML", input: "HTML", expected: true},
		{name: "SQL", input: "SQL", expected: true},
		{name: "TCP", input: "TCP", expected: true},
		{name: "UDP", input: "UDP", expected: true},
		{name: "IP", input: "IP", expected: true},
		{name: "DNS", input: "DNS", expected: true},
		{name: "SSH", input: "SSH", expected: true},
		{name: "FTP", input: "FTP", expected: true},
		{name: "OK", input: "OK", expected: true},
		{name: "EOF", input: "EOF", expected: true},

		// Combinaisons d'initialismes
		{name: "HTTPAPI", input: "HTTPAPI", expected: true},
		{name: "URLID", input: "URLID", expected: true},
		{name: "APIURL", input: "APIURL", expected: true},
		{name: "HTTPJSON", input: "HTTPJSON", expected: true},
		{name: "TCPIP", input: "TCPIP", expected: true},

		// Initialismes avec MixedCaps (valide selon Go)
		{name: "HTTPServer", input: "HTTPServer", expected: true},
		{name: "URLParser", input: "URLParser", expected: true},
		{name: "APIClient", input: "APIClient", expected: true},
		{name: "HTTPOK", input: "HTTPOK", expected: true},
		{name: "HTTPNotFound", input: "HTTPNotFound", expected: true},

		// Cas INVALIDES - avec underscore (viole KTN-VAR-008)
		{name: "HTTP_OK with underscore", input: "HTTP_OK", expected: false},
		{name: "MAX_BUFFER with underscore", input: "MAX_BUFFER", expected: false},
		{name: "API_URL with underscore", input: "API_URL", expected: false},
		{name: "HTTP_SERVER with underscore", input: "HTTP_SERVER", expected: false},

		// Cas INVALIDES - pas d'initialisme
		{name: "MaxSize", input: "MaxSize", expected: false},
		{name: "MAXSIZE", input: "MAXSIZE", expected: false},
		{name: "MAX", input: "MAX", expected: false},
		{name: "BUFFER", input: "BUFFER", expected: false},
		{name: "TIMEOUT", input: "TIMEOUT", expected: false},

		// Cas INVALIDES - initialismes partiels
		{name: "HTP partial", input: "HTP", expected: false},
		{name: "HTT partial", input: "HTT", expected: false},

		// Cas spéciaux
		{name: "empty", input: "", expected: false},
		{name: "lowercase", input: "http", expected: false},
		{name: "mixed wrong", input: "HttpOk", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.IsValidInitialism(tt.input)
			if result != tt.expected {
				t.Errorf("naming.IsValidInitialism(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsValidInitialismAllInitialisms teste IsValidInitialism AllInitialisms.
//
// Params:
//   - t: instance de test
func TestIsValidInitialismAllInitialisms(t *testing.T) {
	// Test tous les initialismes de la liste pour s'assurer qu'ils sont reconnus
	initialisms := []string{
		"HTTP", "HTTPS", "URL", "URI", "ID", "API", "JSON", "XML", "HTML",
		"SQL", "TLS", "SSL", "TCP", "UDP", "IP", "DNS", "SSH", "FTP",
		"OK", "EOF", "UID", "UUID", "ASCII", "UTF", "CPU", "RAM", "IO",
		"DB", "RPC", "CDN", "AWS", "GCP", "TTL", "ACL", "CORS", "CSRF",
	}

	for _, init := range initialisms {
		t.Run("initialism_"+init, func(t *testing.T) {
			result := naming.IsValidInitialism(init)
			if !result {
				t.Errorf("naming.IsValidInitialism(%q) = false, want true (should be valid initialism)", init)
			}
		})
	}
}

// TestIsValidInitialismCombinationsWithMixedCaps teste IsValidInitialism CombinationsWithMixedCaps.
//
// Params:
//   - t: instance de test
func TestIsValidInitialismCombinationsWithMixedCaps(t *testing.T) {
	// Test des combinaisons courantes avec MixedCaps
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "HTTPHandler", input: "HTTPHandler", expected: true},
		{name: "URLPath", input: "URLPath", expected: true},
		{name: "IDGenerator", input: "IDGenerator", expected: true},
		{name: "APIResponse", input: "APIResponse", expected: true},
		{name: "JSONEncoder", input: "JSONEncoder", expected: true},
		{name: "XMLParser", input: "XMLParser", expected: true},
		{name: "HTMLTemplate", input: "HTMLTemplate", expected: true},
		{name: "SQLQuery", input: "SQLQuery", expected: true},
		{name: "TCPListener", input: "TCPListener", expected: true},
		{name: "UDPSocket", input: "UDPSocket", expected: true},
		{name: "HTTPSConnection", input: "HTTPSConnection", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.IsValidInitialism(tt.input)
			if result != tt.expected {
				t.Errorf("naming.IsValidInitialism(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsAllCapsUnicode teste IsAllCaps Unicode.
//
// Params:
//   - t: instance de test
func TestIsAllCapsUnicode(t *testing.T) {
	// Test avec des caractères Unicode
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "greek uppercase", input: "ΑΒΓ", expected: true},
		{name: "greek lowercase", input: "αβγ", expected: false},
		{name: "cyrillic uppercase", input: "АБВ", expected: true},
		{name: "cyrillic lowercase", input: "абв", expected: false},
		{name: "mixed latin greek", input: "ABΓ", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.IsAllCaps(tt.input)
			if result != tt.expected {
				t.Errorf("naming.IsAllCaps(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsMixedCaps teste IsMixedCaps.
//
// Params:
//   - t: instance de test
func TestIsMixedCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valides
		{name: "MixedCaps exporté", input: "ParseHTTPRequest", expected: true},
		{name: "mixedCaps privé", input: "calculateTotal", expected: true},
		{name: "Simple exporté", input: "Process", expected: true},
		{name: "Simple privé", input: "validate", expected: true},
		{name: "HTTPServer", input: "HTTPServer", expected: true},
		{name: "URLParser", input: "URLParser", expected: true},
		{name: "UserID", input: "UserID", expected: true},

		// Invalides - snake_case
		{name: "snake_case", input: "parse_http_request", expected: false},
		{name: "Snake_Case", input: "Calculate_Total", expected: false},
		{name: "underscore début", input: "_private", expected: false},
		{name: "underscore milieu", input: "parse_request", expected: false},

		// Invalides - ALL_CAPS non initialisme
		{name: "ALL_CAPS", input: "MAX_SIZE", expected: false},
		{name: "MAXSIZE non initialisme", input: "MAXSIZE", expected: false},

		// Valides - initialismes complets
		{name: "HTTP seul", input: "HTTP", expected: true},
		{name: "URL seul", input: "URL", expected: true},
		{name: "ID seul", input: "ID", expected: true},

		// Cas spéciaux
		{name: "empty", input: "", expected: false},
		{name: "underscore seul", input: "_", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.IsMixedCaps(tt.input)
			if result != tt.expected {
				t.Errorf("naming.IsMixedCaps(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHasGetterPrefix teste HasGetterPrefix.
//
// Params:
//   - t: instance de test
func TestHasGetterPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Doit signaler
		{name: "GetUserName", input: "GetUserName", expected: true},
		{name: "GetEmail", input: "GetEmail", expected: true},
		{name: "GetHTTPClient", input: "GetHTTPClient", expected: true},
		{name: "GetID", input: "GetID", expected: true},
		{name: "GetValue", input: "GetValue", expected: true},

		// Ne doit PAS signaler - exceptions
		{name: "GetOrCreate", input: "GetOrCreate", expected: false},
		{name: "GetOrSet", input: "GetOrSet", expected: false},
		{name: "GetOrDefault", input: "GetOrDefault", expected: false},
		{name: "GetAndSet", input: "GetAndSet", expected: false},
		{name: "GetAndUpdate", input: "GetAndUpdate", expected: false},
		{name: "GetOrElse", input: "GetOrElse", expected: false},

		// Ne doit PAS signaler - pas de préfixe Get
		{name: "UserName", input: "UserName", expected: false},
		{name: "Email", input: "Email", expected: false},
		{name: "Parse", input: "Parse", expected: false},
		{name: "Calculate", input: "Calculate", expected: false},

		// Ne doit PAS signaler - Get sans majuscule après
		{name: "Getter", input: "Getter", expected: false},
		{name: "Get seul", input: "Get", expected: false},

		// Cas spéciaux
		{name: "empty", input: "", expected: false},
		{name: "G", input: "G", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.HasGetterPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("naming.HasGetterPrefix(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// containsSuggestion vérifie si une suggestion est présente dans les résultats.
//
// Params:
//   - results: les suggestions retournées
//   - expected: la suggestion attendue
//
// Returns:
//   - bool: true si la suggestion est trouvée
func containsSuggestion(results []string, expected string) bool {
	for _, r := range results {
		if r == expected {
			// Retourne true car la suggestion attendue a été trouvée
			return true
		}
	}
	// Retourne false car la suggestion n'a pas été trouvée
	return false
}

// TestFixInitialisms teste FixInitialisms.
//
// Params:
//   - t: instance de test
func TestFixInitialisms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// Corrections simples
		{name: "HttpServer", input: "HttpServer", expected: []string{"HTTPServer"}},
		{name: "UrlParser", input: "UrlParser", expected: []string{"URLParser"}},
		{name: "IdGenerator", input: "IdGenerator", expected: []string{"IDGenerator"}},
		{name: "ApiClient", input: "ApiClient", expected: []string{"APIClient"}},
		{name: "JsonEncoder", input: "JsonEncoder", expected: []string{"JSONEncoder"}},
		{name: "XmlParser", input: "XmlParser", expected: []string{"XMLParser"}},

		// Déjà correct - pas de suggestion
		{name: "HTTPServer correct", input: "HTTPServer", expected: []string{}},
		{name: "URLParser correct", input: "URLParser", expected: []string{}},
		{name: "UserName correct", input: "UserName", expected: []string{}},

		// Multiples initialismes incorrects
		{name: "HttpApiClient", input: "HttpApiClient", expected: []string{"HTTPAPIClient"}},
		{name: "UrlIdParser", input: "UrlIdParser", expected: []string{"URLIDParser"}},

		// Cas spéciaux
		{name: "empty", input: "", expected: []string{}},
		{name: "no initialism", input: "UserName", expected: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := naming.FixInitialisms(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("naming.FixInitialisms(%q) returned %d suggestions, want %d", tt.input, len(result), len(tt.expected))
				// Retourne car le nombre de suggestions ne correspond pas
				return
			}

			if len(tt.expected) > 0 && !containsSuggestion(result, tt.expected[0]) {
				t.Errorf("naming.FixInitialisms(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
