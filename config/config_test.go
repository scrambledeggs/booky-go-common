package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDev(t *testing.T) {
	var err error

	// mock data
	validAppEnv := "development"
	validConfigName := "config_sample"

	var configMap = map[string]string{
		"DB_TABLE": "db.table",
		"APP_KEY":  "app.key",
	}

	conf := New(validAppEnv, configMap)
	conf.SetConfig(validConfigName, ".", "yml")
	err = conf.ApplyEnvConfig()

	assert.Nil(t, err)
	assert.Equal(t, "development", conf.GetEnv())
	assert.Equal(t, "dev_table", os.Getenv("DB_TABLE"))
}

func TestConfigStaging(t *testing.T) {
	var err error

	// mock data
	validAppEnv := "staging"
	validConfigName := "config_sample"
	passphrase := "some-secure-passphrase"
	var configMap = map[string]string{
		"DB_TABLE": "db.table",
		"APP_KEY":  "app.key",
	}

	// Valid Configuration - development
	conf := New(validAppEnv, configMap)
	conf.SetConfig(validConfigName, ".", "yml")
	conf.SetCipherPass(passphrase)
	err = conf.ApplyEnvConfig()

	assert.Nil(t, err)
	assert.Equal(t, "staging_table", os.Getenv("DB_TABLE"))

	// Ensuring that encrypted word was successfully decrypted
	assert.Equal(t, "super-secret-key", os.Getenv("APP_KEY"))

	// Invalid passphrase
	conf.SetCipherPass("random-pass")
	err = conf.ApplyEnvConfig()
	assert.NotNil(t, err)
}

func TestConfigOtherInvalids(t *testing.T) {
	var err error

	// mock data
	validAppEnv := "staging"
	validConfigName := "config_sample"
	var configMap = map[string]string{
		"DB_TABLE": "db.table",
		"APP_KEY":  "app.key",
	}

	// Invalid Environment
	conf := New("xenv", configMap)
	assert.Equal(t, "development", conf.GetEnv())

	// Invalid map-config
	conf = New(validAppEnv, nil)
	err = conf.ApplyEnvConfig()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Config Map is not set.")

	// Invalid config properties
	conf = New(validAppEnv, configMap)
	conf.SetConfig("configx", ".", "yml")
	err = conf.ApplyEnvConfig()
	assert.Contains(t, err.Error(), "Config File \"configx\" Not Found")

	conf.SetConfig(validConfigName, ".", "csv")
	err = conf.ApplyEnvConfig()
	assert.Contains(t, err.Error(), "Unsupported Config Type")
}
