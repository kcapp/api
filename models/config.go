package models

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config type
type Config struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
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
	connectionString := config.Username +
		":" +
		config.Password +
		"@(" +
		config.Address +
		":" +
		config.Port +
		")/" +
		config.Schema

	return connectionString, nil
}
