// Bad examples for the var018 test case.
package var018

// Constants to avoid magic numbers
const (
	// BadPortValue is the port value
	BadPortValue int = 8080
	// BadMaxValue is the max value
	BadMaxValue int = 100
	// BadTimeoutValue is the timeout value
	BadTimeoutValue int = 30
)

// Bad: Variables using snake_case (with underscores)

var (
	// http_client uses snake_case
	http_client string = "client" // want "KTN-VAR-018"

	// server_port uses snake_case
	server_port int = BadPortValue // want "KTN-VAR-018"

	// max_connections uses snake_case
	max_connections int = BadMaxValue // want "KTN-VAR-018"

	// api_key uses snake_case
	api_key string = "secret" // want "KTN-VAR-018"

	// user_name uses snake_case
	user_name string = "admin" // want "KTN-VAR-018"

	// is_enabled uses snake_case
	is_enabled bool = true // want "KTN-VAR-018"

	// default_timeout uses snake_case
	default_timeout int = BadTimeoutValue // want "KTN-VAR-018"

	// my_var_name uses snake_case with multiple underscores
	my_var_name string = "test" // want "KTN-VAR-018"
)

// init uses the variables to avoid compilation errors
func init() {
	_ = http_client
	_ = server_port
	_ = max_connections
	_ = api_key
	_ = user_name
	_ = is_enabled
	_ = default_timeout
	_ = my_var_name
}
