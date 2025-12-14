// Good examples for the struct002 test case.
package struct001

// UserService interface reprend toutes les méthodes publiques de userServiceImpl
type UserService interface {
	Create(name string) error
	GetByID(id int) (string, error)
	Update(id int, name string) error
}

// userServiceImpl implémente UserService avec toutes les méthodes
type userServiceImpl struct {
	users map[int]string
}

// useUserService utilise l'interface UserService.
//
// Params:
//   - us: service utilisateur
func useUserService(us UserService) {
	// Utilise l'interface
	_ = us
}

// Create crée un utilisateur.
//
// Params:
//   - name: nom de l'utilisateur
//
// Returns:
//   - error: erreur éventuelle
func (u *userServiceImpl) Create(name string) error {
	// Implementation
	return nil
}

// GetByID récupère un utilisateur par ID.
//
// Params:
//   - id: identifiant utilisateur
//
// Returns:
//   - string: nom de l'utilisateur
//   - error: erreur éventuelle
func (u *userServiceImpl) GetByID(id int) (string, error) {
	// Implementation
	return "", nil
}

// Update met à jour un utilisateur.
//
// Params:
//   - id: identifiant utilisateur
//   - name: nouveau nom
//
// Returns:
//   - error: erreur éventuelle
func (u *userServiceImpl) Update(id int, name string) error {
	// Implementation
	return nil
}

// helper est une fonction privée - pas besoin dans l'interface
func (u *userServiceImpl) helper() {
	// Private method
}

// Config est une struct simple sans méthode - PAS BESOIN D'INTERFACE.
// Contient la configuration de connexion réseau.
type Config struct {
	Host string
	Port int
}

// DataModel est une struct DTO sans méthode - PAS BESOIN D'INTERFACE.
// Représente un modèle de données avec identifiant, nom et tags.
type DataModel struct {
	ID   int
	Name string
	Tags []string
}

// ExternalInterface représente une interface externe (simulée pour le test).
type ExternalInterface interface {
	Process() error
}

// RepositoryImpl implémente une interface externe (DDD pattern).
// L'interface est ici mais dans un vrai cas DDD elle serait dans domain/.
type RepositoryImpl struct {
	data map[string]string
}

// Process implémente ExternalInterface.
//
// Returns:
//   - error: erreur éventuelle
func (r *RepositoryImpl) Process() error {
	// Implementation
	return nil
}

// Compile-time interface check - prouve que RepositoryImpl implémente ExternalInterface
// Dans un vrai cas DDD, ce serait: var _ domain.Repository = (*RepositoryImpl)(nil)
var _ ExternalInterface = (*RepositoryImpl)(nil)

// init utilise les fonctions privées
func init() {
	// Appel de useUserService avec une implémentation concrète
	useUserService(&userServiceImpl{users: map[int]string{}})
}

// --- CONSUMER PATTERN ---
// Un consommateur est une struct qui utilise l'injection de dépendances.
// Ces structs orchestrent les dépendances injectées et n'ont pas besoin de leur propre interface.

// UserRepository est une interface pour la persistance.
type UserRepository interface {
	Save(name string) (int, error)
	FindByID(id int) (string, error)
}

// EmailSender est une interface pour l'envoi d'emails.
type EmailSender interface {
	Send(to, subject, body string) error
}

// ConsumerService est un consommateur - il utilise des interfaces injectées.
// PAS BESOIN D'INTERFACE car c'est un consommateur (orchestrateur).
type ConsumerService struct {
	repo   UserRepository
	mailer EmailSender
}

// NewConsumerService crée un nouveau ConsumerService.
//
// Params:
//   - repo: repository utilisateur
//   - mailer: service d'envoi d'emails
//
// Returns:
//   - *ConsumerService: nouvelle instance
func NewConsumerService(repo UserRepository, mailer EmailSender) *ConsumerService {
	// Retour de la nouvelle instance
	return &ConsumerService{repo: repo, mailer: mailer}
}

// CreateUser crée un utilisateur et envoie un email de bienvenue.
//
// Params:
//   - name: nom de l'utilisateur
//   - email: adresse email
//
// Returns:
//   - int: identifiant créé
//   - error: erreur éventuelle
func (s *ConsumerService) CreateUser(name, email string) (int, error) {
	// Sauvegarde via le repository injecté
	id, err := s.repo.Save(name)
	// Vérifier l'erreur
	if err != nil {
		// Retour avec erreur
		return 0, err
	}
	// Envoi email via le mailer injecté (ignore l'erreur)
	_ = s.mailer.Send(email, "Welcome", "Hello "+name)
	// Retour de l'ID
	return id, nil
}

// GetUser récupère un utilisateur par ID.
//
// Params:
//   - id: identifiant utilisateur
//
// Returns:
//   - string: nom de l'utilisateur
//   - error: erreur éventuelle
func (s *ConsumerService) GetUser(id int) (string, error) {
	// Récupération via le repository injecté
	return s.repo.FindByID(id)
}
