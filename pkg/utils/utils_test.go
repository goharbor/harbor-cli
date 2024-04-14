package utils

import (
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
)

func TestGetClientByConfig(t *testing.T) {
	clientConfig := &harbor.ClientSetConfig{
		URL:      "testURL",
		Username: "testUsername",
		Password: "testPassword",
	}
	client := GetClientByConfig(clientConfig)
	if client == nil {
		t.Errorf("Expected client not to be nil")
	}
}

func TestGetClientByCredentialName(t *testing.T) {
	client := GetClientByCredentialName("127.0.0.1-testuser")
	if client == nil {
		t.Errorf("Expected client not to be nil")
	}
}

func TestPrintPayloadInJSONFormat(t *testing.T) {
	PrintPayloadInJSONFormat(nil)
	
	payload := map[string]interface{}{
		"testkey": "test-value",
	}
	PrintPayloadInJSONFormat(payload)
}