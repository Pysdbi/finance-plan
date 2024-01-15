package database

import (
	"finance/internal/finance/findate"
	"gorm.io/gorm"
)

type FinDate struct {
	gorm.Model

	ID           int               `gorm:"primary_key"`
	StartYear    int               `gorm:"type:int;"`
	StartMonth   int               `gorm:"type:int;"`
	StartDay     int               `gorm:"type:int;"`
	EndYear      *int              `gorm:"type:int;"`
	EndMonth     *int              `gorm:"type:int;"`
	EndDay       *int              `gorm:"type:int;"`
	PeriodDays   int               `gorm:"type:int;"`
	PeriodMonths int               `gorm:"type:int;"`
	PeriodYears  int               `gorm:"type:int;"`
	Type         findate.EventType `gorm:"type:varchar(50);"`
}
