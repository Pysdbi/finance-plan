package database

import (
	"gorm.io/gorm"
)

// Expense Расход
type Expense struct {
	gorm.Model

	ID          int     `gorm:"primary_key"`
	ClientID    int     `gorm:"index;"`
	FinDateID   int     `gorm:"index;"` // Внешний ключ для связи с FinDate
	FinDate     FinDate `gorm:"foreignKey:FinDateID"`
	Amount      int64   `gorm:"type:bigint;"`
	Name        string  `gorm:"type:varchar(100);"`
	Description string  `gorm:"type:text;"`
}
