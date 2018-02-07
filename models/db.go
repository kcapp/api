package models

import (
	"database/sql"
	"log"
	// Blank import used to be able to register DB driver
	_ "github.com/go-sql-driver/mysql"
)

// DB point to our database
var DB *sql.DB

// InitDB will initialize the database with the given datasource
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}
