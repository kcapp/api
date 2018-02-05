package models

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// DBConfig stuct config
type DBConfig struct {
	Address  string `yaml:"address"`
	Dbport   string `yaml:"dbport"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

// APIConfig struct config
type APIConfig struct {
	Port string `yaml:"port"`
}

// Config type
type Config struct {
	DBConfig  DBConfig  `yaml:"db"`
	APIConfig APIConfig `yaml:"api"`
}

// GetConfig loads configuration from yaml file
func GetConfig() (*Config, error) {
	config := new(Config)
	yamlFile, err := ioutil.ReadFile("config/configdb.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetMysqlConnectionString returns mysql connection string
func (config *Config) GetMysqlConnectionString() (string, error) {
	return fmt.Sprintf(
			"%s:%s@(%s:%s)/%s",
			config.DBConfig.Username,
			config.DBConfig.Password,
			config.DBConfig.Address,
			config.DBConfig.Dbport,
			config.DBConfig.Schema),
		nil
}

// GetAPIPort returns API configuration
func (config *Config) GetAPIPort() (string, error) {
	return fmt.Sprintf("%s", config.APIConfig.Port), nil
}
