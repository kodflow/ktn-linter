// Package gospec_bad_naming montre des conventions de nommage non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_naming

import "fmt"

// ❌ BAD PRACTICE: Using snake_case instead of MixedCaps
var user_count int
var max_retry_count int

const default_timeout = 30
const api_version = "v1"

func get_user_name() string { return "" }
func calculate_total_price() int { return 0 }

type user_data struct {
	first_name string
	last_name  string
}

// ❌ BAD PRACTICE: Using ALL_CAPS for non-constants (not Go convention)
var MAX_SIZE = 100 // Devrait être: MaxSize ou maxSize

// ❌ BAD PRACTICE: Redundant package name in identifier
type UserUser struct{} // package user + User = UserUser (redondant)

func UserGetUser() {} // Redondant si dans package user

// ❌ BAD PRACTICE: Stutter in naming (repetitive names)
type ConfigConfig struct {
	ConfigValue string
}

// ❌ BAD PRACTICE: Using generic/unclear names
func DoStuff() {}      // Trop vague
func Process1() {}     // Que traite-t-on?
func Handle() {}       // Que gère-t-on?
func Manager1() {}     // Que gère-t-on?
func Data() {}         // Quelles données?
func Info() string { return "" }

var result int // Trop générique au niveau package
var data []byte
var temp string

// ❌ BAD PRACTICE: Using Hungarian notation
var strName string  // 'str' prefix inutile
var iCount int      // 'i' prefix inutile
var bIsValid bool   // 'b' prefix inutile
var ptrUser *string // 'ptr' prefix inutile

// ❌ BAD PRACTICE: Overly long names without clear benefit
func GetUserAccountInformationFromDatabaseByUniqueIdentifier() {}

var thisIsAVeryLongVariableNameThatDoesNotAddMuchValue int

// ❌ BAD PRACTICE: Using single letter names where clarity needed
func ProcessItem(d []byte, c int, f bool) error { // d? c? f? peu clair
	_, _, _ = d, c, f
	return nil
}

// ❌ BAD PRACTICE: Not using conventional names for common types
func ProcessString(s string) { // 'data' serait plus clair que 's' ici
	fmt.Println(s)
}

// ❌ BAD PRACTICE: Using 'get' prefix unnecessarily
func GetName() string { return "" }  // En Go, souvent simplement Name()
func GetCount() int { return 0 }     // Souvent simplement Count()

// ❌ BAD PRACTICE: Method names don't follow receiver naming
type Handler struct{}

// Receiver devrait être cohérent (h ou hand, pas les deux)
func (handler Handler) HandleProcess() {}
func (h Handler) Execute() {}
func (hdlr Handler) Run() {}

// ❌ BAD PRACTICE: Interface names without 'er' suffix (when single method)
type Processor interface { // Devrait être nommé différemment ou utiliser convention
	Process()
}

type Manager2 interface { // Trop vague
	Manage()
}

// ❌ BAD PRACTICE: Using abbreviations inconsistently
var usrCfg int  // usr abrégé
var userConfig int // user complet - incohérent
var cfgMgr int     // double abréviation

// ❌ BAD PRACTICE: Not using conventional abbreviations
var configuration int // 'config' ou 'cfg' est conventionnel
var identification string // 'id' est conventionnel
var synchronization bool // 'sync' est conventionnel

// ❌ BAD PRACTICE: Exported names without clear purpose
var Flag bool        // Trop vague pour export
var Value int        // Trop générique pour export
var Status string    // Status de quoi?

// ❌ BAD PRACTICE: Using 'is', 'has' inconsistently with booleans
var EnableFeature bool  // Devrait être: FeatureEnabled
var isActive bool       // Incohérent (minuscule après 'is')
var HasPermission bool  // OK mais devrait suivre convention

// ❌ BAD PRACTICE: Using verb prefixes for non-boolean values
var GetResult int    // Get implique une action, pas un état
var ProcessOutput string // Process implique une action

// ❌ BAD PRACTICE: Package name doesn't match directory or is plural
// Si le répertoire est 'user', le package devrait être 'user', pas 'users'

// ❌ BAD PRACTICE: Using underscores in package names
// package user_service // Devrait être: package userservice

// ❌ BAD PRACTICE: Not using conventional receiver names
type Service struct{}

// Receiver devrait être 's', pas 'service', 'srv', 'this', 'self'
func (service Service) Method1() {}
func (this Service) Method2() {}
func (self Service) Method3() {}

// ❌ BAD PRACTICE: Inconsistent acronym casing
var HttpClient int  // Devrait être: HTTPClient (acronymes en majuscules)
var XmlParser int   // Devrait être: XMLParser
var JsonData int    // Devrait être: JSONData
var IdValue int     // Devrait être: IDValue

// ❌ BAD PRACTICE: Using 'new' or 'make' in function names (reserved feel)
func NewCreate() {} // Redondant
func MakeData() {}  // Confus avec built-in make

// ❌ BAD PRACTICE: Error variable not starting with 'Err'
var NotFoundError = fmt.Errorf("not found")      // Devrait être: ErrNotFound
var InvalidInputError = fmt.Errorf("invalid")    // Devrait être: ErrInvalidInput
