package migrations

import (
	"record-services/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Section{},
		&models.Employee{},
	)
	return err
}
