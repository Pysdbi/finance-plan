package database

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	err = db.AutoMigrate(&Client{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&Client{},
		&Income{},
		&Expense{},
		&FinDate{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
