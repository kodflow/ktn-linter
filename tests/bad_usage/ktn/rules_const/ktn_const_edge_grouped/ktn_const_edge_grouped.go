package badconstgrouped

// Violations multiples dans const groupées

const (
	MaxRetries     = 3           // OK
	default_Value  = "test"      // Snake_case violation
	API_KEY        = "secret"    // Screaming snake case
	minTimeout     = 10          // Privé sans commentaire
	HTTPStatusCode = 200         // OK - HTTP est un initialisme
	htmlParser     = "parser"    // html devrait être HTML
	urlEndpoint    = "/api/test" // url devrait être URL
)

// Const sans groupement alors qu'elles sont liées
const DatabaseHost = "localhost"
const DatabasePort = 5432
const DatabaseName = "mydb"

// Mélange de types sans organisation
const (
	ConfigVersion  int    = 1
	ConfigName     string = "app"
	ConfigEnabled  bool   = true
	CONFIG_TIMEOUT int    = 30 // Screaming snake
)
