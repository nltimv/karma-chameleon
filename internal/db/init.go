package db

import (
	"nltimv.com/karma-chameleon/internal/db/migrate"
	"nltimv.com/karma-chameleon/internal/log"
)

func Migrate() {
	if db == nil {
		log.Error.Fatal("Database has not been initialized yet!")
	}

	migrate.Migrate(db)
}
