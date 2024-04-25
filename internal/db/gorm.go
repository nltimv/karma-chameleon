package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db    *gorm.DB
	sqldb *sql.DB
)

func Open(dbHost string, dbPort string, dbUser string, dbPassword string, dbName string, sslmode string) {
	var err error
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbHost, dbUser, dbPassword, dbName, dbPort, sslmode)

	ctr := 0
	for {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}

		if ctr == 5 {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		log.Println("Failed to connect to database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		ctr++
	}

	sqldb, _ = db.DB()

	if err = sqldb.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

}

func Close() error {
	return sqldb.Close()
}
