// Good examples for the func008 test case.
package func008

const (
	// DOUBLE_MULTIPLIER représente le multiplicateur pour doubler une valeur.
	DOUBLE_MULTIPLIER int = 2
	// ARRAY_SIZE représente la taille du tableau local.
	ARRAY_SIZE int = 5
	// INITIAL_COUNT représente le compteur initial.
	INITIAL_COUNT int = 10
	// SLICE_CAPACITY représente la capacité du slice.
	SLICE_CAPACITY int = 10
)

// MyStruct représente une structure de test pour les getters.
// Elle fournit des méthodes pour accéder et modifier ses champs privés.
type MyStruct struct {
	value int
	name  string
}

// MyStructInterface définit les méthodes publiques de MyStruct.
type MyStructInterface interface {
	GetValue() int
	IsValid() bool
	HasName() bool
	GetDoubleValue() int
	SetValue(v int)
	UpdateValue(v int)
	GetProcessed() []int
	GetMap() map[string]int
	GetLocalArray() []int
	GetLocalMapValue() int
	GetIncrementedValue() int
	GetDecrementedValue() int
}

// NewMyStruct crée une nouvelle instance de MyStruct.
//
// Params:
//   - value: valeur initiale
//   - name: nom initial
//
// Returns:
//   - *MyStruct: nouvelle instance
func NewMyStruct(value int, name string) *MyStruct {
	// Retour de la nouvelle instance
	return &MyStruct{
		value: value,
		name:  name,
	}
}

// GetValue retourne la valeur du champ value.
// Returns:
//   - int: valeur actuelle
func (m *MyStruct) GetValue() int {
	// Retourne la valeur du champ
	return m.value
}

// IsValid vérifie si la valeur est positive.
// Returns:
//   - bool: true si valeur > 0
func (m *MyStruct) IsValid() bool {
	// Vérifie si la valeur est strictement positive
	if m.value > 0 {
		// Valeur positive
		return true
	}
	// Valeur négative ou nulle
	return false
}

// HasName vérifie si le nom n'est pas vide.
// Returns:
//   - bool: true si nom non vide
func (m *MyStruct) HasName() bool {
	// Vérifie si le nom n'est pas une chaîne vide
	if m.name != "" {
		// Nom non vide
		return true
	}
	// Nom vide
	return false
}

// GetDoubleValue retourne le double de la valeur.
// Returns:
//   - int: valeur multipliée par 2
func (m *MyStruct) GetDoubleValue() int {
	// Calcule le double de la valeur
	result := m.value * DOUBLE_MULTIPLIER
	// Retourne le résultat
	return result
}

// SetValue définit une nouvelle valeur.
// Params:
//   - v: nouvelle valeur à assigner
func (m *MyStruct) SetValue(v int) {
	// Assigne la nouvelle valeur
	m.value = v
}

// UpdateValue met à jour la valeur.
// Params:
//   - v: nouvelle valeur à assigner
func (m *MyStruct) UpdateValue(v int) {
	// Met à jour la valeur du champ
	m.value = v
}

// GetProcessed retourne un tableau avec la valeur traitée.
// Returns:
//   - []int: tableau contenant la valeur
func (m *MyStruct) GetProcessed() []int {
	// Crée un tableau local
	local := [SLICE_CAPACITY]int{}
	// Assigne la valeur au premier élément
	local[0] = m.value
	// Retourne le tableau local
	return local[:]
}

// GetMap retourne une map contenant la valeur.
// Returns:
//   - map[string]int: map avec la valeur
func (m *MyStruct) GetMap() map[string]int {
	// Crée une map locale
	result := make(map[string]int, 1)
	// Stocke la valeur dans la map
	result["value"] = m.value
	// Retourne la map
	return result
}

// TestGetValue fonction de test avec effet de bord.
// Params:
//   - m: structure à tester
//
// Returns:
//   - int: valeur incrémentée
func TestGetValue(m *MyStruct) int {
	// Incrémente la valeur pour le test
	m.value++
	// Retourne la nouvelle valeur
	return m.value
}

// BenchmarkGetValue fonction de benchmark avec effet de bord.
// Params:
//   - m: structure à benchmarker
//
// Returns:
//   - int: valeur incrémentée
func BenchmarkGetValue(m *MyStruct) int {
	// Incrémente la valeur pour le benchmark
	m.value++
	// Retourne la nouvelle valeur
	return m.value
}

// GetLocalArray retourne un tableau avec les valeurs calculées.
// Returns:
//   - []int: tableau de valeurs
func (m *MyStruct) GetLocalArray() []int {
	// Crée un tableau local de 5 éléments
	arr := [ARRAY_SIZE]int{}
	// Stocke la valeur actuelle
	arr[0] = m.value
	// Stocke le double de la valeur
	arr[1] = m.value * DOUBLE_MULTIPLIER
	// Retourne le tableau
	return arr[:]
}

// GetLocalMapValue retourne une valeur depuis une map locale.
// Returns:
//   - int: valeur extraite de la map
func (m *MyStruct) GetLocalMapValue() int {
	// Crée une map locale
	localMap := make(map[string]int, 1)
	// Stocke la valeur dans la map
	localMap["key"] = m.value
	// Retourne la valeur depuis la map
	return localMap["key"]
}

// DataReader interface pour la lecture de données.
type DataReader interface {
	GetData() string
	IsReady() bool
	HasItems() bool
}

// dataReaderImpl implémente DataReader.
type dataReaderImpl struct {
	data  string
	ready bool
}

// GetData retourne les données.
//
// Returns:
//   - string: les données stockées
func (d *dataReaderImpl) GetData() string {
	// Retourne les données
	return d.data
}

// IsReady indique si prêt.
//
// Returns:
//   - bool: état de préparation
func (d *dataReaderImpl) IsReady() bool {
	// Retourne l'état ready
	return d.ready
}

// HasItems vérifie si des items existent.
//
// Returns:
//   - bool: true si données non vides
func (d *dataReaderImpl) HasItems() bool {
	// Vérifie si la chaîne n'est pas vide
	if d.data != "" {
		// Données présentes
		return true
	}
	// Pas de données
	return false
}

// useDataReader utilise l'interface DataReader.
//
// Params:
//   - dr: interface de lecture de données
//
// Returns:
//   - string: données lues
func useDataReader(dr DataReader) string {
	// Vérifie si prêt
	if dr.IsReady() && dr.HasItems() {
		// Retourne les données
		return dr.GetData()
	}
	// Retourne chaîne vide
	return ""
}

// GetIncrementedValue retourne la valeur incrémentée localement.
// Returns:
//   - int: valeur + 1
func (m *MyStruct) GetIncrementedValue() int {
	// Copie la valeur dans une variable locale
	local := m.value
	// Incrémente la variable locale
	local++
	// Retourne la valeur incrémentée
	return local
}

// GetDecrementedValue retourne une valeur décrémentée.
// Returns:
//   - int: constante décrémentée
func (m *MyStruct) GetDecrementedValue() int {
	// Initialise un compteur
	count := INITIAL_COUNT
	// Décrémente le compteur
	count--
	// Retourne le compteur décrémenté
	return count
}

// init utilise les fonctions privées
func init() {
	// Appel de useDataReader
	_ = useDataReader(DataReader{})
}
