package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dbDsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
