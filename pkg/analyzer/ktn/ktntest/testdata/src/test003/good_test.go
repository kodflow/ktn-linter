package test003

import "testing"

// TestNewGoodResource teste le constructeur
func TestNewGoodResource(t *testing.T) {
	r := NewGoodResource()
	if r == nil {
		t.Error("Expected non-nil resource")
	}
}

// TestGoodResource_Metadata teste la méthode Metadata
func TestGoodResource_Metadata(t *testing.T) {
	r := NewGoodResource()
	result := r.Metadata()
	if result != "good_resource" {
		t.Errorf("Expected 'good_resource', got '%s'", result)
	}
}

// TestGoodResource_Schema teste la méthode Schema
func TestGoodResource_Schema(t *testing.T) {
	r := NewGoodResource()
	schema := r.Schema()
	if len(schema) == 0 {
		t.Error("Expected non-empty schema")
	}
}

// TestGoodResource_Configure teste la méthode Configure
func TestGoodResource_Configure(t *testing.T) {
	r := NewGoodResource()
	err := r.Configure("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
