package expense

import "finance/internal/finance/findate"

// Expense Расход
type Expense struct {
	FinDate     findate.FinDate
	Amount      int64
	Name        string
	Description string
}
