package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckAndUpdateCredentialName(t *testing.T) {
	// Test case 1: When credential name is already specified
	credential := Credential{
		Name:          "testname",
		Username:      "testuser",
		Password:      "testpassword",
		ServerAddress: "http://test-server",
	}
	checkAndUpdateCredentialName(&credential)
	require.Equal(t, "testname", credential.Name, "Credential name should remain unchanged")

	// Test case 2: When server address is a valid URL
	validURLCredential := Credential{
		Username:      "testuser",
		Password:      "testpassword",
		ServerAddress: "http://test-server",
	}
	checkAndUpdateCredentialName(&validURLCredential)
	require.Equal(t, "test-server-testuser", validURLCredential.Name, "Credential name not updated correctly")
}

func TestReadCredentialStore(t *testing.T) {
	// Test case 1: When config file does not exist
	_, err := readCredentialStore()
	require.Error(t, err, "Should return error when config file does not exist")

	configFile = t.TempDir() + "/config"

	// Test case 2: When config file exists
	sampleCredentialYaml := "current-credential-name: test-credential\ncredentials:\n- name: test-credential\n  username: testuser\n  password: testpassword\n  serveraddress: http://test-server"
	os.WriteFile(configFile, []byte(sampleCredentialYaml), 0644)

	credentialStore, err := readCredentialStore()
	require.NoError(t, err, "Should not return error when config file exists")
	require.Equal(t, 1, len(credentialStore.Credentials))
	require.Equal(t, credentialStore.CurrentCredentialName, "test-credential")
	require.Equal(t, credentialStore.Credentials[0], Credential{
		Name:          "test-credential",
		Username:      "testuser",
		Password:      "testpassword",
		ServerAddress: "http://test-server",
	})
}

func TestCheckAndCreateConfigFile(t *testing.T) {
	// Test case 1: When config file does not exist
	configFile = t.TempDir() + "/config"
	checkAndCreateConfigFile()
	_, err := os.Stat(configFile)
	require.NoError(t, err, "Config file should be created")

	// Test case 2: When config file already exists
	checkAndCreateConfigFile()
	_, err = os.Stat(configFile)
	require.NoError(t, err, "Config file should not be created again")
}

func TestStoreCredential(t *testing.T) {
	configFile = t.TempDir() + "/config"

	// Test case 1: Should store the credential and set it as the current credential
	credential := Credential{
		Name:          "test-credential-1",
		Username:      "testuser1",
		Password:      "testpassword1",
		ServerAddress: "http://test-server1",
	}
	StoreCredential(credential, true)
	credentialStore, err := readCredentialStore()
	require.NoError(t, err)
	require.Equal(t, "test-credential-1", credentialStore.CurrentCredentialName)
	require.Equal(t, 1, len(credentialStore.Credentials))
	require.Equal(t, credential, credentialStore.Credentials[0])

	// Test case 2: Should store the credential but not set it as the current credential
	credential = Credential{
		Name:          "test-credential-2",
		Username:      "testuser2",
		Password:      "testpassword2",
		ServerAddress: "http://test-server2",
	}
	StoreCredential(credential, false)
	credentialStore, err = readCredentialStore()
	require.NoError(t, err)
	require.Equal(t, "test-credential-1", credentialStore.CurrentCredentialName)
	require.Equal(t, 2, len(credentialStore.Credentials))
	require.Equal(t, credential, credentialStore.Credentials[1])

	// Test case 3: Should update the name and password of stored credentials
	credential = Credential{
		Name:          "new-name-1",
		Username:      "testuser1",
		Password:      "new-password-1",
		ServerAddress: "http://test-server1",
	}
	StoreCredential(credential, true)
	credentialStore, err = readCredentialStore()
	require.NoError(t, err)
	require.Equal(t, "new-name-1", credentialStore.CurrentCredentialName)
	require.Equal(t, 2, len(credentialStore.Credentials))
	require.Equal(t, credential, credentialStore.Credentials[1])
}

func TestResolveCredential(t *testing.T) {
	// Define the credential that we expect to be resolved
	requiredCrendential := Credential{
		Name:          "test-credential-1",
		Username:      "testuser1",
		Password:      "testpassword1",
		ServerAddress: "http://test-server1",
	}

	// Populate the config file with some sample credentials of two different users
	configFile = t.TempDir() + "/config"
	sampleCredentialYaml := "current-credential-name: test-credential-1\ncredentials:\n- name: test-credential-1\n  username: testuser1\n  password: testpassword1\n  serveraddress: http://test-server1 \n- name: test-credential-2\n  username: testuser2\n  password: testpassword2\n  serveraddress: http://test-server2"
	os.WriteFile(configFile, []byte(sampleCredentialYaml), 0644)

	// Test case 1: When credential name is specified
	resolvedCredential, err := resolveCredential("test-credential-1")
	require.NoError(t, err)
	require.Equal(t, requiredCrendential, resolvedCredential)

	// Test case 2: When credential name is not specified, we should fetch it from the environment variable
	os.Setenv("HARBOR_CREDENTIAL", "test-credential-1")
	resolvedCredential, err = resolveCredential("")
	require.NoError(t, err)
	require.Equal(t, requiredCrendential, resolvedCredential)

	// Test case 3: When credential name is not specified and environment variable is not set, we should use the current active credential
	os.Unsetenv("HARBOR_CREDENTIAL")
	resolvedCredential, err = resolveCredential("")
	require.NoError(t, err)
	require.Equal(t, requiredCrendential, resolvedCredential)

	// Test case 4: When all of the above are not specified, we should return an error
	os.WriteFile(configFile, []byte(""), 0644)
	resolvedCredential, err = resolveCredential("")
	require.Error(t, err)
	require.Equal(t, Credential{}, resolvedCredential)
}
