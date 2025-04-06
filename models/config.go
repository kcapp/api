package models

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetMysqlConnectionString returns mysql connection string
func GetMysqlConnectionString() string {
	// Need to add ?parseTime=true here to support time.Time in queries
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?parseTime=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.address"),
		viper.GetInt("db.port"),
		viper.GetString("db.schema"))
}
