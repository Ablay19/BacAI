package serpapi

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")

	if client == nil {
		t.Error("Expected client to be created, got nil")
	}

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected apiKey to be 'test-api-key', got '%s'", client.apiKey)
	}

	if client.baseURL != "https://serpapi.com/search" {
		t.Errorf("Expected baseURL to be 'https://serpapi.com/search', got '%s'", client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be created, got nil")
	}
}
