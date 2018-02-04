package models

import (
	"database/sql"
	"log"
	// Blank import used to be able to register DB driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB will initialize the database with the given datasource
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
