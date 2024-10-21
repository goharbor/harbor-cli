package e2e

import (
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/stretchr/testify/assert"
)

func initialize(t *testing.T) {
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
}

func Test_Login_Success(t *testing.T) {
	initialize(t) // Initialize the root command

	cmd := root.LoginCommand()
	validServerAddresses := []string{
		"http://demo.goharbor.io:80",
		"https://demo.goharbor.io:443",
		"http://demo.goharbor.io",
		"https://demo.goharbor.io",
		// "demo.goharbor.io",
	}

	for _, serverAddress := range validServerAddresses {
		t.Run("ValidServer_"+serverAddress, func(t *testing.T) {
			args := []string{serverAddress}
			cmd.SetArgs(args)

			assert.NoError(t, cmd.Flags().Set("name", "test"))
			assert.NoError(t, cmd.Flags().Set("username", "admin"))
			assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

			err := cmd.Execute()
			assert.NoError(t, err, "Expected no error for server: %s", serverAddress)
		})
	}
}

func Test_Login_Failure_WrongServer(t *testing.T) {
	cmd := root.LoginCommand()
	args := []string{"wrongserver"}
	cmd.SetArgs(args)

	assert.NoError(t, cmd.Flags().Set("name", "test"))
	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for invalid server")
}

func Test_Login_Failure_WrongUsername(t *testing.T) {
	cmd := root.LoginCommand()
	args := []string{"http://demo.goharbor.io"}
	cmd.SetArgs(args)

	assert.NoError(t, cmd.Flags().Set("name", "test"))
	assert.NoError(t, cmd.Flags().Set("username", "wrong"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong username")
}

func Test_Login_Failure_WrongPassword(t *testing.T) {
	cmd := root.LoginCommand()
	args := []string{"http://demo.goharbor.io"}
	cmd.SetArgs(args)

	assert.NoError(t, cmd.Flags().Set("name", "test"))
	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "wrong"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong password")
}
