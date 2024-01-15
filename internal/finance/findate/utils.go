package findate

import (
	"fmt"
	"time"
)

const formatDate = "02.01.2006"

// IsDate проверяет, соответствует ли строка формату dd.mm.yyyy
func IsDate(v string) bool {
	_, err := time.Parse(formatDate, v)
	return err == nil
}

// ParseDate анализирует строку формата "dd.mm.yyyy" и возвращает день, месяц и год
func ParseDate(v string) (day, month, year int) {
	t, err := time.Parse(formatDate, v)
	if err != nil {
		return 0, 0, 0
	}
	return t.Day(), int(t.Month()), t.Year()
}

// ParseDatePtr анализирует строку формата "dd.mm.yyyy" и возвращает день, месяц и год
func ParseDatePtr(v string) (day, month, year *int) {
	t, err := time.Parse(formatDate, v)
	if err != nil {
		return nil, nil, nil
	}
	d, m, y := t.Day(), int(t.Month()), t.Year()
	return &d, &m, &y
}

const formatPeriod = "%d:%d:%d"

// IsPeriod проверяет, соответствует ли строка формату "%d:%d:%d",
// где %d представляет целочисленное значение.
func IsPeriod(v string) bool {
	var days, months, years int
	_, err := fmt.Sscanf(v, formatPeriod, &days, &months, &years)
	return err == nil
}

// ParsePeriod анализирует строку формата "%d:%d:%d" и возвращает дни, месяцы и годы.
// Возвращает 0 для каждого числа, если строка не соответствует формату.
func ParsePeriod(v string) (days, months, years int) {
	_, err := fmt.Sscanf(v, formatPeriod, &days, &months, &years)
	if err != nil {
		return 0, 0, 0
	}
	return days, months, years
}
