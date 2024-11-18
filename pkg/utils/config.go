package utils

import (
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

var (
	CurrentHarborData    *HarborData
	CurrentConfig        *HarborConfig
	configMutex          sync.RWMutex
	configInitialization sync.Once
	configInitError      error
)

type HarborData struct {
	ConfigPath string `yaml:"configpath"`
}

var (
	HarborConfigDir   string
	HarborDataDir     string
	DefaultConfigPath string
	DefaultDataPath   string
)

func SetLocation() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user home directory: %v", err)
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home, ".config")
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(home, ".local", "share")
	}

	HarborConfigDir = filepath.Join(xdgConfigHome, "harbor-cli")
	HarborDataDir = filepath.Join(xdgDataHome, "harbor-cli")

	DefaultConfigPath = filepath.Join(HarborConfigDir, "config.yaml")
	DefaultDataPath = filepath.Join(HarborDataDir, "data.yaml")
}

// Check for ERROR management
func GetCurrentConfig() *HarborConfig {
	configInitialization.Do(func() {
		// No action needed; InitConfig should have been called before
	})

	if configInitError != nil {
		return nil
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	if CurrentConfig == nil {
		return nil
	}

	return CurrentConfig
}

func GetHarborData() *HarborData {
	return CurrentHarborData
}

func (hc *HarborConfig) GetCurrentCredentialName() string {
	return hc.CurrentCredentialName
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
	currentConfig := GetCurrentConfig()
	for _, cred := range currentConfig.Credentials {
		if cred.Name == credentialName {
			return cred, nil
		}
	}

	return Credential{}, fmt.Errorf("credential with name LOL '%s' not found", credentialName)
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

func InitConfig(cfgFile string, userSpecifiedConfig bool) {
	configInitialization.Do(func() {
		SetLocation()
		dataFilePath := DefaultDataPath

		// Ensure data directory exists
		if err := os.MkdirAll(HarborDataDir, os.ModePerm); err != nil {
			configInitError = fmt.Errorf("failed to create data directory: %w", err)
			log.Fatalf("failed to create data directory: %v", err)
		}

		var configPath string

		// Check if data file exists and read it
		dataFile, err := ReadDataFile(dataFilePath)
		if err == nil && dataFile.ConfigPath != "" {
			configPath = dataFile.ConfigPath
		} else {
			// Data file does not exist or failed to read
			if userSpecifiedConfig && cfgFile != "" {
				// Create data file with the specified config path
				if err := CreateDataFile(dataFilePath, cfgFile); err != nil {
					configInitError = fmt.Errorf("failed to create data file: %w", err)
					log.Fatalf("failed to create data file: %v", err)
				}
				configPath = cfgFile
			} else {
				// Create data file with the default config path
				if err := CreateDataFile(dataFilePath, DefaultConfigPath); err != nil {
					configInitError = fmt.Errorf("failed to create data file: %w", err)
					log.Fatalf("failed to create data file: %v", err)
				}
				configPath = DefaultConfigPath
			}
		}

		// If user specified --config, override and update data file
		if userSpecifiedConfig && cfgFile != "" {
			configPath, err = filepath.Abs(cfgFile)
			if err != nil {
				configInitError = fmt.Errorf("failed to resolve absolute path for config file: %w", err)
				log.Fatalf("failed to resolve absolute path for config file: %v", err)
			}
			// Update the data file with the new config path
			if err := UpdateDataFile(dataFilePath, configPath); err != nil {
				configInitError = fmt.Errorf("failed to update data file: %w", err)
				log.Fatalf("failed to update data file: %v", err)
			}
		}

		// If configPath is still not set, use the default config path
		if configPath == "" {
			configPath = DefaultConfigPath
		}
		// Initialize Viper for the main config
		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")
		// Handle default config path
		if configPath == DefaultConfigPath {
			stat, err := os.Stat(configPath)
			if err == nil && stat.Size() == 0 {
				log.Println("Config file is empty, creating a new one")
			}

			if os.IsNotExist(err) {
				log.Infof("Config file not found at %s, creating a new one", configPath)
			}

			if os.IsNotExist(err) || (err == nil && stat.Size() == 0) {
				// Ensure the config directory exists
				if _, err := os.Stat(HarborConfigDir); os.IsNotExist(err) {
					log.Println("Creating config directory:", HarborConfigDir)
					if err := os.MkdirAll(HarborConfigDir, os.ModePerm); err != nil {
						configInitError = fmt.Errorf("failed to create config directory: %w", err)
						return
					}
				}

				// Create the config file
				if err := CreateConfigFile(configPath); err != nil {
					configInitError = fmt.Errorf("failed to create config file: %w", err)
					return
				}

				log.Infof("Config file created at %s", configPath)
			}
		}

		// Read in the main config file
		if err := v.ReadInConfig(); err != nil {
			configInitError = fmt.Errorf("error reading config file: %w. Please ensure the config file exists.", err)
			log.Fatalf("error reading config file: %v. Please ensure the config file exists.", err)
		}
		// Unmarshal into HarborConfig struct
		var harborConfig HarborConfig
		if err := v.Unmarshal(&harborConfig); err != nil {
			configInitError = fmt.Errorf("failed to unmarshal config file: %w", err)
			log.Fatalf("failed to unmarshal config file: %v", err)
		}
		// Assign to global variable with thread safety
		configMutex.Lock()
		CurrentConfig = &harborConfig
		configMutex.Unlock()

		CurrentHarborData = &HarborData{
			ConfigPath: configPath,
		}

		log.Infof("Using config file: %s", v.ConfigFileUsed())
	})
}
