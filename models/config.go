package models

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// DBConfig stuct config
type DBConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

// APIConfig struct config
type APIConfig struct {
	Port int `yaml:"port"`
}

// Config type
type Config struct {
	DBConfig  DBConfig  `yaml:"db"`
	APIConfig APIConfig `yaml:"api"`
}

// GetConfig loads configuration from yaml file
func GetConfig(configFileParam string) (*Config, error) {
	// Default location
	configFilePath := "config/config.yaml"
	if len(configFileParam) > 0 {
		configFilePath = configFileParam
	}
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetMysqlConnectionString returns mysql connection string
func (config *Config) GetMysqlConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s",
		config.DBConfig.Username,
		config.DBConfig.Password,
		config.DBConfig.Address,
		config.DBConfig.Port,
		config.DBConfig.Schema)
}
