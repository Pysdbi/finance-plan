package database

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	ID       int       `gorm:"primary_key"`
	Name     string    `gorm:"type:varchar(100);"`
	Incomes  []Income  `gorm:"foreignKey:ClientID"`
	Expenses []Expense `gorm:"foreignKey:ClientID"`
}
