package income

import "finance/internal/finance/findate"

// Income Доход
type Income struct {
	FinDate findate.FinDate
	// Amount сумма поступления
	Amount int64

	// Name название поступления
	Name string
	// Description описание поступления
	Description string
}
