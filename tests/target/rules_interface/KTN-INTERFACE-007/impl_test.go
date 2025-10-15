package goodinterfaces_test

import "testing"

func TestNewService(t *testing.T) {
	svc := NewService("test-service")
	if svc == nil {
		t.Error("NewService() returned nil")
	}
}

func TestNewHelper(t *testing.T) {
	helper := NewHelper()
	if helper == nil {
		t.Error("NewHelper() returned nil")
	}
}
