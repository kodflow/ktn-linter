# Exemples de Code Conforme KTN-Linter v1.3.1

## Fonction avec Params et Returns

```go
// ProcessData processes the input data and returns the result.
//
// Params:
//   - ctx: Operation context
//   - data: Input data to process
//   - options: Processing options
//
// Returns:
//   - string: Processed result
//   - error: Error if processing fails
func ProcessData(ctx context.Context, data string, options map[string]any) (string, error) {
	// Validate input data
	if data == "" {
		return "", errors.New("data cannot be empty")
	}

	// Process the data
	return data, nil
}
```

## Fonction sans Params

```go
// GetDefaultConfig returns the default configuration.
//
// Returns:
//   - Config: Default configuration object
func GetDefaultConfig() Config {
	return Config{}
}
```

## Fonction sans Returns

```go
// LogMessage logs a message to the console.
//
// Params:
//   - msg: Message to log
func LogMessage(msg string) {
	fmt.Println(msg)
}
```

## Paramètres inutilisés avec underscore

```go
// Metadata populates provider metadata.
//
// Params:
//   - _: Context (unused)
//   - _: Metadata request (unused)
//   - resp: Response object to populate
func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}
```

## Struct avec Interface (KTN-STRUCT-002 + KTN-INTERFACE-001)

```go
// UserService manages user operations.
type UserService struct {
	db Database
}

// UserServiceInterface defines the contract for user operations.
type UserServiceInterface interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
}

// GetUser retrieves a user by ID.
//
// Params:
//   - id: User identifier
//
// Returns:
//   - *User: User object
//   - error: Error if user not found
func (s *UserService) GetUser(id string) (*User, error) {
	return nil, nil
}

// CreateUser creates a new user.
//
// Params:
//   - user: User to create
//
// Returns:
//   - error: Error if creation fails
func (s *UserService) CreateUser(user *User) error {
	return nil
}

// NewUserService creates a UserService instance.
//
// Params:
//   - db: Database connection
//
// Returns:
//   - UserServiceInterface: User service interface
func NewUserService(db Database) UserServiceInterface {
	return &UserService{db: db}
}
```

**✅ L'interface est utilisée** : retournée par `NewUserService()`, donc pas d'erreur KTN-INTERFACE-001.

## Return Slice Vide vs Nil

**Pour RETOUR de fonction** : Utiliser slice vide (KTN-RETURN-002)
```go
// GetItems returns available items.
//
// Returns:
//   - []string: List of items
func GetItems() []string {
	return []string{} // ✅ Slice vide
}
```

**Pour VARIABLE globale** : Utiliser nil (KTN-VAR-010)
```go
// Ignore KTN-VAR-010 here if you want, but this is the intended pattern:
var defaultItems []string // ✅ nil (économise allocation)
```

## Contraintes Framework Terraform (Exemption)

Pour `schema.Schema` qui vient du framework Terraform :

```go
// Schema defines the provider configuration schema.
//
// Params:
//   - _: Context (unused)
//   - _: Schema request (unused)
//   - resp: Response to populate
func (p *N8nProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// EXEMPTION KTN-VAR-014: schema.Schema est une contrainte du framework Terraform
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for n8n",
	}
}
```

**Note** : Ajouter un commentaire `// EXEMPTION KTN-XXX: raison` pour documenter l'exception.
