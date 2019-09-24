package config

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	enc "github.com/scrambledeggs/booky-go-common/encryption"
	"github.com/spf13/viper"
)

// Environment variables
const (
	ENV_DEV     string = "development"
	ENV_STAGING string = "staging"
	ENV_TEST    string = "test"
	ENV_PROD    string = "production"
)

type Config interface {
	SetConfig(configFile string, configFilePath string, configFileType string)
	SetCipherPass(passphrase string)
	ApplyEnvConfig() error
	GetEnv() string
}

type config struct {
	environment    *string
	configMap      *map[string]string
	configFile     *string
	configFilePath *string
	configFileType *string
	cipherpass     *string
}

func New(appEnv string, configMap map[string]string) Config {
	// Default config values
	cf := "config"
	cfp := "."
	cft := "yml"

	conf := config{
		environment:    &appEnv,
		configMap:      &configMap,
		configFile:     &cf,
		configFilePath: &cfp,
		configFileType: &cft,
		cipherpass:     &cfp,
	}

	conf.checkEnvValidity()

	return &conf
}

func (c *config) SetConfig(configFile string, configFilePath string, configFileType string) {
	c.configFile = &configFile
	c.configFilePath = &configFilePath
	c.configFileType = &configFileType
}

func (c *config) SetCipherPass(passphrase string) {
	c.cipherpass = &passphrase
}

func (c *config) ApplyEnvConfig() error {
	if len(*c.configMap) == 0 {
		return errors.New("Config Map is not set.")
	}

	err := setViperConfig(*c.configFile, *c.configFilePath, *c.configFileType)
	if err != nil {
		return err
	}

	// Set environment variables
	prefix := *c.environment + "."
	for k, v := range *c.configMap {
		err := c.setKeyVarEnv(k, viper.GetString(prefix+v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *config) GetEnv() string {
	return *c.environment
}

func (c *config) setKeyVarEnv(key string, val string) error {
	finalVal := val

	// Handle encrypted values
	isEncrypted, _ := regexp.MatchString(`^ENC(.*)$`, val)
	if isEncrypted {
		trimmed := strings.TrimSuffix(strings.TrimPrefix(val, "ENC("), ")")
		temp, err := enc.DecryptB64(trimmed, *c.cipherpass)
		if err != nil {
			return err
		}

		finalVal = string(temp)
	}

	os.Setenv(key, finalVal)
	return nil
}

// Checks passed environment variable
func (c *config) checkEnvValidity() bool {
	validEnvs := map[int]string{
		1: ENV_DEV,
		2: ENV_TEST,
		3: ENV_STAGING,
		4: ENV_PROD,
	}

	for _, v := range validEnvs {
		if *c.environment == v {
			return true
		}
	}

	log.Print("Could not detect valid environment. Setting Config ENVIRONMENT to 'development'.")
	defaultEnv := ENV_DEV
	c.environment = &defaultEnv
	return false
}

// Reads config values and stores in viper
func setViperConfig(configFile string, configPath string, configType string) error {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetConfigType(configType)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
