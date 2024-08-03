package project

import (
	"bytes"
	"encoding/json"
	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to read JSON file
func readJSONFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// CaptureOutput : captures output printed to stdout via `PrintPayloadInJSONFormat`
func CaptureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old // restoring the real stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestViewCommandOutput(t *testing.T) {
	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	// mock the test-config file
	testConfigFile := filepath.Join(workingDir, "test_config.yaml")
	if _, err := os.Stat(testConfigFile); os.IsNotExist(err) {
		t.Fatalf("Test configuration file does not exist: %s", testConfigFile)
	}

	os.Setenv("HARBOR_CLI_CONFIG", testConfigFile)
	defer os.Unsetenv("HARBOR_CLI_CONFIG")

	viper.SetConfigFile(testConfigFile)
	err = viper.ReadInConfig()
	assert.NoError(t, err)

	// Construct the absolute path to the expected JSON output file
	expectedFilePath := filepath.Join(workingDir, "view_test.json")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected JSON output file does not exist: %s", expectedFilePath)
	}

	// Read the expected JSON output from file
	expectedOutput, err := readJSONFile(expectedFilePath)
	assert.NoError(t, err)

	// Unmarshal the expected JSON
	var expectedJSON map[string]interface{}
	err = json.Unmarshal(expectedOutput, &expectedJSON)
	assert.NoError(t, err)
	
	rootCmd := root.RootCmd()
	rootCmd.SetArgs([]string{"project", "view", "my-project"})

	// Capture the output
	output := CaptureOutput(func() {
		err = rootCmd.Execute()
		assert.NoError(t, err)
	})

	// Unmarshal the captured output
	var actualJSON map[string]interface{}
	err = json.Unmarshal([]byte(output), &actualJSON)
	assert.NoError(t, err)

	assert.Equal(t, expectedJSON, actualJSON, "The output JSON does not match the expected JSON")
}
