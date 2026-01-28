package gcpvision

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Note: This test will fail without proper Google Cloud credentials
	// It's primarily to verify the code compiles correctly
	t.Skip("Skipping test that requires Google Cloud credentials")

	_, err := NewClient()
	if err != nil {
		t.Logf("Expected error due to missing credentials: %v", err)
	}
}
