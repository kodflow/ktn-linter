// Package api001 contains test cases for KTN-API-001.
package api001

import (
	"bytes"
	"net/http"
	"os"
)

// badHTTPClientWithOneMethod uses external concrete type with one method call.
// This should trigger KTN-API-001.
func badHTTPClientWithOneMethod(client *http.Client) (*http.Response, error) { // want "KTN-API-001"
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	return client.Do(req)
}

// badHTTPClientWithMultipleMethods uses external concrete type with multiple method calls.
// This should trigger KTN-API-001.
func badHTTPClientWithMultipleMethods(client *http.Client) (*http.Response, error) { // want "KTN-API-001"
	_, err := client.Get("http://example.com")
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return client.Head("http://example.com")
}

// badFileWithMethods uses os.File with method calls.
// This should trigger KTN-API-001.
func badFileWithMethods(f *os.File) ([]byte, error) { // want "KTN-API-001"
	buf := make([]byte, 100)
	_, err := f.Read(buf)
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// badBufferWithMethods uses bytes.Buffer with method calls.
// This should trigger KTN-API-001.
func badBufferWithMethods(buf *bytes.Buffer) string { // want "KTN-API-001"
	return buf.String()
}

// badMultipleParams tests multiple external concrete params with method calls.
// Both should trigger KTN-API-001.
func badMultipleParams(client *http.Client, f *os.File) error { // want "KTN-API-001" "KTN-API-001"
	_, err := client.Get("http://example.com")
	// Vérification de la condition
	if err != nil {
		return err
	}
	return f.Close()
}
