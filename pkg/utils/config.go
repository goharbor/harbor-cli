package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Credential struct {
	Name          string `yaml:"name"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	ServerAddress string `yaml:"serveraddress"`
}

type HarborConfig struct {
	CurrentCredentialName string       `mapstructure:"current-credential-name" yaml:"current-credential-name"`
	Credentials           []Credential `mapstructure:"credentials" yaml:"credentials"`
}

type HarborData struct {
	ConfigPath string `mapstructure:"configpath" yaml:"configpath"`
}

type Once struct {
	once sync.Once
}

func (o *Once) Do(f func()) {
	o.once.Do(f)
}

func (o *Once) Reset() {
	o.once = sync.Once{}
}

var (
	CurrentHarborData   *HarborData
	CurrentHarborConfig *HarborConfig
	configMutex         sync.RWMutex
	configInitError     error
)

var ConfigInitialization = &Once{}

func InitConfig(cfgFile string, userSpecifiedConfig bool) {
	ConfigInitialization.Do(func() {
		harborDataPath, harborDataDir := GetDataPaths()
		harborConfigPath, err := DetermineConfigPath(cfgFile, userSpecifiedConfig)
		if err != nil {
			configInitError = err
			log.Fatalf("%v", err)
		}

		// Ensure data directory exists
		if err := os.MkdirAll(harborDataDir, os.ModePerm); err != nil {
			configInitError = fmt.Errorf("failed to create data directory: %w", err)
			log.Fatalf("%v", configInitError)
		}

		// Update or create data file
		if err := ApplyDataFile(harborDataPath, harborConfigPath); err != nil {
			configInitError = err
			log.Fatalf("%v", err)
		}

		// Ensure config file exists
		if err := EnsureConfigFileExists(harborConfigPath); err != nil {
			configInitError = err
			log.Fatalf("%v", err)
		}

		// Read and unmarshal the config file
		v, err := ReadConfig(harborConfigPath)
		if err != nil {
			configInitError = err
			log.Fatalf("%v", err)
		}

		var harborConfig HarborConfig
		if err := v.Unmarshal(&harborConfig); err != nil {
			configInitError = fmt.Errorf("failed to unmarshal config file: %w", err)
			log.Fatalf("%v", configInitError)
		}

		configMutex.Lock()
		defer configMutex.Unlock()
		CurrentHarborConfig = &harborConfig
		CurrentHarborData = &HarborData{ConfigPath: harborConfigPath}
	})
}

// Helper function to get data paths
func GetDataPaths() (harborDataPath string, harborDataDir string) {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to determine user home directory: %v", err)
		}
		xdgDataHome = filepath.Join(home, ".local", "share")
	}
	harborDataDir = filepath.Join(xdgDataHome, "harbor-cli")
	harborDataPath = filepath.Join(harborDataDir, "data.yaml")
	return
}

// Helper function to determine the config path
func DetermineConfigPath(cfgFile string, userSpecifiedConfig bool) (string, error) {
	var harborConfigPath string
	var err error

	// 1. Check if user specified --config
	if userSpecifiedConfig && cfgFile != "" {
		harborConfigPath, err = filepath.Abs(cfgFile)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path for config file: %w", err)
		}
		return harborConfigPath, nil
	}
	// 2. Check HARBOR_CLI_CONFIG environment variable
	harborEnvVar := os.Getenv("HARBOR_CLI_CONFIG")
	if harborEnvVar != "" {
		harborConfigPath, err = filepath.Abs(harborEnvVar)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path for config file from HARBOR_CLI_CONFIG: %w", err)
		}
		return harborConfigPath, nil
	}

	// 3. Use default XDG config path
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to determine user home directory: %w", err)
		}
		xdgConfigHome = filepath.Join(home, ".config")
	}
	harborConfigPath, err = filepath.Abs(filepath.Join(xdgConfigHome, "harbor-cli", "config.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for default config file: %w", err)
	}
	return harborConfigPath, nil
}

// Helper function to ensure config file exists
func EnsureConfigFileExists(harborConfigPath string) error {
	// Ensure parent directory exists
	configDir := filepath.Dir(harborConfigPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config file if it doesn't exist
	if _, err := os.Stat(harborConfigPath); os.IsNotExist(err) {
		if err := CreateConfigFile(harborConfigPath); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}
	return nil
}

// Helper function to read the config file using Viper
func ReadConfig(harborConfigPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(harborConfigPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w. Please ensure the config file exists.", err)
	}
	return v, nil
}

func GetCurrentHarborConfig() (*HarborConfig, error) {
	ConfigInitialization.Do(func() {
		// No action needed; InitConfig should have been called before
	})

	if configInitError != nil {
		return nil, fmt.Errorf("initialization error: %w", configInitError)
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	if CurrentHarborConfig == nil {
		return nil, errors.New("configuration is not yet initialized")
	}

	return CurrentHarborConfig, nil
}

func GetCurrentHarborData() (*HarborData, error) {
	ConfigInitialization.Do(func() {
		// No action needed; initialization should have been called before
	})

	if configInitError != nil {
		return nil, fmt.Errorf("initialization error: %w", configInitError)
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	if CurrentHarborData == nil {
		return nil, errors.New("HarborData is not yet initialized")
	}

	return CurrentHarborData, nil
}

func CreateDataFile(dataFilePath string, initialConfigPath string) error {
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		dataDir := filepath.Dir(dataFilePath)
		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}

		absConfigPath, err := filepath.Abs(initialConfigPath)
		if err != nil {
			log.Fatalf("Failed to resolve absolute path for config file: %v", err)
		}

		dataFile := HarborData{
			ConfigPath: absConfigPath,
		}

		v := viper.New()
		v.SetConfigType("yaml")
		v.Set("configPath", dataFile.ConfigPath)

		if err := v.WriteConfigAs(dataFilePath); err != nil {
			log.Fatalf("Failed to write data file: %v", err)
		}

		log.Infof("Data file created at %s with configPath: %s", dataFilePath, dataFile.ConfigPath)
	} else if err != nil {
		log.Fatalf("Error checking data file: %v", err)
	}

	return nil
}

func ReadDataFile(dataFilePath string) (HarborData, error) {
	var dataFile HarborData

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(dataFilePath)

	if err := v.ReadInConfig(); err != nil {
		return dataFile, fmt.Errorf("failed to read data file: %v", err)
	}

	dataFile.ConfigPath = v.GetString("configPath")

	return dataFile, nil
}

// Helper function to update or create the data file
func ApplyDataFile(harborDataPath, harborConfigPath string) error {
	dataFileContent, err := ReadDataFile(harborDataPath)
	if err == nil {
		if dataFileContent.ConfigPath != harborConfigPath {
			if err := UpdateDataFile(harborDataPath, harborConfigPath); err != nil {
				return fmt.Errorf("failed to update data file: %w", err)
			}
		} else if dataFileContent.ConfigPath == "" {
			if err := CreateDataFile(harborDataPath, harborConfigPath); err != nil {
				return fmt.Errorf("failed to create data file: %w", err)
			}
		} else {
			log.Debugf("Data file already exists with the same config path: %s", harborConfigPath)
		}
	} else {
		// Data file does not exist, create it
		if err := CreateDataFile(harborDataPath, harborConfigPath); err != nil {
			return fmt.Errorf("failed to create data file: %w", err)
		}
	}
	return nil
}

func UpdateDataFile(dataFilePath string, newConfigPath string) error {
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		log.Fatalf("data file does not exist at %s", dataFilePath)
	} else if err != nil {
		log.Fatalf("error checking data file: %v", err)
	}

	absConfigPath, err := filepath.Abs(newConfigPath)
	if err != nil {
		log.Fatalf("failed to resolve absolute path for new config file: %v", err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(dataFilePath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read existing data file: %v", err)
	}

	v.Set("configPath", absConfigPath)

	if err := v.WriteConfig(); err != nil {
		log.Fatalf("failed to write updated data file: %v", err)
	}

	log.Infof("Data file at %s updated with new configPath: %s", dataFilePath, absConfigPath)
	return nil
}

func CreateConfigFile(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			log.Fatalf("failed to create config directory: %v", err)
		}

		v := viper.New()
		v.SetConfigType("yaml")

		defaultConfig := HarborConfig{
			CurrentCredentialName: "",
			Credentials:           []Credential{},
		}

		v.Set("current-credential-name", defaultConfig.CurrentCredentialName)
		v.Set("credentials", defaultConfig.Credentials)

		if err := v.WriteConfigAs(configPath); err != nil {
			log.Fatalf("failed to write config file: %v", err)
		}

		log.Infof("Config file created at %s", configPath)
	} else if err != nil {
		log.Fatalf("error checking config file: %v", err)
	}

	return nil
}

func GetCredentials(credentialName string) (Credential, error) {
	currentConfig, err := GetCurrentHarborConfig()
	if err != nil {
		return Credential{}, fmt.Errorf("failed to get current Harbor configuration: %w", err)
	}

	if currentConfig == nil || currentConfig.Credentials == nil {
		return Credential{}, errors.New("current Harbor configuration or credentials are not initialized")
	}

	for _, cred := range currentConfig.Credentials {
		if cred.Name == credentialName {
			return cred, nil
		}
	}

	return Credential{}, fmt.Errorf("credential with name '%s' not found", credentialName)
}

func AddCredentialsToConfigFile(credential Credential, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at %s", configPath)
	} else if err != nil {
		log.Fatalf("error checking config file: %v", err)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var c HarborConfig
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	c.Credentials = append(c.Credentials, credential)
	c.CurrentCredentialName = credential.Name

	v.Set("current-credential-name", c.CurrentCredentialName)
	v.Set("credentials", c.Credentials)

	if err := v.WriteConfig(); err != nil {
		log.Fatalf("failed to write updated config file: %v", err)
	}

	log.Infof("Added credential '%s' to config file at %s", credential.Name, configPath)
	return nil
}

func UpdateCredentialsInConfigFile(updatedCredential Credential, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at %s", configPath)
	} else if err != nil {
		log.Fatalf("error checking config file: %v", err)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var c HarborConfig
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	updated := false
	for i, cred := range c.Credentials {
		if cred.Name == updatedCredential.Name {
			c.Credentials[i] = updatedCredential
			updated = true
			break
		}
	}

	if !updated {
		log.Fatalf("credential with name '%s' not found", updatedCredential.Name)
	}

	v.Set("current-credential-name", c.CurrentCredentialName)
	v.Set("credentials", c.Credentials)

	if err := v.WriteConfig(); err != nil {
		log.Fatalf("failed to write updated config file: %v", err)
	}

	log.Infof("Updated credential '%s' in config file at %s", updatedCredential.Name, configPath)
	return nil
}
