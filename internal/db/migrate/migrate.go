package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"nltimv.com/karma-chameleon/internal/log"
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
				if !d.Migrator().HasTable(&User{}) {
					return d.Migrator().CreateTable(&User{})
				}
				return nil
			},
			Rollback: func(d *gorm.DB) error {
				if d.Migrator().HasTable("users") {
					return d.Migrator().DropTable("users")
				}
				return nil
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
				if !d.Migrator().HasTable(&Group{}) {
					return d.Migrator().CreateTable(&Group{})
				}
				return nil
			},
			Rollback: func(d *gorm.DB) error {
				if d.Migrator().HasTable("groups") {
					return d.Migrator().DropTable("groups")
				}
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Error.Fatalf("Failed to migrate database: %v", err)
	} else {
		log.Default.Println("Database migration successful!")
	}
}
