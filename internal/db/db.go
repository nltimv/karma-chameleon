package db

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func Open(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string) error {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	return err
}

func Close() error {
	return db.Close()
}
