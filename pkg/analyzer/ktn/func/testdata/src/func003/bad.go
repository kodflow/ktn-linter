package func003

// Bad: Starts with a noun
func User() {} // want "KTN-FUNC-003"

// Bad: Starts with adjective
func ValidData() {} // want "KTN-FUNC-003"

// Bad: Generic name without verb
func Handler() {} // want "KTN-FUNC-003"

// Bad: Name without verb prefix
func TotalAmount() {} // want "KTN-FUNC-003"

// Bad: Property-like name
func CurrentUser() {} // want "KTN-FUNC-003"

// Bad: Noun-based name
func Database() {} // want "KTN-FUNC-003"

// Bad: Configuration without verb
func Settings() {} // want "KTN-FUNC-003"

// Bad: Result without verb
func Sum() {} // want "KTN-FUNC-003"

// Bad: Single word noun
func Data() {} // want "KTN-FUNC-003"

// Bad: Single word non-verb
func Result() {} // want "KTN-FUNC-003"

// Bad: Single word adjective
func Valid() {} // want "KTN-FUNC-003"

// Bad: Single word generic name
func Helper() {} // want "KTN-FUNC-003"
