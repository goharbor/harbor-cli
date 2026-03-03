package oauth

import (
	"testing"
	"time"
)

func TestOAuthConfigCreation(t *testing.T) {
	config := &OAuthConfig{
		IssuerURL: "https://example.com",
		Scopes:    []string{"openid", "profile"},
	}

	if config.IssuerURL != "https://example.com" {
		t.Errorf("Expected issuer URL to be https://example.com, got %s", config.IssuerURL)
	}

	if len(config.Scopes) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(config.Scopes))
	}
}

func TestTokenResponseExpiry(t *testing.T) {
	token := &TokenResponse{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	if token.AccessToken != "test-token" {
		t.Errorf("Expected access token to be test-token, got %s", token.AccessToken)
	}

	if token.Expiry.Before(time.Now()) {
		t.Errorf("Token expiry should be in the future, got %s", token.Expiry)
	}
}

func TestOAuthErrorWithDescription(t *testing.T) {
	err := &OAuthError{
		Code:        "invalid_grant",
		Description: "The provided authorization grant is invalid",
	}

	expected := "invalid_grant: The provided authorization grant is invalid"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestOAuthErrorWithoutDescription(t *testing.T) {
	err := &OAuthError{
		Code: "access_denied",
	}

	if err.Error() != "access_denied" {
		t.Errorf("Expected error message 'access_denied', got '%s'", err.Error())
	}
}
