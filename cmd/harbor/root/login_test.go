package root

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/views/login"
	"github.com/stretchr/testify/assert"
)

func TestRunLogin_Successful(t *testing.T) {
	// Mock Harbor server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v2.0/users/current", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()
	opts := login.LoginView{
		Server:   mockServer.URL,
		Username: "testuser",
		Password: "testpassword",
	}
	err := runLogin(opts)
	assert.NoError(t, err)
}

func TestRunLogin_Failed(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v2.0/users/current", r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer mockServer.Close()
	opts := login.LoginView{
		Server:   mockServer.URL,
		Username: "testuser",
		Password: "testpassword",
	}
	err := runLogin(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "login failed")
}

func TestLoginCommand(t *testing.T) {
	cmd := LoginCommand()
	flags := cmd.Flags()

	// Test case 1: required flags value should be set
	cmd.SetArgs([]string{"http://testgoharbor.io"})
	flags.Set("username", "testuser")
	flags.Set("password", "testpass")

	assert.NotEmpty(t, cmd.Args, "serverAddress cmd argument has already set")
	usernameFlag, err := flags.GetString("username")
	assert.NoError(t, err, "Expected no error getting username flag")
	assert.NotEmpty(t, usernameFlag, "Expected username flag has already set")
	passwordFlag, err := flags.GetString("password")
	assert.NoError(t, err, "Expected no error getting password flag")
	assert.NotEmpty(t, passwordFlag, "Expected password flag has already set")

	// Test case 2: required flags value should not be set
	flags.Set("username", "")
	flags.Set("password", "")
	usernameFlag, _ = flags.GetString("username")
	if usernameFlag == "" {
		assert.Empty(t, usernameFlag, "Expected username flag has to be set")
	}
	passwordFlag, _ = flags.GetString("password")
	if passwordFlag == "" {
		assert.Empty(t, passwordFlag, "Expected password flag has to be set")
	}
}
