package e2e

import (
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	var err error

	cds := root.RootCmd()
	err = cds.Execute()

	assert.NoError(t, err)
	cmd := root.LoginCommand()

	validServerAddresses := []string{
		"http://demo.goharbor.io:80",
		"https://demo.goharbor.io:443",
		"http://demo.goharbor.io",
		"https://demo.goharbor.io",
	}

	for _, serverAddress := range validServerAddresses {
		args := []string{
			serverAddress,
		}
		cmd.SetArgs(args)

		err = cmd.Flags().Set("name", "test")
		if err != nil {
			t.Fatal(err)
		}
		err = cmd.Flags().Set("username", "admin")
		if err != nil {
			t.Fatal(err)
		}
		err = cmd.Flags().Set("password", "Harbor12345")
		if err != nil {
			t.Fatal(err)
		}

		err = cmd.Execute()

		assert.NoError(t, err)
	}
}

func Test_Failure_WrongServer(t *testing.T) {
	cmd := root.LoginCommand()
	var err error

	args := []string{
		"demo.goharbor.io",
	}
	cmd.SetArgs(args)

	err = cmd.Flags().Set("name", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("username", "admin")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("password", "Harbor12345")
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Execute()

	assert.Error(t, err)
}

func Test_Failure_WrongUsername(t *testing.T) {
	cmd := root.LoginCommand()
	var err error

	args := []string{
		"http://demo.goharbor.io",
	}
	cmd.SetArgs(args)

	err = cmd.Flags().Set("name", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("username", "wrong")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("password", "Harbor12345")
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Execute()

	assert.Error(t, err)
}

func Test_Failure_WrongPassword(t *testing.T) {
	cmd := root.LoginCommand()
	var err error

	args := []string{
		"http://demo.goharbor.io",
	}
	cmd.SetArgs(args)

	err = cmd.Flags().Set("name", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("username", "admin")
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Flags().Set("password", "wrong")
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Execute()

	assert.Error(t, err)
}
