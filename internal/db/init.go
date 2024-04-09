package db

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func CreateTables() {
	fmt.Println("Migrating database...")
	io, err := iofs.New(fs, "migrations")
	handleError(err)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	handleError(err)
	m, err := migrate.NewWithInstance("iofs", io, "postgres", driver)
	handleError(err)
	err = m.Up()
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
