package balance

import "finance/internal/finance/findate"

type TransactionIncome struct {
	FinDate findate.FinDate
	Amount  int64
}

type TransactionExpense struct {
	FinDate findate.FinDate
	Amount  int64
}
