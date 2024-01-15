package findate

import (
	"time"
)

type EventType int

const (
	Once EventType = iota
	Periodic
)

type FinDate struct {
	StartYear    int
	StartMonth   int
	StartDay     int
	EndYear      *int
	EndMonth     *int
	EndDay       *int
	PeriodDays   int
	PeriodMonths int
	PeriodYears  int
	Type         EventType
}

func (fd *FinDate) InDate(d time.Time) bool {
	startDate := time.Date(fd.StartYear, time.Month(fd.StartMonth), fd.StartDay, 0, 0, 0, 0, time.UTC)
	var endDate time.Time
	if fd.EndYear != nil && fd.EndMonth != nil && fd.EndDay != nil {
		endDate = time.Date(*fd.EndYear, time.Month(*fd.EndMonth), *fd.EndDay, 0, 0, 0, 0, time.UTC)
	}

	if fd.Type == Once && d.Equal(startDate) {
		return true
	} else if fd.Type == Periodic {
		periodMatch := false

		// Проверяем дневную периодичность
		if fd.PeriodDays > 0 {
			daysSinceStart := int(d.Sub(startDate).Hours() / 24)
			periodMatch = daysSinceStart >= 0 && daysSinceStart%fd.PeriodDays == 0
		}

		// Проверяем месячную периодичность
		if !periodMatch && fd.PeriodMonths > 0 {
			monthsSinceStart := (d.Year()-startDate.Year())*12 + int(d.Month()) - int(startDate.Month())

			if monthsSinceStart >= 0 && monthsSinceStart%fd.PeriodMonths == 0 {
				// Подсчитываем предполагаемую дату следующего события
				expectedMonth := startDate.Month() + time.Month((monthsSinceStart/fd.PeriodMonths)*fd.PeriodMonths)
				expectedYear := startDate.Year()
				for expectedMonth > 12 {
					expectedMonth -= 12
					expectedYear++
				}

				// Определяем количество дней в ожидаемом месяце
				daysInMonth := time.Date(expectedYear, expectedMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()

				// Корректируем день, если он выходит за пределы месяца
				expectedDay := startDate.Day()
				if expectedDay > daysInMonth {
					expectedDay = daysInMonth
				}

				// Создаем ожидаемую дату
				expectedDate := time.Date(expectedYear, expectedMonth, expectedDay, 0, 0, 0, 0, time.UTC)

				// Проверяем соответствует ли текущая дата ожидаемой дате события
				if d.Equal(expectedDate) {
					periodMatch = true
				}
			}
		}

		// Проверяем годовую периодичность
		if !periodMatch && fd.PeriodYears > 0 {
			yearsSinceStart := d.Year() - startDate.Year()
			periodMatch = d.Day() == startDate.Day() && d.Month() == startDate.Month() && yearsSinceStart >= 0 && yearsSinceStart%fd.PeriodYears == 0
		}

		if periodMatch && (fd.EndYear == nil || d.Before(endDate) || d.Equal(endDate)) {
			return true
		}
		return periodMatch
	}

	return false
}
