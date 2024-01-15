package findate_b

import (
	"finance/internal/finance/findate"
)

type finDateBuilder struct {
	startYear    int
	startMonth   int
	startDay     int
	endYear      *int
	endMonth     *int
	endDay       *int
	periodDays   int
	periodMonths int
	periodYears  int
	_type        findate.EventType
}

func NewFinDateBuilder() IFinDateBuilder {
	return &finDateBuilder{}
}

type IFinDateBuilder interface {
	StartAt(day, month, year int) IFinDateBuilder
	EndAt(day, month, year int) IFinDateBuilder // only for Periodic

	IFinDateBuilderPeriods

	Once() IFinDateBuilder
	Periodic() IFinDateBuilder

	Build() *findate.FinDate
}

type IFinDateBuilderPeriods interface {
	SetPeriod(days, months, years int) IFinDateBuilder

	EveryDay() IFinDateBuilder
	Every14Days() IFinDateBuilder

	EveryMonth() IFinDateBuilder
	Every6Months() IFinDateBuilder

	EveryYear() IFinDateBuilder
}

func (b *finDateBuilder) StartAt(day, month, year int) IFinDateBuilder {
	b.startDay, b.startMonth, b.startYear = day, month, year
	return b
}

func (b *finDateBuilder) EndAt(day, month, year int) IFinDateBuilder {
	b.endDay, b.endMonth, b.endYear = &day, &month, &year
	return b
}

func (b *finDateBuilder) Once() IFinDateBuilder {
	b._type = findate.Once
	return b
}

func (b *finDateBuilder) Periodic() IFinDateBuilder {
	b._type = findate.Periodic
	return b
}

// ========= Set Periods =========

func (b *finDateBuilder) SetPeriod(days, months, years int) IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = days, months, years
	b.Periodic()
	return b
}

func (b *finDateBuilder) EveryDay() IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = 1, 0, 0
	b.Periodic()
	return b
}

func (b *finDateBuilder) Every14Days() IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = 14, 0, 0
	b.Periodic()
	return b
}

func (b *finDateBuilder) EveryMonth() IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = 0, 1, 0
	b.Periodic()
	return b
}

func (b *finDateBuilder) Every6Months() IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = 1, 0, 0
	b.Periodic()
	return b
}

func (b *finDateBuilder) EveryYear() IFinDateBuilder {
	b.periodDays, b.periodMonths, b.periodYears = 0, 0, 1
	b.Periodic()
	return b
}

// ========= END Set Periods =========

func (b *finDateBuilder) Build() *findate.FinDate {
	return &findate.FinDate{
		StartYear:    b.startYear,
		StartMonth:   b.startMonth,
		StartDay:     b.startDay,
		EndYear:      b.endYear,
		EndMonth:     b.endMonth,
		EndDay:       b.endDay,
		PeriodDays:   b.periodDays,
		PeriodMonths: b.periodMonths,
		PeriodYears:  b.periodYears,
		Type:         b._type,
	}
}
