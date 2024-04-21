package db

import (
	"log"

	"nltimv.com/karma-chameleon/internal/db/migrate"
)

func Migrate() {
	if db == nil {
		log.Fatal("Database has not been initialized yet!")
	}

	migrate.Migrate(db)
}
