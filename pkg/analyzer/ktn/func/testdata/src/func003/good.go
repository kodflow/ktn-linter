package func003

// Good: Starts with Get
func GetUser() {}

// Good: Starts with Set
func SetValue() {}

// Good: Starts with Create
func CreateAccount() {}

// Good: Starts with Update
func UpdateRecord() {}

// Good: Starts with Delete
func DeleteFile() {}

// Good: Starts with Is
func IsValid() bool { return true }

// Good: Starts with Has
func HasPermission() bool { return true }

// Good: Starts with Can
func CanAccess() bool { return true }

// Good: Starts with Should
func ShouldRetry() bool { return true }

// Good: Starts with Handle
func HandleRequest() {}

// Good: Starts with Process
func ProcessData() {}

// Good: Starts with Validate
func ValidateInput() {}

// Good: Starts with Convert
func ConvertToJSON() {}

// Good: Starts with Calculate
func CalculateTotal() {}

// Good: unexported functions don't need to start with a verb
func helper() {}

func utilityFunction() {}

// Good: init and main are exempt
func init() {}

func main() {}

// Good: Test functions are exempt
func TestSomething() {}

func BenchmarkPerformance() {}
