package db

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	sqldb *sql.DB
)

func Open(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string, sslmode string) {
	var err error
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbHost, dbUser, dbPassword, dbName, dbPort, sslmode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqldb, _ = db.DB()

	if err = sqldb.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

}

func Close() error {
	return sqldb.Close()
}
