package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetConfig will test reading a configuration file from disk
func TestGetConfig(t *testing.T) {
	config, err := GetConfig("../config/config.yaml")
	assert.Equal(t, err, nil, "err should be nil")
	assert.Equal(t, config.APIConfig.Port, 8001, "port should be 8001")
	assert.Equal(t, config.DBConfig.Address, "localhost", "address should be localhost")
	assert.Equal(t, config.DBConfig.Port, 3306, "port should be 3306")
	assert.Equal(t, config.DBConfig.Username, "developer", "username should be developer")
	assert.Equal(t, config.DBConfig.Password, "abcd1234", "password should be abcd1234")
	assert.Equal(t, config.DBConfig.Schema, "kcapp", "schema should be kcapp")
}

// TestGetMysqlConnectionString will check that we create a correct MySQL connection string
func TestGetMysqlConnectionString(t *testing.T) {
	config := new(Config)
	config.DBConfig.Address = "localhost"
	config.DBConfig.Port = 3306
	config.DBConfig.Username = "developer"
	config.DBConfig.Password = "abcd1234"
	config.DBConfig.Schema = "kcapp"

	conn := config.GetMysqlConnectionString()
	assert.Equal(t, conn, "developer:abcd1234@(localhost:3306)/kcapp?parseTime=true")
}
