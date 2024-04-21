package migrate

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0001",
			Migrate: func(d *gorm.DB) error {
				type User struct {
					ID     uint `gorm:"primaryKey"`
					UserId string
					TeamId string
					Karma  int
				}

				return d.Migrator().CreateTable(&User{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("users")
			},
		},
		{
			ID: "0002",
			Migrate: func(d *gorm.DB) error {
				type Group struct {
					ID      uint `gorm:"primaryKey"`
					GroupId string
					TeamId  string
					Karma   int
				}

				return d.Migrator().CreateTable(&Group{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("groups")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
