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

// Transaction runs the given function and calls Commit/Rollback as needed
func Transaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
