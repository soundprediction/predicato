package gliner2

import (
	"testing"
)

func TestProvider(t *testing.T) {
	// Test provider constants
	if ProviderLocal != 0 {
		t.Errorf("Expected ProviderLocal to be 0, got %d", ProviderLocal)
	}
	if ProviderFastino != 1 {
		t.Errorf("Expected ProviderFastino to be 1, got %d", ProviderFastino)
	}
	if ProviderNative != 2 {
		t.Errorf("Expected ProviderNative to be 2, got %d", ProviderNative)
	}
}

func TestConfig(t *testing.T) {
	config := Config{
		Provider: ProviderLocal,
		Local: &LocalConfig{
			Endpoint: "http://localhost:8000",
			Timeout:  0,
		},
	}

	if config.Provider != ProviderLocal {
		t.Error("Expected provider to be ProviderLocal")
	}
	if config.Local.Endpoint != "http://localhost:8000" {
		t.Error("Expected endpoint to match")
	}
}

func TestNewClient(t *testing.T) {
	// Test Local provider creation
	localConfig := Config{
		Provider: ProviderLocal,
		Local: &LocalConfig{
			Endpoint: "http://localhost:8000",
			Timeout:  0,
		},
	}

	client, err := NewClient(localConfig)
	if err != nil {
		t.Fatalf("Failed to create local client: %v", err)
	}
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.GetCapabilities() == nil {
		t.Error("Expected capabilities to be set")
	}

	// Test cleanup
	if err := client.Close(); err != nil {
		t.Errorf("Expected clean close, got error: %v", err)
	}
}
