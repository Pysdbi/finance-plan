package findate

import (
	"time"
)

type DateRangeItem struct {
	Time time.Time
	Date string
}

const format = "02.01.2006"

func DefineRange(from, to time.Time) []DateRangeItem {
	rangeItems := make([]DateRangeItem, 0)

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		rangeItems = append(rangeItems, DateRangeItem{
			Time: d,
			Date: d.Format(format),
		})
	}
	return rangeItems
}
