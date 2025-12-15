// Package api001 contains test cases for KTN-API-001.
package api001

import (
	"net/http"
	"os"
)

// HTTPClient is an alias to test alias handling.
type HTTPClient = http.Client

// badHTTPClientWithOneMethod uses external concrete type with one method call.
// This should trigger KTN-API-001.
func badHTTPClientWithOneMethod(client *http.Client) (*http.Response, error) { // want `KTN-API-001:.*client.*http\.Client.*client.*Do`
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	return client.Do(req)
}

// badHTTPClientWithMultipleMethods uses external concrete type with multiple method calls.
// This should trigger KTN-API-001.
func badHTTPClientWithMultipleMethods(client *http.Client) (*http.Response, error) { // want `KTN-API-001:.*client.*http\.Client.*client.*Get.*Head`
	_, err := client.Get("http://example.com")
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return client.Head("http://example.com")
}

// badFileWithMethods uses os.File with method calls.
// This should trigger KTN-API-001.
func badFileWithMethods(f *os.File) ([]byte, error) { // want `KTN-API-001:.*f.*os\.File.*file.*Read`
	buf := make([]byte, 100)
	_, err := f.Read(buf)
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// badMultipleParams tests multiple external concrete params with method calls.
// Both should trigger KTN-API-001.
func badMultipleParams(client *http.Client, f *os.File) error { // want `KTN-API-001:.*client.*http\.Client.*client.*Get` `KTN-API-001:.*f.*os\.File.*file.*Close`
	_, err := client.Get("http://example.com")
	// Vérification de la condition
	if err != nil {
		return err
	}
	return f.Close()
}

// badAliasExternalType uses a type alias to external type.
// Should still trigger KTN-API-001 because alias resolves to http.Client.
func badAliasExternalType(c *HTTPClient) (*http.Response, error) { // want `KTN-API-001:.*c.*http\.Client.*client.*Do`
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	return c.Do(req)
}
